package list

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
)

// Credential contains details from credentials table.
type Credential struct {
	ID        uuid.UUID    `json:"id"`
	DeletedAt *time.Time   `json:"-"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	Role      general.Role `json:"role"`
	RoleID    uuid.UUID    `json:"-"`
}

// TableName will create the table for Credential model with name credentials.
func (*Credential) TableName() string {
	return "credentials"
}
