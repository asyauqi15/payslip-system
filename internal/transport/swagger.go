package transport

import (
	"embed"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/chi/v5"
	"net/http"
)

//go:embed swagger
var swaggerFiles embed.FS

func swaggerRoutes(r *chi.Mux) {
	r.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.FS(swaggerFiles)),
	))

	w, err := v1.GetSwagger()
	if err != nil {
		return
	}
	b, _ := w.MarshalJSON()

	r.Get("/openapi3.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	})
}
