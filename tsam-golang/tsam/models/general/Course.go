package general

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

//Course contain information about courses
type Course struct {
	TenantBase
	Code             string       `json:"code" gorm:"type:varchar(20)"`
	Name             string       `json:"name" gorm:"type:varchar(100)"`
	DurationInMonths int64        `json:"durationInMonths" gorm:"type:TINYINT"`
	TotalHours       float32      `json:"totalHours" gorm:"type:decimal(6,2)"`
	Eligibility      *Eligibility `json:"courseEligibility"`
	EligibilityID    *uuid.UUID   `json:"-" gorm:"type:varchar(36)"`
}

// ValidateCourse validate the course details
func (course *Course) ValidateCourse() error {
	// Coursecode can contain any characters
	if util.IsEmpty(course.Code) {
		return errors.NewValidationError("Course code must be specified")
	}
	// Course name can contain any characters in the name
	if util.IsEmpty(course.Name) {
		return errors.NewValidationError("Course name must be specified")
	}
	if course.DurationInMonths <= 0 {
		return errors.NewValidationError("Course duration must be specified")
	}
	if course.TotalHours <= 0 {
		return errors.NewValidationError("Course total hours must be specified")
	}
	// if util.IsUUIDValid(course.EligibilityID) {
	// 	return errors.NewValidationError("EligibilityID must be specified")
	// }
	if course.Eligibility != nil {
		if err := course.Eligibility.ValidateEligibility(); err != nil {
			return err
		}
	}
	return nil
}
