package college

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Branch defines the college branch structure and fields in it.
type Branch struct {
	general.TenantBase
	general.Address     `json:"address,omitempty"`
	BranchName          string              `json:"branchName,omitempty" gorm:"type:varchar(150)"`
	CollegeID           uuid.UUID           `json:"collegeID,omitempty" gorm:"type:varchar(36)"`
	Code                string              `json:"code,omitempty" gorm:"type:varchar(10);not null"`
	AllIndiaRanking     *uint32             `json:"allIndiaRanking,omitempty" gorm:"type:int(6)"`
	TPOName             *string             `json:"tpoName,omitempty" gorm:"type:varchar(100)"`
	TPOContact          *string             `json:"tpoContact,omitempty" gorm:"type:varchar(15)"`
	TPOAlternateContact *string             `json:"tpoAlternateContact,omitempty" gorm:"type:varchar(15)"`
	TPOEmail            *string             `json:"tpoEmail,omitempty" gorm:"type:varchar(100);"`
	CollegeRating       *uint8              `json:"collegeRating,omitempty" gorm:"type:varchar(2)"`
	Email               *string             `json:"email,omitempty" gorm:"type:varchar(100);"`
	SalesPersonID       *uuid.UUID          `json:"-" gorm:"type:varchar(36)"`
	SalesPerson         *SalesPerson        `json:"salesPerson,omitempty" gorm:"foreignkey:SalesPersonID"`
	UniversityID        uuid.UUID           `json:"-" gorm:"type:varchar(36)"`
	University          *general.University `json:"university,omitempty" gorm:"foreignkey:UniversityID"`
}

// TableName will name the table of branch model as "college_branches"
func (*Branch) TableName() string {
	return "college_branches"
}

// ValidateCollegeBranch validates all fields of the college's branch
func (branch *Branch) ValidateCollegeBranch() error {

	if isEmpty := util.IsEmpty(branch.BranchName); isEmpty {
		return errors.NewValidationError("College's branch name must be specified")
	}

	// Just checking for nil as it would be tough to check for name or id as it is a pointer
	// Can check in another condition
	if branch.University == nil {
		return errors.NewValidationError("University must be specified")
	}

	if err := branch.MandatoryValidation(); err != nil {
		return err
	}
	return nil
}
