package blog

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// BlogReply contains add update fields required for blog reply.
type BlogReply struct {
	general.TenantBase

	// Related table IDs.
	ReplierID *uuid.UUID `json:"replierID" gorm:"type:varchar(36)"`
	BlogID 	*uuid.UUID `json:"blogID" gorm:"type:varchar(36)"`
	ReplyID 	*uuid.UUID `json:"replyID" gorm:"type:varchar(36)"`

	// Other fields.
	Reply       string  `json:"reply" gorm:"type:varchar(2000)"`

	// Flags.
	IsVerified  bool `json:"isVerified"`
}

// Validate validates compulsary fields of Blog Reply.
func (reply *BlogReply) Validate() error {

	// Check if reply is blank or not.
	if util.IsEmpty(reply.Reply) {
		return errors.NewValidationError("Reply must be specified")
	}

	// Reply maximum characters.
	if len(reply.Reply) > 2000 {
		return errors.NewValidationError("Reply can have maximum 2000 characters")
	}

	// Either blog ID or reply ID must be present.
	if ((reply.BlogID == nil && reply.ReplyID == nil) || (reply.BlogID != nil && reply.ReplyID != nil)){
		return errors.NewValidationError("Either blog ID or reply ID must be present")
	}

	// Blog ID.
	if reply.BlogID != nil && !util.IsUUIDValid(*reply.BlogID) {
		return errors.NewValidationError("Blog ID must be a proper uuid")
	}

	// Reply ID.
	if reply.ReplyID != nil && !util.IsUUIDValid(*reply.ReplyID) {
		return errors.NewValidationError("Reply ID must be a proper uuid")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// BlogReplyDTO contains all fields required for blog reply.
type BlogReplyDTO struct {
	general.TenantBase

	// Related table IDs.
	BlogID 	*uuid.UUID `json:"blogID"`
	ReplyID 	*uuid.UUID `json:"replyID"`
	Replier   list.Credential `json:"replier" gorm:"foreignkey:ReplierID"`
	ReplierID *uuid.UUID      `json:"-"`

	// Other fields.
	Reply       string  `json:"reply"`

	// Mutiple field.
	Replies []BlogReplyDTO  `json:"replies" gorm:"foreignkey:ReplyID"`
	Reactions []BlogReactionDTO  `json:"reactions" gorm:"foreignkey:ReplyID"`

	// Flags.
	IsVerified  bool `json:"isVerified"`
}

// TableName defines table name of the struct.
func (*BlogReplyDTO) TableName() string {
	return "blog_replies"
}

