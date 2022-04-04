package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	crs "github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

// MappedSession consists of map for sessions of each batch
type MappedSession struct {
	general.TenantBase
	BatchID         uuid.UUID          `json:"batchID" gorm:"type:varchar(36)"`
	CourseSessionID uuid.UUID          `json:"courseSessionID" gorm:"type:varchar(36)"`
	StartDate       *string            `json:"startDate" gorm:"type:datetime"`
	Order           uint               `json:"order" gorm:"type:int"`
	IsCompleted     *bool              `json:"isCompleted" gorm:"type:tinyint(1)"`
	Session         *crs.CourseSession `json:"session" gorm:"foreignkey:SessionID;association_autocreate:false;association_autoupdate:false"`
}

func (batchSession *MappedSession) ValidateSession() error {

	if !util.IsUUIDValid(batchSession.BatchID) {
		return errors.NewValidationError("Batch ID must be specified")
	}
	if !util.IsUUIDValid(batchSession.CourseSessionID) {
		return errors.NewValidationError("Session ID must be specified")
	}
	return nil
}

// TableName indicates name for table in DB.
func (*MappedSession) TableName() string {
	return "old_batch_sessions"
}

// MappedSessionDTO consists of map for sessions of each batch
type MappedSessionDTO struct {
	ID                         uuid.UUID                     `json:"id"`
	DeletedAt                  *time.Time                    `json:"-"`
	BatchID                    uuid.UUID                     `json:"batchID"`
	CourseSessionID            uuid.UUID                     `json:"courseSessionID"`
	StartDate                  *string                       `json:"startDate"`
	Order                      uint                          `json:"order"`
	IsCompleted                *bool                         `json:"isCompleted"`
	Session                    *crs.CourseSessionDTO         `json:"session" gorm:"foreignkey:CourseSessionID"`
	SessionAssignment          []*TopicAssignmentDTO         `json:"sessionAssignment"`
	TalentBatchSessionFeedback []*TalentBatchSessionFeedback `json:"talentBatchSessionFeedback"`
	BatchStartTime             *string                       `json:"batchStartTime"`
	Faculty                    *list.Faculty                 `json:"faculty"`

	// Flag
	IsFeedbackGiven   bool `json:"isFeedbackGiven"`
	IsAttendanceGiven bool `json:"isAttendanceGiven"`
	// TakenBy   *uuid.UUID `json:"takenBy" gorm:"type:varchar(36)"`
	// FacultyCredential list.FacultyCredentialDTO `json:"facultyCredential"`
}

// TableName indicates name for table in DB.
func (*MappedSessionDTO) TableName() string {
	return "old_batch_sessions"
}
