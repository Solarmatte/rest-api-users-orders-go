// order_handler.go
// Этот файл реализует HTTP-слой для работы с заказами.
// Содержит обработчики маршрутов, связанных с заказами.

package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"kvant_task/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderHandler — HTTP-слой для заказов.
type OrderHandler struct {
	svc *services.OrderService
}

// NewOrderHandler конструктор для создания нового OrderHandler.
func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{svc: services.NewOrderService(db)}
}

// CreateForUser создаёт заказ для пользователя.
// @Summary      Создание заказа
// @Description  Создаёт новый заказ для указанного пользователя.
// @Tags         Заказы
// @Accept       json
// @Produce      json
// @Param        id     path      int                     true  "ID пользователя"
// @Param        input  body      services.CreateOrderRequest true "Данные заказа"
// @Success      201    {object} services.OrderResponse "Заказ успешно создан"
// @Failure      400    {object} handlers.ErrorResponse "Некорректный ID пользователя"
// @Failure      422    {object} handlers.ValidationErrorResponse "Ошибка валидации данных заказа"
// @Failure      500    {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Security     BearerAuth
// @Router       /users/{id}/orders [post]
func (h *OrderHandler) CreateForUser(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || uid <= 0 {
		HandleError(c, fmt.Errorf("ID должен быть положительным целым числом"), nil, "ID должен быть положительным целым числом")
		return
	}
	userSvc := services.NewUserService(h.svc.GetDB(), "")
	if _, err := userSvc.GetByID(c.Request.Context(), uint(uid)); err != nil {
		HandleError(c, err, gorm.ErrRecordNotFound, "пользователь не найден")
		return
	}
	var req services.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err, nil, "некорректные данные заказа")
		return
	}
	o, err := h.svc.Create(c.Request.Context(), uint(uid), &req)
	if err != nil {
		HandleError(c, err, nil, "ошибка сервера при создании заказа")
		return
	}
	c.JSON(http.StatusCreated, o)
}

// ListByUser возвращает заказы пользователя.
// @Summary      Список заказов
// @Description  Возвращает все заказы указанного пользователя.
// @Tags         Заказы
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {array}   services.OrderResponse "Список заказов"
// @Failure      400  {object}  handlers.ErrorResponse "Некорректный ID пользователя"
// @Failure      500  {object}  handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Security     BearerAuth
// @Router       /users/{id}/orders [get]
func (h *OrderHandler) ListByUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		HandleError(c, fmt.Errorf("ID должен быть положительным целым числом"), nil, "ID должен быть положительным целым числом")
		return
	}
	list, err := h.svc.ListByUser(c.Request.Context(), uint(id))
	if err != nil {
		HandleError(c, err, nil, "ошибка сервера при получении заказов")
		return
	}
	c.JSON(http.StatusOK, list)
}
