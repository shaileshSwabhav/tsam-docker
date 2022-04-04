package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// NextActionTypeController provides method to update, delete, add, get method for nextActionType.
type NextActionTypeController struct {
	NextActionTypeService *service.NextActionTypeService
}

// NewNextActionTypeController creates new instance of NextActionTypeController.
func NewNextActionTypeController(nextActionTypeService *service.NextActionTypeService) *NextActionTypeController {
	return &NextActionTypeController{
		NextActionTypeService: nextActionTypeService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *NextActionTypeController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all nextActionTypes.
	router.HandleFunc("/tenant/{tenantID}/next-action-type",
		controller.GetNextActionTypes).Methods(http.MethodGet)

	// Add one nextActionType.
	router.HandleFunc("/tenant/{tenantID}/next-action-type/credential/{credentialID}",
		controller.AddNextActionType).Methods(http.MethodPost)

	// Add multiple nextActionTypes.
	router.HandleFunc("/tenant/{tenantID}/next-action-types/credential/{credentialID}",
		controller.AddNextActionTypes).Methods(http.MethodPost)

	// Update one nextActionType.
	router.HandleFunc("/tenant/{tenantID}/next-action-type/{nextActionTypeID}/credential/{credentialID}",
		controller.UpdateNextActionType).Methods(http.MethodPut)

	// Delete one nextActionType.
	router.HandleFunc("/tenant/{tenantID}/next-action-type/{nextActionTypeID}/credential/{credentialID}",
		controller.DeleteNextActionType).Methods(http.MethodDelete)

	log.NewLogger().Info("NextActionType Routes Registered")
}

// AddNextActionType will add nextActionType to the record.
func (controller *NextActionTypeController) AddNextActionType(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddNextActionType called==============================")
	nextActionType := &talent.NextActionType{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, nextActionType)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	nextActionType.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	nextActionType.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = nextActionType.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.NextActionTypeService.AddNextActionType(nextActionType)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Next action type added successfully")
}

// AddNextActionTypes will add multiple nextActionTypes to the table.
func (controller *NextActionTypeController) AddNextActionTypes(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddNextActionTypes called==============================")
	nextActionTypes := &[]talent.NextActionType{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, nextActionTypes)
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

	for _, nextActionType := range *nextActionTypes {
		fmt.Println("***nextActionType controller ->", nextActionType)
		err := nextActionType.Validate()
		fmt.Println("***err controller ->", err)
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = controller.NextActionTypeService.AddNextActionTypes(nextActionTypes, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Next action types added successfully")
}

// UpdateNextActionType will update the specified record in the table.
func (controller *NextActionTypeController) UpdateNextActionType(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateNextActionType called==============================")
	nextActionType := &talent.NextActionType{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, nextActionType)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	nextActionType.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	nextActionType.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	nextActionType.ID, err = util.ParseUUID(param["nextActionTypeID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = nextActionType.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.NextActionTypeService.UpdateNextActionType(nextActionType)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Next action type updated successfully")
}

// DeleteNextActionType will delete the specified record from table.
func (controller *NextActionTypeController) DeleteNextActionType(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteNextActionType called==============================")
	nextActionType := &talent.NextActionType{}
	param := mux.Vars(r)
	var err error

	nextActionType.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	nextActionType.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	nextActionType.ID, err = util.ParseUUID(param["nextActionTypeID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.NextActionTypeService.DeleteNextActionType(nextActionType)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Next action type deleted successfully")
}

// GetNextActionTypes will return all nextActionTypes from the table.
func (controller *NextActionTypeController) GetNextActionTypes(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetNextActionTypes called==============================")
	nextActionTypes := &[]talent.NextActionType{}
	param := mux.Vars(r)

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.NextActionTypeService.GetNextActionTypes(nextActionTypes, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, nextActionTypes)
}
