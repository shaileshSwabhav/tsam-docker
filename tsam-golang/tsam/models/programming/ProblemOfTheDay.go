package programming

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// ProblemOfTheDay contains add update fields required for problem of the day.
type ProblemOfTheDay struct {
	general.TenantBase
	Date                  string    `json:"date" gorm:"type:date"`
	ProgrammingQuestionID uuid.UUID `json:"programmingQuestionID" gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*ProblemOfTheDay) TableName() string {
	return "problem_of_the_day"
}

// Validate validates compulsary fields of ProblemOfTheDay.
func (problemOfTheDay *ProblemOfTheDay) Validate() error {

	// Check if date is blank or not.
	if util.IsEmpty(problemOfTheDay.Date) {
		return errors.NewValidationError("Date must be specified")
	}

	// Programming Question ID.
	if !util.IsUUIDValid(problemOfTheDay.ProgrammingQuestionID) {
		return errors.NewValidationError("Programming Question ID must ne a proper uuid")
	}

	return nil
}
