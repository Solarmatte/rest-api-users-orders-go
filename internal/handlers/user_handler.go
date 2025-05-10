package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"kvant_task/internal/services"
	"kvant_task/internal/middleware"

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
// @Success 201 {object} services.UserResponse "Пользователь успешно создан"
// @Failure 400 {object} handlers.ErrorResponse "Пользователь с таким email уже существует"
// @Failure 422 {object} handlers.ValidationErrorResponse "Ошибка валидации данных"
// @Failure 500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, fmt.Errorf("некорректные данные: %w", err))
		return
	}
	user, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		HandleError(c, err, services.ErrUserExists, "пользователь с таким email уже существует")
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
// @Success 200 {object} services.TokenResponse "Успешная аутентификация"
// @Failure 400 {object} handlers.ErrorResponse "Некорректные данные для входа"
// @Failure 401 {object} handlers.ErrorResponse "Неверный email или пароль"
// @Failure 422 {object} handlers.ValidationErrorResponse "Ошибка валидации данных"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, fmt.Errorf("некорректные данные: %w", err))
		return
	}
	tok, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		HandleError(c, err, services.ErrInvalidCredentials, "неверный email или пароль")
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
// @Success      200      {object} handlers.UserListResponse "Список пользователей"
// @Failure      401      {object} handlers.ErrorResponse "Неавторизованный доступ"
// @Failure      500      {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Security     BearerAuth
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		HandleError(c, fmt.Errorf("некорректный номер страницы"), nil, "номер страницы должен быть положительным целым числом")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		HandleError(c, fmt.Errorf("некорректный размер страницы"), nil, "размер страницы должен быть положительным целым числом")
		return
	}
	minAge := c.DefaultQuery("min_age", "")
	if minAge != "" {
		if age, err := strconv.Atoi(minAge); err != nil || age < 0 {
			HandleError(c, fmt.Errorf("некорректный минимальный возраст"), nil, "минимальный возраст должен быть неотрицательным целым числом")
			return
		}
	}
	maxAge := c.DefaultQuery("max_age", "")
	if maxAge != "" {
		if age, err := strconv.Atoi(maxAge); err != nil || age < 0 {
			HandleError(c, fmt.Errorf("некорректный максимальный возраст"), nil, "максимальный возраст должен быть неотрицательным целым числом")
			return
		}
	}
	filters := services.UserFilter{
		MinAge: c.Query("min_age"),
		MaxAge: c.Query("max_age"),
		Page:   page,
		Limit:  limit,
	}

	// Считаем общее число
	total, err := h.svc.Count(c.Request.Context(), &filters)
	if err != nil {
		HandleError(c, err, nil, "ошибка при подсчете пользователей")
		return
	}

	// Получаем срез пользователей
	users, err := h.svc.List(c.Request.Context(), &filters)
	if err != nil {
		HandleError(c, err, nil, "ошибка при получении списка пользователей")
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
		HandleError(c, fmt.Errorf("некорректный ID"), nil, "ID должен быть положительным целым числом")
		return
	}
	u, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		HandleError(c, err, services.ErrNotFound, "пользователь не найден")
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
		HandleError(c, err, services.ErrNotFound, "пользователь не найден")
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
		// Логирование для отладки
		log.Printf("[Delete] Пользователь удаляет свою учетную запись: user_id=%d", userID)
		// Автоматический выход из авторизации
		c.Header("Authorization", "")
	} else {
		log.Printf("[Delete] Удаление учетной записи другого пользователя: user_id=%d, target_id=%d", userID, id)
	}

	// Check if user exists before attempting to delete
	if _, err := h.svc.GetByID(c.Request.Context(), uint(id)); err != nil {
		HandleError(c, err, gorm.ErrRecordNotFound, "пользователь не найден")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		HandleError(c, err, services.ErrNotFound, "пользователь не найден")
		return
	}
	// Invalidate token logic
	token := c.GetHeader("Authorization")
	if token != "" {
		parts := strings.SplitN(token, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			middleware.InvalidateToken(parts[1])
		}
	}
	log.Printf("[Delete] Учетная запись успешно удалена: id=%d", id)
	c.Status(http.StatusNoContent)
}
