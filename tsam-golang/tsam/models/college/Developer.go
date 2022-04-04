package college

import (
	uuid "github.com/satori/go.uuid"
)

// Developer contains limited fields of employee model
type Developer struct {
	ID        uuid.UUID `json:"id" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	FirstName string    `json:"firstName" example:"Ravi"`
	LastName  string    `json:"lastName" example:"Sharma"`
}

// TableName will refer table "faculties" for model references.
func (*Developer) TableName() string {
	return "employees"
}
