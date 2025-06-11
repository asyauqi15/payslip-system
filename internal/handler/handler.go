package handler

import (
	"github.com/asyauqi15/payslip-system/internal/handler/admin"
	"github.com/asyauqi15/payslip-system/internal/handler/auth"
	"github.com/asyauqi15/payslip-system/internal/usecase"
)

type Registry struct {
	Auth             auth.Handler
	AttendancePeriod admin.AttendancePeriodHandler
}

func InitializeHandler(usecase *usecase.Registry) *Registry {
	return &Registry{
		Auth:             auth.NewHandler(usecase.Auth),
		AttendancePeriod: admin.NewAttendancePeriodHandler(usecase.CreateAttendancePeriod),
	}
}
