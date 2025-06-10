package entity

import "time"

type Reimbursement struct {
	Base
	EmployeeID  string    `gorm:"not null;index"`
	Amount      int64     `gorm:"not null"`
	Date        time.Time `gorm:"not null"`
	Description string    `gorm:"not null"`
}
