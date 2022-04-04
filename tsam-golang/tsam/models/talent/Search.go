package talent

import (
	uuid "github.com/satori/go.uuid"
)

// Search contains fields for giving search criteria to search talents.
type Search struct {
	FirstName                 *string     `json:"firstName"`
	LastName                  *string     `json:"lastName"`
	Email                     *string     `json:"email"`
	College                   *string     `josn:"college"`
	PersonalityType           *string     `json:"personalityType"`
	IsActive                  *bool       `json:"isActive"`
	IsExperience              *bool       `json:"isExperience"`
	IsSwabhavTalent           *bool       `json:"isSwabhavTalent"`
	IsMastersAbroad           *bool       `json:"isMastersAbroad"`
	Passout                   *int        `json:"passout"`
	TalentType                *int8       `json:"talentType"`
	TotalExperience           *int8       `json:"totalExperience"`
	LifetimeValue             *uint       `json:"lifetimeValue"`
	MinimumExperience         *uint8      `json:"minimumExperience"`
	MaximumExperience         *uint8      `json:"maximumExperience"`
	Percentage                *float32    `json:"percentage"`
	CourseID                  *uuid.UUID  `json:"courseID"`
	FacultyID                 *uuid.UUID  `json:"facultyID"`
	BatchID                   *uuid.UUID  `json:"batchID"`
	CallRecordPurposeID       *uuid.UUID  `json:"callRecordPurposeID"`
	CallRecordOutcomeID       *uuid.UUID  `json:"callRecordOutcomeID"`
	NextActionTypeID          *uuid.UUID  `json:"nextActionTypeID"`
	AcademicYears             []*uint8    `json:"academicYears"`
	Qualifications            []uuid.UUID `json:"qualifications"`
	Designations              []uuid.UUID `json:"designations"`
	SalesPersonIDs            []uuid.UUID `json:"salesPersonIDs"`
	Technologies              []uuid.UUID `json:"technologies"`
	ExperienceTechnologies    []uuid.UUID `json:"experienceTechnologies"`
	CompanyName    			  *string     `json:"companyName"`
	SearchAllTalents          *bool       `json:"searchAllTalents"`
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
}
