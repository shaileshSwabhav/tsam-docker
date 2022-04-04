package company

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Company strcut contain company information
type Company struct {
	general.TenantBase
	CompanyName string    `json:"companyName" gorm:"type:varchar(200)"`
	About       string    `json:"about" gorm:"type:varchar(2000)"`
	Logo        string    `json:"logo" gorm:"type:varchar(200)"`
	Code        string    `json:"code" gorm:"type:varchar(10);not null"`
	Website     *string   `json:"website,omitempty" gorm:"type:varchar(100)"`
	Branches    []*Branch `json:"branches,omitempty" gorm:"foreignkey:CompanyID"`
}

// ValidateCompany validates company
func (company *Company) ValidateCompany() error {
	if util.IsEmpty(company.CompanyName) {
		return errors.NewValidationError("Company Name must be specified.")
	}
	if util.IsEmpty(company.Logo) {
		return errors.NewValidationError("Company Logo must be specified.")
	}
	if company.ID != uuid.Nil {
		if company.Branches == nil {
			return errors.NewValidationError("Company must have atleast 1 branch")
		}
		for _, branch := range company.Branches {
			if err := branch.ValidateCompanyBranch(); err != nil {
				return err
			}
		}
	}
	return nil
}

// CompanyDTO strcut contain company information
type CompanyDTO struct {
	ID          uuid.UUID  `json:"id"`
	DeletedAt   *time.Time `json:"-"`
	CompanyName string     `json:"companyName"`
	About       string     `json:"about"`
	Logo        string     `json:"logo"`
	Code        string     `json:"code"`
	Website     *string    `json:"website"`
	Branches    []*Branch  `json:"branches" gorm:"foreignkey:CompanyID"`
}

func (*CompanyDTO) TableName() string {
	return "companies"
}
