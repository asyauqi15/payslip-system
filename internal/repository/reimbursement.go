package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

//go:generate mockgen -destination=./mock/mock_reimbursement_repository.go -package=mock github.com/asyauqi15/payslip-system/internal/repository ReimbursementRepository
type ReimbursementRepository interface {
	BaseRepository[entity.Reimbursement]
}

type ReimbursementRepositoryImpl struct {
	BaseRepositoryImpl[entity.Reimbursement]
}

func NewReimbursementRepository(db *BaseRepositoryImpl[entity.Reimbursement]) ReimbursementRepository {
	return &ReimbursementRepositoryImpl{
		BaseRepositoryImpl: *db,
	}
}
