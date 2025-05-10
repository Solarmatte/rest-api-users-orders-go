package tests

import (
	"fmt"
	"kvant_task/internal/models"
	"kvant_task/internal/repositories"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// getEnv возвращает значение переменной окружения или дефолт.
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// getTestDB создаёт подключение к PostgreSQL и мигрирует модели User и Order.
func getTestDB(t *testing.T) *gorm.DB {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	pass := getEnv("POSTGRES_PASSWORD", "qwerty")
	dbname := getEnv("POSTGRES_DB", "kvant_db")
	sslmode := getEnv("POSTGRES_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err, "не удалось подключиться к PostgreSQL")

	// Мигрируем сразу обе модели
	require.NoError(t, db.AutoMigrate(&models.User{}, &repositories.Order{}))

	return db
}

// cleanUsers очищает таблицы users и orders и сбрасывает последовательности.
func cleanUsers(t *testing.T, db *gorm.DB) {
	err := db.Exec("TRUNCATE TABLE orders, users RESTART IDENTITY CASCADE").Error
	require.NoError(t, err, "не удалось очистить таблицы users и orders")
}

// generateTestToken создаёт JWT токен для тестов сервисов и хендлеров.
func generateTestToken(userID uint, secret string) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tok, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret))
	return tok
}
