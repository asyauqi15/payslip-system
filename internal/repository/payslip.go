package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

//go:generate mockgen -destination=./mock/mock_payslip_repository.go -package=mock github.com/asyauqi15/payslip-system/internal/repository PayslipRepository
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
