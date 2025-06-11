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

func TestPayrollRepository_Create(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Payroll]{DB: gormDB}
	repo := repository.NewPayrollRepository(baseRepo)

	payroll := &entity.Payroll{
		AttendancePeriodID: 1,
		TotalEmployees:     10,
		TotalReimbursement: 500000,
		TotalOvertime:      200000,
		TotalPayroll:       5000000,
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
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "payrolls" ("created_at","updated_at","attendance_period_id","employees_count","total_reimbursement","total_overtime","total_payroll") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1), int64(10), int64(500000), int64(200000), int64(5000000)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "duplicate attendance period constraint violation",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "payrolls"`)).
					WillReturnError(gorm.ErrDuplicatedKey)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			ctx := context.WithValue(context.Background(), "skip_audit", true)
			result, err := repo.Create(ctx, payroll, nil)

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

func TestPayrollRepository_FindByID(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Payroll]{DB: gormDB}
	repo := repository.NewPayrollRepository(baseRepo)

	tests := []struct {
		name        string
		payrollID   uint
		setupMock   func()
		expectError bool
		expectNil   bool
	}{
		{
			name:      "successful find",
			payrollID: 1,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "attendance_period_id", "employees_count", "total_reimbursement", "total_overtime", "total_payroll"}).
					AddRow(1, time.Now(), time.Now(), 1, 10, 500000, 200000, 5000000)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payrolls" WHERE "payrolls"."id" = $1 ORDER BY "payrolls"."id" LIMIT $2`)).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name:      "not found",
			payrollID: 999,
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payrolls" WHERE "payrolls"."id" = $1 ORDER BY "payrolls"."id" LIMIT $2`)).
					WithArgs(999, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.FindByID(context.Background(), tt.payrollID, nil)

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

func TestPayrollRepository_FindByTemplate(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Payroll]{DB: gormDB}
	repo := repository.NewPayrollRepository(baseRepo)

	template := &entity.Payroll{
		AttendancePeriodID: 1,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
		expectCount int
	}{
		{
			name: "successful find by attendance period",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "attendance_period_id", "employees_count", "total_reimbursement", "total_overtime", "total_payroll"}).
					AddRow(1, time.Now(), time.Now(), 1, 10, 500000, 200000, 5000000)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payrolls" WHERE "payrolls"."attendance_period_id" = $1`)).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "no payroll records",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "attendance_period_id", "employees_count", "total_reimbursement", "total_overtime", "total_payroll"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payrolls" WHERE "payrolls"."attendance_period_id" = $1`)).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectError: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			results, err := repo.FindByTemplate(context.Background(), template, nil)

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
