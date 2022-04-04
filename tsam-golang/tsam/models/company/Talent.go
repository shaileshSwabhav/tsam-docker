package company

import (
	uuid "github.com/satori/go.uuid"
)

// Talent Contains ID of the talent.
type Talent struct {
	ID uuid.UUID `json:"id" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	// model.base
	// model.Address
	// FirstName string `json:"firstName" gorm:"type:varchar(100)"`
	// LastName    string              `json:"lastName,omitempty" gorm:"type:varchar(100)"`
	// Experiences []*model.Experience `json:"experiences,omitempty" gorm:"many2many:talent_experiences"`
}
