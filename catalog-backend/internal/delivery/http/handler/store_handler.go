package handler

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeUC *usecase.StoreUseCase
}

func NewStoreHandler(storeUC *usecase.StoreUseCase) *StoreHandler {
	return &StoreHandler{storeUC: storeUC}
}

func (h *StoreHandler) GetStoreInfo(c *gin.Context) {
	info, err := h.storeUC.GetStoreInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *StoreHandler) UpdateStoreInfo(c *gin.Context) {
	var info entity.StoreInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.storeUC.UpdateStoreInfo(c.Request.Context(), &info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}
