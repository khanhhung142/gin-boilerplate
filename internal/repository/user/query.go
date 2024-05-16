package user_repository

import (
	"database/sql"
	"errors"
	"gin-boilerplate/consts"
	sqlAbs "gin-boilerplate/database/sql"
	"gin-boilerplate/internal/model"
)

// sqlAbs.Database is an interface that abstracts the sql.DB and sql.Tx
// So we can use either one in our repository

func createUser(db sqlAbs.Database, user model.User) (id int64, err error) {
	err = db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// scanUsers scans the rows from the database and returns a slice of users
func scanUsers(rows sqlAbs.Rows) ([]model.User, error) {
	var users []model.User
	err := sqlAbs.ScanRows(rows, func(row sqlAbs.Scannable) error {
		var user model.User
		err := row.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			return err
		}
		users = append(users, user)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func getById(db sqlAbs.Database, id int64) (model.User, error) {
	var user model.User
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, consts.CodeUserNotFound
		}
		return model.User{}, err
	}
	return user, nil
}

func getByUsername(db sqlAbs.Database, username string) (model.User, error) {
	var user model.User
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, consts.CodeUserNotFound
		}
		return model.User{}, err
	}
	return user, nil
}

func searchByUsername(db sqlAbs.Database, username string) ([]model.User, error) {
	rows, err := db.Query("SELECT * FROM users WHERE username LIKE $1", username)
	if err != nil {
		return nil, err
	}
	users, err := scanUsers(rows)
	if err != nil {
		return nil, err
	}
	return users, nil
}
