package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONMap is a type that can be stored as JSONB in PostgreSQL
type JSONMap map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

type AuditLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	TableName  string    `gorm:"not null"`
	RecordID   int64     `gorm:"not null"`
	Action     string    `gorm:"not null"`
	DataBefore JSONMap   `gorm:"type:jsonb"`
	DataAfter  JSONMap   `gorm:"type:jsonb"`
	UserID     string    `gorm:"not null"`
	IPAddress  string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime;not null"`
}

const (
	AuditLogActionCreate = "create"
	AuditLogActionUpdate = "update"
	AuditLogActionDelete = "delete"
)
