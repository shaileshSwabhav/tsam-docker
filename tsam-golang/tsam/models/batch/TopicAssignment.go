package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/util"
)

// TopicAssignment contains assignment for specified sub_topics.
//  table: batch_topic_assignments
type TopicAssignment struct {
	general.TenantBase
	TopicID               uuid.UUID  `json:"topicID" gorm:"type:varchar(36)"`
	BatchID               uuid.UUID  `json:"batchID" gorm:"type:varchar(36)"`
	ProgrammingQuestionID uuid.UUID  `json:"programmingQuestionID" gorm:"type:varchar(36)"`
	DueDate               *string    `json:"dueDate" gorm:"type:date"`
	AssignedDate          *time.Time `json:"assignedDate" gorm:"type:date"`
	BatchSessionID        *uuid.UUID `json:"batchSessionID" gorm:"type:varchar(36)"`
	ModuleID              uuid.UUID  `json:"moduleID" gorm:"type:varchar(36)"`

	// BatchAssignmentTask   *BatchAssignmentTask `json:"batchAssignmentTask" gorm:"foreignkey:BatchAssignmentTaskID"`
	// BatchAssignmentTaskID *uuid.UUID           `json:"-" gorm:"type:varchar(36)"`
	// DueDate               string               `json:"dueDate" gorm:"type:datetime"`
	// Order                 uint                 `json:"order" gorm:"type:int"`
}

// TableName defines table name of the struct.
func (*TopicAssignment) TableName() string {
	return "batch_topic_assignments"
}

// Validate will validate all the fields of batch_topic_assignments.
func (assignment *TopicAssignment) Validate() error {

	if !util.IsUUIDValid(assignment.BatchID) {
		return errors.NewValidationError("Batch ID must be specified.")
	}

	if !util.IsUUIDValid(assignment.TopicID) {
		return errors.NewValidationError("Batch topic ID must be specified.")
	}

	if !util.IsUUIDValid(assignment.ProgrammingQuestionID) {
		return errors.NewValidationError("Programming assignment ID must be specified.")
	}

	// if util.IsEmpty(assignment.DueDate) {
	// 	return errors.NewValidationError("Due date must be specified.")
	// }

	// if assignment.BatchAssignmentTask != nil {
	// 	if err := assignment.BatchAssignmentTask.Validate(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// TopicAssignmentDTO contains assignments for specified sub_topics.
//  table: batch_topic_assignments
type TopicAssignmentDTO struct {
	general.BaseDTO
	BatchID               uuid.UUID                           `json:"batchID"`
	TopicID               uuid.UUID                           `json:"topicID"`
	ProgrammingQuestionID *uuid.UUID                          `json:"-"`
	ProgrammingQuestion   *programming.ProgrammingQuestionDTO `json:"programmingQuestion" gorm:"foreignkey:ProgrammingQuestionID"`
	DueDate               string                              `json:"dueDate,omitempty"`
	AssignedDate          *time.Time                          `json:"assignedDate,omitempty"`
	Submissions           []talentAssignmentSubmission        `json:"submissions" gorm:"foreignkey:BatchTopicAssignmentID"`
	BatchSessionID        *uuid.UUID                          `json:"batchSessionID"`
	ModuleID              uuid.UUID                           `json:"moduleID"`
	Topic                 course.ModuleTopic                  `json:"topic"`
	// Order                 uint                                `json:"order"`
	// BatchAssignmentTaskID *uuid.UUID                          `json:"-"`
	// BatchAssignmentTask   *BatchAssignmentTask                `json:"batchAssignmentTask" gorm:"foreignkey:BatchAssignmentTaskID"`
	// BatchTopicID          uuid.UUID                           `json:"batchTopicID"`
	// BatchTopic            *BatchTopic                         `json:"batchTopic"`
}

// TableName defines table name of the struct.
func (*TopicAssignmentDTO) TableName() string {
	return "batch_topic_assignments"
}

// talentAssignmentSubmission contains all fields required for GET operation of talents programming_assignment subsmission.
//  table: talent_assignment_submissions
type talentAssignmentSubmission struct {
	general.BaseDTO
	TalentID                    uuid.UUID                              `json:"-"`
	Talent                      *TalentDTO                             `json:"talent,omitempty"`
	IsLatestSubmission          bool                                   `json:"isLatestSubmission"`
	IsAccepted                  *bool                                  `json:"isAccepted,omitempty"`
	IsChecked                   *bool                                  `json:"isChecked,omitempty"`
	SubmittedOn                 time.Time                              `json:"submittedOn,omitempty"`
	AcceptanceDate              *string                                `json:"acceptanceDate,omitempty"`
	FacultyRemarks              *string                                `json:"facultyRemarks,omitempty"`
	Score                       *float32                               `json:"score,omitempty"`
	Solution                    *string                                `json:"solution,omitempty"`
	GithubURL                   *string                                `json:"githubURL,omitempty"`
	FacultyVoiceNote            *string                                `json:"facultyVoiceNote,omitempty"`
	BatchTopicAssignmentID      uuid.UUID                              `json:"batchTopicAssignmentID"`
	AssignmentSubmissionUploads []*talentAssignmentSubmissionUploadDTO `json:"assignmentSubmissionUploads,omitempty" gorm:"foreignkey:AssignmentSubmissionID"`
	TalentConceptRatings        []*talentTalentConceptRatingDTO        `json:"talentConceptRatings,omitempty" gorm:"foreignkey:TalentSubmissionID"`
}

// talentAssignmentSubmissionUploadDTO will contain images for assignment submission.
type talentAssignmentSubmissionUploadDTO struct {
	general.TenantBase
	AssignmentSubmissionID uuid.UUID `json:"assignmentSubmissionID" gorm:"type:varchar(36)"`
	ImageURL               string    `json:"imageURL" gorm:"type:varchar(500)"`
	Description            *string   `json:"description" gorm:"type:varchar(1000)"`
}

// TableName defines table name of the struct.
func (*talentAssignmentSubmissionUploadDTO) TableName() string {
	return "talent_assignment_submission_uploads"
}

// talentTalentConceptRatingDTO will contain concept modules for assignment submission.
type talentTalentConceptRatingDTO struct {
	ModuleProgrammingConcept   programming.ModuleProgrammingConceptsDTO `json:"programmingConceptModule" gorm:"foreignkey:ModuleProgrammingConceptID"`
	ID                         uuid.UUID                                `json:"id"`
	DeletedAt                  *time.Time                               `json:"-"`
	ModuleProgrammingConceptID uuid.UUID                                `json:"-"`
	TalentID                   uuid.UUID                                `json:"talentID"`
	TalentSubmissionID         uuid.UUID                                `json:"-"`
	Score                      float32                                  `json:"score"`
}

// TableName defines table name of the struct.
func (*talentTalentConceptRatingDTO) TableName() string {
	return "talent_concept_ratings"
}

// // TableName defines table name of the struct.
// func (*talentAssignmentSubmission) TableName() string {
// 	return "talent_assignment_submissions"
// }
