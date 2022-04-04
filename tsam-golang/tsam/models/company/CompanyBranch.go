package company

import (
	uuid "github.com/satori/go.uuid"

	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
)

// Branch contains all fields required in a company branch.
type Branch struct {
	general.TenantBase
	general.Address
	CompanyID              uuid.UUID             `json:"companyID,omitempty" gorm:"type:varchar(36)"`
	Code                   string                `json:"code" gorm:"type:varchar(10);not null"`
	BranchName             string                `json:"branchName" gorm:"type:varchar(150)"`
	MainBranch             *bool                 `json:"mainBranch,omitempty"`
	Domains                []*Domain             `json:"domains" gorm:"many2many:company_branches_domains;association_autoupdate:false"`
	Technologies           []*general.Technology `json:"technologies,omitempty" gorm:"many2many:company_branches_technologies;association_autoupdate:false"`
	HRHeadName             *string               `json:"hrHeadName,omitempty" gorm:"type:varchar(100)"`
	HRHeadContact          *string               `json:"hrHeadContact,omitempty" gorm:"type:varchar(15)"`
	HRHeadEmail            *string               `json:"hrHeadEmail,omitempty" gorm:"type:varchar(100)"`
	TechnologyHeadName     *string               `json:"technologyHeadName,omitempty" gorm:"type:varchar(100)"`
	TechnologyHeadContact  *string               `json:"technologyHeadContact,omitempty" gorm:"type:varchar(15)"`
	TechnologyHeadEmail    *string               `json:"technologyHeadEmail,omitempty" gorm:"type:varchar(100)"`
	UnitHeadName           *string               `json:"unitHeadName,omitempty" gorm:"type:varchar(100)"`
	UnitHeadContact        *string               `json:"unitHeadContact,omitempty" gorm:"type:varchar(15)"`
	UnitHeadEmail          *string               `json:"unitHeadEmail,omitempty" gorm:"type:varchar(100)"`
	FinanceHeadName        *string               `json:"financeHeadName,omitempty" gorm:"type:varchar(100)"`
	FinanceHeadContact     *string               `json:"financeHeadContact,omitempty" gorm:"type:varchar(15)"`
	FinanceHeadEmail       *string               `json:"financeHeadEmail,omitempty" gorm:"type:varchar(100)"`
	RecruitmentHeadName    *string               `json:"recruitmentHeadName,omitempty" gorm:"type:varchar(100)"`
	RecruitmentHeadContact *string               `json:"recruitmentHeadContact,omitempty" gorm:"type:varchar(15)"`
	RecruitmentHeadEmail   *string               `json:"recruitmentHeadEmail,omitempty" gorm:"type:varchar(100)"`
	NumberOfEmployees      *uint                 `json:"numberOfEmployees,omitempty"`
	CompanyRating          *uint8                `json:"companyRating,omitempty" gorm:"type:varchar(2)"`
	SalesPersonID          *uuid.UUID            `json:"-" gorm:"type:varchar(36)"`
	SalesPerson            *SalesPerson          `json:"salesPerson" gorm:"foreignkey:SalesPersonID;association_autoupdate:false"`
	OnePager               *string               `json:"onePager" gorm:"type:varchar(200)"`
	TermsAndConditions     *string               `json:"termsAndConditions" gorm:"type:varchar(200)"`
}

// TableName will name the table of branch model as "company_branches"
func (*Branch) TableName() string {
	return "company_branches"
}

// ValidateCompanyBranch validates all fields of the company's branch
func (branch *Branch) ValidateCompanyBranch() error {

	if err := branch.Address.ValidateAddress(); err != nil {
		return err
	}
	if branch.Technologies == nil {
		return errors.NewValidationError("Atleast 1 technology must be specified.")
	}
	// for _, technology := range branch.Technologies {
	// 	if err := technology.Validate(); err != nil {
	// 		return err
	// 	}
	// }
	if branch.Domains != nil {
		for _, domain := range branch.Domains {
			if err := domain.ValidateDomain(); err != nil {
				return err
			}
		}
	}
	return nil
}
