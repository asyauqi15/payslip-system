package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asyauqi15/payslip-system/internal/handler/auth"
	authpkg "github.com/asyauqi15/payslip-system/internal/usecase/auth"
	"github.com/asyauqi15/payslip-system/internal/usecase/auth/mock"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mock.NewMockUsecase(ctrl)
	handler := auth.NewHandler(mockAuthUsecase)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful login",
			requestBody: v1.AuthRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMock: func() {
				authResult := &authpkg.Result{
					AccessToken:  "access_token_123",
					RefreshToken: "refresh_token_123",
				}
				mockAuthUsecase.EXPECT().
					Auth(gomock.Any(), "testuser", "password123").
					Return(authResult, nil)
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "invalid credentials",
			requestBody: v1.AuthRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			setupMock: func() {
				mockAuthUsecase.EXPECT().
					Auth(gomock.Any(), "testuser", "wrongpassword").
					Return(nil, httppkg.NewUnauthorizedError("username or password is incorrect"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name: "user not found",
			requestBody: v1.AuthRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			setupMock: func() {
				mockAuthUsecase.EXPECT().
					Auth(gomock.Any(), "nonexistent", "password123").
					Return(nil, httppkg.NewNotFoundError("username or password is incorrect"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			setupMock:      func() {}, // No usecase call expected
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "empty username",
			requestBody: v1.AuthRequest{
				Username: "",
				Password: "password123",
			},
			setupMock: func() {
				mockAuthUsecase.EXPECT().
					Auth(gomock.Any(), "", "password123").
					Return(nil, httppkg.NewBadRequestError("username is required"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name: "empty password",
			requestBody: v1.AuthRequest{
				Username: "testuser",
				Password: "",
			},
			setupMock: func() {
				mockAuthUsecase.EXPECT().
					Auth(gomock.Any(), "testuser", "").
					Return(nil, httppkg.NewBadRequestError("password is required"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Create request body
			var requestBody []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				requestBody = []byte(str)
			} else {
				requestBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal("Failed to marshal request body:", err)
				}
			}

			// Create HTTP request
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(context.Background())

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.Login(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				// For successful requests, verify response structure
				if w.Code == http.StatusOK {
					var response v1.AuthResponse
					err := json.Unmarshal(w.Body.Bytes(), &response)
					if err != nil {
						t.Error("Failed to unmarshal response:", err)
					}

					if response.AccessToken == "" {
						t.Error("Expected access token in response")
					}
					if response.RefreshToken == "" {
						t.Error("Expected refresh token in response")
					}

					// Check content type
					contentType := w.Header().Get("Content-Type")
					if contentType != "application/json" {
						t.Errorf("Expected Content-Type application/json but got %s", contentType)
					}
				}
			}
		})
	}
}
