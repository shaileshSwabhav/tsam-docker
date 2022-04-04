package test

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	common "github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Question have Subject, Difficulty, Options
type Question struct {
	common.TenantBase
	Question   *string   `json:"question"`
	Subject    *string   `json:"subject" gorm:"type:varchar(80)"`
	Difficulty *string   `json:"difficulty" gorm:"type:varchar(80)"`
	Options    []*Option `json:"options" gorm:"foreignkey:QuestionID;"`
}

// CreateNewQuestion Create New Instance Of Question.
func CreateNewQuestion(id *uuid.UUID) *Question {
	return &Question{
		TenantBase: common.TenantBase{
			ID: *id,
		},
	}
}

// MakeQuestionValidOrError Add ID To Option Or Return Error On Invalid Data
func (question *Question) MakeQuestionValidOrError() error {
	if util.IsNil(question.Question) {
		return errors.NewValidationError("Question must be specified")
	}
	if util.IsNil(question.Subject) {
		return errors.NewValidationError("Subject must be specified")
	}
	if util.IsNil(question.Difficulty) {
		return errors.NewValidationError("Difficulty must be specified")
	}

	for _, option := range question.Options {
		if err := option.ValidateOption(); err != nil {
			return err
		}
		if !util.IsUUIDValid(option.ID) {
			option.ID = util.GenerateUUID()
		}
	}
	return nil
}
