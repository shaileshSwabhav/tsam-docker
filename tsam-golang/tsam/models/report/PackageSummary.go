package report

import "github.com/techlabs/swabhav/tsam/models/general"

// PackageSummary will return count for different packages and experiences.
type PackageSummary struct {
	Experience         string `json:"experience"`
	LessThanThree      uint   `json:"lessThanThree"`
	ThreeToFive        uint   `json:"threeToFive"`
	FiveToTen          uint   `json:"fiveToTen"`
	TenToFifteen       uint   `json:"tenToFifteen"`
	GreaterThanFifteen uint   `json:"greaterThanFifteen"`
}

type TechnologyPackageSummary struct {
	Technology   *general.Technology `json:"technology"`
	TechLanguage string              `json:"techLanguage"`
	TotalCount   uint                `json:"totalCount"`
}

type ExperienceTechnologySummary struct {
	Experience        string                     `json:"experience"`
	TechnologySummary []TechnologyPackageSummary `json:"technologySummary"`
}
