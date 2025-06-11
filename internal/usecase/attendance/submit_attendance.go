package attendance

import (
	"context"
	"log/slog"
	"time"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/entity"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/spf13/cast"
)

func (u *UsecaseImpl) SubmitAttendance(ctx context.Context, attendanceType v1.PostEmployeeAttendanceJSONBodyAttendanceType) error {
	// Get the user ID from context
	userIDStr := ctx.Value(constant.ContextKeyUserID)
	if userIDStr == nil {
		return httppkg.NewUnauthorizedError("user not authenticated")
	}

	userID := cast.ToInt64(userIDStr)

	// Find the employee by user ID
	employee, err := u.employeeRepo.FindOneByTemplate(ctx, &entity.Employee{UserID: userID}, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find employee", "user_id", userID, "error", err)
		return httppkg.NewInternalServerError("failed to find employee")
	}
	if employee == nil {
		return httppkg.NewNotFoundError("employee not found")
	}

	currentTime := time.Now().Format(time.RFC3339)
	today := time.Now().Format("2006-01-02")
	now := time.Now()

	// Check if today is a weekday (Monday to Friday)
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return httppkg.NewBadRequestError("attendance can only be submitted on weekdays")
	}

	switch attendanceType {
	case v1.CheckIn:
		return u.handleCheckIn(ctx, employee.ID, currentTime, today)
	case v1.CheckOut:
		return u.handleCheckOut(ctx, employee.ID, currentTime, today)
	default:
		return httppkg.NewBadRequestError("invalid attendance type")
	}
}

func (u *UsecaseImpl) handleCheckIn(ctx context.Context, employeeID int64, currentTime, today string) error {
	// Check if there's already a check-in for today
	existingAttendances, err := u.attendanceRepo.FindByTemplate(ctx, &entity.Attendance{EmployeeID: employeeID}, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find existing attendances", "employee_id", employeeID, "error", err)
		return httppkg.NewInternalServerError("failed to check existing attendance")
	}

	// Check if there's already a check-in for today
	for _, attendance := range existingAttendances {
		if attendance.ClockInTime != "" {
			// Parse the clock-in time to check if it's today
			clockInTime, parseErr := time.Parse(time.RFC3339, attendance.ClockInTime)
			if parseErr == nil && clockInTime.Format("2006-01-02") == today {
				return httppkg.NewConflictError("already checked in today")
			}
		}
	}

	// Create new attendance record with check-in
	attendance := &entity.Attendance{
		EmployeeID:  employeeID,
		ClockInTime: currentTime,
	}

	_, err = u.attendanceRepo.Create(ctx, attendance, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create attendance record", "employee_id", employeeID, "error", err)
		return httppkg.NewInternalServerError("failed to create attendance record")
	}

	return nil
}

func (u *UsecaseImpl) handleCheckOut(ctx context.Context, employeeID int64, currentTime, today string) error {
	// Find today's attendance record that has check-in but no check-out
	existingAttendances, err := u.attendanceRepo.FindByTemplate(ctx, &entity.Attendance{EmployeeID: employeeID}, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find existing attendances", "employee_id", employeeID, "error", err)
		return httppkg.NewInternalServerError("failed to check existing attendance")
	}

	var todayAttendance *entity.Attendance
	for _, attendance := range existingAttendances {
		if attendance.ClockInTime != "" {
			// Parse the clock-in time to check if it's today
			clockInTime, parseErr := time.Parse(time.RFC3339, attendance.ClockInTime)
			if parseErr == nil && clockInTime.Format("2006-01-02") == today {
				todayAttendance = &attendance
				break
			}
		}
	}

	if todayAttendance == nil {
		return httppkg.NewBadRequestError("cannot check out without checking in first")
	}

	if todayAttendance.ClockOutTime != "" {
		return httppkg.NewConflictError("already checked out today")
	}

	// Update the attendance record with check-out time
	_, err = u.attendanceRepo.Updates(ctx, todayAttendance, entity.Attendance{ClockOutTime: currentTime}, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update attendance record", "attendance_id", todayAttendance.ID, "error", err)
		return httppkg.NewInternalServerError("failed to update attendance record")
	}

	return nil
}
