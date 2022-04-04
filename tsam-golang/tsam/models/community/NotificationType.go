package community

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"

	"github.com/techlabs/swabhav/tsam/models/general"
)

// NotificationType contains differnt types of notifaction names.
type NotificationType struct {
	general.TenantBase
	NotificationTypeName string `json:"notificationTypeName" example:"New Reply Added"`
}

// Validate validates fields in NotificationType.
func (notificationType *NotificationType) Validate() error {
	if util.IsEmpty(notificationType.NotificationTypeName) {
		return errors.NewValidationError("Notification type name must be specified")
	}
	return nil
}
