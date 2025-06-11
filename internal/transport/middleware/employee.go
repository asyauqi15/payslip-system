package middleware

import (
	"net/http"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/entity"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/render"
)

func RequireEmployeeRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value(constant.ContextKeyUserRole)
		if userRole == nil {
			resp := &v1.DefaultErrorResponse{}
			resp.Error.Message = "access denied: authentication required"
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, resp)
			return
		}

		// Allow both admin and non-admin users to access employee endpoints
		// But you can also restrict it to only non-admin users if needed
		role := userRole.(string)
		if role != entity.UserRoleAdmin && role != entity.UserRoleDefault {
			resp := &v1.DefaultErrorResponse{}
			resp.Error.Message = "access denied: employee role required"
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}
