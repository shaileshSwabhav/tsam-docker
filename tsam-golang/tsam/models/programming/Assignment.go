package programming

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingAssignment defines fields required for programming_assignments.
type ProgrammingAssignment struct {
	general.TenantBase
	Title                        string                          `json:"title" gorm:"type:varchar(100)"`
	TaskDescription              string                          `json:"taskDescription" gorm:"type:varchar(1000)"`
	ProgrammingAssignmentType    string                          `json:"programmingAssignmentType" gorm:"type:varchar(20)"`
	TimeRequired                 uint                            `json:"timeRequired" gorm:"type:int"`
	ComplexityLevel              uint8                           `json:"complexityLevel" gorm:"type:tinyint"`
	Score                        uint                            `json:"score" gorm:"type:int"`
	AdditionalComments           *string                         `json:"additionalComments" gorm:"type:varchar(1000)"`
	ProgrammingQuestion          []*ProgrammingQuestion          `json:"programmingQuestion" gorm:"many2many:programming_assignments_programming_questions;association_autocreate:false;association_autoupdate:false"`
	ProgrammingAssignmentSubTask []*ProgrammingAssignmentSubTask `json:"programmingAssignmentSubTask" gorm:"foreignkey:ProgrammingAssignmentID"`
}

// Validate validates all fields of programming_assignments.
func (assignment *ProgrammingAssignment) Validate() error {

	if util.IsEmpty(assignment.Title) {
		return errors.NewValidationError("Programming assignment title must be specified.")
	}

	// if util.IsEmpty(assignment.TaskDescription) {
	// 	return errors.NewValidationError("Programming assignment task description must be specified.")
	// }

	if util.IsEmpty(assignment.ProgrammingAssignmentType) {
		return errors.NewValidationError("Programming assignment type must be specified.")
	}

	// if assignment.TimeRequired > 0 {
	// 	return errors.NewValidationError("Programming assignment time must be specified.")
	// }

	if assignment.ProgrammingAssignmentSubTask != nil {
		for _, subTask := range assignment.ProgrammingAssignmentSubTask {
			err := subTask.Validate()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ProgrammingAssignmentDTO defines fields required for programming_assignments.
type ProgrammingAssignmentDTO struct {
	ID                           uuid.UUID                          `json:"id"`
	DeletedAt                    *time.Time                         `json:"-"`
	Title                        string                             `json:"title"`
	TimeRequired                 uint                               `json:"timeRequired"`
	ComplexityLevel              uint8                              `json:"complexityLevel"`
	Score                        uint                               `json:"score"`
	AdditionalComments           *string                            `json:"additionalComments"`
	ProgrammingQuestion          []*ProgrammingQuestionDTO          `json:"programmingQuestion" gorm:"many2many:programming_assignments_programming_questions;association_jointable_foreignkey:programming_question_id;jointable_foreignkey:programming_assignment_id"`
	ProgrammingAssignmentSubTask []*ProgrammingAssignmentSubTaskDTO `json:"programmingAssignmentSubTask" gorm:"foreignkey:ProgrammingAssignmentID"`
	TaskDescription              string                             `json:"-"`
	ProgrammingAssignmentType    string                             `json:"-"`
	SourceURL                    string                             `json:"-"`
	Source                       *string                            `json:"-"`
}

// resourceID above table

// TableName defines table name of the struct.
func (*ProgrammingAssignmentDTO) TableName() string {
	return "programming_assignments"
}
