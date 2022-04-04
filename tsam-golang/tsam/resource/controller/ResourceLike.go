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

// ResourceLikeController Provide method to Update, Delete, Add, Get Method For resource_likes.
type ResourceLikeController struct {
	ResourceLikeService *service.ResourceLikeService
}

// NewResourceLikeController Create New Instance Of ResourceLikeController.
func NewResourceLikeController(resourceService *service.ResourceLikeService) *ResourceLikeController {
	return &ResourceLikeController{
		ResourceLikeService: resourceService,
	}
}

// RegisterRoutes Register All Endpoint To Router
func (con *ResourceLikeController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/resource/like/credential/{credentialID}",
		con.AddLike).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/resource/like/{likeID}/credential/{credentialID}",
		con.UpdateLike).Methods(http.MethodPut)

	// get
	router.HandleFunc("/tenant/{tenantID}/resource/{resourceID}/like/credential/{credentialID}",
		con.GetResourceLike).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/resource/{resourceID}/like",
		con.GetAllResourceLikes).Methods(http.MethodGet)

	log.NewLogger().Info("Resource Like Route Registered")

}

// AddLike will add/update entry in the resource-like table
func (con *ResourceLikeController) AddLike(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddLike call==============================")

	resourceLike := res.Like{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &resourceLike)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resourceLike.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resourceLike.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = resourceLike.ValidateResourceLike()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.ResourceLikeService.AddLike(&resourceLike)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Like successfully added")
}

// UpdateLike will add/update entry in the resource-like table
func (con *ResourceLikeController) UpdateLike(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateLike call==============================")

	resourceLike := res.Like{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &resourceLike)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resourceLike.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resourceLike.ID, err = util.ParseUUID(param["likeID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	resourceLike.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = resourceLike.ValidateResourceLike()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.ResourceLikeService.UpdateLike(&resourceLike)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Like successfully updated")
}

// GetResourceLike will return count of like of specified resource
func (con *ResourceLikeController) GetResourceLike(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetResourceLike call==============================")

	resourceLike := res.LikeDTO{}
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

	credentialID, err := util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.ResourceLikeService.GetResourceLike(tenantID, resourceID, credentialID, &resourceLike)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, resourceLike)
}

// GetAllResourceLikes will return all the like details for specified resource.
func (con *ResourceLikeController) GetAllResourceLikes(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllResourceLikes call==============================")

	resourceLike := []res.LikeDTO{}
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

	err = con.ResourceLikeService.GetAllResourceLikes(tenantID, resourceID, &resourceLike)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, resourceLike)
}
