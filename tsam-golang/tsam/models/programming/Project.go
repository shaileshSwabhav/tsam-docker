package programming

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingProject will contain details of programming-project
type ProgrammingProject struct {
	general.TenantBase
	ProjectName     string               `json:"projectName" gorm:"type:varchar(100)"`
	ProjectType     string               `json:"projectType"`
	Description     string               `json:"description" gorm:"type:varchar(250)"`
	Code            string               `json:"code" gorm:"type:varchar(10);not null"`
	IsActive        *bool                `json:"isActive" gorm:"default:true"`
	ComplexityLevel uint8                `json:"complexityLevel" gorm:"type:int"`
	RequiredHours   uint                 `json:"requiredHours" gorm:"type:int"`
	SampleURL       *string              `json:"sampleUrl" gorm:"type:varchar(250)"`
	ResourceType    *string              `json:"resourceType" gorm:"type:varchar(50)"`
	Document        *string              `json:"document" gorm:"type:varchar(200)"`
	Score           *float32             `json:"score"`
	Technologies    []general.Technology `json:"technologies" gorm:"many2many:programming_projects_technologies;association_autocreate:false;association_autoupdate:false;"`
	Resources       []resource.Resource  `json:"resources" gorm:"many2many:programming_projects_resources;association_autocreate:false;association_autoupdate:false;"`
}

// Validate will validate all the compuslory fields.
func (programmingProject *ProgrammingProject) Validate() error {

	if util.IsEmpty(programmingProject.ProjectName) {
		return errors.NewValidationError("Project name must be specified")
	}
	if util.IsEmpty(programmingProject.Description) {
		return errors.NewValidationError("Project description must be specified")
	}
	if programmingProject.ComplexityLevel == 0 {
		return errors.NewValidationError("Complexity level must be specified")
	}
	if programmingProject.RequiredHours == 0 {
		return errors.NewValidationError("Required hours must be specified")
	}
	if programmingProject.Technologies == nil {
		return errors.NewValidationError("Technologies must be specified")
	}

	return nil
}

// ProgrammingProjectDTO will contain details of programming-project.
type ProgrammingProjectDTO struct {
	ID              uuid.UUID              `json:"id"`
	DeletedAt       *time.Time             `json:"-"`
	ProjectName     string                 `json:"projectName"`
	ProjectType     string                 `json:"projectType"`
	Description     string                 `json:"description"`
	Code            string                 `json:"code"`
	IsActive        *bool                  `json:"isActive"`
	ComplexityLevel uint8                  `json:"complexityLevel"`
	RequiredHours   uint                   `json:"requiredHours"`
	SampleURL       *string                `json:"sampleUrl"`
	ResourceType    *string                `json:"resourceType"`
	Score           *float32               `json:"score"`
	Document        *string                `json:"document" gorm:"type:varchar(200)"`
	Technologies    []general.Technology   `json:"technologies" gorm:"many2many:programming_projects_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:programming_project_id"`
	Resources       []resource.ResourceDTO `json:"resources" gorm:"many2many:programming_projects_resources;association_jointable_foreignkey:resource_id;jointable_foreignkey:programming_project_id"`
}

// TableName defines table name of the struct.
func (*ProgrammingProjectDTO) TableName() string {
	return "programming_projects"
}
