package attendance_period

import (
	"context"
	"time"

	"github.com/asyauqi15/payslip-system/internal/entity"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	"github.com/asyauqi15/payslip-system/pkg/logger"
)

func (u *UsecaseImpl) CreateAttendancePeriod(ctx context.Context, startDate, endDate time.Time) (*entity.AttendancePeriod, error) {
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
		logger.Error(ctx, "failed to create attendance period", "error", err)
		return nil, httppkg.NewInternalServerError("failed to create attendance period")
	}

	return result, nil
}
