package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"kvant_task/internal/handlers"
	"kvant_task/internal/middleware"
	"kvant_task/internal/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// setupUserRouter инициализирует тестовую БД, роутер и возвращает Gin-Engine.
func setupUserRouter(t *testing.T) *gin.Engine {
	db := getTestDB(t)
	cleanUsers(t, db)

	userH := handlers.NewUserHandler(db, "test-secret")

	r := gin.New()
	// Public
	r.POST("/users", userH.CreateUser)
	r.POST("/auth/login", userH.Login)

	// Protected
	auth := r.Group("/")
	auth.Use(middleware.Auth("test-secret"))
	auth.GET("/users", userH.List)
	auth.GET("/users/:id", userH.GetByID)
	auth.PUT("/users/:id", userH.Update)
	auth.DELETE("/users/:id", userH.Delete)

	return r
}

func Test_CreateUser_Success(t *testing.T) {
	r := setupUserRouter(t)

	body := map[string]interface{}{
		"name":     "Test User",
		"email":    "test@example.com",
		"age":      25,
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var resp services.UserResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "Test User", resp.Name)
	require.Equal(t, "test@example.com", resp.Email)
	require.Equal(t, 25, resp.Age)
}

func Test_Login_Success(t *testing.T) {
	r := setupUserRouter(t)
	db := getTestDB(t)
	cleanUsers(t, db)

	// Создаём пользователя напрямую через сервис
	svc := services.NewUserService(db, "test-secret")
	created, err := svc.Create(context.Background(), &services.RegisterRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "pass1234",
		Age:      28,
	})
	require.NoError(t, err)

	// Логинимся
	loginBody := map[string]string{
		"email":    "john@example.com",
		"password": "pass1234",
	}
	jsonLogin, _ := json.Marshal(loginBody)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var tokResp services.TokenResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &tokResp))
	require.NotEmpty(t, tokResp.Token)

	// Парсим JWT напрямую
	parsed, err := jwt.Parse(tokResp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)
	claims := parsed.Claims.(jwt.MapClaims)
	userIDf, ok := claims["user_id"].(float64)
	require.True(t, ok)
	require.Equal(t, float64(created.ID), userIDf)
}

func Test_ListUsers(t *testing.T) {
	r := setupUserRouter(t)

	// Подготовка чистой БД и создание двух пользователей
	db := getTestDB(t)
	cleanUsers(t, db)
	svc := services.NewUserService(db, "test-secret")
	_, _ = svc.Create(context.Background(), &services.RegisterRequest{
		Name:     "A",
		Email:    "a@example.com",
		Password: "p1",
		Age:      20,
	})
	_, _ = svc.Create(context.Background(), &services.RegisterRequest{
		Name:     "B",
		Email:    "b@example.com",
		Password: "p2",
		Age:      30,
	})

	// Получаем валидный токен
	token := generateTestToken(1, "test-secret")

	req, _ := http.NewRequest("GET", "/users?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// Парсим в новую структуру
	type listResp struct {
		Page  int                     `json:"page"`
		Limit int                     `json:"limit"`
		Total int64                   `json:"total"`
		Users []services.UserResponse `json:"users"`
	}
	var resp listResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	require.Equal(t, 1, resp.Page)
	require.Equal(t, 10, resp.Limit)
	require.Equal(t, int64(2), resp.Total)
	require.Len(t, resp.Users, 2)
}

func Test_GetUpdateDeleteUser(t *testing.T) {
	r := setupUserRouter(t)
	db := getTestDB(t)
	cleanUsers(t, db)

	// создаём пользователя
	svc := services.NewUserService(db, "test-secret")
	created, err := svc.Create(context.Background(), &services.RegisterRequest{
		Name:     "C",
		Email:    "c@example.com",
		Password: "p3",
		Age:      40,
	})
	require.NoError(t, err)
	token := generateTestToken(created.ID, "test-secret")

	// GET /users/:id
	reqGet, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(int(created.ID)), nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, reqGet)
	require.Equal(t, http.StatusOK, wGet.Code)

	// PUT /users/:id
	update := map[string]interface{}{"name": "C2", "age": 41}
	jsonUpd, _ := json.Marshal(update)
	reqUpd, _ := http.NewRequest("PUT", "/users/"+strconv.Itoa(int(created.ID)), bytes.NewBuffer(jsonUpd))
	reqUpd.Header.Set("Content-Type", "application/json")
	reqUpd.Header.Set("Authorization", "Bearer "+token)
	wUpd := httptest.NewRecorder()
	r.ServeHTTP(wUpd, reqUpd)
	require.Equal(t, http.StatusOK, wUpd.Code)

	// DELETE /users/:id
	reqDel, _ := http.NewRequest("DELETE", "/users/"+strconv.Itoa(int(created.ID)), nil)
	reqDel.Header.Set("Authorization", "Bearer "+token)
	wDel := httptest.NewRecorder()
	r.ServeHTTP(wDel, reqDel)
	require.Equal(t, http.StatusNoContent, wDel.Code)
}

