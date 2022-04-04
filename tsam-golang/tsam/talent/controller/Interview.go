package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	general "github.com/techlabs/swabhav/tsam/models/general"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	service "github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// InterviewController provides method to update, delete, add, get all, get one for interviews.
type InterviewController struct {
	InterviewService *service.InterviewService
}

// NewInterviewController creates new instance of InterviewController.
func NewInterviewController(interviewservice *service.InterviewService) *InterviewController {
	return &InterviewController{
		InterviewService: interviewservice,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *InterviewController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all interviews by interview schedule id.
	router.HandleFunc("/tenant/{tenantID}/interview/interview-schedule/{interviewScheduleID}",
		controller.GetInterviews).Methods(http.MethodGet)

	// Add one interview.
	router.HandleFunc("/tenant/{tenantID}/interview/interview-schedule/{interviewScheduleID}/credential/{credentialID}",
		controller.AddInterview).Methods(http.MethodPost)

	// Get one interview.
	router.HandleFunc("/tenant/{tenantID}/interview",
		controller.GetInterview).Methods(http.MethodGet)

	// Update one interview.
	router.HandleFunc("/tenant/{tenantID}/interview/{interviewID}/interview-schedule/{interviewScheduleID}/credential/{credentialID}",
		controller.UpdateInterview).Methods(http.MethodPut)

	// Delete one interview.
	router.HandleFunc("/tenant/{tenantID}/interview/{interviewID}/interview-schedule/{interviewScheduleID}/credential/{credentialID}",
		controller.DeleteInterview).Methods(http.MethodDelete)

	// Get interviewers list.
	router.HandleFunc("/tenant/{tenantID}/interview/interviewers-list",
		controller.GetAllInterviewersList).Methods(http.MethodGet)

	log.NewLogger().Info("Interview Routes Registered")
}

// AddInterview adds one interview.
func (controller *InterviewController) AddInterview(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddInterview called=======================================")

	// Create bucket.
	interview := tal.Interview{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the interview variable with given data.
	if err := web.UnmarshalJSON(r, &interview); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err := interview.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	interview.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of interview.
	interview.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set interview schedule ID.
	interview.ScheduleID, err = util.ParseUUID(params["interviewScheduleID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview schedule id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.InterviewService.AddInterview(&interview); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Interview added successfully")
}

//GetInterviews gets all interviews.
func (controller *InterviewController) GetInterviews(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetInterviews called=======================================")

	// Create bucket.
	interviews := []tal.Interview{}

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting interview schedule id from param and parsing it to uuid.
	interviewScheduleID, err := util.ParseUUID(params["interviewScheduleID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview schedule id", http.StatusBadRequest))
		return
	}

	// Call get interviews method.
	if err := controller.InterviewService.GetInterviews(&interviews, tenantID, interviewScheduleID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, interviews)
}

//GetInterview gets one interview.
func (controller *InterviewController) GetInterview(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetInterview called=======================================")

	// Create bucket.
	interview := tal.Interview{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parsing for query params.
	r.ParseForm()

	// Parse and set tenant ID.
	interview.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.InterviewService.GetInterview(&interview,r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, interview)
}

//UpdateInterview updates interview.
func (controller *InterviewController) UpdateInterview(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateInterview called=======================================")

	// Create bucket.
	interview := tal.Interview{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the interview variable with given data.
	err := web.UnmarshalJSON(r, &interview)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := interview.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set interview ID to interview.
	interview.ID, err = util.ParseUUID(params["interviewID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to interview.
	interview.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of interview.
	interview.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set interview schedule ID to interview.
	interview.ScheduleID, err = util.ParseUUID(params["interviewScheduleID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview schedule id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.InterviewService.UpdateInterview(&interview); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Interview updated successfully")
}

//DeleteInterview deletes one interview.
func (controller *InterviewController) DeleteInterview(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteInterview called=======================================")

	// Create bucket.
	interview := tal.Interview{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set interview ID.
	interview.ID, err = util.ParseUUID(params["interviewID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	interview.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to interview's DeletedBy field.
	interview.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set interview schedule ID to interview.
	interview.ScheduleID, err = util.ParseUUID(mux.Vars(r)["interviewScheduleID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview schedule id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.InterviewService.DeleteInterview(&interview); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Interview deleted successfully")
}

//GetAllInterviewersList gets all interviewers list.
func (controller *InterviewController) GetAllInterviewersList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetAllInterviewersList called=======================================")

	// Create bucket.
	interviewers := []general.Credential{}

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get interviews method.
	if err := controller.InterviewService.GetAllInterviewersList(&interviewers, tenantID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, interviewers)
}
