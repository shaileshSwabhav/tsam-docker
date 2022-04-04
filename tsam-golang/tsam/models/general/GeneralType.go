package general

import (
	"github.com/techlabs/swabhav/tsam/errors"

	"github.com/techlabs/swabhav/tsam/util"
)

// CommonType contains ID, Key, Value which can be used by any module.
type CommonType struct {
	TenantBase
	Key   int8   `json:"key" gorm:"type:tinyint(3)"`
	Type  string `json:"type" gorm:"type:varchar(50)"`
	Value string `json:"value" gorm:"type:varchar(50)"`
}

// TableName will make the model refer "general_types" table.
func (*CommonType) TableName() string {
	return "general_types"
}

// Validate returns error if GeneralType is invalid.
func (generalType *CommonType) Validate() error {
	// Check if 0 key should be allowed.
	if generalType.Key <= 0 {
		return errors.NewValidationError("Key must be specified")
	}
	if util.IsEmpty(generalType.Type) {
		return errors.NewValidationError("Type must be specified")
	}
	if util.IsEmpty(generalType.Value) {
		return errors.NewValidationError("Value must be specified")
	}
	return nil
}
