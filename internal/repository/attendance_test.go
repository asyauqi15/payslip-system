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

func TestAttendanceRepository_Create(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Attendance]{DB: gormDB}
	repo := repository.NewAttendanceRepository(baseRepo)

	attendance := &entity.Attendance{
		EmployeeID:   1,
		ClockInTime:  "2025-01-01T09:00:00Z",
		ClockOutTime: "",
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
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "attendances" ("created_at","updated_at","employee_id","clock_in_time","clock_out_time") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1), "2025-01-01T09:00:00Z", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "attendances"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.Create(context.Background(), attendance, nil)

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

func TestAttendanceRepository_CountAttendanceInPeriod(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Attendance]{DB: gormDB}
	repo := repository.NewAttendanceRepository(baseRepo)

	employeeID := int64(1)
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name          string
		setupMock     func()
		expectError   bool
		expectedCount int64
	}{
		{
			name: "successful count",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "attendances" WHERE employee_id = $1 AND clock_in_time != '' AND clock_in_time >= $2 AND clock_in_time <= $3`)).
					WithArgs(employeeID, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))
			},
			expectError:   false,
			expectedCount: 20,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "attendances"`)).
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectError:   true,
			expectedCount: 0,
		},
		{
			name: "no attendance records",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "attendances" WHERE employee_id = $1 AND clock_in_time != '' AND clock_in_time >= $2 AND clock_in_time <= $3`)).
					WithArgs(employeeID, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectError:   false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			count, err := repo.CountAttendanceInPeriod(context.Background(), employeeID, startDate, endDate, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if count != tt.expectedCount {
					t.Errorf("Expected count %d but got %d", tt.expectedCount, count)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestAttendanceRepository_FindByTemplate(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Attendance]{DB: gormDB}
	repo := repository.NewAttendanceRepository(baseRepo)

	template := &entity.Attendance{
		EmployeeID: 1,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
		expectCount int
	}{
		{
			name: "successful find by employee",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "employee_id", "clock_in_time", "clock_out_time"}).
					AddRow(1, time.Now(), time.Now(), 1, "2025-01-01T09:00:00Z", "2025-01-01T17:00:00Z").
					AddRow(2, time.Now(), time.Now(), 1, "2025-01-02T09:00:00Z", "2025-01-02T17:00:00Z")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "attendances" WHERE "attendances"."employee_id" = $1`)).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "no attendance records",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "employee_id", "clock_in_time", "clock_out_time"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "attendances" WHERE "attendances"."employee_id" = $1`)).
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

func TestAttendanceRepository_Updates(t *testing.T) {
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

	baseRepo := &repository.BaseRepositoryImpl[entity.Attendance]{DB: gormDB}
	repo := repository.NewAttendanceRepository(baseRepo)

	attendance := &entity.Attendance{
		Base:         entity.Base{ID: 1},
		EmployeeID:   1,
		ClockInTime:  "2025-01-01T09:00:00Z",
		ClockOutTime: "",
	}

	updateData := entity.Attendance{
		ClockOutTime: "2025-01-01T17:00:00Z",
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "successful clock out update",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "attendances" SET "updated_at"=$1,"clock_out_time"=$2 WHERE "id" = $3`)).
					WithArgs(sqlmock.AnyArg(), "2025-01-01T17:00:00Z", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "attendances"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.Updates(context.Background(), attendance, updateData, nil)

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
