package talent

import uuid "github.com/satori/go.uuid"

// ProfessionalSummaryReport is used for getting company name related professional summay report.
// FirstCount is for 1-2 years in experience.
// SecondCount is for 2-5 years in experience.
// ThirdCount is for 5-7 years in experience.
// FourthCount is for above 7 years in experience.
// All data is for only currently working talents.
type ProfessionalSummaryReport struct {
	Company     string  `json:"company"`
	FirstCount  *uint16 `json:"firstCount"`
	SecondCount *uint16 `json:"secondCount"`
	ThirdCount  *uint16 `json:"thirdCount"`
	FourthCount *uint16 `json:"fourthCount"`
}

// ProfessionalSummaryReportCounts is used for getting company name related professional summay report.
// FirstCountTotal is the total of FirstCount from ProfessionalSummaryReport.
// SecondCountTotal is the total of SecondCount from ProfessionalSummaryReport.
// ThirdCountTotal is the total of ThirdCount from ProfessionalSummaryReport.
// FourthCountTotal is the total of FourthCount from ProfessionalSummaryReport.
// All data is for only currently working talents.
type ProfessionalSummaryReportCounts struct {
	TotalCount       int   `json:"totalCount"`
	FirstCountTotal  *uint `json:"firstCountTotal"`
	SecondCountTotal *uint `json:"secondCountTotal"`
	ThirdCountTotal  *uint `json:"thirdCountTotal"`
	FourthCountTotal *uint `json:"fourthCountTotal"`
}

// TechnologyTalent provides model for getting count of talents for each technology and company.
type CompanyTechnologyTalent struct {
	Company     string    `json:"company"`
	TechID      uuid.UUID `json:"techID"`
	TechName    string    `json:"techName"`
	TalentCount uint      `json:"talentCount"`
}
