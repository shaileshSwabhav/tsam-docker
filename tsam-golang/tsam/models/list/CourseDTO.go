package list

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Course is DTO used for listing of all the courses
type Course struct {
	ID         uuid.UUID `json:"id" gorm:"type:varchar(36)"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	CourseType string    `json:"courseType"`
}

// ValidateCourse validate the course details
func (course *Course) ValidateCourse() error {
	if !util.IsUUIDValid(course.ID) {
		return errors.NewValidationError("Invalid course id")
	}
	return nil
}

// TableName will create the table for Course model with name courses.
func (*Course) TableName() string {
	return "courses"
}
