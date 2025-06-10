package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Entity interface {
	AfterCreate(tx *gorm.DB) (err error)
	AfterUpdate(tx *gorm.DB) (err error)
	BeforeUpdate(tx *gorm.DB) (err error)
}

type BaseRepositoryImpl[T Entity] struct {
	DB *gorm.DB
}

type BaseRepository[T Entity] interface {
	Create(ctx context.Context, o *T, tx *gorm.DB) (*T, error)
	Updates(ctx context.Context, o *T, u T, tx *gorm.DB) (*T, error)
	Save(ctx context.Context, o *T, tx *gorm.DB) error
	FindByID(ctx context.Context, i uint, tx *gorm.DB) (*T, error)
	FindByTemplate(ctx context.Context, t *T, tx *gorm.DB) ([]T, error)
	FindOneByTemplate(ctx context.Context, o *T, tx *gorm.DB) (*T, error)
}

func (b *BaseRepositoryImpl[T]) UseTransaction(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}

	return b.DB
}

func (b *BaseRepositoryImpl[T]) Create(ctx context.Context, o *T, tx *gorm.DB) (*T, error) {
	conn := b.UseTransaction(tx)
	err := conn.WithContext(ctx).Omit(clause.Associations).Create(o).Error
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (b *BaseRepositoryImpl[T]) Updates(ctx context.Context, o *T, u T, tx *gorm.DB) (*T, error) {
	conn := b.UseTransaction(tx)
	err := conn.WithContext(ctx).Model(o).Updates(u).Error
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (b *BaseRepositoryImpl[T]) Save(ctx context.Context, o *T, tx *gorm.DB) error {
	conn := b.UseTransaction(tx)
	return conn.WithContext(ctx).Omit(clause.Associations).Save(o).Error
}

func (b *BaseRepositoryImpl[T]) FindByID(ctx context.Context, i uint, tx *gorm.DB) (*T, error) {
	conn := b.UseTransaction(tx)
	var result T
	err := conn.WithContext(ctx).First(&result, i).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *BaseRepositoryImpl[T]) FindByTemplate(ctx context.Context, t *T, tx *gorm.DB) ([]T, error) {
	conn := b.UseTransaction(tx)
	var results []T
	err := conn.WithContext(ctx).Where(t).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (b *BaseRepositoryImpl[T]) FindOneByTemplate(ctx context.Context, o *T, tx *gorm.DB) (*T, error) {
	conn := b.UseTransaction(tx)
	var result T
	err := conn.WithContext(ctx).Where(o).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
