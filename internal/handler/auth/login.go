package auth

import (
	"encoding/json"
	"net/http"

	"github.com/asyauqi15/payslip-system/pkg/logger"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/render"
)

func (h *HandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req v1.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "failed to decode login request", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	logger.Info(ctx, "user login attempt", "username", req.Username)

	result, err := h.authUsecase.Auth(ctx, req.Username, req.Password)
	if err != nil {
		logger.Error(ctx, "login failed", "username", req.Username, "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logger.Info(ctx, "login successful", "username", req.Username)

	response := v1.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response)
}
