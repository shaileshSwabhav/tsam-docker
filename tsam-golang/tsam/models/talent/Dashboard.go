package talent

import uuid "github.com/satori/go.uuid"

type PerformanceDetails struct {
	TalentID      uuid.UUID         `json:"talentID"`
	FirstName     string            `json:"firstName"`
	LastName      string            `json:"lastName"`
	AverageRating float64           `json:"averageRating,omitempty"`
	Score         []WeeklyAvgRating `json:"weeklyAvgRating,omitempty"`
	Feedback      []FeedbackScore   `json:"feedbackScore,omitempty"`
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

// ConceptRatingWithAssignment will talent concept ratings for talent.
type ConceptRatingWithAssignment struct {
	TalentID      uuid.UUID         `json:"talentID"`
	FirstName     string            `json:"firstName"`
	LastName      string            `json:"lastName"`
	Assignments      []AssignmentScore   `json:"assignments"`
}

// AssignmentScore will score for each batch topic assignment.
type AssignmentScore struct {
	AssignmentID      string  `json:"assignmentID"`
	Score 			  float64 `json:"score"`
}
