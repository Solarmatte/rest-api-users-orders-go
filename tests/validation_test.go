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

func setupTestRouter() *gin.Engine {
	db := GetTestDB(nil)
	CleanUsers(nil, db)

	userHandler := handlers.NewUserHandler(db, "test-secret")
	orderHandler := handlers.NewOrderHandler(db)

	r := gin.New()
	r.POST("/users", userHandler.CreateUser)
	r.POST("/users/:id/orders", orderHandler.CreateForUser)
	r.GET("/users/:id/orders", orderHandler.ListByUser)

	protected := r.Group("/")
	protected.Use(middleware.Auth("test-secret"))
	protected.GET("/users/:id", userHandler.GetByID)
	protected.DELETE("/users/:id", userHandler.Delete)

	return r
}

// TestNegativeIDValidation проверяет валидацию отрицательных идентификаторов.
// Тест включает проверку эндпоинтов GetByID и ListByUser с отрицательным ID.
func TestNegativeIDValidation(t *testing.T) {
	r := setupTestRouter()

	// Проверяем, что запрос к GetByID с отрицательным ID возвращает ошибку 400.
	req, _ := http.NewRequest(http.MethodGet, "/users/-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	fmt.Printf("Response: %s\n", w.Body.String())
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "ID должен быть положительным целым числом")

	// Проверяем, что запрос к ListByUser с отрицательным ID возвращает ошибку 400.
	req, _ = http.NewRequest(http.MethodGet, "/users/-1/orders", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "ID должен быть положительным целым числом")

	// Test for CreateForUser with negative ID
	req, _ = http.NewRequest(http.MethodPost, "/users/-1/orders", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "ID должен быть положительным целым числом")
}

func TestCreateOrderForNonExistentUser(t *testing.T) {
	r := setupTestRouter()

	// Test for creating an order for a non-existent user
	req, _ := http.NewRequest(http.MethodPost, "/users/999/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "пользователь не найден")
}

func TestInvalidToken(t *testing.T) {
	r := setupTestRouter()

	// Test for accessing protected route with invalid token
	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Contains(t, w.Body.String(), "недействительный токен")
}
