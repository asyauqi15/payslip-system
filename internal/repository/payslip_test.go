package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestPayslipRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}

	baseRepo := &repository.BaseRepositoryImpl[entity.Payslip]{DB: gormDB}
	repo := repository.NewPayslipRepository(baseRepo)

	payslip := &entity.Payslip{
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

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "successful creation",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "payslips" ("created_at","updated_at","employee_id","payroll_id","base_salary","attendance_count","total_working_days","prorated_salary","overtime_total_hours","overtime_total_amount","reimbursement_total","total_take_home") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1), int64(1), int64(5000000), 20, 22, int64(4545454), 10, int64(500000), int64(100000), int64(5145454)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "payslips"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			ctx := context.WithValue(context.Background(), "skip_audit", true)
			result, err := repo.Create(ctx, payslip, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if result == nil {
					t.Error("Expected result but got nil")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPayslipRepository_FindByTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}

	baseRepo := &repository.BaseRepositoryImpl[entity.Payslip]{DB: gormDB}
	repo := repository.NewPayslipRepository(baseRepo)

	tests := []struct {
		name        string
		template    *entity.Payslip
		setupMock   func()
		expectError bool
		expectCount int
	}{
		{
			name: "successful find by employee",
			template: &entity.Payslip{
				EmployeeID: 1,
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "employee_id", "payroll_id", "base_salary", "attendance_count", "total_working_days", "prorated_salary", "overtime_total_hours", "overtime_total_amount", "reimbursement_total", "total_take_home"}).
					AddRow(1, time.Now(), time.Now(), 1, 1, 5000000, 20, 22, 4545454, 10, 500000, 100000, 5145454).
					AddRow(2, time.Now(), time.Now(), 1, 2, 5000000, 21, 22, 4772727, 5, 250000, 50000, 5072727)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payslips" WHERE "payslips"."employee_id" = $1`)).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "successful find by payroll",
			template: &entity.Payslip{
				PayrollID: 1,
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "employee_id", "payroll_id", "base_salary", "attendance_count", "total_working_days", "prorated_salary", "overtime_total_hours", "overtime_total_amount", "reimbursement_total", "total_take_home"}).
					AddRow(1, time.Now(), time.Now(), 1, 1, 5000000, 20, 22, 4545454, 10, 500000, 100000, 5145454).
					AddRow(2, time.Now(), time.Now(), 2, 1, 4000000, 22, 22, 4000000, 0, 0, 0, 4000000)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payslips" WHERE "payslips"."payroll_id" = $1`)).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "no payslip records",
			template: &entity.Payslip{
				EmployeeID: 999,
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "employee_id", "payroll_id", "base_salary", "attendance_count", "total_working_days", "prorated_salary", "overtime_total_hours", "overtime_total_amount", "reimbursement_total", "total_take_home"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payslips" WHERE "payslips"."employee_id" = $1`)).
					WithArgs(int64(999)).
					WillReturnRows(rows)
			},
			expectError: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			results, err := repo.FindByTemplate(context.Background(), tt.template, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if len(results) != tt.expectCount {
					t.Errorf("Expected %d results but got %d", tt.expectCount, len(results))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPayslipRepository_FindOneByTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}

	baseRepo := &repository.BaseRepositoryImpl[entity.Payslip]{DB: gormDB}
	repo := repository.NewPayslipRepository(baseRepo)

	template := &entity.Payslip{
		EmployeeID: 1,
		PayrollID:  1,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
		expectNil   bool
	}{
		{
			name: "successful find",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "employee_id", "payroll_id", "base_salary", "attendance_count", "total_working_days", "prorated_salary", "overtime_total_hours", "overtime_total_amount", "reimbursement_total", "total_take_home"}).
					AddRow(1, time.Now(), time.Now(), 1, 1, 5000000, 20, 22, 4545454, 10, 500000, 100000, 5145454)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payslips" WHERE "payslips"."employee_id" = $1 AND "payslips"."payroll_id" = $2 ORDER BY "payslips"."id" LIMIT $3`)).
					WithArgs(int64(1), int64(1), 1).
					WillReturnRows(rows)
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name: "not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payslips" WHERE "payslips"."employee_id" = $1 AND "payslips"."payroll_id" = $2 ORDER BY "payslips"."id" LIMIT $3`)).
					WithArgs(int64(1), int64(1), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: false,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.FindOneByTemplate(context.Background(), template, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}

			if tt.expectNil && result != nil {
				t.Error("Expected nil result but got value")
			}

			if !tt.expectNil && !tt.expectError && result == nil {
				t.Error("Expected result but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
