package admin

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *HandlerImpl) RunPayroll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req v1.PostAdminPayrollsJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		resp := &v1.DefaultErrorResponse{}
		resp.Error.Message = "invalid request payload"
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp)
		return
	}

	err := h.payrollUsecase.RunPayroll(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to run payroll", "error", err)
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

	w.WriteHeader(http.StatusCreated)
}

func (h *HandlerImpl) GetPayrollSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payroll ID from URL path
	payrollIDStr := chi.URLParam(r, "id")
	payrollID, err := strconv.ParseInt(payrollIDStr, 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "invalid payroll ID", "id", payrollIDStr, "error", err)
		resp := &v1.DefaultErrorResponse{}
		resp.Error.Message = "invalid payroll ID"
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp)
		return
	}

	summary, err := h.payrollUsecase.GetPayrollSummary(ctx, payrollID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get payroll summary", "payroll_id", payrollID, "error", err)
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
	render.JSON(w, r, summary)
}
