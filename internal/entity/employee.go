package entity

type Employee struct {
	Base
	UserID     string `gorm:"not null;uniqueIndex"`
	BaseSalary int64  `gorm:"not null"`
}
