package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// RespondError единообразно возвращает JSON с описанием ошибки.
// Если это ошибка валидации, вернёт список человекочитаемых сообщений.
func RespondError(c *gin.Context, status int, err error) {
	if err == nil {
		c.Status(status)
		return
	}

	// Обработка ошибок валидации from Gin (validator.ValidationErrors)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msgs := make([]string, len(ve))
		for i, fe := range ve {
			// fe.Field() — имя поля, fe.ActualTag() — правило валидации, fe.Param() — параметр (если есть)
			var msg string
			switch fe.ActualTag() {
			case "required":
				msg = fmt.Sprintf("Поле '%s' обязательно для заполнения", fe.Field())
			case "email":
				msg = fmt.Sprintf("Поле '%s' должно быть корректным email", fe.Field())
			case "min":
				msg = fmt.Sprintf("Поле '%s' должно содержать минимум %s символов", fe.Field(), fe.Param())
			case "gt":
				msg = fmt.Sprintf("Поле '%s' должно быть больше %s", fe.Field(), fe.Param())
			default:
				msg = fmt.Sprintf("Поле '%s' не прошло проверку '%s'", fe.Field(), fe.ActualTag())
			}
			msgs[i] = msg
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": msgs})
		return
	}

	// Обычная ошибка
	c.JSON(status, TokenErrorResponse{Error: err.Error()})
}

// ErrorResponse — стандартное сообщение об ошибке.
type ErrorResponse struct {
	Error string `json:"error"`
}

// TokenErrorResponse нужен для унификации в формате json при не-валидации.
type TokenErrorResponse struct {
	Error string `json:"error"`
}

// ValidationErrorResponse — для ошибок валидации.
type ValidationErrorResponse struct {
	Errors []string `json:"errors"`
}
