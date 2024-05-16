package server

import (
	"gin-boilerplate/database/sql/postgres"
	user_repository "gin-boilerplate/internal/repository/user"
	auth_usecase "gin-boilerplate/internal/usecase/auth"
)

func Register() {
	// noSqlDB := mongodb.MongoDBClient()

	// user_repository.InitUserNosqlRepository(noSqlDB)
	// auth_usecase.InitAuthUsecase(user_repository.UserNosqlRepository())

	sqlDB := postgres.PostgresClient()
	user_repository.InitUserSqlRepository(sqlDB)
	auth_usecase.InitAuthUsecase(user_repository.UserSqlRepository())

}
