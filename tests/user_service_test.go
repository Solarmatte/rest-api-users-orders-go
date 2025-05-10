package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"kvant_task/internal/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

// TestUserService covers Create, DuplicateEmail, Login, List, GetByID, Update, Delete.
func TestUserService(t *testing.T) {
	db := getTestDB(t)
	cleanUsers(t, db)

	svc := services.NewUserService(db, "test-secret")

	// 1. Create success
	t.Run("Create_Success", func(t *testing.T) {
		u, err := svc.Create(context.Background(), &services.RegisterRequest{
			Name:     "Alice",
			Email:    "alice@example.com",
			Password: "password123",
			Age:      25,
		})
		require.NoError(t, err)
		require.NotZero(t, u.ID)
		require.Equal(t, "Alice", u.Name)
		require.Equal(t, "alice@example.com", u.Email)
		require.Equal(t, 25, u.Age)
	})

	// 2. Create duplicate
	t.Run("Create_DuplicateEmail", func(t *testing.T) {
		_, err := svc.Create(context.Background(), &services.RegisterRequest{
			Name:     "Alice",
			Email:    "alice@example.com",
			Password: "password123",
			Age:      25,
		})
		require.ErrorIs(t, err, services.ErrUserExists)
	})

	// 3. Login success
	t.Run("Login_Success", func(t *testing.T) {
		tokResp, err := svc.Login(context.Background(), &services.LoginRequest{
			Email:    "alice@example.com",
			Password: "password123",
		})
		require.NoError(t, err)
		require.NotEmpty(t, tokResp.Token)

		parsed, err := jwt.Parse(tokResp.Token, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		})
		require.NoError(t, err)
		require.True(t, parsed.Valid)
		claims := parsed.Claims.(jwt.MapClaims)
		uid, ok := claims["user_id"].(float64)
		require.True(t, ok)
		require.Equal(t, float64(1), uid)
		// exp should be ~now+72h
		exp, ok := claims["exp"].(float64)
		require.True(t, ok)
		require.True(t, int64(exp) > time.Now().Unix())
	})

	// 4. Login invalid
	t.Run("Login_Invalid", func(t *testing.T) {
		_, err := svc.Login(context.Background(), &services.LoginRequest{
			Email:    "alice@example.com",
			Password: "wrongpass",
		})
		require.ErrorIs(t, err, services.ErrInvalidCredentials)
	})

	// 5. List and pagination/filter
	t.Run("List_Filter_Pagination", func(t *testing.T) {
		// create more users with ages 30,40,50
		for i, age := range []int{30, 40, 50} {
			_, err := svc.Create(context.Background(), &services.RegisterRequest{
				Name:     fmt.Sprintf("User%d", i+1),
				Email:    fmt.Sprintf("user%d@example.com", i+1),
				Password: "pwd",
				Age:      age,
			})
			require.NoError(t, err)
		}
		// no filters
		all, err := svc.List(context.Background(), &services.UserFilter{
			MinAge: "",
			MaxAge: "",
			Page:   1,
			Limit:  10,
		})
		require.NoError(t, err)
		require.Len(t, all, 4)

		// min_age=35 → two users (40,50)
		minList, err := svc.List(context.Background(), &services.UserFilter{
			MinAge: "35",
			MaxAge: "",
			Page:   1,
			Limit:  10,
		})
		require.NoError(t, err)
		require.Len(t, minList, 2)
		for _, u := range minList {
			require.GreaterOrEqual(t, u.Age, 35)
		}

		// pagination: page2 limit=2 → two users
		page2, err := svc.List(context.Background(), &services.UserFilter{
			MinAge: "",
			MaxAge: "",
			Page:   2,
			Limit:  2,
		})
		require.NoError(t, err)
		require.Len(t, page2, 2)
	})

	// 6. GetByID and NotFound
	t.Run("GetByID", func(t *testing.T) {
		u, err := svc.GetByID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, uint(1), u.ID)

		_, err = svc.GetByID(context.Background(), 999)
		require.ErrorIs(t, err, services.ErrNotFound)
	})

	// 7. Update and NotFound
	t.Run("Update", func(t *testing.T) {
		newName := "Alice Updated"
		newAge := 26
		updated, err := svc.Update(context.Background(), 1, &services.UpdateRequest{
			Name: &newName,
			Age:  &newAge,
		})
		require.NoError(t, err)
		require.Equal(t, "Alice Updated", updated.Name)
		require.Equal(t, 26, updated.Age)

		_, err = svc.Update(context.Background(), 999, &services.UpdateRequest{
			Name: &newName,
		})
		require.ErrorIs(t, err, services.ErrNotFound)
	})

	// 8. Delete and NotFound
	t.Run("Delete", func(t *testing.T) {
		err := svc.Delete(context.Background(), 1)
		require.NoError(t, err)
		_, err = svc.GetByID(context.Background(), 1)
		require.ErrorIs(t, err, services.ErrNotFound)
	})
}
