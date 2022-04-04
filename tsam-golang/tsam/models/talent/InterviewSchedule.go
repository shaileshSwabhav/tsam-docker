package talent

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// InterviewSchedule contains the interview schedule info of talent.
type InterviewSchedule struct {
	general.TenantBase
	TalentID      uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	ScheduledDate string    `json:"scheduledDate" gorm:"type:date"`
	Status        string    `json:"status" gorm:"type:varchar(100)"`
}

// TableName will create the table for InterviewSchedule model with name talent_interview_schedules.
func (*InterviewSchedule) TableName() string {
	return "talent_interview_schedules"
}

// Validate Validates fields of talent interview schedule
func (interviewSchedule *InterviewSchedule) Validate() error {
	// Check if Interview Schedule Date is empty or not.
	if util.IsEmpty(interviewSchedule.ScheduledDate) {
		return errors.NewValidationError("Interview Schedule Date must be specified")
	}

	// Check if Interview Schedule status is empty or not.
	if util.IsEmpty(interviewSchedule.Status) {
		return errors.NewValidationError("Interview Schedule Status must be specified")
	}

	return nil
}
