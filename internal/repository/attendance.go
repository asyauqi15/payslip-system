package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

type AttendanceRepository interface {
	BaseRepository[entity.Attendance]
}

type AttendanceRepositoryImpl struct {
	BaseRepositoryImpl[entity.Attendance]
}

func NewAttendanceRepository(db *BaseRepositoryImpl[entity.Attendance]) AttendanceRepository {
	return &AttendanceRepositoryImpl{
		BaseRepositoryImpl: *db,
	}
}
