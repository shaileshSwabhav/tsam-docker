package list

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// User contains both salesperson and admin.
type User struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
}

// ValidateUser validate salesperson details
func (user *User) ValidateUser() error {
	if util.IsEmpty(user.FirstName) {
		return errors.NewValidationError("FirstName is required")
	}
	if util.IsEmpty(user.LastName) {
		return errors.NewValidationError("LastName is required")
	}
	return nil
}
