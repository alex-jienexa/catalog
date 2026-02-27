package app

import (
	"catalog-backend/config"
	"catalog-backend/internal/delivery/http"
	"catalog-backend/internal/delivery/http/middleware"
	"catalog-backend/internal/infrastructure/database"
	infraRepo "catalog-backend/internal/infrastructure/repository"
	"catalog-backend/internal/usecase"
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type App struct {
	config *config.Config
	db     *sqlx.DB
	router *gin.Engine
}

func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Run() error {
	// Инициализация базы данных
	db, err := initDatabase(a.config.Database.Path)
	if err != nil {
		return err
	}
	a.db = db

	if err := os.MkdirAll(a.config.Upload.Path+"/products", 0755); err != nil {
		return err
	}

	// Инициализация репозиториев
	productRepo := infraRepo.NewSQLiteProductRepository(db)
	sectionRepo := infraRepo.NewSQLiteSectionRepository(db)
	contactRepo := infraRepo.NewSQLiteContactRepository(db)

	// Инициализация use cases
	productUC := usecase.NewProductUseCase(productRepo, sectionRepo, a.config.Upload.Path)
	sectionUC := usecase.NewSectionUseCase(sectionRepo, productRepo)
	contactUC := usecase.NewContactUseCase(contactRepo)

	// Настройка Gin
	gin.SetMode(a.config.Server.Mode)
	a.router = gin.Default()

	// Настройка middleware
	setupMiddleware(a.router)

	// Настройка маршрутов
	router := http.NewRouter(productUC, sectionUC, contactUC, &a.config.Upload)
	router.SetupRoutes(a.router, a.config)

	if err := http.ServeFrontend(a.router); err != nil {
		return err
	}

	// Запуск сервера
	addr := ":" + a.config.Server.Port
	log.Printf("Server starting on %s", addr)
	return a.router.Run(addr)
}

func (a *App) Close() {
	if a.db != nil {
		a.db.Close()
	}
}

func initDatabase(path string) (*sqlx.DB, error) {
	db, err := NewSQLiteDB(path)
	if err != nil {
		return nil, err
	}

	// Запуск миграций
	if err := RunMigrations(db.DB); err != nil {
		return nil, err
	}

	return db, nil
}

// Вспомогательные функции для инициализации базы данных
func NewSQLiteDB(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Включаем foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	return database.RunMigrations(db)
}

func setupMiddleware(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger()) // наш кастомный логер
	router.Use(CORSMiddleware())
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
