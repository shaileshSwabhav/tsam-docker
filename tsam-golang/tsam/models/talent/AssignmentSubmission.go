package talent

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

// TableName defines table name of the struct.
func (*AssignmentSubmission) TableName() string {
	return "talent_assignment_submissions"
}

// AssignmentSubmission contains all fields required for talents programming_assignment subsmission.
//  table: talent_assignment_submissions
// IsAccepted change to bool #niranjan
type AssignmentSubmission struct {
	general.TenantBase
	TalentID                    uuid.UUID                     `json:"talentID" gorm:"type:varchar(36)"`
	BatchTopicAssignmentID      uuid.UUID                     `json:"batchTopicAssignmentID" gorm:"type:varchar(36)"`
	IsAccepted                  *bool                         `json:"isAccepted"`
	IsChecked                   *bool                         `json:"isChecked"`
	SubmittedOn                 time.Time                     `json:"submittedOn" gorm:"type:datetime"`
	AcceptanceDate              *time.Time                    `json:"acceptanceDate" gorm:"type:date"`
	FacultyRemarks              *string                       `json:"facultyRemarks" gorm:"type:varchar(1000)"`
	Score                       *float32                      `json:"score" gorm:"type:decimal(3,2)"`
	Solution                    *string                       `json:"solution" gorm:"type:varchar(10000)"`
	GithubURL                   *string                       `json:"githubURL" gorm:"type:varchar(500)"`
	AssignmentSubmissionUploads []*AssignmentSubmissionUpload `json:"assignmentSubmissionUploads"`
	TalentConceptRatings        []*TalentConceptRating        `json:"talentConceptRatings"`

	// Faculty who checks and updates
	FacultyVoiceNote *string    `json:"facultyVoiceNote,omitempty" gorm:"type:varchar(500)"`
	FacultyID        *uuid.UUID `json:"facultyID,omitempty" gorm:"type:varchar(36)"`

	// PreviousSubmissionID *uuid.UUID `json:"previousSubmissionID" gorm:"type:varchar(36)"`
	// BatchSessionProgrammingAssignmentID uuid.UUID                           `json:"batchSessionProgrammingAssignmentID" gorm:"type:varchar(36)"`
}

// Validate will validate all the fields of talent_programming_assignment_submissions.
func (submission *AssignmentSubmission) Validate() error {

	if !util.IsUUIDValid(submission.TalentID) {
		return errors.NewValidationError("Talent ID must be specified.")
	}

	if !util.IsUUIDValid(submission.BatchTopicAssignmentID) {
		return errors.NewValidationError("Session Assignment ID must be specified.")
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
	if (submission.GithubURL == nil || util.IsEmpty(*submission.GithubURL)) && len(submission.AssignmentSubmissionUploads) == 0 {
		return errors.NewValidationError("Either github URL or images must be provided.")
	}

	if submission.AssignmentSubmissionUploads != nil {
		for _, upload := range submission.AssignmentSubmissionUploads {
			err := upload.Validate()
			if err != nil {
				return err
			}
		}
	}

	if submission.TalentConceptRatings == nil && *submission.IsAccepted {
		return errors.NewValidationError("Concept Rating must be speciified")
	}

	return nil
}

// Validate will validate all the fields of talent_programming_assignment_submissions. #Niranjan
func (submission *AssignmentSubmission) ValidateScore() error {

	if submission.FacultyRemarks == nil && submission.FacultyVoiceNote == nil {
		return errors.NewValidationError("Text or audio feedback must be provided.")
	}

	if !util.IsUUIDValid(submission.TalentID) {
		return errors.NewValidationError("Talent ID must be specified.")
	}

	if !util.IsUUIDValid(submission.BatchTopicAssignmentID) {
		return errors.NewValidationError("Session Assignment ID must be specified.")
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

	if submission.TalentConceptRatings == nil && *submission.IsAccepted {
		return errors.NewValidationError("Concept Rating must be speciified")
	}

	return nil
}

// AssignmentSubmissionDTO contains all fields required for GET operation of talents programming_assignment subsmission.
//  table: talent_assignment_submissions
type AssignmentSubmissionDTO struct {
	general.BaseDTO
	TalentID                    uuid.UUID                        `json:"-"`
	Talent                      *list.Talent                     `json:"talent"`
	IsLatestSubmission          bool                             `json:"isLatestSubmission"`
	BatchTopicAssignmentID      uuid.UUID                        `json:"batchTopicAssignmentID"`
	BatchTopicAssignment        *batch.TopicAssignmentDTO        `json:"batchTopicAssignment,omitempty" gorm:"foreignkey:BatchTopicAssignmentID"`
	IsAccepted                  *bool                            `json:"isAccepted,omitempty"`
	IsChecked                   *bool                            `json:"isChecked,omitempty"`
	SubmittedOn                 time.Time                        `json:"submittedOn"`
	AcceptanceDate              *time.Time                       `json:"acceptanceDate,omitempty"`
	FacultyRemarks              *string                          `json:"facultyRemarks,omitempty"`
	Score                       *float32                         `json:"score,omitempty"`
	Solution                    *string                          `json:"solution,omitempty"`
	GithubURL                   *string                          `json:"githubURL,omitempty"`
	AssignmentSubmissionUploads []*AssignmentSubmissionUploadDTO `json:"assignmentSubmissionUploads" gorm:"foreignkey:AssignmentSubmissionID"`
	TalentConceptRatings        []*TalentConceptRatingDTO        `json:"talentConceptRatings" gorm:"foreignkey:TalentSubmissionID"`
	FacultyID                   uuid.UUID                        `json:"facultyID,omitempty"`
	FacultyVoiceNote            *string                          `json:"facultyVoiceNote,omitempty"`
}

// TableName defines table name of the struct.
func (*AssignmentSubmissionDTO) TableName() string {
	return "talent_assignment_submissions"
}

// TalentAssignmentScoreDTO will consist of talent score for all the assignments.
type TalentAssignmentScoreDTO struct {
	Talent           *list.Talent               `json:"talent"`
	TalentSubmission []*AssignmentSubmissionDTO `json:"talentSubmission"`

	// Get session-assignment list thru api cl
	// SessionAssignmentList []*batch.BatchSessionsProgrammingAssignmentDTO `json:"sessionAssignmentList" gorm:"foreignkey:BatchSessionsProgrammingAssignmentID"`
}
