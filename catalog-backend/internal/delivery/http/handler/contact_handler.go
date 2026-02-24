package handler

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/ZeRg0912/logger"
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
	logger.Info("CreateContact: получен запрос")

	var input entity.ContactCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("CreateContact: ошибка при парсинге JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("CreateContact: данные запроса: %+v", input)

	contact, err := h.contactUC.CreateContact(c.Request.Context(), input)
	if err != nil {
		logger.Error("CreateContact: ошибка создания товара: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("CreateContact: контакт успешно создан, ID=%d", contact.ID)
	c.JSON(http.StatusCreated, contact)
}

func (h *ContactHandler) GetContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("GetContact: неверный ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	logger.Debug("GetContact: запрос контакта с ID=%d", id)

	contact, err := h.contactUC.GetContact(c.Request.Context(), id)
	if err != nil {
		logger.Error("GetContact: ошибка получения контакта ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if contact == nil {
		logger.Warn("GetContact: контакт с ID=%d не найден", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	logger.Debug("GetContact: получен контакт: %+v", contact)
	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) GetContacts(c *gin.Context) {
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active_only", "true"))

	contacts, err := h.contactUC.GetContacts(c.Request.Context(), activeOnly)
	if err != nil {
		logger.Error("GetContacts: ошибка получения контактов: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("GetContacts: получено %d контактов: %+v", len(contacts), contacts)
	c.JSON(http.StatusOK, contacts)
}

func (h *ContactHandler) UpdateContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("UpdateContact: неверный ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	var input entity.ContactUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("UpdateContact: ошибка при парсинге JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.contactUC.UpdateContact(c.Request.Context(), id, input)
	if err != nil {
		logger.Error("UpdateContact: ошибка обновлении контакта ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("UpdateContact: контакт ID=%d успешно изменён", contact.ID)
	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) DeleteContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("DeleteContact: неверный ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	err = h.contactUC.DeleteContact(c.Request.Context(), id)
	if err != nil {
		logger.Error("DeleteContact: ошибка удаления контакта ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("DeleteContact: контакт ID=%d успешно удалён", id)
	c.Status(http.StatusNoContent)
}
