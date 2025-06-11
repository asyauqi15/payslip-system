package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

type PayslipRepository interface {
	BaseRepository[entity.Payslip]
}

type PayslipRepositoryImpl struct {
	BaseRepositoryImpl[entity.Payslip]
}

func NewPayslipRepository(db *BaseRepositoryImpl[entity.Payslip]) PayslipRepository {
	return &PayslipRepositoryImpl{
		BaseRepositoryImpl: *db,
	}
}
