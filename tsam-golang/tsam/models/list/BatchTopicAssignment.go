package list

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/programming"
)

type BatchTopicAssignmentDTO struct {
	ID                    uuid.UUID                        `json:"id"`
	DeletedAt             *time.Time                       `json:"-"`
	TopicID               uuid.UUID                        `json:"topicID"`
	ProgrammingQuestionID *uuid.UUID                       `json:"programmingQuestionID"`
	BatchID               uuid.UUID                        `json:"batchID"`
	ProgrammingQuestion   *programming.ProgrammingQuestion `json:"programmingQuestion" gorm:"foreignkey:ProgrammingQuestionID"`
	DueDate               *string                          `json:"dueDate" gorm:"type:date"`
	AssignedDate          *string                          `json:"assignedDate" gorm:"type:date"`

	// Days                  uint                             `json:"days"`
	// Order                 uint                    `json:"order"`
	// BatchAssignmentTaskID *uuid.UUID              `json:"-" gorm:"type:varchar(36)"`
	// BatchAssignmentTask   *BatchAssignmentTaskDTO `json:"batchAssignmentTask" gorm:"foreignkey:BatchAssignmentTaskID"`
}

// TableName defines table name of the struct.
func (*BatchTopicAssignmentDTO) TableName() string {
	return "batch_topic_assignments"
}

// BatchAssignmentTaskDTO will store task to be done for the given assignment.
// type BatchAssignmentTaskDTO struct {
// 	ID        uuid.UUID  `json:"id"`
// 	DeletedAt *time.Time `json:"-"`
// 	Task      string     `json:"task"`
// }

// // TableName defines table name of the struct.
// func (*BatchAssignmentTaskDTO) TableName() string {
// 	return "batch_assignment_tasks"
// }
