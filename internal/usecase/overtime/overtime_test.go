package overtime_test

import (
	"context"
	"testing"
	"time"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository/mock"
	"github.com/asyauqi15/payslip-system/internal/usecase/overtime"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestOvertimeUsecase_SubmitOvertime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOvertimeRepo := mock.NewMockOvertimeRepository(ctrl)
	mockEmployeeRepo := mock.NewMockEmployeeRepository(ctrl)
	mockAttendanceRepo := mock.NewMockAttendanceRepository(ctrl)

	usecase := overtime.NewUsecase(mockOvertimeRepo, mockEmployeeRepo, mockAttendanceRepo)

	startTime := time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC) // 6 PM
	endTime := time.Date(2025, 1, 1, 20, 0, 0, 0, time.UTC)   // 8 PM

	tests := []struct {
		name         string
		request      v1.OvertimeRequest
		setupContext func() context.Context
		setupMock    func()
		expectError  bool
	}{
		{
			name: "successful overtime submission",
			request: v1.OvertimeRequest{
				StartTime:   startTime,
				EndTime:     endTime,
				Description: "Working on urgent project",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}

				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)

				// Mock finding no conflicting overtime
				mockOvertimeRepo.EXPECT().
					FindByTemplate(gomock.Any(), gomock.Any(), nil).
					Return([]entity.Overtime{}, nil)

				// Mock creating overtime
				mockOvertimeRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.Overtime{
						Base:        entity.Base{ID: 1},
						EmployeeID:  1,
						StartAt:     startTime,
						EndAt:       endTime,
						Description: "Working on urgent project",
					}, nil)
			},
			expectError: false,
		},
		{
			name: "user not authenticated",
			request: v1.OvertimeRequest{
				StartTime:   startTime,
				EndTime:     endTime,
				Description: "Working on urgent project",
			},
			setupContext: func() context.Context {
				return context.Background() // No user ID in context
			},
			setupMock:   func() {},
			expectError: true,
		},
		{
			name: "employee not found",
			request: v1.OvertimeRequest{
				StartTime:   startTime,
				EndTime:     endTime,
				Description: "Working on urgent project",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(999))
			},
			setupMock: func() {
				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(999)}, nil).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name: "invalid time range (end before start)",
			request: v1.OvertimeRequest{
				StartTime:   endTime,   // Later time
				EndTime:     startTime, // Earlier time
				Description: "Working on urgent project",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}

				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)
			},
			expectError: true,
		},
		{
			name: "overtime too long (more than 3 hours)",
			request: v1.OvertimeRequest{
				StartTime:   startTime,
				EndTime:     startTime.Add(4 * time.Hour), // 4 hours
				Description: "Working on urgent project",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}

				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)
			},
			expectError: true,
		},
		{
			name: "weekend overtime without attendance",
			request: v1.OvertimeRequest{
				StartTime:   time.Date(2025, 1, 4, 18, 0, 0, 0, time.UTC), // Saturday
				EndTime:     time.Date(2025, 1, 4, 20, 0, 0, 0, time.UTC), // Saturday
				Description: "Weekend work",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}

				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)

				// Mock finding no conflicting overtime
				mockOvertimeRepo.EXPECT().
					FindByTemplate(gomock.Any(), gomock.Any(), nil).
					Return([]entity.Overtime{}, nil)

				// Mock creating overtime (weekend work is allowed)
				mockOvertimeRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.Overtime{
						Base:        entity.Base{ID: 1},
						EmployeeID:  1,
						StartAt:     time.Date(2025, 1, 4, 18, 0, 0, 0, time.UTC),
						EndAt:       time.Date(2025, 1, 4, 20, 0, 0, 0, time.UTC),
						Description: "Weekend work",
					}, nil)
			},
			expectError: false,
		},
		{
			name: "conflicting overtime",
			request: v1.OvertimeRequest{
				StartTime:   startTime,
				EndTime:     endTime,
				Description: "Working on urgent project",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}

				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)

				// Mock finding conflicting overtime
				conflictingOvertime := entity.Overtime{
					Base:        entity.Base{ID: 2},
					EmployeeID:  1,
					StartAt:     startTime.Add(-30 * time.Minute), // Overlapping time
					EndAt:       startTime.Add(30 * time.Minute),
					Description: "Existing overtime",
				}
				mockOvertimeRepo.EXPECT().
					FindByTemplate(gomock.Any(), gomock.Any(), nil).
					Return([]entity.Overtime{conflictingOvertime}, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			ctx := tt.setupContext()

			err := usecase.SubmitOvertime(ctx, tt.request)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
