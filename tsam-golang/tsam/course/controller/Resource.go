package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/course/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/course"
	res "github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// ModuleResourceController Provide method to Update, Delete, Add, Get Method For resource.
type ModuleResourceController struct {
	ModuleResourceService *service.ModuleResourceService
	log                   log.Logger
	auth                  *security.Authentication
}

// NewModuleResourceController Create New Instance Of CourseSessionResourceController.
func NewModuleResourceController(resourceService *service.ModuleResourceService,
	log log.Logger, auth *security.Authentication) *ModuleResourceController {
	return &ModuleResourceController{
		ModuleResourceService: resourceService,
		log:                   log,
		auth:                  auth,
	}
}

// RegisterRoutes Register All Endpoint To Router
func (controller *ModuleResourceController) RegisterRoutes(router *mux.Router) {

	// add
	router.HandleFunc("/tenant/{tenantID}/modules/{moduleID}/resources",
		controller.AddResource).Methods(http.MethodPost)

	// delete
	router.HandleFunc("/tenant/{tenantID}/modules/{moduleID}/resources/{resourceID}",
		controller.DeleteResources).Methods(http.MethodDelete)

	// router.HandleFunc("/tenant/{tenantID}/modules/module/{moduleID}/resources",
	// 	controller.AddResources).Methods(http.MethodPost)

	// update
	// router.HandleFunc("/tenant/{tenantID}/session/{sessionID}/resources",
	// 	controller.UpdateResource).Methods(http.MethodPut)

	// get
	router.HandleFunc("/tenant/{tenantID}/modules/{moduleID}/resources",
		controller.GetModuleResources).Methods(http.MethodGet)

	router.Use(controller.auth.Middleware)

	// router.HandleFunc("/tenant/{tenantID}/resource", controller.GetAllResources).Methods(http.MethodGet)

	log.NewLogger().Info("Module Resource Route Registered")

}

// AddResource will add the resource for the specified session to DB
func (controller *ModuleResourceController) AddResource(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddModuleResource call==============================")
	resource := course.ModuleResource{}

	err := web.UnmarshalJSON(r, &resource)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resource.ModuleID, err = parser.GetUUID("moduleID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = resource.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ModuleResourceService.AddResource(&resource, tenantID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Module Resource successfully added")
}

// DeleteResources will update the resource for the specified resource to DB
func (controller *ModuleResourceController) DeleteResources(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteModulesResources call==============================")
	resource := course.ModuleResource{}
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resource.ModuleID, err = parser.GetUUID("moduleID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resource.ResourceID, err = parser.GetUUID("resourceID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ModuleResourceService.DeleteResources(&resource, tenantID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Module Resource successfully deleted")
}

// GetModuleResources returns all the resource for the specified session
func (controller *ModuleResourceController) GetModuleResources(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetModuleResources call==============================")
	var resources []res.Resource
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	moduleID, err := parser.GetUUID("moduleID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ModuleResourceService.GetModuleResources(&resources, tenantID, moduleID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, resources)
}

// // GetAllResource will return all the resources.
// func (controller *ModuleResourceController) GetAllResources(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetAllCourseSessionResources call==============================")
// 	param := mux.Vars(r)
// 	resource := []res.Resource{}

// 	tenantID, err := util.ParseUUID(param["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = controller.ModuleResourceService.GetAllResources(&resource, tenantID)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, resource)
// }

// // UpdateResource will update the resource for the specified resource to DB
// func (controller *CourseSessionResourceController) UpdateResource(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================UpdateResource call==============================")
// 	resources := []res.Resource{}

// 	parser := web.NewParser(r)

// 	err := web.UnmarshalJSON(r, &resources)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	// param := mux.Vars(r)
// 	tenantID, err := parser.GetTenantID()
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	sessionID, err := parser.GetUUID("sessionID")
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	for _, resource := range resources {
// 		err = resource.ValidateResource()
// 		if err != nil {
// 			log.NewLogger().Error(err.Error())
// 			web.RespondError(w, err)
// 			return
// 		}
// 	}

// 	err = controller.CourseSessionResourceService.UpdateResource(&resources, tenantID, sessionID, credentialID)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, "Resource updated")
// }
