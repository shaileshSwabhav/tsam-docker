package admin

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Timesheet consists of the fields of weekly timesheet.
type Timesheet struct {
	general.TenantBase
	Date         string              `json:"date" gorm:"type:date"`
	DepartmentID uuid.UUID           `json:"departmentID" gorm:"type:varchar(36)"` // should not be visible in UI (auto)
	CredentialID uuid.UUID           `json:"credentialID" gorm:"type:varchar(36)"`
	IsOnLeave    bool                `json:"isOnLeave"`
	Activities   []TimesheetActivity `json:"activities" gorm:"foreignkey:TimesheetID;association_autoupdate:false"`
}

// TimesheetDTO consists of the fields of weekly timesheet.
type TimesheetDTO struct {
	general.BaseDTO
	Date         string                  `json:"date"`
	DepartmentID uuid.UUID               `json:"departmentID"` // should not be visible in UI (auto)
	CredentialID uuid.UUID               `json:"credentialID"`
	IsOnLeave    bool                    `json:"isOnLeave"`
	Activities   []*TimesheetActivityDTO `json:"activities" gorm:"foreignkey:TimesheetID"`
}

// TimesheetHeader consists of fields whose values might be put in headers.
type TimesheetHeader struct {
	Limit      int
	Offset     int
	TotalCount int
	TotalHours float32
	FreeHours  float32
}

// ValidateTimesheet fields of Timesheet.
func (timesheet *Timesheet) ValidateTimesheet() error {
	if util.IsEmpty(timesheet.Date) {
		return errors.NewValidationError("Date must be specified")
	}
	if !util.IsUUIDValid(timesheet.CredentialID) {
		return errors.NewValidationError("Credential id must be specified")
	}
	if !util.IsUUIDValid(timesheet.DepartmentID) {
		return errors.NewValidationError("Department id must be specified")
	}

	if timesheet.Activities == nil && !timesheet.IsOnLeave {
		return errors.NewValidationError("Timesheet activity must be specified")
	}

	if timesheet.IsOnLeave {
		timesheet.Activities = nil
	}

	for _, activity := range timesheet.Activities {
		err := activity.ValidateTimesheetActivity()
		if err != nil {
			return errors.NewValidationError(err.Error())
		}
		if activity.NextEstimatedDate != nil {
			err = timesheet.ValidateNextEstimateDate(*activity.NextEstimatedDate)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateNextEstimateDate is a check to see if next estimated date is a future date.
func (timesheet *Timesheet) ValidateNextEstimateDate(nextEstimtedDate string) error {
	date, err := time.Parse("2006-01-02", timesheet.Date)
	if err != nil {
		return errors.NewValidationError("Invalid date")
	}
	nextEstDate, err := time.Parse("2006-01-02", nextEstimtedDate)
	if err != nil {
		return errors.NewValidationError("Invalid next estimated date")
	}
	if nextEstDate.Before(date) || nextEstDate == date {
		return errors.NewValidationError("Next estimated date cannot be on or before the actual date")
	}
	return nil
}

// TableName will refer it to timesheets.
func (*TimesheetDTO) TableName() string {
	return "timesheets"
}
