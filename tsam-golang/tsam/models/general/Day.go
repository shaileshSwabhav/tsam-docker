package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Day contains information on week days.
type Day struct {
	TenantBase
	Day   string `json:"day" gorm:"type:varchar(15)"`
	Order uint   `json:"order" gorm:"type:tinyint(1)"`
	Type  string `json:"type" gorm:"type:varchar(20)"`
}

// ValidateDay validates day
func (d *Day) ValidateDay() error {
	if util.IsEmpty(d.Day) {
		return errors.NewValidationError("Day must be specified")
	}
	if d.Order <= 0 {
		return errors.NewValidationError("Order should be 1 or greater")
	}
	if util.IsEmpty(d.Type) {
		return errors.NewValidationError("Type must be specified")
	}
	return nil
}
