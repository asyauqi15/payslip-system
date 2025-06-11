package overtime

import (
	"context"
	"time"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/entity"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	"github.com/asyauqi15/payslip-system/pkg/logger"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/spf13/cast"
)

func (u *UsecaseImpl) SubmitOvertime(ctx context.Context, req v1.OvertimeRequest) error {
	// Get the user ID from context
	userIDStr := ctx.Value(constant.ContextKeyUserID)
	if userIDStr == nil {
		return httppkg.NewUnauthorizedError("user not authenticated")
	}

	userID := cast.ToInt64(userIDStr)

	// Find the employee by user ID
	employee, err := u.employeeRepo.FindOneByTemplate(ctx, &entity.Employee{UserID: userID}, nil)
	if err != nil {
		logger.Error(ctx, "failed to find employee", "user_id", userID, "error", err)
		return httppkg.NewInternalServerError("failed to find employee")
	}
	if employee == nil {
		return httppkg.NewNotFoundError("employee not found")
	}

	// Validate overtime times
	if req.EndTime.Before(req.StartTime) {
		return httppkg.NewBadRequestError("end time must be after start time")
	}

	// Check maximum 3 hours per day
	duration := req.EndTime.Sub(req.StartTime)
	if duration > 3*time.Hour {
		return httppkg.NewBadRequestError("maximum overtime per day is 3 hours")
	}

	// Get the date of the overtime
	overtimeDate := req.StartTime.Format("2006-01-02")

	// Check if it's weekend
	isWeekend := req.StartTime.Weekday() == time.Saturday || req.StartTime.Weekday() == time.Sunday

	if !isWeekend {
		// Weekday rules: must have checked out and overtime must start after 5PM
		if err := u.validateWeekdayOvertime(ctx, req.StartTime); err != nil {
			return err
		}
	}

	// Check for overlapping overtime records and daily limit
	if err := u.validateOvertimeConflicts(ctx, employee.ID, req.StartTime, req.EndTime, overtimeDate); err != nil {
		return err
	}

	// Create overtime record
	overtime := &entity.Overtime{
		EmployeeID:  employee.ID,
		StartAt:     req.StartTime,
		EndAt:       req.EndTime,
		Description: req.Description,
	}

	_, err = u.overtimeRepo.Create(ctx, overtime, nil)
	if err != nil {
		logger.Error(ctx, "failed to create overtime", "employee_id", employee.ID, "error", err)
		return httppkg.NewInternalServerError("failed to submit overtime")
	}

	logger.Info(ctx, "overtime submitted successfully",
		"employee_id", employee.ID,
		"start_time", req.StartTime,
		"end_time", req.EndTime,
		"is_weekend", isWeekend)

	return nil
}

func (u *UsecaseImpl) validateWeekdayOvertime(ctx context.Context, startTime time.Time) error {
	// Check if overtime starts after 5PM (17:00)
	hour := startTime.Hour()
	if hour < 17 {
		return httppkg.NewBadRequestError("overtime on weekdays can only start after 5PM")
	}

	return nil
}

func (u *UsecaseImpl) validateOvertimeConflicts(ctx context.Context, employeeID int64, startTime, endTime time.Time, date string) error {
	// Get all existing overtimes for the employee
	existingOvertimes, err := u.overtimeRepo.FindByTemplate(ctx, &entity.Overtime{EmployeeID: employeeID}, nil)
	if err != nil {
		logger.Error(ctx, "failed to find existing overtimes", "employee_id", employeeID, "error", err)
		return httppkg.NewInternalServerError("failed to validate overtime period")
	}

	var totalDailyOvertimeDuration time.Duration
	for _, existing := range existingOvertimes {
		// Check for time overlap
		if startTime.Before(existing.EndAt) && endTime.After(existing.StartAt) {
			return httppkg.NewBadRequestError("overtime period overlaps with existing overtime")
		}

		// Calculate total overtime for the same date
		if existing.StartAt.Format("2006-01-02") == date {
			totalDailyOvertimeDuration += existing.EndAt.Sub(existing.StartAt)
		}
	}

	// Check if adding this overtime would exceed 3 hours daily limit
	newOvertimeDuration := endTime.Sub(startTime)
	if totalDailyOvertimeDuration+newOvertimeDuration > 3*time.Hour {
		return httppkg.NewBadRequestError("total overtime for the day cannot exceed 3 hours")
	}

	return nil
}
