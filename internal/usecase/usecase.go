package usecase

import (
	"github.com/asyauqi15/payslip-system/internal/repository"
	"github.com/asyauqi15/payslip-system/internal/usecase/attendance"
	"github.com/asyauqi15/payslip-system/internal/usecase/attendance_period"
	authusecase "github.com/asyauqi15/payslip-system/internal/usecase/auth"
	"github.com/asyauqi15/payslip-system/internal/usecase/overtime"
	"github.com/asyauqi15/payslip-system/internal/usecase/payroll"
	"github.com/asyauqi15/payslip-system/internal/usecase/reimbursement"
	jwtauth "github.com/asyauqi15/payslip-system/pkg/jwt-auth"
)

type Registry struct {
	Auth                   authusecase.Usecase
	CreateAttendancePeriod attendance_period.CreateAttendancePeriodUsecase
	SubmitAttendance       attendance.SubmitAttendanceUsecase
	SubmitOvertime         overtime.SubmitOvertimeUsecase
	SubmitReimbursement    reimbursement.SubmitReimbursementUsecase
	RunPayroll             payroll.RunPayrollUsecase
}

func InitializeUseCase(repository *repository.Registry, jwt *jwtauth.JWTAuthentication) *Registry {
	return &Registry{
		Auth:                   authusecase.NewUsecase(repository.UserRepository, jwt),
		CreateAttendancePeriod: attendance_period.NewCreateAttendancePeriodUsecase(repository.AttendancePeriodRepository),
		SubmitAttendance:       attendance.NewSubmitAttendanceUsecase(repository.AttendanceRepository, repository.EmployeeRepository),
		SubmitOvertime:         overtime.NewSubmitOvertimeUsecase(repository.OvertimeRepository, repository.EmployeeRepository, repository.AttendanceRepository),
		SubmitReimbursement:    reimbursement.NewSubmitReimbursementUsecase(repository.ReimbursementRepository, repository.EmployeeRepository),
		RunPayroll:             payroll.NewRunPayrollUsecase(repository.PayrollRepository, repository.PayslipRepository, repository.EmployeeRepository, repository.AttendanceRepository, repository.AttendancePeriodRepository, repository.OvertimeRepository, repository.ReimbursementRepository),
	}
}
