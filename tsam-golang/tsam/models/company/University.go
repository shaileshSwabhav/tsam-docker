package company

import (
	uuid "github.com/satori/go.uuid"
)

// University contain infomation about University
type University struct {
	ID          uuid.UUID `json:"id" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	UniversityName string `json:"universityName"`
}
