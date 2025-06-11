package usecase

import (
	"github.com/asyauqi15/payslip-system/internal/repository"
	"github.com/asyauqi15/payslip-system/internal/usecase/attendance"
	"github.com/asyauqi15/payslip-system/internal/usecase/attendance_period"
	authusecase "github.com/asyauqi15/payslip-system/internal/usecase/auth"
	"github.com/asyauqi15/payslip-system/internal/usecase/overtime"
	jwtauth "github.com/asyauqi15/payslip-system/pkg/jwt-auth"
)

type Registry struct {
	Auth                   authusecase.Usecase
	CreateAttendancePeriod attendance_period.CreateAttendancePeriodUsecase
	SubmitAttendance       attendance.SubmitAttendanceUsecase
	SubmitOvertime         overtime.SubmitOvertimeUsecase
}

func InitializeUseCase(repository *repository.Registry, jwt *jwtauth.JWTAuthentication) *Registry {
	return &Registry{
		Auth:                   authusecase.NewUsecase(repository.UserRepository, jwt),
		CreateAttendancePeriod: attendance_period.NewCreateAttendancePeriodUsecase(repository.AttendancePeriodRepository),
		SubmitAttendance:       attendance.NewSubmitAttendanceUsecase(repository.AttendanceRepository, repository.EmployeeRepository),
		SubmitOvertime:         overtime.NewSubmitOvertimeUsecase(repository.OvertimeRepository, repository.EmployeeRepository, repository.AttendanceRepository),
	}
}
