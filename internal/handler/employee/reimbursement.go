package employee

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/asyauqi15/payslip-system/internal/usecase/reimbursement"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/render"
)

type ReimbursementHandler interface {
	SubmitReimbursement(w http.ResponseWriter, r *http.Request)
}

type ReimbursementHandlerImpl struct {
	submitReimbursementUsecase reimbursement.SubmitReimbursementUsecase
}

func NewReimbursementHandler(submitReimbursementUsecase reimbursement.SubmitReimbursementUsecase) ReimbursementHandler {
	return &ReimbursementHandlerImpl{
		submitReimbursementUsecase: submitReimbursementUsecase,
	}
}

func (h *ReimbursementHandlerImpl) SubmitReimbursement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req v1.PostEmployeeReimbursementJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		resp := &v1.DefaultErrorResponse{}
		resp.Error.Message = "invalid request payload"
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp)
		return
	}

	err := h.submitReimbursementUsecase.SubmitReimbursement(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to submit reimbursement", "error", err)
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
