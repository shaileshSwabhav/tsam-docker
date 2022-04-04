package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// FeedbackQuestion contains fields of a question in feedback.
type FeedbackQuestion struct {
	TenantBase
	Type                  string                 `json:"type" gorm:"varchar(100)"`
	Question              string                 `json:"question" gorm:"varchar(250)"`
	HasOptions            *bool                  `json:"hasOptions"`
	Keyword               string                 `json:"keyword" gorm:"type:varchar(30)"`
	Order                 uint                   `json:"order" gorm:"type:tinyint(1)"`
	MaxScore              *int                   `json:"maxScore" gorm:"type:smallint(2)"`
	Options               []FeedbackOption       `json:"options" gorm:"foreignkey:QuestionID"`
	IsActive              *bool                  `json:"isActive"`
	GroupID               *uuid.UUID             `json:"-" gorm:"type:varchar(36)"`
	FeedbackQuestionGroup *FeedbackQuestionGroup `json:"feedbackQuestionGroup" gorm:"foreignkey:GroupID;association_autocreate:false;association_autoupdate:false"`
}

// Validate returns error if GeneralType is invalid.
func (feedbackQuestion *FeedbackQuestion) Validate() error {

	if util.IsEmpty(feedbackQuestion.Type) {
		return errors.NewValidationError("Type must be specified")
	}
	if util.IsEmpty(feedbackQuestion.Question) {
		return errors.NewValidationError("Question must be specified")
	}
	if feedbackQuestion.Order <= 0 {
		return errors.NewValidationError("Order cannot be less than or equal to zero")
	}
	if feedbackQuestion.IsActive == nil {
		return errors.NewValidationError("Is active must be specified")
	}

	if *feedbackQuestion.HasOptions {
		if util.IsEmpty(feedbackQuestion.Keyword) {
			return errors.NewValidationError("Keyword must be specified")
		}
		if feedbackQuestion.Options == nil || len(feedbackQuestion.Options) == 0 {
			return errors.NewValidationError("Options must be specified")
		}
	}

	return nil
}

// FeedbackQuestionDTO contains fields of question in feedback.
type FeedbackQuestionDTO struct {
	ID                    uuid.UUID                 `json:"id"`
	DeletedAt             *time.Time                `json:"-"`
	Type                  string                    `json:"type"`
	Question              string                    `json:"question"`
	HasOptions            *bool                     `json:"hasOptions"`
	Keyword               *string                   `json:"keyword"`
	Order                 uint                      `json:"order"`
	MaxScore              *int                      `json:"maxScore"`
	Options               []FeedbackOptionDTO       `json:"options" gorm:"foreignkey:QuestionID;"`
	IsActive              *bool                     `json:"isActive"`
	GroupID               *uuid.UUID                `json:"-"`
	FeedbackQuestionGroup *FeedbackQuestionGroupDTO `json:"feedbackQuestionGroup" gorm:"foreignkey:GroupID;"`
}

// TableName defines table name of the struct.
func (*FeedbackQuestionDTO) TableName() string {
	return "feedback_questions"
}
