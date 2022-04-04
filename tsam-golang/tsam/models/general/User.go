package general

import (
	"regexp"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// User Contain ID, Name, Code.
type User struct {
	TenantBase
	Address       `json:"address"`
	Code          string    `json:"code" gorm:"varchar(10)"`
	FirstName     string    `json:"firstName" example:"Ravi" gorm:"type:varchar(50)"`
	LastName      string    `json:"lastName" example:"Sharma" gorm:"type:varchar(50)"`
	Email         string    `json:"email" example:"abc@gmail.com" gorm:"type:varchar(50)"`
	Contact       string    `json:"contact" example:"9700795509" gorm:"type:varchar(15)"`
	DateOfBirth   *string   `json:"dateOfBirth" example:"2000-02-30" gorm:"type:date"`
	DateOfJoining *string   `json:"dateOfJoining" example:"2019-02-23" gorm:"type:date"`
	RoleID        uuid.UUID `json:"-" gorm:"type:varchar(36)"`
	Role          Role      `json:"role" gorm:"association_autocreate:false;association_autoupdate:false"`
	IsActive      bool      `json:"isActive" gorm:"DEFAULT:true"`
	Resume        *string   `json:"resume" gorm:"type:varchar(200)"`
}

// Validate (ADD REGEX)
func (user *User) Validate() error {
	datePattern := regexp.MustCompile("/^\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])$/g")
	// accpets date in yyyy-mm-dd pattern
	if util.IsEmpty(user.FirstName) || !util.ValidateString(user.FirstName) {
		return errors.NewValidationError("User FirstName must be specified and must have characters only")
	}
	if util.IsEmpty(user.LastName) || !util.ValidateString(user.LastName) {
		return errors.NewValidationError("User LastName must be specified and must have characters only")
	}
	if util.IsEmpty(user.Email) || !util.ValidateEmail(user.Email) {
		return errors.NewValidationError("User Email must be specified and should be of the type abc@domain.com")
	}
	if util.IsEmpty(user.Contact) || !util.ValidateContact(user.Contact) {
		return errors.NewValidationError("User Contact must be specified and have 10 digits")
	}
	if !util.IsUUIDValid(user.Role.ID) {
		return errors.NewValidationError("RoleID must be specified")
	}
	if user.DateOfBirth != nil && datePattern.MatchString(*user.DateOfBirth) {
		return errors.NewValidationError("Date of birth should be in the form YYYY-MM-DD")
	}
	if user.DateOfJoining != nil && datePattern.MatchString(*user.DateOfJoining) {
		return errors.NewValidationError("Date of joining should be in the form YYYY-MM-DD")
	}
	return nil
}
