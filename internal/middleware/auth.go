// auth.go
// Этот файл содержит middleware для аутентификации.
// Реализует проверку JWT токенов для защиты маршрутов.

package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var invalidatedTokens = make(map[string]struct{})

// InvalidateToken добавляет токен в список аннулированных.
func InvalidateToken(token string) {
	invalidatedTokens[token] = struct{}{}
}

// IsTokenInvalidated проверяет, аннулирован ли токен.
func IsTokenInvalidated(token string) bool {
	_, exists := invalidatedTokens[token]
	return exists
}

// Auth middleware для проверки JWT токенов.
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		fmt.Printf("[DEBUG] Authorization header: %s\n", header)
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("[DEBUG] Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			return
		}
		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			fmt.Printf("[DEBUG] Token parsing error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "некорректный токен"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("[DEBUG] Invalid token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "некорректные данные"})
			return
		}
		userID, ok := claims["user_id"].(float64)
		if !ok {
			fmt.Println("[DEBUG] Missing or invalid user_id in token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "некорректные данные"})
			return
		}
		fmt.Printf("[DEBUG] Token valid, user_id: %v\n", userID)
		c.Set("user_id", uint(userID))
		c.Next()
	}
}
