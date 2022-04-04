package faculty

import (
	"regexp"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Academic contains single FacultyAcadmeic details.
type Academic struct {
	general.TenantBase
	Specialization   *general.Specialization `json:"specialization" gorm:"foreignkey:SpecializationID;association_autocreate:false;association_autoupdate:false;"`
	Degree           *general.Degree         `json:"degree" gorm:"foreignkey:DegreeID;association_autocreate:false;association_autoupdate:false;"`
	College          string                  `json:"college" gorm:"type:varchar(200)"`
	Percentage       float32                 `json:"percentage" gorm:"type:decimal(4,2)"`
	Passout          int                     `json:"passout" gorm:"type:SMALLINT UNSIGNED"`
	DegreeID         uuid.UUID               `json:"-" gorm:"type:varchar(36)"`
	SpecializationID uuid.UUID               `json:"-" gorm:"type:varchar(36)"`
	FacultyID        uuid.UUID               `json:"facultyID" gorm:"type:varchar(36)"`
}

// ValidateFacultyAcademic validates fields of faculty academics.
func (academic *Academic) ValidateFacultyAcademic() error {
	yearPattern := regexp.MustCompile(`^(?:19|20)\d{2}`)

	if util.IsEmpty(academic.College) {
		return errors.NewValidationError("College name is required and should contain only characters")
	}

	if !yearPattern.MatchString(strconv.Itoa(academic.Passout)) {
		return errors.NewValidationError("Passout year can be from 1900 to 2099")
	}

	if academic.Specialization == nil {
		return errors.NewValidationError("Specialization is required")
	}

	if err := academic.Specialization.Validate(); err != nil {
		return err
	}

	if academic.Degree == nil {
		return errors.NewValidationError("Degree is required")
	}

	if err := academic.Degree.Validate(); err != nil {
		return err
	}

	return nil
}

// TableName overrides name of the table.
func (*Academic) TableName() string {
	return "faculty_academics"
}

// AcademicDTO will provide faculty's acadmeic details.
type AcademicDTO struct {
	ID               uuid.UUID               `json:"id"`
	DeletedAt        *time.Time              `json:"-"`
	Degree           *general.Degree         `json:"degree" gorm:"foreignkey:DegreeID;"`
	DegreeID         uuid.UUID               `json:"-"`
	Specialization   *general.Specialization `json:"specialization" gorm:"foreignkey:SpecializationID"`
	SpecializationID uuid.UUID               `json:"-"`

	College    string     `json:"college"`
	CollegeID  *uuid.UUID `json:"collegeID" gorm:"column:college_branch_id"`
	Percentage float32    `json:"percentage"`
	Passout    int        `json:"passout"`

	FacultyID uuid.UUID `json:"facultyID"`
}

// TableName overrides name of the table.
func (*AcademicDTO) TableName() string {
	return "faculty_academics"
}
