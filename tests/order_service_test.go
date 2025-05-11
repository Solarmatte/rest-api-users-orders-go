package tests

import (
	"context"
	"testing"

	"kvant_task/internal/repositories"
	"kvant_task/internal/services"

	"github.com/stretchr/testify/require"
)

// TestOrderService проверяет методы сервиса заказов, включая создание заказа
// и получение списка заказов для пользователя.
func TestOrderService(t *testing.T) {
	db := getTestDB(t)
	cleanUsers(t, db)

	// First, create a user to attach orders to
	userSvc := services.NewUserService(db, "test-secret")
	user, err := userSvc.Create(context.Background(), &services.RegisterRequest{
		Name:     "Order Tester",
		Email:    "ordertester@example.com",
		Password: "pass123",
		Age:      30,
	})
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	orderSvc := services.NewOrderService(db)

	t.Run("CreateOrder_Success", func(t *testing.T) {
		// Проверяем успешное создание заказа через сервисный слой.
		// Убедимся, что данные заказа корректно сохраняются в базе данных.
		req := &services.CreateOrderRequest{
			Product:  "Gadget",
			Quantity: 3,
			Price:    19.95,
		}
		o, err := orderSvc.Create(context.Background(), user.ID, req)
		require.NoError(t, err)
		require.NotZero(t, o.ID)
		require.Equal(t, user.ID, o.UserID)
		require.Equal(t, "Gadget", o.Product)
		require.Equal(t, 3, o.Quantity)
		require.Equal(t, 19.95, o.Price)

		// verify in DB
		var dbOrder repositories.Order
		err = db.First(&dbOrder, o.ID).Error
		require.NoError(t, err)
		require.Equal(t, user.ID, dbOrder.UserID)
		require.Equal(t, "Gadget", dbOrder.Product)
		require.Equal(t, 3, dbOrder.Quantity)
		require.Equal(t, 19.95, dbOrder.Price)
	})

	t.Run("ListByUser_ReturnsAll", func(t *testing.T) {
		// create additional orders
		orders := []services.CreateOrderRequest{
			{Product: "Widget", Quantity: 1, Price: 5.00},
			{Product: "Thing", Quantity: 2, Price: 12.50},
		}
		for _, req := range orders {
			_, err := orderSvc.Create(context.Background(), user.ID, &req)
			require.NoError(t, err)
		}

		list, err := orderSvc.ListByUser(context.Background(), user.ID)
		require.NoError(t, err)
		// should have 3 orders now
		require.Len(t, list, 3)

		// check that each product from our requests appears exactly once
		found := make(map[string]bool)
		for _, o := range list {
			found[o.Product] = true
		}
		require.True(t, found["Gadget"])
		require.True(t, found["Widget"])
		require.True(t, found["Thing"])
	})
}
