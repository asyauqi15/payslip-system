package entity

type Attendance struct {
	Base
	EmployeeID   string `gorm:"not null;index"`
	ClockInTime  string `gorm:"not null"`
	ClockOutTime string
}
