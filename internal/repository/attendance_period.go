package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

//go:generate mockgen -destination=./mock/mock_attendance_period_repository.go -package=mock github.com/asyauqi15/payslip-system/internal/repository AttendancePeriodRepository
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
