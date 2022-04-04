package programming

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// ProgrammingQuestion contains add update fields required for programming questions.
type ProgrammingQuestion struct {
	general.TenantBase
	Label                    string                         `json:"label" gorm:"type:varchar(100)"`
	Question                 string                         `json:"question" gorm:"type:varchar(1000)"`
	Example                  *string                        `json:"example" gorm:"type:varchar(2000)"`
	Constraints              *string                        `json:"constraints" gorm:"type:varchar(500)"`
	Comment                  *string                        `json:"comment" gorm:"type:varchar(500)"`
	HasOptions               bool                           `json:"hasOptions"`
	IsActive                 *bool                          `json:"isActive"`
	IsLanguageSpecific       bool                           `json:"isLanguageSpecific"`
	Level                    uint8                          `json:"level" gorm:"type:tinyint(2)"`
	Score                    uint16                         `json:"score" gorm:"type:smallint(4)"`
	TimeRequired             *uint                          `json:"timeRequired" gorm:"type:smallint(4)"`
	ProgrammingQuestionTypes []ProgrammingQuestionType      `json:"programmingQuestionTypes" gorm:"many2many:programming_questions_programming_question_types;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	ProgrammingConcepts      []*ProgrammingConcept          `json:"programmingConcept" gorm:"many2many:programming_questions_programming_concepts;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Options                  []ProgrammingQuestionOption    `json:"options"`
	TestCases                []*ProgrammingQuestionTestCase `json:"testCases"`
	ProgrammingLanguages     []*general.ProgrammingLanguage `json:"programmingLanguages" gorm:"many2many:programming_questions_programming_languages;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	// ProgrammingProblemType string `json:"programmingProblemType" gorm:"type:varchar(100)"`
}

// Validate validates compulsary fields of ProgrammingQuestion.
func (question *ProgrammingQuestion) Validate() error {

	// Check if label is blank or not.
	if util.IsEmpty(question.Label) {
		return errors.NewValidationError("Label must be specified")
	}

	// Label maximum characters.
	if len(question.Label) > 100 {
		return errors.NewValidationError("Label can have maximum 100 characters")
	}

	// Check if question is blank or not.
	if util.IsEmpty(question.Question) {
		return errors.NewValidationError("Question must be specified")
	}

	// Question maximum characters.
	if len(question.Question) > 1000 {
		return errors.NewValidationError("Question can have maximum 1000 characters")
	}

	// Check if example is blank or not when has options is false.
	if !question.HasOptions && question.Example == nil {
		return errors.NewValidationError("Example must be specified when options are not given")
	}

	// Check if example is blank or not when has options is false.
	if !question.HasOptions && util.IsEmpty(*question.Example) {
		return errors.NewValidationError("Example must be specified when options are not given")
	}

	// Example maximum characters when has options is false.
	if !question.HasOptions && len(*question.Example) > 2000 {
		return errors.NewValidationError("Example can have maximum 2000 characters")
	}

	// Check if constraints is blank or not when has options is false.
	if !question.HasOptions && question.Constraints == nil {
		return errors.NewValidationError("Constraints must be specified when options are not given")
	}

	// Check if constraints is blank or not when has options is false.
	if !question.HasOptions && util.IsEmpty(*question.Constraints) {
		return errors.NewValidationError("Constraints must be specified when options are not given")
	}

	// Constraints maximum characters when has options is false.
	if !question.HasOptions && len(*question.Constraints) > 500 {
		return errors.NewValidationError("Constraints can have maximum 500 characters")
	}

	// Comment maximum characters.
	if question.Comment != nil && len(*question.Comment) > 500 {
		return errors.NewValidationError("Comment can have maximum 500 characters")
	}

	// Level must be specified.
	if question.Level == 0 {
		return errors.NewValidationError("Level must be specified")
	}

	// Level maximum.
	if question.Level > 99 {
		return errors.NewValidationError("Level cannot be above 99")
	}

	// Score minimum.
	if question.Score <= 0 {
		return errors.NewValidationError("Score cannot be below 1")
	}

	// Score maximum.
	if question.Score > 100 {
		return errors.NewValidationError("Score cannot be above 100")
	}

	// Time required minimum.
	if question.TimeRequired != nil && *question.TimeRequired <= 0 {
		return errors.NewValidationError("Time required cannot be below 1")
	}

	// Time required maximum.
	if question.TimeRequired != nil && *question.TimeRequired > 1440 {
		return errors.NewValidationError("Time required cannot be above 1440")
	}

	// Programming Question Types.
	if question.ProgrammingQuestionTypes == nil {
		return errors.NewValidationError("Programming Question Types must be specified")
	}

	// Programming Question Types.
	if question.ProgrammingQuestionTypes != nil && len(question.ProgrammingQuestionTypes) <= 0 {
		return errors.NewValidationError("Programming Question Types must be specified")
	}

	// Check if all options have unique order.
	optionMap := make(map[uint8]uint)
	isCorrectCount := 0
	for _, option := range question.Options {
		optionMap[option.Order]++
		if *option.IsCorrect && *option.IsActive {
			isCorrectCount++
		}
		if isCorrectCount > 1 {
			return errors.NewValidationError("Question cannot have multiple correct answers.")
		}

		if optionMap[option.Order] > 1 {
			return errors.NewValidationError("Options with same order not allowed.")
		}
	}

	if isCorrectCount == 0 && question.HasOptions {
		return errors.NewValidationError("Option should consist of one active correct answer.")
	}

	// Check if test cases exist, if true then validate.
	if question.TestCases != nil {
		for _, testCase := range question.TestCases {
			if err := testCase.Validate(); err != nil {
				return err
			}
		}
	}

	// If is language specific then languages must be specified.
	if question.IsLanguageSpecific && (question.ProgrammingLanguages == nil || (question.ProgrammingLanguages != nil && len(question.ProgrammingLanguages) == 0)) {
		return errors.NewValidationError("If question is language specific then languages must be specified")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// ProgrammingQuestionDTO contains all fields required for programming questions.
type ProgrammingQuestionDTO struct {
	ID                       uuid.UUID                      `json:"id"`
	DeletedAt                *time.Time                     `json:"-"`
	Label                    string                         `json:"label"`
	Question                 string                         `json:"question"`
	Example                  *string                        `json:"example"`
	Constraints              *string                        `json:"constraints"`
	Comment                  *string                        `json:"comment"`
	HasOptions               bool                           `json:"hasOptions"`
	IsActive                 bool                           `json:"isActive"`
	IsLanguageSpecific       bool                           `json:"isLanguageSpecific"`
	Level                    uint8                          `json:"level"`
	Score                    uint16                         `json:"score"`
	TimeRequired             *uint                          `json:"timeRequired"`
	ProgrammingQuestionTypes []ProgrammingQuestionType      `json:"programmingQuestionTypes" gorm:"many2many:programming_questions_programming_question_types;association_jointable_foreignkey:programming_question_type_id;jointable_foreignkey:programming_question_id"`
	ProgrammingConcepts      []*ProgrammingConcept          `json:"programmingConcept" gorm:"many2many:programming_questions_programming_concepts;association_jointable_foreignkey:programming_concept_id;jointable_foreignkey:programming_question_id"`
	Options                  []*ProgrammingQuestionOption   `json:"options" gorm:"foreignkey:ProgrammingQuestionID"`
	SolutionCount            uint16                         `json:"solutionCount"`
	HasAnyTalentAnswered     bool                           `json:"hasAnyTalentAnswered"`
	ProgrammingLanguages     []*general.ProgrammingLanguage `json:"programmingLanguages" gorm:"many2many:programming_questions_programming_languages;association_jointable_foreignkey:programming_language_id;jointable_foreignkey:programming_question_id"`
	// ProgrammingProblemType string `json:"programmingProblemType" gorm:"type:varchar(100)"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionDTO) TableName() string {
	return "programming_questions"
}

