package entity

type Attendance struct {
	Base
	EmployeeID   int64  `gorm:"not null;index"`
	ClockInTime  string `gorm:"not null"`
	ClockOutTime string
}
