package talent

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProjectSubmissionUpload will contain images for project submission.
//  table: talent_project_submission_uploads
type ProjectSubmissionUpload struct {
	general.TenantBase
	ProjectSubmissionID uuid.UUID `json:"projectSubmissionID" gorm:"type:varchar(36)"`
	ImageURL            string    `json:"imageURL" gorm:"type:varchar(500)"`
	Description         *string   `json:"description" gorm:"type:varchar(1000)"`
}

// TableName defines table name of the struct.
func (*ProjectSubmissionUpload) TableName() string {
	return "talent_project_submission_uploads"
}

// Validate will validate all the fields of talent_project_submission_uploads.
func (upload *ProjectSubmissionUpload) Validate() error {

	if util.IsEmpty(upload.ImageURL) {
		return errors.NewValidationError("Image URL must be specified.")
	}

	return nil
}

// ProjectSubmissionUpload will contain images for project submission.
//  table: talent_project_submission_uploads
type ProjectSubmissionUploadDTO struct {
	ID                  uuid.UUID  `json:"id"`
	DeletedAt           *time.Time `json:"-"`
	ProjectSubmissionID uuid.UUID  `json:"projectSubmissionID"`
	ImageURL            string     `json:"imageURL"`
	Description         *string    `json:"description"`
}

// TableName defines table name of the struct.
func (*ProjectSubmissionUploadDTO) TableName() string {
	return "talent_project_submission_uploads"
}
