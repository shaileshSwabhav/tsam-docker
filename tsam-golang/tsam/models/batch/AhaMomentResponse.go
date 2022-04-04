package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// AhaMomentResponse will contains response for the aha moments
type AhaMomentResponse struct {
	general.TenantBase
	BatchID        uuid.UUID                `json:"-" gorm:"type:varchar(36)"`
	BatchSessionID uuid.UUID                `json:"-" gorm:"type:varchar(36)"`
	FeelingID      uuid.UUID                `json:"-" gorm:"type:varchar(36)"`
	FeelingLevelID uuid.UUID                `json:"-" gorm:"type:varchar(36)"`
	QuestionID     uuid.UUID                `json:"-" gorm:"type:varchar(36)"`
	FacultyID      uuid.UUID                `json:"facultyID" gorm:"type:varchar(36)"`
	TalentID       uuid.UUID                `json:"talentID" gorm:"type:varchar(36)"`
	Question       general.FeedbackQuestion `json:"question" gorm:"foreignkey:QuestionID;association_autocreate:false;association_autoupdate:false"`
	Batch          Batch                    `json:"batch" gorm:"foreignkey:BatchID;association_autocreate:false;association_autoupdate:false"`
	// BatchTopic     BatchTopic               `json:"batchTopic" gorm:"foreignkey:BatchTopicID;association_autocreate:false;association_autoupdate:false"`
	Talent       TalentDTO            `json:"talent" gorm:"foreignkey:TalentID;association_autocreate:false;association_autoupdate:false"`
	Faculty      faculty.Faculty      `json:"faculty" gorm:"foreignkey:FacultyID;association_autocreate:false;association_autoupdate:false"`
	Feeling      general.Feeling      `json:"feeling" gorm:"foreignkey:FeelingID;association_autocreate:false;association_autoupdate:false"`
	FeelingLevel general.FeelingLevel `json:"feelingLevel" gorm:"foreignkey:FeelingLevelID;association_autocreate:false;association_autoupdate:false"`
	Response     string               `json:"response" gorm:"type:varchar(200)"`
	AhaMomentID  uuid.UUID            `json:"-" gorm:"type:varchar(36)"`

	// BatchSessionID uuid.UUID                `json:"-" gorm:"type:varchar(36)"`
	// BatchSession   MappedSession            `json:"session" gorm:"foreignkey:BatchSessionID;association_autocreate:false;association_autoupdate:false"`
}

// ValidateAhaMomentResponse will check if all the ID's and response are valid
func (response *AhaMomentResponse) ValidateAhaMomentResponse() error {

	if !util.IsUUIDValid(response.Question.ID) {
		return errors.NewValidationError("Question ID must be specified")
	}
	if util.IsEmpty(response.Response) {
		return errors.NewValidationError("Response must be specified")
	}
	return nil
}

// AhaMomentResponse will contains response for the aha moments
type AhaMomentResponseDTO struct {
	ID             uuid.UUID                    `json:"id"`
	DeletedAt      *time.Time                   `json:"-"`
	BatchID        uuid.UUID                    `json:"-"`
	BatchSessionID uuid.UUID                    `json:"-"`
	FeelingID      uuid.UUID                    `json:"-"`
	FeelingLevelID uuid.UUID                    `json:"-"`
	QuestionID     uuid.UUID                    `json:"-"`
	FacultyID      uuid.UUID                    `json:"-"`
	TalentID       uuid.UUID                    `json:"-"`
	Question       *general.FeedbackQuestionDTO `json:"question" gorm:"foreignkey:QuestionID"`
	Batch          *BatchDTO                    `json:"batch"`
	// BatchTopic     *BatchTopic                  `json:"batchTopic"`
	Talent       *TalentDTO            `json:"talent"`
	Faculty      *faculty.FacultyDTO   `json:"faculty"`
	Feeling      *general.Feeling      `json:"feeling"`
	FeelingLevel *general.FeelingLevel `json:"feelingLevel"`
	Response     string                `json:"response"`
	AhaMomentID  uuid.UUID             `json:"-"`

	// BatchSession   *MappedSession               `json:"session"`
}

// TableName defines table name of the struct.
func (*AhaMomentResponseDTO) TableName() string {
	return "aha_moment_responses"
}
