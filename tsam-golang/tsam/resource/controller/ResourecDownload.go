package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	res "github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/resource/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ResourceDownloadController Provide method to Update, Delete, Add, Get Method For resource_downloads.
type ResourceDownloadController struct {
	ResourceDownloadService *service.ResourceDownloadService
}

// NewResourceDownloadController Create New Instance Of ResourceDownloadController.
func NewResourceDownloadController(resourceService *service.ResourceDownloadService) *ResourceDownloadController {
	return &ResourceDownloadController{
		ResourceDownloadService: resourceService,
	}
}

// RegisterRoutes Register All Endpoint To Router
func (con *ResourceDownloadController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/resource/download/credential/{credentialID}",
		con.AddResourceDownload).Methods(http.MethodPost)

	// get
	router.HandleFunc("/tenant/{tenantID}/resource/{resourceID}/download", con.GetResourceDownload).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/resource/{resourceID}/download/count", con.GetResourceCount).Methods(http.MethodGet)

	log.NewLogger().Info("Resource Download Route Registered")

}

// AddResourceDownload will add download details.
func (con *ResourceDownloadController) AddResourceDownload(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddResourceDownload call==============================")

	resourceDownload := res.Download{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &resourceDownload)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resourceDownload.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resourceDownload.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = resourceDownload.ValidateResourceDownload()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.ResourceDownloadService.AddResourceDownload(&resourceDownload)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Resource download successfully added")
}

// GetResourceDownload will return all the downloaded resources and its credential
func (con *ResourceDownloadController) GetResourceDownload(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllResourceDownload call==============================")

	resourceDownloads := []res.DownloadDTO{}
	param := mux.Vars(r)

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resourceID, err := util.ParseUUID(param["resourceID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.ResourceDownloadService.GetResourceDownload(tenantID, resourceID, &resourceDownloads)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, resourceDownloads)
}

// GetResourceCount will return count of download of specified resource
func (con *ResourceDownloadController) GetResourceCount(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetResourceCount call==============================")

	resourceDownload := res.DownloadCount{}
	param := mux.Vars(r)

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resourceID, err := util.ParseUUID(param["resourceID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.ResourceDownloadService.GetResourceCount(tenantID, resourceID, &resourceDownload)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, resourceDownload)
}
