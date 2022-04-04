package list

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Requirement struct {
	ID              uuid.UUID  `json:"id"`
	DeletedAt       *time.Time `json:"-"`
	Code            string     `json:"code"`
	CompanyBranchID uuid.UUID  `json:"companyBranchID" gorm:"type:varchar(36)"`
	Branch          BranchName `json:"branch" gorm:"FOREIGNKEY:CompanyBranchID"`
}

// BranchName contains the name and ID of a company branch.
type BranchName struct {
	ID              uuid.UUID `json:"id"`
	BranchName      string    `json:"branchName"`
	UnitHeadName    *string   `json:"unitHeadName"`
	UnitHeadContact *string   `json:"unitHeadContact"`
	UnitHeadEmail   *string   `json:"unitHeadEmail"`
}

// TableName will name the table of requirement model as "company_requirements"
func (*Requirement) TableName() string {
	return "company_requirements"
}

// TableName will name the table of requirement model as "company_requirements"
func (*BranchName) TableName() string {
	return "company_branches"
}
