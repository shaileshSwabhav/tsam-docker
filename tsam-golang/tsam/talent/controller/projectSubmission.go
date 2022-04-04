package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentProjectSubmissionController provides methods to do Update, Delete, Add, Get operations on talent_assignment_submissions.
type TalentProjectSubmissionController struct {
	log     log.Logger
	service *service.TalentProjectSubmissionService
	auth    *security.Authentication
}

// NewTalentProjectSubmissionController creates new instance of TalentProjectSubmissionController.
func NewTalentProjectSubmissionController(service *service.TalentProjectSubmissionService,
	log log.Logger, auth *security.Authentication) *TalentProjectSubmissionController {
	return &TalentProjectSubmissionController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *TalentProjectSubmissionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add talent project submission.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/project/{projectID}/talent/{talentID}/project-submission",
		controller.AddTalentProjectSubmission).Methods(http.MethodPost)

	//Update talent project submission
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/project/{projectID}/talent/{talentID}/project-submission/{projectSubmissionID}",
		controller.UpdateTalentProjectSubmission).Methods(http.MethodPut)

	//Update talent project submission with score
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/project/{projectID}/talent/{talentID}/project-submission/{projectSubmissionID}/score",
		controller.ScoreTalentProjectSubmission).Methods(http.MethodPut)

	// Get batch project for one talent.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}/project-submission",
		controller.GetAllProjectsSubmissionsOfTalent).Methods(http.MethodGet)

	// get talent project submission
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/project/{projectID}/talent/{talentID}/project-submission",
		controller.GetTalentProjectSubmissions).Methods(http.MethodGet)

	// future endpoint : ("/tenant/{tenantID}/batch/{batchID}/faculty/{facultyID}/topic-assignments"
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/projects",
		controller.GetAllProjectsWithSubmissions).Methods(http.MethodGet)

	controller.log.Info("Talent Project Submissions Routes Registered")
}

// AddTalentProjectSubmission will add project submission of talent.
func (controller *TalentProjectSubmissionController) AddTalentProjectSubmission(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add Talent Project Submission Called==============================")
	talentProjectSubmission := talent.ProjectSubmission{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &talentProjectSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	talentProjectSubmission.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentProjectSubmission.TalentID, err = parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set batch ID.
	talentProjectSubmission.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Parse and set project ID.
	talentProjectSubmission.BatchProjectID, err = parser.GetUUID("projectID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	talentProjectSubmission.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = talentProjectSubmission.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddTalentProjectSubmission(&talentProjectSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Project successfully submitted")
}

//UpdateTalentSubmissions
func (controller *TalentProjectSubmissionController) UpdateTalentProjectSubmission(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Update Talent Submission Called==============================")
	talentProjectSubmission := talent.ProjectSubmission{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &talentProjectSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	talentProjectSubmission.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentProjectSubmission.TalentID, err = parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set batch ID.
	talentProjectSubmission.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Parse and set project ID.
	talentProjectSubmission.BatchProjectID, err = parser.GetUUID("projectID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Parse and set ID.
	talentProjectSubmission.ID, err = parser.GetUUID("projectSubmissionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse submission id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	talentProjectSubmission.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = talentProjectSubmission.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdateTalentProjectSubmission(&talentProjectSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Project successfully updated")
}

// GetTalentSubmissions will return talent submissions.
func (controller *TalentProjectSubmissionController) GetTalentProjectSubmissions(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Talent Submission Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch ID.
	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	projectID, err := parser.GetUUID("projectID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	var talentSubmissions []talent.ProjectSubmissionDTO
	err = controller.service.GetTalentProjectSubmissions(tenantID, batchID, talentID, projectID, &talentSubmissions, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, talentSubmissions)
}

// GetAllProjectsSubmissionsOfTalent will return batch project for one talent.
func (controller *TalentProjectSubmissionController) GetAllProjectsSubmissionsOfTalent(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("============================== GetAllProjectsSubmissionsOfTalent Call ==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	var projects []batch.ProjectDTO
	err = controller.service.GetAllProjectsSubmissionsOfTalent(tenantID, batchID, talentID, &projects, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, projects)
}

// ScoreTalentSubmission *not tested once.
func (controller *TalentProjectSubmissionController) ScoreTalentProjectSubmission(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================ScoreTalentSubmission Called==============================")
	talentProjectSubmission := talent.ProjectSubmission{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &talentProjectSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	talentProjectSubmission.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentProjectSubmission.TalentID, err = parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set batch ID.
	talentProjectSubmission.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Parse and set project ID.
	talentProjectSubmission.BatchProjectID, err = parser.GetUUID("projectID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Parse and set ID.
	talentProjectSubmission.ID, err = parser.GetUUID("projectSubmissionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse submission id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	talentProjectSubmission.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = talentProjectSubmission.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	facultyID, err := controller.auth.ExtractIDFromToken(r, "loginID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	talentProjectSubmission.FacultyID = &facultyID

	err = talentProjectSubmission.ValidateScore()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.ScoreTalentProjectSubmission(&talentProjectSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusNoContent, nil)
}

// GetAllProjectsWithSubmissions will return session_project list.
func (controller *TalentProjectSubmissionController) GetAllProjectsWithSubmissions(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("============================== GetAllProjectsWithSubmissions Call ==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	var projects []batch.ProjectDTO
	err = controller.service.GetAllProjectsWithSubmissions(tenantID, batchID, &projects, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, projects)
}
