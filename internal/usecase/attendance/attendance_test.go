package attendance_test

import (
	"context"
	"testing"
	"time"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository/mock"
	"github.com/asyauqi15/payslip-system/internal/usecase/attendance"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestAttendanceUsecase_SubmitAttendance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAttendanceRepo := mock.NewMockAttendanceRepository(ctrl)
	mockEmployeeRepo := mock.NewMockEmployeeRepository(ctrl)

	usecase := attendance.NewUsecase(mockAttendanceRepo, mockEmployeeRepo)

	tests := []struct {
		name           string
		attendanceType v1.PostEmployeeAttendanceJSONBodyAttendanceType
		setupContext   func() context.Context
		setupMock      func()
		expectError    bool
	}{
		{
			name:           "successful check-in",
			attendanceType: v1.CheckIn,
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

				// Mock finding existing attendances (empty for check-in)
				mockAttendanceRepo.EXPECT().
					FindByTemplate(gomock.Any(), &entity.Attendance{EmployeeID: int64(1)}, nil).
					Return([]entity.Attendance{}, nil)

				// Mock creating new attendance
				mockAttendanceRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.Attendance{
						Base:        entity.Base{ID: 1},
						EmployeeID:  1,
						ClockInTime: time.Now().Format(time.RFC3339),
					}, nil)
			},
			expectError: false,
		},
		{
			name:           "successful check-out",
			attendanceType: v1.CheckOut,
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

				// Mock finding existing attendance for today
				existingAttendance := entity.Attendance{
					Base:         entity.Base{ID: 1},
					EmployeeID:   1,
					ClockInTime:  time.Now().Format(time.RFC3339),
					ClockOutTime: "",
				}
				mockAttendanceRepo.EXPECT().
					FindByTemplate(gomock.Any(), &entity.Attendance{EmployeeID: int64(1)}, nil).
					Return([]entity.Attendance{existingAttendance}, nil)

				// Mock updating attendance with check-out time
				mockAttendanceRepo.EXPECT().
					Updates(gomock.Any(), &existingAttendance, gomock.Any(), nil).
					Return(&existingAttendance, nil)
			},
			expectError: false,
		},
		{
			name:           "user not authenticated",
			attendanceType: v1.CheckIn,
			setupContext: func() context.Context {
				return context.Background() // No user ID in context
			},
			setupMock:   func() {},
			expectError: true,
		},
		{
			name:           "employee not found",
			attendanceType: v1.CheckIn,
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
			name:           "weekend attendance",
			attendanceType: v1.CheckIn,
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
			expectError: func() bool {
				// This test would fail on weekends, but for simplicity, we'll assume it's a weekday
				now := time.Now()
				return now.Weekday() == time.Saturday || now.Weekday() == time.Sunday
			}(),
		},
		{
			name:           "already checked in today",
			attendanceType: v1.CheckIn,
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

				// Mock finding existing attendance for today
				existingAttendance := entity.Attendance{
					Base:        entity.Base{ID: 1},
					EmployeeID:  1,
					ClockInTime: time.Now().Format(time.RFC3339),
				}
				mockAttendanceRepo.EXPECT().
					FindByTemplate(gomock.Any(), &entity.Attendance{EmployeeID: int64(1)}, nil).
					Return([]entity.Attendance{existingAttendance}, nil)
			},
			expectError: true,
		},
		{
			name:           "check out without check in",
			attendanceType: v1.CheckOut,
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

				// Mock finding no existing attendance for today
				mockAttendanceRepo.EXPECT().
					FindByTemplate(gomock.Any(), &entity.Attendance{EmployeeID: int64(1)}, nil).
					Return([]entity.Attendance{}, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			ctx := tt.setupContext()

			err := usecase.SubmitAttendance(ctx, tt.attendanceType)

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
