package talent

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProjectSubmission contains all fields required for talents programming_project_subsmissions.
//  table: talent_project_submissions
type ProjectSubmission struct {
	general.TenantBase
	TalentID       uuid.UUID  `json:"talentID" gorm:"type:varchar(36)"`
	BatchID        uuid.UUID  `json:"batchID" gorm:"type:varchar(36)"`
	BatchProjectID uuid.UUID  `json:"batchProjectID" gorm:"type:varchar(36)"`
	IsAccepted     *bool      `json:"isAccepted"`
	IsChecked      *bool      `json:"isChecked"`
	SubmittedOn    time.Time  `json:"submittedOn" gorm:"type:datetime"`
	AcceptanceDate *time.Time `json:"acceptanceDate" gorm:"type:date"`
	FacultyRemarks *string    `json:"facultyRemarks" gorm:"type:varchar(1000)"`
	Score          *float32   `json:"score" gorm:"type:decimal(3,2)"`
	Solution       *string    `json:"solution" gorm:"type:varchar(5000)"`
	ProjectUpload  *string    `json:"projectUpload" gorm:"type:varchar(5000)"`
	GithubURL      *string    `json:"githubURL" gorm:"type:varchar(500)"`
	WebsiteLink    *string    `json:"websiteLink" gorm:"type:varchar(500)"`

	ProjectSubmissionUploads  []*ProjectSubmissionUpload        `json:"projectSubmissionUploads"`
	ProgrammingProjectRatings []*batch.ProgrammingProjectRating `json:"programmingProjectRatings"`
	// TalentConceptRatings     []*TalentConceptRating     `json:"talentConceptRatings"`

	// Faculty who checks and updates
	FacultyVoiceNote *string    `json:"facultyVoiceNote,omitempty" gorm:"type:varchar(100)"`
	FacultyID        *uuid.UUID `json:"facultyID,omitempty" gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*ProjectSubmission) TableName() string {
	return "talent_project_submissions"
}

// Validate will validate all the fields of talent_project_submissions.
func (submission *ProjectSubmission) Validate() error {

	if !util.IsUUIDValid(submission.TalentID) {
		return errors.NewValidationError("Talent ID must be specified.")
	}

	if !util.IsUUIDValid(submission.BatchID) {
		return errors.NewValidationError("Batch ID must be specified.")
	}

	if !util.IsUUIDValid(submission.BatchProjectID) {
		return errors.NewValidationError("Batch Project ID must be specified.")
	}

	if submission.IsAccepted != nil {
		if !(*submission.IsAccepted) {
			submission.AcceptanceDate = nil
			submission.Score = nil
		} else {
			// if submission.AcceptanceDate == nil {
			// 	return errors.NewValidationError("If solution is accepted then acceptance date must be specified.")
			// }
			if submission.Score == nil {
				return errors.NewValidationError("If solution is accepted then score must be specified.")
			}
		}
	}

	if submission.WebsiteLink == nil && submission.GithubURL == nil && submission.ProjectSubmissionUploads == nil {
		return errors.NewValidationError("Either solution or github URL or images must be provided.")
	}

	if submission.ProjectSubmissionUploads != nil {
		for _, upload := range submission.ProjectSubmissionUploads {
			err := upload.Validate()
			if err != nil {
				return err
			}
		}
	}

	// if submission.TalentConceptRatings == nil && *submission.IsAccepted {
	// 	return errors.NewValidationError("Concept Rating must be speciified")
	// }

	return nil
}

// Validate will validate all the fields of talent_programming_assignment_submissions. #Niranjan
func (submission *ProjectSubmission) ValidateScore() error {

	if submission.FacultyRemarks == nil && submission.FacultyVoiceNote == nil {
		return errors.NewValidationError("Text or audio feedback must be provided.")
	}

	if !util.IsUUIDValid(submission.TalentID) {
		return errors.NewValidationError("Talent ID must be specified.")
	}

	if !util.IsUUIDValid(submission.BatchProjectID) {
		return errors.NewValidationError("Project ID must be specified.")
	}

	if submission.IsAccepted != nil {
		if !(*submission.IsAccepted) {
			submission.AcceptanceDate = nil
			submission.Score = nil
		} else {
			// if submission.AcceptanceDate == nil {
			// 	return errors.NewValidationError("If solution is accepted then acceptance date must be specified.")
			// }
			if submission.Score == nil {
				return errors.NewValidationError("If solution is accepted then score must be specified.")
			}
		}
	}
	fmt.Println("submission", submission.ProgrammingProjectRatings)

	if submission.ProgrammingProjectRatings == nil && *submission.IsAccepted {
		return errors.NewValidationError("Project Rating must be speciified")
	}

	return nil
}

// ProjectSubmissionDTO contains all fields required for GET operation of talents programming_assignment subsmission.
//  table: talent_assignment_submissions
type ProjectSubmissionDTO struct {
	general.BaseDTO
	TalentID                  uuid.UUID                            `json:"-"`
	Talent                    *list.Talent                         `json:"talent"`
	BatchID                   uuid.UUID                            `json:"batchID"`
	BatchProjectID            uuid.UUID                            `json:"batchProjectID"`
	Batch                     *list.Batch                          `json:"batch,omitempty" gorm:"foreignkey:BatchID"`
	IsAccepted                *bool                                `json:"isAccepted,omitempty"`
	IsChecked                 *bool                                `json:"isChecked,omitempty"`
	SubmittedOn               time.Time                            `json:"submittedOn"`
	AcceptanceDate            *time.Time                           `json:"acceptanceDate,omitempty"`
	FacultyRemarks            *string                              `json:"facultyRemarks,omitempty"`
	Score                     *float32                             `json:"score,omitempty"`
	Solution                  *string                              `json:"solution,omitempty"`
	ProjectUpload             *string                              `json:"projectUpload,omitempty"`
	GithubURL                 *string                              `json:"githubURL,omitempty"`
	WebsiteLink               *string                              `json:"websiteLink,omitempty"`
	ProjectSubmissionUpload   []*ProjectSubmissionUploadDTO        `json:"projectSubmissionUpload" gorm:"foreignkey:ProjectSubmissionID"`
	FacultyID                 uuid.UUID                            `json:"facultyID,omitempty"`
	ProgrammingProjectRatings []*batch.ProgrammingProjectRatingDTO `json:"programmingProjectRatings"`

	// TalentConceptRating        []*TalentConceptRatingDTO        `json:"talentConceptRating" gorm:"foreignkey:TalentSubmissionID"`
	FacultyVoiceNote *string `json:"facultyVoiceNote,omitempty"`
}

// TableName defines table name of the struct.
func (*ProjectSubmissionDTO) TableName() string {
	return "talent_project_submissions"
}

// ProgrammingProjectRatings will store map for project rating
// type ProgrammingProjectRating struct {
// 	general.TenantBase
// 	TalentID                            uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
// 	BatchID                             uuid.UUID `json:"batchID" gorm:"type:varchar(36)"`
// 	TalentSubmissionID                  uuid.UUID `json:"talentSubmissionID" gorm:"type:varchar(36)"`
// 	ProgrammingProjectRatingParameterID uuid.UUID `json:"programmingProjectRatingParameterID" gorm:"type:varchar(36)"`
// 	ProgrammingProjectRatingParameter   batch.ProgrammingProjectRatingParameter `json:"programmingProjectRatingParameter"`

// 	Score                               uint      `json:"score" gorm:"type:int"`
// }

// type ProgrammingProjectRatingDTO struct {
// 	ID                                  uuid.UUID  `json:"id"`
// 	DeletedAt                           *time.Time `json:"-"`
// 	TalentID                            uuid.UUID  `json:"talentID"`
// 	BatchID                             uuid.UUID  `json:"-"`
// 	TalentSubmissionID                  uuid.UUID  `json:"-"`
// 	ProgrammingProjectRatingParameterID uuid.UUID  `json:"-"`
// 	Score                               uint       `json:"score"`
// }

// func (programmingprojectRating *ProgrammingProjectRating) Validate() error {

// 	if !util.IsUUIDValid(programmingprojectRating.BatchID) {
// 		return errors.NewValidationError("Batch ID must be specified.")
// 	}

// 	if !util.IsUUIDValid(programmingprojectRating.TalentID) {
// 		return errors.NewValidationError("Talent ID must be specified.")
// 	}

// 	if !util.IsUUIDValid(programmingprojectRating.ProgrammingProjectRatingParameterID) {
// 		return errors.NewValidationError("Programming Project Rating Parameter ID must be specified.")
// 	}

// 	if !util.IsUUIDValid(programmingprojectRating.TalentSubmissionID) {
// 		return errors.NewValidationError("Talent submission ID must be specified.")
// 	}

// 	return nil
// }
