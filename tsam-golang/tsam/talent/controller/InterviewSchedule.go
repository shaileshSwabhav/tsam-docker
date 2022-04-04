package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	service "github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// InterviewShceduleController provides method to update, delete, add, get all, get one for interview schedules.
type InterviewShceduleController struct {
	InterviewScheduleService *service.InterviewScheduleService
}

// NewInterviewShceduleController creates new instance of InterviewShceduleController.
func NewInterviewShceduleController(interviewScheduleService *service.InterviewScheduleService) *InterviewShceduleController {
	return &InterviewShceduleController{
		InterviewScheduleService: interviewScheduleService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *InterviewShceduleController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all interview schedules by talent id.
	router.HandleFunc("/tenant/{tenantID}/interview-schedule/talent/{talentID}",
		controller.GetInterviewSchedules).Methods(http.MethodGet)

	// Add one interview schedule.
	router.HandleFunc("/tenant/{tenantID}/interview-schedule/talent/{talentID}/credential/{credentialID}",
		controller.AddInterviewSchedule).Methods(http.MethodPost)

	// Get one interview schedule.
	router.HandleFunc("/tenant/{tenantID}/interview-schedule/{interviewScheduleID}/talent/{talentID}",
		controller.GetInterviewSchedule).Methods(http.MethodGet)

	// Update one interview schedule.
	router.HandleFunc("/tenant/{tenantID}/interview-schedule/{interviewScheduleID}/talent/{talentID}/credential/{credentialID}",
		controller.UpdateInterviewSchedule).Methods(http.MethodPut)

	// Delete one interview schedule.
	router.HandleFunc("/tenant/{tenantID}/interview-schedule/{interviewScheduleID}/talent/{talentID}/credential/{credentialID}",
		controller.DeleteInterviewSchedule).Methods(http.MethodDelete)

	log.NewLogger().Info("InterviewSchedule Routes Registered")
}

// AddInterviewSchedule adds one interview schedule.
func (controller *InterviewShceduleController) AddInterviewSchedule(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddInterviewSchedule called=======================================")

	// Create bucket.
	interviewSchedule := tal.InterviewSchedule{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the interviewSchedule variable with given data.
	if err := web.UnmarshalJSON(r, &interviewSchedule); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err := interviewSchedule.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	interviewSchedule.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of interviewSchedule.
	interviewSchedule.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	interviewSchedule.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.InterviewScheduleService.AddInterviewSchedule(&interviewSchedule); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Interview schedule added successfully")
}

//GetInterviewSchedules gets all interview schedules.
func (controller *InterviewShceduleController) GetInterviewSchedules(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetInterviewSchedules called=======================================")

	// Create bucket.
	interviewSchedules := []tal.InterviewSchedule{}

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting talent id from param and parsing it to uuid.
	talentID, err := util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call get interview schedules method.
	if err := controller.InterviewScheduleService.GetInterviewSchedules(&interviewSchedules, tenantID, talentID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, interviewSchedules)
}

//GetInterviewSchedule gets one interview schedule.
func (controller *InterviewShceduleController) GetInterviewSchedule(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetInterviewSchedule called=======================================")

	// Create bucket.
	interviewSchedule := tal.InterviewSchedule{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set interviewSchedule ID.
	interviewSchedule.ID, err = util.ParseUUID(params["interviewScheduleID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview schedule id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	interviewSchedule.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	interviewSchedule.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.InterviewScheduleService.GetInterviewSchedule(&interviewSchedule); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, interviewSchedule)
}

//UpdateInterviewSchedule updates interview schedule.
func (controller *InterviewShceduleController) UpdateInterviewSchedule(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateInterviewSchedule called=======================================")

	// Create bucket.
	interviewSchedule := tal.InterviewSchedule{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the interviewSchedule variable with given data.
	err := web.UnmarshalJSON(r, &interviewSchedule)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := interviewSchedule.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set interviewSchedule ID to interviewSchedule.
	interviewSchedule.ID, err = util.ParseUUID(params["interviewScheduleID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview schedule id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to interviewSchedule.
	interviewSchedule.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of interviewSchedule.
	interviewSchedule.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set interviewSchedule ID to interviewSchedule.
	interviewSchedule.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.InterviewScheduleService.UpdateInterviewSchedule(&interviewSchedule); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Interview schedule updated successfully")
}

//DeleteInterviewSchedule deletes one interview schedule.
func (controller *InterviewShceduleController) DeleteInterviewSchedule(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteInterviewSchedule called=======================================")

	// Create bucket.
	interviewSchedule := tal.InterviewSchedule{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set interviewSchedule ID.
	interviewSchedule.ID, err = util.ParseUUID(params["interviewScheduleID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse interview schedule id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	interviewSchedule.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to interviewSchedule's DeletedBy field.
	interviewSchedule.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set interviewSchedule ID to interviewSchedule.
	interviewSchedule.TalentID, err = util.ParseUUID(mux.Vars(r)["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call delete service method
	if err := controller.InterviewScheduleService.DeleteInterviewSchedule(&interviewSchedule); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Interview schedule deleted successfully")
}
