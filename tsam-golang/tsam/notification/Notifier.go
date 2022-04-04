package notification

import (
	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/repository"
)

// Notifier must be implemented by custom notifications
type Notifier interface {
	Notify(db *gorm.DB, repo repository.Repository)
}
