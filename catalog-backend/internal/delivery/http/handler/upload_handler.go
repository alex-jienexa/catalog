package handler

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/usecase"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	productUC  *usecase.ProductUseCase
	uploadPath string
	maxSize    int64
}

func NewUploadHandler(productUC *usecase.ProductUseCase, uploadPath string, maxSize int64) *UploadHandler {
	return &UploadHandler{
		productUC:  productUC,
		uploadPath: uploadPath,
		maxSize:    maxSize,
	}
}

func (h *UploadHandler) UploadProductImage(c *gin.Context) {
	// Получаем ID товара
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Проверяем существование товара
	product, err := h.productUC.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Получаем файл из формы
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}
	defer file.Close()

	// Проверка размера
	if header.Size > h.maxSize {
		fmt.Printf("FileSize: %d", header.Size)
		fmt.Printf("MaxFileSize: %d", h.maxSize)
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
	// Возвращаем указатель в начало файла
	file.Seek(0, io.SeekStart)

	// Генерируем уникальное имя файла
	ext := filepath.Ext(header.Filename)
	newFilename := uuid.New().String() + ext
	relativePath := "/uploads/products/" + newFilename
	fullPath := filepath.Join(h.uploadPath, "products", newFilename)

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

	// Обновляем товар в БД
	updateData := entity.ProductUpdate{
		ImageURL: &relativePath,
	}
	updatedProduct, err := h.productUC.UpdateProduct(c.Request.Context(), id, updateData)
	if err != nil {
		// Если ошибка, удаляем загруженный файл
		os.Remove(fullPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}
