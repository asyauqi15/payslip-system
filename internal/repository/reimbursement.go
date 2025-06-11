package repository

import "github.com/asyauqi15/payslip-system/internal/entity"

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
