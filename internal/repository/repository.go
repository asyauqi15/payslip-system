package repository

import (
	"github.com/asyauqi15/payslip-system/internal/entity"
	"gorm.io/gorm"
)

type Registry struct {
	UserRepository             UserRepository
	EmployeeRepository         EmployeeRepository
	AttendanceRepository       AttendanceRepository
	AttendancePeriodRepository AttendancePeriodRepository
	OvertimeRepository         OvertimeRepository
	PayrollRepository          PayrollRepository
	PayslipRepository          PayslipRepository
	ReimbursementRepository    ReimbursementRepository
}

func InitializeRepository(db *gorm.DB) *Registry {
	return &Registry{
		UserRepository:             NewUserRepository(&BaseRepositoryImpl[entity.User]{DB: db}),
		EmployeeRepository:         NewEmployeeRepository(&BaseRepositoryImpl[entity.Employee]{DB: db}),
		AttendanceRepository:       NewAttendanceRepository(&BaseRepositoryImpl[entity.Attendance]{DB: db}),
		AttendancePeriodRepository: NewAttendancePeriodRepository(&BaseRepositoryImpl[entity.AttendancePeriod]{DB: db}),
		OvertimeRepository:         NewOvertimeRepository(&BaseRepositoryImpl[entity.Overtime]{DB: db}),
		PayrollRepository:          NewPayrollRepository(&BaseRepositoryImpl[entity.Payroll]{DB: db}),
		PayslipRepository:          NewPayslipRepository(&BaseRepositoryImpl[entity.Payslip]{DB: db}),
		ReimbursementRepository:    NewReimbursementRepository(&BaseRepositoryImpl[entity.Reimbursement]{DB: db}),
	}
}
