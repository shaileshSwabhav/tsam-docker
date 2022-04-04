package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	service "github.com/techlabs/swabhav/tsam/college/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	college "github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CandidateController provides method to update, delete, add, get all, get one for candidate.
type CandidateController struct {
	CandidateService *service.CandidateService
}

// NewCandidateController creates new instance of CandidateController.
func NewCandidateController(candidateService *service.CandidateService) *CandidateController {
	return &CandidateController{
		CandidateService: candidateService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *CandidateController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all candidates by campus drive id.
	router.HandleFunc("/tenant/{tenantID}/candidate/campus-drive/{campusDriveID}/limit/{limit}/offset/{offset}",
		controller.GetCandidates).Methods(http.MethodGet)

	// Add one candidate.
	router.HandleFunc("/tenant/{tenantID}/candidate/campus-drive/{campusDriveID}/credential/{credentialID}",
		controller.AddCandidate).Methods(http.MethodPost)

	// Add one candidate from campus registration form.
	addCandidate := router.HandleFunc("/tenant/{tenantID}/candidate-reg-form/campus-drive/{campusDriveID}",
		controller.AddCandidateFromRegForm).Methods(http.MethodPost)

	// Update one candidate.
	router.HandleFunc("/tenant/{tenantID}/candidate/{candidateID}/campus-drive/{campusDriveID}/credential/{credentialID}",
		controller.UpdateCandidate).Methods(http.MethodPut)

	// Delete one candidate.
	router.HandleFunc("/tenant/{tenantID}/candidate/{candidateID}/campus-drive/{campusDriveID}/credential/{credentialID}",
		controller.DeleteCandidate).Methods(http.MethodDelete)

	// Update all candidates campus talent registration field.
	router.HandleFunc("/tenant/{tenantID}/candidate/campus-drive/{campusDriveID}/credential/{credentialID}",
		controller.UpdateMultipleCandidate).Methods(http.MethodPut)

	// Exculde routes.
	*exclude = append(*exclude, addCandidate)

	log.NewLogger().Info("Candidate Routes Registered")
}

// AddCandidate adds one candidate.
func (controller *CandidateController) AddCandidate(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddCandidate called=======================================")

	// Create bucket.
	candidate := college.Candidate{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the candidate variable with given data.
	if err := web.UnmarshalJSON(r, &candidate); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Declare err.
	var err error

	// Parse and set tenant ID.
	candidate.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of candidate.
	candidate.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set campus drive ID.
	candidate.CampusDriveID, err = util.ParseUUID(params["campusDriveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.CandidateService.AddCandidateWithCredential(&candidate); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Candidate added successfully")
}

// AddCandidateFromRegForm adds one candidate from campus registration form.
func (controller *CandidateController) AddCandidateFromRegForm(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddCandidateFromRegForm called=======================================")

	// Create bucket.
	candidate := college.Candidate{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the candidate variable with given data.
	if err := web.UnmarshalJSON(r, &candidate); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Declare err.
	var err error

	// Parse and set tenant ID.
	candidate.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set campus drive ID.
	candidate.CampusDriveID, err = util.ParseUUID(params["campusDriveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.CandidateService.AddCandidate(&candidate); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Candidate added successfully")
}

// GetCandidates gets all candidates.
func (controller *CandidateController) GetCandidates(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetCandidates called=======================================")

	// Create bucket.
	candidates := []college.CandidateDTO{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the r.Form.
	r.ParseForm()

	// Create bucket for total campus drive count.
	var totalCount int

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting campus drive id from param and parsing it to uuid.
	campusDriveID, err := util.ParseUUID(params["campusDriveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive id", http.StatusBadRequest))
		return
	}

	// Get limit and offset from param and convert it to int.
	limit, offset := web.GetLimitAndOffset(r)

	// Call get candidates method.
	if err := controller.CandidateService.GetCandidates(&candidates, tenantID, campusDriveID, limit, offset, &totalCount, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, candidates)
}

// UpdateCandidate updates candidate.
func (controller *CandidateController) UpdateCandidate(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateCandidate called=======================================")

	// Create bucket.
	candidate := college.Candidate{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the candidate variable with given data.
	err := web.UnmarshalJSON(r, &candidate)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := candidate.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set candidate ID to candidate.
	candidate.CampusTalentRegistrationID, err = util.ParseUUID(params["candidateID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse candidate id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to candidate.
	candidate.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of candidate.
	candidate.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set candidate ID to candidate.
	candidate.CampusDriveID, err = util.ParseUUID(params["campusDriveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.CandidateService.UpdateCandidate(&candidate); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Candidate updated successfully")
}

// DeleteCandidate deletes one candidate.
func (controller *CandidateController) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteCandidate called=======================================")

	// Create bucket
	candidate := college.Candidate{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set candidate ID.
	candidate.CampusTalentRegistrationID, err = util.ParseUUID(params["candidateID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse candidate id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	candidate.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to candidate's DeletedBy field.
	candidate.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set candidate ID to candidate.
	candidate.CampusDriveID, err = util.ParseUUID(mux.Vars(r)["campusDriveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.CandidateService.DeleteCandidate(&candidate); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Candidate deleted successfully")
}

// UpdateMultipleCandidate updates multiple candidates' multiple fields.
func (controller *CandidateController) UpdateMultipleCandidate(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateMultipleCandidate called=======================================")

	// Create bucket.
	updateMultipleCandidate := college.UpdateMultipleCandidate{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the updateMultipleCandidate variable with given data.
	err := web.UnmarshalJSON(r, &updateMultipleCandidate)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set candidate ID to candidate.
	updateMultipleCandidate.CampusDriveID, err = util.ParseUUID(mux.Vars(r)["campusDriveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	updateMultipleCandidate.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	updateMultipleCandidate.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.CandidateService.UpdateMultipleCandidate(&updateMultipleCandidate)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Candidates updated successfully")
}
