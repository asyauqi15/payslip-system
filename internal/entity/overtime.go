package entity

import "time"

type Overtime struct {
	Base
	EmployeeID  int64     `gorm:"not null;index"`
	StartAt     time.Time `gorm:"not null;index"`
	EndAt       time.Time `gorm:"not null"`
	Description string    `gorm:"not null"`
}
