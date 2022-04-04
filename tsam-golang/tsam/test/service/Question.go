package service

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/repository"
	tst "github.com/techlabs/swabhav/tsam/models/test"
	"github.com/techlabs/swabhav/tsam/util"
)

// QuestionService Provide method to Update, Delete, Add, Get Method For Question.
type QuestionService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewQuestionService creates a new instance of QuestionService
func NewQuestionService(db *gorm.DB, repository repository.Repository) *QuestionService {
	return &QuestionService{
		DB:         db,
		Repository: repository,
	}

}

// QuestionAssociationNames have all association names
// use for preload Option in Question
var QuestionAssociationNames []string = []string{
	"Options",
}

// AddQuestion Add New Question to Database.
func (ser *QuestionService) AddQuestion(out *tst.Question) error {

	// Add ID to Association if Not Present Or Return Error
	err := out.MakeQuestionValidOrError()
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Create New Instace of UnitOfWork For Write Operation
	// Add Question to Database
	uow := repository.NewUnitOfWork(ser.DB, false)
	err = ser.Repository.Add(uow, out)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// AddQuestions Add Multiple Question to Database.
func (ser *QuestionService) AddQuestions(questions *[]tst.Question, questionsIDs *[]uuid.UUID) error {

	// Add individual Question To Database
	for _, question := range *questions {

		// Add ID To Question
		question.ID = util.GenerateUUID()

		// Add Question To Database
		err := ser.AddQuestion(&question)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
		*questionsIDs = append(*questionsIDs, question.ID)
	}
	return nil
}

// UpdateQuestion Update Question data By Taking
// Question Struct & QueryProcessor
func (ser *QuestionService) UpdateQuestion(out *tst.Question, queryProcessor ...repository.QueryProcessor) error {

	// MakeQuestionValidOrError function validate Question struct
	// for all required field
	err := out.MakeQuestionValidOrError()
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// GetQuestion Call for check Question with provided ID
	// Exist or not in database
	tempQuestion := tst.Question{}
	err = ser.GetQuestion(&tempQuestion, &out.ID)
	if err != nil {
		return err
	}

	// Create New UnitOfWork Instance For Update Question
	// By Repository
	uow := repository.NewUnitOfWork(ser.DB, false)
	err = ser.Repository.Update(uow, out)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteQuestion Delete Question & also there Options
// By Providing Question ID
func (ser *QuestionService) DeleteQuestion(questionID *uuid.UUID, queryProcessor ...repository.QueryProcessor) error {

	// Get Question From Database
	question := tst.Question{}
	err := ser.GetQuestion(&question, questionID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Deleting redudant options in Question
	uow := repository.NewUnitOfWork(ser.DB, false)
	for _, option := range question.Options {
		err = ser.Repository.Delete(uow, option)
		if err != nil {
			return err
		}
	}

	// Delete Question From Database
	err = ser.Repository.Delete(uow, &question)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetQuestion Add New Question to Database.
func (ser *QuestionService) GetQuestion(question *tst.Question, questionID *uuid.UUID,
	queryProcessor ...repository.QueryProcessor) error {
	// Get Question By ID From Database
	uow := repository.NewUnitOfWork(ser.DB, true)
	queryProcessor = append(queryProcessor, repository.PreloadAssociations(QuestionAssociationNames))
	err := ser.Repository.Get(uow, *questionID, question, queryProcessor...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// GetQuestions Return All Question From Database.
func (ser *QuestionService) GetQuestions(questions *[]tst.Question, queryProcessor ...repository.QueryProcessor) error {
	uow := repository.NewUnitOfWork(ser.DB, true)

	// Preload Query
	queryProcessor = append(queryProcessor, repository.PreloadAssociations(QuestionAssociationNames))

	// Get All Question Questions From Database
	err := ser.Repository.GetAll(uow, questions, queryProcessor...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// GetCount Return Total Count Of Result Set.
func (ser *QuestionService) GetCount(out interface{}, count *int, queryProcessor []repository.QueryProcessor) error {
	uow := repository.NewUnitOfWork(ser.DB, true)
	return ser.Repository.GetCount(uow, out, count, queryProcessor...)
}
