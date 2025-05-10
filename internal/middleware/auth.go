// auth.go
// Этот файл содержит middleware для аутентификации.
// Реализует проверку JWT токенов для защиты маршрутов.

package middleware

import (
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
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			return
		}
		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "некорректный токен"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "некорректные данные"})
			return
		}
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "некорректные данные"})
			return
		}
		c.Set("user_id", uint(userID))
		c.Next()
	}
}
