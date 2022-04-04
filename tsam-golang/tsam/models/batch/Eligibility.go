package batch

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
)

// Eligibility Contain Technology, StudentRating, Experience and passoutyear.
type Eligibility struct {
	general.TenantBase
	Technologies  []*general.Technology `json:"technologies" gorm:"many2many:eligibilities_technologies;association_autoupdate:false"`
	StudentRating *string               `json:"studentRating" gorm:"type:varchar(10)"`
	IsExperience  *bool                 `json:"experience" gorm:"type:TINYINT(1)"`
	AcademicYear  *string               `json:"academicYear" gorm:"type:varchar(200)"`
}

// TableName gives table name
func (*Eligibility) TableName() string {
	return "batch_eligibilities"
}

// ValidateEligibility will validate all the fields of eligibility
func (eligibility *Eligibility) ValidateEligibility() error {

	if eligibility.Technologies == nil {
		return errors.NewValidationError("Technology must be specified.")
	}

	// if util.IsEmpty(*eligibility.StudentRating) {
	// 	return errors.NewValidationError("Student Rating must be specified")
	// }

	// if util.IsEmpty(*eligibility.AcademicYear) {
	// 	return errors.NewValidationError("Academic Year must be specified")
	// }

	// for _, technology := range eligibility.Technologies {
	// 	if err := technology.Validate(); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
