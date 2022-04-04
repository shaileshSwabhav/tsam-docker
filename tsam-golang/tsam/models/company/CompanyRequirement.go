package company

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Requirement contains information of company requirements
type Requirement struct {
	general.TenantBase
	general.Address `json:"jobLocation"`

	// Maps.
	Qualifications  []*general.Degree    `json:"qualifications,omitempty" gorm:"many2many:company_requirements_qualifications;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Universities    []*University        `json:"universities" gorm:"many2many:company_requirements_universities;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Technologies    []general.Technology `json:"technologies" gorm:"many2many:company_requirements_technologies;association_autocreate:false;association_autoupdate:false"`
	SelectedTalents []*Talent            `json:"selectedTalents,omitempty" gorm:"many2many:company_requirements_talents;association_autocreate:false;association_autoupdate:false"`

	// Related table IDs.
	SalesPersonID *uuid.UUID `json:"salesPersonID" gorm:"type:varchar(36)"`
	CompanyID     uuid.UUID  `json:"companyID" gorm:"type:varchar(36)"`
	DesignationID uuid.UUID  `json:"designationID" gorm:"type:varchar(36)"`

	// Other fields.
	IsActive          bool    `json:"isActive" gorm:"DEFAULT:true"`
	Code              string  `json:"code" gorm:"type:varchar(10);not null"`
	TalentRating      *uint8  `json:"talentRating" gorm:"type:tinyint(2)"`
	PersonalityType   *string `json:"personalityType" gorm:"type:varchar(50)"`
	MinimumExperience *uint8  `json:"minimumExperience" gorm:"type:tinyint(2)"`
	MaximumExperience *uint8  `json:"maximumExperience" gorm:"type:tinyint(2)"`
	JobDescription    string  `json:"jobDescription" gorm:"type:varchar(2000)"`
	JobRequirement    string  `json:"jobRequirement" gorm:"type:varchar(2000)"`
	JobType           *string `json:"jobType" gorm:"type:varchar(200)"`
	// JobRole           string  `json:"jobRole" gorm:"varchar(200)"`
	// PackageOffered     uint64  `json:"packageOffered" gorm:"type:int(10)"`
	MinimumPackage     uint64  `json:"minimumPackage" gorm:"type:int(10)"`
	MaximumPackage     uint64  `json:"maximumPackage" gorm:"type:int(10)"`
	RequiredBefore     string  `json:"requiredBefore" gorm:"type:date"`
	RequiredFrom       string  `json:"requiredFrom" gorm:"type:date"`
	Vacancy            uint16  `json:"vacancy" gorm:"type:smallint(5)"`
	Comment            *string `json:"comment" gorm:"type:varchar(4000)"`
	TermsAndConditions *string `json:"termsAndConditions" gorm:"type:varchar(200)"`
	SampleOfferLetter  *string `json:"sampleOfferLetter" gorm:"type:varchar(200)"`

	// Criteria.
	Increment        *string `json:"increment" gorm:"type:varchar(50)"`
	WeeklyHoliday    *string `json:"weeklyHoliday" gorm:"type:varchar(50)"`
	Qualification    *string `json:"qualification" gorm:"type:varchar(50)"`
	BondPeriod       *string `json:"bondPeriod" gorm:"type:varchar(50)"`
	TalentLocation   *string `json:"talentLocation" gorm:"type:varchar(50)"`
	WorkShift        *string `json:"workShift" gorm:"type:varchar(50)"`
	CompanyType      *string `json:"companyType" gorm:"type:varchar(50)"`
	JoiningPeriod    *string `json:"joiningPeriod" gorm:"type:varchar(50)"`
	GenderPreference *string `json:"genderPreference" gorm:"type:varchar(50)"`
	// Criteria10     *string `json:"criteria10" gorm:"type:varchar(50)"`

	// Rating.
	Rating *float64 `json:"rating" gorm:"type:decimal(4,2)"`

	// Undecided fields.
	// MarksCriteria      *uint8  `json:"marksCriteria,omitempty" gorm:"type:tinyint(2)"`
	// Colleges        []*College           `json:"colleges,omitempty" gorm:"many2many:company_requirements_colleges;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`

}

// TableName will name the table of requirement model as "company_requirements"
func (*Requirement) TableName() string {
	return "company_requirements"
}

