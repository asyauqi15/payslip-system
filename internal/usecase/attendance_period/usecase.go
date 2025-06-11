package attendance_period

import (
	"context"
	"time"

	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository"
)

type Usecase interface {
	CreateAttendancePeriod(ctx context.Context, startDate, endDate time.Time) (*entity.AttendancePeriod, error)
}

type UsecaseImpl struct {
	attendancePeriodRepo repository.AttendancePeriodRepository
}

func NewUsecase(attendancePeriodRepo repository.AttendancePeriodRepository) Usecase {
	return &UsecaseImpl{
		attendancePeriodRepo: attendancePeriodRepo,
	}
}
