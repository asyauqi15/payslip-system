package employee

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/asyauqi15/payslip-system/internal/usecase/overtime"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/render"
)

type OvertimeHandler interface {
	SubmitOvertime(w http.ResponseWriter, r *http.Request)
}

type OvertimeHandlerImpl struct {
	submitOvertimeUsecase overtime.SubmitOvertimeUsecase
}

func NewOvertimeHandler(submitOvertimeUsecase overtime.SubmitOvertimeUsecase) OvertimeHandler {
	return &OvertimeHandlerImpl{
		submitOvertimeUsecase: submitOvertimeUsecase,
	}
}

func (h *OvertimeHandlerImpl) SubmitOvertime(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req v1.PostEmployeeOvertimeJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		resp := &v1.DefaultErrorResponse{}
		resp.Error.Message = "invalid request payload"
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp)
		return
	}

	err := h.submitOvertimeUsecase.SubmitOvertime(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to submit overtime", "error", err)
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
