package transport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asyauqi15/payslip-system/internal"
	"github.com/asyauqi15/payslip-system/internal/handler"
	"github.com/asyauqi15/payslip-system/internal/transport/middleware"
	jwtauth "github.com/asyauqi15/payslip-system/pkg/jwt-auth"
	"github.com/go-chi/chi/v5"
)

type RESTServer struct {
	srv *http.Server
}

func NewRESTServer(
	config internal.Config,
	h *handler.Registry,
	jwt *jwtauth.JWTAuthentication,
) (*RESTServer, error) {
	routes := chi.NewRouter()
	routes.Use(middleware.RequestID)
	routes.Use(middleware.HTTPLogger)
	routes.Use(middleware.Recoverer)
	routes.Use(middleware.CheckIPAddress)

	// Public routes
	routes.Post("/auth/login", h.Auth.Login)
	routes.Get("/ping", pingHandler)
	swaggerRoutes(routes)

	// Admin routes (require authentication and admin role)
	routes.Route("/admin", func(r chi.Router) {
		r.Use(jwt.Authenticator)
		r.Use(middleware.RequireAdminRole)

		r.Post("/attendance-periods", h.Admin.CreateAttendancePeriod)
		r.Post("/payrolls", h.Admin.RunPayroll)
		r.Get("/payrolls/{id}", h.Admin.GetPayrollSummary)
	})

	// Employee routes (require authentication)
	routes.Route("/employee", func(r chi.Router) {
		r.Use(jwt.Authenticator)
		r.Use(middleware.RequireEmployeeRole)

		r.Post("/attendance", h.Employee.SubmitAttendance)
		r.Post("/overtime", h.Employee.SubmitOvertime)
		r.Post("/reimbursement", h.Employee.SubmitReimbursement)
		r.Get("/payroll/{id}", h.Employee.GetPayslip)
	})

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
