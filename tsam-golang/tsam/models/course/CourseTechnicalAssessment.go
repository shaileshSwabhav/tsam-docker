package course

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

// CourseTechnicalAssessment contains technical assessment for faculty
type CourseTechnicalAssessment struct {
	general.TenantBase
	CourseID  uuid.UUID `json:"courseID" gorm:"type:varchar(36)"`
	FacultyID uuid.UUID `json:"facultyID" gorm:"type:varchar(36)"`
	Rating    *uint64   `json:"rating" gorm:"type:int"`
	// Course    list.Course     `json:"course" gorm:"foreignkey:CourseID;association_autocreate:false;association_autoupdate:false"`
	// Faculty   faculty.Faculty `json:"faculty" gorm:"foreignkey:FacultyID;association_autocreate:false;association_autoupdate:false"`
}

func (assessment *CourseTechnicalAssessment) Validate() error {

	if !util.IsUUIDValid(assessment.CourseID) {
		return errors.NewValidationError("Course ID must be specified.")
	}

	// if !util.IsUUIDValid(assessment.Faculty.ID) {
	// 	return errors.NewValidationError("Faculty ID must be specified.")
	// }

	if assessment.Rating != nil {
		if *assessment.Rating <= 0 && *assessment.Rating > 10 {
			return errors.NewValidationError("Rating must be specified and should be greater than 0 and less than 10.")
		}
	}

	return nil
}

// CourseTechnicalAssessmentDTO contains technical assessment for faculty
type CourseTechnicalAssessmentDTO struct {
	ID        uuid.UUID        `json:"id"`
	DeletedAt *time.Time       `json:"-"`
	CourseID  uuid.UUID        `json:"-"`
	FacultyID uuid.UUID        `json:"-"`
	Course    list.Course      `json:"course" gorm:"foreignkey:CourseID"`
	Rating    *uint64          `json:"rating"`
	Faculty   *faculty.Faculty `json:"faculty" gorm:"foreignkey:FacultyID"`
}

func (*CourseTechnicalAssessmentDTO) TableName() string {
	return "course_technical_assessments"
}
