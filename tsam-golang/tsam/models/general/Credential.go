package general

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Credential contains information about all the credentials in the system.
type Credential struct {
	TenantBase
	FirstName     string     `json:"firstName" example:"Ravi" gorm:"type:varchar(50)"`
	LastName      *string    `json:"lastName" example:"Sharma" gorm:"type:varchar(50)"`
	Email         string     `json:"email" example:"abc@gmail.com" gorm:"type:varchar(100)"`
	Password      string     `json:"password" example:"7huw2J,X" gorm:"type:varchar(255)"`
	Contact       string     `json:"contact" example:"9700795509" gorm:"type:varchar(15)"`
	IsActive      *bool      `json:"isActive" example:"1" gorm:"DEFAULT:true"`
	Role          Role       `json:"role" gorm:"association_autocreate:false;association_autoupdate:false"`
	RoleID        uuid.UUID  `json:"roleID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	TalentID      *uuid.UUID `json:"talentID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	FacultyID     *uuid.UUID `json:"facultyID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	CollegeID     *uuid.UUID `json:"collegeID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	CompanyID     *uuid.UUID `json:"companyID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	EmployeeID    *uuid.UUID `json:"employeeID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	SalesPersonID *uuid.UUID `json:"salesPersonID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	UserID        *uuid.UUID `json:"userID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	// DepartmentID  *uuid.UUID `json:"-" gorm:"-"`
}

// ValidateCredentials validates all the compulsory credential fields.
func (cred *Credential) ValidateCredentials() error {

	// First Name.
	if util.IsEmpty(cred.FirstName) || !util.ValidateString(cred.FirstName) {
		return errors.NewValidationError("First name must be specified and should contain only alphabets")
	}

	// First name maximum characters.
	if len(cred.FirstName) > 50 {
		return errors.NewValidationError("First name must can have maximum 50 characters")
	}

	// Last Name.
	if cred.LastName != nil && !util.ValidateString(cred.FirstName) {
		return errors.NewValidationError("Last name should contain only alphabets")
	}

	// Last name maximum characters.
	if cred.LastName != nil && len(*cred.LastName) > 50 {
		return errors.NewValidationError("Last name can have maximum 50 characters")
	}

	// Email.
	if util.IsEmpty(cred.Email) || !util.ValidateEmail(cred.Email) {
		return errors.NewValidationError("Email must be specified and should be of type example@domain.com")
	}

	// Email maximum characters.
	if len(cred.Email) > 100 {
		return errors.NewValidationError("Email can have maximum 100 characters")
	}

	// Password.
	if util.IsEmpty(cred.Password) {
		return errors.NewValidationError("Password must be specified")
	}

	// Contact.
	if util.IsEmpty(cred.Contact) || !util.ValidateContact(cred.Contact) {
		return errors.NewValidationError("Contact must be specified and should be 10 digits")
	}

	// Contact maximum characters.
	if len(cred.Contact) > 15 {
		return errors.NewValidationError("Contact can have maximum 15 characters")
	}

	// Role ID.
	if cred.RoleID == uuid.Nil {
		return errors.NewValidationError("Role ID must be specified")
	}

	// Tenant ID.
	if cred.TenantID == uuid.Nil {
		return errors.NewValidationError("Tenant ID must be specified")
	}

	return nil
}

// PasswordChange contains information for verifying and chnaging the password.
type PasswordChange struct {
	TenantID     uuid.UUID `json:"tenantID"`
	CredentialID uuid.UUID `json:"credentialID"`
	Email        string    `json:"email"`
	RoleID       uuid.UUID `json:"roleID"`
	Password     string    `json:"password"`
}

// Validate validates all the compulsory of PasswordChange model.
func (passwordChange *PasswordChange) Validate() error {
	// Email.
	if util.IsEmpty(passwordChange.Email) {
		return errors.NewValidationError("Email must be specified")
	}

	// Role ID.
	if !util.IsUUIDValid(passwordChange.RoleID) {
		return errors.NewValidationError("Role ID must be specified")
	}

	// Password.
	if util.IsEmpty(passwordChange.Password) {
		return errors.NewValidationError("Password must be specified")
	}

	return nil
}
