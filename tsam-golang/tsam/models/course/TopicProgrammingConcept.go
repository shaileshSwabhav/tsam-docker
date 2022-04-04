package course

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/util"
)

// TopicProgrammingConcept is a map of programming_concept and module_topics.
type TopicProgrammingConcept struct {
	general.TenantBase
	TopicID              uuid.UUID `json:"topicID" gorm:"type:varchar(36)"`
	ProgrammingConceptID uuid.UUID `json:"programmingConceptID" gorm:"type:varchar(36)"`
	// CourseID             uuid.UUID `json:"courseID" gorm:"type:varchar(36)"`
	// IsActive             *bool     `json:"isActive"`
	// CourseSessionID         uuid.UUID `json:"courseSessionID" gorm:"type:varchar(36)"`
	// Order                uint      `json:"order" gorm:"type:int"`
}

// Validate will validate all the fields of topic_programming_concepts table.
func (concept *TopicProgrammingConcept) Validate() error {

	// if !util.IsUUIDValid(concept.CourseID) {
	// 	return errors.NewValidationError("Course ID must be specified.")
	// }

	if !util.IsUUIDValid(concept.ProgrammingConceptID) {
		return errors.NewValidationError("Programming concept ID must be specified.")
	}

	if !util.IsUUIDValid(concept.TopicID) {
		return errors.NewValidationError("Sub topic ID must be specified.")
	}

	return nil
}

// TopicProgrammingConceptDTO defines fields that are used for get operation.
type TopicProgrammingConceptDTO struct {
	ID                   uuid.UUID                          `json:"id"`
	DeletedAt            *time.Time                         `json:"-"`
	TopicID              uuid.UUID                          `json:"-"`
	ProgrammingConceptID uuid.UUID                          `json:"programmingConceptID"`
	Topic                *ModuleTopicDTO                    `json:"topic"`
	ProgrammingConcept   *programming.ProgrammingConceptDTO `json:"programmingConcept" gorm:"foreignkey:ProgrammingConceptID"`
	// IsActive             *bool                           `json:"isActive"`
	// Order                uint                            `json:"order"`
}

// TableName defines table name of the struct.
func (*TopicProgrammingConceptDTO) TableName() string {
	return "topic_programming_concepts"
}
