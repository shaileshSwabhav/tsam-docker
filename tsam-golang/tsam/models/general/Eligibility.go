package general

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Eligibility Contain Technology, StudentRating, Experience and passoutyear.
type Eligibility struct {
	TenantBase
	Technologies  []*Technology `json:"technologies" gorm:"many2many:eligibilities_technologies;association_autocreate:false;association_autoupdate:false;"`
	StudentRating string        `json:"studentRating" gorm:"type:varchar(10)"`
	IsExperience  bool          `json:"experience" gorm:"type:TINYINT(1)"`
	AcademicYear  string        `json:"academicYear" gorm:"type:varchar(200)"`
}

// ValidateEligibility will validate all the fields of eligibility
func (eligibility *Eligibility) ValidateEligibility() error {
	if util.IsEmpty(eligibility.StudentRating) {
		return errors.NewValidationError("Student Rating must be specified")
	}
	if util.IsEmpty(eligibility.AcademicYear) {
		return errors.NewValidationError("Academic Year must be specified")
	}
	for _, technology := range eligibility.Technologies {
		if err := technology.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// IsEligibilityIDValid Check weather eligibility ID Valid Or Not.
func (eligibility *Eligibility) IsEligibilityIDValid() bool {
	// if eligibility.ID != uuid.Nil {
	// 	return true
	// }
	return eligibility.ID != uuid.Nil
}

// MakeEligibiltyValidOrError Either Make Eligibilty Valid Or Return Error
func (eligibility *Eligibility) MakeEligibiltyValidOrError() error {
	if !eligibility.IsEligibilityIDValid() {
		eligibility.AddEligibilityID()
	}

	if len(eligibility.Technologies) == 0 {
		return errors.NewValidationError("Atleast One Technology must be specified in Eligibility")
	}

	for _, tech := range eligibility.Technologies {
		if err := tech.Validate(); err != nil {
			return err
		}
	}

	if util.IsEmpty(eligibility.StudentRating) {
		return errors.NewValidationError("Student Rating must be specified")
	}

	if util.IsEmpty(eligibility.AcademicYear) {
		return errors.NewValidationError("Academic Year must be specified")
	}

	return nil
}

// AddEligibilityID Add ID To Eligibility.
func (eligibility *Eligibility) AddEligibilityID() {
	eligibility.ID = util.GenerateUUID()
}
