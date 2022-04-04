package course

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/util"
)

// TopicProgrammingQuestion is a map of programming_questions and module_topics.
type TopicProgrammingQuestion struct {
	general.TenantBase
	ProgrammingQuestionID uuid.UUID `json:"programmingQuestionID" gorm:"type:varchar(36)"`
	TopicID               uuid.UUID `json:"topicID" gorm:"type:varchar(36)"`
	IsActive              *bool     `json:"isActive" gorm:"default:true"`
	// Order                 uint       `json:"order" gorm:"type:int"`
}

// Validate will validate all the fields of topic_programming_questions table.
func (assignment *TopicProgrammingQuestion) Validate() error {

	if !util.IsUUIDValid(assignment.ProgrammingQuestionID) {
		return errors.NewValidationError("Programming question ID must be specified.")
	}

	if !util.IsUUIDValid(assignment.TopicID) {
		return errors.NewValidationError("Topic ID must be specified.")
	}

	return nil
}

// TopicProgrammingQuestionDTO defines fields that are used for get operation.
type TopicProgrammingQuestionDTO struct {
	ID                    uuid.UUID                           `json:"id"`
	DeletedAt             *time.Time                          `json:"-"`
	TopicID               uuid.UUID                           `json:"-"`
	Topic                 *ModuleTopicDTO                     `json:"topic"`
	ProgrammingQuestionID uuid.UUID                           `json:"-"`
	ProgrammingQuestion   *programming.ProgrammingQuestionDTO `json:"programmingQuestion"`
	ProgrammingConceptID  uuid.UUID                           `json:"-"`
	ProgrammingConcept    *programming.ProgrammingConceptDTO  `json:"programmingConcept"`
	IsActive              *bool                               `json:"isActive"`
	// Order                 uint                                `json:"order"`
}

// TableName defines table name of the struct.
func (*TopicProgrammingQuestionDTO) TableName() string {
	return "topic_programming_questions"
}
