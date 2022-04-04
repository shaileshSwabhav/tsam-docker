package test

import (
	"github.com/techlabs/swabhav/tsam/util"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	common "github.com/techlabs/swabhav/tsam/models/general"
)

// Option have Option & Status
type Option struct {
	common.TenantBase
	Option     *string   `json:"option"`
	Status     *bool     `json:"status" gorm:"type:tinyint(1)"`
	QuestionID uuid.UUID `json:"-" gorm:"type:varchar(36);"`
}

// ValidateOption Check Option data & Return Error On Invalid Data
func (option *Option) ValidateOption() error {
	if util.IsNil(option.Option) {
		return errors.NewValidationError("Option must be specified")
	}

	if option.Status == nil {
		return errors.NewValidationError("Status must be specified")
	}
	return nil
}
