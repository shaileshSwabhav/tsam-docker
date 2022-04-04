package list

import (
	uuid "github.com/satori/go.uuid"
)

// College is used for listing of colleges
type College struct {
	ID          uuid.UUID `json:"id"`
	CollegeName string    `json:"collegeName"`
	Code        string    `json:"code"`
}
