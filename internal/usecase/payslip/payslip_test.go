package payslip_test

import (
	"context"
	"testing"
	"time"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository/mock"
	"github.com/asyauqi15/payslip-system/internal/usecase/payslip"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestPayslipUsecase_GetPayslip(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPayslipRepo := mock.NewMockPayslipRepository(ctrl)
	mockPayrollRepo := mock.NewMockPayrollRepository(ctrl)
	mockEmployeeRepo := mock.NewMockEmployeeRepository(ctrl)
	mockReimbursementRepo := mock.NewMockReimbursementRepository(ctrl)
	mockAttendancePeriodRepo := mock.NewMockAttendancePeriodRepository(ctrl)

	usecase := payslip.NewUsecase(
		mockPayslipRepo,
		mockPayrollRepo,
		mockEmployeeRepo,
		mockReimbursementRepo,
		mockAttendancePeriodRepo,
	)

	tests := []struct {
		name         string
		payrollID    int64
		setupContext func() context.Context
		setupMock    func()
		expectError  bool
	}{
		{
			name:      "successful payslip retrieval",
			payrollID: 1,
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				// Mock employee lookup
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}
				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)

				// Mock payslip lookup
				payslipEntity := &entity.Payslip{
					Base:               entity.Base{ID: 1},
					EmployeeID:         1,
					PayrollID:          1,
					BaseSalary:         5000000,
					AttendanceCount:    20,
					TotalWorkingDays:   22,
					ProratedSalary:     4545454,
					OvertimeTotalHours: 10,
					OvertimeTotalPay:   500000,
					ReimbursementTotal: 100000,
					TotalTakeHome:      5145454,
				}
				mockPayslipRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Payslip{EmployeeID: int64(1), PayrollID: int64(1)}, nil).
					Return(payslipEntity, nil)

				// Mock payroll lookup
				payrollEntity := &entity.Payroll{
					Base:               entity.Base{ID: 1},
					AttendancePeriodID: 1,
					TotalEmployees:     10,
					TotalReimbursement: 500000,
					TotalOvertime:      200000,
					TotalPayroll:       5000000,
				}
				mockPayrollRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(payrollEntity, nil)

				// Mock attendance period lookup
				attendancePeriod := &entity.AttendancePeriod{
					Base:      entity.Base{ID: 1},
					StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
				}
				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(attendancePeriod, nil)

				// Mock reimbursement lookup
				reimbursements := []entity.Reimbursement{
					{
						Base:        entity.Base{ID: 1},
						EmployeeID:  1,
						Amount:      50000,
						Date:        time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
						Description: "Business travel",
					},
					{
						Base:        entity.Base{ID: 2},
						EmployeeID:  1,
						Amount:      50000,
						Date:        time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
						Description: "Meal expenses",
					},
				}
				mockReimbursementRepo.EXPECT().
					FindByTemplate(gomock.Any(), gomock.Any(), nil).
					Return(reimbursements, nil)
			},
			expectError: false,
		},
		{
			name:      "user not authenticated",
			payrollID: 1,
			setupContext: func() context.Context {
				return context.Background() // No user ID in context
			},
			setupMock:   func() {},
			expectError: true,
		},
		{
			name:      "employee not found",
			payrollID: 1,
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
			name:      "payslip not found",
			payrollID: 999,
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				payrollEntity := &entity.Payroll{
					Base:               entity.Base{ID: 999},
					AttendancePeriodID: 1,
					TotalEmployees:     10,
					TotalReimbursement: 500000,
					TotalOvertime:      200000,
					TotalPayroll:       5000000,
				}
				mockPayrollRepo.EXPECT().
					FindByID(gomock.Any(), uint(999), nil).
					Return(payrollEntity, nil)

				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(&entity.AttendancePeriod{
						Base: entity.Base{ID: 1},
					}, nil)

				// Mock employee lookup
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}
				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)

				// Mock payslip not found
				mockPayslipRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Payslip{EmployeeID: int64(1), PayrollID: int64(999)}, nil).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:      "payroll not found",
			payrollID: 1,
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				// Mock employee lookup
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}
				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)

				// Mock payroll not found
				mockPayrollRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			ctx := tt.setupContext()

			result, err := usecase.GetPayslip(ctx, tt.payrollID)

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
					// Verify result structure
					if result.BaseSalary == 0 || result.BaseSalary != 5000000 {
						t.Error("Expected valid base salary in result")
					}
					if result.PayrollId == 0 || result.PayrollId != 1 {
						t.Error("Expected valid payroll ID in result")
					}
				}
			}
		})
	}
}
