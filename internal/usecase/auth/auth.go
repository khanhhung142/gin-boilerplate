package auth_usecase

import (
	"context"
	"errors"
	"fmt"
	"habbit-tracker/consts"
	"habbit-tracker/internal/model"
	user_repository "habbit-tracker/internal/repository/user"
	"habbit-tracker/pkg/logger"
)

type IAuthUsecase interface {
	SignUp(ctx context.Context, user model.User) error
	SignIn(ctx context.Context, username, password string) (SignInOutput, error)
}

type authUsecase struct {
	userRepo user_repository.IUserRepository
}

var localAuthUsecase IAuthUsecase

func InitAuthUsecase(userRepo user_repository.IUserRepository) {
	localAuthUsecase = &authUsecase{
		userRepo: userRepo,
	}
}

func AuthUsecase() IAuthUsecase {
	return localAuthUsecase
}

func (u *authUsecase) SignUp(ctx context.Context, user model.User) error {
	// check if user already exists
	dbUser, err := u.userRepo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		if errors.Is(err, consts.CodeUserNotFound) {
			err = nil
		}
	}

	if dbUser.Username == user.Username {
		return consts.CodeUserAlreadyExists
	}

	// Add hash password
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		logger.Error(ctx, "hashPassword", err)
		return consts.CodeInternalError
	}

	_, err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		logger.Error(ctx, "CreateUser", err)
		return consts.CodeInternalError
	}

	return nil
}

func (u *authUsecase) SignIn(ctx context.Context, username, password string) (SignInOutput, error) {
	var output SignInOutput
	dbUser, err := u.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return output, err
	}

	if !checkPasswordHash(password, dbUser.Password) {
		return output, consts.CodeWrongPassword
	}

	// Gen Access Token
	acToken := AccessToken{
		Sub: fmt.Sprintf("%v", dbUser.ID),
		Iss: fmt.Sprintf("%v", dbUser.ID),
	}
	acTokenString, err := acToken.Gen()
	if err != nil {
		return output, err
	}

	output.Token = acTokenString
	output.Exp = acToken.Exp
	return output, nil
}
