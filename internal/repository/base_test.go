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

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

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

	return gormDB, mock
}

func TestBaseRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := &repository.BaseRepositoryImpl[entity.User]{DB: db}

	user := &entity.User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Role:         "admin",
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
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","username","password_hash","role") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser", "hashedpassword", "admin").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.Create(context.Background(), user, nil)

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

func TestBaseRepository_FindByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := &repository.BaseRepositoryImpl[entity.User]{DB: db}

	tests := []struct {
		name        string
		userID      uint
		setupMock   func()
		expectError bool
		expectNil   bool
	}{
		{
			name:   "successful find",
			userID: 1,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password_hash", "role"}).
					AddRow(1, time.Now(), time.Now(), "testuser", "hashedpassword", "admin")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name:   "not found",
			userID: 999,
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.FindByID(context.Background(), tt.userID, nil)

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

func TestBaseRepository_Updates(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := &repository.BaseRepositoryImpl[entity.User]{DB: db}

	user := &entity.User{
		Base:         entity.Base{ID: 1},
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Role:         "admin",
	}

	updateData := entity.User{
		Username: "updateduser",
		Role:     "user",
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
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"username"=$2,"role"=$3 WHERE "id" = $4`)).
					WithArgs(sqlmock.AnyArg(), "updateduser", "user", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.Updates(context.Background(), user, updateData, nil)

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

func TestBaseRepository_Save(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := &repository.BaseRepositoryImpl[entity.User]{DB: db}

	user := &entity.User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Role:         "admin",
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "successful save",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser", "hashedpassword", "admin").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := repo.Save(context.Background(), user, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestBaseRepository_FindByTemplate(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := &repository.BaseRepositoryImpl[entity.User]{DB: db}

	template := &entity.User{
		Role: "admin",
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
		expectCount int
	}{
		{
			name: "successful find multiple",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password_hash", "role"}).
					AddRow(1, time.Now(), time.Now(), "user1", "hash1", "admin").
					AddRow(2, time.Now(), time.Now(), "user2", "hash2", "admin")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."role" = $1`)).
					WithArgs("admin").
					WillReturnRows(rows)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "no results",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password_hash", "role"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."role" = $1`)).
					WithArgs("admin").
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

func TestBaseRepository_FindOneByTemplate(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := &repository.BaseRepositoryImpl[entity.User]{DB: db}

	template := &entity.User{
		Username: "testuser",
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
		expectNil   bool
	}{
		{
			name: "successful find one",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password_hash", "role"}).
					AddRow(1, time.Now(), time.Now(), "testuser", "hashedpassword", "admin")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."username" = $1`)).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name: "not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."username" = $1`)).
					WithArgs("testuser").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
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
