package user_repository

import (
	"context"
	"database/sql"
	"gin-boilerplate/internal/model"
)

type userSqlRepository struct {
	sqlDB *sql.DB
}

// Singleton pattern
var localUserSqlRepository IUserRepository

func InitUserSqlRepository(sqlDB *sql.DB) {
	localUserSqlRepository = &userSqlRepository{
		sqlDB: sqlDB,
	}
}

func UserSqlRepository() IUserRepository {
	return localUserSqlRepository
}

func (r *userSqlRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	id, err := createUser(r.sqlDB, user)
	if err != nil {
		return model.User{}, err
	}
	return getById(r.sqlDB, id)
}

func (r *userSqlRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	return getByUsername(r.sqlDB, username)
}
