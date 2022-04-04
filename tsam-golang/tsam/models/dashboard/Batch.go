package dashboard

import uuid "github.com/satori/go.uuid"

// BatchDashboard contains all details about the batch dashboard.
type BatchDashboard struct {
	Live      batchType `json:"live"`
	Completed batchType `json:"completed"`
	Upcoming  batchType `json:"upcoming"`
	Total     batchType `json:"total"`
}

type batchType struct {
	B2BBatches    uint `json:"b2bBatches"`
	B2CBatches    uint `json:"b2cBatches"`
	TotalStudents uint `json:"totalStudents"`
	AllBatches    uint `json:"allBatches"`
}

// BatchPerformance will contain different paramter counts and talent details.
type BatchPerformance struct {
	TotalBatches      uint                   `json:"totalBatches"`
	Outstanding       uint                   `json:"outstanding"`
	Good              uint                   `json:"good"`
	Average           uint                   `json:"average"`
	KeywordNames      []GroupWiseKeywordName `json:"keywordNames"`
	OutstandingTalent []TalentFeedbackScore  `json:"outstandingTalent"`
	GoodTalent        []TalentFeedbackScore  `json:"goodTalent"`
	AverageTalent     []TalentFeedbackScore  `json:"averageTalent"`
	// KeywordNames      []KeywordName          `json:"keywordNames"`
}

// BatchScore will return scores of different parameters.
type BatchScore struct {
	TotalBatches uint `json:"totalBatches"`
	Outstanding  uint `json:"outstanding"`
	Good         uint `json:"good"`
	Average      uint `json:"average"`
}

// TalentFeedbackScore.
type TalentFeedbackScore struct {
	TalentID        uuid.UUID `json:"talentID"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	PersonalityType *string   `json:"personalityType"`
	TalentType      *uint8    `json:"talentType"`
	BatchName       string    `json:"batchName"`
	BatchID         uuid.UUID `json:"batchID"`
	Score           float64   `json:"score"`
	InterviewRating float64   `json:"interviewRating"`
	// Keywords        []FeedbackKeywords `json:"feedbackKeywords"`
	KeywordNames  []GroupWiseKeywordName `json:"keywordNames"`
	FeedbackGroup []GroupScore           `json:"feedbackKeywords"`
}

// TalentSessionFeedbackScore.
type TalentSessionFeedbackScore struct {
	TalentID        uuid.UUID                `json:"talentID"`
	FirstName       string                   `json:"firstName"`
	LastName        string                   `json:"lastName"`
	PersonalityType *string                  `json:"personalityType"`
	TalentType      *uint8                   `json:"talentType"`
	BatchName       string                   `json:"batchName"`
	BatchID         uuid.UUID                `json:"batchID"`
	KeywordNames    []GroupWiseKeywordName   `json:"keywordNames"`
	SessionFeedback []SessionKeywordFeedback `json:"sessionFeedback"`
	// KeywordNames          []KeywordName            `json:"keywordNames"`
	// Group           []SessionGroupScore      `json:"group"`
}

// SessionKeywordFeedback will contain group wise feedback score and feeling details for the specified session.
type SessionKeywordFeedback struct {
	BatchSessionID uuid.UUID    `json:"batchSessionID"`
	SessionName    string       `json:"sessionName"`
	Order          uint         `json:"order"`
	Score          float64      `json:"score"`
	SessionDate    string       `json:"sessionDate"`
	FeelingName    string       `json:"feelingName"`
	LevelNumber    uint         `json:"levelNumber"`
	Description    string       `json:"description"`
	FeedbackGroup  []GroupScore `json:"feedbackGroup"`
	// Keywords       []FeedbackKeywords  `json:"feedbackKeywords"`
}

// GroupScore will contain group name and scores for all the feedback question of that groups.
type GroupScore struct {
	GroupID       uuid.UUID          `json:"groupID"`
	GroupName     string             `json:"groupName"`
	FeedbackScore []FeedbackKeywords `json:"feedbackScore"`
}

// FeedbackKeywords will contain keyword name and its score.
type FeedbackKeywords struct {
	Keyword      string  `json:"keyword"`
	KeywordScore float64 `json:"keywordScore"`
}

// BatchDetails will contain details of specified batch.
type BatchDetails struct {
	BatchID       string `json:"batchID"`
	BatchName     string `json:"batchName"`
	BatchType     bool   `json:"batchType"`
	TotalStudents uint   `json:"totalStudents"`
}

// AhaMoment will give details of aha-moment feelings and aha-moment responses.
type AhaMoment struct {
	FeelingName       string              `json:"feelingName"`
	LevelNumber       uint                `json:"levelNumber"`
	Description       string              `json:"description"`
	AhaMomentResponse []AhaMomentResponse `json:"ahaMomentResponse"`
}

// AhaMomentResponse will contain question and response detais
type AhaMomentResponse struct {
	Question string `json:"question"`
	Response string `json:"response"`
}

// GroupWiseKeywordName will contain group name and keywords for that question.
type GroupWiseKeywordName struct {
	GroupName string        `json:"groupName"`
	Keywords  []KeywordName `json:"keywords"`
}

// KeywordName will contain keyword names.
type KeywordName struct {
	Name string `json:"name"`
}