// тест на 422 Unprocessable Entity при неполном JSON POST /users
func Test_CreateUser_BadRequest(t *testing.T) {
	r := setupUserRouter(t)

	testCases := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{"Missing email", map[string]interface{}{"name": "User 1", "password": "pass1", "age": 20}, http.StatusUnprocessableEntity},
		{"Missing password", map[string]interface{}{"name": "User 2", "email": "u2@example.com", "age": 30}, http.StatusUnprocessableEntity},
		{"Missing name", map[string]interface{}{"email": "u3@example.com", "password": "pass3", "age": 25}, http.StatusUnprocessableEntity},
		{"Missing age", map[string]interface{}{"email": "u4@example.com", "password": "pass4", "name": "User 4"}, http.StatusUnprocessableEntity},
		{"Invalid email format", map[string]interface{}{"name": "User 5", "email": "invalid", "password": "pass5", "age": 25}, http.StatusUnprocessableEntity},
		{"Age zero", map[string]interface{}{"name": "User 6", "email": "u6@example.com", "password": "pass6", "age": 0}, http.StatusUnprocessableEntity},
		{"Password too short", map[string]interface{}{"name": "User 7", "email": "u7@example.com", "password": "123", "age": 30}, http.StatusUnprocessableEntity},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			require.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

// тесты на 401 Unauthorized при обращении к защищённым эндпоинтам без/с неверным токеном
func Test_Endpoints_Unauthorized(t *testing.T) {
	r := setupUserRouter(t)
	// создаём пользователя, чтобы знать id
	db := getTestDB(t)
	cleanUsers(t, db)
	svc := services.NewUserService(db, "test-secret")
	user, err := svc.Create(context.Background(), &services.RegisterRequest{
		Name:     "ForAuth",
		Email:    "auth@example.com",
		Password: "password",
		Age:      30,
	})
	require.NoError(t, err)

	paths := []struct {
		method string
		route  string
		body   map[string]interface{}
	}{
		{"GET", "/users", nil},
		{"GET", "/users/" + strconv.Itoa(int(user.ID)), nil},
		{"PUT", "/users/" + strconv.Itoa(int(user.ID)), map[string]interface{}{"name": "Updated"}},
		{"DELETE", "/users/" + strconv.Itoa(int(user.ID)), nil},
	}

	for _, p := range paths {
		t.Run(p.method+" "+p.route+" without token", func(t *testing.T) {
			var req *http.Request
			if p.body != nil {
				b, _ := json.Marshal(p.body)
				req, _ = http.NewRequest(p.method, p.route, bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(p.method, p.route, nil)
			}
			// без заголовка Authorization
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			require.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run(p.method+" "+p.route+" with invalid token", func(t *testing.T) {
			var req *http.Request
			if p.body != nil {
				b, _ := json.Marshal(p.body)
				req, _ = http.NewRequest(p.method, p.route, bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(p.method, p.route, nil)
			}
			req.Header.Set("Authorization", "Bearer invalidtoken")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			require.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}
