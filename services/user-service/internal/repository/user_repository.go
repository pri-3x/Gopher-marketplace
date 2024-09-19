package repository

import (
	"context"
	"gopher-market/services/user-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *UserRepository) List(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}
