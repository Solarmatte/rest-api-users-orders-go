package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"kvant_task/internal/handlers"
	"kvant_task/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// теперь принимает t *testing.T
func setupTestRouter(t *testing.T) *gin.Engine {
	// прокидываем настоящий t
	db := GetTestDB(t)
	CleanUsers(t, db)

	userHandler := handlers.NewUserHandler(db, "test-secret")
	orderHandler := handlers.NewOrderHandler(db)

	r := gin.New()
	// эндпоинты без авторизации
	r.POST("/users", userHandler.CreateUser)

	// Группа с авторизацией
	auth := r.Group("/")
	auth.Use(middleware.Auth("test-secret"))
	auth.GET("/users/:id", userHandler.GetByID)
	auth.DELETE("/users/:id", userHandler.Delete)
	auth.POST("/users/:id/orders", orderHandler.CreateForUser)
	auth.GET("/users/:id/orders", orderHandler.ListByUser)

	return r
}

func TestNegativeIDValidation(t *testing.T) {
	r := setupTestRouter(t)

	token := generateTestToken(1, "test-secret")

	// GetByID c отрицательным ID
	req, _ := http.NewRequest(http.MethodGet, "/users/-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	fmt.Printf("[DEBUG] GET /users/-1: %s\n", w.Body.String())
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "ID должен быть положительным целым числом")

	// ListByUser c отрицательным ID
	req, _ = http.NewRequest(http.MethodGet, "/users/-1/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	fmt.Printf("[DEBUG] GET /users/-1/orders: %s\n", w.Body.String())
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "ID должен быть положительным целым числом")

	// CreateForUser c отрицательным ID
	req, _ = http.NewRequest(http.MethodPost, "/users/-1/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	fmt.Printf("[DEBUG] POST /users/-1/orders: %s\n", w.Body.String())
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "ID должен быть положительным целым числом")
}

func TestCreateOrderForNonExistentUser(t *testing.T) {
	r := setupTestRouter(t)

	token := generateTestToken(999, "test-secret")

	req, _ := http.NewRequest(http.MethodPost, "/users/999/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "пользователь не найден")
}

func TestInvalidToken(t *testing.T) {
	r := setupTestRouter(t)

	// Создаем запрос с недействительным токеном
	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Отладочная информация
	fmt.Printf("[DEBUG] GET /users/1 with invalid token: %s\n", w.Body.String())

	// Проверяем статус-код и сообщение об ошибке
	require.Equal(t, http.StatusUnauthorized, w.Code, "Ожидался статус 401 Unauthorized")
	require.Contains(t, w.Body.String(), "некорректный токен", "Ожидалось сообщение о некорректном токене")
}
