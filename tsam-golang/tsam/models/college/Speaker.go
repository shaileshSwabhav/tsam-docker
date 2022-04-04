package college

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Speaker defines the fields of a speaker required for adding and updating the model.
type Speaker struct {
	general.TenantBase

	// Related table IDs.
	DesignationID *uuid.UUID `json:"designationID" gorm:"type:varchar(36)"`

	// Other fields.
	FirstName         string  `json:"firstName" example:"Ravi" gorm:"type:varchar(50)"`
	LastName          string  `json:"lastName" example:"Sharma" gorm:"type:varchar(50)"`
	Company           *string `json:"company" gorm:"type:varchar(200)"`
	ExperienceInYears *uint16 `json:"experienceInYears" gorm:"type:int"`
}

// Validate validates the fields of speaker.
func (speaker *Speaker) Validate() error {

	// Check if first name is blank or not.
	if util.IsEmpty(speaker.FirstName) {
		return errors.NewValidationError("First name must be specified")
	}

	// First name should consist of only alphabets.
	if !util.ValidateString(speaker.FirstName) {
		return errors.NewValidationError("First name should consist of only alphabets")
	}

	// First name maximum characters.
	if len(speaker.FirstName) > 50 {
		return errors.NewValidationError("First name can have maximum 50 characters")
	}

	// Check if last name is blank or not.
	if util.IsEmpty(speaker.LastName) {
		return errors.NewValidationError("Last name must be specified")
	}

	// Last name should consist of only alphabets.
	if !util.ValidateString(speaker.LastName) {
		return errors.NewValidationError("Last name should consist of only alphabets")
	}

	// Last name maximum characters.
	if len(speaker.LastName) > 50 {
		return errors.NewValidationError("Last name can have maximum 50 characters")
	}

	// Company name maximum length.
	if speaker.Company != nil && len(*speaker.Company) > 255 {
		return errors.NewValidationError("Company name can have maximum 200 characters")
	}

	// Experience In Years minimum.
	if speaker.ExperienceInYears != nil && *speaker.ExperienceInYears > 60 {
		return errors.NewValidationError("Experience In Years cannot be greater than 60")
	}

	// Experience In Years maximum.
	if speaker.ExperienceInYears != nil && *speaker.ExperienceInYears < 0 {
		return errors.NewValidationError("Experience In Years cannot be lesser than 0")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// SpeakerDTO defines the fields of a speaker for getting the model.
type SpeakerDTO struct {
	general.TenantBase

	// Related tables.
	Designation   *general.Designation `json:"designation" gorm:"foreignkey:DesignationID"`
	DesignationID *uuid.UUID           `json:"-"`

	// Other fields.
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
	Company           *string `json:"company"`
	ExperienceInYears *uint16 `json:"experienceInYears"`
}

// TableName defines table name of the struct.
func (*SpeakerDTO) TableName() string {
	return "speakers"
}
