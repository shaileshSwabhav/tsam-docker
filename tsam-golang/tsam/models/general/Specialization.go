package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Specialization contains specialization details which are enough for adding and updating specialization.
type Specialization struct {
	TenantBase
	BranchName string    `json:"branchName" gorm:"type:varchar(200)"`
	DegreeID   uuid.UUID `json:"degreeID" gorm:"type:varchar(36)"`
}

// Validate compusalry filds of specialization.
func (specialization *Specialization) Validate() error {

	// Check if branch name exists or not.
	if util.IsEmpty(specialization.BranchName) {
		return errors.NewValidationError("Branch name must be specified")
	}

	// Branch name maximum characters.
	if len(specialization.BranchName) > 200 {
		return errors.NewValidationError("Branch name can have maximum 200 characters")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// SpecializationDTO contains the complete information of specialization which is needed to display.
type SpecializationDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	BranchName string    `json:"branchName"`
	DegreeID   uuid.UUID `json:"-"`
	Degree     Degree    `json:"degree" gorm:"foreignkey:DegreeID"`
}

// TableName will name the table of Experience model as "specializations"
func (*SpecializationDTO) TableName() string {
	return "specializations"
}
