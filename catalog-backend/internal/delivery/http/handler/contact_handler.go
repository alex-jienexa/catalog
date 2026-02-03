package handler

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	contactUC *usecase.ContactUseCase
}

func NewContactHandler(contactUC *usecase.ContactUseCase) *ContactHandler {
	return &ContactHandler{
		contactUC: contactUC,
	}
}

func (h *ContactHandler) CreateContact(c *gin.Context) {
	var input entity.ContactCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.contactUC.CreateContact(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

func (h *ContactHandler) GetContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	contact, err := h.contactUC.GetContact(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if contact == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) GetContacts(c *gin.Context) {
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active_only", "true"))

	contacts, err := h.contactUC.GetContacts(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

func (h *ContactHandler) UpdateContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	var input entity.ContactUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.contactUC.UpdateContact(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) DeleteContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	err = h.contactUC.DeleteContact(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
