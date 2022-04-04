package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"

	uuid "github.com/satori/go.uuid"
)

// State Contain information Example id, State name.
type State struct {
	TenantBase
	Name      string    `json:"name" gorm:"type:varchar(50)"`
	CountryID uuid.UUID `json:"countryID" gorm:"type:varchar(36)"`
}

// ValidateState returns error if state is invalid.
func (state *State) ValidateState() error {
	if util.IsEmpty(state.Name) {
		return errors.NewValidationError("Invalid State")
	}
	if !util.IsUUIDValid(state.CountryID) {
		return errors.NewValidationError("CountryID must be specified")
	}
	return nil
}
