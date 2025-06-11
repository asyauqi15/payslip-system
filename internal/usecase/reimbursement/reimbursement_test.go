package reimbursement_test

import (
	"context"
	"testing"
	"time"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/entity"
	"github.com/asyauqi15/payslip-system/internal/repository/mock"
	"github.com/asyauqi15/payslip-system/internal/usecase/reimbursement"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestReimbursementUsecase_SubmitReimbursement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReimbursementRepo := mock.NewMockReimbursementRepository(ctrl)
	mockEmployeeRepo := mock.NewMockEmployeeRepository(ctrl)

	usecase := reimbursement.NewUsecase(mockReimbursementRepo, mockEmployeeRepo)

	reimbursementDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		request      v1.ReimbursementRequest
		setupContext func() context.Context
		setupMock    func()
		expectError  bool
	}{
		{
			name: "successful reimbursement submission",
			request: v1.ReimbursementRequest{
				Amount:      50000,
				Date:        openapi_types.Date{Time: reimbursementDate},
				Description: "Business travel expenses",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}

				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)

				// Mock creating reimbursement
				mockReimbursementRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.Reimbursement{
						Base:        entity.Base{ID: 1},
						EmployeeID:  1,
						Amount:      50000,
						Date:        reimbursementDate,
						Description: "Business travel expenses",
					}, nil)
			},
			expectError: false,
		},
		{
			name: "user not authenticated",
			request: v1.ReimbursementRequest{
				Amount:      50000,
				Date:        openapi_types.Date{Time: reimbursementDate},
				Description: "Business travel expenses",
			},
			setupContext: func() context.Context {
				return context.Background() // No user ID in context
			},
			setupMock:   func() {},
			expectError: true,
		},
		{
			name: "employee not found",
			request: v1.ReimbursementRequest{
				Amount:      50000,
				Date:        openapi_types.Date{Time: reimbursementDate},
				Description: "Business travel expenses",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(999))
			},
			setupMock: func() {
				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(999)}, nil).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name: "repository error during employee lookup",
			request: v1.ReimbursementRequest{
				Amount:      50000,
				Date:        openapi_types.Date{Time: reimbursementDate},
				Description: "Business travel expenses",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(nil, gorm.ErrInvalidDB)
			},
			expectError: true,
		},
		{
			name: "repository error during reimbursement creation",
			request: v1.ReimbursementRequest{
				Amount:      50000,
				Date:        openapi_types.Date{Time: reimbursementDate},
				Description: "Business travel expenses",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}

				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)

				// Mock repository error during creation
				mockReimbursementRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(nil, gorm.ErrInvalidDB)
			},
			expectError: true,
		},
		{
			name: "zero amount reimbursement",
			request: v1.ReimbursementRequest{
				Amount:      0,
				Date:        openapi_types.Date{Time: reimbursementDate},
				Description: "Test reimbursement",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}

				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)

				// Mock creating reimbursement with zero amount
				mockReimbursementRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.Reimbursement{
						Base:        entity.Base{ID: 1},
						EmployeeID:  1,
						Amount:      0,
						Date:        reimbursementDate,
						Description: "Test reimbursement",
					}, nil)
			},
			expectError: false, // Zero amount might be valid in some business cases
		},
		{
			name: "negative amount reimbursement",
			request: v1.ReimbursementRequest{
				Amount:      -10000,
				Date:        openapi_types.Date{Time: reimbursementDate},
				Description: "Correction reimbursement",
			},
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), constant.ContextKeyUserID, int64(1))
			},
			setupMock: func() {
				employee := &entity.Employee{
					Base:       entity.Base{ID: 1},
					UserID:     1,
					BaseSalary: 5000000,
				}

				mockEmployeeRepo.EXPECT().
					FindOneByTemplate(gomock.Any(), &entity.Employee{UserID: int64(1)}, nil).
					Return(employee, nil)

				// Mock creating reimbursement with negative amount
				mockReimbursementRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), nil).
					Return(&entity.Reimbursement{
						Base:        entity.Base{ID: 1},
						EmployeeID:  1,
						Amount:      -10000,
						Date:        reimbursementDate,
						Description: "Correction reimbursement",
					}, nil)
			},
			expectError: false, // Negative amount might be valid for corrections
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			ctx := tt.setupContext()

			err := usecase.SubmitReimbursement(ctx, tt.request)

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
