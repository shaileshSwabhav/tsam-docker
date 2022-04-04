package course

import uuid "github.com/satori/go.uuid"

// Search Contain Technology, CourseType and CourseStatus.
type Search struct {
	CreatedAt    string      `json:"createdAt"`
	Technologies []uuid.UUID `json:"technologies"`
	CourseType   string      `json:"courseType"`
	CourseLevel  string      `json:"courseLevel"`
}
