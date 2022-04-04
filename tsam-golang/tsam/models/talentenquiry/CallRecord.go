package talentenquiry

import (
	"time"

	uuid "github.com/satori/go.uuid"
	// "github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
	// "github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// ExpectedCTCLatest is used to get a single value from the talent call records table.
type ExpectedCTCLatest struct {
	ExpectedCTC int64
}

// CallRecord contains the call record details of enquiry which are enough for adding and updating call record.
type CallRecord struct {
	general.TenantBase

	// Related table IDs.
	PurposeID uuid.UUID `json:"purposeID" gorm:"type:varchar(36)"`
	OutcomeID uuid.UUID `json:"outcomeID" gorm:"type:varchar(36)"`
	EnquiryID uuid.UUID `json:"enquiryID" gorm:"type:varchar(36)"`

	// Other fields.
	DateTime     string  `json:"dateTime" gorm:"type:datetime"`
	Comment      *string `json:"comment" gorm:"type:varchar(500)"`
	ExpectedCTC  *uint   `json:"expectedCTC" gorm:"type:int"`
	NoticePeriod *uint8  `json:"noticePeriod" example:"1" gorm:"type:tinyint(2)"`
	TargetDate   *string `json:"targetDate" gorm:"type:date"`
}

// TableName will name the table of CallRecord model as "talent_enquiry_call_records".
func (*CallRecord) TableName() string {
	return "talent_enquiry_call_records"
}

// Validate Validates the required fields of call record.
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

	// Expected CTC.
	if callRecord.ExpectedCTC != nil && *callRecord.ExpectedCTC < 100000 {
		return errors.NewValidationError("Expected CTC should be equal to or above 100000")
	}

	if callRecord.ExpectedCTC != nil && *callRecord.ExpectedCTC > 100000000 {
		return errors.NewValidationError("Expected CTC should be equal to or less than 100000000")
	}

	// Notice Period.
	if callRecord.NoticePeriod != nil && *callRecord.NoticePeriod < 0 {
		return errors.NewValidationError("Notice Period should be equal to or above 0")
	}

	if callRecord.NoticePeriod != nil && *callRecord.NoticePeriod > 9 {
		return errors.NewValidationError("Notice Period should be equal to or less than 10")
	}

	// Comment maximum length.
	if callRecord.Comment != nil && len(*(*callRecord).Comment) > 500 {
		return errors.NewValidationError("Comment can have maximum 500 characters")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// CallRecordDTO contains the complete information of enquiry call record which is needed to display.
type CallRecordDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Related tables.
	Purpose   *general.Purpose `json:"purpose" gorm:"foreignkey:PurposeID"`
	PurposeID uuid.UUID        `json:"-"`
	Outcome   *general.Outcome `json:"outcome" gorm:"foreignkey:OutcomeID"`
	OutcomeID uuid.UUID        `json:"-"`

	// Other fields.
	DateTime     string    `json:"dateTime"`
	Comment      *string   `json:"comment"`
	ExpectedCTC  *uint     `json:"expectedCTC"`
	NoticePeriod *uint8    `json:"noticePeriod" example:"1"`
	EnquiryID    uuid.UUID `json:"enquiryID"`
	TargetDate   *string   `json:"targetDate"`
}

// TableName will name the table of CallRecord model as "talent_call_records".
func (*CallRecordDTO) TableName() string {
	return "talent_enquiry_call_records"
}
