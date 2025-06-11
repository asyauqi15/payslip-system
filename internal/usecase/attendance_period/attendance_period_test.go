package attendance_period_test

import (
	"context"
	"testing"
	"time"

	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository/mock"
	"github.com/asyauqi15/payslip-system/internal/usecase/attendance_period"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestAttendancePeriodUsecase_CreateAttendancePeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAttendancePeriodRepo := mock.NewMockAttendancePeriodRepository(ctrl)

	usecase := attendance_period.NewUsecase(mockAttendancePeriodRepo)

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name        string
		startDate   time.Time
		endDate     time.Time
		setupMock   func()
		expectError bool
	}{
		{
			name:      "successful creation",
			startDate: startDate,
			endDate:   endDate,
			setupMock: func() {
				mockAttendancePeriodRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.AttendancePeriod{
						Base:      entity.Base{ID: 1},
						StartDate: startDate,
						EndDate:   endDate,
					}, nil)
			},
			expectError: false,
		},
		{
			name:      "repository error",
			startDate: startDate,
			endDate:   endDate,
			setupMock: func() {
				mockAttendancePeriodRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(nil, gorm.ErrInvalidDB)
			},
			expectError: true,
		},
		{
			name:      "overlapping period error",
			startDate: startDate,
			endDate:   endDate,
			setupMock: func() {
				mockAttendancePeriodRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(nil, gorm.ErrDuplicatedKey) // Simulating overlap constraint violation
			},
			expectError: true,
		},
		{
			name:      "invalid date range (end before start)",
			startDate: endDate,   // Later date
			endDate:   startDate, // Earlier date
			setupMock: func() {
				// No mock expectation since validation should fail before repository call
			},
			expectError: true,
		},
		{
			name:      "same start and end date",
			startDate: startDate,
			endDate:   startDate,
			setupMock: func() {
				mockAttendancePeriodRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.AttendancePeriod{
						Base:      entity.Base{ID: 1},
						StartDate: startDate,
						EndDate:   startDate,
					}, nil)
			},
			expectError: false,
		},
		{
			name:      "future dates",
			startDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
			setupMock: func() {
				mockAttendancePeriodRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.AttendancePeriod{
						Base:      entity.Base{ID: 1},
						StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
					}, nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := usecase.CreateAttendancePeriod(context.Background(), tt.startDate, tt.endDate)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if result != nil {
					t.Error("Expected nil result but got value")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if result == nil {
					t.Error("Expected result but got nil")
				} else {
					if !result.StartDate.Equal(tt.startDate) {
						t.Errorf("Expected start date %v but got %v", tt.startDate, result.StartDate)
					}
					if !result.EndDate.Equal(tt.endDate) {
						t.Errorf("Expected end date %v but got %v", tt.endDate, result.EndDate)
					}
				}
			}
		})
	}
}
