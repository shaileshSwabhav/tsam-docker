package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// BatchSessionHoliday will specified holiday dates for batch.
type BatchSessionHoliday struct {
	general.TenantBase
	HolidayDate string    `json:"holidayDate" gorm:"date"`
	BatchID     uuid.UUID `json:"batchID" gorm:"type:varchar(36)"`
}

func (session *BatchSessionHoliday) Validate() error {

	if util.IsEmpty(session.HolidayDate) {
		return errors.NewValidationError("holiday dates must be specified")
	}

	if session.BatchID == uuid.Nil {
		return errors.NewValidationError("batch ID must be specified")
	}

	return nil
}

// BatchSessionHolidayDTO will specified holiday dates for batch.
type BatchSessionHolidayDTO struct {
	ID          uuid.UUID  `json:"id"`
	DeletedAt   *time.Time `json:"-"`
	HolidayDate string     `json:"holidayDate"`
	BatchID     uuid.UUID  `json:"batchID"`
}

// TableName defines table name of the struct.
func (*BatchSessionHolidayDTO) TableName() string {
	return "batch_session_holidays"
}
