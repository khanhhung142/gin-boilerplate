package user_repository

import (
	"context"
	"gin-boilerplate/internal/model"
)

// This is example. Choose your database
type IUserRepository interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
}
