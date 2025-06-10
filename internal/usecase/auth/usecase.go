package auth

import (
	"context"
	"github.com/asyauqi15/payslip-system/internal/repository"
	jwt_auth "github.com/asyauqi15/payslip-system/pkg/jwt-auth"
)

type Usecase interface {
	Auth(ctx context.Context, email string, password string) (*Result, error)
}

type UsecaseImpl struct {
	userRepo repository.UserRepository
	jwtAuth  *jwt_auth.JWTAuthentication
}

func NewUsecase(userRepo repository.UserRepository, jwtAuth *jwt_auth.JWTAuthentication) Usecase {
	return &UsecaseImpl{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}
