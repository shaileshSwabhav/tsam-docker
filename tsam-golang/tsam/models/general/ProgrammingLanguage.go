package general

import (
	"github.com/techlabs/swabhav/tsam/util"

	"github.com/techlabs/swabhav/tsam/errors"
)

// ProgrammingLanguage contains add update fields of programming language.
type ProgrammingLanguage struct {
	TenantBase
	Name          string  `json:"name" gorm:"type:varchar(100)"`
	Rating        uint8   `json:"rating" gorm:"type:int"`
	FileExtension *string `json:"fileExtension" gorm:"type:varchar(20)"`
}

// Validate returns error if language or id fields are not initialized with proper values.
func (language *ProgrammingLanguage) Validate() error {

	// Check if name is present or not.
	if util.IsEmpty(language.Name) {
		return errors.NewValidationError("Programming language name must be specified")
	}

	// Check length of name.
	if len(language.Name) > 100 {
		return errors.NewValidationError("Programming language name can have maximum 100 characters")
	}

	// Rating cannot be below 1
	if language.Rating <= 0 {
		return errors.NewValidationError("Pragramming language rating cannot be below 1")
	}

	// Rating cannot be above 5
	if language.Rating > 5 {
		return errors.NewValidationError("Pragramming language rating cannot be above 5")
	}

	// Check length of file extension.
	if language.FileExtension != nil && len(*language.FileExtension) > 20 {
		return errors.NewValidationError("File Extension can have maximum 20 characters")
	}

	return nil
}
