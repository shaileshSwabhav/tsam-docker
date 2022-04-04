package programming

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Option contains add update fields required for ProgrammingQuestionOption.
type ProgrammingQuestionOption struct {
	general.TenantBase
	ProgrammingQuestionID uuid.UUID `json:"programmingQuestionID" gorm:"type:varchar(36)"`
	Option                string    `json:"option" gorm:"type:varchar(100)"`
	IsCorrect             *bool     `json:"isCorrect"`
	Order                 uint8     `json:"order" gorm:"type:tinyint(2)"`
	IsActive              *bool     `json:"isActive"`
}

// Validate validates compulsary fields of ProgrammingQuestionOption.
func (option *ProgrammingQuestionOption) Validate() error {

	// Question ID.
	if !util.IsUUIDValid(option.ProgrammingQuestionID) {
		return errors.NewValidationError("Qusetione ID must ne a proper uuid")
	}

	// Check if option is blank or not.
	if util.IsEmpty(option.Option) {
		return errors.NewValidationError("Option must be specified")
	}

	// Option maximum characters.
	if len(option.Option) > 100 {
		return errors.NewValidationError("Option can have maximum 100 characters")
	}

	// Order must be specified.
	if option.Order == 0 {
		return errors.NewValidationError("Order must be specified")
	}

	// Order maximum.
	if option.Order > 99 {
		return errors.NewValidationError("Order cannot be above 99")
	}

	return nil
}

// Option contains add update fields required for ProgrammingQuestionOption.
type ProgrammingQuestionOptionDTO struct {
	ID                    uuid.UUID  `json:"id"`
	DeletedAt             *time.Time `json:"-"`
	ProgrammingQuestionID uuid.UUID  `json:"programmingQuestionID" gorm:"type:varchar(36)"`
	Option                string     `json:"option" gorm:"type:varchar(100)"`
	IsCorrect             bool       `json:"isCorrect"`
	Order                 uint8      `json:"order" gorm:"type:tinyint(2)"`
	IsActive              bool       `json:"isActive"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionOptionDTO) TableName() string {
	return "programming_question_options"
}
