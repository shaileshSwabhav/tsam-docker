package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Examination fields for the exams the talent or enquiry has given.
type Examination struct {
	TenantBase
	Name       string  `json:"name" gorm:"type:varchar(50)"`
	TotalMarks float64 `json:"totalMarks" gorm:"type:decmial(10, 2)"`
	// TotalMarks uint16 `json:"totalMarks" gorm:"type:smallint"` -> ielts score has decimal
}

// Validate will validate all the fields of examination.
func (examination *Examination) Validate() error {

	// Check if name is present or not.
	if util.IsEmpty(examination.Name) {
		return errors.NewValidationError("Examination Name must be specified")
	}

	// Check if name is alphabets or not.
	if !util.ValidateStringWithSpace(examination.Name) {
		return errors.NewValidationError("Examination Name must contain alphabets and space only")
	}

	// Name maximum characters.
	if len(examination.Name) > 50 {
		return errors.NewValidationError("Examination Name can have maximum 50 characters")
	}

	// Check if total marks is 0 or not.
	if examination.TotalMarks == 0 {
		return errors.NewValidationError("Examination total marks must be greater than 0")
	}

	return nil
}
