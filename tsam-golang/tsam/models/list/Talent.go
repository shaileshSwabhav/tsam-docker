package list

import uuid "github.com/satori/go.uuid"

// Talent is DTO which is used for listing purpose.
type Talent struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

// TableName will create the table for Talent model with name talents.
func (*Talent) TableName() string {
	return "talents"
}
