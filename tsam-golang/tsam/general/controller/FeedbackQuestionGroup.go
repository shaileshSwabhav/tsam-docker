package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// FeedbackQuestionGroupController provides method to update, delete, add, get method for FeedbackQuestionGroup.
type FeedbackQuestionGroupController struct {
	FeedbackQuestionGroupService *service.FeedbackQuestionGroupService
}

// NewFeedbackQuestionGroupController creates new instance of FeedbackQuestionGroupController.
func NewFeedbackQuestionGroupController(feedbackQuestionGroupService *service.FeedbackQuestionGroupService) *FeedbackQuestionGroupController {
	return &FeedbackQuestionGroupController{
		FeedbackQuestionGroupService: feedbackQuestionGroupService,
	}
}

/// RegisterRoutes registers all endpoints to router.
func (con *FeedbackQuestionGroupController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Groupwise keyword names
	router.HandleFunc("/tenant/{tenantID}/feedback-question-group/group-wise-keyword",
		con.GetGroupwiseKeywordName).Methods(http.MethodGet)

	// Get all feedbackQuestionGroups with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/feedback-question-group/limit/{limit}/offset/{offset}",
		con.GetFeedbackQuestionGroups).Methods(http.MethodGet)

	// Get all feedbackQuestionGroup list by type.
	router.HandleFunc("/tenant/{tenantID}/feedback-question-group/type/{type}",
		con.GetFeedbackQuestionGroupListByType).Methods(http.MethodGet)

	// Get feedbackQuestionGroup by question type
	router.HandleFunc("/tenant/{tenantID}/feedback-question-group/feedback-question/type/{type}",
		con.GetFeedbackQuestionGroupByType).Methods(http.MethodGet)

	// Get all feedbackQuestionGroup list.
	router.HandleFunc("/tenant/{tenantID}/feedback-question-group",
		con.GetFeedbackQuestionGroupList).Methods(http.MethodGet)

	// Get one feedbackQuestionGroup.
	router.HandleFunc("/tenant/{tenantID}/feedback-question-group/{feedbackQuestionGroupID}",
		con.GetFeedbackQuestionGroup).Methods(http.MethodGet)

	// Add one feedbackQuestionGroup.
	router.HandleFunc("/tenant/{tenantID}/feedback-question-group/credential/{credentialID}",
		con.AddFeedbackQuestionGroup).Methods(http.MethodPost)

	// Add multiple feedbackQuestionGroups.
	router.HandleFunc("/tenant/{tenantID}/feedback-question-group/credential/{credentialID}",
		con.AddFeedbackQuestionGroups).Methods(http.MethodPost)

	// Update feedbackQuestionGroup.
	router.HandleFunc("/tenant/{tenantID}/feedback-question-group/{feedbackQuestionGroupID}/credential/{credentialID}",
		con.UpdateFeedbackQuestionGroup).Methods(http.MethodPut)

	// Delete feedbackQuestionGroup.
	router.HandleFunc("/tenant/{tenantID}/feedback-question-group/{feedbackQuestionGroupID}/credential/{credentialID}",
		con.DeleteFeedbackQuestionGroup).Methods(http.MethodDelete)

	log.NewLogger().Info("FeedbackQuestionGroup Route Registered")
}

// GetFeedbackQuestionGroups returns all feedbackQuestionGroups.
func (con *FeedbackQuestionGroupController) GetFeedbackQuestionGroups(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetFeedbackQuestionGroups call********************************")

	// Create bucket.
	feedbackQuestionGroups := []general.FeedbackQuestionGroupDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parsinf for query params.
	r.ParseForm()

	// Fpr pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all feedbackQuestionGroups service method.
	err = con.FeedbackQuestionGroupService.GetFeedbackQuestionGroups(&feedbackQuestionGroups, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, feedbackQuestionGroups)
}

// GetFeedbackQuestionGroupList returns feedbackQuestionGroup list.
func (con *FeedbackQuestionGroupController) GetFeedbackQuestionGroupList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetFeedbackQuestionGroupList call********************************")

	// Create bucket.
	feedbackQuestionGroups := []general.FeedbackQuestionGroupDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get feedbackQuestionGroup list method.
	err = con.FeedbackQuestionGroupService.GetFeedbackQuestionGroupList(&feedbackQuestionGroups, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, feedbackQuestionGroups)
}

// GetFeedbackQuestionGroupListByType returns feedbackQuestionGroup list.
func (con *FeedbackQuestionGroupController) GetFeedbackQuestionGroupListByType(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetFeedbackQuestionGroupListByType call********************************")

	// Create bucket.
	feedbackQuestionGroups := []general.FeedbackQuestionGroupDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	questionType := mux.Vars(r)["type"]

	// Call get feedbackQuestionGroup list method.
	err = con.FeedbackQuestionGroupService.GetFeedbackQuestionGroupListByType(&feedbackQuestionGroups,
		tenantID, questionType)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, feedbackQuestionGroups)
}

