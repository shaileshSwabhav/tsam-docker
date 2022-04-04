package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/web"
)

// Session will consists of records of topics that are completed from session_plan.
type Session struct {
	general.TenantBase
	BatchID           uuid.UUID      `json:"batchID" gorm:"type:varchar(36)"`
	FacultyID         uuid.UUID      `json:"facultyID" gorm:"type:varchar(36)"`
	Date              string         `json:"date" gorm:"date"`
	IsCompleted       *bool          `json:"isCompleted" gorm:"type:TINYINT"`
	IsSessionTaken    *bool          `json:"isSessionTaken" gorm:"type:TINYINT"`
	BatchSessionTopic []SessionTopic `json:"batchSessionTopic"`
}

// TableName overrides name of the table
func (*Session) TableName() string {
	return "batch_sessions"
}

// Validate will verify fields of batch-session-plan.
func (session *Session) Validate() error {

	// if util.IsEmpty(session.Date) {
	// 	return errors.NewValidationError("Completed date must be specified")
	// }

	if session.IsCompleted == nil {
		return errors.NewValidationError("is completed must be specified")
	}

	return nil
}

// SessionDTO will store batch_sessions.
type SessionDTO struct {
	general.BaseDTO
	BatchID                     uuid.UUID                    `json:"-"`
	Batch                       *list.Batch                  `json:"batch"`
	FacultyID                   uuid.UUID                    `json:"-"`
	Faculty                     *list.Batch                  `json:"faculty"`
	Date                        string                       `json:"date"`
	IsCompleted                 *bool                        `json:"isCompleted"`
	IsSessionTaken              *bool                        `json:"isSessionTaken"`
	BatchSessionPrerequisiteDTO *BatchSessionPrerequisiteDTO `json:"batchSessionPrerequisite" gorm:"foreignkey:BatchSessionID"`
	PendingModule               []course.ModuleDTO           `json:"pendingModule"`
	Module                      []course.ModuleDTO           `json:"module"`
	BatchSessionTopic           []*SessionTopicDTO           `json:"batchSessionTopic" gorm:"foreignkey:BatchSessionID"`
	IsAttendanceMarked          bool                         `json:"isAttendanceMarked"`
	IsFeedbackGiven             bool                         `json:"isFeedbackGiven"`
}

// TableName overrides name of the table
func (*SessionDTO) TableName() string {
	return "batch_sessions"
}

type SessionParameter struct {
	BatchSessions       *[]Session
	SessionTopics       *[]SessionTopic
	ModuleBatchSessions *[]Session
	ModuleSessionTopics *[]SessionTopic
	Modules             *[]ModuleDTO
	BatchID             uuid.UUID
	Parser              *web.Parser
	InitialDate         *string
	FacultyID           uuid.UUID
	TenantID            uuid.UUID
	CredentialID        uuid.UUID
	CurrentModuleID     uuid.UUID
	IsCreate            bool
	IsSkip              bool
	// DeletedDates  *[]string
}

type AllocateParameter struct {
	ModuleSessionTopics    *[]SessionTopic
	SessionTopics          *[]SessionTopic
	BatchSession           *Session
	ModuleFacultyMap       map[uuid.UUID]uuid.UUID
	CurrentFacultyID       uuid.UUID
	CurrentDate            time.Time
	AllocationDurationLeft float64
}

type GetModuleParams struct {
	TenantID    uuid.UUID
	BatchID     uuid.UUID
	Modules     *[]course.ModuleDTO
	Parser      *web.Parser
	IsPending   bool
	Channel     chan error
	SessionDate string
	// Queryprocessors []repository.QueryProcessor
}

// SessionCounts will return counts for modules, topics, assignments etc.
type SessionCounts struct {
	ModuleCount     uint `json:"moduleCount"`
	TopicCount      uint `json:"topicCount"`
	AssignmentCount uint `json:"assignmentCount"`
	SessionCount    uint `json:"sessionCount"`
	TotalBatchHours uint `json:"totalBatchHours"`
	ProjectCount    uint `json:"projectCount"`
}

// ************************** BATCH SESSION LIST WITH TOPIC ANS SUB TOPIC NAME ***************************

// BatchSessionWithTopicNameDTO will get batch session with topic names.
type SessionWithTopicNameDTO struct {
	general.BaseDTO
	BatchID            uuid.UUID             `json:"batchID"`
	Date               string                `json:"date"`
	IsCompleted        *bool                 `json:"isCompleted"`
	IsSessionTaken     *bool                 `json:"isSessionTaken"`
	FacultyID          uuid.UUID             `json:"-"`
	Faculty            list.Faculty          `json:"faculty"`
	BatchSessionTopics []SessionTopicNameDTO `json:"batchSessionTopics" gorm:"foreignkey:BatchSessionID"`
	ModuleID           uuid.UUID             `json:"-"`
	ModuleTiming       ModuleTiming          `json:"moduleTiming"`
}

// TableName overrides name of the table
func (*SessionWithTopicNameDTO) TableName() string {
	return "batch_sessions"
}

// SessionTopicNameDTO will consist of name regarding a specific topic.
type SessionTopicNameDTO struct {
	general.BaseDTO
	Topic          ModuleTopicNameDTO `json:"topic" gorm:"foreignkey:TopicID"`
	TopicID        uuid.UUID          `json:"-"`
	BatchSessionID uuid.UUID          `json:"-"`
	SubTopicID     uuid.UUID          `json:"-"`
	SubTopic       ModuleTopicNameDTO `json:"subTopic" gorm:"foreignkey:SubTopicID"`
}

// TableName overrides name of the table
func (*SessionTopicNameDTO) TableName() string {
	return "batch_session_topics"
}

// ModuleTopicNameDTO will consist of details regarding a specific topic.
type ModuleTopicNameDTO struct {
	general.BaseDTO
	TopicName string `json:"topicName"`
	Order     uint   `json:"order"`
	// SubTopics                 []*ModuleTopicDTO              `json:"subTopics" gorm:"foreignkey:TopicID;association_autoupdate:false"`
	// TopicID                   *uuid.UUID                     `json:"topicID"`
	// ModuleID                  uuid.UUID                      `json:"-"`
	// Module                    *ModuleDTO                     `json:"module" gorm:"foreignkey:ModuleID"`
	// ProgrammingConceptID      uuid.UUID                      `json:"-"`
	// TopicProgrammingConcept   []*TopicProgrammingConceptDTO  `json:"topicProgrammingConcept" gorm:"foreignkey:TopicID"`
	// TopicProgrammingQuestions []*TopicProgrammingQuestionDTO `json:"topicProgrammingQuestions" gorm:"foreignkey:TopicID"`

	// BatchTopicAssignment []*list.BatchTopicAssignmentDTO `json:"batchTopicAssignment" gorm:"foreignkey:TopicID"`
	// BatchSessionTopic    *list.BatchSessionTopic         `json:"batchSessionTopic" gorm:"foreignkey:SubTopicID"`
}

// TableName overrides name of the table
func (*ModuleTopicNameDTO) TableName() string {
	return "module_topics"
}
