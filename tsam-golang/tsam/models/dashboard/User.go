package dashboard

// SalesPersonDashboard contains fields of salesPerson to be displayed on the dashboard.
type SalesPersonDashboard struct {
	TotalSalesPeople   uint `json:"totalSalesPeople"`
	EnquiriesAssigned  uint `json:"enquiriesAssigned"`
	EnquiriesConverted uint `json:"enquiriesConverted"`
	// EnquiriesNotAssigned  uint `json:"enquiriesNotAssigned"`
	EnquiriesNotHandled           uint `json:"enquiriesNotHandled"`
	TalentsApproached             uint `json:"talentsApproached"`
	JoinedBatches                 uint `json:"joinedBatches"`
	TrainingEnquiries             uint `json:"trainingEnquiries"`
	PlacementEnquiries            uint `json:"placementEnquiries"`
	TrainingAndPlacementEnquiries uint `json:"trainingAndPlacementEnquiries"`
	CampusDrivesCompleted         uint `json:"campusDrivesCompleted"`
	SeminarsCompleted             uint `json:"seminarsCompleted"`
}

// AdminDashboard defines all fields in admin dashboard.
type AdminDashboard struct {
	TotalCourses      uint `json:"totalCourses"`
	TotalTechnologies uint `json:"totalTechnologies"`
}

// Course id and name, student => talent, B2C => total count,  totalTalentsJoined & TotalIntake
// faculty talents Joined ( no first letter caps in Course)
