package entity

type Employee struct {
	Base
	UserID     int64 `gorm:"not null;uniqueIndex"`
	BaseSalary int64 `gorm:"not null"`
}
