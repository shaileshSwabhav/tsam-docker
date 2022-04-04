package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"

	uuid "github.com/satori/go.uuid"
)

// City contains information about city.
type City struct {
	TenantBase
	Name    string    `json:"name" gorm:"type:varchar(50)"`
	StateID uuid.UUID `json:"stateID" gorm:"type:varchar(36)"`
}

// Validate returns error if city is invalid.
func (city *City) Validate() error {
	if util.IsEmpty(city.Name) {
		return errors.NewValidationError("City name must be specified")
	}
	if len(city.Name) > 50 {
		return errors.NewValidationError("City name can have maximum 50 characters")
	}
	if !util.IsUUIDValid(city.StateID) {
		return errors.NewValidationError("StateID must be specified")
	}
	return nil
}