// Validate validate fields of company requirement
func (requirement *Requirement) Validate() error {

	// requiredfrom and requiredBefore validation.
	requiredFrom, err := time.Parse(time.RFC3339, requirement.RequiredFrom)
	if err != nil {
		return err
	}

	requiredBefore, err := time.Parse(time.RFC3339, requirement.RequiredBefore)
	if err != nil {
		return err
	}

	// After reports whether the time instant requiredFrom is after requiredBefore.
	if requiredFrom.After(requiredBefore) {
		return errors.NewValidationError("Required before must be greater than required from.")
	}

	requirement.RequiredFrom = requiredFrom.Format("2006-01-02")
	requirement.RequiredBefore = requiredBefore.Format("2006-01-02")

	// Validate address.
	if err := requirement.Address.MandatoryValidation(); err != nil {
		return err
	}

	// Validate company branch.
	if requirement.CompanyID == uuid.Nil {
		return errors.NewValidationError("Company Branch ID must be specified.")
	}

	// Validate minimum experience.
	if requirement.MinimumExperience != nil && *requirement.MinimumExperience < 1 {
		return errors.NewValidationError("Minimum Experience cannot be below 1.")
	}

	if requirement.MinimumExperience != nil && *requirement.MinimumExperience > 30 {
		return errors.NewValidationError("Minimum Experience cannot be above 30.")
	}
	// Validate maximum experience.
	if requirement.MaximumExperience != nil && *requirement.MaximumExperience < 1 {
		return errors.NewValidationError("Maximum Experience cannot be below 1.")
	}

	if requirement.MaximumExperience != nil && *requirement.MaximumExperience > 30 {
		return errors.NewValidationError("Maximum Experience cannot be above 30.")
	}

	// Validate designation.
	if requirement.DesignationID == uuid.Nil {
		return errors.NewValidationError("Designation must be specified.")
	}

	// Validate job description.
	if util.IsEmpty(requirement.JobDescription) {
		return errors.NewValidationError("Job description must be specified.")
	}

	// Validate job type.
	if util.IsEmpty(*requirement.JobType) {
		return errors.NewValidationError("Job type must be specified.")
	}

	// Validate job type.
	if util.IsEmpty(requirement.JobRequirement) {
		return errors.NewValidationError("Job requirement must be specified.")
	}

	// Validate pacakage offered.
	if requirement.MinimumPackage < 100000 {
		return errors.NewValidationError("Minimum package cannot be below 100000.")
	}

	if requirement.MaximumPackage > 1000000000 {
		return errors.NewValidationError("Maximum package offered cannot be above 1000000000.")
	}

	if requirement.MaximumPackage < requirement.MinimumPackage {
		return errors.NewValidationError("Maximum package should be greater than minimum package.")
	}

	// Validate required before.
	if util.IsEmpty(requirement.RequiredBefore) {
		return errors.NewValidationError("Required before date must be specified.")
	}

	// Validate required from.
	if util.IsEmpty(requirement.RequiredFrom) {
		return errors.NewValidationError("Required from date must be specified")
	}

	// Validate vacancy.
	if requirement.Vacancy <= 0 {
		return errors.NewValidationError("Vacancy cannot be below 0.")
	}

	if requirement.Vacancy > 10000 {
		return errors.NewValidationError("Vacancy cannot be above 10000.")
	}

	// Validate comment.
	if requirement.Comment != nil && len(*requirement.Comment) > 4000 {
		return errors.NewValidationError("Comment can have maximum 4000 characters.")
	}

	// Validate ratingg.
	if requirement.Rating != nil && *requirement.Rating < 0.0 {
		return errors.NewValidationError("Rating must be specified.")
	}

	return nil
}

// RequirementUpdate is the struct for updating company requirement' salesperson.
type RequirementUpdate struct {
	RequirementID uuid.UUID `json:"requirementID"`
}

// ===========Defining many to many structs===========

// CompanyRequirementQualifications is the map of company requorement and qualification.
type CompanyRequirementQualifications struct {
	RequirementID uuid.UUID `gorm:"type:varchar(36)"`
	DegreeID      uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CompanyRequirementQualifications) TableName() string {
	return "company_requirements_qualifications"
}

