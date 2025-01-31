package user_repository

import (
	"context"
	"habbit-tracker/internal/model"

	"gorm.io/gorm"
)

type UserSqlRepository struct {
	DB *gorm.DB
}

func InitUserSqlRepository(sqlDB *gorm.DB) IUserRepository {
	return &UserSqlRepository{
		DB: sqlDB,
	}
}

func (r *UserSqlRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	id, err := r.createUser(user)
	if err != nil {
		return model.User{}, err
	}
	return r.getById(id)
}

func (r *UserSqlRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	return r.getByUsername(username)
}
