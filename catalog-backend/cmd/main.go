package main

import (
	"catalog-backend/config"
	"catalog-backend/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZeRg0912/logger"
)

func main() {
	// Инициализация логгера
	err := logger.InitBoth(logger.LevelInfo, logger.LevelDebug, "logs/app.log", 10*1024*1024)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	// Закрывать логгер при выходе из программы, чтобы логи были записаны
	defer logger.Close()

	logger.Info("Launching programm...")

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Создание приложения
	application := app.NewApp(cfg)

	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := application.Run(); err != nil {
			logger.Error("Failed to run application: %v", err)
			os.Exit(1)
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	log.Println("Shutting down server...")
	application.Close()
	log.Println("Server exited properly")
}
