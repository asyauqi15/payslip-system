package middleware

import (
	"context"
	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/spf13/cast"
	"net/http"
	"strings"
)

func CheckIPAddress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string

		xff := r.Header.Get("X-Forwarded-For")
		xRealIP := r.Header.Get("X-Real-IP")

		if xff != "" {
			ips := strings.Split(xff, ",")
			ip = strings.TrimSpace(ips[0])
		} else if xRealIP != "" {
			ip = xRealIP
		} else {
			ip = r.RemoteAddr
			if colon := strings.LastIndex(ip, ":"); colon != -1 {
				ip = ip[:colon]
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, constant.ContextKeyIPAddress, cast.ToString(ip))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
