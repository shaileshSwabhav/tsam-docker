package list

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type ModuleTopic struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`
	TopicName string     `json:"topicName"`
	TotalTime *uint      `json:"totalTime"`
	Order     uint       `json:"order"`
}

// TableName overrides name of the table
func (*ModuleTopic) TableName() string {
	return "module_topics"
}
