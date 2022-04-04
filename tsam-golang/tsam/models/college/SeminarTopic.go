package college

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Topic defines the fields of a topic required for adding and updating the model.
type Topic struct {
	general.TenantBase

	// Related table IDs.
	SpeakerID *uuid.UUID `json:"speakerID" gorm:"type:varchar(36)"`
	SeminarID uuid.UUID  `json:"seminarID" gorm:"type:varchar(36)"`

	// Other fields.
	TopicName   string  `json:"topicName" example:"Ravi" gorm:"type:varchar(100)"`
	Date        string  `json:"date" gorm:"type:date"`
	FromTime    string  `json:"fromTime" gorm:"type:time"`
	ToTime      string  `json:"toTime" gorm:"type:time"`
	Description *string `json:"description" gorm:"type:varchar(500)"`
}

// TableName defines table name of the struct.
func (*Topic) TableName() string {
	return "seminar_topics"
}

// Validate validates the fields of topic.
func (topic *Topic) Validate() error {

	// Check if topic name is blank or not.
	if util.IsEmpty(topic.TopicName) {
		return errors.NewValidationError("Topic name must be specified")
	}

	// Topic name maximum characters.
	if len(topic.TopicName) > 50 {
		return errors.NewValidationError("Topic name can have maximum 100 characters")
	}

	// Topic date.
	if util.IsEmpty(topic.Date) {
		return errors.NewValidationError("Topic date must be specified")
	}

	// Topic from time.
	if util.IsEmpty(topic.FromTime) {
		return errors.NewValidationError("Topic from time must be specified")
	}

	// Topic to time.
	if util.IsEmpty(topic.ToTime) {
		return errors.NewValidationError("Topic to time must be specified")
	}

	// Description.
	if topic.Description != nil && len(*topic.Description) > 500 {
		return errors.NewValidationError("Topic description can have maximum 500 characters")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// TopicDTO defines the fields of a topic for getting the model.
type TopicDTO struct {
	general.TenantBase

	// Related tables.
	SeminarID uuid.UUID  `json:"seminarID" gorm:"type:varchar(36)"`
	Speaker   *Speaker   `json:"speaker" gorm:"foreignkey:SpeakerID"`
	SpeakerID *uuid.UUID `json:"-"`

	// Other fields.
	TopicName   string  `json:"topicName" example:"Ravi" gorm:"type:varchar(100)"`
	Date        string  `json:"date" gorm:"type:date"`
	FromTime    string  `json:"fromTime" gorm:"type:time"`
	ToTime      string  `json:"toTime" gorm:"type:time"`
	Description *string `json:"description" gorm:"type:varchar(500)"`
}

// TableName defines table name of the struct.
func (*TopicDTO) TableName() string {
	return "seminar_topics"
}
