// jwt.go
// Этот файл содержит утилиты для работы с JWT токенами.
// Реализует функции для генерации и проверки токенов.

package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenTTL — время жизни токена
const TokenTTL = 72 * time.Hour

// CustomClaims — структура для JWT-клаимов
type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken генерирует JWT токен для пользователя.
func GenerateToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET не задан")
	}

	claims := CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken разбирает JWT и возвращает userID, либо ошибку
func ParseToken(tokenString string) (uint, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, errors.New("JWT_SECRET не задан")
	}

	parsed, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи — HMAC
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("неверный метод подписи токена")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := parsed.Claims.(*CustomClaims)
	if !ok || !parsed.Valid {
		return 0, errors.New("неверные данные токена")
	}

	return claims.UserID, nil
}
