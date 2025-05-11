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
// Этот файл является точкой входа в приложение.
// Реализует запуск сервера, подключение к базе данных и обработку сигналов завершения.

func main() {
	// Загрузка конфигурации из файла или переменных окружения
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[main] конфигурация: %v", err)
	}

	// Подключение к базе данных
	db, err := bootstrap.Database(cfg)
	if err != nil {
		log.Fatalf("[main] БД: %v", err)
	}
	log.Println("[main] успешно подключились к БД")

	// Создание нового роутера с конфигурацией и подключением к БД
	r := router.New(db, cfg)

	// Настройка HTTP-сервера
	srv := &http.Server{
		Addr:    cfg.Server.Address, // Адрес сервера
		Handler: r,                  // Роутер для обработки запросов
	}

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("[main] сервер запущен на %s", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[main] ListenAndServe: %v", err)
		}
	}()

	// Ожидание сигнала завершения (SIGINT или SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[main] получен сигнал завершения")

	// Завершение работы сервера с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[main] Shutdown: %v", err)
	}
	log.Println("[main] сервер остановлен")
}
