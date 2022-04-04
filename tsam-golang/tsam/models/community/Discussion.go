package community

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Discussion contains discussion details in db
type Discussion struct {
	general.TenantBase
	Question    string     `json:"question" example:"Question string" gorm:"type:varchar(1000)"`
	Description string     `json:"description" example:"Question Description" gorm:"type:varchar(2000)"`
	IsClosed    *bool      `json:"isClosed" example:"true" gorm:"type:tinyint(1)"`
	IsApproved  *bool      `json:"isApproved" example:"true" gorm:"type:tinyint(1);DEFAULT:false"`
	ChannelID   uuid.UUID  `json:"channelID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	AuthorID    uuid.UUID  `json:"authorID" example:"fga09865-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	BestReplyID *uuid.UUID `json:"bestReplyID" example:"fgb09865-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
}

// DiscussionDTO contains data transfer object of discussion
type DiscussionDTO struct {
	general.BaseDTO
	Question    string     `json:"question"`
	Description string     `json:"description"`
	IsClosed    bool       `json:"isClosed"`
	IsApproved  bool       `json:"isApproved"`
	Author      Talent     `json:"author"`
	Channel     Channel    `json:"channel"`
	BestReply   *Reply     `json:"bestReply,omitempty"`
	ChannelID   uuid.UUID  `json:"-"`
	AuthorID    uuid.UUID  `json:"-"`
	BestReplyID *uuid.UUID `json:"-"`
}

// TableName will refer it to timesheets.
func (*Discussion) TableName() string {
	return "community_discussions"
}

// TableName will refer it to timesheets.
func (*DiscussionDTO) TableName() string {
	return "community_discussions"
}

// Validate Check weather Discussion Valid or Not
func (discussion *Discussion) Validate() error {
	if util.IsEmpty(discussion.Question) {
		return errors.NewValidationError("Question must be specified")
	}

	if util.IsEmpty(discussion.Description) {
		return errors.NewValidationError("Description must be specified")
	}

	if !util.IsUUIDValid(discussion.ChannelID) {
		return errors.NewValidationError("Channel must be specified")
	}

	if !util.IsUUIDValid(discussion.AuthorID) {
		return errors.NewValidationError("Author must be specified")
	}
	return nil
}
