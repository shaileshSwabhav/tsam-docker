package general

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// FeelingLevel will consist of different levels for each feeling.
type FeelingLevel struct {
	TenantBase
	FeelingID   uuid.UUID `json:"feelingID" gorm:"type:varchar(36)"`
	Description string    `json:"description" gorm:"type:varchar(200)"`
	LevelNumber uint      `json:"levelNumber" gorm:"INT"`
	// Feeling     Feeling   `json:"feeling" gorm:"foreignkey:FeelingID"`
}

// Validate will check if all fields are valid.
func (feelingLevel *FeelingLevel) Validate() error {

	// Description must be specified.
	if util.IsEmpty(feelingLevel.Description) {
		return errors.NewValidationError("Description must be specified")
	}

	// Description maximum characters.
	if len(feelingLevel.Description) > 200 {
		return errors.NewValidationError("Description can have maximum 200 characters")
	}

	// Level number must be greater than 0.
	if feelingLevel.LevelNumber <= 0 {
		return errors.NewValidationError("Level number must be specified and cannot be zero")
	}

	return nil
}
