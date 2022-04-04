package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/course/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// TopicQuestionController provides methods to do Update, Delete, Add, Get operations on topic_programming_questions.
type TopicQuestionController struct {
	log     log.Logger
	service *service.TopicQuestionService
	auth    *security.Authentication
}

// NewCourseTopicQuestionController creates new instance of CourseProgrammingQuestionController.
func NewCourseTopicQuestionController(service *service.TopicQuestionService,
	log log.Logger, auth *security.Authentication) *TopicQuestionController {
	return &TopicQuestionController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *TopicQuestionController) RegisterRoutes(router *mux.Router) {

	// add
	router.HandleFunc("/tenant/{tenantID}/topic/{topicID}/programming-question",
		controller.AddTopicProgrammingQuestion).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/topic/{topicID}/programming-question/{topicQuestionID}",
		controller.UpdateTopicProgrammingQuestion).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/topic/{topicID}/programming-question/{topicQuestionID}",
		controller.DeleteTopicProgrammingQuestion).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/topic/{topicID}/programming-question-list",
		controller.GetTopicProgrammingQuestionList).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/topic/{topicID}/programming-question",
		controller.GetTopicProgrammingQuestions).Methods(http.MethodGet)

	router.Use(controller.auth.Middleware)

	controller.log.Info("Topic Programming Question Routes Registered")
}

// AddTopicProgrammingQuestion will add new question to course.
func (controller *TopicQuestionController) AddTopicProgrammingQuestion(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add Topic Programming Question Called==============================")
	topicQuestion := course.TopicProgrammingQuestion{}

	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &topicQuestion)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	topicQuestion.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// // Parse and set topic ID.
	topicQuestion.TopicID, err = parser.GetUUID("topicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	topicQuestion.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = topicQuestion.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddTopicProgrammingQuestion(&topicQuestion)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Question successfully added for topic")
}

// UpdateTopicProgrammingQuestion will update question to course.
func (controller *TopicQuestionController) UpdateTopicProgrammingQuestion(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Update Topic Programming Question Called==============================")
	topicQuestion := course.TopicProgrammingQuestion{}

	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &topicQuestion)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	topicQuestion.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set topic ID.
	topicQuestion.TopicID, err = parser.GetUUID("topicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic id", http.StatusBadRequest))
		return
	}

	// Parse and set topic Question ID.
	topicQuestion.ID, err = parser.GetUUID("topicQuestionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic question id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field.
	topicQuestion.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = topicQuestion.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.TopicCourseProgrammingQuestion(&topicQuestion)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Question successfully update for topic")
}

// DeleteTopicProgrammingQuestion will delete assignment to topic.
func (controller *TopicQuestionController) DeleteTopicProgrammingQuestion(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Delete Topic Programming Question Called==============================")
	topicQuestion := course.TopicProgrammingQuestion{}

	parser := web.NewParser(r)
	var err error

	// Parse and set tenant ID.
	topicQuestion.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// // Parse and set topic ID.
	topicQuestion.TopicID, err = parser.GetUUID("topicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic id", http.StatusBadRequest))
		return
	}

	// Parse and set topic question ID.
	topicQuestion.ID, err = parser.GetUUID("topicQuestionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic assignment id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in DeletedBy field.
	topicQuestion.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.DeleteCourseProgrammingQuestion(&topicQuestion)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Question successfully deleted for topic")
}

// GetTopicProgrammingQuestionList will return all the topic programming questions.
func (controller *TopicQuestionController) GetTopicProgrammingQuestionList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Topic Programming Question List Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set course ID.
	topicID, err := parser.GetUUID("topicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic id", http.StatusBadRequest))
		return
	}

	topicQuestions := &[]course.TopicProgrammingQuestionDTO{}
	err = controller.service.GetTopicProgrammingQuestionList(tenantID, topicID, topicQuestions, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, topicQuestions)
}

// GetTopicProgrammingQuestions will fetch the programming questions for the specified topic with limit and offset.
func (controller *TopicQuestionController) GetTopicProgrammingQuestions(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Topic Programming Question Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set topic ID.
	topicID, err := parser.GetUUID("topicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic id", http.StatusBadRequest))
		return
	}

	topicQuestions := &[]course.TopicProgrammingQuestionDTO{}
	var totalCount int

	err = controller.service.GetTopicProgrammingQuestions(tenantID, topicID, topicQuestions, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, topicQuestions)
}
