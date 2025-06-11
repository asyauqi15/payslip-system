package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

type PayrollRepository interface {
	BaseRepository[entity.Payroll]
}

type PayrollRepositoryImpl struct {
	BaseRepositoryImpl[entity.Payroll]
}

func NewPayrollRepository(db *BaseRepositoryImpl[entity.Payroll]) PayrollRepository {
	return &PayrollRepositoryImpl{
		BaseRepositoryImpl: *db,
	}
}
