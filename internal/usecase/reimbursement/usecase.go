package reimbursement

import (
	"context"

	"github.com/asyauqi15/payslip-system/internal/repository"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
)

type Usecase interface {
	SubmitReimbursement(ctx context.Context, req v1.ReimbursementRequest) error
}

type UsecaseImpl struct {
	reimbursementRepo repository.ReimbursementRepository
	employeeRepo      repository.EmployeeRepository
}

func NewUsecase(
	reimbursementRepo repository.ReimbursementRepository,
	employeeRepo repository.EmployeeRepository,
) Usecase {
	return &UsecaseImpl{
		reimbursementRepo: reimbursementRepo,
		employeeRepo:      employeeRepo,
	}
}
