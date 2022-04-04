package list

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Specialization contains the details for specialization of a degree
type Specialization struct {
	ID         uuid.UUID  `json:"id"`
	DeletedAt  *time.Time `json:"-"`
	BranchName string     `json:"branchName"`
	DegreeID   uuid.UUID  `json:"degreeID"`
	// Degree     general.Degree `json:"degree"`
}
