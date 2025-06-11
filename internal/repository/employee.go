package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

//go:generate mockgen -destination=./mock/mock_employee_repository.go -package=mock github.com/asyauqi15/payslip-system/internal/repository EmployeeRepository
type EmployeeRepository interface {
	BaseRepository[entity.Employee]
}

type EmployeeRepositoryImpl struct {
	BaseRepositoryImpl[entity.Employee]
}

func NewEmployeeRepository(db *BaseRepositoryImpl[entity.Employee]) EmployeeRepository {
	return &EmployeeRepositoryImpl{
		BaseRepositoryImpl: *db,
	}
}
