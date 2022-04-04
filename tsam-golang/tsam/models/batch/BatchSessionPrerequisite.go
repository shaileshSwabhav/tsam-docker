package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// BatchSessionPrerequisite will store pre-requisite for batch session.
type BatchSessionPrerequisite struct {
	general.TenantBase
	BatchID        uuid.UUID `json:"batchID" gorm:"type:varchar(36)"`
	BatchSessionID uuid.UUID `json:"batchSessionID" gorm:"type:varchar(36)"`
	Prerequisite   string    `json:"prerequisite" gorm:"type:varchar(1000)"`
}

// Validate will checks if table fields are valid.
func (prerequistie *BatchSessionPrerequisite) Validate() error {

	if util.IsEmpty(prerequistie.Prerequisite) {
		return errors.NewValidationError("pre-requisite must be specified")
	}

	return nil
}

// BatchSessionPrerequisiteDTO will store prerequistie for the given batch session.
type BatchSessionPrerequisiteDTO struct {
	ID             uuid.UUID  `json:"id"`
	DeletedAt      *time.Time `json:"-"`
	BatchID        uuid.UUID  `json:"-"`
	BatchSessionID uuid.UUID  `json:"-"`
	Prerequisite   string     `json:"prerequisite"`
}

// TableName defines table name of the struct.
func (*BatchSessionPrerequisiteDTO) TableName() string {
	return "batch_session_prerequisites"
}
