package batch

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Timing will store the schedule for sessions of batch
type Timing struct {
	general.TenantBase
	BatchID  uuid.UUID    `json:"batchID" gorm:"type:varchar(36)"`
	Day      *general.Day `json:"day" gorm:"type:foreignkey:DayID;association_autocreate:false;association_autoupdate:false"`
	DayID    uuid.UUID    `json:"-" gorm:"type:varchar(36)"`
	FromTime *string      `json:"fromTime" gorm:"type:time"`
	ToTime   *string      `json:"toTime" gorm:"type:time"`

	// Day      *string   `json:"day" gorm:"type:varchar(10)"`
	// ModuleID  uuid.UUID    `json:"moduleID" gorm:"type:varchar(36)"`
	// FacultyID uuid.UUID    `json:"facultyID" gorm:"type:varchar(36)"`
}

// ValidateBatchTiming will check if all fields are valid.
func (batch *Timing) ValidateBatchTiming() error {

	if batch.Day == nil {
		return errors.NewValidationError("Session day must be specified")
	}

	if batch.FromTime != nil && util.IsEmpty(*batch.FromTime) {
		return errors.NewValidationError("Start time must be specified")
	}

	if batch.ToTime != nil && util.IsEmpty(*batch.ToTime) {
		return errors.NewValidationError("End time must be specified")
	}

	return nil
}

// TableName overrides name of the table
func (*Timing) TableName() string {
	return "batch_timing"
}
