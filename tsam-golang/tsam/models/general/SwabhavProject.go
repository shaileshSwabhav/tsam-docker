package general

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// SwabhavProject will contain names of all projects.
type SwabhavProject struct {
	TenantBase
	Name        string           `json:"name" gorm:"type:varchar(100)"`
	ProjectID   *uuid.UUID       `json:"projectID" gorm:"type:varchar(36)"`
	SubProjects []SwabhavProject `json:"subProjects" gorm:"foreignkey:ProjectID"`
}

// Validate validates fields of project.
func (project *SwabhavProject) Validate() error {
	if util.IsEmpty(project.Name) {
		return errors.NewValidationError("Project name must be specified")
	}

	// Check if sub pojects are having same name or not.
	projectMap := make(map[string]uint)
	for _, subProject := range project.SubProjects {
		projectMap[subProject.Name]++
		if projectMap[subProject.Name] > 1 {
			return errors.NewValidationError("Sub Projects projects cannot have same name")
		}
	}

	// Check if parent project id is valid or not.
	if project.ProjectID != nil && !util.IsUUIDValid(*project.ProjectID) {
		return errors.NewValidationError("Parent Project ID is invalid")
	}
	return nil
}
