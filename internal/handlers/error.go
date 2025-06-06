// error.go
// Этот файл содержит функции для обработки ошибок в HTTP-запросах.
// Основная цель — унификация формата ошибок и логирование.

package handlers

import (
	"errors"
	"fmt"
	"log"
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

	log.Printf("RespondError: status=%d, error=%s", status, err.Error())

	// Обработка ошибок валидации from Gin (validator.ValidationErrors)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msgs := make([]string, len(ve))
		for i, fe := range ve {
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
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{Errors: msgs})
		return
	}

	// Обычная ошибка
	c.JSON(status, ErrorResponse{Error: err.Error()})
}

// HandleError обрабатывает ошибки и возвращает соответствующий HTTP статус и сообщение.
func HandleError(c *gin.Context, err error, notFoundErr error, notFoundMsg string) {
	if errors.Is(err, notFoundErr) {
		RespondError(c, http.StatusNotFound, fmt.Errorf("%s", notFoundMsg))
		return
	}
	if err.Error() == "некорректный ID" || err.Error() == "ID должен быть положительным целым числом" {
		RespondError(c, http.StatusBadRequest, err)
		return
	}
	if err.Error() == "ID должен быть положительным целым числом" {
		RespondError(c, http.StatusBadRequest, err)
		return
	}
	RespondError(c, http.StatusInternalServerError, fmt.Errorf("внутренняя ошибка сервера: %w", err))
}

// ErrorResponse — стандартное сообщение об ошибке.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ValidationErrorResponse — для ошибок валидации.
type ValidationErrorResponse struct {
	Errors []string `json:"errors"`
}
