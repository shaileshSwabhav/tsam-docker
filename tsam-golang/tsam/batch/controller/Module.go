package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

type ModuleController struct {
	log     log.Logger
	service *service.ModuleService
	auth    *security.Authentication
}

func NewModuleController(service *service.ModuleService, log log.Logger, auth *security.Authentication) *ModuleController {
	return &ModuleController{
		log:     log,
		service: service,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *ModuleController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/module",
		controller.addModule).Methods(http.MethodPost)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/modules",
		controller.addModules).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/modules/{batchModuleID}",
		controller.updateModule).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/modules/{batchModuleID}",
		controller.deleteModule).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/modules",
		controller.getBatchModules).Methods(http.MethodGet)

	// Get all batch modules with all fields.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/modules-all-fields",
		controller.GetBatchModulesWithAllFields).Methods(http.MethodGet)

	controller.log.Info("Batch Module Routes Registered")
}

// addModule will add module for specified batch.
func (controller *ModuleController) addModule(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddModule Called==============================")
	parser := web.NewParser(r)
	module := batch.Module{}

	err := web.UnmarshalJSON(r, &module)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	module.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	module.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	module.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = module.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddModule(&module)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusCreated, nil)
}

// addModules will add multiple modules to specified batch.
func (controller *ModuleController) addModules(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddModules Called==============================")
	parser := web.NewParser(r)
	modules := []batch.Module{}

	err := web.UnmarshalJSON(r, &modules)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, module := range modules {
		err = module.Validate()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = controller.service.AddModules(&modules, tenantID, credentialID, batchID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusCreated, nil)
}

// updateModule will update specified batch_module.
func (controller *ModuleController) updateModule(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateModule Called==============================")
	parser := web.NewParser(r)
	module := batch.Module{}

	err := web.UnmarshalJSON(r, &module)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	module.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	module.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	module.ID, err = parser.GetUUID("batchModuleID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	module.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = module.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdateModule(&module)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusCreated, nil)
}

// deleteModule will delete specified batch_module.
func (controller *ModuleController) deleteModule(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteModule Called==============================")
	parser := web.NewParser(r)
	module := batch.Module{}
	var err error

	module.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	module.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	module.ID, err = parser.GetUUID("batchModuleID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	module.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.DeleteModule(&module)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusCreated, nil)
}

// getBatchModules will get all the modules for specified batch.
func (controller *ModuleController) getBatchModules(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetBatchModules Called==============================")
	parser := web.NewParser(r)
	modules := []batch.ModuleDTO{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	var totalCount int

	err = controller.service.GetBatchModules(&modules, tenantID, batchID, &totalCount, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, modules)
}

// GetBatchModulesWithAllFields will get all the modules with all the preloads.
func (controller *ModuleController) GetBatchModulesWithAllFields(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetBatchModulesWithAllFields Called==============================")
	parser := web.NewParser(r)
	modules := []batch.ModuleDTO{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	var totalCount int

	err = controller.service.GetBatchModulesWithAllFields(&modules, tenantID, batchID, &totalCount, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, modules)
}
