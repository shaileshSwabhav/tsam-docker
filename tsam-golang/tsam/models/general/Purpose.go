package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Purpose defines fields for purpose which is used in call records.
type Purpose struct {
	TenantBase
	Purpose     string `json:"purpose" gorm:"varchar(50)"`
	PurposeType string `json:"purposeType" gorm:"varchar(20)"`
}

// Validate checks all the fields of purpose table.
func (purpose *Purpose) Validate() error {
	if util.IsEmpty(purpose.Purpose) {
		return errors.NewValidationError("Purpose must be specified")
	}
	if util.IsEmpty(purpose.PurposeType) {
		return errors.NewValidationError("Purpose type must be specified")
	}
	return nil
}
