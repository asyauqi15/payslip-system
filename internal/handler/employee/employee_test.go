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
	"github.com/go-chi/chi/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
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
				return context.Background() // No user ID in context
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
			req := httptest.NewRequest(http.MethodPost, "/employee/attendance", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(tt.setupContext())

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.SubmitAttendance(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError {
				// Verify error response structure
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
			name: "invalid time format",
			requestBody: v1.PostEmployeeOvertimeJSONRequestBody{
				StartTime:   time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC), // invalid
				EndTime:     time.Date(2025, 6, 10, 20, 0, 0, 0, time.UTC),
				Description: "Invalid time",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				req := v1.PostEmployeeOvertimeJSONRequestBody{
					StartTime:   time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC),
					EndTime:     time.Date(2025, 6, 10, 20, 0, 0, 0, time.UTC),
					Description: "Invalid time",
				}
				mockOvertimeUsecase.EXPECT().
					SubmitOvertime(gomock.Any(), req).
					Return(httppkg.NewBadRequestError("invalid time format"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "weekend overtime error",
			requestBody: v1.PostEmployeeOvertimeJSONRequestBody{
				StartTime:   time.Date(2025, 6, 14, 18, 0, 0, 0, time.UTC), // Saturday
				EndTime:     time.Date(2025, 6, 14, 20, 0, 0, 0, time.UTC),
				Description: "Weekend overtime",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				req := v1.PostEmployeeOvertimeJSONRequestBody{
					StartTime:   time.Date(2025, 6, 14, 18, 0, 0, 0, time.UTC),
					EndTime:     time.Date(2025, 6, 14, 20, 0, 0, 0, time.UTC),
					Description: "Weekend overtime",
				}
				mockOvertimeUsecase.EXPECT().
					SubmitOvertime(gomock.Any(), req).
					Return(httppkg.NewBadRequestError("overtime not allowed on weekends"))
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
			req := httptest.NewRequest(http.MethodPost, "/employee/overtime", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(tt.setupContext())

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.SubmitOvertime(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError {
				// Verify error response structure
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

func TestEmployeeHandler_GetPayslip(t *testing.T) {
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
		payrollID      string
		setupContext   func() context.Context
		setupMock      func()
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "successful payslip retrieval",
			payrollID: "1",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				payslipResponse := &v1.PayslipResponse{
					PayrollId:           int64(1),
					EmployeeId:          int64(1),
					BaseSalary:          5000000,
					AttendanceCount:     20,
					TotalWorkingDays:    22,
					ProratedSalary:      4545454,
					OvertimeTotalHours:  8,
					OvertimePayment:     500000,
					ReimbursementsTotal: 100000,
					TotalTakeHome:       5145454,
				}
				mockPayslipUsecase.EXPECT().
					GetPayslip(gomock.Any(), int64(1)).
					Return(payslipResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid payroll ID",
			payrollID:      "invalid",
			setupContext:   func() context.Context { return context.Background() },
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "payslip not found",
			payrollID: "999",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				mockPayslipUsecase.EXPECT().
					GetPayslip(gomock.Any(), int64(999)).
					Return(nil, httppkg.NewNotFoundError("payslip not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:      "unauthorized access",
			payrollID: "1",
			setupContext: func() context.Context {
				return context.Background() // No user ID in context
			},
			setupMock: func() {
				mockPayslipUsecase.EXPECT().
					GetPayslip(gomock.Any(), int64(1)).
					Return(nil, httppkg.NewUnauthorizedError("user not authenticated"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Create HTTP request with URL parameter
			req := httptest.NewRequest(http.MethodGet, "/employee/payslip/"+tt.payrollID, nil)
			req = req.WithContext(tt.setupContext())

			// Setup chi URL parameter
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.payrollID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.GetPayslip(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				// For successful requests, verify response structure
				if w.Code == http.StatusOK {
					var response v1.PayslipResponse
					err := json.Unmarshal(w.Body.Bytes(), &response)
					if err != nil {
						t.Error("Failed to unmarshal response:", err)
					}

					if response.PayrollId == 0 {
						t.Error("Expected payroll ID in response")
					}
					if response.EmployeeId == 0 {
						t.Error("Expected employee ID in response")
					}
				}
			} else {
				// Verify error response structure
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

func TestEmployeeHandler_SubmitReimbursement(t *testing.T) {
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
			name: "successful reimbursement submission",
			requestBody: v1.PostEmployeeReimbursementJSONRequestBody{
				Amount:      100000,
				Date:        openapi_types.Date{Time: time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)},
				Description: "Business travel expenses",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				req := v1.PostEmployeeReimbursementJSONRequestBody{
					Amount:      100000,
					Date:        openapi_types.Date{Time: time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)},
					Description: "Business travel expenses",
				}
				mockReimbursementUsecase.EXPECT().
					SubmitReimbursement(gomock.Any(), req).
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
			name: "negative amount error",
			requestBody: v1.PostEmployeeReimbursementJSONRequestBody{
				Amount:      -100000,
				Date:        openapi_types.Date{Time: time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)},
				Description: "Invalid amount",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				req := v1.PostEmployeeReimbursementJSONRequestBody{
					Amount:      -100000,
					Date:        openapi_types.Date{Time: time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)},
					Description: "Invalid amount",
				}
				mockReimbursementUsecase.EXPECT().
					SubmitReimbursement(gomock.Any(), req).
					Return(httppkg.NewBadRequestError("amount must be positive"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid date format",
			requestBody: v1.PostEmployeeReimbursementJSONRequestBody{
				Amount:      100000,
				Date:        openapi_types.Date{Time: time.Time{}}, // invalid/zero time
				Description: "Travel expenses",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				req := v1.PostEmployeeReimbursementJSONRequestBody{
					Amount:      100000,
					Date:        openapi_types.Date{Time: time.Time{}},
					Description: "Travel expenses",
				}
				mockReimbursementUsecase.EXPECT().
					SubmitReimbursement(gomock.Any(), req).
					Return(httppkg.NewBadRequestError("invalid date format"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "unauthorized user",
			requestBody: v1.PostEmployeeReimbursementJSONRequestBody{
				Amount:      100000,
				Date:        openapi_types.Date{Time: time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)},
				Description: "Travel expenses",
			},
			setupContext: func() context.Context {
				return context.Background() // No user ID in context
			},
			setupMock: func() {
				req := v1.PostEmployeeReimbursementJSONRequestBody{
					Amount:      100000,
					Date:        openapi_types.Date{Time: time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)},
					Description: "Travel expenses",
				}
				mockReimbursementUsecase.EXPECT().
					SubmitReimbursement(gomock.Any(), req).
					Return(httppkg.NewUnauthorizedError("user not authenticated"))
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
			req := httptest.NewRequest(http.MethodPost, "/employee/reimbursement", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(tt.setupContext())

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.SubmitReimbursement(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError {
				// Verify error response structure
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
