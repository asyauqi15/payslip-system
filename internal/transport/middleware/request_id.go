package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/asyauqi15/payslip-system/internal/constant"
)

// RequestID middleware generates or extracts a unique request ID for each request
// and adds it to the request context for tracing purposes across services
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestID string

		// Check if request ID is already present in headers (for distributed tracing)
		if existingID := r.Header.Get("X-Request-ID"); existingID != "" {
			requestID = existingID
		} else {
			// Generate a new request ID
			requestID = generateRequestID()
		}

		// Add request ID to response headers for client visibility
		w.Header().Set("X-Request-ID", requestID)

		// Add request ID to request context
		ctx := context.WithValue(r.Context(), constant.ContextKeyRequestID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// generateRequestID creates a new random request ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a simple timestamp-based ID if crypto/rand fails
		return "req_fallback"
	}
	return hex.EncodeToString(bytes)
}
