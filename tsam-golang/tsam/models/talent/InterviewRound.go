package talent

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// InterviewRound contains the details about the round for an interview.
type InterviewRound struct {
	general.TenantBase
	Name string `json:"name" gorm:"type:varchar(100)"`
}

// TableName will create the table for InterviewRound model with name talent_interview_rounds.
func (*InterviewRound) TableName() string {
	return "talent_interview_rounds"
}

// Validate Validates fields of talent interview round.
func (interviewRound *InterviewRound) Validate() error {
	// Check if name is empty or not.
	if util.IsEmpty(interviewRound.Name) {
		return errors.NewValidationError("Interview name must be specified")
	}
	return nil
}
