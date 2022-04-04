package general

import (
	"github.com/techlabs/swabhav/tsam/util"

	"github.com/techlabs/swabhav/tsam/errors"
)

// Technology contains computer languages.
type Technology struct {
	TenantBase
	Language string `json:"language" gorm:"type:varchar(50)"`
	Rating   uint8  `json:"rating" gorm:"type:int"`
}

// Validate returns error if language or id fields are not initialized with proper values.
func (technology *Technology) Validate() error {
	if util.IsEmpty(technology.Language) {
		return errors.NewValidationError("Technology Language must be specified")
	}
	if technology.Rating <= 0 || technology.Rating > 5 {
		return errors.NewValidationError("Technology rating must be between 1 to 5")
	}

	return nil
}
