package event

import (
	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/notification"
)

// Event consists of notification as of now. (beta)
type Event struct {
	Name     Name
	DB       *gorm.DB
	Notifier notification.Notifier
}
