package transport

import (
	"github.com/go-chi/render"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	})
}
