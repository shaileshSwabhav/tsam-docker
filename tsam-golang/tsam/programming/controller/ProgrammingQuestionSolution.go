package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/programming"
	service "github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingQuestionSolutionController provides method to update, delete, add, get all, get one for programming question solutions.
type ProgrammingQuestionSolutionController struct {
	ProgrammingQuestionSolutionService *service.ProgrammingQuestionSolutionService
}

// NewProgrammingQuestionSolutionController creates new instance of ProgrammingQuestionSolutionController.
func NewProgrammingQuestionSolutionController(programmingQuestionSolutionService *service.ProgrammingQuestionSolutionService) *ProgrammingQuestionSolutionController {
	return &ProgrammingQuestionSolutionController{
		ProgrammingQuestionSolutionService: programmingQuestionSolutionService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *ProgrammingQuestionSolutionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all programming question solution by programming question id.
	router.HandleFunc("/tenant/{tenantID}/programming-question-solution/programming-question/{questionID}",
		controller.GetProgrammingQuestionSolutions).Methods(http.MethodGet)

	// Add one programming question solution.
	router.HandleFunc("/tenant/{tenantID}/programming-question-solution/programming-question/{questionID}/credential/{credentialID}",
		controller.AddProgrammingQuestionSolution).Methods(http.MethodPost)

	// Get one programming question solution.
	router.HandleFunc("/tenant/{tenantID}/programming-question-solution/{solutionID}/programming-question/{questionID}",
		controller.GetProgrammingQuestionSolution).Methods(http.MethodGet)

	// Update one programming question solution.
	router.HandleFunc("/tenant/{tenantID}/programming-question-solution/{solutionID}/programming-question/{questionID}/credential/{credentialID}",
		controller.UpdateProgrammingQuestionSolution).Methods(http.MethodPut)

	// Delete one programming question solution.
	router.HandleFunc("/tenant/{tenantID}/programming-question-solution/{solutionID}/programming-question/{questionID}/credential/{credentialID}",
		controller.DeleteProgrammingQuestionSolution).Methods(http.MethodDelete)

	log.NewLogger().Info("Programming Question Solution Routes Registered")
}

// AddProgrammingQuestionSolution adds one programming question solution.
func (controller *ProgrammingQuestionSolutionController) AddProgrammingQuestionSolution(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddProgrammingQuestionSolution called=======================================")

	// Create bucket.
	solution := programming.ProgrammingQuestionSolution{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the programming question solution variable with given data.
	if err := web.UnmarshalJSON(r, &solution); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err := solution.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	solution.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of programming question solution.
	solution.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set programming question ID.
	solution.ProgrammingQuestionID, err = util.ParseUUID(params["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.ProgrammingQuestionSolutionService.AddProgrammingQuestionSolution(&solution); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming question solution added successfully")
}

//GetProgrammingQuestionSolutions gets all programming question solutions.
func (controller *ProgrammingQuestionSolutionController) GetProgrammingQuestionSolutions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetProgrammingQuestionSolutions called=======================================")

	// Create bucket.
	solutions := []programming.ProgrammingQuestionSolutionDTO{}

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting programming question id from param and parsing it to uuid.
	questionID, err := util.ParseUUID(params["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question id", http.StatusBadRequest))
		return
	}

	// Call get programming question solutions method.
	if err := controller.ProgrammingQuestionSolutionService.GetProgrammingQuestionSolutions(&solutions, tenantID, questionID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, solutions)
}

//GetProgrammingQuestionSolution gets one programming question solution.
func (controller *ProgrammingQuestionSolutionController) GetProgrammingQuestionSolution(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetProgrammingQuestionSolution called=======================================")

	// Create bucket.
	solution := programming.ProgrammingQuestionSolutionDTO{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set programming question solution ID.
	solution.ID, err = util.ParseUUID(params["solutionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question solution id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set programming question ID.
	solution.ProgrammingQuestionID, err = util.ParseUUID(params["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.ProgrammingQuestionSolutionService.GetProgrammingQuestionSolution(&solution, tenantID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, solution)
}

//UpdateProgrammingQuestionSolution updates programming question solution.
func (controller *ProgrammingQuestionSolutionController) UpdateProgrammingQuestionSolution(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateProgrammingQuestionSolution called=======================================")

	// Create bucket.
	solution := programming.ProgrammingQuestionSolution{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the programming question solution variable with given data.
	err := web.UnmarshalJSON(r, &solution)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := solution.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set programming question solution ID to programming question solution.
	solution.ID, err = util.ParseUUID(params["solutionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question solution id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to programming question solution.
	solution.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of programming question solution.
	solution.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set programming question solution ID to programming question solution.
	solution.ProgrammingQuestionID, err = util.ParseUUID(params["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.ProgrammingQuestionSolutionService.UpdateProgrammingQuestionSolution(&solution); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Programming question solution updated successfully")
}

//DeleteProgrammingQuestionSolution deletes one programming question solution.
func (controller *ProgrammingQuestionSolutionController) DeleteProgrammingQuestionSolution(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteProgrammingQuestionSolution called=======================================")

	// Create bucket.
	solution := programming.ProgrammingQuestionSolution{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set programming question solution ID.
	solution.ID, err = util.ParseUUID(params["solutionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question solution id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	solution.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to programming question solution's DeletedBy field.
	solution.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set programming question solution ID to programming question solution.
	solution.ProgrammingQuestionID, err = util.ParseUUID(mux.Vars(r)["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.ProgrammingQuestionSolutionService.DeleteProgrammingQuestionSolution(&solution); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming question solution deleted successfully")
}
