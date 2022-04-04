package programming

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingQuestionSolutionIsViewed contains add update fields required for programming question solution is viewed.
type ProgrammingQuestionSolutionIsViewed struct {
	general.TenantBase

	ProgrammingQuestionID uuid.UUID `json:"programmingQuestionID" gorm:"type:varchar(36)"`
	TalentID              uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionSolutionIsViewed) TableName() string {
	return "programming_solution_is_viewed"
}

// Validate validates compulsary fields of ProgrammingQuestionSolutionIsViewed.
func (solutionIsViewed *ProgrammingQuestionSolutionIsViewed) Validate() error {

	// Programming Question ID.
	if !util.IsUUIDValid(solutionIsViewed.ProgrammingQuestionID) {
		return errors.NewValidationError("Programming Question ID must be a proper uuid")
	}

	// Programming Talent ID.
	if !util.IsUUIDValid(solutionIsViewed.TalentID) {
		return errors.NewValidationError("Talent ID must be a proper uuid")
	}

	return nil
}
