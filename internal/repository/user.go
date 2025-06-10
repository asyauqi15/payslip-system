package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

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
