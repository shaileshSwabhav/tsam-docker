package blog

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// BlogView contains add update fields required for blog view.
type BlogView  struct {
	general.TenantBase
	BlogID uuid.UUID `json:"blogID" gorm:"type:varchar(36)"`
	ViewerID uuid.UUID `json:"viewerID" gorm:"type:varchar(36)"`
}

// Validate validates compulsary fields of BlogView.
func (viewer *BlogView) Validate() error {

	// Blog ID.
	if !util.IsUUIDValid(viewer.BlogID) {
		return errors.NewValidationError("Blog ID must be a proper uuid")
	}
	
	// Viewer ID.
	if !util.IsUUIDValid(viewer.ViewerID) {
		return errors.NewValidationError("Viewer ID must be a proper uuid")
	}

	return nil
}