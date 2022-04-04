package batch

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// FacultyTalentBatchSessionFeedback contains the fields inside faculty_talent_batch_session_feedback table.
type FacultyTalentBatchSessionFeedback struct {
	general.TenantBase
	BatchID        uuid.UUID                 `json:"batchID" gorm:"type:varchar(36)"`
	TalentID       uuid.UUID                 `json:"talentID" gorm:"type:varchar(36)"`
	QuestionID     uuid.UUID                 `json:"questionID" gorm:"type:varchar(36)"`
	FacultyID      uuid.UUID                 `json:"facultyID" gorm:"type:varchar(36)"`
	Question       *general.FeedbackQuestion `json:"question" gorm:"foreignkey:QuestionID"`
	OptionID       *uuid.UUID                `json:"optionID" gorm:"type:varchar(36)"`
	Answer         string                    `json:"answer" gorm:"type:varchar(2000)"`
	Option         *general.FeedbackOption   `json:"option" gorm:"foreignkey:OptionID"`
	BatchSessionID uuid.UUID                 `json:"batchSessionID" gorm:"type:varchar(36)"`
	// BatchTopicID uuid.UUID                 `json:"batchTopicID" gorm:"type:varchar(36)"`
	// Talent     *talent.Talent            `json:"talent" gorm:"foreignkey:talentID"`
}

// TableName to faculty_talent_batch_session_feedback (plural : feedback)
func (*FacultyTalentBatchSessionFeedback) TableName() string {
	return "faculty_talent_batch_session_feedback"
}

// Validate will check the fields of TalentSessionFeedback
func (sessionFeedback *FacultyTalentBatchSessionFeedback) Validate() error {
	if sessionFeedback.OptionID == nil {
		if util.IsEmpty(sessionFeedback.Answer) {
			return errors.NewValidationError("Answer must be specified")
		}
	}
	if util.IsEmpty(sessionFeedback.Answer) {
		if sessionFeedback.OptionID == nil {
			return errors.NewValidationError("Answer must be specified")
		}
	}
	// if !util.IsUUIDValid(sessionFeedback.BatchID) {
	// 	return errors.NewValidationError("Invalid batch specified")
	// }
	// if !util.IsUUIDValid(sessionFeedback.SessionID) {
	// 	return errors.NewValidationError("Invalid session specified")
	// }
	// if !util.IsUUIDValid(sessionFeedback.QuestionID) {
	// 	return errors.NewValidationError("Invalid question specified")
	// }
	return nil
}
