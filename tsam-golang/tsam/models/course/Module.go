package course

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/util"
)

// Module will have fields related to module table.
type Module struct {
	general.TenantBase
	ModuleName string               `json:"moduleName" gorm:"type:varchar(100)"`
	Logo       *string              `json:"logo" gorm:"type:varchar(1000)"`
	Resources  []*resource.Resource `json:"resources" gorm:"many2many:modules_resources;association_autoupdate:false"`

	// ComplexityLevel int `json:"complexityLevel" gorm:"type:int"`
	// technology map
	// minutes
}

// Validate will verfiy if compulsory fields of module struct are specified.
func (module *Module) Validate() error {

	if util.IsEmpty(module.ModuleName) {
		return errors.NewValidationError("Module name must be specified")
	}

	return nil
}

// ModuleDTO will get all the fields of modules table.
type ModuleDTO struct {
	ID                        uuid.UUID            `json:"id"`
	DeletedAt                 *time.Time           `json:"-"`
	ModuleName                string               `json:"moduleName"`
	Logo                      *string              `json:"logo"`
	TotalModuleTopics         int                  `json:"totalTopics,omitempty"`
	TotalSubTopics            int                  `json:"totalSubTopics,omitempty"`
	TotalProgrammingQuestions int                  `json:"totalProgrammingQuestions,omitempty"`
	ModuleTopics              []*ModuleTopicDTO    `json:"moduleTopics" gorm:"foreignkey:ModuleID"`
	Resources                 []*resource.Resource `json:"resources" gorm:"many2many:modules_resources;association_jointable_foreignkey:resource_id;jointable_foreignkey:module_id"`
}

// TableName overrides name of the table
func (*ModuleDTO) TableName() string {
	return "modules"
}
