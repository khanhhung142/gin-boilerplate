package userdb

import (
	"habbit-tracker/internal/model"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `db:"password" gorm:"unique"`
	Password string `db:"password"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToBusiness() model.User {
	return model.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}
}
