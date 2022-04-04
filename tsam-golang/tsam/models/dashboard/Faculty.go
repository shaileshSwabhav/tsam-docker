package dashboard

// FacultyDashboard includes fields to be displayed in faculty dashboard
type FacultyDashboard struct {
	TotalFaculty     uint                 `json:"totalFaculty"`
	Active           uint                 `json:"active"`
	InActive         uint                 `json:"inActive"`
	PartTime         uint                 `json:"partTime"`
	FullTime         uint                 `json:"fullTime"`
	BatchesCompleted standardBatchDetails `json:"batchesCompleted"`
	LiveBatches      standardBatchDetails `json:"liveBatches"`
	UpcomingBatches  standardBatchDetails `json:"upcomingBatches"`
}
