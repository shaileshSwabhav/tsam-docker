package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Source contain Name and ID.
type Source struct {
	TenantBase
	Name        string `json:"name" gorm:"type:varchar(40)"`
	URLName     string `json:"urlName" gorm:"type:varchar(50)"`
	Description string `json:"description" gorm:"type:varchar(150)"`
	// Removed fields
	// TalentID    uuid.UUID `gorm:"type:varchar(36)"`
}

// Validate validates all the fields.
func (source *Source) Validate() error {
	if util.IsEmpty(source.Name) {
		return errors.NewValidationError("Source name must be specified")
	}
	// if util.IsEmpty(source.URLName) {
	// 	return errors.NewValidationError("Url name must be specified")
	// }
	if util.IsEmpty(source.Description) {
		return errors.NewValidationError("Description must be specified")
	}
	return nil
}
