package http

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed all:web/dist
var staticFiles embed.FS

func ServeFrontend(router *gin.Engine) error {
	// Создаем под-файловую систему из встроенной папки dist
	dist, err := fs.Sub(staticFiles, "web/dist")
	if err != nil {
		return err
	}
	// Отдаем статику для всех маршрутов, кроме /api/
	router.NoRoute(gin.WrapH(http.FileServer(http.FS(dist))))
	return nil
}
