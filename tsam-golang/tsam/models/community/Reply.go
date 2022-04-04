package community

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Reply contains the fields of a reply.
type Reply struct {
	general.TenantBase
	DiscussionID uuid.UUID  `json:"discussionID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	ReplyID      *uuid.UUID `json:"replyID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	Reply        string     `json:"reply" example:"Answer Against Question" gorm:"type:varchar(4000)"`
	IsBestReply  *bool      `json:"isBestReply" example:"true" gorm:"type:tinyint(1)"`
	IsVerified   *bool      `json:"isVerified" example:"true" gorm:"type:tinyint(1)"`
	ReplierID    uuid.UUID  `json:"replierID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	// Comments     *[]Reply   `json:"comments,omitempty" gorm:"foreignkey:ReplyID;association_autocreate:false;association_autoupdate:false"`
}

// ReplyDTO contains fields to be shown in a get operation
type ReplyDTO struct {
	general.BaseDTO
	Reply        string      `json:"reply" example:"Answer Against Question"`
	IsBestReply  *bool       `json:"isBestReply" example:"true"`
	IsVerified   *bool       `json:"isVerified" example:"true"`
	DiscussionID uuid.UUID   `json:"discussionID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	ReplyID      *uuid.UUID  `json:"replyID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	Comments     *[]ReplyDTO `json:"comments,omitempty"`
	Replier      Credential  `json:"replier"`
	ReplierID    uuid.UUID   `json:"-" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	// change to uint if int is not needed.
	NumberOfLikes    uint16 `json:"numberOfLikes" example:"20"`
	NumberOfDislikes uint16 `json:"numberOfDislikes" example:"20"`
	IsLiked          *bool  `json:"isLiked"`
}

// TableName will refer it to timesheets.
func (*Reply) TableName() string {
	return "community_replies"
}

// TableName will refer it to timesheets.
func (*ReplyDTO) TableName() string {
	return "community_replies"
}

// // TobeDecided Contain Discussion, ReplyDTO
// type TobeDecided struct {
// 	DiscussionDTO *Discussion `json:"discussionDTO,omitempty"`
// 	ReplyDTO      []Reply     `json:"replyDTO,omitempty"`
// }

// Validate Return Error On Invalid Reply
func (reply *Reply) Validate() error {
	if util.IsEmpty(reply.Reply) {
		return errors.NewValidationError("Reply must be specified")
	}
	if !util.IsUUIDValid(reply.DiscussionID) {
		return errors.NewValidationError("DiscussionID must be specified")
	}
	if !util.IsUUIDValid(reply.ReplierID) {
		return errors.NewValidationError("ReplierID must be specified")
	}
	return nil
}
