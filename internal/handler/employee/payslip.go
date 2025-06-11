package employee

import (
	"net/http"
	"strconv"

	"github.com/asyauqi15/payslip-system/pkg/logger"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *HandlerImpl) GetPayslip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payroll ID from URL path
	payrollIDStr := chi.URLParam(r, "id")
	payrollID, err := strconv.ParseInt(payrollIDStr, 10, 64)
	if err != nil {
		logger.Error(ctx, "invalid payroll ID", "id", payrollIDStr, "error", err)
		resp := &v1.DefaultErrorResponse{}
		resp.Error.Message = "invalid payroll ID"
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp)
		return
	}

	payslipResponse, err := h.payslipUsecase.GetPayslip(ctx, payrollID)
	if err != nil {
		logger.Error(ctx, "failed to get payslip", "payroll_id", payrollID, "error", err)
		resp := &v1.DefaultErrorResponse{}
		resp.Error.Message = err.Error()

		if httpErr, ok := err.(interface{ HTTPStatus() int }); ok {
			render.Status(r, httpErr.HTTPStatus())
		} else {
			render.Status(r, http.StatusInternalServerError)
		}

		render.JSON(w, r, resp)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, payslipResponse)
}
