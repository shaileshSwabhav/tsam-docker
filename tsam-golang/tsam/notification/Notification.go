package notification

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
)

type Notification struct {
	general.TenantBase
	TypeID     uuid.UUID
	SeenTime   time.Time
	NotifierID uuid.UUID
	NotifiedID uuid.UUID
}
