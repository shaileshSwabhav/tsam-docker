package list

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type BatchSessionTopic struct {
	ID            uuid.UUID  `json:"id"`
	DeletedAt     *time.Time `json:"-"`
	SubTopicID    uuid.UUID  `json:"-"`
	Order         uint       `json:"order"`
	CompletedDate string     `json:"completedDate"`
	IsCompleted   *bool      `json:"isCompleted"`
	TotalTime     uint       `json:"totalTime"`
}

// TableName overrides name of the table
func (*BatchSessionTopic) TableName() string {
	return "batch_session_topics"
}
