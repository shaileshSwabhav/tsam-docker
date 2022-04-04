package talent

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
)

// TotalCount contains information about total entries in database.
type TotalCount struct {
	TotalCount int
}

// WaitingListCompanyBranchDTO is used for getting company branch related waiting list report.
type WaitingListCompanyBranchDTO struct {
	CompanyBranchID uuid.UUID          `json:"-"`
	CompanyBranch   list.CompanyBranch `json:"companyBranch" gorm:"foreignkey:CompanyBranchID"`
	TalentCount     *uint16            `json:"talentCount"`
	EnquiryCount    *uint16            `json:"enquiryCount"`
}

// WaitingListCourseDTO is used for getting course related waiting list report.
type WaitingListCourseDTO struct {
	CourseID     uuid.UUID   `json:"-"`
	Course       list.Course `json:"course" gorm:"foreignkey:CourseID"`
	TalentCount  *uint16     `json:"talentCount"`
	EnquiryCount *uint16     `json:"enquiryCount"`
}

// WaitingListRequirementDTO is used for getting requirement related waiting list report.
type WaitingListRequirementDTO struct {
	CompanyRequirementID uuid.UUID        `json:"-"`
	CompanyRequirement   list.Requirement `json:"requirement" gorm:"foreignkey:CompanyRequirementID"`
	TalentCount          *uint16          `json:"talentCount"`
	EnquiryCount         *uint16          `json:"enquiryCount"`
}

// WaitingListBatchDTO is used for getting batch related waiting list report.
type WaitingListBatchDTO struct {
	BatchID      uuid.UUID  `json:"-"`
	Batch        list.Batch `json:"batch" gorm:"foreignkey:BatchID"`
	TalentCount  *uint16    `json:"talentCount"`
	EnquiryCount *uint16    `json:"enquiryCount"`
}

// WaitingListTechnologyDTO is used for getting technology related waiting list report.
type WaitingListTechnologyDTO struct {
	TechnologyID uuid.UUID          `json:"-"`
	Technology   general.Technology `json:"technology" gorm:"foreignkey:TechnologyID"`
	TalentCount  *uint16            `json:"talentCount"`
	EnquiryCount *uint16            `json:"enquiryCount"`
}
