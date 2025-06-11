package payroll

import (
	"context"
	"fmt"
	"time"

	"github.com/asyauqi15/payslip-system/internal/entity"
	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	"github.com/asyauqi15/payslip-system/pkg/logger"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"gorm.io/gorm"
)

type PayrollTotals struct {
	TotalPayroll       int64
	TotalReimbursement int64
	TotalOvertime      int64
}

func (u *UsecaseImpl) RunPayroll(ctx context.Context, req v1.PostAdminPayrollsJSONRequestBody) error {
	attendancePeriod, err := u.validateAndGetAttendancePeriod(ctx, int64(req.AttendancePeriodId))
	if err != nil {
		return err
	}

	if err := u.validatePayrollNotExists(ctx, int64(req.AttendancePeriodId)); err != nil {
		return err
	}

	employees, err := u.getEmployees(ctx)
	if err != nil {
		return err
	}

	totalWorkingDays := u.calculateWorkingDays(attendancePeriod.StartDate, attendancePeriod.EndDate)

	createdPayroll, err := u.createPayrollRecord(ctx, int64(req.AttendancePeriodId), len(employees))
	if err != nil {
		return err
	}

	payrollTotals, err := u.processAllEmployeePayslips(ctx, employees, createdPayroll.ID, attendancePeriod, totalWorkingDays)
	if err != nil {
		return err
	}

	if err := u.updatePayrollTotals(ctx, createdPayroll, payrollTotals); err != nil {
		return err
	}

	u.logPayrollSuccess(ctx, createdPayroll.ID, int64(req.AttendancePeriodId), len(employees), payrollTotals.TotalPayroll)
	return nil
}

func (u *UsecaseImpl) processEmployeePayslip(ctx context.Context, employee entity.Employee, payrollID int64, attendancePeriod *entity.AttendancePeriod, totalWorkingDays int) (*entity.Payslip, error) {
	// Count employee attendance
	attendanceCount, err := u.countEmployeeAttendance(ctx, employee.ID, attendancePeriod.StartDate, attendancePeriod.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to count attendance for employee %d: %w", employee.ID, err)
	}

	// Calculate prorated salary based on attendance
	proratedSalary := u.calculateProratedSalary(employee.BaseSalary, attendanceCount, totalWorkingDays)

	// Calculate overtime pay
	overtimeHours, overtimePay, err := u.calculateOvertimePay(ctx, employee.ID, employee.BaseSalary, attendancePeriod.StartDate, attendancePeriod.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate overtime for employee %d: %w", employee.ID, err)
	}

	// Calculate reimbursement total
	reimbursementTotal, err := u.calculateReimbursementTotal(ctx, employee.ID, attendancePeriod.StartDate, attendancePeriod.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate reimbursement for employee %d: %w", employee.ID, err)
	}

	// Calculate total take home
	totalTakeHome := proratedSalary + overtimePay + reimbursementTotal

	// Create payslip
	payslip := &entity.Payslip{
		EmployeeID:         employee.ID,
		PayrollID:          payrollID,
		BaseSalary:         employee.BaseSalary,
		AttendanceCount:    attendanceCount,
		TotalWorkingDays:   totalWorkingDays,
		ProratedSalary:     proratedSalary,
		OvertimeTotalHours: overtimeHours,
		OvertimeTotalPay:   overtimePay,
		ReimbursementTotal: reimbursementTotal,
		TotalTakeHome:      totalTakeHome,
	}

	createdPayslip, err := u.payslipRepo.Create(ctx, payslip, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create payslip for employee %d: %w", employee.ID, err)
	}

	return createdPayslip, nil
}

func (u *UsecaseImpl) calculateWorkingDays(startDate, endDate time.Time) int {
	count := 0
	current := startDate

	for current.Before(endDate) || current.Equal(endDate) {
		// Only count weekdays (Monday to Friday)
		if current.Weekday() >= time.Monday && current.Weekday() <= time.Friday {
			count++
		}
		current = current.AddDate(0, 0, 1)
	}

	return count
}

