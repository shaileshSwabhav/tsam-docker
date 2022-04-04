package college

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// College defines College structure
type College struct {
	general.TenantBase
	CollegeName     string    `json:"collegeName" gorm:"type:varchar(100);not null"`
	Code            string    `json:"code" gorm:"type:varchar(10);not null"`
	ChairmanName    *string   `json:"chairmanName,omitempty" gorm:"type:varchar(100)"`
	ChairmanContact *string   `json:"chairmanContact,omitempty" gorm:"type:varchar(15)"`
	CollegeBranches []*Branch `json:"collegeBranches,omitempty" gorm:"foreignkey:CollegeID"`
}

// ValidateCollege validates the college. Returns error if College has invalid fields
func (college *College) ValidateCollege() error {

	if util.IsEmpty(college.CollegeName) {
		return errors.NewValidationError("College Name must be specified")
	}

	if college.CollegeBranches == nil {
		return errors.NewValidationError("College must have atleast 1 branch")
	}

	if !util.IsUUIDValid(college.ID) {
		for _, branch := range college.CollegeBranches {
			if err := branch.ValidateCollegeBranch(); err != nil {
				return err
			}
		}
	}
	return nil
}
