package entity

type User struct {
	Base
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null"`
}

const (
	UserRoleAdmin   = "admin"
	UserRoleDefault = "default"
)

const (
	UserStatusActive   = 1
	UserStatusInactive = 0
)