// GetFeedbackQuestionGroup return the specifed feedbackQuestionGroup
func (con *FeedbackQuestionGroupController) GetFeedbackQuestionGroup(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************Get FeedbackQuestionGroup Call********************************")
	feedbackQuestionGroup := general.FeedbackQuestionGroupDTO{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	feedbackQuestionGroupID, err := util.ParseUUID(mux.Vars(r)["feedbackQuestionGroupID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse feedbackQuestionGroup id", http.StatusBadRequest))
		return
	}
	err = con.FeedbackQuestionGroupService.GetFeedbackQuestionGroup(&feedbackQuestionGroup, tenantID, feedbackQuestionGroupID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, feedbackQuestionGroup)
}

// GetFeedbackQuestionGroupByType will return group-wise feedback question based on question-type
func (con *FeedbackQuestionGroupController) GetFeedbackQuestionGroupByType(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetFeedbackQuestionGroupByType call********************************")

	feedbackQuestionGroup := []general.FeedbackQuestionGroupDTO{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	questionType := mux.Vars(r)["type"]

	err = con.FeedbackQuestionGroupService.GetFeedbackQuestionGroupByType(&feedbackQuestionGroup, tenantID, questionType)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, feedbackQuestionGroup)
}

// AddFeedbackQuestionGroup adds new feedbackQuestionGroup.
func (con *FeedbackQuestionGroupController) AddFeedbackQuestionGroup(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddFeedbackQuestionGroup call********************************")

	// Create bucket.
	feedbackQuestionGroup := general.FeedbackQuestionGroup{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &feedbackQuestionGroup)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = feedbackQuestionGroup.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	feedbackQuestionGroup.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	feedbackQuestionGroup.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = con.FeedbackQuestionGroupService.AddFeedbackQuestionGroup(&feedbackQuestionGroup)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Feddback question group added successfully")
}

// AddFeedbackQuestionGroups adds multiple feedbackQuestionGroups.
func (con *FeedbackQuestionGroupController) AddFeedbackQuestionGroups(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddFeedbackQuestionGroups call**************************************")

	// Create bucket.
	feedbackQuestionGroups := []general.FeedbackQuestionGroup{}
	feedbackQuestionGroupIDs := []uuid.UUID{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &feedbackQuestionGroups)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate feedbackQuestionGroup fields.
	for _, feedbackQuestionGroup := range feedbackQuestionGroups {
		err = feedbackQuestionGroup.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add multiple service method.
	err = con.FeedbackQuestionGroupService.AddFeedbackQuestionGroups(&feedbackQuestionGroups, &feedbackQuestionGroupIDs, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Feedback question groups added successfully")
}

// UpdateFeedbackQuestionGroup updates the specified feedbackQuestionGroup.
func (con *FeedbackQuestionGroupController) UpdateFeedbackQuestionGroup(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateFeedbackQuestionGroup call**************************************")

	// Create bucket.
	feedbackQuestionGroup := general.FeedbackQuestionGroup{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &feedbackQuestionGroup)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate feedbackQuestionGroup fields.
	err = feedbackQuestionGroup.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	feedbackQuestionGroup.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting feedbackQuestionGroup id from param and parsing it to uuid.
	feedbackQuestionGroup.ID, err = util.ParseUUID(mux.Vars(r)["feedbackQuestionGroupID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse feedbackQuestionGroup id", http.StatusBadRequest))
		return
	}

	// Getting cresential id from param and parsing it to uuid.
	feedbackQuestionGroup.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = con.FeedbackQuestionGroupService.UpdateFeedbackQuestionGroup(&feedbackQuestionGroup)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Feedback question group updated successfully")
}

// DeleteFeedbackQuestionGroup deletes FeedbackQuestionGroup.
func (con *FeedbackQuestionGroupController) DeleteFeedbackQuestionGroup(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteFeedbackQuestionGroup call**************************************")

	// Create bucket.
	feedbackQuestionGroup := general.FeedbackQuestionGroup{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	feedbackQuestionGroup.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting feedbackQuestionGroup id from param and parsing it to uuid.
	feedbackQuestionGroup.ID, err = util.ParseUUID(mux.Vars(r)["feedbackQuestionGroupID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse feedbackQuestionGroup id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	feedbackQuestionGroup.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = con.FeedbackQuestionGroupService.DeleteFeedbackQuestionGroup(&feedbackQuestionGroup)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Feedback question group deleted successfully")
}

// GetGroupwiseKeywordName will get all the feedback question groups for faculty-session-feedback.
func (controller *FeedbackQuestionGroupController) GetGroupwiseKeywordName(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetGroupwiseKeywordName call==============================")

	keywordNames := new([]dashboard.GroupWiseKeywordName)

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackQuestionGroupService.GetGroupwiseKeywordName(tenantID, keywordNames)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, keywordNames)
}
