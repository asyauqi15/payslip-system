package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

//go:generate mockgen -destination=./mock/mock_payroll_repository.go -package=mock github.com/asyauqi15/payslip-system/internal/repository PayrollRepository
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
