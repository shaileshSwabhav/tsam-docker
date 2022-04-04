package college

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Student acts as DTO for fields of talent and seminar talent registration.
type Student struct {

	// Common details.
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
	DeletedBy uuid.UUID
	TenantID  uuid.UUID

	// Seminar talent registration details.
	RegistrationDate            string    `json:"registrationDate"`
	HasVisited                  bool      `json:"hasVisited"`
	SeminarID                   uuid.UUID `json:"seminarID"`
	SeminarTalentRegistrationID uuid.UUID `json:"seminarTalentRegistrationID"`

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

// Validate validates all fields of the Student.
func (student *Student) Validate() error {

	// Check if first name is blank or not.
	if util.IsEmpty(student.FirstName) {
		return errors.NewValidationError("First name must be specified")
	}

	// First name should consist of only alphabets.
	if !util.ValidateString(student.FirstName) {
		return errors.NewValidationError("First name should consist of only alphabets")
	}

	// Check if last name is blank or not.
	if util.IsEmpty(student.LastName) {
		return errors.NewValidationError("Last name must be specified")
	}

	// Last name should consist of only alphabets.
	if !util.ValidateString(student.LastName) {
		return errors.NewValidationError("Last name should consist of only alphabets")
	}

	// Check if email is blank or not.
	if util.IsEmpty(student.Email) {
		return errors.NewValidationError("Email must be specified")
	}

	// Email should be in proper format.
	if !util.ValidateEmail(student.Email) {
		return errors.NewValidationError("Email must in the format : email@example.com")
	}

	// Check if contact is blank or not.
	if util.IsEmpty(student.Contact) {
		return errors.NewValidationError("Contact must be specified")
	}

	// Contact should consist of only 10 numbers.
	if !util.ValidateContact(student.Contact) {
		return errors.NewValidationError("Contact should consist of only 10 numbers")
	}

	// Check if academic year is present or not.
	if student.AcademicYear == 0 {
		return errors.NewValidationError("Academic Year must be specified")
	}

	// Check if registration date is blank or not.
	if util.IsEmpty(student.RegistrationDate) {
		return errors.NewValidationError("Registration Date must be specified")
	}

	// Validate address.
	if err := student.Address.ValidateAddress(); err != nil {
		return err
	}

	// Validate Degree ID.
	if !util.IsUUIDValid(student.DegreeID) {
		return errors.NewValidationError("Degree ID must ne a proper uuid")
	}

	// Check if college name is specified or not.
	if util.IsEmpty(student.College) {
		return errors.NewValidationError("College Name must be specified")
	}

	// Check if percentage is specified or not.
	if student.Percentage == 0 {
		return errors.NewValidationError("Percentage must be specified")
	}

	// Check if passout year is specified or not.
	if student.Passout == 0 {
		return errors.NewValidationError("Passout must be specified")
	}

	// Validate Specialization ID.
	if !util.IsUUIDValid(student.SpecializationID) {
		return errors.NewValidationError("Specialization ID must ne a proper uuid")
	}

	return nil
}

// UpdateMultipleStudent has fields for updating multiple seminar talent registrations.
type UpdateMultipleStudent struct {
	UpdatedBy                    uuid.UUID
	TenantID                     uuid.UUID
	HasVisited                   *bool       `json:"hasVisited"`
	SeminarID                    uuid.UUID   `json:"seminarID"`
	SeminarTalentRegistrationIDs []uuid.UUID `json:"seminarTalentRegistrationIDs"`
}

//************************************* DTO MODEL *************************************************************

// StudentDTO acts as DTO for fields of talent and seminar talent registration.
type StudentDTO struct {

	// Common details.
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Campus talent registration details.
	RegistrationDate            string    `json:"registrationDate"`
	HasVisited                  bool      `json:"hasVisited"`
	SeminarID                   uuid.UUID `json:"seminarID"`
	SeminarTalentRegistrationID uuid.UUID `json:"seminarTalentRegistrationID"`

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
