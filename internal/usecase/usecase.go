package usecase

import (
	"github.com/asyauqi15/payslip-system/internal/repository"
	"github.com/asyauqi15/payslip-system/internal/usecase/attendance"
	"github.com/asyauqi15/payslip-system/internal/usecase/attendance_period"
	authusecase "github.com/asyauqi15/payslip-system/internal/usecase/auth"
	jwtauth "github.com/asyauqi15/payslip-system/pkg/jwt-auth"
)

type Registry struct {
	Auth                   authusecase.Usecase
	CreateAttendancePeriod attendance_period.CreateAttendancePeriodUsecase
	SubmitAttendance       attendance.SubmitAttendanceUsecase
}

func InitializeUseCase(repository *repository.Registry, jwt *jwtauth.JWTAuthentication) *Registry {
	return &Registry{
		Auth:                   authusecase.NewUsecase(repository.UserRepository, jwt),
		CreateAttendancePeriod: attendance_period.NewCreateAttendancePeriodUsecase(repository.AttendancePeriodRepository),
		SubmitAttendance:       attendance.NewSubmitAttendanceUsecase(repository.AttendanceRepository, repository.EmployeeRepository),
	}
}
