package community

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// NotificationDTO contains fields to sent as JSON on request.
type NotificationDTO struct {
	ID               uuid.UUID         `json:"id,omitempty"`
	Discussion       *discussionDTO    `json:"discussion,omitempty"`
	Reply            *replyDTO         `json:"reply,omitempty"`
	NotificationType *NotificationType `json:"notificationType,omitempty"`
	NotifierTalent   *Talent           `json:"notifierTalent,omitempty"`
	NotifiedTalent   *Talent           `json:"notifiedTalent,omitempty"`
	SeenTime         *time.Time        `json:"seenTime,omitempty"`
	IsSeen           *bool             `json:"isSeen,omitempty"`
	Added            *string           `json:"added,omitempty"`
	NotificationText *string           `json:"notificationText,omitempty"`
}

type discussionDTO struct {
	DiscussionID uuid.UUID `json:"discussionID,omitempty"`
	Question     *string   `json:"question,omitempty"`
}

type replyDTO struct {
	ReplyID uuid.UUID `json:"replyID,omitempty"`
	Replier *Talent   `json:"replier,omitempty"`
	Reply   string    `json:"reply,omitempty"`
}

// NewNotificationDTO is used to create an instance of NotificationDTO.
func NewNotificationDTO() *NotificationDTO {
	return &NotificationDTO{
		Discussion: &discussionDTO{},
		Reply: &replyDTO{
			Replier: &Talent{},
		},
		NotificationType: &NotificationType{},
	}
}
