package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Role will contain different roles names which will determine what access will be given to which user.
type Role struct {
	TenantBase
	RoleName   string `json:"roleName" example:"administrator" gorm:"varchar(50)"`
	Level      uint8  `json:"level" gorm:"SMALLINT(2)"`
	IsEmployee bool   `json:"isEmployee"`
}

// Validate Validates role name.
func (role *Role) Validate() error {
	if util.IsEmpty(role.RoleName) {
		return errors.NewValidationError("Rolename must be specified")
	}
	if role.Level == 0 {
		return errors.NewValidationError("Level of role must be specified")
	}
	return nil
}
