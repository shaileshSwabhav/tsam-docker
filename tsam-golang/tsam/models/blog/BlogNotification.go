package blog

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
)

type BlogNotification struct {
	general.TenantBase
	BlogID         *uuid.UUID `json:"blogID" gorm:"type:varchar(36)"`
	BlogTopicID    *uuid.UUID `json:"blogTopicID" gorm:"type:varchar(36)"`
	BlogReactionID *uuid.UUID `json:"blogReactionID" gorm:"type:varchar(36)"`
	BlogReplyID    *uuid.UUID `json:"blogReplyID" gorm:"type:varchar(36)"`
}
