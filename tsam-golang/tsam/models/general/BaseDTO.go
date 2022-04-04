package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// BaseDTO contains master fields for DTO specifically.
// Should only be used for reading operations.
type BaseDTO struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	TenantID  uuid.UUID  `json:"tenantID"`
	DeletedAt *time.Time `json:"-"`
}
