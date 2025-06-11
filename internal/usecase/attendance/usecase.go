package attendance

import (
	"context"

	"github.com/asyauqi15/payslip-system/internal/repository"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
)

type Usecase interface {
	SubmitAttendance(ctx context.Context, attendanceType v1.PostEmployeeAttendanceJSONBodyAttendanceType) error
}

type UsecaseImpl struct {
	attendanceRepo repository.AttendanceRepository
	employeeRepo   repository.EmployeeRepository
}

func NewUsecase(
	attendanceRepo repository.AttendanceRepository,
	employeeRepo repository.EmployeeRepository,
) Usecase {
	return &UsecaseImpl{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
	}
}
