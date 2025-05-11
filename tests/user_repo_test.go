package tests

import (
	"context"
	"testing"

	"kvant_task/internal/models"
	"kvant_task/internal/repositories"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestUserRepository проверяет основные операции CRUD с репозиторием пользователей.
// Тест включает создание, получение, обновление и удаление пользователей.
func TestUserRepository(t *testing.T) {
	db := getTestDB(t)
	cleanUsers(t, db)

	repo := repositories.NewUserRepo(db)

	// 1. Create
	// Проверяем успешное создание пользователя в базе данных.
	// Убедимся, что идентификатор пользователя не равен нулю.
	u := &models.User{
		Name:         "RepoUser",
		Email:        "repo@example.com",
		Age:          42,
		PasswordHash: "hashedpwd",
	}
	require.NoError(t, repo.Create(context.Background(), u))
	require.NotZero(t, u.ID)

	// 2. GetByEmail
	// Проверяем, что метод GetByEmail возвращает корректного пользователя
	// по указанному email.
	gotByEmail, err := repo.GetByEmail(context.Background(), "repo@example.com")
	require.NoError(t, err)
	require.Equal(t, u.ID, gotByEmail.ID)
	require.Equal(t, "RepoUser", gotByEmail.Name)

	// 3. GetByID
	// Проверяем, что метод GetByID возвращает корректного пользователя
	// по указанному идентификатору.
	gotByID, err := repo.GetByID(context.Background(), u.ID)
	require.NoError(t, err)
	require.Equal(t, u.Email, gotByID.Email)

	// 4. List without filters
	// Проверяем, что метод List возвращает список пользователей без фильтров.
	users, err := repo.List(context.Background(), "", "", 1, 10)
	require.NoError(t, err)
	require.Len(t, users, 1)

	// 5. Update
	// Проверяем успешное обновление данных пользователя.
	// Убедимся, что изменения корректно сохраняются в базе данных.
	gotByID.Name = "RepoUserUpdated"
	require.NoError(t, repo.Update(context.Background(), gotByID))
	fresh, err := repo.GetByID(context.Background(), u.ID)
	require.NoError(t, err)
	require.Equal(t, "RepoUserUpdated", fresh.Name)

	// 6. Delete
	// Проверяем успешное удаление пользователя из базы данных.
	// Убедимся, что пользователь больше не существует.
	require.NoError(t, repo.Delete(context.Background(), u.ID))
	_, err = repo.GetByID(context.Background(), u.ID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
