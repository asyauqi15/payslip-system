package usecase

import (
	"github.com/asyauqi15/payslip-system/internal/repository"
	authusecase "github.com/asyauqi15/payslip-system/internal/usecase/auth"
	jwtauth "github.com/asyauqi15/payslip-system/pkg/jwt-auth"
)

type Registry struct {
	Auth authusecase.Usecase
}

func InitializeUseCase(repository *repository.Registry, jwt *jwtauth.JWTAuthentication) *Registry {
	return &Registry{
		Auth: authusecase.NewUsecase(repository.UserRepository, jwt),
	}
}
