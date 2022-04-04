package general

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Outcome contains different outcomes w.r.t. purposes for call records.
type Outcome struct {
	TenantBase
	PurposeID uuid.UUID `json:"purposeID" gorm:"type:varchar(36)"`
	Outcome   string    `json:"outcome" gorm:"type:varchar(50)"`
}

// Validate validates fields of Outcome.
func (outcome *Outcome) Validate() error {
	if util.IsEmpty(outcome.Outcome) {
		return errors.NewValidationError("Outcome must be specified")
	}
	return nil
}
