package batch

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// FacultyTalentFeedback contains the fields inside batch_feedback table.
type FacultyTalentFeedback struct {
	general.TenantBase
	BatchID    uuid.UUID                 `json:"batchID" gorm:"type:varchar(36)"`
	TalentID   uuid.UUID                 `json:"talentID" gorm:"type:varchar(36)"`
	FacultyID  uuid.UUID                 `json:"facultyID" gorm:"type:varchar(36)"`
	QuestionID uuid.UUID                 `json:"questionID" gorm:"type:varchar(36)"`
	Question   *general.FeedbackQuestion `json:"question" gorm:"foreignkey:QuestionID"`
	OptionID   *uuid.UUID                `json:"optionID" gorm:"type:varchar(36)"` // check if it will contain optionID
	Answer     string                    `json:"answer" gorm:"type:varchar(2000)"`
	Option     *general.FeedbackOption   `json:"option" gorm:"foreignkey:OptionID"`
	// Talent     *talent.Talent            `json:"talent" gorm:"foreignkey:TalentID"`
}

// TableName to batch_feedback (plural : feedback)
func (*FacultyTalentFeedback) TableName() string {
	return "faculty_talent_batch_feedback"
}

// Validate TalentFeedback fields
func (talentFeedback *FacultyTalentFeedback) Validate() error {
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
