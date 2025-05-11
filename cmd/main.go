// @title           API для управления пользователями и заказами
// @description     REST API на Go + PostgreSQL с авторизацией, Swagger и фильтрацией

// @contact.name   Junior Golang Developer
// @contact.email  you@example.com

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "kvant_task/docs" // swagger docs

	"kvant_task/internal/bootstrap"
	"kvant_task/internal/config"
	"kvant_task/internal/router"
)

// main.go
// Точка входа в приложение Kvant Task API.
// Запускает HTTP-сервер, подключает БД, обрабатывает сигналы завершения.

// main запускает сервер Kvant Task API.
func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[main] ошибка конфигурации: %v", err)
	}

	// Подключение к базе данных
	db, err := bootstrap.Database(cfg)
	if err != nil {
		log.Fatalf("[main] ошибка БД: %v", err)
	}
	log.Println("[main] Подключение к БД успешно")

	// Инициализация роутера
	r := router.New(db, cfg)

	// HTTP-сервер
	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: r,
	}

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("[main] Сервер запущен на %s", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[main] ListenAndServe: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[main] Получен сигнал завершения, останавливаем сервер...")

	// Корректное завершение работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[main] Shutdown: %v", err)
	}
	log.Println("[main] Сервер остановлен корректно")
}
