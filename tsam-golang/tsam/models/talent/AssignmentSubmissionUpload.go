package talent

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// AssignmentSubmissionUpload will contain images for assignment submission.
//  table: talent_assignment_submission_uploads
type AssignmentSubmissionUpload struct {
	general.TenantBase
	AssignmentSubmissionID uuid.UUID `json:"assignmentSubmissionID" gorm:"type:varchar(36)"`
	ImageURL               string    `json:"imageURL" gorm:"type:varchar(500)"`
	Description            *string   `json:"description" gorm:"type:varchar(1000)"`
}

// TableName defines table name of the struct.
func (*AssignmentSubmissionUpload) TableName() string {
	return "talent_assignment_submission_uploads"
}

// Validate will validate all the fields of talent_assignment_submission_uploads.
func (upload *AssignmentSubmissionUpload) Validate() error {

	if util.IsEmpty(upload.ImageURL) {
		return errors.NewValidationError("Image URL must be specified.")
	}

	return nil
}

// AssignmentSubmissionUpload will contain images for assignment submission.
//  table: talent_assignment_submission_uploads
type AssignmentSubmissionUploadDTO struct {
	ID                     uuid.UUID  `json:"id"`
	DeletedAt              *time.Time `json:"-"`
	AssignmentSubmissionID uuid.UUID  `json:"AssignmentSubmissionID"`
	ImageURL               string     `json:"imageURL"`
	Description            *string    `json:"description"`
}

// TableName defines table name of the struct.
func (*AssignmentSubmissionUploadDTO) TableName() string {
	return "talent_assignment_submission_uploads"
}
