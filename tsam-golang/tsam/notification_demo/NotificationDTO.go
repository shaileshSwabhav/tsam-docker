package notificationdemo

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/blog"
	"github.com/techlabs/swabhav/tsam/models/general"
)

type Notification_Test_DTO struct {
	general.TenantBase
	TypeID             *uuid.UUID             `json:"typeID" gorm:"type:varchar(36)"`
	SeenTime           *time.Time             `json:"seenTime" gorm:"type:varchar(36)"`
	NotifierID         uuid.UUID              `json:"notifierID" gorm:"type:varchar(36)"`
	NotifiedID         uuid.UUID              `json:"notifiedID" gorm:"type:varchar(36)"`
	BlogNotificationID *uuid.UUID             `json:"blogNotificationID" gorm:"type:varchar(36)"`
	BlogNotifications  *blog.BlogNotification `json:"blogNotifications" gorm:"foreignkey:BlogNotificationID"`
}

func (*Notification_Test_DTO) TableName() string {
	return "notification_tests"
}
