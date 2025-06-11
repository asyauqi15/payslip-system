package repository

import (
	"context"
	"time"

	"github.com/asyauqi15/payslip-system/internal/entity"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=./mock/mock_attendance_repository.go -package=mock github.com/asyauqi15/payslip-system/internal/repository AttendanceRepository
type AttendanceRepository interface {
	BaseRepository[entity.Attendance]
	CountAttendanceInPeriod(ctx context.Context, employeeID int64, startDate, endDate time.Time, tx *gorm.DB) (int64, error)
}

type AttendanceRepositoryImpl struct {
	BaseRepositoryImpl[entity.Attendance]
}

func NewAttendanceRepository(db *BaseRepositoryImpl[entity.Attendance]) AttendanceRepository {
	return &AttendanceRepositoryImpl{
		BaseRepositoryImpl: *db,
	}
}

func (r *AttendanceRepositoryImpl) CountAttendanceInPeriod(ctx context.Context, employeeID int64, startDate, endDate time.Time, tx *gorm.DB) (int64, error) {
	conn := r.UseTransaction(tx)
	var count int64

	err := conn.WithContext(ctx).Model(&entity.Attendance{}).
		Where("employee_id = ? AND clock_in_time >= ? AND clock_in_time <= ?",
			employeeID,
			startDate.Format("2006-01-02 00:00:00"),
			endDate.Format("2006-01-02 23:59:59")).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}
