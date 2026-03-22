package http

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:web/dist
var staticFiles embed.FS

func ServeFrontend(router *gin.Engine) error {
	// Получаем файловую систему из встроенной папки dist
	dist, err := fs.Sub(staticFiles, "web/dist")
	if err != nil {
		return err
	}

	// Хендлер для отдачи статических файлов (css, js, изображения)
	staticHandler := http.FileServer(http.FS(dist))

	// Обслуживаем статику
	router.GET("/assets/*filepath", func(c *gin.Context) {
		staticHandler.ServeHTTP(c.Writer, c.Request)
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		staticHandler.ServeHTTP(c.Writer, c.Request)
	})

	// Все остальные маршруты, кроме API, отдаём index.html
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Если запрос к API — возвращаем 404 с JSON
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		// Отдаём index.html
		indexFile, err := dist.Open("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}
		defer indexFile.Close()

		// Устанавливаем правильный Content-Type и отдаём файл
		c.Data(http.StatusOK, "text/html; charset=utf-8", func() []byte {
			data, _ := io.ReadAll(indexFile)
			return data
		}())
	})

	return nil
}
