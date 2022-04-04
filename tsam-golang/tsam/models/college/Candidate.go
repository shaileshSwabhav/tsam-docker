package college

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Candidate acts as DTO for fields of talent and campus talent registration.
type Candidate struct {

	// Common details.
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
	DeletedBy uuid.UUID
	TenantID  uuid.UUID

	// Campus talent registration details.
	RegistrationDate           string    `json:"registrationDate"`
	IsTestLinkSent             bool      `json:"isTestLinkSent"`
	HasAttempted               bool      `json:"hasAttempted"`
	Result                     *string   `json:"result"`
	CampusDriveID              uuid.UUID `json:"campusDriveID"`
	CampusTalentRegistrationID uuid.UUID `json:"campusTalentRegistrationID"`

	// Talent personal details.
	general.Address
	TalentID        uuid.UUID `json:"talentID"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email" example:"john@doe.com"`
	Contact         string    `json:"contact" example:"1234567890"`
	AcademicYear    uint8     `json:"academicYear" example:"1"`
	Resume          *string   `json:"resume"`
	IsSwabhavTalent bool      `json:"isSwabhavTalent"`

	// Talent college details.
	CollegeID        *uuid.UUID `json:"collegeID"`
	College          string     `json:"college"`
	Percentage       float32    `json:"percentage"`
	Passout          uint16     `json:"passout"`
	DegreeID         uuid.UUID  `json:"degreeID"`
	SpecializationID uuid.UUID  `json:"specializationID"`
	TalentAcademicID uuid.UUID  `json:"talentAcademicID"`
}

// Validate validates all fields of the Candidate.
func (candidate *Candidate) Validate() error {

	// Check if first name is blank or not.
	if util.IsEmpty(candidate.FirstName) {
		return errors.NewValidationError("First name must be specified")
	}

	// First name should consist of only alphabets.
	if !util.ValidateString(candidate.FirstName) {
		return errors.NewValidationError("First name should consist of only alphabets")
	}

	// Check if last name is blank or not.
	if util.IsEmpty(candidate.LastName) {
		return errors.NewValidationError("Last name must be specified")
	}

	// Last name should consist of only alphabets.
	if !util.ValidateString(candidate.LastName) {
		return errors.NewValidationError("Last name should consist of only alphabets")
	}

	// Check if email is blank or not.
	if util.IsEmpty(candidate.Email) {
		return errors.NewValidationError("Email must be specified")
	}

	// Email should be in proper format.
	if !util.ValidateEmail(candidate.Email) {
		return errors.NewValidationError("Email must in the format : email@example.com")
	}

	// Check if contact is blank or not.
	if util.IsEmpty(candidate.Contact) {
		return errors.NewValidationError("Contact must be specified")
	}

	// Contact should consist of only 10 numbers.
	if !util.ValidateContact(candidate.Contact) {
		return errors.NewValidationError("Contact should consist of only 10 numbers")
	}

	// Check if academic year is present or not.
	if candidate.AcademicYear == 0 {
		return errors.NewValidationError("Academic Year must be specified")
	}

	// Check if registration date is blank or not.
	if util.IsEmpty(candidate.RegistrationDate) {
		return errors.NewValidationError("Registration Date must be specified")
	}

	// Validate address.
	if err := candidate.Address.ValidateAddress(); err != nil {
		return err
	}

	// Validate Degree ID.
	if !util.IsUUIDValid(candidate.DegreeID) {
		return errors.NewValidationError("Degree ID must ne a proper uuid")
	}

	// Check if college name is specified or not.
	if util.IsEmpty(candidate.College) {
		return errors.NewValidationError("College Name must be specified")
	}

	// Check if percentage is specified or not.
	if candidate.Percentage == 0 {
		return errors.NewValidationError("Percentage must be specified")
	}

	// Check if passout year is specified or not.
	if candidate.Passout == 0 {
		return errors.NewValidationError("Passout must be specified")
	}

	// Validate Specialization ID.
	if !util.IsUUIDValid(candidate.SpecializationID) {
		return errors.NewValidationError("Specialization ID must ne a proper uuid")
	}

	return nil
}

// UpdateMultipleCandidate has fields for updating multiple campus talent registrations.
type UpdateMultipleCandidate struct {
	UpdatedBy                   uuid.UUID
	TenantID                    uuid.UUID
	IsTestLinkSent              *bool       `json:"isTestLinkSent"`
	HasAttempted                *bool       `json:"hasAttempted"`
	Result                      *string     `json:"result"`
	CampusDriveID               uuid.UUID   `json:"campusDriveID"`
	CampusTalentRegistrationIDs []uuid.UUID `json:"campusTalentRegistrationIDs"`
}

//************************************* DTO MODEL *************************************************************

// CandidateDTO acts as DTO for fields of talent and campus talent registration.
type CandidateDTO struct {

	// Common details.
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Campus talent registration details.
	RegistrationDate           string    `json:"registrationDate"`
	IsTestLinkSent             bool      `json:"isTestLinkSent"`
	HasAttempted               bool      `json:"hasAttempted"`
	Result                     *string   `json:"result"`
	CampusDriveID              uuid.UUID `json:"campusDriveID"`
	CampusTalentRegistrationID uuid.UUID `json:"campusTalentRegistrationID"`

	// Talent personal details.
	general.Address
	TalentID        uuid.UUID `json:"talentID"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email" example:"john@doe.com"`
	Contact         string    `json:"contact" example:"1234567890"`
	AcademicYear    uint8     `json:"academicYear" example:"1"`
	Resume          *string   `json:"resume"`
	IsSwabhavTalent *bool     `json:"isSwabhavTalent"`

	// Talent college details.
	College          string                 `json:"college"`
	Percentage       float32                `json:"percentage"`
	Passout          uint16                 `json:"passout"`
	Degree           general.Degree         `json:"degree" gorm:"foreignkey:DegreeID"`
	Specialization   general.Specialization `json:"specialization" gorm:"foreignkey:SpecializationID"`
	TalentAcademicID uuid.UUID              `json:"talentAcademicID"`
}
