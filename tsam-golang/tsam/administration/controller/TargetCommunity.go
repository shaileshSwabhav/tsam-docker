package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/administration/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TargetCommunityController provides methods to do CRUD operations.
type TargetCommunityController struct {
	TargetCommunityService *service.TargetCommunityService
}

// NewTargetCommunityController creates new instance of target community controller.
func NewTargetCommunityController(generalService *service.TargetCommunityService) *TargetCommunityController {
	return &TargetCommunityController{
		TargetCommunityService: generalService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *TargetCommunityController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add one target community.
	router.HandleFunc("/tenant/{tenantID}/target-community/credential/{credentialID}",
		controller.AddTargetCommunity).Methods(http.MethodPost)

	// Update target community.
	router.HandleFunc("/tenant/{tenantID}/target-community/{targetCommunityID}/credential/{credentialID}",
		controller.UpdateTargetCommunity).Methods(http.MethodPut)

	// Update is target achieved for target community.
	router.HandleFunc("/tenant/{tenantID}/target-community-achieved/credential/{credentialID}",
		controller.UpdateIsAchievedTargetCommunity).Methods(http.MethodPut)

	// Delete target community.
	router.HandleFunc("/tenant/{tenantID}/target-community/{targetCommunityID}/credential/{credentialID}",
		controller.DeleteTargetCommunity).Methods(http.MethodDelete)

	// Get all target community by limit and offset.
	router.HandleFunc("/tenant/{tenantID}/target-community/limit/{limit}/offset/{offset}",
		controller.GetTargetCommunities).Methods(http.MethodGet)

	// Get target community list.
	router.HandleFunc("/tenant/{tenantID}/target-community",
		controller.GetTargetCommunityList).Methods(http.MethodGet)

	// Get one target community.
	router.HandleFunc("/tenant/{tenantID}/target-community/{targetCommunityID}",
		controller.GetTargetCommunity).Methods(http.MethodGet)

	log.NewLogger().Info("Target Community Route Registered")
}

// AddTargetCommunity will add the target_community record in the table.
func (controller *TargetCommunityController) AddTargetCommunity(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddTargetCommunity called==============================")
	param := mux.Vars(r)
	community := admin.TargetCommunity{}

	err := web.UnmarshalJSON(r, &community)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	community.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	community.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = community.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TargetCommunityService.AddTargetCommunity(&community)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Target community added successfully")
}

// UpdateTargetCommunity will update the specified tagert community record in the table.
func (controller *TargetCommunityController) UpdateTargetCommunity(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateTargetCommunity called==============================")
	param := mux.Vars(r)
	community := admin.TargetCommunity{}

	err := web.UnmarshalJSON(r, &community)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	community.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	community.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	community.ID, err = util.ParseUUID(param["targetCommunityID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = community.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TargetCommunityService.UpdateTargetCommunity(&community)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Target community updated successfully")
}

// UpdateIsAchievedTargetCommunity will update the is target achieved field of specified tagert community record in the table.
func (controller *TargetCommunityController) UpdateIsAchievedTargetCommunity(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateIsAchievedTargetCommunity called==============================")
	param := mux.Vars(r)
	community := admin.TargetCommunityUpdate{}

	err := web.UnmarshalJSON(r, &community)
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

	err = controller.TargetCommunityService.UpdateIsAchievedTargetCommunity(&community, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Target community achieved successfully")
}

// DeleteTargetCommunity will delete the specified target community record from the table.
func (controller *TargetCommunityController) DeleteTargetCommunity(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteTargetCommunity called==============================")
	param := mux.Vars(r)
	community := admin.TargetCommunity{}
	var err error

	community.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	community.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	community.ID, err = util.ParseUUID(param["targetCommunityID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TargetCommunityService.DeleteTargetCommunity(&community)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Target community deleted successfully")
}

// GetTargetCommunities will return all the records from target_community table.
func (controller *TargetCommunityController) GetTargetCommunities(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllTargetCommunity called==============================")
	param := mux.Vars(r)
	communities := []admin.TargetCommunityDTO{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TargetCommunityService.GetTargetCommunities(&communities, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, communities)
}

// GetTargetCommunityList will return all the records from target_community table.
func (controller *TargetCommunityController) GetTargetCommunityList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTargetCommunities called==============================")
	param := mux.Vars(r)
	communities := []admin.TargetCommunity{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TargetCommunityService.GetTargetCommunityList(&communities, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, communities)
}

// GetTargetCommunity will return specified record from target_community table.
func (controller *TargetCommunityController) GetTargetCommunity(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTargetCommunity called==============================")
	param := mux.Vars(r)
	community := admin.TargetCommunity{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	targetCommunityID, err := util.ParseUUID(param["targetCommunityID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TargetCommunityService.GetTargetCommunity(&community, tenantID, targetCommunityID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, community)
}
