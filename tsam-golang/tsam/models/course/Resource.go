package course

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// ModuleResource contain ResourceID and CourseSessionID.
type ModuleResource struct {
	ModuleID   uuid.UUID `json:"moduleID" gorm:"type:varchar(36)"`
	ResourceID uuid.UUID `json:"resourceID" gorm:"type:varchar(36)"`
}

// TableName changes the default table name
func (*ModuleResource) TableName() string {
	return "modules_resources"
}

// Validate module resource.
func (resource *ModuleResource) Validate() error {
	if !util.IsUUIDValid(resource.ModuleID) {
		return errors.NewValidationError("Resource ID must be specified")
	}
	if !util.IsUUIDValid(resource.ResourceID) {
		return errors.NewValidationError("Resource ID must be specified")
	}
	return nil

}
