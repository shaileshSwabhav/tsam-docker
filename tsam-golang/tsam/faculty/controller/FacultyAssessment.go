package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/faculty/service"
	"github.com/techlabs/swabhav/tsam/log"
	fct "github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// FacultyAssessmentController Provide method to Update, Delete, Add, Get Method For faculty.
type FacultyAssessmentController struct {
	FacultyAssessmentService *service.FacultyAssessmentService
	log                      log.Logger
	auth                     *security.Authentication
}

// NewFacultyAssessmentController Create New Instance Of FacultyAssessmentController.
func NewFacultyAssessmentController(ser *service.FacultyAssessmentService, log log.Logger, auth *security.Authentication) *FacultyAssessmentController {
	return &FacultyAssessmentController{
		FacultyAssessmentService: ser,
		log:                      log,
		auth:                     auth,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *FacultyAssessmentController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/faculty/assessments",
		controller.AddFacultyAssessments).Methods(http.MethodPost)

	// delete
	router.HandleFunc("/tenant/{tenantID}/faculty/{facultyID}/assessment",
		controller.DeleteFacultyAssessment).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/faculty/{facultyID}/assessment",
		controller.GetFacultyAssessment).Methods(http.MethodGet)

	controller.log.Info("Faculty Assessment Route Registered")
}

// AddFacultyAssessment will add assessment of faculty for specified group.
func (controller *FacultyAssessmentController) AddFacultyAssessments(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================AddFacultyAssessments called===========================")

	// param := mux.Vars(r)
	parser := web.NewParser(r)

	facultyAssessments := []fct.FacultyAssessment{}

	err := web.UnmarshalJSON(r, &facultyAssessments)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	// util.ParseUUID(param["credentialID"])
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	// util.ParseUUID(param["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	for _, assessment := range facultyAssessments {
		err = assessment.ValidateFacultyAssessment()
		if err != nil {
			controller.log.Error(err)
			web.RespondError(w, err)
			return
		}
	}

	err = controller.FacultyAssessmentService.AddFacultyAssessments(&facultyAssessments, tenantID, credentialID)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Faculty assessment successfully added")
}

// DeleteFacultyAssessment will delete specified assessment of a faculty in the table.
func (controller *FacultyAssessmentController) DeleteFacultyAssessment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================DeleteFacultyAssessment called===========================")

	// param := mux.Vars(r)
	parser := web.NewParser(r)

	var err error

	facultyAssessment := fct.FacultyAssessment{}

	// util.ParseUUID(param["credentialID"])
	facultyAssessment.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	// util.ParseUUID(param["tenantID"])
	facultyAssessment.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(param["facultyID"])
	facultyAssessment.FacultyID, err = parser.GetUUID("facultyID")
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("Unable to parse faculty id", http.StatusBadRequest))
		return
	}

	err = controller.FacultyAssessmentService.DeleteFacultyAssessment(&facultyAssessment)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Faculty assessment successfully deleted")
}

// GetFacultyAssessment will get all the assessment for specified faculty
func (controller *FacultyAssessmentController) GetFacultyAssessment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================GetFacultyAssessment called===========================")
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	facultyAssessment := []fct.FacultyAssessmentDTO{}

	// util.ParseUUID(param["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(param["facultyID"])
	facultyID, err := parser.GetUUID("facultyID")
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("Unable to parse faculty id", http.StatusBadRequest))
		return
	}

	err = controller.FacultyAssessmentService.GetFacultyAssessment(&facultyAssessment, tenantID, facultyID)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, facultyAssessment)
}
