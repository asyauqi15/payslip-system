package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

//go:generate mockgen -destination=./mock/mock_overtime_repository.go -package=mock github.com/asyauqi15/payslip-system/internal/repository OvertimeRepository
type OvertimeRepository interface {
	BaseRepository[entity.Overtime]
}

type OvertimeRepositoryImpl struct {
	BaseRepositoryImpl[entity.Overtime]
}

func NewOvertimeRepository(db *BaseRepositoryImpl[entity.Overtime]) OvertimeRepository {
	return &OvertimeRepositoryImpl{
		BaseRepositoryImpl: *db,
	}
}
