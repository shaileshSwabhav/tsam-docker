package resource

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Resource will contain resource details.
type Resource struct {
	general.TenantBase
	ResourceName       string     `json:"resourceName" gorm:"type:varchar(50)"`
	IsExistingResource bool       `json:"isExistingResource" gorm:"-"`
	ResourceType       string     `json:"resourceType" gorm:"type:varchar(50)"`
	ResourceSubType    *string    `json:"resourceSubType" gorm:"type:varchar(50)"`
	FileType           string     `json:"fileType" gorm:"type:varchar(50)"`
	ResourceURL        string     `json:"resourceURL" gorm:"type:varchar(255)"`
	PreviewURL         string     `json:"previewURL" gorm:"type:varchar(255)"`
	Description        *string    `json:"description" gorm:"type:varchar(250)"`
	TechnologyID       *uuid.UUID `json:"technologyID" gorm:"type:varchar(36)"`

	// Book
	IsBook      *bool   `json:"isBook" gorm:"type:tinyint(1)"`
	Author      string  `json:"author" gorm:"type:varchar(100)"`
	Publication *string `json:"publication" gorm:"type:varchar(100)"`
	// Book   *Book `json:"book" gorm:"foreignkey:ID;association_autocreate:true;association_autoupdate:false;"`
}

// Book fields if resources.
type Book struct {
	general.TenantBase
	Author      string  `json:"author" gorm:"type:varchar(100)"`
	Publication *string `json:"publication" gorm:"type:varchar(100)"`
}

// ValidateResource validates resource fields.
func (resource *Resource) ValidateResource() error {

	if !resource.IsExistingResource {
		if util.IsEmpty(resource.ResourceName) {
			return errors.NewValidationError("Resource Name must be specified")
		}

		if util.IsEmpty(resource.ResourceURL) {
			return errors.NewValidationError("Resource must be specified")
		}
	} else {
		if !util.IsUUIDValid(resource.ID) {
			return errors.NewValidationError("Resource ID must be specified")
		}
	}
	if util.IsEmpty(resource.ResourceType) {
		return errors.NewValidationError("Resource Type must be specified")
	}
	if util.IsEmpty(resource.FileType) {
		return errors.NewValidationError("File Type must be specified")
	}
	// if util.IsEmpty(resource.Description) {
	// 	return errors.NewValidationError("Resource Description must be specified")
	// }
	return nil
}

// ResourceDTO will contain resource details.
type ResourceDTO struct {
	ID              uuid.UUID           `json:"id"`
	DeletedAt       *time.Time          `json:"-"`
	ResourceName    string              `json:"resourceName"`
	ResourceType    string              `json:"resourceType"`
	ResourceSubType *string             `json:"resourceSubType"`
	FileType        string              `json:"fileType"`
	ResourceURL     string              `json:"resourceURL"`
	PreviewURL      string              `json:"previewURL"`
	Description     *string             `json:"description"`
	TechnologyID    *uuid.UUID          `json:"-" gorm:"type:varchar(36)"`
	Technology      *general.Technology `json:"technology" gorm:"foreignkey:TechnologyID"`

	TotalDownload uint     `json:"totalDownload"`
	TotalLike     uint     `json:"totalLike"`
	ResourceLike  *LikeDTO `json:"resourceLike" gorm:"foreignkey:ResourceID"`

	// Book
	IsBook      *bool   `json:"isBook" gorm:"type:tinyint(1)"`
	Author      string  `json:"author" gorm:"type:varchar(100)"`
	Publication *string `json:"publication" gorm:"type:varchar(100)"`
	// Book   *Book `json:"book" gorm:"foreignkey:ID"`
}

// TableName defines table name of the struct.
func (*ResourceDTO) TableName() string {
	return "resources"
}

// Count will return count of resources available
type Count struct {
	FileType   string `json:"fileType"`
	TotalCount uint   `json:"totalCount"`
}
