package talent

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// TalentEventRegistration will consist of registration details of talent.
type TalentEventRegistration struct {
	general.TenantBase
	TalentID         uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	EventID          uuid.UUID `json:"eventID" gorm:"type:varchar(36)"`
	RegistrationDate string    `json:"registrationDate" gorm:"type:date"`
	HasAttended      *bool     `json:"hasAttended"`
}

// Validate will validate all the fields of talent event registration struct.
func (registration *TalentEventRegistration) Validate() error {

	if registration.TalentID == uuid.Nil {
		return errors.NewValidationError("Talent must be specified")
	}

	if registration.EventID == uuid.Nil {
		return errors.NewValidationError("Event must be specified")
	}

	if util.IsEmpty(registration.RegistrationDate) {
		return errors.NewValidationError("Registration date must be specified")
	}

	return nil
}

// TalentEventRegistrationDTO will consist of registration details of talent.
type TalentEventRegistrationDTO struct {
	ID                 uuid.UUID             `json:"id"`
	DeletedAt          *time.Time            `json:"-"`
	TalentID           uuid.UUID             `json:"talentID"`
	EventID            uuid.UUID             `json:"eventID"`
	Talent             DTO                   `json:"talent" gorm:"foreignkey:TalentID"`
	Event              admin.SwabhavEventDTO `json:"event" gorm:"foreignkey:EventID"`
	RegistrationDate   string                `json:"registrationDate"`
	HasAttended        *bool                 `json:"hasAttended"`
	IsTalentRegistered bool                  `json:"isTalentRegistered"`
}

// TableName will refer it to events.
func (*TalentEventRegistrationDTO) TableName() string {
	return "talent_event_registrations"
}
