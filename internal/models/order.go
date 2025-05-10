package models

import "time"

// Order — модель заказа.
// @Description Заказ, привязанный к пользователю.
type Order struct {
	// ID заказа
	// required: true
	ID uint `gorm:"primaryKey" json:"id"`

	// ID пользователя, сделавшего заказ
	// required: true
	UserID uint `gorm:"not null" json:"user_id"`

	// Наименование продукта
	// required: true
	Product string `gorm:"not null" json:"product"`

	// Количество единиц
	// required: true
	Quantity int `gorm:"not null" json:"quantity"`

	// Цена за единицу
	// required: true
	Price float64 `gorm:"not null" json:"price"`

	// Время создания заказа
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
