package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// FeedbackQuestionGroup will consist of group name and its description.
type FeedbackQuestionGroup struct {
	TenantBase
	GroupName         string             `json:"groupName" gorm:"type:varchar(50)"`
	GroupDescription  string             `json:"groupDescription" gorm:"type:varchar(200)"`
	Order             uint               `json:"order" gorm:"type:SMALLINT(2)"`
	Type              string             `json:"type" gorm:"type:varchar(100)"`
	MaxScore          uint               `json:"maxScore" gorm:"-"`
	MinScore          uint               `json:"minScore" gorm:"-"`
	FeedbackQuestions []FeedbackQuestion `json:"feedbackQuestions" gorm:"-"`
}

// Validate will check if group name and description are valid.
func (group *FeedbackQuestionGroup) Validate() error {

	// Group name must be specified.
	if util.IsEmpty(group.GroupName) {
		return errors.NewValidationError("Group name must be specified")
	}

	// Group name maximum characters.
	if len(group.GroupName) > 50 {
		return errors.NewValidationError("Group name can have maximum 50 characters")
	}

	// Group description must be specified.
	if util.IsEmpty(group.GroupDescription) {
		return errors.NewValidationError("Group description must be specified")
	}

	// Group description maximum characters.
	if len(group.GroupDescription) > 200 {
		return errors.NewValidationError("Group description can have maximum 200 characters")
	}

	// Group order must be greater than 0.
	if group.Order <= 0 {
		return errors.NewValidationError("Order must be above 0")
	}
	return nil
}

// FeedbackQuestionGroupDTO will consist of group name and its description.
type FeedbackQuestionGroupDTO struct {
	ID                uuid.UUID             `json:"id"`
	DeletedAt         *time.Time            `json:"-"`
	GroupName         string                `json:"groupName"`
	GroupDescription  string                `json:"groupDescription"`
	Order             uint                  `json:"order"`
	Type              string                `json:"type"`
	MaxScore          uint                  `json:"maxScore"`
	MinScore          uint                  `json:"minScore"`
	FeedbackQuestions []FeedbackQuestionDTO `json:"feedbackQuestions" gorm:"foreignkey:GroupID"`
}

// TableName defines table name of the struct.
func (FeedbackQuestionGroupDTO) TableName() string {
	return "feedback_question_groups"
}
