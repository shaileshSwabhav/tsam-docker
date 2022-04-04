package course

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

// ModuleTopic will consist of details regarding a specific topic.
type ModuleTopic struct {
	general.TenantBase
	TopicName               string                     `json:"topicName" gorm:"type:varchar(1000)"`
	TotalTime               *uint                      `json:"totalTime" gorm:"type:INT"`
	Order                   uint                       `json:"order" gorm:"type:INT"`
	SubTopics               []ModuleTopic              `json:"subTopics" gorm:"foreignkey:TopicID"`
	TopicID                 *uuid.UUID                 `json:"topicID" gorm:"type:varchar(36)"`
	ModuleID                uuid.UUID                  `json:"moduleID" gorm:"type:varchar(36)"`
	TopicProgrammingConcept []*TopicProgrammingConcept `json:"topicProgrammingConcept" gorm:"foreignkey:SubTopicID;association_autocreate:true;association_autoupdate:false"`

	// StudentOutput           *string                    `json:"studentOutput" gorm:"type:varchar(100)"`
}

// Validate will verify if all compuslory fields of topic are specifed or not.
func (topic *ModuleTopic) Validate() error {

	if util.IsEmpty(topic.TopicName) {
		return errors.NewValidationError("Topic name must be specified")
	}

	if topic.TotalTime != nil {
		if *topic.TotalTime <= 0 {
			return errors.NewValidationError("Topic time should not be zero")
		}
	}

	if topic.Order <= 0 {
		return errors.NewValidationError("Topic order should not be zero")
	}

	// if util.IsEmpty(topic.StudentOutput) {
	// 	return errors.NewValidationError("Student output must be specified")
	// }

	if topic.SubTopics != nil {
		for _, subTopic := range topic.SubTopics {
			if err := subTopic.Validate(); err != nil {
				return err
			}
		}
	}

	// for _, resource := range session.Resources {
	// 	if err := resource.ValidateResource(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// ModuleTopicDTO will consist of details regarding a specific topic.
type ModuleTopicDTO struct {
	ID                        uuid.UUID                      `json:"id"`
	DeletedAt                 *time.Time                     `json:"-"`
	TopicName                 string                         `json:"topicName"`
	TotalTime                 *uint                          `json:"totalTime"`
	Order                     uint                           `json:"order"`
	SubTopics                 []*ModuleTopicDTO              `json:"subTopics" gorm:"foreignkey:TopicID;association_autoupdate:false"`
	TopicID                   *uuid.UUID                     `json:"topicID"`
	ModuleID                  uuid.UUID                      `json:"-"`
	Module                    *ModuleDTO                     `json:"module" gorm:"foreignkey:ModuleID"`
	ProgrammingConceptID      uuid.UUID                      `json:"-"`
	TopicProgrammingConcept   []*TopicProgrammingConceptDTO  `json:"topicProgrammingConcept" gorm:"foreignkey:TopicID"`
	TopicProgrammingQuestions []*TopicProgrammingQuestionDTO `json:"topicProgrammingQuestions" gorm:"foreignkey:TopicID"`

	BatchTopicAssignment []*list.BatchTopicAssignmentDTO `json:"batchTopicAssignment" gorm:"foreignkey:TopicID"`
	BatchSessionTopic    *list.BatchSessionTopic         `json:"batchSessionTopic" gorm:"foreignkey:SubTopicID"`
}

// TableName overrides name of the table
func (*ModuleTopicDTO) TableName() string {
	return "module_topics"
}
