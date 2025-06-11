package employee

import (
	"net/http"

	"github.com/asyauqi15/payslip-system/internal/usecase/attendance"
	"github.com/asyauqi15/payslip-system/internal/usecase/overtime"
	"github.com/asyauqi15/payslip-system/internal/usecase/payslip"
	"github.com/asyauqi15/payslip-system/internal/usecase/reimbursement"
)

type Handler interface {
	SubmitAttendance(w http.ResponseWriter, r *http.Request)
	SubmitOvertime(w http.ResponseWriter, r *http.Request)
	GetPayslip(w http.ResponseWriter, r *http.Request)
	SubmitReimbursement(w http.ResponseWriter, r *http.Request)
}

type HandlerImpl struct {
	attendanceUsecase    attendance.Usecase
	overtimeUsecase      overtime.Usecase
	payslipUsecase       payslip.Usecase
	reimbursementUsecase reimbursement.Usecase
}

func NewHandler(
	attendanceUsecase attendance.Usecase,
	overtimeUsecase overtime.Usecase,
	payslipUsecase payslip.Usecase,
	reimbursementUsecase reimbursement.Usecase,
) Handler {
	return &HandlerImpl{
		attendanceUsecase:    attendanceUsecase,
		overtimeUsecase:      overtimeUsecase,
		payslipUsecase:       payslipUsecase,
		reimbursementUsecase: reimbursementUsecase,
	}
}
