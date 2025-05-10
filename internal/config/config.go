// config.go
// Этот файл содержит функции для работы с конфигурацией приложения.
// Реализует загрузку переменных окружения и их валидацию.

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config — все настройки приложения.
type Config struct {
	Server struct {
		Address string
	}
	DB struct {
		DSN string
	}
	JWTSecret string
}

// LoadConfig загружает конфигурацию из переменных окружения.
func LoadConfig() (*Config, error) {
	// пытаемся загрузить .env из корня проекта
	if err := godotenv.Load(); err != nil {
		log.Println("[config] .env не найден, читаем системные переменные")
	} else {
		log.Println("[config] переменные загружены из .env")
	}

	cfg := &Config{}
	// Server
	cfg.Server.Address = getEnv("SERVER_ADDRESS", ":8080")

	// Postgres DSN из отдельных переменных
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	pass := getEnv("POSTGRES_PASSWORD", "qwerty")
	dbname := getEnv("POSTGRES_DB", "rest-api-db")
	ssl := getEnv("POSTGRES_SSLMODE", "disable")
	cfg.DB.DSN = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, ssl,
	)

	// JWT
	cfg.JWTSecret = getEnv("JWT_SECRET", "secret")
	return cfg, nil
}

// getEnv возвращает значение переменной окружения или default
func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
