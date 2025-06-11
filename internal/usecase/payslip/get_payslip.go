package payslip

import (
	"context"
	"time"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/entity"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	"github.com/asyauqi15/payslip-system/pkg/logger"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cast"
)

func (u *UsecaseImpl) GetPayslip(ctx context.Context, payrollID int64) (*v1.PayslipResponse, error) {
	// Get employee ID from context
	userID := ctx.Value(constant.ContextKeyUserID)
	if userID == nil {
		return nil, httppkg.NewUnauthorizedError("user not authenticated")
	}

	// Find employee by user ID
	employee, err := u.employeeRepo.FindOneByTemplate(ctx, &entity.Employee{
		UserID: cast.ToInt64(userID),
	}, nil)
	if err != nil {
		logger.Error(ctx, "failed to find employee", "user_id", userID, "error", err)
		return nil, httppkg.NewInternalServerError("failed to find employee")
	}
	if employee == nil {
		return nil, httppkg.NewNotFoundError("employee not found")
	}

	// Get payroll record
	payroll, err := u.payrollRepo.FindByID(ctx, uint(payrollID), nil)
	if err != nil {
		logger.Error(ctx, "failed to find payroll", "payroll_id", payrollID, "error", err)
		return nil, httppkg.NewInternalServerError("failed to find payroll")
	}
	if payroll == nil {
		return nil, httppkg.NewNotFoundError("payroll not found")
	}

	// Get attendance period
	attendancePeriod, err := u.attendancePeriodRepo.FindByID(ctx, uint(payroll.AttendancePeriodID), nil)
	if err != nil {
		logger.Error(ctx, "failed to find attendance period", "attendance_period_id", payroll.AttendancePeriodID, "error", err)
		return nil, httppkg.NewInternalServerError("failed to find attendance period")
	}
	if attendancePeriod == nil {
		return nil, httppkg.NewInternalServerError("attendance period not found")
	}

	// Find payslip for this employee and payroll
	payslip, err := u.payslipRepo.FindOneByTemplate(ctx, &entity.Payslip{
		EmployeeID: employee.ID,
		PayrollID:  payrollID,
	}, nil)
	if err != nil {
		logger.Error(ctx, "failed to find payslip", "employee_id", employee.ID, "payroll_id", payrollID, "error", err)
		return nil, httppkg.NewInternalServerError("failed to find payslip")
	}
	if payslip == nil {
		return nil, httppkg.NewNotFoundError("payslip not found for this payroll")
	}

	// Get reimbursements for this employee in the attendance period
	reimbursements, err := u.reimbursementRepo.FindByTemplate(ctx, &entity.Reimbursement{
		EmployeeID: employee.ID,
	}, nil)
	if err != nil {
		logger.Error(ctx, "failed to find reimbursements", "employee_id", employee.ID, "error", err)
		return nil, httppkg.NewInternalServerError("failed to find reimbursements")
	}

	// Filter reimbursements by date range
	var reimbursementItems []v1.ReimbursementItem
	for _, reimbursement := range reimbursements {
		if reimbursement.Date.After(attendancePeriod.StartDate.AddDate(0, 0, -1)) &&
			reimbursement.Date.Before(attendancePeriod.EndDate.AddDate(0, 0, 1)) {
			reimbursementItems = append(reimbursementItems, v1.ReimbursementItem{
				Date:        datePtr(reimbursement.Date),
				Amount:      intPtr(int(reimbursement.Amount)),
				Description: stringPtr(reimbursement.Description),
			})
		}
	}

	// Build response
	response := &v1.PayslipResponse{
		PayrollId: payslip.PayrollID,
		AttendancePeriod: v1.AttendancePeriod{
			StartDate: openapi_types.Date{Time: attendancePeriod.StartDate},
			EndDate:   openapi_types.Date{Time: attendancePeriod.EndDate},
		},
		EmployeeId:          payslip.EmployeeID,
		BaseSalary:          payslip.BaseSalary,
		AttendanceCount:     payslip.AttendanceCount,
		TotalWorkingDays:    payslip.TotalWorkingDays,
		ProratedSalary:      payslip.ProratedSalary,
		OvertimeTotalHours:  payslip.OvertimeTotalHours,
		OvertimePayment:     payslip.OvertimeTotalPay,
		Reimbursements:      reimbursementItems,
		ReimbursementsTotal: payslip.ReimbursementTotal,
		TotalTakeHome:       payslip.TotalTakeHome,
	}

	return response, nil
}

// Helper functions to convert values to pointers
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func datePtr(t time.Time) *openapi_types.Date {
	date := openapi_types.Date{Time: t}
	return &date
}
