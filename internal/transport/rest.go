package transport

import (
	"context"
	"fmt"
	"github.com/asyauqi15/payslip-system/internal"
	"github.com/asyauqi15/payslip-system/internal/handler"
	"github.com/asyauqi15/payslip-system/internal/transport/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type RESTServer struct {
	srv *http.Server
}

func NewRESTServer(
	config internal.Config,
	h *handler.Registry,
) (*RESTServer, error) {
	routes := chi.NewRouter()
	routes.Use(middleware.Recoverer)
	routes.Use(middleware.CheckIPAddress)

	routes.Post("/auth/login", h.Auth.Login)

	routes.Get("/ping", pingHandler)
	swaggerRoutes(routes)

	return &RESTServer{
		srv: &http.Server{
			Addr:              fmt.Sprintf("0.0.0.0:%d", config.HTTPServer.Port),
			Handler:           routes,
			ReadTimeout:       config.HTTPServer.ReadTimeout,
			ReadHeaderTimeout: config.HTTPServer.ReadHeaderTimeout,
			IdleTimeout:       config.HTTPServer.IdleTimeout,
			WriteTimeout:      config.HTTPServer.WriteTimeout,
		},
	}, nil
}

func (r *RESTServer) Start() error {
	return r.srv.ListenAndServe()
}

func (r *RESTServer) Stop(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}
