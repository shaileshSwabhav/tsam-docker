package course

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

// CourseModule contains all of fields of course_modules table.
type CourseModule struct {
	general.TenantBase
	CourseID uuid.UUID `json:"courseID" gorm:"type:varchar(36)"`
	ModuleID uuid.UUID `json:"moduleID" gorm:"type:varchar(36)"`
	Order    uint      `json:"order" gorm:"type:int"`
	IsActive *bool     `json:"isActive" gorm:"DEFAULT:true"`
}

// Validate will check if all fields of course_modules struct.
func (module *CourseModule) Validate() error {

	if !util.IsUUIDValid(module.CourseID) {
		return errors.NewValidationError("CourseID must be specified")
	}

	if !util.IsUUIDValid(module.ModuleID) {
		return errors.NewValidationError("ModuleID must be specified")
	}

	return nil
}

// CourseModuleDTO will get all the details of a course_modules.
type CourseModuleDTO struct {
	ID           uuid.UUID         `json:"id"`
	DeletedAt    *time.Time        `json:"-"`
	CourseID     uuid.UUID         `json:"-"`
	ModuleID     uuid.UUID         `json:"-"`
	Course       list.Course       `json:"course"`
	Module       ModuleDTO         `json:"module"`
	Order        uint              `json:"order"`
	IsActive     *bool             `json:"isActive"`
	ModuleTopics []*ModuleTopicDTO `json:"moduleTopics,omitempty" gorm:"foreignkey:CourseID"`
}

// TableName overrides name of the table
func (*CourseModuleDTO) TableName() string {
	return "course_modules"
}
