package payslip

import (
	"context"

	"github.com/asyauqi15/payslip-system/internal/repository"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/asyauqi15/payslip-system/internal/usecase/payslip Usecase
type Usecase interface {
	GetPayslip(ctx context.Context, payrollID int64) (*v1.PayslipResponse, error)
}

type UsecaseImpl struct {
	payslipRepo          repository.PayslipRepository
	payrollRepo          repository.PayrollRepository
	employeeRepo         repository.EmployeeRepository
	reimbursementRepo    repository.ReimbursementRepository
	attendancePeriodRepo repository.AttendancePeriodRepository
}

func NewUsecase(
	payslipRepo repository.PayslipRepository,
	payrollRepo repository.PayrollRepository,
	employeeRepo repository.EmployeeRepository,
	reimbursementRepo repository.ReimbursementRepository,
	attendancePeriodRepo repository.AttendancePeriodRepository,
) Usecase {
	return &UsecaseImpl{
		payslipRepo:          payslipRepo,
		payrollRepo:          payrollRepo,
		employeeRepo:         employeeRepo,
		reimbursementRepo:    reimbursementRepo,
		attendancePeriodRepo: attendancePeriodRepo,
	}
}
