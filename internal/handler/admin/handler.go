package admin

import (
	"net/http"

	"github.com/asyauqi15/payslip-system/internal/usecase/attendance_period"
	"github.com/asyauqi15/payslip-system/internal/usecase/payroll"
)

type Handler interface {
	CreateAttendancePeriod(w http.ResponseWriter, r *http.Request)
	RunPayroll(w http.ResponseWriter, r *http.Request)
	GetPayrollSummary(w http.ResponseWriter, r *http.Request)
}

type HandlerImpl struct {
	attendancePeriodUsecase attendance_period.Usecase
	payrollUsecase          payroll.Usecase
}

func NewHandler(
	attendancePeriodUsecase attendance_period.Usecase,
	payrollUsecase payroll.Usecase,
) Handler {
	return &HandlerImpl{
		attendancePeriodUsecase: attendancePeriodUsecase,
		payrollUsecase:          payrollUsecase,
	}
}
