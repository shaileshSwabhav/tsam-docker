package faculty

import (
	uuid "github.com/satori/go.uuid"
)

type FacultyBatch struct {
	OngoingBatches       int    `json:"ongoingBatches"`
	UpcomingBatches      int    `json:"upcomingBatches"`
	FinishedBatches      int    `json:"finishedBatches"`
	TotalStudents        int    `json:"totalStudents"`
	CompletedTrainingHrs string `json:"completedTrainingHrs"`
}

type OngoingBatchDetails struct {
	BatchID        uuid.UUID `json:"batchID"`
	CourseName     string    `json:"courseName"`
	BatchName      string    `json:"batchName"`
	TotalSession   uint      `json:"totalSession"`
	PendingSession uint      `json:"pendingSession"`
	TotalStudents  uint      `json:"totalStudents"`
}

type PiechartData struct {
	ProjectName string  `json:"projectName"`
	TotalCount  int     `json:"totalCount"`
	Hours       float64 `json:"hours"`
}

type BarchartData struct {
	TotalStudents  int `json:"totalStudents"`
	Fresher        int `json:"fresher"`
	Professional   int `json:"professional"`
	StudentsPlaced int `json:"studentsPlaced"`
}

type Feedback struct {
	TalentFeedback []TalentFeedbackScore `json:"talentFeedback"`
	Keywords       []KeywordName         `json:"keywords"`
}

// TalentFeedbackScore.
type TalentFeedbackScore struct {
	TalentID        uuid.UUID       `json:"talentID"`
	FirstName       string          `json:"firstName"`
	LastName        string          `json:"lastName"`
	PersonalityType *string         `json:"personalityType"`
	TalentType      *uint8          `json:"talentType"`
	BatchName       string          `json:"batchName"`
	BatchID         uuid.UUID       `json:"batchID"`
	Score           float64         `json:"score"`
	InterviewRating float64         `json:"interviewRating"`
	SessionFeedback []FeedbackScore `json:"sessionFeedback"`
	// KeywordNames    []KeywordName   `json:"keywordNames"`
}

// KeywordName will contain keyword names.
type KeywordName struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// FeedbackScore will contain keyword name and its score.
type FeedbackScore struct {
	Keyword      string  `json:"keyword"`
	KeywordScore float64 `json:"keywordScore"`
}

type WeeklyAvgRating struct {
	Rating   float64         `json:"rating"`
	Feedback []FeedbackScore `json:"feedbackScore,omitempty"`
}
