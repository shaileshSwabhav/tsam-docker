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

// FeedbackQuestionController provides methods to do CRUD operations.
type FeedbackQuestionController struct {
	FeedbackQuestionService *service.FeedbackQuestionService
}

// NewFeedbackQuestionController creates new instance of feedbackQuestion type controller.
func NewFeedbackQuestionController(feedbackQuestionService *service.FeedbackQuestionService) *FeedbackQuestionController {
	return &FeedbackQuestionController{
		FeedbackQuestionService: feedbackQuestionService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *FeedbackQuestionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// take confirmation for the endpoint from sir

	// add
	router.HandleFunc("/tenant/{tenantID}/feedback-question/credential/{credentialID}",
		controller.AddFeedbackQuestion).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/feedback-questions/credential/{credentialID}",
		controller.AddFeedbackQuestions).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/feedback-question/{feedbackQuestionID}/credential/{credentialID}",
		controller.UpdateFeedbackQuestion).Methods(http.MethodPut)
	router.HandleFunc("/tenant/{tenantID}/feedback-question/{feedbackQuestionID}/status/credential/{credentialID}",
		controller.UpdateFeedbackQuestionStatus).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/feedback-question/{feedbackQuestionID}/credential/{credentialID}",
		controller.DeleteFeedbackQuestion).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/feedback-question/{questionID}", controller.GetFeedbackQuestion).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/feedback-question/limit/{limit}/offset/{offset}", controller.GetFeedbackQuestions).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/feedback-question", controller.GetAllFeedbackQuestion).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/feedback-question/type/{type}", controller.GetFeedbackQuestionsByType).Methods(http.MethodGet)

	log.NewLogger().Info("Feedback Question Routes Registered")
}

// AddFeedbackQuestion will add feedback question to the table
func (controller *FeedbackQuestionController) AddFeedbackQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddFeedbackQuestion called==============================")
	param := mux.Vars(r)
	feedbackQuestion := &general.FeedbackQuestion{}

	err := web.UnmarshalJSON(r, feedbackQuestion)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackQuestion.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	feedbackQuestion.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = feedbackQuestion.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackQuestionService.AddFeedbackQuestion(feedbackQuestion)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback question added successfully")
}

// AddFeedbackQuestions will add multiple feedback questions to the table
func (controller *FeedbackQuestionController) AddFeedbackQuestions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddFeedbackQuestions called==============================")
	param := mux.Vars(r)
	feedbackQuestions := &[]general.FeedbackQuestion{}

	err := web.UnmarshalJSON(r, feedbackQuestions)
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

	for _, feedbackQuestion := range *feedbackQuestions {
		err = feedbackQuestion.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = controller.FeedbackQuestionService.AddFeedbackQuestions(feedbackQuestions, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback questions added successfully")
}

// UpdateFeedbackQuestion will update specified feedback question
func (controller *FeedbackQuestionController) UpdateFeedbackQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateFeedbackQuestion called==============================")
	param := mux.Vars(r)
	feedbackQuestion := &general.FeedbackQuestion{}

	err := web.UnmarshalJSON(r, feedbackQuestion)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackQuestion.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	feedbackQuestion.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackQuestion.ID, err = util.ParseUUID(param["feedbackQuestionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = feedbackQuestion.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackQuestionService.UpdateFeedbackQuestion(feedbackQuestion)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback question updated successfully")
}

// UpdateFeedbackQuestionStatus will update the status of the specified feedback question
func (controller *FeedbackQuestionController) UpdateFeedbackQuestionStatus(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateFeedbackQuestionStatus called==============================")

	param := mux.Vars(r)
	feedbackQuestion := &general.FeedbackQuestion{}

	err := web.UnmarshalJSON(r, feedbackQuestion)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackQuestion.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	feedbackQuestion.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackQuestion.ID, err = util.ParseUUID(param["feedbackQuestionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = feedbackQuestion.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackQuestionService.UpdateFeedbackQuestionStatus(feedbackQuestion)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback questions' status updated successfully")
}

// DeleteFeedbackQuestion will delete specified feedback question
func (controller *FeedbackQuestionController) DeleteFeedbackQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteFeedbackQuestion called==============================")
	param := mux.Vars(r)
	var err error
	feedbackQuestion := &general.FeedbackQuestion{}

	feedbackQuestion.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	feedbackQuestion.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackQuestion.ID, err = util.ParseUUID(param["feedbackQuestionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackQuestionService.DeleteFeedbackQuestion(feedbackQuestion)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback question deleted successfully")
}

// GetFeedbackQuestion will return specified feedback question
func (controller *FeedbackQuestionController) GetFeedbackQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetFeedbackQuestion called==============================")
	param := mux.Vars(r)
	var err error
	feedbackQuestion := &general.FeedbackQuestion{}

	feedbackQuestion.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedbackQuestion.ID, err = util.ParseUUID(param["questionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackQuestionService.GetFeedbackQuestion(feedbackQuestion)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, feedbackQuestion)
}

// GetFeedbackQuestions will return all the feedback questions
func (controller *FeedbackQuestionController) GetFeedbackQuestions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetFeedbackQuestions called==============================")
	param := mux.Vars(r)
	feedbackQuestions := &[]general.FeedbackQuestionDTO{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Parse form
	r.ParseForm()

	err = controller.FeedbackQuestionService.GetFeedbackQuestions(feedbackQuestions, r.Form, tenantID, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, feedbackQuestions)
}

// GetAllFeedbackQuestion will return all the feedback questions
func (controller *FeedbackQuestionController) GetAllFeedbackQuestion(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllFeedbackQuestion called==============================")
	param := mux.Vars(r)
	feedbackQuestions := &[]general.FeedbackQuestionDTO{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeedbackQuestionService.GetAllFeedbackQuestion(feedbackQuestions, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, feedbackQuestions)
}

// GetFeedbackQuestionsByType will return all the feedback questions
func (controller *FeedbackQuestionController) GetFeedbackQuestionsByType(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetFeedbackQuestionsByType called==============================")
	param := mux.Vars(r)
	feedbackQuestions := &[]general.FeedbackQuestionDTO{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	typeName := param["type"]

	err = controller.FeedbackQuestionService.GetFeedbackQuestionsByType(feedbackQuestions, tenantID, typeName)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, feedbackQuestions)
}
