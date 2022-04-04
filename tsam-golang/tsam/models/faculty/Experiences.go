package faculty

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Experience contain single FacultyExperience details.
type Experience struct {
	general.TenantBase
	CompanyName   string                `json:"companyName"`
	Designation   *general.Designation  `json:"designation" gorm:"foreignkey:DesignationID;association_autocreate:false;association_autoupdate:false;"`
	Technologies  []*general.Technology `json:"technologies" gorm:"many2many:faculty_experiences_technologies;association_autocreate:false;association_autoupdate:false;"`
	FromDate      *string               `json:"fromDate" gorm:"type:date"`
	ToDate        *string               `json:"toDate" gorm:"type:date"`
	Package       *uint                 `json:"package" gorm:"type:int(10)"`
	DesignationID uuid.UUID             `json:"-" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	FacultyID     uuid.UUID             `json:"facultyID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	// YearOfExperience float32               `json:"yearOfExperience" gorm:"type:float"`
}

// TableName overrides name of the table.
func (*Experience) TableName() string {
	return "faculty_experiences"
}

// ValidateFacultyExperiences validates the fields of FacultyExperiences.
func (experience *Experience) ValidateFacultyExperiences() error {

	if util.IsEmpty(experience.CompanyName) {
		return errors.NewValidationError("Company name must be specified")
	}

	if experience.Designation == nil {
		return errors.NewValidationError("Designation must be specified")
		// if err := experience.Designation.Validate(); err != nil {
		// 	return errors.NewValidationError(err.Error())
		// }
	}
	if experience.Technologies == nil {
		return errors.NewValidationError("Technologies must be specified")

	}
	// if experience.Technologies != nil {
	// 	for _, technology := range experience.Technologies {
	// 		if err := technology.Validate(); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// if experience.YearOfExperience == 0 {
	// 	return errors.NewValidationError("Years of experience must be specified")
	// }

	return nil
}

// ExperienceDTO will provide faculty's experence details.
type ExperienceDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	Technologies []*general.Technology `json:"technologies" gorm:"many2many:faculty_experiences_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:experience_id"`

	Designation   *general.Designation `json:"designation" gorm:"foreignkey:DesignationID"`
	DesignationID uuid.UUID            `json:"-"`

	CompanyName string    `json:"companyName"`
	FromDate    *string   `json:"fromDate"`
	ToDate      *string   `json:"toDate"`
	Package     *uint     `json:"package"`
	FacultyID   uuid.UUID `json:"facultyID"`
}

// TableName overrides name of the table.
func (*ExperienceDTO) TableName() string {
	return "faculty_experiences"
}
