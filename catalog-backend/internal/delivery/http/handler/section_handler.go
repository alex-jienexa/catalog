package handler

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/usecase"
	"net/http"
	"strconv"

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
	var input entity.SectionCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section, err := h.sectionUC.CreateSection(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, section)
}

func (h *SectionHandler) GetSection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	section, err := h.sectionUC.GetSection(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if section == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	c.JSON(http.StatusOK, section)
}

func (h *SectionHandler) GetSections(c *gin.Context) {
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active_only", "true"))

	sections, err := h.sectionUC.GetSections(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sections)
}

func (h *SectionHandler) UpdateSection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	var input entity.SectionUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section, err := h.sectionUC.UpdateSection(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, section)
}

func (h *SectionHandler) DeleteSection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	err = h.sectionUC.DeleteSection(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
