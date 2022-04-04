package list

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// BatchSession contains the fields for sessions in batch as well as session name from course_sessions.
type BatchSession struct {
	ID          uuid.UUID  `json:"id"`
	DeletedAt   *time.Time `json:"-"`
	BatchID     uuid.UUID  `json:"-"`
	Batch       Batch      `json:"batch"`
	Date        string     `json:"date"`
	IsCompleted *bool      `json:"isCompleted"`
}
