package bootstrap

import (
	"fmt"

	"kvant_task/internal/config"
	"kvant_task/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database открывает соединение и выполняет миграции.
func Database(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("подключение к БД: %w", err)
	}
	// Авто-миграция моделей
	if err := db.AutoMigrate(&models.User{}, &models.Order{}); err != nil {
		return nil, fmt.Errorf("миграция БД: %w", err)
	}
	return db, nil
}
