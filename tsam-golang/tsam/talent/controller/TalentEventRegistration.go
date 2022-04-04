package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentEventRegistrationController Provide method to Update, Delete, Add, Get Method For talent_registration_event.
type TalentEventRegistrationController struct {
	TalentEventRegistrationService *service.TalentEventRegistrationService
}

// NewTalentEventRegistrationController creates new instance TalentEventRegistrationController.
func NewTalentEventRegistrationController(talentEventRegistrationService *service.TalentEventRegistrationService) *TalentEventRegistrationController {
	return &TalentEventRegistrationController{
		TalentEventRegistrationService: talentEventRegistrationService,
	}
}

// RegisterRoutes Register All Endpoints To Router excluding a few endpoints from token check.
func (controller *TalentEventRegistrationController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/talent-event-registration/credential/{credentialID}",
		controller.AddTalentRegistration).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/talent-event-registration/{registrationID}/credential/{credentialID}",
		controller.UpdateTalentRegistration).Methods(http.MethodPut)

	// get
	router.HandleFunc("/tenant/{tenantID}/talent-event-registration/limit/{limit}/offset/{offset}",
		controller.GetTalentRegistrations).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/talent-event-registration",
		controller.GetTalentRegistration).Methods(http.MethodGet)

	log.NewLogger().Info("Events Route Registered")
}

// AddTalentRegistration will add new event in the table.
func (controller *TalentEventRegistrationController) AddTalentRegistration(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddTalentRegistration call==============================")

	var registration talent.TalentEventRegistration
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &registration)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	registration.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	registration.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = registration.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TalentEventRegistrationService.AddTalentRegistration(&registration)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "You have successfully registered for the event.")
}

// UpdateTalentRegistration will add new event in the table.
func (controller *TalentEventRegistrationController) UpdateTalentRegistration(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateTalentRegistration call==============================")

	var registration talent.TalentEventRegistration
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &registration)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	registration.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	registration.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	registration.ID, err = util.ParseUUID(param["registrationID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = registration.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TalentEventRegistrationService.UpdateTalentRegistration(&registration)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Registration details updated")
}

// GetTalentRegistration will for specified talent and event.
func (controller *TalentEventRegistrationController) GetTalentRegistration(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentRegistration call==============================")

	registrations := &talent.TalentEventRegistrationDTO{}

	param := mux.Vars(r)
	var err error

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	err = controller.TalentEventRegistrationService.GetTalentRegistration(registrations, tenantID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, registrations)
}

// GetTalentRegistrations will return all the talents registered with limit and offset.
func (controller *TalentEventRegistrationController) GetTalentRegistrations(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentRegistrations call==============================")

	var registrations []talent.TalentEventRegistrationDTO

	param := mux.Vars(r)
	var err error

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	err = controller.TalentEventRegistrationService.GetTalentRegistrations(&registrations, tenantID, r.Form,
		limit, offset, &totalCount)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, registrations)
}

// // GetAllTalentRegistration will return all the talents registered.
// func (controller *TalentEventRegistrationController) GetAllTalentRegistration(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetAllTalentRegistration call==============================")

// 	var registrations []talent.TalentEventRegistrationDTO

// 	param := mux.Vars(r)
// 	var err error

// 	tenantID, err := util.ParseUUID(param["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	r.ParseForm()

// 	err = controller.TalentEventRegistrationService.GetAllTalentRegistration(&registrations, tenantID, r.Form)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, registrations)
// }
