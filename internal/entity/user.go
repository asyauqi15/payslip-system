package entity

type User struct {
	Base
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null"`
	Status       int    `gorm:"not null;default:1"`
}

const (
	UserRoleAdmin   = "admin"
	UserRoleDefault = "default"
)

const (
	UserStatusActive   = 1
	UserStatusInactive = 0
)
