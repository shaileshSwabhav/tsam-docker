package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingQuestionTalentAnswerController provides method to update, delete, add, get method for programming question talent answer.
type ProgrammingQuestionTalentAnswerController struct {
	ProgrammingQuestionTalentAnswerService *service.ProgrammingQuestionTalentAnswerService
}

// NewProgrammingQuestionTalentAnswerController creates new instance of ProgrammingQuestionTalentAnswerController.
func NewProgrammingQuestionTalentAnswerController(programmingQuestionTalentAnswerService *service.ProgrammingQuestionTalentAnswerService) *ProgrammingQuestionTalentAnswerController {
	return &ProgrammingQuestionTalentAnswerController{
		ProgrammingQuestionTalentAnswerService: programmingQuestionTalentAnswerService,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *ProgrammingQuestionTalentAnswerController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all programming question talent answers with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/programming-question-talent-answer/limit/{limit}/offset/{offset}",
		controller.GetProgrammingQuestionTalentAnswers).Methods(http.MethodGet)

	// Get one programming question talent answer.
	router.HandleFunc("/tenant/{tenantID}/programming-question-talent-answer/{answerID}",
		controller.GetProgrammingQuestionTalentAnswer).Methods(http.MethodGet)

	// Add one programming question talent answer.
	router.HandleFunc("/tenant/{tenantID}/programming-question-talent-answer/credential/{credentialID}",
		controller.AddProgrammingQuestionTalentAnswer).Methods(http.MethodPost)

	// Update programming question talent answer.
	router.HandleFunc("/tenant/{tenantID}/programming-question-talent-answer/{answerID}/credential/{credentialID}",
		controller.UpdateProgrammingQuestionTalentAnswer).Methods(http.MethodPut)

	// Update score and isCorrect of programming question talent answer.
	router.HandleFunc("/tenant/{tenantID}/programming-question-talent-answer/score/{answerID}/credential/{credentialID}",
		controller.UpdateProgrammingQuestionTalentAnswerScore).Methods(http.MethodPut)

	// Delete programming question talent answer.
	router.HandleFunc("/tenant/{tenantID}/programming-question-talent-answer/{answerID}/credential/{credentialID}",
		controller.DeleteProgrammingQuestionTalentAnswer).Methods(http.MethodDelete)

	// Add one solution is viewed.
	router.HandleFunc("/tenant/{tenantID}/programming-question-solution-is-viewed/credential/{credentialID}",
		controller.AddProgrammingQuestionSolutionIsViewed).Methods(http.MethodPost)

	log.NewLogger().Info("ProgrammingQuestionTalentAnswer Route Registered")
}

// GetProgrammingQuestionTalentAnswers returns all programming question talent answers.
func (controller *ProgrammingQuestionTalentAnswerController) GetProgrammingQuestionTalentAnswers(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetProgrammingQuestionTalentAnswers call**************************************")

	// Create bucket.
	answers := []programming.ProgrammingQuestionTalentAnswerDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parsing for query params.
	r.ParseForm()

	// For pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all  programming question talent answers service method.
	err = controller.ProgrammingQuestionTalentAnswerService.GetProgrammingQuestionTalentAnswers(&answers, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, answers)
}

// GetProgrammingQuestionTalentAnswer return the specifed  programming question talent answer by id.
func (controller *ProgrammingQuestionTalentAnswerController) GetProgrammingQuestionTalentAnswer(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetProgrammingQuestionTalentAnswer call**************************************")

	// Create bucket.
	answer := programming.ProgrammingQuestionTalentAnswerWithFullQuestionDTO{}

	// Create error variable.
	var err error

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting programming question talent answer id from param and parsing it to uuid.
	answer.ID, err = util.ParseUUID(mux.Vars(r)["answerID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question talent answer id", http.StatusBadRequest))
		return
	}

	// Call get programming question talent answer service method.
	err = controller.ProgrammingQuestionTalentAnswerService.GetProgrammingQuestionTalentAnswer(&answer, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, answer)
}

// AddProgrammingQuestionTalentAnswer adds new programming question talent answer.
func (controller *ProgrammingQuestionTalentAnswerController) AddProgrammingQuestionTalentAnswer(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddProgrammingQuestionTalentAnswer call**************************************")

	// Create bucket.
	answer := programming.ProgrammingQuestionTalentAnswer{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &answer)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = answer.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	answer.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	answer.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.ProgrammingQuestionTalentAnswerService.AddProgrammingQuestionTalentAnswer(&answer)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming question talent answer added successfully")
}

// AddProgrammingQuestionSolutionIsViewed adds new solution is viewed entry to database.
func (controller *ProgrammingQuestionTalentAnswerController) AddProgrammingQuestionSolutionIsViewed(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddProgrammingQuestionSolutionIsViewed call**************************************")

	// Create bucket.
	solutionIsviewed := programming.ProgrammingQuestionSolutionIsViewed{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &solutionIsviewed)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = solutionIsviewed.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	solutionIsviewed.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	solutionIsviewed.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.ProgrammingQuestionTalentAnswerService.AddProgrammingQuestionSolutionIsViewed(&solutionIsviewed)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming question solution is viewed added successfully")
}

// UpdateProgrammingQuestionTalentAnswer updates the specified programming question talent answer by id.
func (controller *ProgrammingQuestionTalentAnswerController) UpdateProgrammingQuestionTalentAnswer(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateProgrammingQuestionTalentAnswer call**************************************")

	// Create bucket.
	answer := programming.ProgrammingQuestionTalentAnswer{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &answer)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate programming question talent answer fields.
	err = answer.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	answer.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting programming question talent answer id from param and parsing it to uuid.
	answer.ID, err = util.ParseUUID(mux.Vars(r)["answerID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question talent answer id", http.StatusBadRequest))
		return
	}

	// Getting cresential id from param and parsing it to uuid.
	answer.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.ProgrammingQuestionTalentAnswerService.UpdateProgrammingQuestionTalentAnswer(&answer)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming question talent answer updated successfully")
}

// UpdateProgrammingQuestionTalentAnswerScore updates programming question talent answer's score and isCorrect field to database.
func (controller *ProgrammingQuestionTalentAnswerController) UpdateProgrammingQuestionTalentAnswerScore(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateProgrammingQuestionTalentAnswerScore call**************************************")

	// Create bucket.
	answer := programming.ProgrammingQuestionTalentAnswerScore{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &answer)
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

	// Getting programming question talent answer id from param and parsing it to uuid.
	answer.ID, err = util.ParseUUID(mux.Vars(r)["answerID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question talent answer id", http.StatusBadRequest))
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
	err = controller.ProgrammingQuestionTalentAnswerService.UpdateProgrammingQuestionTalentAnswerScore(&answer, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming question talent answer updated successfully")
}

// DeleteProgrammingQuestionTalentAnswer delete specofic programming question talent answer by id.
func (controller *ProgrammingQuestionTalentAnswerController) DeleteProgrammingQuestionTalentAnswer(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteProgrammingQuestionTalentAnswer call**************************************")

	// Create bucket.
	answer := programming.ProgrammingQuestionTalentAnswer{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	answer.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting programming question talent answer id from param and parsing it to uuid.
	answer.ID, err = util.ParseUUID(mux.Vars(r)["answerID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question talent answer id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	answer.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.ProgrammingQuestionTalentAnswerService.DeleteProgrammingQuestionTalentAnswer(&answer)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Programming question talent answer deleted successfully")
}
