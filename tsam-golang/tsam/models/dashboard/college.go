package dashboard

// CollegeDashboard contains all details about the college dashboard.
type CollegeDashboard struct {
	TotalColleges       uint                `json:"totalColleges"`
	TotalActiveColleges uint                `json:"totalActiveColleges"`
	CampusDetails       CampusDashboardDTO  `json:"campusDetails"`
	SeminarDetails      SeminarDashboardDTO `json:"seminarDetails"`
}

// CampusDashboardDTO contains all campus details.
type CampusDashboardDTO struct {
	TotalCampusDrives        uint `json:"totalCampusDrives"`
	AllTimeRegisteredTalents uint `json:"allTimeRegisteredTalents"`
	Ongoing                  uint `json:"ongoing"`
	Upcoming                 uint `json:"upcoming"`
	// TalentsSelected          uint `json:"studentsSelected"`
}

// SeminarDashboardDTO contains details related to seminar.
type SeminarDashboardDTO struct {
	TotalSeminars            uint `json:"totalSeminar"`
	AllTimeRegisteredTalents uint `json:"allTimeRegisteredTalents"`
	Ongoing                  uint `json:"ongoing"`
	Upcoming                 uint `json:"upcoming"`
}
