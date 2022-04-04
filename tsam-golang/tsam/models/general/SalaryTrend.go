package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// SalaryTrend contains salary trend details for different experience.
type SalaryTrend struct {
	TenantBase
	Date              string       `json:"date" gorm:"type:date"`
	CompanyRating     uint         `json:"companyRating" gorm:"type:int"`
	MinimumExperience uint         `json:"minimumExperience" gorm:"type:int"`
	MaximumExperience uint         `json:"maximumExperience" gorm:"type:int"`
	MinimumSalary     uint         `json:"minimumSalary" gorm:"type:int"`
	MaximumSalary     uint         `json:"maximumSalary" gorm:"type:int"`
	MedianSalary      uint         `json:"medianSalary" gorm:"type:int"`
	Technology        *Technology  `json:"technology" gorm:"foreignkey:TechnologyID;association_autoupdate:false;association_autocreate:false;"`
	Designation       *Designation `json:"designation" gorm:"foreignkey:DesignationID;association_autoupdate:false;association_autocreate:false;"`
	TechnologyID      uuid.UUID    `json:"-" gorm:"type:varchar(36)"`
	DesignationID     uuid.UUID    `json:"-" gorm:"type:varchar(36)"`
}

// ValidateSalaryTrend will validate all the compuslory fields of salary-trend.
func (salaryTrend *SalaryTrend) ValidateSalaryTrend() error {

	if salaryTrend.CompanyRating == 0 {
		return errors.NewValidationError("Company rating must be specified")
	}

	if util.IsEmpty(salaryTrend.Date) {
		return errors.NewValidationError("Date must be specified")
	}

	if salaryTrend.MaximumExperience == 0 {
		return errors.NewValidationError("Maximum exprerience must be greater than 0")
	}

	if salaryTrend.MaximumExperience <= salaryTrend.MinimumExperience {
		return errors.NewValidationError("Maximum exprerience must be greater than minimum exprerience")
	}

	if salaryTrend.MinimumSalary == 0 {
		return errors.NewValidationError("Minimum salary must be specified")
	}

	if salaryTrend.MaximumSalary == 0 {
		return errors.NewValidationError("Maximum salary must be specified")
	}

	if salaryTrend.MedianSalary == 0 {
		return errors.NewValidationError("Median salary must be specified")
	}

	if salaryTrend.MaximumSalary < salaryTrend.MinimumSalary {
		return errors.NewValidationError("Maximum salary must be greater than minimum salary")
	}

	if salaryTrend.Technology == nil {
		return errors.NewValidationError("Technology must be specified")
	}

	if salaryTrend.Designation == nil {
		return errors.NewValidationError("Designation must be specified")
	}

	return nil
}

// SalaryTrendDTO contains salary trend details for different experience.
type SalaryTrendDTO struct {
	ID                uuid.UUID    `json:"id"`
	DeletedAt         *time.Time   `json:"-"`
	Date              string       `json:"date"`
	CompanyRating     string       `json:"companyRating"`
	MinimumExperience uint         `json:"minimumExperience"`
	MaximumExperience uint         `json:"maximumExperience"`
	MinimumSalary     uint         `json:"minimumSalary"`
	MaximumSalary     uint         `json:"maximumSalary"`
	MedianSalary      uint         `json:"medianSalary"`
	Technology        *Technology  `json:"technology" gorm:"foreignkey:TechnologyID"`
	Designation       *Designation `json:"designation" gorm:"foreignkey:DesignationID"`
	TechnologyID      uuid.UUID    `json:"-"`
	DesignationID     uuid.UUID    `json:"-"`
}

// TableName defines table name of the struct.
func (*SalaryTrendDTO) TableName() string {
	return "salary_trends"
}
