package talent

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Experience contains the experience details of talent which are enough for adding and updating talent experience.
type Experience struct {
	general.TenantBase

	// Maps.
	Technologies []*general.Technology `json:"technologies" gorm:"many2many:talent_experiences_technologies;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`

	// Related table IDs.
	DesignationID uuid.UUID `json:"designationID" gorm:"type:varchar(36)"`
	TalentID      uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`

	// Other fields.
	Company  string  `json:"company"`
	FromDate string  `json:"fromDate" gorm:"type:date"`
	ToDate   *string `json:"toDate" gorm:"type:date"`
	Package  *uint   `json:"package" gorm:"type:int"`
}

// TableName will name the table of Experience model as "talent_experiences".
func (*Experience) TableName() string {
	return "talent_experiences"
}

// Validate Validates fields of talent experience.
func (experience *Experience) Validate() error {
	// Designation ID.
	if !util.IsUUIDValid(experience.DesignationID) {
		return errors.NewValidationError("Designation ID must ne a proper uuid")
	}

	// Check if company name is blank or not.
	if util.IsEmpty(experience.Company) {
		return errors.NewValidationError("Company Name must be specified")
	}

	// Company name maximum length.
	if len(experience.Company) > 255 {
		return errors.NewValidationError("Company name can have maximum 255 characters")
	}

	// Check if from date is blank or not.
	if util.IsEmpty(experience.FromDate) {
		return errors.NewValidationError("From date must be specified")
	}

	// Package.
	if experience.Package != nil && *experience.Package < 100000 {
		return errors.NewValidationError("Package cannot be less than 100000")
	}

	if experience.Package != nil && *experience.Package > 100000000 {
		return errors.NewValidationError("Package cannot be more than 100000000")
	}

	return nil
}

// ===========Defining many to many structs===========

// ExperienceTechnology is the map of experience and technology.
type ExperienceTechnology struct {
	ExperienceID uuid.UUID `gorm:"type:varchar(36)"`
	TechnologyID uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*ExperienceTechnology) TableName() string {
	return "talent_experiences_technologies"
}

//************************************* DTO MODEL *************************************************************

// ExperienceDTO contains the complete information of talent experience which is needed to display.
type ExperienceDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Maps.
	Technologies []*general.Technology `json:"technologies" gorm:"many2many:talent_experiences_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:experience_id"`

	// Related tables.
	Designation   general.Designation `json:"designation" gorm:"foreignkey:DesignationID"`
	DesignationID uuid.UUID           `json:"-"`

	// Other fields.
	Company  string    `json:"company"`
	FromDate string    `json:"fromDate"`
	ToDate   *string   `json:"toDate"`
	Package  *uint     `json:"package"`
	TalentID uuid.UUID `json:"talentID"`
}

// TableName will name the table of Experience model as "talent_experiences".
func (*ExperienceDTO) TableName() string {
	return "talent_experiences"
}
