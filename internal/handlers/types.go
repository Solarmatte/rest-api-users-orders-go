package handlers

import "kvant_task/internal/services"

// TokenResponse — ответ с JWT токеном.
type TokenResponse struct {
	Token string `json:"token"`
}

type UserListResponse struct {
	Page  int                     `json:"page"`
	Limit int                     `json:"limit"`
	Total int64                   `json:"total"`
	Users []services.UserResponse `json:"users"`
}
