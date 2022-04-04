package programming

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// ProgrammingQuestionSolution contains add update fields required for programming question solution.
type ProgrammingQuestionSolution struct {
	general.TenantBase
	Solution              string    `json:"solution" gorm:"type:varchar(2000)"`
	ProgrammingQuestionID uuid.UUID `json:"programmingQuestionID" gorm:"type:varchar(36)"`
	ProgrammingLanguageID uuid.UUID `json:"programmingLanguageID" gorm:"type:varchar(36)"`
}

// Validate validates compulsary fields of ProgrammingQuestionSolution.
func (solution *ProgrammingQuestionSolution) Validate() error {

	// Check if solution is blank or not.
	if util.IsEmpty(solution.Solution) {
		return errors.NewValidationError("Solution must be specified")
	}

	// Solution maximum characters.
	if len(solution.Solution) > 100 {
		return errors.NewValidationError("Solution can have maximum 2000 characters")
	}

	// Programming Language ID.
	if !util.IsUUIDValid(solution.ProgrammingLanguageID) {
		return errors.NewValidationError("Programming Language ID must ne a proper uuid")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// ProgrammingQuestionSolutionDTO contains all fields required for programming question solution.
type ProgrammingQuestionSolutionDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	ProgrammingLanguageID uuid.UUID                   `json:"-"`
	ProgrammingLanguage   general.ProgrammingLanguage `json:"programmingLanguage" gorm:"foreignkey:ProgrammingLanguageID"`

	Solution              string    `json:"solution"`
	ProgrammingQuestionID uuid.UUID `json:"programmingQuestionID"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionSolutionDTO) TableName() string {
	return "programming_question_solutions"
}
