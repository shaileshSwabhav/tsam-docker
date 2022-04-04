package dashboard

// CompanyDashboard contains all details about the company dashboard.
type CompanyDashboard struct {
	TotalCompanies         uint                 `json:"totalCompanies"`
	LiveRequirements       standardRequirements `json:"liveRequirements"`
	FresherRequirements    standardRequirements `json:"fresherRequirements"`
	ExperienceRequirements standardRequirements `json:"experienceRequirements"`
}

// StandardRequirements contains the definition for most requirements.
type standardRequirements struct {
	TotalRequirements    uint `json:"totalRequirements"`
	TotalTalentsRequired uint `json:"totalTalentsRequired"`
	PlacementsInProcess  uint `json:"placementsInProcess"`
	TalentsNeeded        uint `json:"talentsNeeded"`
}
