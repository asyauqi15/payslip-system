package reimbursement

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

func (u *UsecaseImpl) SubmitReimbursement(ctx context.Context, req v1.ReimbursementRequest) error {
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

	// Validate reimbursement amount
	if req.Amount <= 0 {
		return httppkg.NewBadRequestError("reimbursement amount must be greater than 0")
	}

	// Validate that the date is not in the future
	reimbursementDate := time.Time(req.Date.Time)
	if reimbursementDate.After(time.Now()) {
		return httppkg.NewBadRequestError("reimbursement date cannot be in the future")
	}

	// Create reimbursement record
	reimbursement := &entity.Reimbursement{
		EmployeeID:  employee.ID,
		Amount:      int64(req.Amount),
		Date:        reimbursementDate,
		Description: req.Description,
	}

	_, err = u.reimbursementRepo.Create(ctx, reimbursement, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create reimbursement", "employee_id", employee.ID, "error", err)
		return httppkg.NewInternalServerError("failed to submit reimbursement")
	}

	slog.InfoContext(ctx, "reimbursement submitted successfully",
		"employee_id", employee.ID,
		"amount", req.Amount,
		"date", reimbursementDate.Format("2006-01-02"),
		"description", req.Description)

	return nil
}
