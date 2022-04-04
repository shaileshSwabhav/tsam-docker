package dashboard

// TechnologyDashboard contains all details about the technology dashboard.
type TechnologyDashboard struct {
	TotalCount       uint              `json:"totalCount"`
	TechnologiesData []*TechnologyData `json:"technologiesData"`
	technologyTotalCount
}

// TechnologyData consists of data required for the dashboard.
type TechnologyData struct {
	Name string `json:"name"`
	technologyTotalCount
}

// TechnologyData consists of data required for the dashboard.
type technologyTotalCount struct {
	TotalFreshers    uint `json:"totalFreshers"`
	TotalExperienced uint `json:"totalExperienced"`
	TotalTalents     uint `json:"totalTalents"`
}
