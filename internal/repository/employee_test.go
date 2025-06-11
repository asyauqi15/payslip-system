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

func TestEmployeeRepository_Create(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Employee]{DB: gormDB}
	repo := repository.NewEmployeeRepository(baseRepo)

	employee := &entity.Employee{
		UserID:     1,
		BaseSalary: 50000,
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
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "employees" ("created_at","updated_at","user_id","base_salary") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1), int64(50000)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "duplicate user_id constraint violation",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "employees"`)).
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
			result, err := repo.Create(ctx, employee, nil)

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

func TestEmployeeRepository_FindByID(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Employee]{DB: gormDB}
	repo := repository.NewEmployeeRepository(baseRepo)

	tests := []struct {
		name        string
		employeeID  uint
		setupMock   func()
		expectError bool
		expectNil   bool
	}{
		{
			name:       "successful find",
			employeeID: 1,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "base_salary"}).
					AddRow(1, time.Now(), time.Now(), 1, 50000)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "employees" WHERE "employees"."id" = $1 ORDER BY "employees"."id" LIMIT $2`)).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name:       "not found",
			employeeID: 999,
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "employees" WHERE "employees"."id" = $1 ORDER BY "employees"."id" LIMIT $2`)).
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

			result, err := repo.FindByID(context.Background(), tt.employeeID, nil)

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

func TestEmployeeRepository_FindOneByTemplate(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Employee]{DB: gormDB}
	repo := repository.NewEmployeeRepository(baseRepo)

	template := &entity.Employee{
		UserID: 1,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
		expectNil   bool
	}{
		{
			name: "successful find by user_id",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "base_salary"}).
					AddRow(1, time.Now(), time.Now(), 1, 50000)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "employees" WHERE "employees"."user_id" = $1 ORDER BY "employees"."id" LIMIT $2`)).
					WithArgs(int64(1), 1).
					WillReturnRows(rows)
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name: "not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "employees" WHERE "employees"."user_id" = $1 ORDER BY "employees"."id" LIMIT $2`)).
					WithArgs(int64(1), 1).
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

func TestEmployeeRepository_Updates(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Employee]{DB: gormDB}
	repo := repository.NewEmployeeRepository(baseRepo)

	employee := &entity.Employee{
		Base:       entity.Base{ID: 1},
		UserID:     1,
		BaseSalary: 50000,
	}

	updateData := entity.Employee{
		BaseSalary: 60000,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "successful salary update",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "employees" SET "updated_at"=$1,"base_salary"=$2 WHERE "id" = $3`)).
					WithArgs(sqlmock.AnyArg(), int64(60000), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "employees"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.Updates(context.Background(), employee, updateData, nil)

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
