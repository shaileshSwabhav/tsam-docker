package programming

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingQuestionType contains type for programming question.
type ProgrammingQuestionType struct {
	general.TenantBase
	ProgrammingType string `json:"programmingType" gorm:"type:varchar(100)"`
}

// Validate validates compulsary fields of ProgrammingQuestionType.
func (questionType *ProgrammingQuestionType) Validate() error {

	// Check if type is blank or not.
	if util.IsEmpty(questionType.ProgrammingType) {
		return errors.NewValidationError("Type must be specified")
	}

	// Type maximum characters.
	if len(questionType.ProgrammingType) > 100 {
		return errors.NewValidationError("Type can have maximum 100 characters")
	}

	return nil
}

// ProgrammingQuestionTypeDTO contains type for programming question.
type ProgrammingQuestionTypeDTO struct {
	ID              uuid.UUID  `json:"id"`
	DeletedAt       *time.Time `json:"-"`
	ProgrammingType string     `json:"programmingType" gorm:"type:varchar(100)"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionTypeDTO) TableName() string {
	return "programming_question_types"
}

// ===========Defining many to many structs===========

// ProgrammingQuestionProgrammingQuestionTypes is the map of programming question and programming question type.
type ProgrammingQuestionProgrammingQuestionTypes struct {
	ProgrammingQuestionID     uuid.UUID `gorm:"type:varchar(36)"`
	ProgrammingQuestionTypeID uuid.UUID `gorm:"type:varchar(36)"`
}
