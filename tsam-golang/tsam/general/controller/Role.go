package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// RoleController provides methods to update, delete, add, get, get all and get all by role id for role
type RoleController struct {
	RoleService *service.RoleService
}

// NewRoleController creates new instance of RoleController.
func NewRoleController(ser *service.RoleService) *RoleController {
	return &RoleController{
		RoleService: ser,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *RoleController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	//get one role
	router.HandleFunc("/tenant/{tenantID}/role/{roleID}", controller.GetRole).Methods(http.MethodGet)

	//get all roles
	router.HandleFunc("/tenant/{tenantID}/role", controller.GetRoles).Methods(http.MethodGet)

	//add one role
	router.HandleFunc("/tenant/{tenantID}/role/credential/{credentialID}", controller.AddRole).Methods(http.MethodPost)

	//update one role
	router.HandleFunc("/tenant/{tenantID}/role/{roleID}/credential/{credentialID}", controller.UpdateRole).Methods(http.MethodPut)

	//delete one role
	router.HandleFunc("/tenant/{tenantID}/role/{roleID}/credential/{credentialID}", controller.DeleteRole).Methods(http.MethodDelete)

	log.NewLogger().Info("Role Route Registered")
}

// GetRole returns one role by specific role id
func (controller *RoleController) GetRole(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Get role API call")

	//create bucket
	role := general.Role{}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//getting role id from param and parsing it to uuid
	roleID, err := util.ParseUUID(mux.Vars(r)["roleID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse role id", http.StatusBadRequest))
		return
	}

	//call get service method
	err = controller.RoleService.GetRole(&role, roleID, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, role)
}

// GetRoles returns all roles
func (controller *RoleController) GetRoles(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Get all roles API call")

	//create bucket
	roles := []general.Role{}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//call get roles method
	err = controller.RoleService.GetRoles(&roles, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, roles)
}

// UpdateRole updates one role by specific role id
func (controller *RoleController) UpdateRole(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Update role API call")

	//create bucket
	role := general.Role{}

	//unmarshal json
	if err := web.UnmarshalJSON(r, &role); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	//validate compulsary fields
	if err := role.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//getting role id from param and parsing it to uuid
	roleID, err := util.ParseUUID(mux.Vars(r)["roleID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse role id", http.StatusBadRequest))
		return
	}

	//getting credential id from param and parsing it to uuid
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	//call update service method
	err = controller.RoleService.UpdateRole(&role, tenantID, roleID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Role updated successfully")
}

// AddRole adds one role
func (controller *RoleController) AddRole(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Add role API call")

	//create bucket
	role := general.Role{}

	//unmarshal json
	if err := web.UnmarshalJSON(r, &role); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	//validate compulsary fields
	if err := role.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//getting credential id from param and parsing it to uuid
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	//call add service method
	err = controller.RoleService.AddRole(&role, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Role added successfully")
}

// DeleteRole deletes one role by specific role id
func (controller *RoleController) DeleteRole(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Delete role API call")

	//create bucket
	role := general.Role{}

	//getting id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//getting id from param and parsing it to uuid
	roleID, err := util.ParseUUID(mux.Vars(r)["roleID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse role id", http.StatusBadRequest))
		return
	}

	//getting credential id from param and parsing it to uuid
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	//call delete service method
	err = controller.RoleService.DeleteRole(&role, tenantID, roleID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Role deleted successfully")
}
