package general

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

type AndroidUser struct {
	TenantBase
	FirstName string `json:"firstName" example:"Ravi" gorm:"type:varchar(50)"`
	LastName  string `json:"lastName" example:"Sharma" gorm:"type:varchar(50)"`
	Email     string `json:"email" example:"abc@gmail.com" gorm:"type:varchar(100)"`
	Password  string `json:"password" example:"7huw2J,X" gorm:"type:varchar(255)"`
	Contact   string `json:"contact" example:"9700795509" gorm:"type:varchar(15)"`
}

func (cred *AndroidUser) Validate() error {

	// First Name.
	if util.IsEmpty(cred.FirstName) || !util.ValidateString(cred.FirstName) {
		return errors.NewValidationError("First name must be specified and should contain only alphabets")
	}

	// First name maximum characters.
	if len(cred.FirstName) > 50 {
		return errors.NewValidationError("First name must can have maximum 50 characters")
	}

	// Last Name.
	if util.IsEmpty(cred.LastName) || !util.ValidateString(cred.FirstName) {
		return errors.NewValidationError("Last name should contain only alphabets")
	}

	// Last name maximum characters.
	if len(cred.LastName) > 50 {
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

	// Tenant ID.
	if cred.TenantID == uuid.Nil {
		return errors.NewValidationError("Tenant ID must be specified")
	}

	return nil
}
