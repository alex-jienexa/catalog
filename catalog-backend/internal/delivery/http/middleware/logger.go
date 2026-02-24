package middleware

import (
	"time"

	"github.com/ZeRg0912/logger"
	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Обработка запроса
		c.Next()

		// После выполнения – собираем информацию
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// Логируем
		logger.Info("[HTTP] %s %s | %d | %v | %s",
			method, path, statusCode, latency, clientIP)

		// Если статус ошибки (>=400), дублируем в error-лог с подробностями
		if statusCode >= 400 {
			// Можно добавить тело ответа, если оно было записано
			// Для этого потребуется буферизация, но для простоты ограничимся
			logger.Error("[HTTP ERROR] %s %s | %d | %v | %s",
				method, path, statusCode, latency, clientIP)
		}
	}
}
