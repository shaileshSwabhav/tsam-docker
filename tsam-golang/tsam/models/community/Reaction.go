package community

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/event"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/notification"
	"github.com/techlabs/swabhav/tsam/util"
)

// Reaction consists of different reactions of user
type Reaction struct {
	general.TenantBase
	DiscussionID *uuid.UUID `json:"discussionID" gorm:"type:varchar(36)"`
	ReplyID      *uuid.UUID `json:"replyID" gorm:"type:varchar(36)"`
	IsLiked      *bool      `json:"isLiked" gorm:"type:tinyint(1)"`
}

// TableName will refer it to timesheets.
func (*Reaction) TableName() string {
	return "community_reactions"
}

// AfterCreate is testing
func (r *Reaction) AfterCreate(db *gorm.DB) error {
	nr := notification.NewReaction(r.TenantID, r.CreatedBy, r.ID, r.DiscussionID, r.ReplyID)
	return event.FireEvent(event.ReactionAdded, nr)
}

// ReactionDTO consists of fields for get.
type ReactionDTO struct {
	general.BaseDTO
	DiscussionID *uuid.UUID `json:"discussionID"`
	ReplyID      *uuid.UUID `json:"replyID"`
	IsLiked      *bool      `json:"isLiked"`
	ReactorID    uuid.UUID  `json:"reactorID"`
}

// TableName will refer it to timesheets.
func (*ReactionDTO) TableName() string {
	return "community_reactions"
}

// Validate validates the reaction.
func (r *Reaction) Validate() error {
	if r.DiscussionID == nil && r.ReplyID == nil {
		return errors.NewValidationError("reply or discussion id must be specified")
	}
	if r.DiscussionID != nil {
		if !util.IsUUIDValid(*r.DiscussionID) {
			return errors.NewValidationError("Improper discussion id")
		}
	}
	if r.ReplyID != nil {
		if !util.IsUUIDValid(*r.ReplyID) {
			return errors.NewValidationError("Improper reply id")
		}
	}
	return nil
}
