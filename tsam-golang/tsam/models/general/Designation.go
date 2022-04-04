package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Designation contain single designation detail example id, position.
type Designation struct {
	TenantBase
	Position string `json:"position" gorm:"type:varchar(200)"`
}

// Validate will validate fields of designation.
func (designation *Designation) Validate() error {

	// Position must be specified.
	if util.IsEmpty(designation.Position) {
		return errors.NewValidationError("Position must be specified")
	}

	// Position must be specified.
	if len(designation.Position) > 200 {
		return errors.NewValidationError("Position can have maximum 200 characters")
	}

	return nil
}
