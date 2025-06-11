package entity

import (
	"time"
)

type AuditLog struct {
	ID         int64                  `gorm:"primaryKey;autoIncrement"`
	TableName  string                 `gorm:"not null"`
	RecordID   int64                  `gorm:"not null"`
	Action     string                 `gorm:"not null"`
	DataBefore map[string]interface{} `gorm:"type:jsonb"`
	DataAfter  map[string]interface{} `gorm:"type:jsonb"`
	UserID     string                 `gorm:"not null"`
	IPAddress  string                 `gorm:"not null"`
	CreatedAt  time.Time              `gorm:"autoCreateTime;not null"`
}

const (
	AuditLogActionCreate = "create"
	AuditLogActionUpdate = "update"
	AuditLogActionDelete = "delete"
)
