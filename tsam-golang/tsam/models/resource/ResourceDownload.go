package resource

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
)

// Download will store credential and resource id's
type Download struct {
	general.TenantBase
	ResourceID   uuid.UUID `json:"resourceID" gorm:"type:varchar(36)"`
	CredentialID uuid.UUID `json:"credentialID" gorm:"type:varchar(36)"`
	// Resource     Resource  `json:"resource" gorm:"foreignkey:ResourceID;association_autoupdate:false;association_autocreate:false"`
	// Credential   Resource  `json:"credential" gorm:"foreignkey:CredentialID;association_autoupdate:false;association_autocreate:false"`
}

// TableName defines table name of the struct.
func (*Download) TableName() string {
	return "resource_downloads"
}

// ValidateResourceDownload will check if all the compulsory fields are specified
func (download *Download) ValidateResourceDownload() error {

	if download.ResourceID == uuid.Nil {
		return errors.NewValidationError("Resource must be specified")
	}

	if download.CredentialID == uuid.Nil {
		return errors.NewValidationError("Credental must be specified")
	}

	return nil
}

// DownloadDTO will store credential and resource id's
type DownloadDTO struct {
	ID           uuid.UUID          `json:"id"`
	DeletedAt    *time.Time         `json:"-"`
	ResourceID   uuid.UUID          `json:"-"`
	CredentialID uuid.UUID          `json:"-"`
	Resource     Resource           `json:"resource" gorm:"foreignkey:ResourceID"`
	Credential   general.Credential `json:"credential" gorm:"foreignkey:CredentialID"`
}

// TableName defines table name of the struct.
func (*DownloadDTO) TableName() string {
	return "resource_downloads"
}

// DownloadCount will contain totalcount of downloads.
type DownloadCount struct {
	TotalCount uint `json:"totalCount"`
}
