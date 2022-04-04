package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// FeedbackOptionController provides methods to do CRUD operations.
type FeedbackOptionController struct {
	FeedbackOptionService *service.FeedbackOptionService
}

// NewFeedbackOptionController creates new instance of feedbackOption type controller.
func NewFeedbackOptionController(feedbackOptionService *service.FeedbackOptionService) *FeedbackOptionController {
	return &FeedbackOptionController{
		FeedbackOptionService: feedbackOptionService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *FeedbackOptionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/feedback/question/{questionID}/option/credential/{credentialID}",
		controller.AddFeedbackOption).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/feedback/question/{questionID}/options/credential/{credentialID}",
		controller.AddFeedbackOptions).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/feedback/question/{questionID}/option/{feedbackOptionID}/credential/{credentialID}",
		controller.UpdateFeedbackOption).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/feedback/question/{questionID}/option/{feedbackOptionID}/credential/{credentialID}",
		controller.DeleteFeedbackOption).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/feedback/question/{questionID}/option/{feedbackOptionID}",
		controller.GetFeedbackOption).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/feedback/question/{questionID}/option",
		controller.GetFeedbackOptionForQuestion).Methods(http.MethodGet)

	log.NewLogger().Info("Feedback Option Routes Registered")
}

// AddFeedbackOption will add feedback options to the table
func (controller *FeedbackOptionController) AddFeedbackOption(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddFeedbackOption called==============================")
	param := mux.Vars(r)
	feedbackOption := &general.FeedbackOption{}

	err := web.UnmarshalJSON(r, feedbackOption)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.QuestionID, err = util.ParseUUID(param["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = feedbackOption.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackOptionService.AddFeedbackOption(feedbackOption)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback option added successfully")
}

// AddFeedbackOptions will add multiple feedback options to the table
func (controller *FeedbackOptionController) AddFeedbackOptions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddFeedbackOptions called==============================")
	param := mux.Vars(r)
	feedbackOptions := &[]general.FeedbackOption{}

	err := web.UnmarshalJSON(r, feedbackOptions)
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

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	questionID, err := util.ParseUUID(param["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, feedbackOption := range *feedbackOptions {
		err = feedbackOption.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = controller.FeedbackOptionService.AddFeedbackOptions(feedbackOptions, tenantID, credentialID, questionID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback options added successfully")
}

// UpdateFeedbackOption will update specified feedback option in the table
func (controller *FeedbackOptionController) UpdateFeedbackOption(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateFeedbackOption called==============================")
	param := mux.Vars(r)
	feedbackOption := &general.FeedbackOption{}

	err := web.UnmarshalJSON(r, feedbackOption)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.QuestionID, err = util.ParseUUID(param["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.ID, err = util.ParseUUID(param["feedbackOptionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = feedbackOption.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackOptionService.UpdateFeedbackOption(feedbackOption)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback option updated successfully")
}

// DeleteFeedbackOption will delete specified feedback option from the table
func (controller *FeedbackOptionController) DeleteFeedbackOption(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteFeedbackOption called==============================")
	param := mux.Vars(r)
	var err error
	feedbackOption := &general.FeedbackOption{}

	feedbackOption.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.QuestionID, err = util.ParseUUID(param["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.ID, err = util.ParseUUID(param["feedbackOptionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackOptionService.DeleteFeedbackOption(feedbackOption)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback option deleted successfully")
}

// GetFeedbackOption will return specified feedback option
func (controller *FeedbackOptionController) GetFeedbackOption(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetFeedbackOption called==============================")
	param := mux.Vars(r)
	var err error
	feedbackOption := &general.FeedbackOption{}

	feedbackOption.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackOption.ID, err = util.ParseUUID(param["feedbackOptionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackOptionService.GetFeedbackOption(feedbackOption)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, feedbackOption)
}

// GetAllFeedbackOption will return specified feedback option
func (controller *FeedbackOptionController) GetAllFeedbackOption(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllFeedbackOption called==============================")
	param := mux.Vars(r)
	feedbackOptions := &[]general.FeedbackOption{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackOptionService.GetAllFeedbackOptions(feedbackOptions, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, feedbackOptions)
}

// GetFeedbackOptionForQuestion will return specified feedback option
func (controller *FeedbackOptionController) GetFeedbackOptionForQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetFeedbackOptionForQuestion called==============================")
	param := mux.Vars(r)
	feedbackOptions := &[]general.FeedbackOption{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	questionID, err := util.ParseUUID(param["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackOptionService.GetFeedbackOptionsForQuestion(feedbackOptions, tenantID, questionID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, feedbackOptions)
}
