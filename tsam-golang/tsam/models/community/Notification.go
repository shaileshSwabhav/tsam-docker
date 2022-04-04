package community

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
)

// Notification contains fields to send notification to talent.
type Notification struct {
	general.TenantBase
	NotifierID         uuid.UUID  `json:"notifierID,omitempty" example:"aca09865-f5fe-48f0-874d-e72cd4edd9b7" gorm:"type:varchar(36)"`
	SubscriberID       uuid.UUID  `json:"subscriberID,omitempty" example:"aca09865-f5fe-48f0-874d-e72cd4edd9b7" gorm:"type:varchar(36)"`
	DiscussionID       uuid.UUID  `json:"discussionID,omitempty" example:"aba09865-f5fe-48f0-874d-e72cd4edd9b8" gorm:"type:varchar(36)"`
	ReplyID            uuid.UUID  `json:"replyID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	NotificationTypeID uuid.UUID  `json:"notificationTypeID,omitempty" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	SeenTime           *time.Time `json:"seenTime,omitempty"`
	Message            string     `json:"message" gorm:"type:varchar(200)"`
	// seen time should suffice. #niranjan
	IsSeen *bool `json:"isSeen,omitempty" example:"true"`
}

// ValidateNotification validates all important fields in Notification.
func (notification *Notification) ValidateNotification() error {
	if notification.DiscussionID == uuid.Nil {
		return errors.NewValidationError("Discussion id must be specified")
	}
	if notification.ReplyID == uuid.Nil {
		return errors.NewValidationError("Reply id must be specified")
	}
	if notification.NotificationTypeID == uuid.Nil {
		return errors.NewValidationError("Notification type id must be specified")
	}
	if notification.NotifierID == uuid.Nil {
		return errors.NewValidationError("Notifier must be specified")
	}
	if notification.SubscriberID == uuid.Nil {
		return errors.NewValidationError("Who is getting notified should be specified")
	}
	return nil
}
