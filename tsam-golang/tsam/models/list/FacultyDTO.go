package list

import (
	uuid "github.com/satori/go.uuid"
)

// Faculty is DTO which is used for listing purpose.
type Faculty struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
}

// FacultyCredentialDTO will contain details of faculty from credential table
type FacultyCredentialDTO struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	FacultyID uuid.UUID `json:"facultyID"`
}

// TableName will create the table for Faculty model with name faculties.
func (*Faculty) TableName() string {
	return "faculties"
}
