package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
)

// BatchSessionTalent will contain talent details for every batch-session.
type BatchSessionTalent struct {
	general.TenantBase
	BatchID        uuid.UUID `json:"batchID" gorm:"type:varchar(36)"`
	TalentID       uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	IsPresent      *bool     `json:"isPresent"`
	AttendedDate   string    `json:"attendedDate" gorm:"type:date"`
	AverageRating  *float32  `json:"averageRating" gorm:"type:decimal(6,2)"`
	BatchSessionID uuid.UUID `json:"batchSessionID" gorm:"type:varchar(36)"`
	// BatchTopicID  uuid.UUID `json:"batchTopicID" gorm:"type:varchar(36)"`
}

// Validate will check if all fields are valid.
func (session *BatchSessionTalent) Validate() error {

	// if !util.IsUUIDValid(session.BatchID) {
	// 	return errors.NewValidationError("Invalid batch ID.")
	// }

	// // if !util.IsUUIDValid(session.BatchSessionID) {
	// // 	return errors.NewValidationError("Invalid batch session ID.")
	// // }

	// if !util.IsUUIDValid(session.BatchSessionID) {
	// 	return errors.NewValidationError("Invalid batch session ID.")
	// }

	// if !util.IsUUIDValid(session.TalentID) {
	// 	return errors.NewValidationError("Invalid talent ID.")
	// }

	if session.IsPresent == nil {
		return errors.NewValidationError("Whether student is present or not must be specified.")
	}

	return nil
}

// BatchSessionTalentDTO will contain detailed information of all fields in batch_sessions_talent table.
type BatchSessionTalentDTO struct {
	ID              uuid.UUID    `json:"id"`
	DeletedAt       *time.Time   `json:"-"`
	BatchID         uuid.UUID    `json:"-"`
	BatchSessionID  uuid.UUID    `json:"-"`
	TalentID        uuid.UUID    `json:"-"`
	Batch           *list.Batch  `json:"batch" gorm:"foreignkey:BatchID"`
	BatchSession    *Session     `json:"batchSession" gorm:"foreignkey:BatchSessionID"`
	Talent          *list.Talent `json:"talent" gorm:"foreignkey:TalentID"`
	IsPresent       *bool        `json:"isPresent"`
	AttendedDate    *string      `json:"attendedDate"`
	AverageRating   *float32     `json:"averageRating"`
	IsFeedbackGiven bool         `json:"isFeedbackGiven"`
	// BatchTopicID uuid.UUID   `json:"-"`
	// BatchTopic      *list.BatchTopic `json:"batchTopic" gorm:"foreignkey:BatchTopicID"`
}

// TableName defines table name of the struct.
func (*BatchSessionTalentDTO) TableName() string {
	return "batch_session_talents"
}

type SessionTotalHours struct {
	TotalHours float32 `json:"totalHours"`
}

type TopicTotalHours struct {
	TotalHours float32 `json:"totalHours"`
}

type TalentDetailsDTO struct {
	Session         `json:"batchSession" gorm:"foreignkey:ID"`
	IsFeedbackGiven bool `json:"isFeedbackGiven"`
	IsPresent       bool `json:"isPresent"`
	// list.BatchTopic `json:"batchTopic" gorm:"foreignkey:ID"`
	// BatchSessionTalent BatchSessionsTalentDTO `json:"batchSessionsTal"`
}

// AverageRatingBatch will get the avergae rating for a talent.
type AverageRatingBatch struct {
	ID              uuid.UUID    `json:"id"`
	DeletedAt       *time.Time   `json:"-"`
	AverageRating  *float32  `json:"averageRating"`
	MaxScore  *float32  `json:"maxScore"`
}

// TableName defines table name of the struct.
func (*AverageRatingBatch) TableName() string {
	return "faculty_talent_batch_session_feedback"
}
