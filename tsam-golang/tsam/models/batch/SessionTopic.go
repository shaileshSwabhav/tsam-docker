package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

// SessionTopic consist of details related to batch_session_topics.
type SessionTopic struct {
	general.TenantBase
	BatchID       uuid.UUID `json:"batchID" gorm:"type:varchar(36)"`
	ModuleID      uuid.UUID `json:"moduleID" gorm:"type:varchar(36)"`
	TopicID       uuid.UUID `json:"topicID" gorm:"type:varchar(36)"`
	SubTopicID    uuid.UUID `json:"subTopicID" gorm:"type:varchar(36)"`
	Order         uint      `json:"order" gorm:"type:INT"`
	InitialDate   *string   `json:"initialDate" gorm:"type:date"`
	CompletedDate *string   `json:"completedDate" gorm:"type:date"`
	IsCompleted   *bool     `json:"isCompleted" gorm:"type:TINYINT"`
	SessionID     uuid.UUID `json:"batchSessionID" gorm:"type:varchar(36);column:batch_session_id"`
	TotalTime     uint      `json:"totalTime" gorm:"type:INT"`
}

// TableName overrides name of the table
func (*SessionTopic) TableName() string {
	return "batch_session_topics"
}

// Validate will verify fields of batch-session-plan.
func (plan *SessionTopic) Validate() error {

	if !util.IsUUIDValid(plan.BatchID) {
		return errors.NewValidationError("Invalid batch ID")
	}

	if !util.IsUUIDValid(plan.TopicID) {
		return errors.NewValidationError("Invalid batch topic ID")
	}

	if plan.IsCompleted == nil {
		return errors.NewValidationError("is completed must be specified")
	}
	return nil
}

// SessionTopicDTO will store topics for batch_session.
type SessionTopicDTO struct {
	ID             uuid.UUID             `json:"id"`
	DeletedAt      *time.Time            `json:"-"`
	BatchID        uuid.UUID             `json:"-"`
	ModuleID       uuid.UUID             `json:"-"`
	TopicID        uuid.UUID             `json:"-"`
	SubTopicID     uuid.UUID             `json:"-"`
	BatchSessionID uuid.UUID             `json:"-"`
	Batch          *list.Batch           `json:"batch"`
	Module         course.ModuleDTO      `json:"module"`
	Topic          course.ModuleTopicDTO `json:"topic"`
	SubTopic       course.ModuleTopicDTO `json:"subTopic"`
	Order          uint                  `json:"order"`
	InitialDate    string                `json:"initialDate"`
	CompletedDate  string                `json:"completedDate"`
	IsCompleted    *bool                 `json:"isCompleted"`
	TotalTime      uint                  `json:"totalTime"`
	BatchSession   []*Session            `json:"batchSession"`
}

// TableName overrides name of the table
func (*SessionTopicDTO) TableName() string {
	return "batch_session_topics"
}

type UpdateBatchSessionPlan struct {
	BatchID                  uuid.UUID
	TenantID                 uuid.UUID
	UpdatedBy                uuid.UUID
	FacultyID                uuid.UUID
	InitialDate              string
	LastCompletedSessionDate string
	ExistingSessionPlan      []SessionTopicDTO
	Batch                    Batch
}

type ModuleTopics struct {
	general.BaseDTO
	TopicName   string         `json:"topicName" gorm:"type:varchar(50)"`
	SubTopics   []ModuleTopics `json:"subTopics" gorm:"foreignkey:TopicID"`
	TopicID     *uuid.UUID     `json:"topicID" gorm:"type:varchar(36)"`
	ModuleID    uuid.UUID      `json:"moduleID" gorm:"type:varchar(36)"`
	Order       uint           `json:"order"`
	IsCompleted *bool          `json:"isCompleted" gorm:"type:TINYINT"`
	TotalTime   uint           `json:"totalTime" gorm:"type:INT"`

	BatchSessionTopicID          *uuid.UUID `json:"batchSessionTopicID"`
	BatchSessionID               *uuid.UUID `json:"batchSessionID" gorm:"type:varchar(36);column:batch_session_id"`
	BatchSessionTopicOrder       *uint      `json:"batchSessionTopicOrder"`
	BatchSessionTopicIsCompleted *bool      `json:"batchSessionTopicIsCompleted" gorm:"type:TINYINT"`
	BatchSessionTopicTotalTime   *uint      `json:"batchSessionTopicTotalTime" gorm:"type:INT"`
	BatchSessionCompletionDate   *string    `json:"batchSessionCompletionDate" gorm:"type:date"`
	BatchSessionInitialDate      *string    `json:"batchSessionInitialDate" gorm:"type:date"`
}
