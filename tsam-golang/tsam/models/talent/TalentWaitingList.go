package talent

import (
	uuid "github.com/satori/go.uuid"
)

// TalentWaitingList is used for getting talent and waiting list realted fields for waiting list report.
type TalentWaitingList struct {

	// Talent related details.
	TalentID             uuid.UUID `json:"talentID"`
	FirstName            string    `json:"firstName"`
	LastName             string    `json:"lastName"`
	Contact              string    `json:"contact"`
	AcademicYear         uint8     `json:"academicYear"`
	College              *string   `json:"college"`
	SalesPersonFirstName *string   `json:"salesPersonFirstName"`
	SalesPersonLastName  *string   `json:"salesPersonLastName"`

	// Waiting list related details.
	CompanyBranch          *string `json:"companyBranch"`
	CompanyRequirement     *string `json:"companyRequirement"`
	CompanyRequirementCode *string `json:"companyRequirementCode"`
	Course                 *string `json:"course"`
	Batch                  *string `json:"batch"`
}
