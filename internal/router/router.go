package router

import (
	"kvant_task/internal/config"
	"kvant_task/internal/handlers"
	"kvant_task/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// New создаёт Gin-Engine и регистрирует маршруты.
func New(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Хендлеры
	userH := handlers.NewUserHandler(db, cfg.JWTSecret)
	orderH := handlers.NewOrderHandler(db)

	// Публичные
	r.POST("/users", userH.CreateUser)
	r.POST("/auth/login", userH.Login) // <- изменённый маршрут

	// Защищённые — все ниже требуют Bearer токен
	auth := r.Group("/")
	auth.Use(middleware.Auth(cfg.JWTSecret))

	// Пользователи
	auth.GET("/users", userH.List)
	auth.GET("/users/:id", userH.GetByID)
	auth.PUT("/users/:id", userH.Update)
	auth.DELETE("/users/:id", userH.Delete)

	// Заказы вложенно
	auth.POST("/users/:id/orders", orderH.CreateForUser)
	auth.GET("/users/:id/orders", orderH.ListByUser)

	return r
}
