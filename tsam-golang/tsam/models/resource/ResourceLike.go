package resource

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
)

// Like will store like information
type Like struct {
	general.TenantBase
	ResourceID   uuid.UUID `json:"resourceID" gorm:"type:varchar(36)"`
	CredentialID uuid.UUID `json:"credentialID" gorm:"type:varchar(36)"`
	IsLiked      *bool     `json:"isLiked"`
	// Resource     Resource  `json:"resource" gorm:"foreignkey:ResourceID;association_autoupdate:false;association_autocreate:false"`
	// Credential   Resource  `json:"credential" gorm:"foreignkey:CredentialID;association_autoupdate:false;association_autocreate:false"`
}

// TableName defines table name of the struct.
func (*Like) TableName() string {
	return "resource_likes"
}

//ValidateResourceLike will validate all fields of resource likes
func (like *Like) ValidateResourceLike() error {

	if like.ResourceID == uuid.Nil {
		return errors.NewValidationError("Resource must be specified")
	}

	if like.CredentialID == uuid.Nil {
		return errors.NewValidationError("Credental must be specified")
	}

	return nil
}

// LikeDTO will store credential and resource id's
type LikeDTO struct {
	ID           uuid.UUID          `json:"id"`
	DeletedAt    *time.Time         `json:"-"`
	ResourceID   uuid.UUID          `json:"resourceID"`
	CredentialID uuid.UUID          `json:"credentialID"`
	IsLiked      *bool              `json:"isLiked"`
	Resource     Resource           `json:"resource" gorm:"foreignkey:ResourceID"`
	Credential   general.Credential `json:"credential" gorm:"foreignkey:CredentialID"`
	TotalCount   uint               `json:"totalCount"`
}

// TableName defines table name of the struct.
func (*LikeDTO) TableName() string {
	return "resource_likes"
}

// type LikeCount struct {
// 	TotalCount uint `json:"totalCount"`
// }
