// user_service.go
// Этот файл содержит бизнес-логику для работы с пользователями.
// Реализует методы для создания, обновления, удаления и получения пользователей.

package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"kvant_task/internal/models"
	"kvant_task/internal/repositories"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	// ErrUserExists ошибка, если пользователь с таким email уже существует.
	ErrUserExists = errors.New("пользователь с таким email уже существует")
	// ErrInvalidCredentials ошибка, если email или пароль неверны.
	ErrInvalidCredentials = errors.New("неверный email или пароль")
	// ErrNotFound ошибка, если пользователь не найден.
	ErrNotFound = gorm.ErrRecordNotFound
)

// RegisterRequest данные для создания пользователя
// Добавлено описание для Swagger
// @Description Данные для создания нового пользователя
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Age      int    `json:"age" binding:"required,gt=0"`
}

// UserResponse данные пользователя в ответе
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// LoginRequest данные для логина
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse возвращает JWT
type TokenResponse struct {
	Token string `json:"token"`
}

// UpdateRequest данные для обновления пользователя
type UpdateRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=2"`
	Email *string `json:"email" binding:"omitempty,email"`
	Age   *int    `json:"age" binding:"omitempty,gt=0"`
}

// UserFilter фильтры при списке пользователей
type UserFilter struct {
	MinAge string
	MaxAge string
	Page   int
	Limit  int
}

// UserService бизнес-логика по пользователям.
type UserService struct {
	repo      *repositories.UserRepo
	jwtSecret string
}

// NewUserService конструктор
func NewUserService(db *gorm.DB, jwtSecret string) *UserService {
	return &UserService{
		repo:      repositories.NewUserRepo(db),
		jwtSecret: jwtSecret,
	}
}

func toUserResponse(u *models.User) *UserResponse {
	return &UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Age:   u.Age,
	}
}

// Create создаёт пользователя и возвращает его данные
func (s *UserService) Create(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	// Add logging for user creation
	log.Printf("Attempting to create user with email: %s", req.Email)
	if _, err := s.repo.GetByEmail(ctx, req.Email); err == nil {
		log.Printf("User with email %s already exists", req.Email)
		return nil, ErrUserExists
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("Error checking user existence: %v", err)
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating password hash: %v", err)
		return nil, err
	}

	u := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: string(hash),
	}
	if err := s.repo.Create(ctx, u); err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}
	log.Printf("User created successfully with ID: %d", u.ID)

	return toUserResponse(u), nil
}

// Login проверяет учётные данные и возвращает JWT
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	u, err := s.repo.GetByEmail(ctx, req.Email)
	if (err != nil) {
		return nil, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		return nil, ErrInvalidCredentials
	}

	// генерируем токен
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}
	return &TokenResponse{Token: tok}, nil
}

// List возвращает срез DTO пользователей по фильтрам и пагинации.
func (s *UserService) List(ctx context.Context, f *UserFilter) ([]UserResponse, error) {
	users, err := s.repo.List(ctx, f.MinAge, f.MaxAge, f.Page, f.Limit)
	if err != nil {
		return nil, err
	}
	out := make([]UserResponse, len(users))
	for i, u := range users {
		out[i] = *toUserResponse(&u)
	}
	return out, nil
}

// Count возвращает общее количество пользователей под фильтрами.
func (s *UserService) Count(ctx context.Context, f *UserFilter) (int64, error) {
	return s.repo.Count(ctx, f.MinAge, f.MaxAge)
}

// GetByID возвращает пользователя по ID.
func (s *UserService) GetByID(ctx context.Context, id uint) (*UserResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("ID должен быть положительным целым числом")
	}
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUserResponse(u), nil
}

// Update обновляет пользователя.
func (s *UserService) Update(ctx context.Context, id uint, req *UpdateRequest) (*UserResponse, error) {
	// Add logging for user update
	log.Printf("Attempting to update user with ID: %d", id)
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Age != nil {
		u.Age = *req.Age
	}
	if err := s.repo.Update(ctx, u); err != nil {
		log.Printf("Error updating user: %v", err)
		return nil, err
	}
	log.Printf("User updated successfully with ID: %d", u.ID)
	return toUserResponse(u), nil
}

// Delete удаляет пользователя.
func (s *UserService) Delete(ctx context.Context, id uint) error {
	// Add logging for user deletion
	log.Printf("Attempting to delete user with ID: %d", id)
	if err := s.repo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}
	log.Printf("User deleted successfully with ID: %d", id)
	return nil
}
