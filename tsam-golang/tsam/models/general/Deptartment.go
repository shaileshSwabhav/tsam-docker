package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Department will contain names of all departments.
type Department struct {
	TenantBase
	Name   string    `json:"name" gorm:"type:varchar(50)"`
	RoleID uuid.UUID `json:"roleID" gorm:"type:varchar(36)"`
}

// Validate validates fields of department.
func (dept *Department) Validate() error {

	// Department name must be specified.
	if util.IsEmpty(dept.Name) {
		return errors.NewValidationError("Department name must be specified")
	}

	// Department name maximum characters.
	if len(dept.Name) > 50 {
		return errors.NewValidationError("Department name cam have maximum 50 characters")
	}

	// Role ID must be secified.
	if !util.IsUUIDValid(dept.RoleID) {
		return errors.NewValidationError("Role ID must be secified")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// Department will contain names of all departments.
type DepartmentDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	Name   string    `json:"name"`
	RoleID uuid.UUID `json:"-"`
	Role   Role      `json:"role" gorm:"foreignkey:RoleID"`
}

// TableName will name the table of Experience model as "departments".
func (*DepartmentDTO) TableName() string {
	return "departments"
}
