package entity

import (
	"context"
	"encoding/json"
	"reflect"

	"time"

	"github.com/asyauqi15/payslip-system/internal/constant"
	"gorm.io/gorm"
)

type Base struct {
	ID        int64     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`
}

func (b Base) BeforeUpdate(tx *gorm.DB) (err error) {
	// Get the current record from database before update
	if tx.Statement.Schema != nil {
		// Create a new instance of the same type
		modelType := tx.Statement.Schema.ModelType
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}

		currentRecord := reflect.New(modelType).Interface()

		// Find the current record by ID
		var recordID interface{}
		if field := tx.Statement.Schema.LookUpField("ID"); field != nil {
			recordID, _ = field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue)
		}

		if recordID != nil {
			err := tx.Where("id = ?", recordID).First(currentRecord).Error
			if err == nil {
				// Convert to map for storage
				dataBytes, _ := json.Marshal(currentRecord)
				var dataMap map[string]interface{}
				json.Unmarshal(dataBytes, &dataMap)

				// Store in context
				tx.Statement.Context = context.WithValue(tx.Statement.Context, "data_before", dataMap)
			}
		}
	}
	return nil
}

func (b Base) AfterUpdate(tx *gorm.DB) (err error) {
	// Get the data before from context
	dataBefore, _ := tx.Statement.Context.Value("data_before").(map[string]interface{})

	// Get the current data after update
	var dataAfter map[string]interface{}
	if tx.Statement.Dest != nil {
		dataBytes, _ := json.Marshal(tx.Statement.Dest)
		json.Unmarshal(dataBytes, &dataAfter)
	}

	// Find differences and exclude UpdatedAt
	differencesBefore := make(map[string]interface{})
	differencesAfter := make(map[string]interface{})

	if dataBefore != nil && dataAfter != nil {
		for key, valueBefore := range dataBefore {
			if key == "UpdatedAt" || key == "updated_at" {
				continue
			}
			if valueAfter, exists := dataAfter[key]; exists {
				// Convert to comparable types
				beforeStr, _ := json.Marshal(valueBefore)
				afterStr, _ := json.Marshal(valueAfter)

				if string(beforeStr) != string(afterStr) {
					differencesBefore[key] = valueBefore
					differencesAfter[key] = valueAfter
				}
			}
		}
	}

	// Only create audit log if there are actual changes
	if len(differencesBefore) > 0 {
		// Get record ID
		var recordID interface{}
		if field := tx.Statement.Schema.LookUpField("ID"); field != nil {
			recordID, _ = field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue)
		}

		// Get user context (with nil checks)
		var userID string
		var ipAddress string

		if userIDVal := tx.Statement.Context.Value(constant.ContextKeyUserID); userIDVal != nil {
			userID, _ = userIDVal.(string)
		}
		if ipVal := tx.Statement.Context.Value(constant.ContextKeyIPAddress); ipVal != nil {
			ipAddress, _ = ipVal.(string)
		}

		// Create audit log entry
		auditLog := AuditLog{
			TableName:  tx.Statement.Table,
			RecordID:   recordID.(int64),
			Action:     AuditLogActionUpdate,
			DataBefore: differencesBefore,
			DataAfter:  differencesAfter,
			UserID:     userID,
			IPAddress:  ipAddress,
		}

		return tx.Create(&auditLog).Error
	}

	return nil
}

func (b Base) AfterCreate(tx *gorm.DB) (err error) {
	// Get record ID
	var recordID interface{}
	if field := tx.Statement.Schema.LookUpField("ID"); field != nil {
		recordID, _ = field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue)
	}

	// Get user context (with nil checks)
	var userID string
	var ipAddress string

	if userIDVal := tx.Statement.Context.Value(constant.ContextKeyUserID); userIDVal != nil {
		userID, _ = userIDVal.(string)
	}
	if ipVal := tx.Statement.Context.Value(constant.ContextKeyIPAddress); ipVal != nil {
		ipAddress, _ = ipVal.(string)
	}

	// Create audit log entry
	auditLog := AuditLog{
		TableName:  tx.Statement.Table,
		RecordID:   recordID.(int64),
		Action:     AuditLogActionCreate,
		DataBefore: nil,
		DataAfter:  nil,
		UserID:     userID,
		IPAddress:  ipAddress,
	}

	return tx.Create(&auditLog).Error
}
