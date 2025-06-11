package employee_test

import (
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
	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

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
				return context.Background()
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

			req := httptest.NewRequest(http.MethodGet, "/employee/payslip/"+tt.payrollID, nil)
			req = req.WithContext(tt.setupContext())

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.payrollID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			handler.GetPayslip(w, req)

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
