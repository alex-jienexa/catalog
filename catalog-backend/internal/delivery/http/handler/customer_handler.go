package handler

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerUC *usecase.CustomerUseCase
}

func NewCustomerHandler(customerUC *usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{customerUC: customerUC}
}

// POST /api/v1/public/customers
func (h *CustomerHandler) GetOrCreate(c *gin.Context) {
	var input entity.CustomerCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer, existed, err := h.customerUC.GetOrCreateCustomer(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"customer": customer,
		"existed":  existed,
	})
}
