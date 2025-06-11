package usecase

import (
	"github.com/asyauqi15/payslip-system/internal/repository"
	"github.com/asyauqi15/payslip-system/internal/usecase/attendance"
	"github.com/asyauqi15/payslip-system/internal/usecase/attendance_period"
	authusecase "github.com/asyauqi15/payslip-system/internal/usecase/auth"
	"github.com/asyauqi15/payslip-system/internal/usecase/overtime"
	"github.com/asyauqi15/payslip-system/internal/usecase/payroll"
	"github.com/asyauqi15/payslip-system/internal/usecase/payslip"
	"github.com/asyauqi15/payslip-system/internal/usecase/reimbursement"
	jwtauth "github.com/asyauqi15/payslip-system/pkg/jwt-auth"
)

type Registry struct {
	Auth                   authusecase.Usecase
	CreateAttendancePeriod attendance_period.Usecase
	SubmitAttendance       attendance.Usecase
	SubmitOvertime         overtime.Usecase
	SubmitReimbursement    reimbursement.Usecase
	PayrollUsecase         payroll.Usecase
	GetPayslip             payslip.Usecase
}

func InitializeUseCase(repository *repository.Registry, jwt *jwtauth.JWTAuthentication) *Registry {
	return &Registry{
		Auth:                   authusecase.NewUsecase(repository.UserRepository, jwt),
		CreateAttendancePeriod: attendance_period.NewUsecase(repository.AttendancePeriodRepository),
		SubmitAttendance:       attendance.NewUsecase(repository.AttendanceRepository, repository.EmployeeRepository),
		SubmitOvertime:         overtime.NewUsecase(repository.OvertimeRepository, repository.EmployeeRepository, repository.AttendanceRepository),
		SubmitReimbursement:    reimbursement.NewUsecase(repository.ReimbursementRepository, repository.EmployeeRepository),
		PayrollUsecase:         payroll.NewUsecase(repository.PayrollRepository, repository.PayslipRepository, repository.EmployeeRepository, repository.AttendanceRepository, repository.AttendancePeriodRepository, repository.OvertimeRepository, repository.ReimbursementRepository, repository.UserRepository),
		GetPayslip:             payslip.NewUsecase(repository.PayslipRepository, repository.PayrollRepository, repository.EmployeeRepository, repository.ReimbursementRepository, repository.AttendancePeriodRepository),
	}
}
