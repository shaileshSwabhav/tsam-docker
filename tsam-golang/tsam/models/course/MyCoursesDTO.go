package course

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/list"
)

//*************************************** MY COURSES DTO ******************************************************

// MyCoursesDTO will contain details of courses taken by talent.
type MyCoursesDTO struct {
	ID                  uuid.UUID    `json:"id"`
	DeletedAt           *time.Time   `json:"-"`
	Name                string       `json:"name"`
	BatchName           string       `json:"batchName"`
	BatchStatus         string       `json:"batchStatus"`
	RequirementID       *uuid.UUID     `json:"requirementID"`
	BatchID             uuid.UUID    `json:"batchID"`
	CourseType          string       `json:"courseType"`
	Description         string       `json:"description"`
	TotalSessions       uint         `json:"totalSessions"`
	CompletedSessions   uint         `json:"completedSessions"`
	TotalCourseTaken    int          `json:"totalCourseTaken"`
	CompletedPercentage float64      `json:"completedPercentage"`
	Logo                *string      `json:"logo"`
	SubscribedDate      string       `json:"subscribedDate"`
	FacultyID           uuid.UUID    `json:"-"`
	Faculty             list.Faculty `json:"faculty"`
	Price               float64      `json:"price"`
	StartDate           string       `json:"startDate"`
	EndDate             string       `json:"endDate"`
}

// TableName defines table name of the struct.
func (*MyCoursesDTO) TableName() string {
	return "courses"
}

//*************************************** COURSE DETAILS ******************************************************

// CourseDetails will contain specific deatils of course.
type CourseDetails struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	Banner           *string `json:"banner"`
	Name             string  `json:"name"`
	DurationInMonths int64   `json:"durationInMonths"`
	Price            float64 `json:"price"`
	SessionsCount    uint    `json:"sessionsCount"`
	Description      string  `json:"description"`
	Brochure         *string `json:"brochure"`
}

// TableName defines table name of the struct.
func (*CourseDetails) TableName() string {
	return "courses"
}

//*************************************** COURSE MINIMIM DETAILS ******************************************************

// CourseMinimumDetails will contain specific deatils of course for student dashboard.
type CourseMinimumDetails struct {
	ID uuid.UUID `json:"id"`

	Logo       *string `json:"logo"`
	Name       string  `json:"name"`
	IsEnrolled uint    `json:"isEnrolled"`
}

// TableName defines table name of the struct.
func (*CourseMinimumDetails) TableName() string {
	return "courses"
}
