package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// DayController provides method to update, delete, add, get method for day.
type DayController struct {
	DayService *service.DayService
}

// NewDayController creates new instance of DayController.
func NewDayController(dayService *service.DayService) *DayController {
	return &DayController{
		DayService: dayService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *DayController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	//get all days
	router.HandleFunc("/tenant/{tenantID}/day", controller.GetDays).Methods(http.MethodGet)

	//add one day
	router.HandleFunc("/tenant/{tenantID}/day/credential/{credentialID}", controller.AddDay).Methods(http.MethodPost)

	//add multiple days
	router.HandleFunc("/tenant/{tenantID}/days/credential/{credentialID}", controller.AddDays).Methods(http.MethodPost)

	//update one day
	router.HandleFunc("/tenant/{tenantID}/day/{dayID}/credential/{credentialID}", controller.UpdateDay).Methods(http.MethodPut)

	//delete one day
	router.HandleFunc("/tenant/{tenantID}/day/{dayID}/credential/{credentialID}", controller.DeleteDay).Methods(http.MethodDelete)

	log.NewLogger().Info("Day Routes Registered")
}

// AddDay will add day to the record.
func (controller *DayController) AddDay(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddDay called==============================")
	day := &general.Day{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, day)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	day.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	day.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = day.ValidateDay()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.DayService.AddDay(day)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Day added successfully")
}

// AddDays will add multiple days to the table.
func (controller *DayController) AddDays(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddDays called==============================")
	days := &[]general.Day{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, days)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	credentialID, err := util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, day := range *days {
		fmt.Println("***day controller ->", day)
		err := day.ValidateDay()
		fmt.Println("***err controller ->", err)
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = controller.DayService.AddDays(days, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Days added successfully")
}

// UpdateDay will update the specified record in the table.
func (controller *DayController) UpdateDay(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateDay called==============================")
	day := &general.Day{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, day)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	day.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	day.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	day.ID, err = util.ParseUUID(param["dayID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = day.ValidateDay()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.DayService.UpdateDay(day)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Day updated successfully")
}

// DeleteDay will delete the specified record from table.
func (controller *DayController) DeleteDay(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteDay called==============================")
	day := &general.Day{}
	param := mux.Vars(r)
	var err error

	day.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	day.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	day.ID, err = util.ParseUUID(param["dayID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.DayService.DeleteDay(day)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Day deleted successfully")
}

// GetDays will return all days from the table.
func (controller *DayController) GetDays(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetDays called==============================")
	days := &[]general.Day{}
	param := mux.Vars(r)

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.DayService.GetDays(days, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, days)
}
