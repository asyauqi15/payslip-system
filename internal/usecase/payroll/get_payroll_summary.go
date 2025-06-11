package payroll

import (
	"context"
	"log/slog"
	"time"

	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type GetPayrollSummaryUsecase interface {
	GetPayrollSummary(ctx context.Context, payrollID int64) (*v1.AdminPayrollSummaryResponse, error)
}

type GetPayrollSummaryUsecaseImpl struct {
	payrollRepo          repository.PayrollRepository
	payslipRepo          repository.PayslipRepository
	employeeRepo         repository.EmployeeRepository
	userRepo             repository.UserRepository
	attendancePeriodRepo repository.AttendancePeriodRepository
}

func NewGetPayrollSummaryUsecase(
	payrollRepo repository.PayrollRepository,
	payslipRepo repository.PayslipRepository,
	employeeRepo repository.EmployeeRepository,
	userRepo repository.UserRepository,
	attendancePeriodRepo repository.AttendancePeriodRepository,
) GetPayrollSummaryUsecase {
	return &GetPayrollSummaryUsecaseImpl{
		payrollRepo:          payrollRepo,
		payslipRepo:          payslipRepo,
		employeeRepo:         employeeRepo,
		userRepo:             userRepo,
		attendancePeriodRepo: attendancePeriodRepo,
	}
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

func (u *GetPayrollSummaryUsecaseImpl) GetPayrollSummary(ctx context.Context, payrollID int64) (*v1.AdminPayrollSummaryResponse, error) {
	// Get payroll record
	payroll, err := u.payrollRepo.FindByID(ctx, uint(payrollID), nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find payroll", "payroll_id", payrollID, "error", err)
		return nil, httppkg.NewInternalServerError("failed to find payroll")
	}
	if payroll == nil {
		return nil, httppkg.NewNotFoundError("payroll not found")
	}

	// Get attendance period
	attendancePeriod, err := u.attendancePeriodRepo.FindByID(ctx, uint(payroll.AttendancePeriodID), nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find attendance period", "attendance_period_id", payroll.AttendancePeriodID, "error", err)
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
		slog.ErrorContext(ctx, "failed to find payslips", "payroll_id", payrollID, "error", err)
		return nil, httppkg.NewInternalServerError("failed to find payslips")
	}

	// Build payslip items with employee information
	payslipItems := make([]v1.PayslipItem, 0, len(payslips))
	for _, payslip := range payslips {
		// Get employee
		employee, err := u.employeeRepo.FindByID(ctx, uint(payslip.EmployeeID), nil)
		if err != nil {
			slog.ErrorContext(ctx, "failed to find employee", "employee_id", payslip.EmployeeID, "error", err)
			continue // Skip this payslip if employee not found
		}
		if employee == nil {
			slog.WarnContext(ctx, "employee not found for payslip", "employee_id", payslip.EmployeeID)
			continue
		}

		// Get user for username
		user, err := u.userRepo.FindByID(ctx, uint(employee.UserID), nil)
		if err != nil {
			slog.ErrorContext(ctx, "failed to find user", "user_id", employee.UserID, "error", err)
			continue
		}
		if user == nil {
			slog.WarnContext(ctx, "user not found for employee", "user_id", employee.UserID)
			continue
		}

		payslipItem := v1.PayslipItem{
			EmployeeId:            intPtr(int(payslip.EmployeeID)),
			Username:              stringPtr(user.Username),
			BaseSalary:            intPtr(int(payslip.BaseSalary)),
			AttendanceCount:       intPtr(payslip.AttendanceCount),
			OvertimeCount:         intPtr(payslip.OvertimeTotalHours),
			ProratedSalary:        intPtr(int(payslip.ProratedSalary)),
			OvertimePayment:       intPtr(int(payslip.OvertimeTotalPay)),
			ReimbursementsPayment: intPtr(int(payslip.ReimbursementTotal)),
			TotalPay:              intPtr(int(payslip.TotalTakeHome)),
		}

		payslipItems = append(payslipItems, payslipItem)
	}

	// Build response
	response := &v1.AdminPayrollSummaryResponse{
		PayrollId: intPtr(int(payroll.ID)),
		AttendancePeriod: &v1.AttendancePeriod{
			StartDate: datePtr(attendancePeriod.StartDate),
			EndDate:   datePtr(attendancePeriod.EndDate),
		},
		EmployeesCount:         intPtr(int(payroll.TotalEmployees)),
		TotalPayroll:           intPtr(int(payroll.TotalPayroll)),
		TotalReimbursementsPay: intPtr(int(payroll.TotalReimbursement)),
		TotalOvertimePay:       intPtr(int(payroll.TotalOvertime)),
		PayslipList:            &payslipItems,
	}

	return response, nil
}
