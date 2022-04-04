package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingQuestionController provides methods to do CRUD operations.
type ProgrammingQuestionController struct {
	ProgrammingQuestionService *service.ProgrammingQuestionService
}

// NewProgrammingQuestionController creates new instance of programming question type controller.
func NewProgrammingQuestionController(programmingQuestionService *service.ProgrammingQuestionService) *ProgrammingQuestionController {
	return &ProgrammingQuestionController{
		ProgrammingQuestionService: programmingQuestionService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *ProgrammingQuestionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add one question.
	router.HandleFunc("/tenant/{tenantID}/programming-question/credential/{credentialID}",
		controller.AddProgrammingQuestion).Methods(http.MethodPost)

	// Add multiple qestions.
	router.HandleFunc("/tenant/{tenantID}/programming-questions/credential/{credentialID}",
		controller.AddProgrammingQuestions).Methods(http.MethodPost)

	// Update one question.
	router.HandleFunc("/tenant/{tenantID}/programming-question/{programmingQuestionID}/credential/{credentialID}",
		controller.UpdateProgrammingQuestion).Methods(http.MethodPut)

	// Update isActive of question.
	router.HandleFunc("/tenant/{tenantID}/programming-question/is-active/{programmingQuestionID}/credential/{credentialID}",
		controller.UpdateProgrammingQuestionIsActive).Methods(http.MethodPut)

	// Delete one question.
	router.HandleFunc("/tenant/{tenantID}/programming-question/{programmingQuestionID}/credential/{credentialID}",
		controller.DeleteProgrammingQuestion).Methods(http.MethodDelete)

	// Get questions by limit and offset.
	router.HandleFunc("/tenant/{tenantID}/programming-question/limit/{limit}/offset/{offset}",
		controller.GetProgrammingQuestions).Methods(http.MethodGet)

	// Get question list.
	router.HandleFunc("/tenant/{tenantID}/programming-question",
		controller.GetProgrammingQuestionList).Methods(http.MethodGet)

	// Get one question.
	router.HandleFunc("/tenant/{tenantID}/programming-question/{programmingQuestionID}",
		controller.GetProgrammingQuestion).Methods(http.MethodGet)

	// Get questions in problem of the day format for practice.
	router.HandleFunc("/tenant/{tenantID}/programming-question-practice/limit/{limit}/offset/{offset}",
		controller.GetProgrammingQuestionsPractice).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/topics/programming-questions",
		controller.GetProgrammingQuestionForTopic).Methods(http.MethodGet)

	log.NewLogger().Info("Programming Question Routes Registered")
}

// AddProgrammingQuestion will add programming question to the table.
func (controller *ProgrammingQuestionController) AddProgrammingQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddProgrammingQuestion called==============================")
	param := mux.Vars(r)
	question := new(programming.ProgrammingQuestion)

	err := web.UnmarshalJSON(r, question)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	question.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	question.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = question.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProgrammingQuestionService.AddProgrammingQuestion(question)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusCreated, question.ID)
}

// AddProgrammingQuestions will add multiple programming questions to the table.
func (controller *ProgrammingQuestionController) AddProgrammingQuestions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddProgrammingQuestions called==============================")
	param := mux.Vars(r)
	questions := new([]programming.ProgrammingQuestion)

	err := web.UnmarshalJSON(r, questions)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, question := range *questions {
		err = question.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = controller.ProgrammingQuestionService.AddProgrammingQuestions(questions, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Programming questions added successfully")
}

// UpdateProgrammingQuestion will update specified programming question.
func (controller *ProgrammingQuestionController) UpdateProgrammingQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateProgrammingQuestion called==============================")
	param := mux.Vars(r)
	question := new(programming.ProgrammingQuestion)

	err := web.UnmarshalJSON(r, question)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	question.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	question.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	question.ID, err = util.ParseUUID(param["programmingQuestionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = question.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProgrammingQuestionService.UpdateProgrammingQuestion(question)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Programming question updated successfully")
}

// UpdateProgrammingQuestionIsActive updates question's isActive field to database.
func (controller *ProgrammingQuestionController) UpdateProgrammingQuestionIsActive(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateProgrammingQuestionIsActive call**************************************")

	// Create bucket.
	question := programming.ProgrammingQuestionIsActive{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &question)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting cresential id from param and parsing it to uuid.
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.ProgrammingQuestionService.UpdateProgrammingQuestionIsActive(&question, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming question updated successfully")
}

// DeleteProgrammingQuestion will delete specified programming question.
func (controller *ProgrammingQuestionController) DeleteProgrammingQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteProgrammingQuestion called==============================")
	param := mux.Vars(r)
	var err error
	question := new(programming.ProgrammingQuestion)

	question.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	question.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	question.ID, err = util.ParseUUID(param["programmingQuestionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProgrammingQuestionService.DeleteProgrammingQuestion(question)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Programming question deleted successfully")
}

// GetProgrammingQuestions will return all the questions
func (controller *ProgrammingQuestionController) GetProgrammingQuestions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetProgrammingQuestions called==============================")
	param := mux.Vars(r)
	questions := new([]programming.ProgrammingQuestionDTO)

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Parse form
	r.ParseForm()
	err = controller.ProgrammingQuestionService.GetProgrammingQuestions(questions, r.Form, tenantID, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, questions)
}

// GetProgrammingQuestionList will return all the questions
func (controller *ProgrammingQuestionController) GetProgrammingQuestionList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetProgrammingQuestionList called==============================")
	param := mux.Vars(r)
	questions := []list.ProgrammingQuestion{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	err = controller.ProgrammingQuestionService.GetProgrammingQuestionList(&questions, tenantID, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, questions)
}

// GetProgrammingQuestionsPractice will return all the questions in problem of the day format for practice.
func (controller *ProgrammingQuestionController) GetProgrammingQuestionsPractice(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetProgrammingQuestionsPractice called==============================")
	param := mux.Vars(r)
	questions := []programming.QuestionProblemOfTheDayDTO{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Fill the r.Form.
	r.ParseForm()

	err = controller.ProgrammingQuestionService.GetProgrammingQuestionsPractice(&questions, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, questions)
}

// GetProgrammingQuestion returns one question by id.
func (controller *ProgrammingQuestionController) GetProgrammingQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetProgrammingQuestion called==============================")
	param := mux.Vars(r)
	question := programming.ProgrammingQuestionWithTalentAnswerDTO{}

	// Get tenant id and parse it to uuid.
	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Get question id and parse it to uuid.
	question.ID, err = util.ParseUUID(param["programmingQuestionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	err = controller.ProgrammingQuestionService.GetProgrammingQuestion(&question, tenantID, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, question)
}

// GetProgrammingQuestionForTopic will get all programming question for specified concepts.
func (controller *ProgrammingQuestionController) GetProgrammingQuestionForTopic(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetProgrammingQuestionForTopic called==============================")
	param := mux.Vars(r)
	programmingQuestions := []programming.ProgrammingQuestionDTO{}
	parser := web.NewParser(r)

	// Get tenant id and parse it to uuid.
	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	var totalCount int

	err = controller.ProgrammingQuestionService.GetProgrammingQuestionForTopic(tenantID, &programmingQuestions, &totalCount, parser)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, programmingQuestions)
}
