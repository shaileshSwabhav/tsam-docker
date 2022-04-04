package batch

import uuid "github.com/satori/go.uuid"

// Search used to perform search operation
type Search struct {
	BatchID        string    `json:"batchID"`
	BatchName      string    `json:"batchName"`
	BatchObjective string    `json:"batchObjective"`
	Status         string    `json:"batchStatus"`
	IsActive       *bool     `json:"isActive"`
	StartDate      string    `json:"startDate"`
	EndDate        string    `json:"endDate"`
	CourseID       uuid.UUID `json:"courseID"`
	FacultyID      uuid.UUID `json:"facultyID"`
	SalesPersonID  uuid.UUID `json:"salesPersonID"`
	TenantID       uuid.UUID `json:"tenantID"`
}
