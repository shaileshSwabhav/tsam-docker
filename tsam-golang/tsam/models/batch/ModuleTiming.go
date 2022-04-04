package batch

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// ModuleTiming will consist of day and its timings for a module.
type ModuleTiming struct {
	general.TenantBase
	BatchID       uuid.UUID `json:"batchID" gorm:"type:varchar(36)"`
	BatchModuleID uuid.UUID `json:"batchModuleID" gorm:"type:varchar(36)"`
	DayID         uuid.UUID `json:"dayID" gorm:"type:varchar(36)"`
	FromTime      string    `json:"fromTime" gorm:"type:time"`
	ToTime        string    `json:"toTime" gorm:"type:time"`
	ModuleID      uuid.UUID `json:"moduleID" gorm:"type:varchar(36)"`
	FacultyID     uuid.UUID `json:"facultyID" gorm:"type:varchar(36)"`

	// BatchModuleID uuid.UUID `json:"batchModuleID" gorm:"type:varchar(36)"`
	// Day       *general.Day `json:"day" gorm:"type:foreignkey:DayID;association_autocreate:false;association_autoupdate:false"`
}

// TableName overrides name of the table
func (*ModuleTiming) TableName() string {
	return "batch_module_timings"
}

// Validate will check if all fields are valid.
func (t *ModuleTiming) Validate() error {

	// if !util.IsUUIDValid(t.BatchID) {
	// 	return errors.NewValidationError("invalid batch ID")
	// }

	// if !util.IsUUIDValid(t.ModuleID) {
	// 	return errors.NewValidationError("invalid module ID")
	// }

	// if !util.IsUUIDValid(t.FacultyID) {
	// 	return errors.NewValidationError("invalid faculty ID")
	// }

	if !util.IsUUIDValid(t.DayID) {
		return errors.NewValidationError("invalid day ID")
	}

	if util.IsEmpty(t.FromTime) {
		return errors.NewValidationError("start time for module must be specified")
	}

	if util.IsEmpty(t.ToTime) {
		return errors.NewValidationError("end time for module must be specified")
	}

	return nil
}

// ModuleTimingDTO will consist of day and its timings for a module.
type ModuleTimingDTO struct {
	general.BaseDTO
	BatchModuleID uuid.UUID   `json:"batchModuleID"`
	Day           general.Day `json:"day" gorm:"type:foreignkey:DayID"`
	FromTime      string      `json:"fromTime"`
	ToTime        string      `json:"toTime"`
	ModuleID      uuid.UUID   `json:"-"`
	FacultyID     uuid.UUID   `json:"-"`
	BatchID       uuid.UUID   `json:"-"`
	DayID         uuid.UUID   `json:"-"`
}

// TableName overrides name of the table
func (*ModuleTimingDTO) TableName() string {
	return "batch_module_timings"
}
