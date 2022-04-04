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

// ProgrammingQuestionTestCaseController provide method to update, delete, add, get method for programming question test case.
type ProgrammingQuestionTestCaseController struct {
	ProgrammingQuestionTestCaseService *service.ProgrammingQuestionTestCaseService
}

// NewProgrammingQuestionTestCaseController creates new instance of ProgrammingQuestionTestCaseController.
func NewProgrammingQuestionTestCaseController(programmingQuestionTestCaseService *service.ProgrammingQuestionTestCaseService) *ProgrammingQuestionTestCaseController {
	return &ProgrammingQuestionTestCaseController{
		ProgrammingQuestionTestCaseService: programmingQuestionTestCaseService,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *ProgrammingQuestionTestCaseController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all programming question test cases by programming question id.
	router.HandleFunc("/tenant/{tenantID}/programming-question-test-case/programming-question/{questionID}",
		controller.GetProgrammingQuestionTestCases).Methods(http.MethodGet)

	// Get one programming question test case.
	router.HandleFunc("/tenant/{tenantID}/programming-question-test-case/{testCaseID}",
		controller.GetProgrammingQuestionTestCase).Methods(http.MethodGet)

	// Add one programming question test case.
	router.HandleFunc("/tenant/{tenantID}/programming-question-test-case/programming-question/{questionID}/credential/{credentialID}",
		controller.AddProgrammingQuestionTestCase).Methods(http.MethodPost)

	// Update programming question test case.
	router.HandleFunc("/tenant/{tenantID}/programming-question-test-case/{testCaseID}/programming-question/{questionID}/credential/{credentialID}",
		controller.UpdateProgrammingQuestionTestCase).Methods(http.MethodPut)

	// Delete programming question test case.
	router.HandleFunc("/tenant/{tenantID}/programming-question-test-case/{testCaseID}/credential/{credentialID}",
		controller.DeleteProgrammingQuestionTestCase).Methods(http.MethodDelete)

	log.NewLogger().Info("Programming question test case Route Registered")
}

// GetProgrammingQuestionTestCases returns all programming question test cases.
func (controller *ProgrammingQuestionTestCaseController) GetProgrammingQuestionTestCases(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetProgrammingQuestionTestCases call**************************************")

	// Create bucket.
	testcases := []programming.ProgrammingQuestionTestCase{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting programming question id from param and parsing it to uuid.
	questionID, err := util.ParseUUID(mux.Vars(r)["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question id", http.StatusBadRequest))
		return
	}

	// Parsing for query params.
	r.ParseForm()

	// Call get all programming question test cases service method.
	err = controller.ProgrammingQuestionTestCaseService.GetProgrammingQuestionTestCases(&testcases, tenantID, questionID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, testcases)
}

// GetProgrammingQuestionTestCase return the specifed programming question test case by id.
func (controller *ProgrammingQuestionTestCaseController) GetProgrammingQuestionTestCase(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetProgrammingQuestionTestCase call**************************************")

	// Create bucket.
	testcase := programming.ProgrammingQuestionTestCase{}

	// Create error variable.
	var err error

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting programming question test case id from param and parsing it to uuid.
	testcase.ID, err = util.ParseUUID(mux.Vars(r)["testCaseID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question test case id", http.StatusBadRequest))
		return
	}

	// Call get programming question test case service method.
	err = controller.ProgrammingQuestionTestCaseService.GetProgrammingQuestionTestCase(&testcase, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, testcase)
}

// AddProgrammingQuestionTestCase adds new programming question test case.
func (controller *ProgrammingQuestionTestCaseController) AddProgrammingQuestionTestCase(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddProgrammingQuestionTestCase call**************************************")

	// Create bucket.
	testCase := programming.ProgrammingQuestionTestCase{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &testCase)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = testCase.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	testCase.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	testCase.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting programming question id from param and parsing it to uuid.
	testCase.ProgrammingQuestionID, err = util.ParseUUID(mux.Vars(r)["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.ProgrammingQuestionTestCaseService.AddProgrammingQuestionTestCase(&testCase)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming question test case added successfully")
}

// UpdateProgrammingQuestionTestCase updates the specified programming question test case by id.
func (controller *ProgrammingQuestionTestCaseController) UpdateProgrammingQuestionTestCase(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateProgrammingQuestionTestCase call**************************************")

	// Create bucket.
	testCase := programming.ProgrammingQuestionTestCase{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &testCase)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate programming question test case fields.
	err = testCase.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	testCase.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting programming question test case id from param and parsing it to uuid.
	testCase.ID, err = util.ParseUUID(mux.Vars(r)["testCaseID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question test case id", http.StatusBadRequest))
		return
	}

	// Getting cresential id from param and parsing it to uuid.
	testCase.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting programming question id from param and parsing it to uuid.
	testCase.ProgrammingQuestionID, err = util.ParseUUID(mux.Vars(r)["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.ProgrammingQuestionTestCaseService.UpdateProgrammingQuestionTestCase(&testCase)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming question test case updated successfully")
}

// DeleteProgrammingQuestionTestCase delete specofic programming question test case by id.
func (controller *ProgrammingQuestionTestCaseController) DeleteProgrammingQuestionTestCase(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteProgrammingQuestionTestCase call**************************************")

	// Create bucket.
	testCase := programming.ProgrammingQuestionTestCase{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	testCase.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting programming question test case id from param and parsing it to uuid.
	testCase.ID, err = util.ParseUUID(mux.Vars(r)["testCaseID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming question test case id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	testCase.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.ProgrammingQuestionTestCaseService.DeleteProgrammingQuestionTestCase(&testCase)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Programming question test case deleted successfully")
}
