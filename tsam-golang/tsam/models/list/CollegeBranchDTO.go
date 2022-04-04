package list

import (
	uuid "github.com/satori/go.uuid"
)

// Branch is used for listing of colleges branches
type Branch struct {
	ID         uuid.UUID `json:"id" gorm:"type:varchar(36)"`
	BranchName string    `json:"branchName"`
	Code       string    `json:"code"`
}

// TableName will create the table for Branch model with name college_branches.
func (*Branch) TableName() string {
	return "college_branches"
}
