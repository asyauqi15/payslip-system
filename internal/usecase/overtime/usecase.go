package overtime

import (
	"context"
	"github.com/asyauqi15/payslip-system/internal/repository"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
)

type Usecase interface {
	SubmitOvertime(ctx context.Context, req v1.OvertimeRequest) error
}

type UsecaseImpl struct {
	overtimeRepo   repository.OvertimeRepository
	employeeRepo   repository.EmployeeRepository
	attendanceRepo repository.AttendanceRepository
}

func NewUsecase(
	overtimeRepo repository.OvertimeRepository,
	employeeRepo repository.EmployeeRepository,
	attendanceRepo repository.AttendanceRepository,
) Usecase {
	return &UsecaseImpl{
		overtimeRepo:   overtimeRepo,
		employeeRepo:   employeeRepo,
		attendanceRepo: attendanceRepo,
	}
}
