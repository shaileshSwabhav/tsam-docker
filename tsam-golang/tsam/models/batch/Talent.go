package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// TalentDTO is used to get talents for a particular batch.
type TalentDTO struct {
	general.Address
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`
	// Academic     []*talent.Academic    `json:"academics"`
	Technologies []*general.Technology `json:"technologies" gorm:"many2many:talents_technologies;"`
	Code         string                `json:"code"`
	FirstName    *string               `json:"firstName"`
	LastName     *string               `json:"lastName"`
	Email        *string               `json:"email"`
	Contact      *string               `json:"contact"`
	Resume       *string               `json:"resume"`
}

// ValidateTalentBatch validates TalentBatch fields
func (talent *TalentDTO) ValidateTalentBatch() error {
	if talent.FirstName != nil && util.ValidateString(*talent.FirstName) {
		return errors.NewValidationError("Firstname should have only characters")
	}
	if talent.LastName != nil && util.ValidateString(*talent.LastName) {
		return errors.NewValidationError("Firstname should have only characters")
	}
	if talent.Email != nil && util.ValidateEmail(*talent.Email) {
		return errors.NewValidationError("Email should be of the format example@domain.com")
	}
	if talent.Contact != nil && util.ValidateContact(*talent.Contact) {
		return errors.NewValidationError("Contact should be 10 digits")
	}
	return nil
}

// TableName defines table name of the struct.
func (*TalentDTO) TableName() string {
	return "talents"
}