func (u *UsecaseImpl) countEmployeeAttendance(ctx context.Context, employeeID int64, startDate, endDate time.Time) (int, error) {
	// Use repository's count method for efficient counting
	count, err := u.attendanceRepo.CountAttendanceInPeriod(ctx, employeeID, startDate, endDate, nil)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (u *UsecaseImpl) calculateProratedSalary(baseSalary int64, attendanceCount, totalWorkingDays int) int64 {
	if totalWorkingDays == 0 {
		return 0
	}
	return baseSalary * int64(attendanceCount) / int64(totalWorkingDays)
}

func (u *UsecaseImpl) calculateOvertimePay(ctx context.Context, employeeID, baseSalary int64, startDate, endDate time.Time) (int, int64, error) {
	// Get all overtime records for the employee in the period
	overtimes, err := u.overtimeRepo.FindByTemplate(ctx, &entity.Overtime{
		EmployeeID: employeeID,
	}, nil)
	if err != nil {
		return 0, 0, err
	}

	totalHours := 0
	totalPay := int64(0)

	// Calculate hourly rate (assuming 8 hours per day, 22 working days per month)
	hourlyRate := baseSalary / (22 * 8)

	for _, overtime := range overtimes {
		// Check if overtime is within the period
		if overtime.StartAt.After(startDate.AddDate(0, 0, -1)) && overtime.StartAt.Before(endDate.AddDate(0, 0, 1)) {
			duration := overtime.EndAt.Sub(overtime.StartAt)
			hours := int(duration.Hours())

			// Overtime is paid twice the hourly rate
			overtimePay := int64(hours) * hourlyRate * 2

			totalHours += hours
			totalPay += overtimePay
		}
	}

	return totalHours, totalPay, nil
}

func (u *UsecaseImpl) calculateReimbursementTotal(ctx context.Context, employeeID int64, startDate, endDate time.Time) (int64, error) {
	// Get all reimbursement records for the employee in the period
	reimbursements, err := u.reimbursementRepo.FindByTemplate(ctx, &entity.Reimbursement{
		EmployeeID: employeeID,
	}, nil)
	if err != nil {
		return 0, err
	}

	total := int64(0)
	for _, reimbursement := range reimbursements {
		// Check if reimbursement is within the period
		if reimbursement.Date.After(startDate.AddDate(0, 0, -1)) && reimbursement.Date.Before(endDate.AddDate(0, 0, 1)) {
			total += reimbursement.Amount
		}
	}

	return total, nil
}

func (u *UsecaseImpl) validateAndGetAttendancePeriod(ctx context.Context, attendancePeriodID int64) (*entity.AttendancePeriod, error) {
	attendancePeriod, err := u.attendancePeriodRepo.FindByID(ctx, uint(attendancePeriodID), nil)
	if err != nil {
		logger.Error(ctx, "failed to find attendance period", "id", attendancePeriodID, "error", err)
		return nil, httppkg.NewInternalServerError("failed to find attendance period")
	}
	if attendancePeriod == nil {
		return nil, httppkg.NewNotFoundError("attendance period not found")
	}
	return attendancePeriod, nil
}

func (u *UsecaseImpl) validatePayrollNotExists(ctx context.Context, attendancePeriodID int64) error {
	existingPayroll, err := u.payrollRepo.FindOneByTemplate(ctx, &entity.Payroll{
		AttendancePeriodID: attendancePeriodID,
	}, nil)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(ctx, "failed to check existing payroll", "attendance_period_id", attendancePeriodID, "error", err)
		return httppkg.NewInternalServerError("failed to check existing payroll")
	}
	if existingPayroll != nil {
		return httppkg.NewBadRequestError("payroll already exists for this attendance period")
	}
	return nil
}

func (u *UsecaseImpl) getEmployees(ctx context.Context) ([]entity.Employee, error) {
	employees, err := u.employeeRepo.FindByTemplate(ctx, &entity.Employee{}, nil)
	if err != nil {
		logger.Error(ctx, "failed to get employees", "error", err)
		return nil, httppkg.NewInternalServerError("failed to get employees")
	}

	if len(employees) == 0 {
		return nil, httppkg.NewBadRequestError("no employees found")
	}

	return employees, nil
}

func (u *UsecaseImpl) createPayrollRecord(ctx context.Context, attendancePeriodID int64, employeeCount int) (*entity.Payroll, error) {
	payroll := &entity.Payroll{
		AttendancePeriodID: attendancePeriodID,
		TotalEmployees:     int64(employeeCount),
		TotalReimbursement: 0,
		TotalOvertime:      0,
		TotalPayroll:       0,
	}

	createdPayroll, err := u.payrollRepo.Create(ctx, payroll, nil)
	if err != nil {
		logger.Error(ctx, "failed to create payroll", "error", err)
		return nil, httppkg.NewInternalServerError("failed to create payroll")
	}

	return createdPayroll, nil
}

func (u *UsecaseImpl) processAllEmployeePayslips(ctx context.Context, employees []entity.Employee, payrollID int64, attendancePeriod *entity.AttendancePeriod, totalWorkingDays int) (*PayrollTotals, error) {
	totals := &PayrollTotals{}

	for _, employee := range employees {
		payslip, err := u.processEmployeePayslip(ctx, employee, payrollID, attendancePeriod, totalWorkingDays)
		if err != nil {
			logger.Error(ctx, "failed to process employee payslip", "employee_id", employee.ID, "error", err)
			return nil, fmt.Errorf("failed to process employee %d payslip: %w", employee.ID, err)
		}

		totals.TotalPayroll += payslip.TotalTakeHome
		totals.TotalReimbursement += payslip.ReimbursementTotal
		totals.TotalOvertime += payslip.OvertimeTotalPay
	}

	return totals, nil
}

func (u *UsecaseImpl) updatePayrollTotals(ctx context.Context, payroll *entity.Payroll, totals *PayrollTotals) error {
	updatedPayroll := entity.Payroll{
		TotalReimbursement: totals.TotalReimbursement,
		TotalOvertime:      totals.TotalOvertime,
		TotalPayroll:       totals.TotalPayroll,
	}

	_, err := u.payrollRepo.Updates(ctx, payroll, updatedPayroll, nil)
	if err != nil {
		logger.Error(ctx, "failed to update payroll totals", "payroll_id", payroll.ID, "error", err)
		return httppkg.NewInternalServerError("failed to update payroll totals")
	}

	return nil
}

func (u *UsecaseImpl) logPayrollSuccess(ctx context.Context, payrollID int64, attendancePeriodID int64, employeeCount int, totalPayout int64) {
	logger.Info(ctx, "payroll generated successfully",
		"payroll_id", payrollID,
		"attendance_period_id", attendancePeriodID,
		"total_employees", employeeCount,
		"total_payout", totalPayout)
}
