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
	dbname := getEnv("POSTGRES_DB", "rest-api-db") 
	sslmode := getEnv("POSTGRES_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("не удалось подключиться к PostgreSQL: %v", err)
	}

	// Проверяем, что объект db не nil
	if db == nil {
		t.Fatalf("gorm.Open вернул nil")
	}

	require.NoError(t, db.AutoMigrate(&models.User{}, &repositories.Order{}))

	return db
}

// GetTestDB экспортируемая обертка для getTestDB.
func GetTestDB(t *testing.T) *gorm.DB {
	return getTestDB(t)
}

// cleanUsers очищает таблицы users и orders и сбрасывает последовательности.
func cleanUsers(t *testing.T, db *gorm.DB) {
	err := db.Exec("TRUNCATE TABLE orders, users RESTART IDENTITY CASCADE").Error
	require.NoError(t, err, "не удалось очистить таблицы users и orders")
}

// CleanUsers экспортируемая обертка для cleanUsers.
func CleanUsers(t *testing.T, db *gorm.DB) {
	cleanUsers(t, db)
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
