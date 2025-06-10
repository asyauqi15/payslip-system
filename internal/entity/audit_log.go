package entity

import "time"

type AuditLog struct {
	ID         string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	TableName  string `gorm:"not null"`
	RecordID   string `gorm:"not null"`
	Action     string `gorm:"not null"`
	DataBefore map[string]any
	DataAfter  map[string]any
	UserID     string    `gorm:"not null"`
	IPAddress  string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime;not null"`
}

const (
	AuditLogActionCreate = "create"
	AuditLogActionUpdate = "update"
)
