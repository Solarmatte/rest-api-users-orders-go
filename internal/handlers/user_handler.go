package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"kvant_task/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler HTTP-слой для пользователей.
type UserHandler struct {
	svc *services.UserService
}

// NewUserHandler конструктор.
func NewUserHandler(db *gorm.DB, jwtSecret string) *UserHandler {
	return &UserHandler{svc: services.NewUserService(db, jwtSecret)}
}

// CreateUser обрабатывает POST /users
// @Summary Создать пользователя
// @Description Создаёт нового пользователя и возвращает его данные.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param input body services.RegisterRequest true "Данные пользователя"
// @Success 201 {object} services.UserResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 422 {object} handlers.ValidationErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Success 400 {object} handlers.ErrorResponse "Пользователь с таким email уже существует"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusUnprocessableEntity, err)
		return
	}
	user, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		if err == services.ErrUserExists {
			RespondError(c, http.StatusBadRequest, fmt.Errorf("пользователь с таким email уже существует"))
			return
		} else {
			RespondError(c, http.StatusInternalServerError, err)
		}
		return
	}
	c.JSON(http.StatusCreated, user)
}

// Login обрабатывает POST /login
// @Summary Аутентификация
// @Description Возвращает JWT по email и паролю.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param input body services.LoginRequest true "Данные для логина"
// @Success 200 {object} services.TokenResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 401 {object} handlers.ErrorResponse
// @Failure 422 {object} handlers.ValidationErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusUnprocessableEntity, err)
		return
	}
	tok, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrInvalidCredentials {
			status = http.StatusUnauthorized
		}
		RespondError(c, status, err)
		return
	}
	c.JSON(http.StatusOK, tok)
}

// List возвращает пользователей с пагинацией и фильтрацией.
// @Summary      Список пользователей
// @Description  Пагинация и фильтрация по возрасту.
// @Tags         Пользователи
// @Produce      json
// @Param        page     query    int     false  "Номер страницы"      default(1)
// @Param        limit    query    int     false  "Размер страницы"     default(10)
// @Param        min_age  query    int     false  "Минимальный возраст"
// @Param        max_age  query    int     false  "Максимальный возраст"
// @Success      200      {object} handlers.UserListResponse
// @Failure      401      {object} handlers.ErrorResponse
// @Failure      500      {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	// Считываем параметры
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	filters := services.UserFilter{
		MinAge: c.Query("min_age"),
		MaxAge: c.Query("max_age"),
		Page:   page,
		Limit:  limit,
	}

	// Считаем общее число
	total, err := h.svc.Count(c.Request.Context(), &filters)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err)
		return
	}

	// Получаем срез пользователей
	users, err := h.svc.List(c.Request.Context(), &filters)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err)
		return
	}

	// Отдаём в формате UserListResponse
	c.JSON(http.StatusOK, UserListResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Users: users,
	})
}

// GetByID возвращает пользователя по ID.
// @Summary      Получить пользователя
// @Description  Возвращает данные пользователя по ID.
// @Tags         Пользователи
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  services.UserResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      401  {object}  handlers.ErrorResponse
// @Failure      404  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		RespondError(c, http.StatusBadRequest, fmt.Errorf("ID должен быть положительным целым числом"))
		return
	}
	u, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrNotFound {
			status = http.StatusNotFound
		}
		RespondError(c, status, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

// Update изменяет данные пользователя.
// @Summary      Обновление пользователя
// @Description  Обновляет имя, email или возраст.
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Param        id     path      int                     true  "ID пользователя"
// @Param        input  body      services.UpdateRequest  true  "Данные для обновления"
// @Success      200    {object}  services.UserResponse
// @Failure      400    {object}  handlers.ErrorResponse
// @Failure      401    {object}  handlers.ErrorResponse
// @Failure      404    {object}  handlers.ErrorResponse
// @Failure      500    {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		RespondError(c, http.StatusBadRequest, fmt.Errorf("ID должен быть положительным целым числом"))
		return
	}
	var req services.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err)
		return
	}
	u, err := h.svc.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrNotFound {
			status = http.StatusNotFound
		}
		RespondError(c, status, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

// Delete удаляет пользователя по ID.
// @Summary      Удаление пользователя
// @Description  Удаляет пользователя по ID.
// @Tags         Пользователи
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      401  {object}  handlers.ErrorResponse
// @Failure      404  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		RespondError(c, http.StatusBadRequest, fmt.Errorf("ID должен быть положительным целым числом"))
		return
	}

	userID, ok := c.Get("user_id")
	if ok && uint(id) == userID.(uint) {
		// Автоматический выход из авторизации
		c.Header("Authorization", "")
	}

	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrNotFound {
			status = http.StatusNotFound
		}
		RespondError(c, status, err)
		return
	}
	c.Status(http.StatusNoContent)
}
