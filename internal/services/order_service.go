package services

import (
	"context"
	"log"
	"time"

	"kvant_task/internal/repositories"

	"gorm.io/gorm"
)

// order_service.go
// Этот файл содержит бизнес-логику для работы с заказами.
// Реализует методы для создания, обновления и получения заказов.

// CreateOrderRequest данные для создания заказа.
// Поля точно соответствуют ТЗ.
type CreateOrderRequest struct {
	Product  string  `json:"product" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,gt=0"`
	Price    float64 `json:"price" binding:"required,gt=0"`
}

// OrderResponse DTO для отправки клиенту.
type OrderResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Product   string    `json:"product"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderService бизнес-логика заказов.
type OrderService struct {
	repo *repositories.OrderRepo
}

// NewOrderService создаёт OrderService.
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{repo: repositories.NewOrderRepo(db)}
}

func toOrderResponse(o *repositories.Order) *OrderResponse {
	return &OrderResponse{
		ID:        o.ID,
		UserID:    o.UserID,
		Product:   o.Product,
		Quantity:  o.Quantity,
		Price:     o.Price,
		CreatedAt: o.CreatedAt,
	}
}

// Create создаёт новый заказ и возвращает его DTO.
func (s *OrderService) Create(ctx context.Context, userID uint, req *CreateOrderRequest) (*OrderResponse, error) {
	// Add logging for order creation
	log.Printf("Attempting to create order for user ID: %d", userID)
	o := &repositories.Order{
		UserID:   userID,
		Product:  req.Product,
		Quantity: req.Quantity,
		Price:    req.Price,
	}
	if err := s.repo.Create(ctx, o); err != nil {
		log.Printf("Error creating order: %v", err)
		return nil, err
	}
	log.Printf("Order created successfully with ID: %d", o.ID)
	return toOrderResponse(o), nil
}

// ListByUser возвращает список заказов пользователя.
func (s *OrderService) ListByUser(ctx context.Context, userID uint) ([]OrderResponse, error) {
	list, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]OrderResponse, len(list))
	for i, o := range list {
		out[i] = *toOrderResponse(&o)
	}
	return out, nil
}

func (s *OrderService) GetDB() *gorm.DB {
	return s.repo.GetDB()
}
