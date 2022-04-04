package list

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type ProgrammingAssignment struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`
	Title     string     `json:"title"`
}

// TableName will create the table for Course model with name courses.
func (*ProgrammingAssignment) TableName() string {
	return "programming_assignments"
}
