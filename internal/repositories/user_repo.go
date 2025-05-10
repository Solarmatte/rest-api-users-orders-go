// user_repo.go
// Этот файл отвечает за взаимодействие с таблицей пользователей в базе данных.
// Реализует методы для создания, обновления и получения пользователей.

package repositories

import (
	"context"
	"strconv"

	"kvant_task/internal/models"

	"gorm.io/gorm"
)

// UserRepo отвечает за работу с таблицей users.
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&u).Error
	return &u, err
}

func (r *UserRepo) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).
		First(&u, id).Error
	return &u, err
}

func (r *UserRepo) Update(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// List возвращает срез пользователей с фильтрацией по возрасту и пагинацией.
func (r *UserRepo) List(ctx context.Context, minAge, maxAge string, page, limit int) ([]models.User, error) {
	q := r.db.WithContext(ctx).Model(&models.User{})
	if minAge != "" {
		if v, err := strconv.Atoi(minAge); err == nil {
			q = q.Where("age >= ?", v)
		}
	}
	if maxAge != "" {
		if v, err := strconv.Atoi(maxAge); err == nil {
			q = q.Where("age <= ?", v)
		}
	}
	offset := (page - 1) * limit
	var users []models.User
	err := q.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// Count возвращает общее число записей пользователей с учётом фильтров.
func (r *UserRepo) Count(ctx context.Context, minAge, maxAge string) (int64, error) {
	q := r.db.WithContext(ctx).Model(&models.User{})
	if minAge != "" {
		if v, err := strconv.Atoi(minAge); err == nil {
			q = q.Where("age >= ?", v)
		}
	}
	if maxAge != "" {
		if v, err := strconv.Atoi(maxAge); err == nil {
			q = q.Where("age <= ?", v)
		}
	}
	var total int64
	err := q.Count(&total).Error
	return total, err
}
