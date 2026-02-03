package main

import (
	"catalog-backend/config"
	"catalog-backend/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создание приложения
	application := app.NewApp(cfg)

	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("Failed to run application: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	log.Println("Shutting down server...")
	application.Close()
	log.Println("Server exited properly")
}
