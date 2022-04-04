package talent

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Interview contains details of every interview round of a talent.
type Interview struct {
	general.TenantBase
	ScheduleID uuid.UUID            `json:"scheduleID" gorm:"type:varchar(36)"`
	RoundID    uuid.UUID            `json:"roundID" gorm:"type:varchar(36)"`
	TalentID   uuid.UUID            `json:"talentID" gorm:"type:varchar(36)"`
	Rating     uint8                `json:"rating" gorm:"type:smallint(2)"`
	Status     string               `json:"status" gorm:"type:varchar(100)"`
	Comment    *string              `json:"comment" gorm:"type:varchar(1000)"`
	TakenBy    []general.Credential `json:"takenBy" gorm:"many2many:talent_interviews_credentials;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
}

// TableName creates the table for model Interview with name talent_interviews.
func (*Interview) TableName() string {
	return "talent_interviews"
}

// Validate Validates fields of talent interview.
func (interview *Interview) Validate() error {
	// Check if rating is above 0.
	if interview.Rating < 1 {
		return errors.NewValidationError("Interview rating must be above 0")
	}

	// Check if rating is not more than 10.
	if interview.Rating > 10 {
		return errors.NewValidationError("Interview rating cannot be above 10")
	}

	// Check if round id is mentioned or not.
	if !util.IsUUIDValid(interview.RoundID) {
		return errors.NewValidationError("Interview round must be specified")
	}

	// Check if interview talent id is mentioned or not.
	if !util.IsUUIDValid(interview.TalentID) {
		return errors.NewValidationError("Interview talent id must be specified")
	}

	// Check if taken by is mentioned or not.
	if interview.TakenBy != nil && len(interview.TakenBy) <= 0 {
		return errors.NewValidationError("Interview Taken by must be specified")
	}

	// Check if comment has more than max characters or not.
	if interview.Comment != nil && len(*interview.Comment) > 1000 {
		return errors.NewValidationError("Comment can have maximum 1000 characters")
	}

	// Check if status is mentioned or not.
	if util.IsEmpty(interview.Status) {
		return errors.NewValidationError("Interview status must be specified")
	}

	return nil
}

// ===========Defining many to many structs===========

// InterviewTakenBy is the map of interview and taken by.
type InterviewTakenBy struct {
	InterviewID uuid.UUID `gorm:"type:varchar(36)"`
	TakenByID   uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*InterviewTakenBy) TableName() string {
	return "talent_interviews_credentials"
}
