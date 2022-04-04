package general

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// CareerObjective contains the definition of career objective.
type CareerObjective struct {
	TenantBase
	Name    string                   `json:"name" gorm:"type:varchar(100)"`
	Courses []CareerObjectivesCourse `json:"courses"`
}

// CareerObjectivesCourse is career_objectives and courses map with order field.
type CareerObjectivesCourse struct {
	TenantBase
	CareerObjectiveID uuid.UUID `json:"careerObjectiveID" gorm:"type:varchar(36)"`
	CourseID          uuid.UUID `json:"courseID" gorm:"type:varchar(36)"`
	Order             uint      `json:"order" gorm:"type:SMALLINT(2)"`
	TechnicalAspect   string    `json:"technicalAspect" gorm:"type:varchar(500)"`
}

// Validate Validates compulsary fields of career objective.
func (careerObjective *CareerObjective) Validate() error {
	// Check if name is blank or not.
	if util.IsEmpty(careerObjective.Name) {
		return errors.NewValidationError("Name must be specified")
	}

	// First name should consist of only alphabets and space.
	if !util.ValidateStringWithSpace(careerObjective.Name) {
		return errors.NewValidationError("Name should consist of only alphabets")
	}

	// First name maximum characters.
	if len(careerObjective.Name) > 100 {
		return errors.NewValidationError("Name can have maximum 100 characters")
	}

	// Check if all career objective courses are having unique order field.
	courseMap := make(map[uint]uint)
	for _, course := range careerObjective.Courses {
		courseMap[course.Order]++
		if courseMap[course.Order] > 1 {
			return errors.NewValidationError("Same course order cannot be assigned to courses")
		}
	}

	// Validate all courses.
	for _, course := range careerObjective.Courses {
		err := course.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// Validate Validates compulsary fields of career objective course.
func (careerObjectiveCourse *CareerObjectivesCourse) Validate() error {
	// Check if order is 0 or below.
	if careerObjectiveCourse.Order <= 0 {
		return errors.NewValidationError("Order cannot be lesser than or equal to 0")
	}

	// Check if technical aspect is blank or not.
	if util.IsEmpty(careerObjectiveCourse.TechnicalAspect) {
		return errors.NewValidationError("Technical Aspect must be specified")
	}

	// Technical aspect maximum characters.
	if len(careerObjectiveCourse.TechnicalAspect) > 500 {
		return errors.NewValidationError("Technical Aspect can have maximum 500 characters")
	}

	// Check if course id is blank or not.
	if !util.IsUUIDValid(careerObjectiveCourse.CourseID) {
		return errors.NewValidationError("Course must be specified")
	}

	return nil
}
