package talentenquiry

import (
	uuid "github.com/satori/go.uuid"
)

// EnquiryWaitingList is used for getting enquiry and waiting list realted fields for waiting list report.
type EnquiryWaitingList struct {

	// Enquiry related details.
	EnquiryID            uuid.UUID `json:"enquiryID"`
	FirstName            string    `json:"firstName"`
	LastName             string    `json:"lastName"`
	Contact              string    `json:"contact"`
	AcademicYear         *uint8    `json:"academicYear"`
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
