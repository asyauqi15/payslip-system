package entity

import (
	"context"
	"github.com/asyauqi15/payslip-system/internal/constant"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`
}

func (b *Base) BeforeUpdate(tx *gorm.DB) (err error) {
	// Capture the state before the update
	dataBefore := make(map[string]any)
	if err := tx.Statement.Preload("").Find(dataBefore).Error; err != nil {
		return err
	}
	tx.Statement.Context = context.WithValue(tx.Statement.Context, "data_before", dataBefore)
	return nil
}

func (b *Base) AfterUpdate(tx *gorm.DB) (err error) {
	// Capture the state after the update
	dataAfter := make(map[string]any)
	if err := tx.Statement.Preload("").Find(dataAfter).Error; err != nil {
		return err
	}

	// Retrieve the dataBefore from the context
	dataBefore, _ := tx.Statement.Context.Value("data_before").(map[string]any)

	// Find differences and exclude UpdatedAt
	differencesBefore := make(map[string]any)
	differencesAfter := make(map[string]any)
	for key, valueBefore := range dataBefore {
		if key == "UpdatedAt" {
			continue
		}
		if valueAfter, exists := dataAfter[key]; exists && valueBefore != valueAfter {
			differencesBefore[key] = valueBefore
			differencesAfter[key] = valueAfter
		}
	}

	// Create an audit log entry
	return tx.Create(AuditLog{
		TableName:  tx.Statement.Table,
		RecordID:   tx.Statement.Schema.PrimaryFields[0].DBName,
		Action:     AuditLogActionUpdate,
		DataBefore: differencesBefore,
		DataAfter:  differencesAfter,
		UserID:     tx.Statement.Context.Value(constant.ContextKeyUserID).(string),
		IPAddress:  tx.Statement.Context.Value(constant.ContextKeyIPAddress).(string),
	}).Error
}

func (b *Base) AfterCreate(tx *gorm.DB) (err error) {
	// Create an audit log entry
	return tx.Create(AuditLog{
		TableName:  tx.Statement.Table,
		RecordID:   tx.Statement.Schema.PrimaryFields[0].DBName,
		Action:     AuditLogActionCreate,
		DataBefore: nil,
		DataAfter:  nil,
		UserID:     tx.Statement.Context.Value(constant.ContextKeyUserID).(string),
		IPAddress:  tx.Statement.Context.Value(constant.ContextKeyIPAddress).(string),
	}).Error
}
