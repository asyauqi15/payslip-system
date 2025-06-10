package repository

import (
	"github.com/asyauqi15/payslip-system/internal/entity"
	"gorm.io/gorm"
)

type Registry struct {
	UserRepository UserRepository
}

func InitializeRepository(db *gorm.DB) *Registry {
	return &Registry{
		UserRepository: NewUserRepository(&BaseRepositoryImpl[entity.User]{DB: db}),
	}
}
