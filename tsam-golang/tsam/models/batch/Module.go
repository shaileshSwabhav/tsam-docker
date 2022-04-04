package batch

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

// Module will store data related to batch_modules.
type Module struct {
	general.TenantBase
	BatchID          uuid.UUID      `json:"batchID" gorm:"type:varchar(36)"`
	ModuleID         uuid.UUID      `json:"moduleID" gorm:"type:varchar(36)"`
	FacultyID        uuid.UUID      `json:"facultyID" gorm:"type:varchar(36)"`
	StartDate        *string        `json:"startDate" gorm:"type:date"`
	EstimatedEndDate *string        `json:"estimatedEndDate" gorm:"type:date"`
	Order            uint           `json:"order" gorm:"type:INT;DEFAULT:false"`
	IsCompleted      *bool          `json:"isCompleted" gorm:"type:TINYINT"`
	ModuleTiming     []ModuleTiming `json:"moduleTimings" gorm:"foreignkey:BatchModuleID;ASSOCIATION_AUTOUPDATE:false;"`

	// Timing      []ModuleTimingDTO `json:"moduleTimings"`
}

// TableName overrides name of the table
func (*Module) TableName() string {
	return "batch_modules"
}

// Validate will verify if all fields of module are valid.
func (m *Module) Validate() error {

	if !util.IsUUIDValid(m.FacultyID) {
		return errors.NewValidationError("Faculty ID must be specified")
	}

	// if util.IsEmpty(m.StartDate) {
	// 	return errors.NewValidationError("Start date must be specified")
	// }

	if m.Order <= 0 {
		return errors.NewValidationError("Order must be specified.")
	}

	if m.ModuleTiming != nil {
		for _, timing := range m.ModuleTiming {
			err := timing.Validate()
			if err != nil {
				return err
			}
		}
	}

	// if m.IsCompleted == nil {
	// 	return errors.NewValidationError("Is completed field must be specified")
	// }

	return nil
}

// ModuleDTO will store data related to batch_modules.
type ModuleDTO struct {
	general.BaseDTO
	BatchID          uuid.UUID          `json:"-"`
	ModuleID         uuid.UUID          `json:"-"`
	FacultyID        uuid.UUID          `json:"-"`
	Module           course.ModuleDTO   `json:"module" gorm:"foreignkey:ModuleID"`
	Faculty          list.Faculty       `json:"faculty"`
	Order            uint               `json:"order"`
	IsCompleted      *bool              `json:"isCompleted"`
	StartDate        *string            `json:"startDate"`
	EstimatedEndDate *string            `json:"estimatedEndDate"`
	ModuleTiming     []*ModuleTimingDTO `json:"moduleTimings" gorm:"foreignkey:BatchModuleID"`
}

// TableName overrides name of the table
func (*ModuleDTO) TableName() string {
	return "batch_modules"
}
