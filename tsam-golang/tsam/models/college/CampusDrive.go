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

// CampusDrive defines fields of campus drive.
type CampusDrive struct {
	general.TenantBase

	// Maps.
	SalesPeople         []*SalesPerson     `json:"salesPeople" gorm:"many2many:campus_drives_sales_people;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Faculties           []*Faculty         `json:"faculties" gorm:"many2many:campus_drives_faculty;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Developers          []*Developer       `json:"developers" gorm:"many2many:campus_drives_developers;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	CompanyRequirements []list.Requirement `json:"companyRequirements" gorm:"many2many:campus_drives_company_requirements;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	CollegeBranches     []list.Branch      `json:"collegeBranches" gorm:"many2many:campus_drives_college_branches;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`

	// Flags.
	Cancelled bool `json:"cancelled"`

	// Other fields.
	CampusName              string  `json:"campusName" gorm:"type:varchar(50)"`
	Description             *string `json:"description" gorm:"type:varchar(500)"`
	Location                *string `json:"location" gorm:"type:varchar(500)"`
	Code                    string  `json:"code" gorm:"varchar(10)"`
	CampusDate              string  `json:"campusDate" gorm:"type:datetime"`
	StudentRegistrationLink *string `json:"studentRegistrationLink" gorm:"type:varchar(100)"`

	// RegistrationLimit *uint   `json:"registrationLimit" gorm:"type:int(7)"`
	// PackageOffered            *uint                  `json:"packageOffered" gorm:"type:int(12)"`
}

// TotalVacancy is used to get total of vacancies of all company requirements.
type TotalVacancy struct {
	TotalVacancy int64
}

// Validate validates all fields of the campus drive.
func (campusDrive *CampusDrive) Validate() error {

	// College branches.
	if campusDrive.CollegeBranches == nil {
		return errors.NewValidationError("College Branch must be specified")
	}

	// Company requirements.
	if campusDrive.CompanyRequirements == nil {
		return errors.NewValidationError("Company Requirements must be specified")
	}

	// Campus date.
	if util.IsEmpty(campusDrive.CampusDate) {
		return errors.NewValidationError("Campus Drive date must be specified")
	}

	// Campus name.
	if util.IsEmpty(campusDrive.CampusName) {
		return errors.NewValidationError("Campus Name must be specified")
	}
	if len(campusDrive.CampusName) > 50 {
		return errors.NewValidationError("Campus Name can have maximum 50 characters")
	}

	// Description.
	if campusDrive.Description != nil && len(*campusDrive.Description) > 500 {
		return errors.NewValidationError("Campus description can have maximum 500 characters")
	}

	// Location.
	if campusDrive.Location != nil && len(*campusDrive.Location) > 500 {
		return errors.NewValidationError("Campus location can have maximum 500 characters")
	}

	return nil
}

// ===========Defining many to many structs===========

// CampusDriveSalesPerson is the map of campus drive and salesperson.
type CampusDriveSalesPerson struct {
	CampusDriveID uuid.UUID `gorm:"type:varchar(36)"`
	SalesPersonID uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CampusDriveSalesPerson) TableName() string {
	return "campus_drives_sales_people"
}

// CampusDriveFaculty is the map of campus drive and faculty.
type CampusDriveFaculty struct {
	CampusDriveID uuid.UUID `gorm:"type:varchar(36)"`
	FacultyID     uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CampusDriveFaculty) TableName() string {
	return "campus_drives_faculty"
}

// CampusDriveDeveloper is the map of campus drive and developer.
type CampusDriveDeveloper struct {
	CampusDriveID uuid.UUID `gorm:"type:varchar(36)"`
	DeveloperID   uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CampusDriveDeveloper) TableName() string {
	return "campus_drives_developers"
}

// CampusDriveCollegeBranch is the map of campus drive and college branch.
type CampusDriveCollegeBranch struct {
	CampusDriveID uuid.UUID `gorm:"type:varchar(36)"`
	BranchID      uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CampusDriveCollegeBranch) TableName() string {
	return "campus_drives_college_branches"
}

// CampusDriveCompanyRequirement is the map of campus drive and company requirement.
type CampusDriveCompanyRequirement struct {
	CampusDriveID uuid.UUID `gorm:"type:varchar(36)"`
	RequirementID uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CampusDriveCompanyRequirement) TableName() string {
	return "campus_drives_company_requirements"
}

//************************************* DTO MODEL *************************************************************

// CampusDrive defines fields of campus drive.
type CampusDriveDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Maps.
	SalesPeople         []*SalesPerson     `json:"salesPeople" gorm:"many2many:campus_drives_sales_people;association_jointable_foreignkey:sales_person_id;jointable_foreignkey:campus_drive_id"`
	Faculties           []*Faculty         `json:"faculties" gorm:"many2many:campus_drives_faculty;association_jointable_foreignkey:faculty_id;jointable_foreignkey:campus_drive_id"`
	Developers          []*Developer       `json:"developers" gorm:"many2many:campus_drives_developers;association_jointable_foreignkey:developer_id;jointable_foreignkey:campus_drive_id"`
	CompanyRequirements []list.Requirement `json:"companyRequirements" gorm:"many2many:campus_drives_company_requirements;association_jointable_foreignkey:requirement_id;jointable_foreignkey:campus_drive_id"`
	CollegeBranches     []list.Branch      `json:"collegeBranches" gorm:"many2many:campus_drives_college_branches;association_jointable_foreignkey:branch_id;jointable_foreignkey:campus_drive_id"`

	// Flags.
	Cancelled bool `json:"cancelled"`

	// Other fields.
	CampusName                string  `json:"campusName"`
	Description               *string `json:"description"`
	Location                  *string `json:"location"`
	Code                      string  `json:"code"`
	TotalRegisteredCandidates *uint16 `json:"totalRegisteredCandidates"`
	TotalAppearedCandidates   *uint16 `json:"totalAppearedCandidates"`
	TotalRequirements         *uint16 `json:"totalRequirements"`
	CampusDate                string  `json:"campusDate"`
	StudentRegistrationLink   *string `json:"studentRegistrationLink"`

	// RegistrationLimit *uint   `json:"registrationLimit" gorm:"type:int(7)"`
	// PackageOffered            *uint                  `json:"packageOffered" gorm:"type:int(12)"`
}

// TableName defines table name of the struct.
func (*CampusDriveDTO) TableName() string {
	return "campus_drives"
}
