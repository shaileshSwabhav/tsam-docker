package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	service "github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// InterviewRound provides method to update, delete, add, get all, get one for interview rounds.
type InterviewRoundController struct {
	InterviewRoundService *service.InterviewRoundService
}

// NewInterviewRoundController creates new instance of InterviewRoundController.
func NewInterviewRoundController(interviewRoundService *service.InterviewRoundService) *InterviewRoundController {
	return &InterviewRoundController{
		InterviewRoundService: interviewRoundService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *InterviewRoundController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all interview rounds.
	router.HandleFunc("/tenant/{tenantID}/interview-round",
		controller.GetInterviewRounds).Methods(http.MethodGet)

	// Add one interview round.
	router.HandleFunc("/tenant/{tenantID}/interview-round/credential/{credentialID}",
		controller.AddInterviewRound).Methods(http.MethodPost)

	// Add multiple interview round.
	router.HandleFunc("/tenant/{tenantID}/interview-rounds/credential/{credentialID}",
		controller.AddInterviewRounds).Methods(http.MethodPost)

	// Get one interview round.
	router.HandleFunc("/tenant/{tenantID}/interview-round/{interviewRoundID}",
		controller.GetInterviewRound).Methods(http.MethodGet)

	// Update one interview round.
	router.HandleFunc("/tenant/{tenantID}/interview-round/{interviewRoundID}/credential/{credentialID}",
		controller.UpdateInterviewRound).Methods(http.MethodPut)

	// Delete one interview round.
	router.HandleFunc("/tenant/{tenantID}/interview-round/{interviewRoundID}/credential/{credentialID}",
		controller.DeleteInterviewRound).Methods(http.MethodDelete)

	log.NewLogger().Info("InterviewRound Routes Registered")
}

// AddInterviewRound adds one interview round.
func (controller *InterviewRoundController) AddInterviewRound(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddInterviewRound called=======================================")

	// Create bucket.
	interviewRound := tal.InterviewRound{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the interviewRound variable with given data.
	if err := web.UnmarshalJSON(r, &interviewRound); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err := interviewRound.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	interviewRound.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of interviewRound.
	interviewRound.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.InterviewRoundService.AddInterviewRound(&interviewRound); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Interview round added successfully")
}

//GetInterviewRounds gets all interview rounds.
func (controller *InterviewRoundController) GetInterviewRounds(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetInterviewRounds called=======================================")

	// Create bucket.
	interviewRounds := []tal.InterviewRound{}

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get interview rounds method.
	if err := controller.InterviewRoundService.GetInterviewRounds(&interviewRounds, tenantID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, interviewRounds)
}

//GetInterviewRound gets one interview round.
func (controller *InterviewRoundController) GetInterviewRound(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetInterviewRound called=======================================")

	// Create bucket.
	interviewRound := tal.InterviewRound{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set interviewRound ID.
	interviewRound.ID, err = util.ParseUUID(params["interviewRoundID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview round id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	interviewRound.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.InterviewRoundService.GetInterviewRound(&interviewRound); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, interviewRound)
}

//UpdateInterviewRound updates interview round.
func (controller *InterviewRoundController) UpdateInterviewRound(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateInterviewRound called=======================================")

	// Create bucket.
	interviewRound := tal.InterviewRound{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the interviewRound variable with given data.
	err := web.UnmarshalJSON(r, &interviewRound)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := interviewRound.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set interviewRound ID to interviewRound.
	interviewRound.ID, err = util.ParseUUID(params["interviewRoundID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview round id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to interviewRound.
	interviewRound.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of interviewRound.
	interviewRound.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.InterviewRoundService.UpdateInterviewRound(&interviewRound); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Interview round updated successfully")
}

//DeleteInterviewRound deletes one interview round.
func (controller *InterviewRoundController) DeleteInterviewRound(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteInterviewRound called=======================================")

	// Create bucket.
	interviewRound := tal.InterviewRound{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set interviewRound ID.
	interviewRound.ID, err = util.ParseUUID(params["interviewRoundID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview round id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	interviewRound.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to interviewRound's DeletedBy field.
	interviewRound.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.InterviewRoundService.DeleteInterviewRound(&interviewRound); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Interview round deleted successfully")
}

// AddInterviewRounds adds multiple interview rounds.
func (controller *InterviewRoundController) AddInterviewRounds(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddInterviewRounds called=======================================")

	// Create bucket for interview round ids.
	interviewRoundIDs := []uuid.UUID{}

	// Create bucket for interview rounds.
	interviewRounds := []tal.InterviewRound{}

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to interviewRound's DeletedBy field.
	credentialID, err := util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Fill the interviewRounds variable with given data.
	err = web.UnmarshalJSON(r, &interviewRounds)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields
	for _, interviewRound := range interviewRounds {
		err = interviewRound.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	// Call add multiple service method.
	err = controller.InterviewRoundService.AddInterviewRounds(&interviewRounds, &interviewRoundIDs, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Interview rounds added successfully")
}
