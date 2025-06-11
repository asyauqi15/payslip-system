package payroll

import (
	"context"

	"github.com/asyauqi15/payslip-system/internal/entity"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	"github.com/asyauqi15/payslip-system/pkg/logger"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (u *UsecaseImpl) GetPayrollSummary(ctx context.Context, payrollID int64) (*v1.AdminPayrollSummaryResponse, error) {
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

	// Get all payslips for this payroll
	payslips, err := u.payslipRepo.FindByTemplate(ctx, &entity.Payslip{
		PayrollID: payrollID,
	}, nil)
	if err != nil {
		logger.Error(ctx, "failed to find payslips", "payroll_id", payrollID, "error", err)
		return nil, httppkg.NewInternalServerError("failed to find payslips")
	}

	// Build payslip items with employee information
	payslipItems := make([]v1.PayslipItem, 0, len(payslips))
	for _, payslip := range payslips {
		// Get employee
		employee, err := u.employeeRepo.FindByID(ctx, uint(payslip.EmployeeID), nil)
		if err != nil {
			logger.Error(ctx, "failed to find employee", "employee_id", payslip.EmployeeID, "error", err)
			continue // Skip this payslip if employee not found
		}
		if employee == nil {
			logger.Warn(ctx, "employee not found for payslip", "employee_id", payslip.EmployeeID)
			continue
		}

		// Get user for username
		user, err := u.userRepo.FindByID(ctx, uint(employee.UserID), nil)
		if err != nil {
			logger.Error(ctx, "failed to find user", "user_id", employee.UserID, "error", err)
			continue
		}
		if user == nil {
			logger.Warn(ctx, "user not found for employee", "user_id", employee.UserID)
			continue
		}

		payslipItem := v1.PayslipItem{
			EmployeeId:            payslip.EmployeeID,
			Username:              user.Username,
			BaseSalary:            payslip.BaseSalary,
			AttendanceCount:       payslip.AttendanceCount,
			OvertimeCount:         payslip.OvertimeTotalHours,
			ProratedSalary:        payslip.ProratedSalary,
			OvertimePayment:       payslip.OvertimeTotalPay,
			ReimbursementsPayment: payslip.ReimbursementTotal,
			TotalPay:              payslip.TotalTakeHome,
		}

		payslipItems = append(payslipItems, payslipItem)
	}

	// Build response
	response := &v1.AdminPayrollSummaryResponse{
		PayrollId: payroll.ID,
		AttendancePeriod: v1.AttendancePeriod{
			StartDate: openapi_types.Date{Time: attendancePeriod.StartDate},
			EndDate:   openapi_types.Date{Time: attendancePeriod.EndDate},
		},
		EmployeesCount:         payroll.TotalEmployees,
		TotalPayroll:           payroll.TotalPayroll,
		TotalReimbursementsPay: payroll.TotalReimbursement,
		TotalOvertimePay:       payroll.TotalOvertime,
		PayslipList:            payslipItems,
	}

	return response, nil
}
