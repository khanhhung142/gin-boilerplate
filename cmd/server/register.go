package server

import (
	"gin-boilerplate/database/nosql/mongodb"
	user_repository "gin-boilerplate/internal/repository/user"
	auth_usecase "gin-boilerplate/internal/usecase/auth"
)

func Register() {
	noSqlDB := mongodb.MongoDBClient()

	user_repository.InitUserRepository(noSqlDB)
	auth_usecase.InitAuthUsecase(user_repository.UserRepository())

}
