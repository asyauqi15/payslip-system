package attendance_period

import (
	"context"
	"log/slog"
	"time"

	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
)

type CreateAttendancePeriodUsecase interface {
	CreateAttendancePeriod(ctx context.Context, startDate, endDate time.Time) (*entity.AttendancePeriod, error)
}

type CreateAttendancePeriodUsecaseImpl struct {
	attendancePeriodRepo repository.AttendancePeriodRepository
}

func NewCreateAttendancePeriodUsecase(attendancePeriodRepo repository.AttendancePeriodRepository) CreateAttendancePeriodUsecase {
	return &CreateAttendancePeriodUsecaseImpl{
		attendancePeriodRepo: attendancePeriodRepo,
	}
}

func (u *CreateAttendancePeriodUsecaseImpl) CreateAttendancePeriod(ctx context.Context, startDate, endDate time.Time) (*entity.AttendancePeriod, error) {
	// Validate dates
	if startDate.After(endDate) {
		return nil, httppkg.NewBadRequestError("start date must be before end date")
	}

	// Create the attendance period
	attendancePeriod := &entity.AttendancePeriod{
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, err := u.attendancePeriodRepo.Create(ctx, attendancePeriod, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create attendance period", "error", err)
		return nil, httppkg.NewInternalServerError("failed to create attendance period")
	}

	return result, nil
}
