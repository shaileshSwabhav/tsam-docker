package batch

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// TalentFeedback contains the fields inside talent_batch_feedback table.
type TalentFeedback struct {
	general.TenantBase
	BatchID    uuid.UUID                 `json:"batchID" gorm:"type:varchar(36)"`
	TalentID   uuid.UUID                 `json:"talentID" gorm:"type:varchar(36)"`
	FacultyID  uuid.UUID                 `json:"facultyID" gorm:"type:varchar(36)"`
	QuestionID uuid.UUID                 `json:"questionID" gorm:"type:varchar(36)"`
	OptionID   *uuid.UUID                `json:"optionID" gorm:"type:varchar(36)"` // check if it will contain optionID
	Answer     string                    `json:"answer" gorm:"type:varchar(2000)"`
	Question   *general.FeedbackQuestion `json:"question" gorm:"foreignkey:QuestionID"`
	Option     *general.FeedbackOption   `json:"option" gorm:"foreignkey:OptionID"`
}

// TableName to talent_batch_feedback (plural : feedback)
func (*TalentFeedback) TableName() string {
	return "talent_batch_feedback"
}

// Validate TalentFeedback fields
func (talentFeedback *TalentFeedback) Validate() error {
	if talentFeedback.OptionID == nil {
		if util.IsEmpty(talentFeedback.Answer) {
			return errors.NewValidationError("Answer must be specified")
		}
	}
	if util.IsEmpty(talentFeedback.Answer) {
		if talentFeedback.OptionID == nil {
			return errors.NewValidationError("Answer must be specified")
		}
	}
	// if talentFeedback.Talent == nil {
	// 	return errors.NewValidationError("Talent must be specified")
	// }
	return nil
}
