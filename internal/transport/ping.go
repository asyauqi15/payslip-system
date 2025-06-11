package transport

import (
	"net/http"

	"github.com/asyauqi15/payslip-system/pkg/logger"
	"github.com/go-chi/render"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Debug(ctx, "ping endpoint accessed")

	render.Status(r, http.StatusOK)
	render.JSON(w, r, struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	})
}
