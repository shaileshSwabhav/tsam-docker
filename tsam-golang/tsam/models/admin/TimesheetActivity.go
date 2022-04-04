package admin

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

// TimesheetActivity will contain activites for specified timesheet
type TimesheetActivity struct {
	general.TenantBase
	TimesheetID    uuid.UUID  `json:"timesheetID" gorm:"type:varchar(36)"`
	ProjectID      *uuid.UUID `json:"projectID" gorm:"type:varchar(36)"`
	Activity       *string    `json:"activity" gorm:"type:varchar(1000)"`
	SubProjectID   *uuid.UUID `json:"subProjectID" gorm:"type:varchar(36)"`
	HoursNeeded    *float32   `json:"hoursNeeded" gorm:"type:decimal(5,2)"`
	BatchID        *uuid.UUID `json:"batchID" gorm:"type:varchar(36)"`
	BatchSessionID *uuid.UUID `json:"batchSessionID" gorm:"type:varchar(36)"`
	// BatchTopicID *uuid.UUID `json:"batchTopicID" gorm:"type:varchar(36)"`
	// Below options should only be avaiable on update.
	IsBillable        *bool   `json:"isBillable"` // Only Admin can update.
	IsCompleted       *bool   `json:"isCompleted"`
	WorkDone          *string `json:"workDone" gorm:"type:varchar(1000)"`
	NextEstimatedDate *string `json:"nextEstimatedDate" gorm:"type:date"` // if isCompleted is false
}

// ValidateTimesheetActivity will validate all the fields of timesheet activity
func (activity *TimesheetActivity) ValidateTimesheetActivity() error {

	if activity.ProjectID == nil || !util.IsUUIDValid(*activity.ProjectID) {
		return errors.NewValidationError("Project must be specified")
	}
	if activity.Activity == nil || util.IsEmpty(*activity.Activity) {
		return errors.NewValidationError("Activity must be specified")
	}
	return nil
}

// TimesheetActivityDTO will contain activites for specified timesheet
type TimesheetActivityDTO struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenantID"`
	DeletedAt      *time.Time `json:"-"`
	TimesheetID    uuid.UUID  `json:"timesheetID"`
	ProjectID      *uuid.UUID `json:"-"`
	Activity       *string    `json:"activity"`
	SubProjectID   *uuid.UUID `json:"-"`
	HoursNeeded    *float32   `json:"hoursNeeded"`
	BatchID        *uuid.UUID `json:"-"`
	BatchSessionID *uuid.UUID `json:"-"`
	// BatchTopicID *uuid.UUID `json:"-"`
	// Below options should only be avaiable on update.
	IsBillable        *bool   `json:"isBillable"` // Only Admin can update.
	IsCompleted       *bool   `json:"isCompleted"`
	WorkDone          *string `json:"workDone"`
	NextEstimatedDate *string `json:"nextEstimatedDate"` // if isCompleted is false

	Batch        *list.Batch             `json:"batch"`
	BatchSession *batch.SessionDTO       `json:"batchSession"`
	Project      *general.SwabhavProject `json:"project"`
	SubProject   *general.SwabhavProject `json:"subProject"`
	// BatchTopic   *batch.BatchTopicDTO `json:"batchTopic"`
}

// TableName will refer it to timesheets.
func (*TimesheetActivityDTO) TableName() string {
	return "timesheet_activities"
}
