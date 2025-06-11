package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

//go:generate mockgen -destination=./mock/mock_user_repository.go -package=mock github.com/asyauqi15/payslip-system/internal/repository UserRepository
type UserRepository interface {
	BaseRepository[entity.User]
}

type UserRepositoryImpl struct {
	BaseRepositoryImpl[entity.User]
}

func NewUserRepository(db *BaseRepositoryImpl[entity.User]) UserRepository {
	return &UserRepositoryImpl{
		BaseRepositoryImpl: *db,
	}
}
