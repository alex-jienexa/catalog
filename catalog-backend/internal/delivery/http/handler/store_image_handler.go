package handler

import (
	"catalog-backend/internal/usecase"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StoreImageHandler struct {
	storeUC    *usecase.StoreUseCase
	uploadPath string
	maxSize    int64
}

func NewStoreImageHandler(storeUC *usecase.StoreUseCase, uploadPath string, maxSize int64) *StoreImageHandler {
	return &StoreImageHandler{
		storeUC:    storeUC,
		uploadPath: uploadPath,
		maxSize:    maxSize,
	}
}

func (h *StoreImageHandler) UploadStoreImage(c *gin.Context) {
	// Получаем текущую информацию о магазине
	currentStore, err := h.storeUC.GetStoreInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get store info"})
		return
	}

	// Получаем файл
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}
	defer file.Close()

	// Проверка размера
	if header.Size > h.maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File too large. Max size: %d bytes", h.maxSize)})
		return
	}

	// Проверка типа (MIME)
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	filetype := http.DetectContentType(buffer)
	if !strings.HasPrefix(filetype, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is not an image"})
		return
	}
	file.Seek(0, io.SeekStart)

	// Создаём папку для изображений магазина
	storeUploadDir := filepath.Join(h.uploadPath, "store")
	if err := os.MkdirAll(storeUploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Генерируем имя файла
	ext := filepath.Ext(header.Filename)
	newFilename := uuid.New().String() + ext
	relativePath := "/uploads/store/" + newFilename
	fullPath := filepath.Join(storeUploadDir, newFilename)

	// Сохраняем файл
	out, err := os.Create(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Удаляем старое изображение, если оно локальное
	if currentStore.ImageURL != "" && strings.HasPrefix(currentStore.ImageURL, "/uploads/") {
		oldPath := filepath.Join(h.uploadPath, strings.TrimPrefix(currentStore.ImageURL, "/uploads"))
		os.Remove(oldPath)
	}

	// Обновляем информацию о магазине
	updatedStore := *currentStore
	updatedStore.ImageURL = relativePath

	_, err = h.storeUC.UpdateStoreInfo(c.Request.Context(), &updatedStore)
	if err != nil {
		// Если ошибка, удаляем загруженный файл
		os.Remove(fullPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedStore)
}
