package general

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// University contains info about university.
type University struct {
	TenantBase
	UniversityName string    `json:"universityName"`
	CountryID      uuid.UUID `json:"-" gorm:"type:varchar(36)"`
	Country        *Country  `json:"country" gorm:"association_autocreate:false;association_autoupdate:false"`
}

// Validate returns error if UniversityName or id fields are not initialized with proper values.
func (university *University) Validate() error {
	if util.IsEmpty(university.UniversityName) {
		return errors.NewValidationError("University name should not be empty.")
	}
	if university.Country == nil {
		return errors.NewValidationError("Country ID or name should be specified.")
	}
	return nil
}
