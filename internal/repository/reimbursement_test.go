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

func TestReimbursementRepository_Create(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Reimbursement]{DB: gormDB}
	repo := repository.NewReimbursementRepository(baseRepo)

	reimbursementDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	reimbursement := &entity.Reimbursement{
		EmployeeID:  1,
		Amount:      50000,
		Date:        reimbursementDate,
		Description: "Business travel expenses",
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
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "reimbursements" ("created_at","updated_at","employee_id","amount","date","description") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1), int64(50000), reimbursementDate, "Business travel expenses").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "reimbursements"`)).
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
			result, err := repo.Create(ctx, reimbursement, nil)

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

func TestReimbursementRepository_FindByTemplate(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Reimbursement]{DB: gormDB}
	repo := repository.NewReimbursementRepository(baseRepo)

	template := &entity.Reimbursement{
		EmployeeID: 1,
	}

	reimbursementDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
		expectCount int
	}{
		{
			name: "successful find by employee",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "employee_id", "amount", "date", "description"}).
					AddRow(1, time.Now(), time.Now(), 1, 50000, reimbursementDate, "Business travel expenses").
					AddRow(2, time.Now(), time.Now(), 1, 25000, reimbursementDate.AddDate(0, 0, 1), "Meal expenses")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reimbursements" WHERE "reimbursements"."employee_id" = $1`)).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "no reimbursement records",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "employee_id", "amount", "date", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reimbursements" WHERE "reimbursements"."employee_id" = $1`)).
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

func TestReimbursementRepository_Updates(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Reimbursement]{DB: gormDB}
	repo := repository.NewReimbursementRepository(baseRepo)

	reimbursementDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	reimbursement := &entity.Reimbursement{
		Base:        entity.Base{ID: 1},
		EmployeeID:  1,
		Amount:      50000,
		Date:        reimbursementDate,
		Description: "Business travel expenses",
	}

	updateData := entity.Reimbursement{
		Amount:      60000,
		Description: "Updated business travel expenses",
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "successful amount and description update",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reimbursements" SET "updated_at"=$1,"amount"=$2,"description"=$3 WHERE "id" = $4`)).
					WithArgs(sqlmock.AnyArg(), int64(60000), "Updated business travel expenses", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reimbursements"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.Updates(context.Background(), reimbursement, updateData, nil)

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
