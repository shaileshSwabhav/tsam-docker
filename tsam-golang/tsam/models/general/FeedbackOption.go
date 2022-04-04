package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// FeedbackOption contains the options which is to be given in FeedbackQuestion.
type FeedbackOption struct {
	TenantBase
	QuestionID uuid.UUID `json:"questionID" gorm:"type:varchar(36)"`
	Order      uint      `json:"order" gorm:"smallint(2)"`
	Key        int       `json:"key" gorm:"type:smallint(2)"`
	Value      string    `json:"value" gorm:"type:varchar(250)"`
}

// Validate will check the fields of FeedbackOption
func (feedbackOption *FeedbackOption) Validate() error {
	if util.IsEmpty(feedbackOption.Value) {
		return errors.NewValidationError("Value must be specified")
	}
	if feedbackOption.Key <= 0 {
		return errors.NewValidationError("Key must be specified")
	}
	return nil
}

// FeedbackOptionDTO contains the options given to FeedbackQuestion.
type FeedbackOptionDTO struct {
	ID         uuid.UUID  `json:"id"`
	DeletedAt  *time.Time `json:"-"`
	QuestionID uuid.UUID  `json:"questionID"`
	Order      uint       `json:"order"`
	Key        int        `json:"key"`
	Value      string     `json:"value"`
}

// TableName defines table name of the struct.
func (*FeedbackOptionDTO) TableName() string {
	return "feedback_options"
}
