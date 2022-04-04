package college

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// SeminarTalentRegistration defines fields of SeminarTalentRegistration.
type SeminarTalentRegistration struct {
	general.TenantBase
	RegistrationDate string    `json:"registrationDate" gorm:"type:date"`
	HasVisited       bool      `json:"hasVisited"`
	SeminarID        uuid.UUID `json:"seminarID" gorm:"type:varchar(36)"`
	TalentID         uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
}

// TableName will refer table "seminar_talent_registrations" for model referrences.
func (*SeminarTalentRegistration) TableName() string {
	return "seminar_talent_registrations"
}

// Validate validates all fields of the SeminarTalentRegistration
func (seminarTalentRegistration *SeminarTalentRegistration) Validate() error {
	// Check if registration date is blank or not.
	if util.IsEmpty(seminarTalentRegistration.RegistrationDate) {
		return errors.NewValidationError("Registration Date name must be specified")
	}

	return nil
}
