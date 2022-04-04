package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/course/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// ModuleController provides methods to do CRUD operations.
type ModuleController struct {
	log     log.Logger
	auth    *security.Authentication
	Service *service.ModuleService
}

// NewModuleController creates new instance of module controller.
func NewModuleController(moduleService *service.ModuleService, log log.Logger, auth *security.Authentication) *ModuleController {
	return &ModuleController{
		Service: moduleService,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *ModuleController) RegisterRoutes(router *mux.Router) {

	// add
	router.HandleFunc("/tenant/{tenantID}/modules", controller.AddModule).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/modules/{moduleID}", controller.UpdateModule).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/modules/{moduleID}", controller.DeleteModule).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/modules", controller.GetAllModules).Methods(http.MethodGet)

	router.Use(controller.auth.Middleware)
	controller.log.Info("Module Routes Registered")
}

// AddModule will add new module to the modules table.
func (controller *ModuleController) AddModule(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddModule called==============================")

	parser := web.NewParser(r)

	module := course.Module{}

	err := web.UnmarshalJSON(r, &module)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	module.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	module.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	err = module.Validate()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.Service.AddModule(&module)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Module successfully added")
}

// UpdateModule will update specified module in the modules table.
func (controller *ModuleController) UpdateModule(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateModule called==============================")

	parser := web.NewParser(r)

	module := course.Module{}

	err := web.UnmarshalJSON(r, &module)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	module.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	module.ID, err = parser.GetUUID("moduleID")
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	module.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	err = module.Validate()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.Service.UpdateModule(&module)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Module successfully updated")
}

// DeleteModule will delete specified module in the modules table.
func (controller *ModuleController) DeleteModule(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteModule called==============================")

	parser := web.NewParser(r)

	module := course.Module{}
	var err error

	module.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	module.ID, err = parser.GetUUID("moduleID")
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	module.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.Service.DeleteModule(&module)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Module successfully deleted")
}

// GetAllModules will get all modules.
func (controller *ModuleController) GetAllModules(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllModules called==============================")
	parser := web.NewParser(r)

	modules := []course.ModuleDTO{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	var totalCount int

	err = controller.Service.GetAllModules(tenantID, &modules, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, modules)
}

// // GetAllModules will get all modules.
// func (controller *ModuleController) GetModules(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================GetModules called==============================")
// 	parser := web.NewParser(r)

// 	modules := []course.ModuleDTO{}

// 	tenantID, err := parser.GetTenantID()
// 	if err != nil {
// 		controller.log.Error(err)
// 		web.RespondError(w, err)
// 		return
// 	}

// 	var totalCount int

// 	err = controller.Service.GetModules(tenantID, &modules, parser, &totalCount)
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, modules)
// }
