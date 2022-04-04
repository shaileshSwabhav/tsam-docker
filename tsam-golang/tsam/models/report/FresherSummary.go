package report

import "github.com/techlabs/swabhav/tsam/models/general"

// FresherSummary contains count count of rating of talents.
type FresherSummary struct {
	ColumnName       string             `json:"columnName"`
	Academic         general.CommonType `json:"academic"`
	OutstandingCount uint               `json:"outstandingCount"`
	ExcellentCount   uint               `json:"excellentCount"`
	AverageCount     uint               `json:"averageCount"`
	UnrankedCount    uint               `json:"unrankedCount"`

	// Type []AcademicYear `json:"type"`
}

// TechnologyFresherSummary contains count count of rating of talents.
type TechnologyFresherSummary struct {
	Technology     general.Technology `json:"technology"`
	FresherSummary []FresherSummary   `json:"fresherSummary"`

	// Type []AcademicYear `json:"type"`
}

// type AcademicYear struct {
// 	AcademicYear     string `json:"academicYear"`
// 	OutstandingCount uint   `json:"outstandingCount"`
// 	ExcellentCount   uint   `json:"excellentCount"`
// 	AverageCount     uint   `json:"averageCount"`
// }

type TechnologySummary struct {
	TechnologyLangugage string              `json:"technologyLangugage"`
	Technology          *general.Technology `json:"technology"`
	TotalCount          uint                `json:"totalCount"`
	ColumnName          string              `json:"columnName"`
	Academic            *general.CommonType `json:"academic"`
}

type AcademicTechnologySummary struct {
	ColumnName        string              `json:"columnName"`
	Academic          *general.CommonType `json:"academic"`
	TechnologySummary []TechnologySummary `json:"technologySummary"`
}
