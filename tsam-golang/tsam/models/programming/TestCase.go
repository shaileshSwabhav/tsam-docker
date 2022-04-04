package programming

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// ProgrammingQuestionTestCase contains add update fields required for test case.
type ProgrammingQuestionTestCase struct {
	general.TenantBase
	ProgrammingQuestionID uuid.UUID `json:"programmingQuestionID" gorm:"type:varchar(36)"`
	Input                 string    `json:"input" gorm:"type:varchar(500)"`
	Output                string    `json:"output" gorm:"type:varchar(500)"`
	Explanation           *string   `json:"explanation" gorm:"type:varchar(500)"`
	IsActive              bool      `json:"isActive"`
	IsHidden              bool      `json:"isHidden"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionTestCase) TableName() string {
	return "programming_question_test_cases"
}

// Validate validates compulsary fields of ProgrammingQuestion.
func (testCase *ProgrammingQuestionTestCase) Validate() error {

	// Check if input is blank or not.
	if util.IsEmpty(testCase.Input) {
		return errors.NewValidationError("Input must be specified")
	}

	// Input maximum characters.
	if len(testCase.Input) > 500 {
		return errors.NewValidationError("Input can have maximum 500 characters")
	}

	// Check if output is blank or not.
	if util.IsEmpty(testCase.Output) {
		return errors.NewValidationError("Output must be specified")
	}

	// Output maximum characters.
	if len(testCase.Output) > 500 {
		return errors.NewValidationError("Output can have maximum 500 characters")
	}

	// Explanation maximum characters.
	if testCase.Explanation != nil && len(*testCase.Explanation) > 500 {
		return errors.NewValidationError("Explanation can have maximum 500 characters")
	}

	return nil
}
