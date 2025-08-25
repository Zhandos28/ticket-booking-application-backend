package repository

import (
	"context"

	"github.com/Zhandos28/ticket-booking/model"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func (r *AuthRepository) RegisterUser(ctx context.Context, registerData *model.AuthCredentials) (*model.User, error) {
	user := &model.User{
		Email:    registerData.Email,
		Password: registerData.Password,
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) GetUser(ctx context.Context, query interface{}, args ...interface{}) (*model.User, error) {
	user := &model.User{}

	if err := r.db.WithContext(ctx).Model(user).Where(query, args...).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}
