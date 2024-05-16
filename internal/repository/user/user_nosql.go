package user_repository

import (
	"context"
	"gin-boilerplate/consts"
	"gin-boilerplate/database/nosql"
	"gin-boilerplate/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

type userNosqlRepository struct {
	noSqlDB nosql.NoSQLInterface
}

// Singleton pattern
var localUserNosqlRepository IUserRepository

func InitUserNosqlRepository(noSqlDB nosql.NoSQLInterface) {
	localUserNosqlRepository = &userNosqlRepository{
		noSqlDB: noSqlDB,
	}
}

func UserNosqlRepository() IUserRepository {
	return localUserNosqlRepository
}

func (r *userNosqlRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	result, err := r.noSqlDB.CreateIfNotExists(ctx, consts.MongoDBCollectionUsers, bson.M{"username": user.Username}, user)
	if err != nil {
		return model.User{}, err
	}
	err = result.Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *userNosqlRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	result, err := r.noSqlDB.FindOne(ctx, consts.MongoDBCollectionUsers, bson.M{"username": username})
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	err = result.Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
