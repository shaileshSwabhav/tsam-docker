package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// MenuController provide methods to update, delete, add, get and get by role id for menu.
type MenuController struct {
	log         log.Logger
	auth        *security.Authentication
	MenuService *service.MenuService
}

// NewMenuController creates new instance of MenuController
func NewMenuController(menuService *service.MenuService, log log.Logger, auth *security.Authentication) *MenuController {
	return &MenuController{
		MenuService: menuService,
		log:         log,
		auth:        auth,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *MenuController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Add one menu.
	router.HandleFunc("/tenant/{tenantID}/menu",
		controller.AddMenu).Methods(http.MethodPost)

	// Get one menu.
	router.HandleFunc("/tenant/{tenantID}/menu/{menuID}",
		controller.GetMenu).Methods(http.MethodGet)

	// Get menus by role id.
	router.HandleFunc("/tenant/{tenantID}/menu/role/{roleID}",
		controller.GetMenusByRole).Methods(http.MethodGet)

	// Update one menu.
	router.HandleFunc("/tenant/{tenantID}/menu/{menuID}",
		controller.UpdateMenu).Methods(http.MethodPut)

	// Delete one menu.
	router.HandleFunc("/tenant/{tenantID}/menu/{menuID}",
		controller.DeleteMenu).Methods(http.MethodDelete)

	controller.log.Info("Menu Routes Registered")
}

// GetMenu returns one menu by specific menu id.
func (controller *MenuController) GetMenu(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetMenu call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	menu := general.MenuDTO{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	menu.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting menu id from param and parsing it to uuid.
	menu.ID, err = parser.GetUUID("menuID")
	if err != nil {
		controller.log.Error("unable to parse menu id")
		web.RespondError(w, errors.NewHTTPError("unable to parse menu id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.MenuService.GetMenu(&menu); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, menu)
}

// GetMenusByRole returns menus by role id.
func (controller *MenuController) GetMenusByRole(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetMenusByRole call**************************************")
	parser := web.NewParser(r)
	// Create bucket
	menus := &[]general.MenuDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting role id from param and parsing it to uuid.
	roleID, err := parser.GetUUID("roleID")
	if err != nil {
		controller.log.Error("unable to parse role id")
		web.RespondError(w, errors.NewHTTPError("unable to parse role id", http.StatusBadRequest))
		return
	}

	// Call get menu by role method.
	if err := controller.MenuService.GetMenusByRole(menus, roleID, tenantID); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, menus)
}

// AddMenu adds one menu.
func (controller *MenuController) AddMenu(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddMenu call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	menu := general.Menu{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &menu)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := menu.Validate(); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	menu.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	menu.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.MenuService.AddMenu(&menu)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Menu added successfully")
}

// UpdateMenu updates one menu by specific menu id.
func (controller *MenuController) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************UpdateMenu call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	menu := general.Menu{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &menu)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = menu.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting id from param and parsing it to uuid.
	menu.ID, err = parser.GetUUID("menuID")
	if err != nil {
		controller.log.Error("unable to parse menu id")
		web.RespondError(w, errors.NewHTTPError("unable to parse menu id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	menu.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	menu.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err := controller.MenuService.UpdateMenu(&menu); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Menu updated successfully")
}

// DeleteMenu deletes one menu by specific menu id.
func (controller *MenuController) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************DeleteMenu call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	menu := general.Menu{}

	var err error

	// Getting id from param and parsing it to uuid.
	menu.ID, err = parser.GetUUID("menuID")
	if err != nil {
		controller.log.Error("unable to parse menu id")
		web.RespondError(w, errors.NewHTTPError("unable to parse menu id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	menu.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	menu.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.MenuService.DeleteMenu(&menu); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Menu deleted successfully")
}
