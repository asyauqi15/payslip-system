package entity

type Payroll struct {
	Base
	AttendancePeriodID int64 `gorm:"not null;uniqueIndex"`
	TotalEmployees     int64 `gorm:"column:employees_count;not null"`
	TotalReimbursement int64 `gorm:"not null"`
	TotalOvertime      int64 `gorm:"not null"`
	TotalPayroll       int64 `gorm:"not null"`
}
