package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	res "github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/resource/service"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ResourceController Provide method to Update, Delete, Add, Get Method For resource.
type ResourceController struct {
	ResourceService *service.ResourceService
	log             log.Logger
	auth            *security.Authentication
}

// NewResourceController Create New Instance Of ResourceController.
func NewResourceController(resourceService *service.ResourceService, log log.Logger, auth *security.Authentication) *ResourceController {
	return &ResourceController{
		ResourceService: resourceService,
		log:             log,
		auth:            auth,
	}
}

// RegisterRoutes Register All Endpoint To Router
func (controller *ResourceController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/resource", controller.AddResource).Methods(http.MethodPost)

	router.HandleFunc("/tenant/{tenantID}/resources", controller.AddResources).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/resource/{resourceID}", controller.UpdateResource).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/resource/{resourceID}", controller.DeleteResource).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/resource/limit/{limit}/offset/{offset}", controller.GetAllResources).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/resource/file-type/count", controller.GetResourceCount).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/resource-list", controller.GetResourcesList).Methods(http.MethodGet)

	// router.HandleFunc("/tenant/{tenantID}/resource/file-type/{fileType}/limit/{limit}/offset/{offset}",
	// 	con.GetResourcesByFileType).Methods(http.MethodGet)

	log.NewLogger().Info("Resource Route Registered")

}

// AddResource will add new resource to the table.
func (controller *ResourceController) AddResource(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddResource call==============================")
	resource := res.Resource{}
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &resource)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resource.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resource.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = resource.ValidateResource()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ResourceService.AddResource(&resource)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Resource successfully added")
}

// AddResources will add multiple resource to the table.
func (controller *ResourceController) AddResources(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddResources call==============================")
	resources := []res.Resource{}
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &resources)
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

	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, resource := range resources {
		err = resource.ValidateResource()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = controller.ResourceService.AddResources(&resources, tenantID, credentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Resources successfully added")
}

// UpdateResource will update the specified resource.
func (controller *ResourceController) UpdateResource(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateResource call==============================")
	resource := res.Resource{}
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &resource)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resource.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resource.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resource.ID, err = parser.GetUUID("resourceID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = resource.ValidateResource()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ResourceService.UpdateResource(&resource)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Resource successfully updated")
}

// DeleteResource will delete the specified resource.
func (controller *ResourceController) DeleteResource(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteResource call==============================")
	resource := res.Resource{}
	// param := mux.Vars(r)
	parser := web.NewParser(r)
	var err error

	resource.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resource.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resource.ID, err = parser.GetUUID("resourceID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ResourceService.DeleteResource(&resource)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Resource successfully deleted")
}

// GetAllResources will return all the resources.
func (con *ResourceController) GetAllResources(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllResources call==============================")
	resource := []res.ResourceDTO{}
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	r.ParseForm()

	err = con.ResourceService.GetAllResources(&resource, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, resource)
}

// GetResourcesByFileType will return all the resources for the specified file-type.
// func (con *ResourceController) GetResourcesByFileType(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetResourcesByFileType call==============================")
// 	resource := []res.Resource{}
// 	param := mux.Vars(r)

// 	tenantID, err := util.ParseUUID(param["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	fileType := param["fileType"]

// 	var totalCount int
// 	limit, offset := web.GetLimitAndOffset(r)

// 	r.ParseForm()

// 	err = con.ResourceService.GetResourcesByFileType(&resource, tenantID, r.Form, fileType, limit, offset, &totalCount)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, resource)
// }

// GetResourceCount will return all the resources.
func (con *ResourceController) GetResourceCount(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetResourceCount call==============================")
	resourceCount := []res.Count{}
	param := mux.Vars(r)

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// fileType := param["fileType"]

	r.ParseForm()

	err = con.ResourceService.GetResourceCount(&resourceCount, tenantID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, resourceCount)
}

// GetResourcesList will return all the resources for specified resource and file type.
func (con *ResourceController) GetResourcesList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetCourseSessionResourcesByType call==============================")
	param := mux.Vars(r)
	resource := []res.Resource{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	err = con.ResourceService.GetResourcesList(&resource, tenantID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, resource)
}
