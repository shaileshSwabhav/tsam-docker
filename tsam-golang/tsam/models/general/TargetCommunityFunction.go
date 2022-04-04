package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// TargetCommunityFunction contains the target community function details which are enough for adding and updating target community function.
type TargetCommunityFunction struct {
	TenantBase
	FunctionName string    `json:"functionName" gorm:"type:varchar(100)"`
	DepartmentID uuid.UUID `json:"departmentID" gorm:"type:varchar(36)"`
}

// Validate will check if all fields are valid in TargetCommunityFunction.
func (targetCommunityFunction *TargetCommunityFunction) Validate() error {

	// Function name must be spcified.
	if util.IsEmpty(targetCommunityFunction.FunctionName) {
		return errors.NewValidationError("Target community function name must be specified")
	}

	// Function name maximum characters.
	if len(targetCommunityFunction.FunctionName) > 100 {
		return errors.NewValidationError("Target community function name can have maximum 100 characters")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// TargetCommunityFunctionDTO contains the complete information of target community function which is needed to display.
type TargetCommunityFunctionDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	FunctionName string        `json:"functionName"`
	Department   DepartmentDTO `json:"department" gorm:"foreignkey:DepartmentID"`
	DepartmentID uuid.UUID     `json:"-"`
}

// TableName will name the table of Experience model as "target_community_functions"
func (*TargetCommunityFunctionDTO) TableName() string {
	return "target_community_functions"
}
