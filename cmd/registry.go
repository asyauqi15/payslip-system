package cmd

import (
	"log"

	"github.com/asyauqi15/payslip-system/internal"
	"github.com/asyauqi15/payslip-system/internal/handler"
	"github.com/asyauqi15/payslip-system/internal/repository"
	"github.com/asyauqi15/payslip-system/internal/usecase"
	jwtauth "github.com/asyauqi15/payslip-system/pkg/jwt-auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type registry struct {
	handler *handler.Registry
	jwt     *jwtauth.JWTAuthentication
}

func initRegistry(cfg internal.Config, jwt *jwtauth.JWTAuthentication) *registry {
	db, err := gorm.Open(postgres.Open(cfg.Database.Source), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	r := repository.InitializeRepository(db)

	u := usecase.InitializeUseCase(r, jwt)

	h := handler.InitializeHandler(u)

	return &registry{
		handler: h,
		jwt:     jwt,
	}
}
