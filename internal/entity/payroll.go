package entity

type Payroll struct {
	Base
	AttendancePeriodID string `gorm:"not null;uniqueIndex"`
	TotalEmployees     int64  `gorm:"not null"`
	TotalReimbursement int64  `gorm:"not null"`
	TotalOvertime      int64  `gorm:"not null"`
	TotalTakeHome      int64  `gorm:"not null"`
	TotalPayroll       int64  `gorm:"not null"`
}
