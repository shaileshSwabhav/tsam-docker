package faculty

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// FacultyAssessment contains the fields inside faculty_assessment table.
type FacultyAssessment struct {
	general.TenantBase
	CredentialID uuid.UUID  `json:"credentialID" gorm:"type:varchar(36)"`
	FacultyID    uuid.UUID  `json:"facultyID" gorm:"type:varchar(36)"`
	QuestionID   uuid.UUID  `json:"questionID" gorm:"type:varchar(36)"`
	GroupID      uuid.UUID  `json:"groupID" gorm:"type:varchar(36)"`
	OptionID     *uuid.UUID `json:"optionID" gorm:"type:varchar(36)"`
	Answer       string     `json:"answer" gorm:"type:varchar(2000)"`
	// AssessmentType string     `json:"assessmentType" gorm:"type:varchar(30)"`

	// Credential            general.Credential            `json:"credential" gorm:"foreignkey:CredentialID"`
	// Faculty               Faculty                       `json:"faculty" gorm:"foreignkey:FacultyID"`
	// Question              general.FeedbackQuestion      `json:"question" gorm:"foreignkey:QuestionID"`
	// Option                *general.FeedbackOption       `json:"option" gorm:"foreignkey:OptionID"`
	// FeedbackQuestionGroup general.FeedbackQuestionGroup `json:"feedbackQuestionGroup" gorm:"foreignkey:GroupID;association_autocreate:false;association_autoupdate:false"`
}

func (assessment *FacultyAssessment) ValidateFacultyAssessment() error {

	if assessment.OptionID == nil {
		if util.IsEmpty(assessment.Answer) {
			return errors.NewValidationError("Answer must be specified")
		}
	}
	if util.IsEmpty(assessment.Answer) {
		if assessment.OptionID == nil {
			return errors.NewValidationError("Answer must be specified")
		}
	}
	return nil
}

// FacultyAssessmentDTO contains the fields inside faculty_assessment table.
type FacultyAssessmentDTO struct {
	ID                    uuid.UUID                     `json:"id"`
	DeletedAt             *time.Time                    `json:"-"`
	CredentialID          uuid.UUID                     `json:"-"`
	FacultyID             uuid.UUID                     `json:"-"`
	QuestionID            uuid.UUID                     `json:"-"`
	OptionID              *uuid.UUID                    `json:"-"`
	GroupID               uuid.UUID                     `json:"-"`
	Answer                string                        `json:"answer"`
	Credential            general.Credential            `json:"credential" gorm:"foreignkey:CredentialID"`
	Faculty               Faculty                       `json:"faculty" gorm:"foreignkey:FacultyID"`
	Question              general.FeedbackQuestion      `json:"question" gorm:"foreignkey:QuestionID"`
	Option                *general.FeedbackOption       `json:"option" gorm:"foreignkey:OptionID"`
	FeedbackQuestionGroup general.FeedbackQuestionGroup `json:"feedbackQuestionGroup" gorm:"foreignkey:GroupID"`

	// AssessmentType string     `json:"assessmentType" gorm:"type:varchar(30)"`
}

// TableName overrides name of the table.
func (*FacultyAssessmentDTO) TableName() string {
	return "faculty_assessments"
}
