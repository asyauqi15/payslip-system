package entity

import (
	"time"

	httppkg "github.com/asyauqi15/payslip-system/pkg/http"
	"gorm.io/gorm"
)

type AttendancePeriod struct {
	Base
	StartDate time.Time `gorm:"not null;index"`
	EndDate   time.Time `gorm:"not null;index"`
}

func (ap AttendancePeriod) BeforeCreate(tx *gorm.DB) (err error) {
	var count int64
	err = tx.Model(&AttendancePeriod{}).
		Where("start_date <= ? AND end_date >= ?", ap.EndDate, ap.StartDate).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return httppkg.NewUnprocessableEntityError("attendance period overlaps with an existing period")
	}

	return nil
}

func (ap AttendancePeriod) BeforeUpdate(tx *gorm.DB) (err error) {
	var count int64
	err = tx.Model(&AttendancePeriod{}).
		Where("id != ? AND start_date =< ? AND end_date >= ?", ap.ID, ap.EndDate, ap.StartDate).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return httppkg.NewUnprocessableEntityError("attendance period overlaps with an existing period")
	}

	return ap.Base.BeforeUpdate(tx)
}