// CompanyRequirementUniversities is the map of company requorement and university.
type CompanyRequirementUniversities struct {
	RequirementID uuid.UUID `gorm:"type:varchar(36)"`
	UniversityID  uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CompanyRequirementUniversities) TableName() string {
	return "company_requirements_universities"
}

// CompanyRequirementTechnologies is the map of company requorement and technology.
type CompanyRequirementTechnologies struct {
	RequirementID uuid.UUID `gorm:"type:varchar(36)"`
	TechnologyID  uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CompanyRequirementTechnologies) TableName() string {
	return "company_requirements_technologies"
}

// CompanyRequirementTalents is the map of company requorement and talent.
type CompanyRequirementTalents struct {
	RequirementID uuid.UUID `gorm:"type:varchar(36)"`
	TalentID      uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CompanyRequirementTalents) TableName() string {
	return "company_requirements_talents"
}

// RequirementDTO contains information of company requirements with extra fields to be displayed in UI.
type RequirementDTO struct {
	// ID        uuid.UUID  `json:"id"`
	// DeletedAt *time.Time `json:"-"`
	general.TenantBase
	general.Address `json:"jobLocation"`

	// Maps.
	Qualifications  []*general.Degree    `json:"qualifications,omitempty" gorm:"many2many:company_requirements_qualifications;association_jointable_foreignkey:degree_id;jointable_foreignkey:requirement_id"`
	Universities    []*University        `json:"universities" gorm:"many2many:company_requirements_universities;association_jointable_foreignkey:university_id;jointable_foreignkey:requirement_id"`
	Technologies    []general.Technology `json:"technologies" gorm:"many2many:company_requirements_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:requirement_id"`
	SelectedTalents []*Talent            `json:"selectedTalents,omitempty" gorm:"many2many:company_requirements_talents;association_jointable_foreignkey:talent_id;jointable_foreignkey:requirement_id"`

	// Related table IDs.
	SalesPerson   *SalesPerson         `json:"salesPerson" gorm:"foreignkey:SalesPersonID"`
	SalesPersonID *uuid.UUID           `json:"-"`
	Company       CompanyDTO           `json:"company" gorm:"foreignkey:companyID"`
	CompanyID     uuid.UUID            `json:"-"`
	Designation   *general.Designation `json:"designation" gorm:"foreignkey:DesignationID"`
	DesignationID uuid.UUID            `json:"-"`

	// Other fields.
	IsActive          bool    `json:"isActive"`
	Code              string  `json:"code"`
	TalentRating      *uint8  `json:"talentRating"`
	PersonalityType   *string `json:"personalityType"`
	MinimumExperience *uint8  `json:"minimumExperience"`
	MaximumExperience *uint8  `json:"maximumExperience"`
	JobRole           string  `json:"jobRole"`
	JobDescription    string  `json:"jobDescription"`
	JobRequirement    string  `json:"jobRequirement"`
	JobType           *string `json:"jobType"`
	MinimumPackage    uint64  `json:"minimumPackage"`
	MaximumPackage    uint64  `json:"maximumPackage"`
	// PackageOffered     uint64   `json:"packageOffered"`
	RequiredBefore     string  `json:"requiredBefore"`
	RequiredFrom       string  `json:"requiredFrom"`
	Vacancy            uint16  `json:"vacancy"`
	Comment            *string `json:"comment"`
	TermsAndConditions *string `json:"termsAndConditions"`
	SampleOfferLetter  *string `json:"sampleOfferLetter"`
	TotalApplicants    *uint16 `json:"totalApplicants"`

	// Criteria
	Increment        *string `json:"increment"`
	WeeklyHoliday    *string `json:"weeklyHoliday"`
	Qualification    *string `json:"qualification"`
	BondPeriod       *string `json:"bondPeriod"`
	TalentLocation   *string `json:"talentLocation"`
	WorkShift        *string `json:"workShift"`
	CompanyType      *string `json:"companyType"`
	JoiningPeriod    *string `json:"joiningPeriod"`
	GenderPreference *string `json:"genderPreference"`
	// Criteria10     *string `json:"criteria10"`

	Rating *float64 `json:"rating"`

	// Salary trend
	SalaryTrend *general.SalaryTrend `json:"salaryTrend"`
}

// TableName defines table name of the struct.
func (*RequirementDTO) TableName() string {
	return "company_requirements"
}
