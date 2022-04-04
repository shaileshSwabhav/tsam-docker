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

// TenantController Provide method to Update, Delete, Add, Get Method For Tenant.
type TenantController struct {
	TenantService *service.TenantService
}

// NewTenantController Create New Instance Of TenantController.
func NewTenantController(service *service.TenantService) *TenantController {
	return &TenantController{
		TenantService: service,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *TenantController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/tenant", controller.GetAllTenants).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}", controller.GetTenant).Methods(http.MethodGet)
	router.HandleFunc("/tenant", controller.AddTenant).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}", controller.UpdateTenant).Methods(http.MethodPut)
	router.HandleFunc("/tenant/{tenantID}", controller.DeleteTenant).Methods(http.MethodDelete)
	router.HandleFunc("/tenant/search", controller.SearchTenant).Methods(http.MethodPost)
	log.NewLogger().Info("Tenant Route Registered")
}

// GetTenant returns specified tenant
func (controller *TenantController) GetTenant(w http.ResponseWriter, r *http.Request) {
	tenant := &general.Tenant{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}
	err = controller.TenantService.GetTenant(tenant, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, tenant)
}

// GetAllTenants returns all the tenants
func (controller *TenantController) GetAllTenants(w http.ResponseWriter, r *http.Request) {
	tenant := &[]general.Tenant{}
	err := controller.TenantService.GetAllTenants(tenant)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, tenant)
}

// SearchTenant returns the searched tenant
func (controller *TenantController) SearchTenant(w http.ResponseWriter, r *http.Request) {
	tenant := &[]general.Tenant{}
	searchTenant := &general.Tenant{}
	err := web.UnmarshalJSON(r, searchTenant)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse requested data", http.StatusBadRequest))
		return
	}
	err = controller.TenantService.SearchTenants(tenant, searchTenant)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, tenant)
}

// AddTenant add new tenant to the database
func (controller *TenantController) AddTenant(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Add Tenant Call")
	var err error
	tenant := &general.Tenant{}
	err = web.UnmarshalJSON(r, tenant)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse requested data", http.StatusBadRequest))
		return
	}
	err = tenant.ValidateTenant()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	err = controller.TenantService.AddTenant(tenant)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Tenant added successfully")
}

// UpdateTenant updates the specified tenant
func (controller *TenantController) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Update Tenant Call")
	var err error
	tenant := &general.Tenant{}
	tenant.ID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}
	err = web.UnmarshalJSON(r, tenant)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse requested data", http.StatusBadRequest))
		return
	}
	// tenant.ID = tenantID
	err = tenant.ValidateTenant()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	err = controller.TenantService.UpdateTenant(tenant)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Tenant updated successfully")
}

// DeleteTenant deletes the specified tenant
func (controller *TenantController) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Delete Tenant Call")
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	err = controller.TenantService.DeleteTenant(tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Tenant deleted successfully")
}
