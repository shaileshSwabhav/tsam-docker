package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Degree contain field id, Degree name.
type Degree struct {
	TenantBase
	Name string `json:"name" gorm:"type:varchar(200)"`
}

// Validate validates all the fields of degree.
func (degree *Degree) Validate() error {

	// Degree name must be specified.
	if util.IsEmpty(degree.Name) {
		return errors.NewValidationError("Degree Name must be specified")
	}

	// Degree name maximum characters.
	if len(degree.Name) > 200 {
		return errors.NewValidationError("Degree Name can have maximum 200 characters")
	}

	return nil
}
