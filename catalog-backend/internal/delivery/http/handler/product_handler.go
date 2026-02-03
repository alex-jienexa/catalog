package handler

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"catalog-backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUC *usecase.ProductUseCase
}

func NewProductHandler(productUC *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUC: productUC,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input entity.ProductCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUC.CreateProduct(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productUC.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sectionID, _ := strconv.Atoi(c.Query("section_id"))
	isActive := c.Query("is_active")

	filters := repository.ProductFilters{
		Limit:  limit,
		Offset: (page - 1) * limit,
	}

	if sectionID > 0 {
		filters.SectionID = &sectionID
	}

	if isActive != "" {
		active, err := strconv.ParseBool(isActive)
		if err == nil {
			filters.IsActive = &active
		}
	}

	products, total, err := h.productUC.GetProducts(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"pages":    (total + limit - 1) / limit,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var input entity.ProductUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUC.UpdateProduct(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = h.productUC.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
