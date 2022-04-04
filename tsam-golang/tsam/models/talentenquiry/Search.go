package talentenquiry

import (
	uuid "github.com/satori/go.uuid"
)

// Search Contains search fields for talent enquiry.
type Search struct {
	FirstName                 *string     `json:"firstName"`
	LastName                  *string     `json:"lastName"`
	Email                     *string     `json:"email"`
	College                   *string     `josn:"college"`
	IsExperience              *bool       `json:"isExperience"`
	IsMastersAbroad           *bool       `json:"isMastersAbroad"`
	Passout                   *int        `json:"passout"`
	TotalExperience           *int8       `json:"totalExperience"`
	MinimumExperience         *uint8      `json:"minimumExperience"`
	MaximumExperience         *uint8      `json:"maximumExperience"`
	CallRecordPurposeID       *uuid.UUID  `json:"callRecordPurposeID"`
	CallRecordOutcomeID       *uuid.UUID  `json:"callRecordOutcomeID"`
	AcademicYears             []*uint8    `json:"academicYears"`
	Qualifications            []uuid.UUID `json:"qualifications"`
	Designations              []uuid.UUID `json:"designations"`
	SalesPersonIDs            []uuid.UUID `json:"salesPersonIDs"`
	Technologies              []uuid.UUID `json:"technologies"`
	ExperienceTechnologies    []uuid.UUID `json:"experienceTechnologies"`
	SearchAllEnquiries        *bool       `json:"searchAllEnquiries"`
	City                      *string     `json:"city"`
	CountryID                 *uuid.UUID  `json:"countryID"`
	YearOfMS                  *uint       `json:"yearOfMS"`
	WaitingFor                *string     `json:"waitingFor"`
	WaitingForCompanyBranchID *string     `json:"waitingForCompanyBranchID"`
	WaitingForRequirementID   *string     `json:"waitingForRequirementID"`
	WaitingForCourseID        *string     `json:"waitingForCourseID"`
	WaitingForBatchID         *string     `json:"waitingForBatchID"`
	WaitingForIsActive        *bool       `json:"waitingForIsActive"`
	WaitingForFromDate        *string     `json:"waitingForFromDate"`
	WaitingForToDate          *string     `json:"waitingForToDate"`
	EnquiryFromDate           *string     `json:"enquiryFromDate"`
	EnquiryToDate             *string     `json:"enquiryToDate"`
	IsLastThirtyDays          *bool       `json:"isLastThirtyDays"`
	EnquirySource             *uuid.UUID  `json:"enquirySource"`
}
