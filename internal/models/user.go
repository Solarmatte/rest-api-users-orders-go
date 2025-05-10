package models

// User — модель пользователя.
// @Description Пользователь системы.
type User struct {
	// ID пользователя
	// required: true
	ID uint `gorm:"primaryKey" json:"id"`

	// Имя пользователя
	// required: true
	Name string `gorm:"not null" json:"name"`

	// Email пользователя
	// required: true
	Email string `gorm:"unique;not null" json:"email"`

	// Возраст пользователя
	// required: true
	Age int `gorm:"not null" json:"age"`

	// Хэш пароля
	// required: true
	PasswordHash string `gorm:"not null" json:"-"`
}
