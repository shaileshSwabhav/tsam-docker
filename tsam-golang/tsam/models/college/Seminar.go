package college

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Seminar defines fields of seminar.
type Seminar struct {
	general.TenantBase

	// Maps.
	CollegeBranches []list.Branch  `json:"collegeBranches" gorm:"many2many:seminars_college_branches;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	SalesPeople     []*SalesPerson `json:"salesPeople" gorm:"many2many:seminars_sales_people;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Speakers        []*Speaker     `json:"speakers" gorm:"many2many:seminars_speakers;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`

	// Flags.
	IsActive bool `json:"isActive"`

	// Other fields.
	SeminarName             string  `json:"seminarName" gorm:"type:varchar(50)"`
	Description             *string `json:"description" gorm:"type:varchar(500)"`
	Location                *string `json:"location" gorm:"type:varchar(500)"`
	Code                    string  `json:"code" gorm:"varchar(10)"`
	SeminarDate             string  `json:"seminarDate" gorm:"type:date"`
	FromTime                string  `json:"fromTime" gorm:"type:time"`
	ToTime                  string  `json:"toTime" gorm:"type:time"`
	StudentRegistrationLink *string `json:"studentRegistrationLink" gorm:"type:varchar(100)"`
}

// Validate validates all fields of the seminar.
func (seminar *Seminar) Validate() error {

	// College branches.
	if seminar.CollegeBranches == nil {
		return errors.NewValidationError("College Branch must be specified")
	}

	// Seminar name.
	if util.IsEmpty(seminar.SeminarName) {
		return errors.NewValidationError("Seminar Name must be specified")
	}

	if len(seminar.SeminarName) > 50 {
		return errors.NewValidationError("Seminar Name can have maximum 50 characters")
	}

	// Description.
	if seminar.Description != nil && len(*seminar.Description) > 500 {
		return errors.NewValidationError("Seminar description can have maximum 500 characters")
	}

	// Location.
	if seminar.Location != nil && len(*seminar.Location) > 500 {
		return errors.NewValidationError("Seminar location can have maximum 500 characters")
	}

	// Seminar date.
	if util.IsEmpty(seminar.SeminarDate) {
		return errors.NewValidationError("Seminar date must be specified")
	}

	// Seminar from time.
	if util.IsEmpty(seminar.FromTime) {
		return errors.NewValidationError("Seminar from time must be specified")
	}

	// Seminar to time.
	if util.IsEmpty(seminar.ToTime) {
		return errors.NewValidationError("Seminar to time must be specified")
	}

	return nil
}

// ===========Defining many to many structs===========

// SeminarCollegeBranch is the map of seminar and college branch.
type SeminarCollegeBranch struct {
	SeminarID uuid.UUID `gorm:"type:varchar(36)"`
	BranchID  uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*SeminarCollegeBranch) TableName() string {
	return "seminars_college_branches"
}

//************************************* DTO MODEL *************************************************************

// SeminarDTO defines fields of seminar.
type SeminarDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Maps.
	CollegeBranches []list.Branch  `json:"collegeBranches" gorm:"many2many:seminars_college_branches;association_jointable_foreignkey:branch_id;jointable_foreignkey:seminar_id"`
	SalesPeople     []*SalesPerson `json:"salesPeople" gorm:"many2many:seminars_sales_people;association_jointable_foreignkey:sales_person_id;jointable_foreignkey:seminar_id"`
	Speakers        []*Speaker     `json:"speakers" gorm:"many2many:seminars_speakers;association_jointable_foreignkey:speaker_id;jointable_foreignkey:seminar_id"`

	// Flags.
	IsActive bool `json:"isActive"`

	// Other fields.
	SeminarName             string  `json:"seminarName"`
	Description             *string `json:"description"`
	Location                *string `json:"location"`
	Code                    string  `json:"code"`
	SeminarDate             string  `json:"seminarDate"`
	FromTime                string  `json:"fromTime"`
	ToTime                  string  `json:"toTime"`
	StudentRegistrationLink *string `json:"studentRegistrationLink"`
	TotalRegisteredStudents *uint16 `json:"totalRegisteredStudents"`
	TotalVisitedStudents    *uint16 `json:"totalVisitedStudents"`
}

// TableName defines table name of the struct.
func (*SeminarDTO) TableName() string {
	return "seminars"
}
