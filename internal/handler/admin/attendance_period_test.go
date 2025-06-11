package admin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/handler/admin"
	attendanceperiodmock "github.com/asyauqi15/payslip-system/internal/usecase/attendance_period/mock"
	payrollmock "github.com/asyauqi15/payslip-system/internal/usecase/payroll/mock"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/oapi-codegen/runtime/types"
	"go.uber.org/mock/gomock"
)

func TestAdminHandler_CreateAttendancePeriod(t *testing.T) {
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
			name: "successful attendance period creation",
			requestBody: v1.AttendancePeriodRequest{
				StartDate: types.Date{Time: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)},
				EndDate:   types.Date{Time: time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)},
			},
			setupMock: func() {
				startDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)

				createdPeriod := &entity.AttendancePeriod{
					Base:      entity.Base{ID: 1},
					StartDate: startDate,
					EndDate:   endDate,
				}

				mockAttendancePeriodUsecase.EXPECT().
					CreateAttendancePeriod(gomock.Any(), startDate, endDate).
					Return(createdPeriod, nil)
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
			name: "overlapping period error",
			requestBody: v1.AttendancePeriodRequest{
				StartDate: types.Date{Time: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)},
				EndDate:   types.Date{Time: time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)},
			},
			setupMock: func() {
				startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

				mockAttendancePeriodUsecase.EXPECT().
					CreateAttendancePeriod(gomock.Any(), startDate, endDate).
					Return(nil, httppkg.NewBadRequestError("attendance period overlaps with existing period"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid date range error",
			requestBody: v1.AttendancePeriodRequest{
				StartDate: types.Date{Time: time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)},
				EndDate:   types.Date{Time: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)},
			},
			setupMock: func() {
				startDate := time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)

				mockAttendancePeriodUsecase.EXPECT().
					CreateAttendancePeriod(gomock.Any(), startDate, endDate).
					Return(nil, httppkg.NewBadRequestError("start date must be before end date"))
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "internal server error",
			requestBody: v1.AttendancePeriodRequest{
				StartDate: types.Date{Time: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)},
				EndDate:   types.Date{Time: time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC)},
			},
			setupMock: func() {
				startDate := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC)

				mockAttendancePeriodUsecase.EXPECT().
					CreateAttendancePeriod(gomock.Any(), startDate, endDate).
					Return(nil, httppkg.NewInternalServerError("database connection failed"))
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

			req := httptest.NewRequest(http.MethodPost, "/admin/attendance-periods", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(context.Background())

			w := httptest.NewRecorder()
			handler.CreateAttendancePeriod(w, req)

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
