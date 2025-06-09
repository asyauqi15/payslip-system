package middleware

import (
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"net/http"

	"github.com/go-chi/render"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler { //nolint: errorlint,goerr113
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				resp := &v1.DefaultErrorResponse{}
				resp.Error.Message = "internal server error"
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, resp)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
