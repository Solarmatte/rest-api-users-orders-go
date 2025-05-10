package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Order — модель заказа для GORM.
type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Product   string    `gorm:"size:255;not null" json:"product"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"type:numeric(10,2);not null" json:"price"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName жёстко задаёт имя таблицы (если нужно).
func (Order) TableName() string {
	return "orders"
}

// OrderRepo предоставляет CRUD-операции для заказов.
type OrderRepo struct {
	db *gorm.DB
}

// NewOrderRepo создаёт новый OrderRepo.
func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// Create сохраняет новый заказ.
func (r *OrderRepo) Create(ctx context.Context, o *Order) error {
	return r.db.WithContext(ctx).Create(o).Error
}

// ListByUser возвращает заказы пользователя.
func (r *OrderRepo) ListByUser(ctx context.Context, userID uint) ([]Order, error) {
	var orders []Order
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}
