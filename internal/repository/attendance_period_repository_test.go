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

func setupAttendancePeriodRepoTest() (*gorm.DB, sqlmock.Sqlmock, repository.AttendancePeriodRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	baseRepo := &repository.BaseRepositoryImpl[entity.AttendancePeriod]{DB: gormDB}
	repo := repository.NewAttendancePeriodRepository(baseRepo)

	return gormDB, mock, repo
}

func TestAttendancePeriodRepository_Create(t *testing.T) {
	_, mock, repo := setupAttendancePeriodRepoTest()

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	attendancePeriod := &entity.AttendancePeriod{
		StartDate: startDate,
		EndDate:   endDate,
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "attendance_periods" WHERE start_date <= $1 AND end_date >= $2`)).
					WithArgs(endDate, startDate).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "attendance_periods" ("created_at","updated_at","start_date","end_date") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), startDate, endDate).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "overlapping period error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "attendance_periods" WHERE start_date <= $1 AND end_date >= $2`)).
					WithArgs(endDate, startDate).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "attendance_periods" WHERE start_date <= $1 AND end_date >= $2`)).
					WithArgs(endDate, startDate).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "attendance_periods"`)).
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
			result, err := repo.Create(ctx, attendancePeriod, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if result != nil {
					if !result.StartDate.Equal(startDate) {
						t.Errorf("Expected start date %v but got %v", startDate, result.StartDate)
					}
					if !result.EndDate.Equal(endDate) {
						t.Errorf("Expected end date %v but got %v", endDate, result.EndDate)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestAttendancePeriodRepository_Update(t *testing.T) {
	_, mock, repo := setupAttendancePeriodRepoTest()

	now := time.Now()
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	attendancePeriod := &entity.AttendancePeriod{
		Base: entity.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		StartDate: startDate,
		EndDate:   endDate,
	}

	updateData := &entity.AttendancePeriod{
		StartDate: startDate,
		EndDate:   endDate,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "successful update",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "attendance_periods" WHERE id != $1 AND start_date <= $2 AND end_date >= $3`)).
					WithArgs(int64(1), endDate, startDate).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "attendance_periods" SET "updated_at"=$1,"start_date"=$2,"end_date"=$3 WHERE "id" = $4`)).
					WithArgs(sqlmock.AnyArg(), startDate, endDate, int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "overlapping period error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "attendance_periods" WHERE id != $1 AND start_date <= $2 AND end_date >= $3`)).
					WithArgs(int64(1), endDate, startDate).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "attendance_periods" WHERE id != $1 AND start_date <= $2 AND end_date >= $3`)).
					WithArgs(int64(1), endDate, startDate).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "attendance_periods"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.Updates(context.Background(), attendancePeriod, *updateData, nil)

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
