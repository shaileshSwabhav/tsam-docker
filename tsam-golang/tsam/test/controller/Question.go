package controller

import (
	"net/http"

	"github.com/techlabs/swabhav/tsam/repository"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	tst "github.com/techlabs/swabhav/tsam/models/test"
	"github.com/techlabs/swabhav/tsam/test/service"
	"github.com/techlabs/swabhav/tsam/util"

	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/web"
)

// QuestionController Provide method to Update, Delete, Add, Get Method For Question.
type QuestionController struct {
	QuestionService *service.QuestionService
}

// NewQuestionController Create New Instance Of QuestionController.
func NewQuestionController(ser *service.QuestionService) *QuestionController {
	return &QuestionController{
		QuestionService: ser,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (questionCon *QuestionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/question", questionCon.AddQuestion).Methods(http.MethodPost)
	router.HandleFunc("/questions", questionCon.AddQuestions).Methods(http.MethodPost)
	router.HandleFunc("/question/{questionid}", questionCon.GetQuestion).Methods(http.MethodGet)
	router.HandleFunc("/question/{questionid}", questionCon.UpdateQuestion).Methods(http.MethodPut)
	router.HandleFunc("/question/{questionid}", questionCon.DeleteQuestion).Methods(http.MethodDelete)
	router.HandleFunc("/question/{limit}/{offset}", questionCon.GetQuestions).Methods(http.MethodGet)
	router.HandleFunc("/question/search/{limit}/{offset}", questionCon.GetQuestions).Methods(http.MethodPost)
	log.NewLogger().Info("Question Route Registered")
}

//GetQuestions Return All questions
func (questionCon *QuestionController) GetQuestions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Get All Question API Call")
	questions := &[]tst.Question{}
	err := questionCon.QuestionService.GetQuestions(questions, repository.Paging(w, r))
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, questions)
}

//GetQuestion Return Specific Question
func (questionCon *QuestionController) GetQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Get Question API Call")

	// Parse Question ID
	param := mux.Vars(r)
	questionID, err := util.ParseUUID(param["questionid"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable To Parse QuestionID", http.StatusBadRequest))
		return
	}

	question := &tst.Question{}
	err = questionCon.QuestionService.GetQuestion(question, &questionID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, question)
}

// UpdateQuestion Update The Question
func (questionCon *QuestionController) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Update Question API Call")
	question := &tst.Question{}

	// Parse Data From Request
	param := mux.Vars(r)
	questionID, err := util.ParseUUID(param["questionid"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse Question ID", http.StatusBadRequest))
		return
	}

	err = web.UnmarshalJSON(r, &question)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}
	question.ID = questionID

	err = questionCon.QuestionService.UpdateQuestion(question)
	if err != nil {
		web.RespondError(w, err)
		return
	}
}

//AddQuestion Add New Question
func (questionCon *QuestionController) AddQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Add Question API Call")

	// Parse Question From Request & Add New ID.
	question := tst.Question{}
	err := web.UnmarshalJSON(r, &question)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}
	question.ID = util.GenerateUUID()

	// Add Question To Database
	err = questionCon.QuestionService.AddQuestion(&question)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, question.ID)
}

// SearchQuestion Search Questions
func (questionCon *QuestionController) SearchQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Search Question API Call")

	// Parse Question From Request & Add New ID.
	searchQuestion := tst.Question{}
	err := web.UnmarshalJSON(r, &searchQuestion)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	questions := []tst.Question{}
	// Add Question To Database
	err = questionCon.QuestionService.GetQuestions(&questions,
		questionCon.questionSearchQueries(&searchQuestion), repository.Paging(w, r))
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, &questions)
}

// AddQuestions Add Multiple New Question
func (questionCon *QuestionController) AddQuestions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Add Question API Call")
	questionsIDs := []uuid.UUID{}
	questions := []tst.Question{}

	// Parse Questions From Request
	err := web.UnmarshalJSON(r, &questions)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse requested data", http.StatusBadRequest))
		return
	}

	err = questionCon.QuestionService.AddQuestions(&questions, &questionsIDs)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, questionsIDs)
}

// DeleteQuestion Delete Question
func (questionCon *QuestionController) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Delete Question API Call")
	param := mux.Vars(r)
	questionID, err := util.ParseUUID(param["questionid"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse Question ID", http.StatusBadRequest))
		return
	}
	err = questionCon.QuestionService.DeleteQuestion(&questionID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
}

// questionSearchQueries adds all search queries by comparing with the Question data recieved from
// Request Body
func (questionCon *QuestionController) questionSearchQueries(question *tst.Question) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if !util.IsNil(question.Subject) {
		util.AddToSlice("subject", "LIKE ?", "AND", "%"+*question.Subject+"%", &columnNames, &conditions, &operators, &values)
	}
	if !util.IsNil(question.Difficulty) {
		util.AddToSlice("difficulty", "LIKE ?", "AND", "%"+*question.Difficulty+"%", &columnNames, &conditions, &operators, &values)
	}
	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
