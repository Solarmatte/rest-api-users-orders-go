package handlers

import (
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

// NewOrderHandler конструктор.
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
// @Success      201    {object} services.OrderResponse
// @Failure      400    {object} handlers.ErrorResponse
// @Failure      401    {object} handlers.ErrorResponse
// @Failure      422    {object} handlers.ValidationErrorResponse
// @Failure      500    {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id}/orders [post]
func (h *OrderHandler) CreateForUser(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, err)
		return
	}
	var req services.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusUnprocessableEntity, err)
		return
	}
	o, err := h.svc.Create(c.Request.Context(), uint(uid), &req)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err)
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
// @Success      200  {array}   services.OrderResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      401  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id}/orders [get]
func (h *OrderHandler) ListByUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, err)
		return
	}
	list, err := h.svc.ListByUser(c.Request.Context(), uint(id))
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, list)
}
