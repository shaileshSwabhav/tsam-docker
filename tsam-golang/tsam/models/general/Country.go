package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Country contains information Example id, country name.
type Country struct {
	TenantBase
	Name string `json:"name" gorm:"type:varchar(50)"`
}

// Validate returns error if country is invalid.
func (country *Country) Validate() error {

	// Name must be specified.
	if util.IsEmpty(country.Name) {
		return errors.NewValidationError("Country name must be specified")
	}

	// Name must be specified.
	if !util.ValidateStringWithSpace(country.Name) {
		return errors.NewValidationError("Country name can have only alphabets and space")
	}

	// Name maximum characters.
	if len(country.Name) > 50 {
		return errors.NewValidationError("Country name cam have maximum 50 characters")
	}

	return nil
}
