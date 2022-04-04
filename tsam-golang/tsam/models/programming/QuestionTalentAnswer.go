package programming

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// ProgrammingQuestionTalentAnswer contains add update fields required for programming question talent answer.
type ProgrammingQuestionTalentAnswer struct {
	general.TenantBase

	Answer *string `json:"answer" gorm:"type:varchar(2000)"`
	Score  uint16  `json:"score" gorm:"type:smallint(4)"`

	// Related table IDs.
	ProgrammingQuestionID       uuid.UUID  `json:"programmingQuestionID" gorm:"type:varchar(36)"`
	TalentID                    uuid.UUID  `json:"talentID" gorm:"type:varchar(36)"`
	ProgrammingQuestionOptionID *uuid.UUID `json:"programmingQuestionOptionID" gorm:"type:varchar(36)"`
	ProgrammingLanguageID       *uuid.UUID `json:"programmingLanguageID" gorm:"type:varchar(36)"`
	ProgrammingQuestionTypeID   uuid.UUID `json:"programmingQuestionTypeID" gorm:"type:varchar(36)"`

	// Flags.
	IsCorrect *bool `json:"isCorrect"`

	// Future change.
	//CheckedBy
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionTalentAnswer) TableName() string {
	return "programming_question_talent_answers"
}

// Validate validates compulsary fields of ProgrammingQuestionTalentAnswer.
func (answer *ProgrammingQuestionTalentAnswer) Validate() error {

	// Programming Question ID.
	if !util.IsUUIDValid(answer.ProgrammingQuestionID) {
		return errors.NewValidationError("Programming Question ID must be a proper uuid")
	}

	// Talent ID.
	if !util.IsUUIDValid(answer.TalentID) {
		return errors.NewValidationError("Talent ID must be a proper uuid")
	}

	// Score minimum.
	if answer.Score < 0 {
		return errors.NewValidationError("Score cannot be below 0")
	}

	// Score maximum.
	if answer.Score > 100 {
		return errors.NewValidationError("Score cannot be above 100")
	}

	// Check if answer and option id both are blank.
	if answer.Answer == nil && answer.ProgrammingQuestionOptionID == nil {
		return errors.NewValidationError("Answer or option must be specified")
	}

	// Check if answer and option id both are blank.
	if answer.Answer != nil && answer.ProgrammingQuestionOptionID != nil {
		return errors.NewValidationError("Answer and option both cannot be specified")
	}

	// If answer is not specified then check if option id is valid.
	if answer.Answer == nil && !util.IsUUIDValid(*answer.ProgrammingQuestionOptionID) {
		return errors.NewValidationError("Option ID must be a proper uuid")
	}

	// If option id is not specified then check if answer is not empty.
	if answer.ProgrammingQuestionOptionID == nil && util.IsEmpty(*answer.Answer) {
		return errors.NewValidationError("Answer must be specified")
	}

	// Answer maximum characters.
	if answer.ProgrammingQuestionOptionID == nil && len(*answer.Answer) > 2000 {
		return errors.NewValidationError("Answer can have maximum 2000 characters")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// ProgrammingQuestionTalentAnswerDTO is used for getting programming question talent answer from database
// including all the programming question talent answer related information also.
type ProgrammingQuestionTalentAnswerDTO struct {
	ID uuid.UUID `json:"id"`

	Answer          *string `json:"answer"`
	Score           uint16  `json:"score"`
	Date            string  `json:"date"`
	TotalAnswers    uint16  `json:"totalAnswers"`
	TotalNotChecked uint16  `json:"totalNotChecked"`

	// Single model.
	ProgrammingQuestion         ProgrammingQuestionTalentAnswerQuestionDTO `json:"programmingQuestion" gorm:"foreignkey:ProgrammingQuestionID"`
	ProgrammingQuestionID       uuid.UUID                                  `json:"-"`
	ProgrammingQuestionOption   ProgrammingQuestionOption                  `json:"programmingQuestionOption" gorm:"foreignkey:ProgrammingQuestionOptionID"`
	ProgrammingQuestionOptionID uuid.UUID                                  `json:"-"`
	ProgrammingLanguage         *general.ProgrammingLanguage               `json:"programmingLanguage" gorm:"foreignkey:ProgrammingLanguageID"`
	ProgrammingLanguageID       *uuid.UUID                                 `json:"-"`
	Talent                      ProgrammingQuestionTalentAnswerTalentDTO   `json:"talent" gorm:"foreignkey:TalentID"`
	TalentID                    uuid.UUID                                  `json:"-"`
	ProgrammingQuestionTypeID   uuid.UUID                                  `json:"-"`
	ProgrammingQuestionType     ProgrammingQuestionType                    `json:"programmingQuestionType" gorm:"foreignkey:ProgrammingQuestionTypeID"`

	// Flags.
	IsCorrect *bool `json:"isCorrect"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionTalentAnswerDTO) TableName() string {
	return "programming_question_talent_answers"
}

// ProgrammingQuestionTalentAnswerQuestionDTO contains programming question information to be given 
// along with programming question talent answer.
type ProgrammingQuestionTalentAnswerQuestionDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	Label           string `json:"label"`
	Level           uint8  `json:"level"`
	Score           uint16 `json:"score"`
	SolutonIsViewed bool   `json:"solutonIsViewed"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionTalentAnswerQuestionDTO) TableName() string {
	return "programming_questions"
}

// ProgrammingQuestionTalentAnswerTalentDTO contains talent information to be given along with 
// programming question talent answer.
type ProgrammingQuestionTalentAnswerTalentDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionTalentAnswerTalentDTO) TableName() string {
	return "talents"
}

// ProgrammingQuestionTalentAnswerWithFullQuestionDTO is used for getting programming question 
// talent answer with whole question from database.
type ProgrammingQuestionTalentAnswerWithFullQuestionDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	Answer *string `json:"answer"`
	Score  uint16  `json:"score"`

	// Single model.
	ProgrammingQuestion   ProgrammingQuestionWithTalentAnswerDTO   `json:"programmingQuestion" gorm:"foreignkey:ProgrammingQuestionID"`
	ProgrammingQuestionID uuid.UUID                                `json:"-"`
	Talent                ProgrammingQuestionTalentAnswerTalentDTO `json:"talent" gorm:"foreignkey:TalentID"`
	TalentID              uuid.UUID                                `json:"-"`
	ProgrammingLanguage   *general.ProgrammingLanguage             `json:"programmingLanguage" gorm:"foreignkey:ProgrammingLanguageID"`
	ProgrammingLanguageID *uuid.UUID                               `json:"-"`

	// Flags.
	IsCorrect *bool `json:"isCorrect"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionTalentAnswerWithFullQuestionDTO) TableName() string {
	return "programming_question_talent_answers"
}

// ProgrammingQuestionTalentAnswerScore is used for updating the score and isCorrect field of 
// programming question talent answer.
type ProgrammingQuestionTalentAnswerScore struct {
	ID        uuid.UUID `json:"id"`
	IsCorrect *bool     `json:"isCorrect"`
	Score     uint16    `json:"score"`
}

// TotalCount is used for getting total count of the programming question talent answers.
type TotalCount struct {
	TotalCount int
}