// QuestionProblemOfTheDayDTO for getting details required for problems of the day.
type QuestionProblemOfTheDayDTO struct {
	ID               uuid.UUID  `json:"id"`
	DeletedAt        *time.Time `json:"-"`
	Label            string     `json:"label"`
	Level            uint8      `json:"level"`
	Score            uint16     `json:"score"`
	AttemptedByCount uint       `json:"attemptedByCount"`
	SolvedByCount    uint       `json:"solvedByCount"`
}

// TableName defines table name of the struct.
func (*QuestionProblemOfTheDayDTO) TableName() string {
	return "programming_questions"
}

// ProgrammingQuestionWithTalentAnswerDTO contains all fields required for programming questions.
// along with talent related details.
type ProgrammingQuestionWithTalentAnswerDTO struct {
	ID                          uuid.UUID                        `json:"id"`
	DeletedAt                   *time.Time                       `json:"-"`
	Label                       string                           `json:"label"`
	Question                    string                           `json:"question"`
	Example                     *string                          `json:"example"`
	Constraints                 *string                          `json:"constraints"`
	TestCases                   []*ProgrammingQuestionTestCase   `json:"testCases" gorm:"foreignkey:ProgrammingQuestionID"`
	Comment                     *string                          `json:"comment"`
	HasOptions                  bool                             `json:"hasOptions"`
	IsActive                    bool                             `json:"isActive"`
	Level                       uint8                            `json:"level"`
	Score                       uint16                           `json:"score"`
	TimeRequired                *uint                            `json:"timeRequired"`
	ProgrammingQuestionTypes    []ProgrammingQuestionType        `json:"programmingQuestionTypes" gorm:"many2many:programming_questions_programming_question_types;association_jointable_foreignkey:programming_question_type_id;jointable_foreignkey:programming_question_id"`
	Options                     []*ProgrammingQuestionOption     `json:"options" gorm:"foreignkey:ProgrammingQuestionID"`
	ProgrammingQuestionOptionID *string                          `json:"programmingQuestionOptionID"`
	Answer                      *string                          `json:"answer"`
	IsAnswered                  bool                             `json:"isAnswered"`
	AttemptedByCount            uint                             `json:"attemptedByCount"`
	SolvedByCount               uint                             `json:"solvedByCount"`
	Solutions                   []ProgrammingQuestionSolutionDTO `json:"solutions" gorm:"foreignkey:ProgrammingQuestionID"`
	ProgrammingLanguageID       uuid.UUID                        `json:"-"`
	ProgrammingLanguage         general.ProgrammingLanguage      `json:"programmingLanguage" gorm:"foreignkey:ProgrammingLanguageID"`
	SolutonIsViewed             bool                             `json:"solutonIsViewed"`
}

// TableName defines table name of the struct.
func (*ProgrammingQuestionWithTalentAnswerDTO) TableName() string {
	return "programming_questions"
}

// ProgrammingQuestionIsActive is used for updating the isActive field of programming question.
type ProgrammingQuestionIsActive struct {
	ID       uuid.UUID `json:"id"`
	IsActive *bool     `json:"isActive"`
}

// TopicProgrammingQuestionDTO defines fields that are used for get operation.
type TopicProgrammingQuestionDTO struct {
	ID                    uuid.UUID  `json:"id"`
	DeletedAt             *time.Time `json:"-"`
	TopicID               uuid.UUID  `json:"-"`
	ProgrammingQuestionID uuid.UUID  `json:"-"`
	IsActive              *bool      `json:"isActive"`
}

// TableName defines table name of the struct.
func (*TopicProgrammingQuestionDTO) TableName() string {
	return "topic_programming_questions"
}
