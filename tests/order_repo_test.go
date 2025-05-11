package tests

import (
	"context"
	"testing"

	"kvant_task/internal/models"
	"kvant_task/internal/repositories"

	"github.com/stretchr/testify/require"
)

// TestOrderRepository проверяет основные операции с репозиторием заказов.
// Тест включает создание пользователя, добавление заказов, получение списка заказов
// и проверку корректности данных.
func TestOrderRepository(t *testing.T) {
	db := getTestDB(t)
	cleanUsers(t, db)

	// First insert a user to satisfy foreign key
	userRepo := repositories.NewUserRepo(db)
	user := &models.User{
		Name:         "RepoOrderUser",
		Email:        "repoorder@example.com",
		Age:          36,
		PasswordHash: "hash",
	}
	require.NoError(t, userRepo.Create(context.Background(), user))
	require.NotZero(t, user.ID)

	orderRepo := repositories.NewOrderRepo(db)

	// 1. Создание заказов
	// Проверяем успешное добавление заказов в базу данных.
	// Убедимся, что идентификаторы заказов не равны нулю.
	o1 := &repositories.Order{
		UserID:   user.ID,
		Product:  "Prod1",
		Quantity: 2,
		Price:    10.5,
	}
	require.NoError(t, orderRepo.Create(context.Background(), o1))
	require.NotZero(t, o1.ID)

	o2 := &repositories.Order{
		UserID:   user.ID,
		Product:  "Prod2",
		Quantity: 5,
		Price:    7.25,
	}
	require.NoError(t, orderRepo.Create(context.Background(), o2))
	require.NotZero(t, o2.ID)

	// 2. ListByUser
	// Проверяем, что метод ListByUser возвращает корректный список заказов
	// для указанного пользователя.
	list, err := orderRepo.ListByUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.Len(t, list, 2)

	// Ensure both products are present
	found := map[string]bool{}
	for _, o := range list {
		found[o.Product] = true
	}
	require.True(t, found["Prod1"])
	require.True(t, found["Prod2"])
}
