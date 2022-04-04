package college

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// CampusTalentRegistration defines fields of CampusTalentRegistration.
type CampusTalentRegistration struct {
	general.TenantBase
	RegistrationDate string    `json:"registrationDate" gorm:"type:date"`
	IsTestLinkSent   bool      `json:"isTestLinkSent"`
	HasAttempted     bool      `json:"hasAttempted"`
	Result           *string   `json:"result" gorm:"type:varchar(100)"`
	CampusDriveID    uuid.UUID `json:"campusDriveID" gorm:"type:varchar(36)"`
	TalentID         uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
}

// TableName will refer table "campus_talent_registrations" for model referrences.
func (*CampusTalentRegistration) TableName() string {
	return "campus_talent_registrations"
}

// Validate validates all fields of the CampusTalentRegistration.
func (campusTalentRegistration *CampusTalentRegistration) Validate() error {
	// Check if registration date is blank or not.
	if util.IsEmpty(campusTalentRegistration.RegistrationDate) {
		return errors.NewValidationError("Registration Date name must be specified")
	}

	return nil
}
