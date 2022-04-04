package list

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// ProgrammingQuestion for programming questions list details.
type ProgrammingQuestion struct {
	ID        uuid.UUID  `json:"id" gorm:"type:varchar(36)"`
	DeletedAt *time.Time `json:"-"`
	Label     string     `json:"label"`
	Level     uint8      `json:"level"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestion) TableName() string {
	return "programming_questions"
}
