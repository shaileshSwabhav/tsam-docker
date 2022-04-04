package list

import (
	uuid "github.com/satori/go.uuid"
)

// CompanyBranch is used for listing of company branches.
type CompanyBranch struct {
	ID         uuid.UUID `json:"id" gorm:"type:varchar(36)"`
	BranchName string    `json:"branchName"`
	CompanyID  uuid.UUID `json:"companyID"`
}

// TableName will create the table for CompanyBranch model with name company_branches.
func (*CompanyBranch) TableName() string {
	return "company_branches"
}
