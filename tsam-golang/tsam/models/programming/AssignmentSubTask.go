package programming

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingAssignmentSubTask contains all fields required for programming_assignment_sub_tasks.
type ProgrammingAssignmentSubTask struct {
	general.TenantBase
	ProgrammingAssignmentID uuid.UUID `json:"programmingAssignmentID" gorm:"type:varchar(36)"`
	ResourceID              uuid.UUID `json:"resourceID" gorm:"type:varchar(36)"`
	Description             *string   `json:"description" gorm:"type:varchar(1000)"`
	// Source                  string     `json:"source" gorm:"type:varchar(100)"`
	// SourceURL               *string    `json:"sourceURL" gorm:"type:varchar(1000)"`
}

// Validate will validate fields of programming_assignment_sub_task table.
func (task *ProgrammingAssignmentSubTask) Validate() error {

	if !util.IsUUIDValid(task.ResourceID) {
		return errors.NewValidationError("valid resource ID must be specified.")
	}

	// if util.IsEmpty(task.Description) {
	// 	return errors.NewValidationError("Sub task description must be specified.")
	// }

	// if util.IsEmpty(task.Source) {
	// 	return errors.NewValidationError("Source must be specified")
	// }

	return nil
}

// ProgrammingAssignmentSubTask contains all fields required for programming_assignment_sub_tasks.
type ProgrammingAssignmentSubTaskDTO struct {
	ID                      uuid.UUID          `json:"id"`
	DeletedAt               *time.Time         `json:"-"`
	ProgrammingAssignmentID uuid.UUID          `json:"programmingAssignmentID"`
	ResourceID              uuid.UUID          `json:"-"`
	Resource                *resource.Resource `json:"resource" gorm:"foreignkey:ResourceID"`
	Description             string             `json:"description"`
	// Source                  string            `json:"source"`
	// SourceURL               *string           `json:"sourceURL"`
}

func (*ProgrammingAssignmentSubTaskDTO) TableName() string {
	return "programming_assignment_sub_tasks"
}
