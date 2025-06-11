package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

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
