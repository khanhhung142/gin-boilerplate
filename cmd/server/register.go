package server

import (
	user_repository "habbit-tracker/internal/repository/user"
	auth_usecase "habbit-tracker/internal/usecase/auth"

	"gorm.io/gorm"
)

func Register(sqlDB *gorm.DB) {
	var (
		userRepo = user_repository.InitUserSqlRepository(sqlDB)
	)

	auth_usecase.InitAuthUsecase(userRepo)
}
