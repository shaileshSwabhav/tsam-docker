package company

import (
	uuid "github.com/satori/go.uuid"
)

// College contain information about college branch.
type College struct {
	ID          uuid.UUID `json:"id" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	CollegeName string    `json:"collegeName" gorm:"type:varchar(100)"`
}
