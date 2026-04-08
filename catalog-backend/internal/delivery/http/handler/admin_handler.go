package handler

import (
	"catalog-backend/internal/delivery/http/middleware"
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/usecase"
	"errors"
	"net/http"
	"strconv"

	"github.com/ZeRg0912/logger"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUC *usecase.AdminUseCase
}

func NewAdminHandler(adminUC *usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUC: adminUC}
}

// IsFirstAdmin — GET /api/v1/auth/is-first
// Фронтенд проверяет, нужно ли показывать форму регистрации
func (h *AdminHandler) IsFirstAdmin(c *gin.Context) {
	isFirst, err := h.adminUC.IsFirstAdmin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_first": isFirst})
}

// Register — POST /api/v1/auth/register
// Только когда администраторов ещё нет
func (h *AdminHandler) Register(c *gin.Context) {
	var input entity.AdminCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, token, err := h.adminUC.Register(c.Request.Context(), input)
	if err != nil {
		logger.Error("Register: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Register: первый администратор '%s' создан", admin.Username)
	c.JSON(http.StatusCreated, gin.H{"admin": admin, "token": token})
}

// Login — POST /api/v1/auth/login
func (h *AdminHandler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, token, err := h.adminUC.Login(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Login: администратор '%s' вошёл в систему", admin.Username)
	c.JSON(http.StatusOK, gin.H{"admin": admin, "token": token})
}

// GetAdmins — GET /api/v1/admin/admins
func (h *AdminHandler) GetAdmins(c *gin.Context) {
	admins, err := h.adminUC.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, admins)
}

// CreateAdmin — POST /api/v1/admin/admins
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var input entity.AdminCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, _, err := h.adminUC.CreateAdmin(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrUsernameTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": "Логин уже занят"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("CreateAdmin: создан администратор '%s'", admin.Username)
	c.JSON(http.StatusCreated, admin)
}

// UpdateAdmin — PUT /api/v1/admin/admins/:id
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	var input entity.AdminUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := h.adminUC.UpdateAdmin(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, usecase.ErrUsernameTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": "Логин уже занят"})
			return
		}
		if errors.Is(err, usecase.ErrAdminNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Администратор не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("UpdateAdmin: обновлён администратор ID=%d", id)
	c.JSON(http.StatusOK, admin)
}

// DeleteAdmin — DELETE /api/v1/admin/admins/:id
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	callerID := middleware.GetCallerID(c)

	err = h.adminUC.DeleteAdmin(c.Request.Context(), id, callerID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrCannotDeleteSelf):
			c.JSON(http.StatusForbidden, gin.H{"error": "Нельзя удалить свой аккаунт"})
		case errors.Is(err, usecase.ErrLastAdmin):
			c.JSON(http.StatusForbidden, gin.H{"error": "Нельзя удалить последнего администратора"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	logger.Info("DeleteAdmin: удалён администратор ID=%d", id)
	c.Status(http.StatusNoContent)
}
