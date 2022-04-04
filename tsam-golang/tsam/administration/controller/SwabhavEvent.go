package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/administration/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// SwabhavEventController Provide method to Update, Delete, Add, Get Method For events.
type SwabhavEventController struct {
	EventService *service.SwabhavEventService
}

// NewSwabhavEventController creates new instance EventController.
func NewSwabhavEventController(eventService *service.SwabhavEventService) *SwabhavEventController {
	return &SwabhavEventController{
		EventService: eventService,
	}
}

// RegisterRoutes Register All Endpoints To Router excluding a few endpoints from token check.
func (controller *SwabhavEventController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/event/credential/{credentialID}",
		controller.AddEvent).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/event/{eventID}/credential/{credentialID}",
		controller.UpdateEvent).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/event/{eventID}/credential/{credentialID}",
		controller.DeleteEvent).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/event/limit/{limit}/offset/{offset}",
		controller.GetEvents).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/event",
		controller.GetAllEvents).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/event/{eventID}",
		controller.GetEvent).Methods(http.MethodGet)

	log.NewLogger().Info("Events Route Registered")
}

// AddEvent will add new event in the table.
func (controller *SwabhavEventController) AddEvent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddEvent call==============================")

	var event admin.SwabhavEvent
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &event)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	event.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	event.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = event.ValidateEvent()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.EventService.AddEvent(&event)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Event successfully added")
}

// UpdateEvent will udpate specified event in the table.
func (controller *SwabhavEventController) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateEvent call==============================")

	var event admin.SwabhavEvent
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &event)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	event.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	event.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	event.ID, err = util.ParseUUID(param["eventID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = event.ValidateEvent()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.EventService.UpdateEvent(&event)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Event successfully updated")
}

// DeleteEvent will delete specified event in the table.
func (controller *SwabhavEventController) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteEvent call==============================")

	var event admin.SwabhavEvent
	param := mux.Vars(r)
	var err error

	event.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	event.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	event.ID, err = util.ParseUUID(param["eventID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.EventService.DeleteEvent(&event)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Event successfully deleted")
}

// GetEvent will return specified event.
func (controller *SwabhavEventController) GetEvent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetEvent call==============================")

	event := admin.SwabhavEventDTO{}

	param := mux.Vars(r)
	var err error

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	eventID, err := util.ParseUUID(param["eventID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	err = controller.EventService.GetEvent(&event, tenantID, eventID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, event)
}

// GetEvents will get all the course-project from the table.
func (controller *SwabhavEventController) GetEvents(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetEvents call==============================")

	var events []admin.SwabhavEventDTO

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

	err = controller.EventService.GetEvents(&events, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, events)
}

// GetAllEvents will get all the course-project from the table.
func (controller *SwabhavEventController) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllEvents call==============================")

	var events []admin.SwabhavEventDTO

	param := mux.Vars(r)
	var err error

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	err = controller.EventService.GetAllEvents(&events, tenantID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, events)
}
