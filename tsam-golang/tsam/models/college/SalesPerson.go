package college

import (
	uuid "github.com/satori/go.uuid"
)

// SalesPerson contains the fields of User model which is need by the college branch.
type SalesPerson struct {
	ID        uuid.UUID `json:"id" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	FirstName string    `json:"firstName" example:"Ravi"`
	LastName  string    `json:"lastName" example:"Sharma"`
}

// TableName will refer table "users" for model referrences.
func (*SalesPerson) TableName() string {
	return "users"
}
