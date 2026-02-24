package handler

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/ZeRg0912/logger"
	"github.com/gin-gonic/gin"
)

type SectionHandler struct {
	sectionUC *usecase.SectionUseCase
}

func NewSectionHandler(sectionUC *usecase.SectionUseCase) *SectionHandler {
	return &SectionHandler{
		sectionUC: sectionUC,
	}
}

func (h *SectionHandler) CreateSection(c *gin.Context) {
	logger.Info("CreateSection: получен запрос")
	var input entity.SectionCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("CreateSection: ошибка при парсинге JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("CreateSection: данные запроса: %+v", input)

	section, err := h.sectionUC.CreateSection(c.Request.Context(), input)
	if err != nil {
		logger.Error("CreateSection: ошибка создания секции: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("CreateSection: секция успешно создана, ID=%d", section.ID)
	c.JSON(http.StatusCreated, section)
}

func (h *SectionHandler) GetSection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("GetSection: неверный ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	logger.Debug("GetSection: запрос секции с ID=%d", id)

	section, err := h.sectionUC.GetSection(c.Request.Context(), id)
	if err != nil {
		logger.Error("GetSection: ошибка получения секции ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if section == nil {
		logger.Warn("GetSection: секция с ID=%d не найдена", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	logger.Debug("GetSection: получена секция: %+v", section)
	c.JSON(http.StatusOK, section)
}

func (h *SectionHandler) GetSections(c *gin.Context) {
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active_only", "true"))

	sections, err := h.sectionUC.GetSections(c.Request.Context(), activeOnly)
	if err != nil {
		logger.Error("GetSections: ошибка получения секций: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("GetSections: получено %d секций: %+v", len(sections), sections)
	c.JSON(http.StatusOK, sections)
}

func (h *SectionHandler) UpdateSection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("UpdateSection: неверный ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	var input entity.SectionUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("UpdateSection: ошибка при парсинге JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section, err := h.sectionUC.UpdateSection(c.Request.Context(), id, input)
	if err != nil {
		logger.Error("UpdateSection: ошибка обновлении секции ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("UpdateSection: секция ID=%d успешно изменена", section.ID)
	c.JSON(http.StatusOK, section)
}

func (h *SectionHandler) DeleteSection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("DeleteSection: неверный ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	err = h.sectionUC.DeleteSection(c.Request.Context(), id)
	if err != nil {
		logger.Error("DeleteSection: ошибка удаления секции ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("DeleteSection: секция ID=%d успешно удалена", id)
	c.Status(http.StatusNoContent)
}
