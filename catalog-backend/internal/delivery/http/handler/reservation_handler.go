package handler

import (
	"catalog-backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	reservationUC *usecase.ReservationUseCase
}

func NewReservationHandler(reservationUC *usecase.ReservationUseCase) *ReservationHandler {
	return &ReservationHandler{reservationUC: reservationUC}
}

type CreateReservationRequest struct {
	CustomerID int `json:"customer_id"`
	ProductID  int `json:"product_id"`
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var req CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reservation, err := h.reservationUC.CreateReservation(c.Request.Context(), req.CustomerID, req.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, reservation)
}

func (h *ReservationHandler) GetReservations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	reservations, total, err := h.reservationUC.GetReservations(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"reservations": reservations,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

func (h *ReservationHandler) UpdateReservationStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reservation, err := h.reservationUC.UpdateReservationStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reservation)
}
