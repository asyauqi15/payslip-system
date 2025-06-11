package admin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/asyauqi15/payslip-system/internal/handler/admin"
	attendanceperiodmock "github.com/asyauqi15/payslip-system/internal/usecase/attendance_period/mock"
	payrollmock "github.com/asyauqi15/payslip-system/internal/usecase/payroll/mock"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime/types"
	"go.uber.org/mock/gomock"
)

func TestAdminHandler_RunPayroll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAttendancePeriodUsecase := attendanceperiodmock.NewMockUsecase(ctrl)
	mockPayrollUsecase := payrollmock.NewMockUsecase(ctrl)

	handler := admin.NewHandler(
		mockAttendancePeriodUsecase,
		mockPayrollUsecase,
	)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful payroll run",
			requestBody: v1.PostAdminPayrollsJSONRequestBody{
				AttendancePeriodId: 1,
			},
			setupMock: func() {
				req := v1.PostAdminPayrollsJSONRequestBody{
					AttendancePeriodId: 1,
				}
				mockPayrollUsecase.EXPECT().
					RunPayroll(gomock.Any(), req).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "attendance period not found",
			requestBody: v1.PostAdminPayrollsJSONRequestBody{
				AttendancePeriodId: 999,
			},
			setupMock: func() {
				req := v1.PostAdminPayrollsJSONRequestBody{
					AttendancePeriodId: 999,
				}
				mockPayrollUsecase.EXPECT().
					RunPayroll(gomock.Any(), req).
					Return(httppkg.NewNotFoundError("attendance period not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name: "payroll already exists",
			requestBody: v1.PostAdminPayrollsJSONRequestBody{
				AttendancePeriodId: 1,
			},
			setupMock: func() {
				req := v1.PostAdminPayrollsJSONRequestBody{
					AttendancePeriodId: 1,
				}
				mockPayrollUsecase.EXPECT().
					RunPayroll(gomock.Any(), req).
					Return(httppkg.NewBadRequestError("payroll already exists for this attendance period"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "no employees found",
			requestBody: v1.PostAdminPayrollsJSONRequestBody{
				AttendancePeriodId: 2,
			},
			setupMock: func() {
				req := v1.PostAdminPayrollsJSONRequestBody{
					AttendancePeriodId: 2,
				}
				mockPayrollUsecase.EXPECT().
					RunPayroll(gomock.Any(), req).
					Return(httppkg.NewBadRequestError("no employees found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "internal server error",
			requestBody: v1.PostAdminPayrollsJSONRequestBody{
				AttendancePeriodId: 3,
			},
			setupMock: func() {
				req := v1.PostAdminPayrollsJSONRequestBody{
					AttendancePeriodId: 3,
				}
				mockPayrollUsecase.EXPECT().
					RunPayroll(gomock.Any(), req).
					Return(httppkg.NewInternalServerError("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
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

			req := httptest.NewRequest(http.MethodPost, "/admin/payrolls", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(context.Background())

			w := httptest.NewRecorder()
			handler.RunPayroll(w, req)

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

func TestAdminHandler_GetPayrollSummary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAttendancePeriodUsecase := attendanceperiodmock.NewMockUsecase(ctrl)
	mockPayrollUsecase := payrollmock.NewMockUsecase(ctrl)

	handler := admin.NewHandler(
		mockAttendancePeriodUsecase,
		mockPayrollUsecase,
	)

	tests := []struct {
		name           string
		payrollID      string
		setupMock      func()
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "successful payroll summary retrieval",
			payrollID: "1",
			setupMock: func() {
				summaryResponse := &v1.AdminPayrollSummaryResponse{
					PayrollId: 1,
					AttendancePeriod: v1.AttendancePeriod{
						StartDate: types.Date{Time: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)},
						EndDate:   types.Date{Time: time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)},
					},
					EmployeesCount:         2,
					TotalPayroll:           8081818,
					TotalReimbursementsPay: 200000,
					TotalOvertimePay:       700000,
					PayslipList: []v1.PayslipItem{
						{
							EmployeeId:            1,
							Username:              "john_doe",
							AttendanceCount:       20,
							OvertimeCount:         4,
							ProratedSalary:        4545454,
							OvertimePayment:       500000,
							ReimbursementsPayment: 100000,
							TotalPay:              5145454,
						},
						{
							EmployeeId:            2,
							Username:              "jane_smith",
							AttendanceCount:       18,
							OvertimeCount:         4,
							ProratedSalary:        3636364,
							OvertimePayment:       200000,
							ReimbursementsPayment: 100000,
							TotalPay:              3936364,
						},
					},
				}
				mockPayrollUsecase.EXPECT().
					GetPayrollSummary(gomock.Any(), int64(1)).
					Return(summaryResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid payroll ID",
			payrollID:      "invalid",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "payroll not found",
			payrollID: "999",
			setupMock: func() {
				mockPayrollUsecase.EXPECT().
					GetPayrollSummary(gomock.Any(), int64(999)).
					Return(nil, httppkg.NewNotFoundError("payroll not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:      "internal server error",
			payrollID: "1",
			setupMock: func() {
				mockPayrollUsecase.EXPECT().
					GetPayrollSummary(gomock.Any(), int64(1)).
					Return(nil, httppkg.NewInternalServerError("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/admin/payrolls/"+tt.payrollID, nil)
			req = req.WithContext(context.Background())

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.payrollID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			handler.GetPayrollSummary(w, req)

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
