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

func TestOvertimeRepository_Create(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Overtime]{DB: gormDB}
	repo := repository.NewOvertimeRepository(baseRepo)

	startTime := time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 1, 20, 0, 0, 0, time.UTC)

	overtime := &entity.Overtime{
		EmployeeID:  1,
		StartAt:     startTime,
		EndAt:       endTime,
		Description: "Working on urgent project",
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
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "overtimes" ("created_at","updated_at","employee_id","start_at","end_at","description") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1), startTime, endTime, "Working on urgent project").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "overtimes"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.Create(context.Background(), overtime, nil)

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

func TestOvertimeRepository_FindByTemplate(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Overtime]{DB: gormDB}
	repo := repository.NewOvertimeRepository(baseRepo)

	template := &entity.Overtime{
		EmployeeID: 1,
	}

	startTime := time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 1, 20, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
		expectCount int
	}{
		{
			name: "successful find by employee",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "employee_id", "start_at", "end_at", "description"}).
					AddRow(1, time.Now(), time.Now(), 1, startTime, endTime, "Working on urgent project").
					AddRow(2, time.Now(), time.Now(), 1, startTime.AddDate(0, 0, 1), endTime.AddDate(0, 0, 1), "Another overtime")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "overtimes" WHERE "overtimes"."employee_id" = $1`)).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "no overtime records",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "employee_id", "start_at", "end_at", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "overtimes" WHERE "overtimes"."employee_id" = $1`)).
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

func TestOvertimeRepository_Updates(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Overtime]{DB: gormDB}
	repo := repository.NewOvertimeRepository(baseRepo)

	startTime := time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 1, 20, 0, 0, 0, time.UTC)

	overtime := &entity.Overtime{
		Base:        entity.Base{ID: 1},
		EmployeeID:  1,
		StartAt:     startTime,
		EndAt:       endTime,
		Description: "Working on urgent project",
	}

	updateData := entity.Overtime{
		Description: "Updated description",
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "successful description update",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "overtimes" SET "updated_at"=$1,"description"=$2 WHERE "id" = $3`)).
					WithArgs(sqlmock.AnyArg(), "Updated description", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "overtimes"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.Updates(context.Background(), overtime, updateData, nil)

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
