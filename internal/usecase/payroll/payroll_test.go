package payroll_test

import (
	"context"
	"testing"
	"time"

	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository/mock"
	"github.com/asyauqi15/payslip-system/internal/usecase/payroll"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestPayrollUsecase_RunPayroll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPayrollRepo := mock.NewMockPayrollRepository(ctrl)
	mockPayslipRepo := mock.NewMockPayslipRepository(ctrl)
	mockEmployeeRepo := mock.NewMockEmployeeRepository(ctrl)
	mockAttendanceRepo := mock.NewMockAttendanceRepository(ctrl)
	mockAttendancePeriodRepo := mock.NewMockAttendancePeriodRepository(ctrl)
	mockOvertimeRepo := mock.NewMockOvertimeRepository(ctrl)
	mockReimbursementRepo := mock.NewMockReimbursementRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	usecase := payroll.NewUsecase(
		mockPayrollRepo,
		mockPayslipRepo,
		mockEmployeeRepo,
		mockAttendanceRepo,
		mockAttendancePeriodRepo,
		mockOvertimeRepo,
		mockReimbursementRepo,
		mockUserRepo,
	)

	tests := []struct {
		name        string
		request     v1.PostAdminPayrollsJSONRequestBody
		setupMock   func()
		expectError bool
	}{
		{
			name: "successful payroll run",
			request: v1.PostAdminPayrollsJSONRequestBody{
				AttendancePeriodId: 1,
			},
			setupMock: func() {
				// Mock attendance period lookup
				attendancePeriod := &entity.AttendancePeriod{
					Base:      entity.Base{ID: 1},
					StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
				}
				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(attendancePeriod, nil)

				// Mock existing payroll check (should return nil)
				mockPayrollRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Payroll{AttendancePeriodID: int64(1)}, nil).
					Return(nil, gorm.ErrRecordNotFound)

				// Mock employees lookup
				employees := []entity.Employee{
					{
						Base:       entity.Base{ID: 1},
						UserID:     1,
						BaseSalary: 5000000,
					},
					{
						Base:       entity.Base{ID: 2},
						UserID:     2,
						BaseSalary: 4000000,
					},
				}
				mockEmployeeRepo.EXPECT().
					FindByTemplate(gomock.Any(), &entity.Employee{}, nil).
					Return(employees, nil)

				// Mock payroll creation
				createdPayroll := &entity.Payroll{
					Base:               entity.Base{ID: 1},
					AttendancePeriodID: 1,
					TotalEmployees:     2,
				}
				mockPayrollRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(createdPayroll, nil)

				// Mock attendance count for each employee
				mockAttendanceRepo.EXPECT().
					CountAttendanceInPeriod(gomock.Any(), int64(1), attendancePeriod.StartDate, attendancePeriod.EndDate, nil).
					Return(int64(20), nil)
				mockAttendanceRepo.EXPECT().
					CountAttendanceInPeriod(gomock.Any(), int64(2), attendancePeriod.StartDate, attendancePeriod.EndDate, nil).
					Return(int64(20), nil)

				// Mock overtime lookup for each employee
				mockOvertimeRepo.EXPECT().
					FindByTemplate(gomock.Any(), gomock.Any(), nil).
					Return([]entity.Overtime{}, nil).Times(2) // No overtime for simplicity

				// Mock reimbursement lookup for each employee
				mockReimbursementRepo.EXPECT().
					FindByTemplate(gomock.Any(), gomock.Any(), nil).
					Return([]entity.Reimbursement{}, nil).Times(2) // No reimbursements for simplicity

				// Mock payslip creation for each employee
				mockPayslipRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.Payslip{Base: entity.Base{ID: 1}, TotalTakeHome: 4545454}, nil).
					Times(1)
				mockPayslipRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.Payslip{Base: entity.Base{ID: 2}, TotalTakeHome: 3636363}, nil).
					Times(1)

				// Mock payroll update
				mockPayrollRepo.EXPECT().
					Updates(gomock.Any(), createdPayroll, gomock.Any(), nil).
					Return(createdPayroll, nil)
			},
			expectError: false,
		},
		{
			name: "attendance period not found",
			request: v1.PostAdminPayrollsJSONRequestBody{
				AttendancePeriodId: 999,
			},
			setupMock: func() {
				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(999), nil).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name: "payroll already exists",
			request: v1.PostAdminPayrollsJSONRequestBody{
				AttendancePeriodId: 1,
			},
			setupMock: func() {
				// Mock attendance period lookup
				attendancePeriod := &entity.AttendancePeriod{
					Base:      entity.Base{ID: 1},
					StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
				}
				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(attendancePeriod, nil)

				// Mock existing payroll check (should return existing payroll)
				existingPayroll := &entity.Payroll{
					Base:               entity.Base{ID: 1},
					AttendancePeriodID: 1,
				}
				mockPayrollRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Payroll{AttendancePeriodID: int64(1)}, nil).
					Return(existingPayroll, nil)
			},
			expectError: true,
		},
		{
			name: "no employees found",
			request: v1.PostAdminPayrollsJSONRequestBody{
				AttendancePeriodId: 1,
			},
			setupMock: func() {
				// Mock attendance period lookup
				attendancePeriod := &entity.AttendancePeriod{
					Base:      entity.Base{ID: 1},
					StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
				}
				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(attendancePeriod, nil)

				// Mock existing payroll check
				mockPayrollRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Payroll{AttendancePeriodID: int64(1)}, nil).
					Return(nil, gorm.ErrRecordNotFound)

				// Mock employees lookup - empty result
				mockEmployeeRepo.EXPECT().
					FindByTemplate(gomock.Any(), &entity.Employee{}, nil).
					Return([]entity.Employee{}, nil)
			},
			expectError: true,
		},
		{
			name: "with overtime and reimbursements",
			request: v1.PostAdminPayrollsJSONRequestBody{
				AttendancePeriodId: 1,
			},
			setupMock: func() {
				// Mock attendance period lookup
				attendancePeriod := &entity.AttendancePeriod{
					Base:      entity.Base{ID: 1},
					StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
				}
				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(attendancePeriod, nil)

				// Mock existing payroll check
				mockPayrollRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Payroll{AttendancePeriodID: int64(1)}, nil).
					Return(nil, gorm.ErrRecordNotFound)

				// Mock employees lookup
				employees := []entity.Employee{
					{
						Base:       entity.Base{ID: 1},
						UserID:     1,
						BaseSalary: 5000000,
					},
				}
				mockEmployeeRepo.EXPECT().
					FindByTemplate(gomock.Any(), &entity.Employee{}, nil).
					Return(employees, nil)

				// Mock payroll creation
				createdPayroll := &entity.Payroll{
					Base:               entity.Base{ID: 1},
					AttendancePeriodID: 1,
					TotalEmployees:     1,
				}
				mockPayrollRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(createdPayroll, nil)

				// Mock attendance count
				mockAttendanceRepo.EXPECT().
					CountAttendanceInPeriod(gomock.Any(), int64(1), attendancePeriod.StartDate, attendancePeriod.EndDate, nil).
					Return(int64(20), nil)

				// Mock overtime lookup (with overtime)
				overtime := []entity.Overtime{
					{
						Base:       entity.Base{ID: 1},
						EmployeeID: 1,
						StartAt:    time.Date(2025, 1, 15, 18, 0, 0, 0, time.UTC),
						EndAt:      time.Date(2025, 1, 15, 20, 0, 0, 0, time.UTC),
					},
				}
				mockOvertimeRepo.EXPECT().
					FindByTemplate(gomock.Any(), gomock.Any(), nil).
					Return(overtime, nil)

				// Mock reimbursement lookup (with reimbursements)
				reimbursements := []entity.Reimbursement{
					{
						Base:       entity.Base{ID: 1},
						EmployeeID: 1,
						Amount:     100000,
						Date:       time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
					},
				}
				mockReimbursementRepo.EXPECT().
					FindByTemplate(gomock.Any(), gomock.Any(), nil).
					Return(reimbursements, nil)

				// Mock payslip creation
				payslip := &entity.Payslip{
					Base:               entity.Base{ID: 1},
					TotalTakeHome:      5145454,
					OvertimeTotalPay:   500000,
					ReimbursementTotal: 100000,
				}
				mockPayslipRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(payslip, nil)

				// Mock payroll update
				mockPayrollRepo.EXPECT().
					Updates(gomock.Any(), createdPayroll, gomock.Any(), nil).
					Return(createdPayroll, nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := usecase.RunPayroll(context.Background(), tt.request)

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

func TestPayrollUsecase_GetPayrollSummary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPayrollRepo := mock.NewMockPayrollRepository(ctrl)
	mockPayslipRepo := mock.NewMockPayslipRepository(ctrl)
	mockEmployeeRepo := mock.NewMockEmployeeRepository(ctrl)
	mockAttendanceRepo := mock.NewMockAttendanceRepository(ctrl)
	mockAttendancePeriodRepo := mock.NewMockAttendancePeriodRepository(ctrl)
	mockOvertimeRepo := mock.NewMockOvertimeRepository(ctrl)
	mockReimbursementRepo := mock.NewMockReimbursementRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	usecase := payroll.NewUsecase(
		mockPayrollRepo,
		mockPayslipRepo,
		mockEmployeeRepo,
		mockAttendanceRepo,
		mockAttendancePeriodRepo,
		mockOvertimeRepo,
		mockReimbursementRepo,
		mockUserRepo,
	)

	tests := []struct {
		name        string
		payrollID   int64
		setupMock   func()
		expectError bool
	}{
		{
			name:      "successful payroll summary retrieval",
			payrollID: 1,
			setupMock: func() {
				// Mock payroll lookup
				payroll := &entity.Payroll{
					Base:               entity.Base{ID: 1},
					AttendancePeriodID: 1,
					TotalEmployees:     2,
					TotalReimbursement: 200000,
					TotalOvertime:      500000,
					TotalPayroll:       8000000,
				}
				mockPayrollRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(payroll, nil)

				// Mock attendance period lookup
				attendancePeriod := &entity.AttendancePeriod{
					Base:      entity.Base{ID: 1},
					StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
				}
				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(attendancePeriod, nil)

				// Mock payslips lookup
				payslips := []entity.Payslip{
					{
						Base:               entity.Base{ID: 1},
						EmployeeID:         1,
						PayrollID:          1,
						BaseSalary:         5000000,
						AttendanceCount:    20,
						TotalWorkingDays:   23,
						ProratedSalary:     4347826,
						OvertimeTotalHours: 8,
						OvertimeTotalPay:   300000,
						ReimbursementTotal: 100000,
						TotalTakeHome:      4747826,
					},
					{
						Base:               entity.Base{ID: 2},
						EmployeeID:         2,
						PayrollID:          1,
						BaseSalary:         4000000,
						AttendanceCount:    22,
						TotalWorkingDays:   23,
						ProratedSalary:     3826087,
						OvertimeTotalHours: 4,
						OvertimeTotalPay:   200000,
						ReimbursementTotal: 100000,
						TotalTakeHome:      4126087,
					},
				}
				mockPayslipRepo.EXPECT().
					FindByTemplate(gomock.Any(), &entity.Payslip{PayrollID: int64(1)}, nil).
					Return(payslips, nil)

				// Mock employee lookups
				employee1 := &entity.Employee{
					Base:   entity.Base{ID: 1},
					UserID: 1,
				}
				employee2 := &entity.Employee{
					Base:   entity.Base{ID: 2},
					UserID: 2,
				}
				mockEmployeeRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(employee1, nil)
				mockEmployeeRepo.EXPECT().
					FindByID(gomock.Any(), uint(2), nil).
					Return(employee2, nil)

				// Mock user lookups
				user1 := &entity.User{
					Base:     entity.Base{ID: 1},
					Username: "john_doe",
				}
				user2 := &entity.User{
					Base:     entity.Base{ID: 2},
					Username: "jane_smith",
				}
				mockUserRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(user1, nil)
				mockUserRepo.EXPECT().
					FindByID(gomock.Any(), uint(2), nil).
					Return(user2, nil)
			},
			expectError: false,
		},
		{
			name:      "payroll not found",
			payrollID: 999,
			setupMock: func() {
				mockPayrollRepo.EXPECT().
					FindByID(gomock.Any(), uint(999), nil).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:      "attendance period not found",
			payrollID: 1,
			setupMock: func() {
				// Mock payroll lookup
				payroll := &entity.Payroll{
					Base:               entity.Base{ID: 1},
					AttendancePeriodID: 999,
					TotalEmployees:     1,
					TotalReimbursement: 0,
					TotalOvertime:      0,
					TotalPayroll:       5000000,
				}
				mockPayrollRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(payroll, nil)

				// Mock attendance period not found
				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(999), nil).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:      "no payslips found",
			payrollID: 1,
			setupMock: func() {
				// Mock payroll lookup
				payroll := &entity.Payroll{
					Base:               entity.Base{ID: 1},
					AttendancePeriodID: 1,
					TotalEmployees:     0,
					TotalReimbursement: 0,
					TotalOvertime:      0,
					TotalPayroll:       0,
				}
				mockPayrollRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(payroll, nil)

				// Mock attendance period lookup
				attendancePeriod := &entity.AttendancePeriod{
					Base:      entity.Base{ID: 1},
					StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
				}
				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(attendancePeriod, nil)

				// Mock payslips lookup - empty slice
				mockPayslipRepo.EXPECT().
					FindByTemplate(gomock.Any(), &entity.Payslip{PayrollID: int64(1)}, nil).
					Return([]entity.Payslip{}, nil)
			},
			expectError: false,
		},
		{
			name:      "employee not found for payslip",
			payrollID: 1,
			setupMock: func() {
				// Mock payroll lookup
				payroll := &entity.Payroll{
					Base:               entity.Base{ID: 1},
					AttendancePeriodID: 1,
					TotalEmployees:     1,
					TotalReimbursement: 0,
					TotalOvertime:      0,
					TotalPayroll:       5000000,
				}
				mockPayrollRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(payroll, nil)

				// Mock attendance period lookup
				attendancePeriod := &entity.AttendancePeriod{
					Base:      entity.Base{ID: 1},
					StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
				}
				mockAttendancePeriodRepo.EXPECT().
					FindByID(gomock.Any(), uint(1), nil).
					Return(attendancePeriod, nil)

				// Mock payslips lookup
				payslips := []entity.Payslip{
					{
						Base:       entity.Base{ID: 1},
						EmployeeID: 999, // Non-existent employee
						PayrollID:  1,
					},
				}
				mockPayslipRepo.EXPECT().
					FindByTemplate(gomock.Any(), &entity.Payslip{PayrollID: int64(1)}, nil).
					Return(payslips, nil)

				// Mock employee not found
				mockEmployeeRepo.EXPECT().
					FindByID(gomock.Any(), uint(999), nil).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: false, // Should continue processing other payslips
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := usecase.GetPayrollSummary(context.Background(), tt.payrollID)

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
					// Basic validation of result structure
					if result.PayrollId != tt.payrollID {
						t.Errorf("Expected payroll ID %d but got %d", tt.payrollID, result.PayrollId)
					}
				}
			}
		})
	}
}
