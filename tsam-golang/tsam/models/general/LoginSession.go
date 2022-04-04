package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// LoginSession will contain login details for every time any user logs in.
type LoginSession struct {
	TenantBase
	StartTime    time.Time  `json:"startTime,omitempty"`
	EndTime      *time.Time `json:"endTime,omitempty"`
	CredentialID uuid.UUID  `json:"credentialID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	// RoleID       uuid.UUID  `json:"roleID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
}

// ValidateLoginSession validates login session fields.
func (loginSession *LoginSession) ValidateLoginSession() error {
	if util.IsUUIDValid(loginSession.CredentialID) {
		return errors.NewValidationError("LoginID must be specified")
	}
	// if util.IsUUIDValid(loginSession.RoleID) {
	// 	return errors.NewValidationError("RoleID must be specified")
	// }
	return nil
}
