package auth

import (
	authusecase "github.com/asyauqi15/payslip-system/internal/usecase/auth"
	"net/http"
)

type Handler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type HandlerImpl struct {
	authUsecase authusecase.Usecase
}

func NewHandler(authUsecase authusecase.Usecase) Handler {
	return &HandlerImpl{
		authUsecase: authUsecase,
	}
}
