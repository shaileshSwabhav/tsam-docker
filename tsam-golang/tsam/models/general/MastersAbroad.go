package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// MastersAbroad contains the masters abroad details of talent which are enough for adding and updating talent masters abroad.
type MastersAbroad struct {
	TenantBase

	// Related table IDs.
	TalentID  *uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	EnquiryID *uuid.UUID `json:"enquiryID" gorm:"type:varchar(36)"`
	DegreeID  uuid.UUID  `json:"degreeID" gorm:"type:varchar(36)"`

	// Maps.
	Universities []University `json:"universities" gorm:"many2many:masters_abroad_universities;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Countries    []Country    `json:"countries" gorm:"many2many:masters_abroad_countries;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`

	// Child tables.
	Scores []*Score `json:"scores"`

	// Other fields.
	YearOfMS *uint `json:"yearOfMS" gorm:"type:SMALLINT(4)"`
}

// TableName will name the table of MastersAbroad model as "masters_abroad".
func (*MastersAbroad) TableName() string {
	return "masters_abroad"
}

// Validate masters in abroad fields.
func (mastersAbroad *MastersAbroad) Validate() error {

	// Check if at least one score is present or not.
	if len(mastersAbroad.Scores) == 0 {
		return errors.NewValidationError("Score(s) must be provided")
	}

	// Validate each score.
	if mastersAbroad.Scores != nil && len(mastersAbroad.Scores) != 0 {
		for _, score := range mastersAbroad.Scores {
			err := score.Validate()
			if err != nil {
				return err
			}
		}
	}

	// Check if countries is present or not.
	if mastersAbroad.Countries != nil && len(mastersAbroad.Countries) == 0 {
		return errors.NewValidationError("Country(s) must be provided")
	}

	// Check if universities is present or not.
	if mastersAbroad.Universities != nil && len(mastersAbroad.Universities) == 0 {
		return errors.NewValidationError("University(s) must be provided")
	}

	// Degree ID must be secified.
	if !util.IsUUIDValid(mastersAbroad.DegreeID) {
		return errors.NewValidationError("Degree ID must be secified")
	}

	// currentYear := time.Now().Year()
	// if mastersAbroad.GRE < uint16(currentYear) {
	// 	return errors.NewValidationError("GRE score should be greater than " + strconv.Itoa(currentYear))
	// }

	// if mastersAbroad.GRE > uint16(currentYear+4) {
	// 	return errors.NewValidationError("GRE score should be less than " + strconv.Itoa(currentYear+4))
	// }
	return nil
}

// ===========Defining many to many structs===========

// MastersAbroadCountries is the map of masters abroad and country.
type MastersAbroadCountries struct {
	MastersAbroadID uuid.UUID `gorm:"type:varchar(36)"`
	CountryID       uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*MastersAbroadCountries) TableName() string {
	return "masters_abroad_countries"
}

// MastersAbroadUniversities is the map of masters abroad and university.
type MastersAbroadUniversities struct {
	MastersAbroadID uuid.UUID `gorm:"type:varchar(36)"`
	UniversityID    uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*MastersAbroadUniversities) TableName() string {
	return "masters_abroad_universities"
}

//************************************* DTO MODEL *************************************************************

// MastersAbroadDTO contains the complete information of talent masters abroad which is needed to display.
type MastersAbroadDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Maps.
	Universities []University `json:"universities" gorm:"many2many:masters_abroad_universities;association_jointable_foreignkey:university_id;jointable_foreignkey:masters_abroad_id"`
	Countries    []Country    `json:"countries" gorm:"many2many:masters_abroad_countries;association_jointable_foreignkey:country_id;jointable_foreignkey:masters_abroad_id"`

	// Related tables and IDs.
	Degree    Degree     `json:"degree" gorm:"foreignkey:DegreeID"`
	DegreeID  uuid.UUID  `json:"-"`
	TalentID  *uuid.UUID `json:"talentID"`
	EnquiryID *uuid.UUID `json:"enquiryID"`

	// Child tables.
	Scores []*ScoreDTO `json:"scores" gorm:"foreignkey:MastersAbroadID"`

	// Other fields.
	YearOfMS *uint `json:"yearOfMS"`
}

// TableName will name the table of MastersAbroad model as "masters_abroad".
func (*MastersAbroadDTO) TableName() string {
	return "masters_abroad"
}
