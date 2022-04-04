package dashboard

// CourseDashboard contains all details about the course-batch dashboard.
type CourseDashboard struct {
	TotalCourses uint          `json:"totalCourses"`
	CourseGroups []CourseGroup `json:"courseGroups"`
}

// CourseGroup contains fields of course as a group.
type CourseGroup struct {
	GroupName       string               `json:"groupName"`
	AllBatches      standardBatchDetails `json:"allBatches"`
	LiveBatches     uint                 `json:"liveBatches"`
	UpcomingBatches upcomingBatchDetails `json:"upcomingBatches"`
	Sessions        sessionDetails       `json:"sessions"`
	CoursesData     []*CourseData        `json:"coursesData"`
}

// CourseData defines fields for necessary course details.
type CourseData struct {
	CourseName      string               `json:"courseName"`
	CourseLevel     string               `json:"courseLevel"`
	AllBatches      standardBatchDetails `json:"allBatches"`
	UpcomingBatches upcomingBatchDetails `json:"upcomingBatches"`
	LiveBatches     uint                 `json:"liveBatches"`
	Sessions        sessionDetails       `json:"sessions"`
}

type upcomingBatchDetails struct {
	standardBatchDetails
	RequiredTalents   uint `json:"requiredTalents"`
	TotalTalentIntake uint `json:"totalTalentIntake"`
}

type sessionDetails struct {
	TotalSessions     uint `json:"totalSessions"`
	RemainingSessions uint `json:"remainingSessions"`
}

type standardBatchDetails struct {
	TotalBatches       uint `json:"totalBatches"`
	TotalTalentsJoined uint `json:"totalTalentsJoined"`
}
