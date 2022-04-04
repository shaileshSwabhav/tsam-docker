package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// BatchAssignmentTask will store task to be done for the given assignment.
type BatchAssignmentTask struct {
	general.TenantBase
	ResourceID *uuid.UUID `json:"resourceID" gorm:"type:varchar(36)"`
	TotalTime  uint       `json:"totalTime" gorm:"type:INT"`
	URL        *string    `json:"url" gorm:"type:varchar(200)"`
	Task       string     `json:"task" gorm:"type:varchar(1000)"`
}

func (task *BatchAssignmentTask) Validate() error {

	if util.IsEmpty(task.Task) {
		return errors.NewValidationError("task must be specified")
	}

	return nil
}

// BatchAssignmentTaskDTO will store task to be done for the given assignment.
type BatchAssignmentTaskDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`
	Task      string     `json:"task"`
}

// TableName defines table name of the struct.
func (*BatchAssignmentTaskDTO) TableName() string {
	return "batch_assignment_tasks"
}
