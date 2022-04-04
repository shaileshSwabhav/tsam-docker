package blog

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// BlogReaction contains add update fields required for blog reaction.
type BlogReaction struct {
	general.TenantBase

	// Related table IDs.
	ReactorID *uuid.UUID `json:"reactorID" gorm:"type:varchar(36)"`
	BlogID 	*uuid.UUID `json:"blogID" gorm:"type:varchar(36)"`
	ReplyID 	*uuid.UUID `json:"replyID" gorm:"type:varchar(36)"`

	// Flags.
	IsClap  bool `json:"isClap"`
}

// Validate validates compulsary fields of Blog Reaction.
func (reaction *BlogReaction) Validate() error {

	// Either blog ID or reply ID must be present.
	if ((reaction.BlogID == nil && reaction.ReplyID == nil) || (reaction.BlogID != nil && reaction.ReplyID != nil)){
		return errors.NewValidationError("Either blog ID or reply ID must be present")
	}

	// Blog ID.
	if reaction.BlogID != nil && !util.IsUUIDValid(*reaction.BlogID) {
		return errors.NewValidationError("Blog ID must be a proper uuid")
	}

	// Reply ID.
	if reaction.ReplyID != nil && !util.IsUUIDValid(*reaction.ReplyID) {
		return errors.NewValidationError("Reply ID must be a proper uuid")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// BlogReactionDTO contains all fields required for blog reaction.
type BlogReactionDTO struct {
	general.TenantBase

	// Related table IDs.
	BlogID 	*uuid.UUID `json:"blogID"`
	ReplyID 	*uuid.UUID `json:"replyID"`
	Reactor   list.Credential `json:"reactor" gorm:"foreignkey:ReactorID"`
	ReactorID *uuid.UUID      `json:"-"`

	// Flags.
	IsClap  bool `json:"isClap"`
}

// TableName defines table name of the struct.
func (*BlogReactionDTO) TableName() string {
	return "blog_reactions"
}

