package batch

import (
	fct "github.com/techlabs/swabhav/tsam/models/faculty"
)

// change struct names

// batch feedback

type BatchFeedbackDTO struct {
	Talent         TalentDTO        `json:"talent"`
	BatchFeedbacks []TalentFeedback `json:"batchFeedbacks"`
}

// FacultyTalentBatchFeedbackDTO will consist of talent and their feedbacks for specified batch
type FacultyTalentBatchFeedbackDTO struct {
	Faculty        fct.Faculty             `json:"faculty"`
	Talent         TalentDTO               `json:"talent"`
	BatchFeedbacks []FacultyTalentFeedback `json:"batchFeedbacks"`
}

// TalentBatchFeedbackDTO will consist of faculty and their feedbacks given by talent for specified batch
type TalentBatchFeedbackDTO struct {
	Faculty        fct.Faculty        `json:"faculty"`
	BatchFeedbacks []TalentFeedback   `json:"batchFeedbacks"`
	Feedbacks      []BatchFeedbackDTO `json:"feedbacks"`
}
