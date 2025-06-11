package payroll

import (
	"context"

	"github.com/asyauqi15/payslip-system/internal/repository"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/asyauqi15/payslip-system/internal/usecase/payroll Usecase
type Usecase interface {
	RunPayroll(ctx context.Context, req v1.PostAdminPayrollsJSONRequestBody) error
	GetPayrollSummary(ctx context.Context, payrollID int64) (*v1.AdminPayrollSummaryResponse, error)
}

type UsecaseImpl struct {
	payrollRepo          repository.PayrollRepository
	payslipRepo          repository.PayslipRepository
	employeeRepo         repository.EmployeeRepository
	attendanceRepo       repository.AttendanceRepository
	attendancePeriodRepo repository.AttendancePeriodRepository
	overtimeRepo         repository.OvertimeRepository
	reimbursementRepo    repository.ReimbursementRepository
	userRepo             repository.UserRepository
}

func NewUsecase(
	payrollRepo repository.PayrollRepository,
	payslipRepo repository.PayslipRepository,
	employeeRepo repository.EmployeeRepository,
	attendanceRepo repository.AttendanceRepository,
	attendancePeriodRepo repository.AttendancePeriodRepository,
	overtimeRepo repository.OvertimeRepository,
	reimbursementRepo repository.ReimbursementRepository,
	userRepo repository.UserRepository,
) Usecase {
	return &UsecaseImpl{
		payrollRepo:          payrollRepo,
		payslipRepo:          payslipRepo,
		employeeRepo:         employeeRepo,
		attendanceRepo:       attendanceRepo,
		attendancePeriodRepo: attendancePeriodRepo,
		overtimeRepo:         overtimeRepo,
		reimbursementRepo:    reimbursementRepo,
		userRepo:             userRepo,
	}
}
