package employee_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestEmployeeHandler_SubmitAttendance(t *testing.T) {
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
			name: "successful check-in",
			requestBody: v1.PostEmployeeAttendanceJSONBody{
				AttendanceType: v1.CheckIn,
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				mockAttendanceUsecase.EXPECT().
					SubmitAttendance(gomock.Any(), v1.CheckIn).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "successful check-out",
			requestBody: v1.PostEmployeeAttendanceJSONBody{
				AttendanceType: v1.CheckOut,
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				mockAttendanceUsecase.EXPECT().
					SubmitAttendance(gomock.Any(), v1.CheckOut).
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
			name: "weekend attendance error",
			requestBody: v1.PostEmployeeAttendanceJSONBody{
				AttendanceType: v1.CheckIn,
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				mockAttendanceUsecase.EXPECT().
					SubmitAttendance(gomock.Any(), v1.CheckIn).
					Return(httppkg.NewBadRequestError("attendance not allowed on weekends"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "already checked in error",
			requestBody: v1.PostEmployeeAttendanceJSONBody{
				AttendanceType: v1.CheckIn,
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				mockAttendanceUsecase.EXPECT().
					SubmitAttendance(gomock.Any(), v1.CheckIn).
					Return(httppkg.NewBadRequestError("already checked in today"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "unauthorized user",
			requestBody: v1.PostEmployeeAttendanceJSONBody{
				AttendanceType: v1.CheckIn,
			},
			setupContext: func() context.Context {
				return context.Background()
			},
			setupMock: func() {
				mockAttendanceUsecase.EXPECT().
					SubmitAttendance(gomock.Any(), v1.CheckIn).
					Return(httppkg.NewUnauthorizedError("user not authenticated"))
			},
			expectedStatus: http.StatusUnauthorized,
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

			req := httptest.NewRequest(http.MethodPost, "/employee/attendance", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(tt.setupContext())

			w := httptest.NewRecorder()
			handler.SubmitAttendance(w, req)

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
