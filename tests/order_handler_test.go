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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// setupOrderRouter инициализирует тестовую БД, создаёт пользователя и возвращает Gin-роутер и его ID.
func setupOrderRouter(t *testing.T) (*gin.Engine, uint) {
	db := getTestDB(t)
	cleanUsers(t, db)

	// создаём пользователя
	userSvc := services.NewUserService(db, "test-secret")
	user, err := userSvc.Create(context.Background(), &services.RegisterRequest{
		Name:     "Order User",
		Email:    "order@example.com",
		Password: "pass1234",
		Age:      33,
	})
	require.NoError(t, err)

	// роутер для заказов (без JWT-мидлвэра)
	orderH := handlers.NewOrderHandler(db)
	r := gin.New()
	r.POST("/users/:id/orders", orderH.CreateForUser)
	r.GET("/users/:id/orders", orderH.ListByUser)

	return r, user.ID
}

func Test_CreateOrder_Success(t *testing.T) {
	r, userID := setupOrderRouter(t)

	// подготовка тела запроса
	order := map[string]interface{}{
		"product":  "Laptop",
		"quantity": 1,
		"price":    1200.50,
	}
	body, _ := json.Marshal(order)

	req, _ := http.NewRequest("POST", "/users/"+strconv.Itoa(int(userID))+"/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "Laptop", resp["product"])
	require.Equal(t, float64(1), resp["quantity"])
	require.Equal(t, 1200.50, resp["price"])
	require.Equal(t, float64(userID), resp["user_id"])
	require.NotEmpty(t, resp["created_at"])
}

func Test_ListOrders_Success(t *testing.T) {
	r, userID := setupOrderRouter(t)

	// создаём два заказа
	toCreate := []map[string]interface{}{
		{"product": "Monitor", "quantity": 2, "price": 350.00},
		{"product": "Mouse", "quantity": 3, "price": 25.50},
	}
	for _, o := range toCreate {
		body, _ := json.Marshal(o)
		req, _ := http.NewRequest("POST", "/users/"+strconv.Itoa(int(userID))+"/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
	}

	// получаем список заказов
	reqList, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(int(userID))+"/orders", nil)
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)

	require.Equal(t, http.StatusOK, wList.Code)

	var list []map[string]interface{}
	require.NoError(t, json.Unmarshal(wList.Body.Bytes(), &list))
	require.Len(t, list, len(toCreate))

	// проверяем, что в ответе есть оба заказа с нужными полями
	found := make(map[string]map[string]interface{})
	for _, o := range list {
		found[o["product"].(string)] = o
	}
	for _, o := range toCreate {
		f, ok := found[o["product"].(string)]
		require.True(t, ok, "заказ %q не найден", o["product"])
		require.Equal(t, float64(o["quantity"].(int)), f["quantity"])
		require.Equal(t, o["price"], f["price"])
		require.Equal(t, float64(userID), f["user_id"])
	}
}

// setupOrderRouterWithAuth инициализирует тестовую БД, создаёт пользователя и возвращает Gin-роутер с JWT middleware и его ID.
func setupOrderRouterWithAuth(t *testing.T) (*gin.Engine, uint, string) {
	db := getTestDB(t)
	cleanUsers(t, db)

	userSvc := services.NewUserService(db, "test-secret")
	user, err := userSvc.Create(context.Background(), &services.RegisterRequest{
		Name:     "Order User",
		Email:    "order@example.com",
		Password: "pass1234",
		Age:      33,
	})
	require.NoError(t, err)

	token := generateTestToken(user.ID, "test-secret")

	orderH := handlers.NewOrderHandler(db)
	r := gin.New()

	// Настраиваем руты с JWT middleware
	auth := r.Group("/")
	auth.Use(middleware.Auth("test-secret"))
	auth.POST("/users/:id/orders", orderH.CreateForUser)
	auth.GET("/users/:id/orders", orderH.ListByUser)

	return r, user.ID, token
}

// тесты на 422 Unprocessable Entity при неполных или некорректных данных
func Test_CreateOrder_BadRequest(t *testing.T) {
	r, userID, token := setupOrderRouterWithAuth(t)

	testCases := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{"Missing product", map[string]interface{}{"quantity": 1, "price": 10.0}, http.StatusUnprocessableEntity},
		{"Missing quantity", map[string]interface{}{"product": "Item", "price": 10.0}, http.StatusUnprocessableEntity},
		{"Missing price", map[string]interface{}{"product": "Item", "quantity": 1}, http.StatusUnprocessableEntity},
		{"Negative quantity", map[string]interface{}{"product": "Item", "quantity": -1, "price": 10.0}, http.StatusUnprocessableEntity},
		{"Negative price", map[string]interface{}{"product": "Item", "quantity": 1, "price": -10.0}, http.StatusUnprocessableEntity},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/users/"+strconv.Itoa(int(userID))+"/orders", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			require.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

// тесты на 401 Unauthorized при отсутствии/некорректном JWT POST /users/:id/orders
func Test_CreateOrder_Unauthorized(t *testing.T) {
	r, userID, _ := setupOrderRouterWithAuth(t)

	order := map[string]interface{}{
		"product":  "Laptop",
		"quantity": 1,
		"price":    1200.50,
	}
	body, _ := json.Marshal(order)

	req, _ := http.NewRequest("POST", "/users/"+strconv.Itoa(int(userID))+"/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// отсутствие Authorization
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// некорректный токен
	req2, _ := http.NewRequest("POST", "/users/"+strconv.Itoa(int(userID))+"/orders", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer invalidtoken")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusUnauthorized, w2.Code)
}

// тесты на 401 Unauthorized при отсутствии/некорректном JWT GET /users/:id/orders
func Test_ListOrders_Unauthorized(t *testing.T) {
	r, userID, _ := setupOrderRouterWithAuth(t)

	req, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(int(userID))+"/orders", nil)
	// отсутствие Authorization
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// некорректный токен
	req2, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(int(userID))+"/orders", nil)
	req2.Header.Set("Authorization", "Bearer badtoken")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusUnauthorized, w2.Code)
}
