package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/util"
)

// Project will store map for batch and programming-project.
type Project struct {
	general.TenantBase
	BatchID              uuid.UUID  `json:"batchID" gorm:"type:varchar(36)"`
	ProgrammingProjectID uuid.UUID  `json:"programmingProjectID" gorm:"type:varchar(36)"`
	DueDate              *string    `json:"dueDate,omitempty"`
	AssignedDate         *time.Time `json:"assignedDate,omitempty"`
}

// TableName overrides name of the table
func (*Project) TableName() string {
	return "batch_projects"
}

// Validate will check if all fields are valid in table.
func (project *Project) Validate() error {

	if !util.IsUUIDValid(project.BatchID) {
		return errors.NewValidationError("invalid batchID")
	}

	if !util.IsUUIDValid(project.ProgrammingProjectID) {
		return errors.NewValidationError("invalid programming project ID")
	}

	return nil
}

type ProjectDTO struct {
	ID                   uuid.UUID                         `json:"id"`
	DeletedAt            *time.Time                        `json:"-"`
	BatchID              uuid.UUID                         `json:"batchID"`
	Batch                *list.Batch                       `json:"batch" gorm:"foreignkey:BatchID"`
	ProgrammingProjectID uuid.UUID                         `json:"programmingProjectID"`
	ProgrammingProject   programming.ProgrammingProjectDTO `json:"programmingProject" gorm:"foreignkey:ProgrammingProjectID"`
	DueDate              *string                           `json:"dueDate"`
	AssignedDate         *time.Time                        `json:"assignedDate"`
	Submissions          []talentProjectSubmissionDTO      `json:"submissions,omitempty" gorm:"foreignkey:BatchProjectID"`
}

// TableName overrides name of the table
func (*ProjectDTO) TableName() string {
	return "batch_projects"
}

type talentProjectSubmissionDTO struct {
	general.BaseDTO
	TalentID                  uuid.UUID                           `json:"-"`
	Talent                    *list.Talent                        `json:"talent"`
	BatchID                   uuid.UUID                           `json:"batchID"`
	BatchProjectID            uuid.UUID                           `json:"batchProjectID"`
	Batch                     *list.Batch                         `json:"batch,omitempty" gorm:"foreignkey:BatchID"`
	IsAccepted                *bool                               `json:"isAccepted,omitempty"`
	IsChecked                 *bool                               `json:"isChecked,omitempty"`
	SubmittedOn               time.Time                           `json:"submittedOn"`
	AcceptanceDate            *time.Time                          `json:"acceptanceDate,omitempty"`
	FacultyRemarks            *string                             `json:"facultyRemarks,omitempty"`
	FacultyVoiceNote          *string                             `json:"facultyVoiceNote,omitempty"`
	Score                     *float32                            `json:"score,omitempty"`
	Solution                  *string                             `json:"solution,omitempty"`
	GithubURL                 *string                             `json:"githubURL,omitempty"`
	WebsiteLink               *string                             `json:"websiteLink,omitempty"`
	ProjectUpload             *string                             `json:"projectUpload,omitempty"`
	DueDate                   *string                             `json:"dueDate,omitempty"`
	AssignedDate              *time.Time                          `json:"assignedDate,omitempty"`
	ProjectSubmissionUpload   []*talentProjectSubmissionUploadDTO `json:"projectSubmissionUpload" gorm:"foreignkey:ProjectSubmissionID"`
	ProgrammingProjectRatings []*ProgrammingProjectRating         `json:"programmingProjectRatings,omitempty" gorm:"foreignkey:TalentSubmissionID"`
	FacultyID                 uuid.UUID                           `json:"facultyID,omitempty"`
	// TalentConceptRating        []*TalentConceptRatingDTO        `json:"talentConceptRating" gorm:"foreignkey:TalentSubmissionID"`
	// FacultyVoiceNote           *string                          `json:"facultyVoiceNote,omitempty"`
}

// TableName defines table name of the struct.
func (*talentProjectSubmissionDTO) TableName() string {
	return "talent_project_submissions"
}

type talentProjectSubmissionUploadDTO struct {
	ID                  uuid.UUID  `json:"id"`
	DeletedAt           *time.Time `json:"-"`
	ProjectSubmissionID uuid.UUID  `json:"projectSubmissionID"`
	ImageURL            string     `json:"imageURL"`
	Description         *string    `json:"description"`
}

// TableName defines table name of the struct.
func (*talentProjectSubmissionUploadDTO) TableName() string {
	return "talent_project_submission_uploads"
}
