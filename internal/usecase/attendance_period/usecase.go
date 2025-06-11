package attendance_period

import (
	"context"
	"time"

	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/asyauqi15/payslip-system/internal/usecase/attendance_period Usecase
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
