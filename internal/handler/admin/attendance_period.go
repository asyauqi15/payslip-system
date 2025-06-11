package admin

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/render"
)

func (h *HandlerImpl) CreateAttendancePeriod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req v1.AttendancePeriodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		resp := &v1.DefaultErrorResponse{}
		resp.Error.Message = "invalid request payload"
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp)
		return
	}

	startDate := time.Time(req.StartDate.Time)
	endDate := time.Time(req.EndDate.Time)

	_, err := h.attendancePeriodUsecase.CreateAttendancePeriod(ctx, startDate, endDate)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create attendance period", "error", err)
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
