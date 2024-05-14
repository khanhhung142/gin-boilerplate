package server

import (
	"emvn/database/nosql/mongodb"
	user_repository "emvn/internal/repository/user"
	auth_usecase "emvn/internal/usecase/auth"
)

func Register() {
	noSqlDB := mongodb.MongoDBClient()

	user_repository.InitUserRepository(noSqlDB)
	auth_usecase.InitAuthUsecase(user_repository.UserRepository())

}
