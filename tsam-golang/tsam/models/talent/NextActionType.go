package talent

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// NextActionType contains details about which type of action is added to talent.
type NextActionType struct {
	general.TenantBase
	Type string `json:"type" gorm:"type:varchar(100)"`
}

// TableName will create a table for model NextActionType with name talent_next_action_types.
func (*NextActionType) TableName() string {
	return "talent_next_action_types"
}

// Validate will validate the fields of next action type.
func (actionType *NextActionType) Validate() error {
	if util.IsEmpty(actionType.Type) {
		return errors.NewValidationError("Action type must be specified.")
	}
	return nil
}
