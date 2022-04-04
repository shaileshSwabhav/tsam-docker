package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Notification_Test struct {
	TenantBase
	TypeID             *uuid.UUID `json:"typeID" gorm:"type:varchar(36)"`
	SeenTime           *time.Time `json:"seenTime" gorm:"type:varchar(36)"`
	NotifierID         uuid.UUID  `json:"notifierID" gorm:"type:varchar(36)"`
	NotifiedID         uuid.UUID  `json:"notifiedID" gorm:"type:varchar(36)"`
	BlogNotificationID *uuid.UUID `json:"blogNotificationID" gorm:"type:varchar(36)"`
}
