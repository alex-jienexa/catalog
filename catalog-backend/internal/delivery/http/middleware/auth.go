package middleware

import (
	"catalog-backend/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const AdminIDKey = "admin_id"

// AuthMiddleware проверяет JWT токен из заголовка Authorization
func AuthMiddleware(adminUC *usecase.AdminUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		adminID, err := adminUC.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set(AdminIDKey, adminID)
		c.Next()
	}
}

// GetCallerID — хелпер для получения ID текущего админа из контекста
func GetCallerID(c *gin.Context) int {
	id, _ := c.Get(AdminIDKey)
	adminID, _ := id.(int)
	return adminID
}
