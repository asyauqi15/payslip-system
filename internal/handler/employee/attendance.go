package employee

import (
	"encoding/json"
	"net/http"

	"github.com/asyauqi15/payslip-system/pkg/logger"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/render"
)

func (h *HandlerImpl) SubmitAttendance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req v1.PostEmployeeAttendanceJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "failed to decode request body", "error", err)
		resp := &v1.DefaultErrorResponse{}
		resp.Error.Message = "invalid request payload"
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp)
		return
	}

	err := h.attendanceUsecase.SubmitAttendance(ctx, req.AttendanceType)
	if err != nil {
		logger.Error(ctx, "failed to submit attendance", "error", err)
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
