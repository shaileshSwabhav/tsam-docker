package company

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// CallRecord contains information realted to enquiry call records needed for add and update.
type CallRecord struct {
	general.TenantBase

	// Related table IDs.
	PurposeID uuid.UUID `json:"purposeID" gorm:"type:varchar(36)"`
	OutcomeID uuid.UUID `json:"outcomeID" gorm:"type:varchar(36)"`
	EnquiryID uuid.UUID `json:"enquiryID" gorm:"type:varchar(36)"`

	// Other fields.
	DateTime string  `json:"dateTime" gorm:"type:datetime"`
	Comment  *string `json:"comment" gorm:"type:varchar(500)"`
}

// TableName will name the table of CallRecord model as "company_enquiry_call_records".
func (*CallRecord) TableName() string {
	return "company_enquiry_call_records"
}

// ValidateCallRecord validates the required fields of call record.
func (callRecord *CallRecord) Validate() error {

	// Purpose ID.
	if !util.IsUUIDValid(callRecord.PurposeID) {
		return errors.NewValidationError("Purpose ID must ne a proper uuid")
	}

	// Outcome ID.
	if !util.IsUUIDValid(callRecord.OutcomeID) {
		return errors.NewValidationError("Outcome ID must ne a proper uuid")
	}

	// Date time.
	if util.IsEmpty(callRecord.DateTime) {
		return errors.NewValidationError("Call record must have a date time.")
	}

	// Comment maximum length.
	if callRecord.Comment != nil && len(*(*callRecord).Comment) > 500 {
		return errors.NewValidationError("Comment can have maximum 500 characters")
	}
	return nil
}

//************************************* DTO MODEL *************************************************************

// CallRecord contains information realted to enquiry call records needed for add and update.
type CallRecordDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Related tables.
	Purpose   *general.Purpose `json:"purpose" gorm:"foreignkey:PurposeID"`
	PurposeID uuid.UUID        `json:"-"`
	Outcome   *general.Outcome `json:"outcome" gorm:"foreignkey:OutcomeID"`
	OutcomeID uuid.UUID        `json:"-"`

	// Other fields.
	DateTime  string    `json:"dateTime"`
	Comment   *string   `json:"comment"`
	EnquiryID uuid.UUID `json:"enquiryID"`
}

// TableName will name the table of CallRecord model as "company_enquiry_call_records"
func (*CallRecordDTO) TableName() string {
	return "company_enquiry_call_records"
}
