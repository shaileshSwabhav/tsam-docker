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

// TalentAssignmentSubmissionController provides methods to do Update, Delete, Add, Get operations on talent_assignment_submissions.
type TalentAssignmentSubmissionController struct {
	log     log.Logger
	service *service.TalentAssignmentSubmissionService
	auth    *security.Authentication
}

// NewTalentAssignmentSubmissionController creates new instance of TalentAssignmentSubmissionController.
func NewTalentAssignmentSubmissionController(service *service.TalentAssignmentSubmissionService,
	log log.Logger, auth *security.Authentication) *TalentAssignmentSubmissionController {
	return &TalentAssignmentSubmissionController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *TalentAssignmentSubmissionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add talent asssignment submission.
	router.HandleFunc("/tenant/{tenantID}/topic-assignment/{topicAssignmentID}/talent/{talentID}/talent-submission",
		controller.AddTalentSubmission).Methods(http.MethodPost)

	// Get batch topic assignments for one talent.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}/talent-submission",
		controller.GetAllAssignmentsSubmissionsOfTalent).Methods(http.MethodGet)

	// get /tenant/{id}/topic-assignments/{id}/talents/{id}/submissions
	router.HandleFunc("/tenant/{tenantID}/topic-assignments/{topicAssignmentID}/talents/{talentID}/submissions",
		controller.GetTalentSubmissions).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent-score",
		controller.GetTalentScores).Methods(http.MethodGet)

	//Update
	// /tenant/{id}/topic-assignments/{id}/talents/{id}/submissions/{id}
	router.HandleFunc("/tenant/{tenantID}/topic-assignments/{topicAssignmentID}/talents/{talentID}/submissions/{submissionID}",
		controller.UpdateTalentSubmissions).Methods(http.MethodPut)

	// /tenant/{id}/topic-assignments/{id}/talents/{id}/submissions/{id}/score
	router.HandleFunc("/tenant/{tenantID}/topic-assignments/{assignmentID}/talents/{talentID}/submissions/{submissionID}/score",
		controller.ScoreTalentSubmission).Methods(http.MethodPut)

	controller.log.Info("Talent Assignment Submissions Routes Registered")
}

// AddTalentSubmission will add assignment submission of talent.
func (controller *TalentAssignmentSubmissionController) AddTalentSubmission(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add Talent Submission Called==============================")
	talentSubmission := talent.AssignmentSubmission{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &talentSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	talentSubmission.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentSubmission.TalentID, err = parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set batch session ID.
	talentSubmission.BatchTopicAssignmentID, err = parser.GetUUID("topicAssignmentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch session assignment id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	talentSubmission.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = talentSubmission.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddTalentSubmission(&talentSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Assignment successfully submitted")
}

// GetTalentScores will return scores of all talents present in batch.
func (controller *TalentAssignmentSubmissionController) GetTalentScores(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Talent Score Called==============================")
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

	var assignments []talent.TalentAssignmentScoreDTO

	err = controller.service.GetTalentScores(tenantID, batchID, &assignments, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, assignments)
}

// GetTalentSubmissions will return talent submissions.
func (controller *TalentAssignmentSubmissionController) GetTalentSubmissions(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Talent Submission Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch ID.
	topicAssignmentID, err := parser.GetUUID("topicAssignmentID")
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

	var talentSubmissions []talent.AssignmentSubmissionDTO
	err = controller.service.GetTalentSubmissions(tenantID, topicAssignmentID, talentID, &talentSubmissions, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, talentSubmissions)
}

//UpdateTalentSubmissions
func (controller *TalentAssignmentSubmissionController) UpdateTalentSubmissions(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Update Talent Submission Called==============================")
	talentSubmission := talent.AssignmentSubmission{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &talentSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	talentSubmission.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch session ID.
	talentSubmission.BatchTopicAssignmentID, err = parser.GetUUID("topicAssignmentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch session assignment id", http.StatusBadRequest))
		return
	}

	// Parse and set ID.
	talentSubmission.ID, err = parser.GetUUID("submissionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse submission id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	talentSubmission.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = talentSubmission.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdateTalentSubmission(&talentSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Assignment successfully updated")
}

// ScoreTalentSubmission *not tested once.
func (controller *TalentAssignmentSubmissionController) ScoreTalentSubmission(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================ScoreTalentSubmission Called==============================")
	talentSubmission := talent.AssignmentSubmission{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &talentSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	talentSubmission.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch session ID.
	talentSubmission.BatchTopicAssignmentID, err = parser.GetUUID("assignmentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch session assignment id", http.StatusBadRequest))
		return
	}

	// Parse and set ID.
	talentSubmission.TalentID, err = parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set ID.
	talentSubmission.ID, err = parser.GetUUID("submissionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse submission id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	talentSubmission.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
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
	talentSubmission.FacultyID = &facultyID

	err = talentSubmission.ValidateScore()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.ScoreTalentSubmission(&talentSubmission)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusNoContent, nil)
}

// GetAllAssignmentsSubmissionsOfTalent will return batch topic assignments for one talent.
func (controller *TalentAssignmentSubmissionController) GetAllAssignmentsSubmissionsOfTalent(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("============================== GetAllAssignmentsSubmissionsOfTalent Call ==============================")
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

	var assignments []batch.TopicAssignmentDTO
	err = controller.service.GetAllAssignmentsSubmissionsOfTalent(tenantID, batchID, talentID, &assignments, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, assignments)
}
