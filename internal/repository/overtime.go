package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

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
