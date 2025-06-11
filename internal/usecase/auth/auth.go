package auth

import (
	"context"

	"github.com/asyauqi15/payslip-system/internal/entity"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	"github.com/asyauqi15/payslip-system/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Result struct {
	AccessToken  string
	RefreshToken string
}

func (u *UsecaseImpl) Auth(ctx context.Context, username string, password string) (*Result, error) {
	user, err := u.userRepo.FindOneByTemplate(ctx, &entity.User{Username: username}, nil)
	if err != nil {
		logger.Error(ctx, "failed to find user by username", "username", username, "error", err)
		return nil, err
	}
	if user == nil {
		return nil, httppkg.NewNotFoundError("username or password is incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, httppkg.NewUnauthorizedError("username or password is incorrect")
	}

	accessToken, _, err := u.jwtAuth.GenerateAccessToken(ctx, user)
	if err != nil {
		logger.Error(ctx, "failed to generate access token", "user_id", user.ID, "error", err)
		return nil, err
	}

	refreshToken, _, err := u.jwtAuth.GenerateRefreshToken(ctx, user)
	if err != nil {
		logger.Error(ctx, "failed to generate refresh token", "user_id", user.ID, "error", err)
		return nil, err
	}

	return &Result{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
