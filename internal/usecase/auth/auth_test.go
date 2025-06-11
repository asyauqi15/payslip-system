package auth_test

import (
	"context"
	"testing"

	"github.com/asyauqi15/payslip-system/internal/entity"
	authmock "github.com/asyauqi15/payslip-system/internal/repository/mock"
	"github.com/asyauqi15/payslip-system/internal/usecase/auth"
	jwtauth "github.com/asyauqi15/payslip-system/pkg/jwt-auth"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestAuthUsecase_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := authmock.NewMockUserRepository(ctrl)
	jwtAuth := &jwtauth.JWTAuthentication{}

	usecase := auth.NewUsecase(mockUserRepo, jwtAuth)

	// Create a test password hash
	testPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		email       string
		password    string
		setupMock   func()
		expectError bool
	}{
		{
			name:     "successful authentication",
			email:    "admin@example.com",
			password: testPassword,
			setupMock: func() {
				user := &entity.User{
					Base:         entity.Base{ID: 1},
					Username:     "admin@example.com",
					PasswordHash: string(hashedPassword),
					Role:         entity.UserRoleAdmin,
				}
				mockUserRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.User{Username: "admin@example.com"}, nil).
					Return(user, nil)
			},
			expectError: false,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: testPassword,
			setupMock: func() {
				mockUserRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.User{Username: "nonexistent@example.com"}, nil).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:     "invalid password",
			email:    "admin@example.com",
			password: "wrongpassword",
			setupMock: func() {
				user := &entity.User{
					Base:         entity.Base{ID: 1},
					Username:     "admin@example.com",
					PasswordHash: string(hashedPassword),
					Role:         entity.UserRoleAdmin,
				}
				mockUserRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.User{Username: "admin@example.com"}, nil).
					Return(user, nil)
			},
			expectError: true,
		},
		{
			name:     "repository error",
			email:    "admin@example.com",
			password: testPassword,
			setupMock: func() {
				mockUserRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.User{Username: "admin@example.com"}, nil).
					Return(nil, gorm.ErrInvalidDB)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := usecase.Auth(context.Background(), tt.email, tt.password)

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
				}
			}
		})
	}
}
