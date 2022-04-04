package company

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
)

// ************************MY OPPURTUNITY DTO FOR MY-OPPORTUNITIES PAGE**********************************
// MyOpportunityDTO will return details about requirement, company branch and company.
type MyOpportunityDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Requirement.
	City           string              `json:"city"`
	JobType        string              `json:"jobType"`
	MinimumPackage uint64              `json:"minimumPackage"`
	MaximumPackage uint64              `json:"maximumPackage"`
	RequiredBefore string              `json:"requiredBefore"`
	Designation    general.Designation `json:"designation" gorm:"foreignkey:DesignationID"`
	DesignationID  uuid.UUID           `json:"-"`

	// company-branch
	CompanyBranch   MyOpportunityDTOCompanyBranch `json:"companyBranch" gorm:"foreignkey:CompanyBranchID"`
	CompanyBranchID uuid.UUID                     `json:"-"`
}

// TableName defines table name of the struct.
func (*MyOpportunityDTO) TableName() string {
	return "company_requirements"
}

// MyOpportunityDTOCompanyBranch is the additional information along with MyOpportunityDTO.
type MyOpportunityDTOCompanyBranch struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	BranchName string `json:"branchName"`

	// Company.
	CompanyID uuid.UUID               `json:"companyID"`
	Company   MyOpportunityDTOCompany `json:"company" gorm:"foreignkey:CompanyID"`
}

// TableName defines table name of the struct.
func (*MyOpportunityDTOCompanyBranch) TableName() string {
	return "company_branches"
}

// MyOpportunityDTOCompany is the additional information along with MyOpportunityDTO.
type MyOpportunityDTOCompany struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	Logo string `json:"logo" gorm:"type:varchar(200)"`
}

// TableName defines table name of the struct.
func (*MyOpportunityDTOCompany) TableName() string {
	return "companies"
}

// ************************COMPANY DETAILS DTO FOR COMPANY DETAILS PAGE**********************************
// CompanyDetails will return details about requirement, company branch and company.
type CompanyDetails struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Requirement.
	City              string              `json:"city"`
	JobDescription    string              `json:"jobDescription"`
	JobRequirement    string              `json:"jobRequirement"`
	JobType           string              `json:"jobType"`
	MinimumPackage    uint64              `json:"minimumPackage"`
	MaximumPackage    uint64              `json:"maximumPackage"`
	MinimumExperience *uint8              `json:"minimumExperience"`
	MaximumExperience *uint8              `json:"maximumExperience"`
	Designation       general.Designation `json:"designation" gorm:"foreignkey:DesignationID"`
	DesignationID     uuid.UUID           `json:"-"`

	// Company-branch.
	CompanyBranch   companyBranch `json:"companyBranch" gorm:"foreignkey:CompanyBranchID"`
	CompanyBranchID uuid.UUID     `json:"-"`
}

// TableName defines table name of the struct.
func (*CompanyDetails) TableName() string {
	return "company_requirements"
}

// companyBranch is the additional information along with CompanyDetails.
type companyBranch struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	BranchName        string `json:"branchName"`
	NumberOfEmployees *uint  `json:"numberOfEmployees,omitempty"`

	// Company.
	CompanyID uuid.UUID `json:"companyID"`
	Company   company   `json:"company" gorm:"foreignkey:CompanyID"`
}

// TableName defines table name of the struct.
func (*companyBranch) TableName() string {
	return "company_branches"
}

// company is the additional information along with CompanyDetails.
type company struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	CompanyName string `json:"companyName" gorm:"type:varchar(200)"`
	About       string `json:"about" gorm:"type:varchar(2000)"`
	Logo        string `json:"logo" gorm:"type:varchar(200)"`
}

// TableName defines table name of the struct.
func (*company) TableName() string {
	return "companies"
}
