package entity

type Payslip struct {
	Base
	EmployeeID         int64 `gorm:"not null;index"`
	PayrollID          int64 `gorm:"not null;index"`
	BaseSalary         int64 `gorm:"not null"`
	ProratedSalary     int64 `gorm:"not null"`
	OvertimePay        int64 `gorm:"not null"`
	ReimbursementTotal int64 `gorm:"not null"`
	TotalTakeHome      int64 `gorm:"not null"`
}
