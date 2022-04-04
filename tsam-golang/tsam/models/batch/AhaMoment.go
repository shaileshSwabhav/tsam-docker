package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// AhaMoment contains all information about aha-moment
type AhaMoment struct {
	general.TenantBase
	BatchID        uuid.UUID `json:"-" gorm:"type:varchar(36)"`
	BatchSessionID uuid.UUID `json:"-" gorm:"type:varchar(36)"`
	FeelingID      uuid.UUID `json:"-" gorm:"type:varchar(36)"`
	FeelingLevelID uuid.UUID `json:"-" gorm:"type:varchar(36)"`
	FacultyID      uuid.UUID `json:"facultyID" gorm:"type:varchar(36)"`
	TalentID       uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	Batch          Batch     `json:"batch" gorm:"foreignkey:BatchID;association_autocreate:false;association_autoupdate:false"`
	// BatchTopic     BatchTopic           `json:"batchTopic" gorm:"foreignkey:BatchTopicID;association_autocreate:false;association_autoupdate:false"`
	Talent       TalentDTO            `json:"talent" gorm:"foreignkey:TalentID;association_autocreate:false;association_autoupdate:false"`
	Faculty      faculty.Faculty      `json:"faculty" gorm:"foreignkey:FacultyID;association_autocreate:false;association_autoupdate:false"`
	Feeling      general.Feeling      `json:"feeling" gorm:"foreignkey:FeelingID;association_autocreate:false;association_autoupdate:false"`
	FeelingLevel general.FeelingLevel `json:"feelingLevel" gorm:"foreignkey:FeelingLevelID;association_autocreate:false;association_autoupdate:false"`

	AhaMomentResponse []AhaMomentResponse `json:"ahaMomentResponse" gorm:"foreignkey:AhaMomentID;association_autoupdate:false"`

	// BatchSessionID uuid.UUID            `json:"-" gorm:"type:varchar(36)"`
	// BatchSession   MappedSession        `json:"batchSession" gorm:"foreignkey:BatchSessionID;association_autocreate:false;association_autoupdate:false"`
}

// ValidateAhaMoment will check if all the ID's are valid
func (ahamoment *AhaMoment) ValidateAhaMoment() error {

	// if !util.IsUUIDValid(ahamoment.BatchID) {
	// 	return errors.NewValidationError("Batch ID must be specified")
	// }
	// if !util.IsUUIDValid(ahamoment.BatchSessionID) {
	// 	return errors.NewValidationError("Session ID must be specified")
	// }
	if !util.IsUUIDValid(ahamoment.FacultyID) {
		return errors.NewValidationError("Faculty ID must be specified")
	}
	// if ahamoment.Talents == nil {
	// 	return errors.NewValidationError("Talents must be specified")
	// }
	if !util.IsUUIDValid(ahamoment.TalentID) {
		return errors.NewValidationError("Talent ID must be specified")
	}
	if !util.IsUUIDValid(ahamoment.Feeling.ID) {
		return errors.NewValidationError("Feeling ID must be specified")
	}
	if !util.IsUUIDValid(ahamoment.FeelingLevel.ID) {
		return errors.NewValidationError("Feeling Level ID must be specified")
	}

	if ahamoment.AhaMomentResponse == nil || len(ahamoment.AhaMomentResponse) == 0 {
		return errors.NewValidationError("Aha Moment Response must be speicified")
	}

	for _, response := range ahamoment.AhaMomentResponse {
		err := response.ValidateAhaMomentResponse()
		if err != nil {
			return err
		}
	}

	return nil
}

// AhaMomentDTO contains information about aha-moment
type AhaMomentDTO struct {
	ID             uuid.UUID  `json:"id"`
	DeletedAt      *time.Time `json:"-"`
	BatchID        uuid.UUID  `json:"-"`
	BatchSessionID uuid.UUID  `json:"-"`
	FeelingID      uuid.UUID  `json:"-"`
	FeelingLevelID uuid.UUID  `json:"-"`
	FacultyID      uuid.UUID  `json:"-"`
	TalentID       uuid.UUID  `json:"-"`
	Batch          BatchDTO   `json:"batch"`
	// BatchTopic     BatchTopicDTO        `json:"batchTopic"`
	Talent       TalentDTO            `json:"talent"`
	Faculty      faculty.FacultyDTO   `json:"faculty"`
	Feeling      general.Feeling      `json:"feeling"`
	FeelingLevel general.FeelingLevel `json:"feelingLevel"`

	AhaMomentResponse []AhaMomentResponseDTO `json:"ahaMomentResponse" gorm:"foreignkey:AhaMomentID"`
	// BatchSession   MappedSession        `json:"batchSession"`
}

// TableName defines table name of the struct.
func (*AhaMomentDTO) TableName() string {
	return "aha_moments"
}
