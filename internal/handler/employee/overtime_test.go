package employee_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/handler/employee"
	attendancemock "github.com/asyauqi15/payslip-system/internal/usecase/attendance/mock"
	overtimemock "github.com/asyauqi15/payslip-system/internal/usecase/overtime/mock"
	payslipmock "github.com/asyauqi15/payslip-system/internal/usecase/payslip/mock"
	reimbursementmock "github.com/asyauqi15/payslip-system/internal/usecase/reimbursement/mock"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"go.uber.org/mock/gomock"
)

func TestEmployeeHandler_SubmitOvertime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAttendanceUsecase := attendancemock.NewMockUsecase(ctrl)
	mockOvertimeUsecase := overtimemock.NewMockUsecase(ctrl)
	mockPayslipUsecase := payslipmock.NewMockUsecase(ctrl)
	mockReimbursementUsecase := reimbursementmock.NewMockUsecase(ctrl)

	handler := employee.NewHandler(
		mockAttendanceUsecase,
		mockOvertimeUsecase,
		mockPayslipUsecase,
		mockReimbursementUsecase,
	)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupContext   func() context.Context
		setupMock      func()
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful overtime submission",
			requestBody: v1.PostEmployeeOvertimeJSONRequestBody{
				StartTime:   time.Date(2025, 6, 10, 18, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2025, 6, 10, 20, 0, 0, 0, time.UTC),
				Description: "Working on urgent project",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				req := v1.PostEmployeeOvertimeJSONRequestBody{
					StartTime:   time.Date(2025, 6, 10, 18, 0, 0, 0, time.UTC),
					EndTime:     time.Date(2025, 6, 10, 20, 0, 0, 0, time.UTC),
					Description: "Working on urgent project",
				}
				mockOvertimeUsecase.EXPECT().
					SubmitOvertime(gomock.Any(), req).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			setupContext:   func() context.Context { return context.Background() },
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "overtime too long error",
			requestBody: v1.PostEmployeeOvertimeJSONRequestBody{
				StartTime:   time.Date(2025, 6, 10, 18, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2025, 6, 10, 23, 0, 0, 0, time.UTC),
				Description: "Working on urgent project",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				req := v1.PostEmployeeOvertimeJSONRequestBody{
					StartTime:   time.Date(2025, 6, 10, 18, 0, 0, 0, time.UTC),
					EndTime:     time.Date(2025, 6, 10, 23, 0, 0, 0, time.UTC),
					Description: "Working on urgent project",
				}
				mockOvertimeUsecase.EXPECT().
					SubmitOvertime(gomock.Any(), req).
					Return(httppkg.NewBadRequestError("overtime duration cannot exceed 3 hours"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "overlapping overtime error",
			requestBody: v1.PostEmployeeOvertimeJSONRequestBody{
				StartTime:   time.Date(2025, 6, 10, 18, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2025, 6, 10, 20, 0, 0, 0, time.UTC),
				Description: "Overlapping overtime",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				req := v1.PostEmployeeOvertimeJSONRequestBody{
					StartTime:   time.Date(2025, 6, 10, 18, 0, 0, 0, time.UTC),
					EndTime:     time.Date(2025, 6, 10, 20, 0, 0, 0, time.UTC),
					Description: "Overlapping overtime",
				}
				mockOvertimeUsecase.EXPECT().
					SubmitOvertime(gomock.Any(), req).
					Return(httppkg.NewBadRequestError("overtime conflicts with existing record"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

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

			req := httptest.NewRequest(http.MethodPost, "/employee/overtime", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(tt.setupContext())

			w := httptest.NewRecorder()
			handler.SubmitOvertime(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError {
				var errorResp v1.DefaultErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResp)
				if err != nil {
					t.Error("Failed to unmarshal error response:", err)
				}
				if errorResp.Error.Message == "" {
					t.Error("Expected error message in response")
				}
			}
		})
	}
}
