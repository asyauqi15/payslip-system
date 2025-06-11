package middleware

import (
	"net/http"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/entity"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/render"
)

func RequireAdminRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value(constant.ContextKeyUserRole)
		if userRole == nil || userRole.(string) != entity.UserRoleAdmin {
			resp := &v1.DefaultErrorResponse{}
			resp.Error.Message = "access denied: admin role required"
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}
