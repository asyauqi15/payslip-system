package auth

import (
	"context"
	"github.com/asyauqi15/payslip-system/internal/entity"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type Result struct {
	AccessToken  string
	RefreshToken string
}

func (u *UsecaseImpl) Auth(ctx context.Context, email string, password string) (*Result, error) {
	user, err := u.userRepo.FindOneByTemplate(ctx, &entity.User{Email: email}, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find user by email", "email", email, "error", err)
		return nil, err
	}
	if user == nil {
		return nil, httppkg.NewNotFoundError("email or password is incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, httppkg.NewUnauthorizedError("email or password is incorrect")
	}

	accessToken, _, err := u.jwtAuth.GenerateAccessToken(ctx, user)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate access token", "user_id", user.ID, "error", err)
		return nil, err
	}

	refreshToken, _, err := u.jwtAuth.GenerateRefreshToken(ctx, user)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate refresh token", "user_id", user.ID, "error", err)
		return nil, err
	}

	return &Result{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
