package list

import (
	uuid "github.com/satori/go.uuid"
)

// Batch is used for listing of batches
type Batch struct {
	ID        uuid.UUID `json:"id"`
	BatchName string    `json:"batchName"`
	Code      string    `json:"code"`
	CourseID  uuid.UUID `json:"courseID"`
}
