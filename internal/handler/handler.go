package handler

import (
	"github.com/asyauqi15/payslip-system/internal/handler/admin"
	"github.com/asyauqi15/payslip-system/internal/handler/auth"
	"github.com/asyauqi15/payslip-system/internal/handler/employee"
	"github.com/asyauqi15/payslip-system/internal/usecase"
)

type Registry struct {
	Auth     auth.Handler
	Employee employee.Handler
	Admin    admin.Handler
}

func InitializeHandler(usecase *usecase.Registry) *Registry {
	return &Registry{
		Auth:     auth.NewHandler(usecase.Auth),
		Employee: employee.NewHandler(usecase.SubmitAttendance, usecase.SubmitOvertime, usecase.GetPayslip, usecase.SubmitReimbursement),
		Admin:    admin.NewHandler(usecase.CreateAttendancePeriod, usecase.PayrollUsecase),
	}
}
