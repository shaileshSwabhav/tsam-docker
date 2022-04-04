package batch

import (
	uuid "github.com/satori/go.uuid"
	fct "github.com/techlabs/swabhav/tsam/models/faculty"
)

type BatchSessionFeedbackDTO struct {
	Talent           TalentDTO                    `json:"talent"`
	SessionFeedbacks []TalentBatchSessionFeedback `json:"sessionFeedbacks"`
}

// TalentBatchSessionFeedbackDTO will consist of faculty and their feedbacks given by talent for specified batch-session
type TalentBatchSessionFeedbackDTO struct {
	Faculty          fct.Faculty                  `json:"faculty"`
	SessionFeedbacks []TalentBatchSessionFeedback `json:"sessionFeedbacks"`
	Feedbacks        []BatchSessionFeedbackDTO    `json:"feedbacks"`
}

// FacultyTalentBatchSessionFeedbackDTO will consist of talent and their feedbacks for specified batch-session
type FacultyTalentBatchSessionFeedbackDTO struct {
	Talent           TalentDTO                           `json:"talent"`
	Faculty          fct.Faculty                         `json:"faculty"`
	SessionFeedbacks []FacultyTalentBatchSessionFeedback `json:"sessionFeedbacks"`
}

// SingleTalentBatchFeedbackDTO will consist of batch session and its related feedbacks.
type SingleTalentBatchFeedbackDTO struct {
	BatchSessionID   uuid.UUID     			  `json:"batchSessionID"`
	Date              string              `json:"date" gorm:"date"`
	SessionFeedbacks []TalentBatchSessionFeedback `json:"sessionFeedbacks"`
}


