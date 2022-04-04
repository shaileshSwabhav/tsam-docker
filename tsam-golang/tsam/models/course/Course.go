package course

import (
	"time"

	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"

	uuid "github.com/satori/go.uuid"
)

// Course struct will contain all details for a particular course
type Course struct {
	general.TenantBase
	Code             string                `json:"code" gorm:"type:varchar(10);not null"`
	Name             string                `json:"name" gorm:"type:varchar(100)"`
	CourseType       string                `json:"courseType" gorm:"type:varchar(20)"`
	CourseLevel      string                `json:"courseLevel" gorm:"type:varchar(20)"`
	Description      *string               `json:"description" gorm:"type:varchar(2000)"`
	PreRequisites    *string               `json:"preRequisites" gorm:"type:varchar(2000)"`
	Price            *float64              `json:"price" gorm:"type:double"`
	DurationInMonths int64                 `json:"durationInMonths" gorm:"type:TINYINT"`
	TotalHours       float32               `json:"totalHours" gorm:"type:decimal(6,2)"`
	Technologies     []*general.Technology `json:"technologies" gorm:"many2many:courses_technologies;association_autocreate:false;association_autoupdate:false;"`
	Eligibility      *general.Eligibility  `json:"eligibility" gorm:"foreignkey:EligibilityID"`
	EligibilityID    *uuid.UUID            `json:"-" gorm:"type:varchar(36)"`
	Brochure         *string               `json:"brochure" gorm:"type:varchar(200)"`
	Logo             *string               `json:"logo" gorm:"type:varchar(200)"`

	// TotalSessions    uint                  `json:"totalSessions" gorm:"-"`
	// Sessions         []*Session            `json:"sessions" gorm:"-"`
}

// Validate validates all the course fields
func (course *Course) Validate() error {
	if util.IsEmpty(course.Name) {
		return errors.NewValidationError("Course Name must be specified")
	}

	if util.IsEmpty(course.CourseType) {
		return errors.NewValidationError("Course Type must be specified")
	}

	if util.IsEmpty(course.CourseLevel) {
		return errors.NewValidationError("Course Level must be specified")
	}

	// if util.IsEmpty(course.Description) {
	// 	return errors.NewValidationError("Description must be specified")
	// }

	if course.TotalHours <= 0 {
		return errors.NewValidationError("Total Hours must be specified")
	}

	if course.Technologies == nil {
		return errors.NewValidationError("Atleat 1 technology must be specified")
	}

	if course.DurationInMonths <= 0 {
		return errors.NewValidationError("Course Duration must be specified")
	}

	// if course.Price <= 0 {
	// 	return errors.NewValidationError("Course Price must be specified")
	// }

	// for _, tech := range course.Technologies {
	// 	if err := tech.Validate(); err != nil {
	// 		return err
	// 	}
	// }

	// if course.Eligibility != nil {
	// 	err := course.Eligibility.ValidateEligibility()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// for _, session := range course.Sessions {
	// 	if err := session.ValidateSession(); err != nil {
	// 		return err
	// 	}
	// 	// course.Sessions[index] = session
	// }

	return nil
}

// CouseDTO will get all the details of a course
type CourseDTO struct {
	ID               uuid.UUID             `json:"id"`
	DeletedAt        *time.Time            `json:"-"`
	Code             string                `json:"code"`
	Name             string                `json:"name"`
	CourseType       string                `json:"courseType"`
	CourseLevel      string                `json:"courseLevel"`
	Description      *string               `json:"description"`
	PreRequisites    *string               `json:"preRequisites"`
	Price            float64               `json:"price"`
	DurationInMonths int64                 `json:"durationInMonths"`
	TotalHours       float32               `json:"totalHours"`
	Technologies     []*general.Technology `json:"technologies" gorm:"many2many:courses_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:course_id"`
	EligibilityID    *uuid.UUID            `json:"-"`
	Eligibility      *general.Eligibility  `json:"eligibility"`
	TotalSessions    uint                  `json:"totalSessions"`
	Brochure         *string               `json:"brochure"`
	Logo             *string               `json:"logo"`

	// extra rules
	TotalModules   int `json:"totalModules"`
	TotalTopics    int `json:"totalTopics"`
	TotalConcepts  int `json:"totalConcepts"`
	TotalQuestions int `json:"totalQuestions"`
}

// TableName defines table name of the struct.
func (*CourseDTO) TableName() string {
	return "courses"
}
