package user_repository

import (
	"database/sql"
	"errors"
	"habbit-tracker/consts"
	userdb "habbit-tracker/database/user"
	"habbit-tracker/internal/model"
)

// sqlAbs.Database is an interface that abstracts the sql.DB and sql.Tx
// So we can use either one in our repository

func (r *UserSqlRepository) createUser(user model.User) (id int64, err error) {
	var userDB = userdb.User{
		Username: user.Username,
		Password: user.Password,
	}

	err = r.DB.Create(&userDB).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserSqlRepository) getById(id int64) (model.User, error) {
	var userDB userdb.User
	if err := r.DB.First(&userDB, id).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, consts.CodeUserNotFound
		}
		return model.User{}, err
	}
	return userDB.ToBusiness(), nil
}

func (r *UserSqlRepository) getByUsername(username string) (model.User, error) {
	var (
		userDB     userdb.User
		userFilter = userdb.User{Username: username}
	)
	if err := r.DB.Where(userFilter).First(&userDB).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, consts.CodeUserNotFound
		}
		return model.User{}, err
	}
	return userDB.ToBusiness(), nil
}
