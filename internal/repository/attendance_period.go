package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

type AttendancePeriodRepository interface {
	BaseRepository[entity.AttendancePeriod]
}

type AttendancePeriodRepositoryImpl struct {
	BaseRepositoryImpl[entity.AttendancePeriod]
}

func NewAttendancePeriodRepository(db *BaseRepositoryImpl[entity.AttendancePeriod]) AttendancePeriodRepository {
	return &AttendancePeriodRepositoryImpl{
		BaseRepositoryImpl: *db,
	}
}
