package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/course/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// TopicController Provide method to Update, Delete, Add, Get Method For topic.
type TopicController struct {
	log          log.Logger
	auth         *security.Authentication
	TopicService *service.TopicService
}

// NewTopicController Create New Instance Of TopicController.
func NewTopicController(topicService *service.TopicService, log log.Logger, auth *security.Authentication) *TopicController {
	return &TopicController{
		TopicService: topicService,
		log:          log,
		auth:         auth,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *TopicController) RegisterRoutes(router *mux.Router) {

	// add
	router.HandleFunc("/tenant/{tenantID}/modules/{moduleID}/topic", controller.AddTopic).Methods(http.MethodPost)

	router.HandleFunc("/tenant/{tenantID}/modules/{moduleID}/topics", controller.AddTopics).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/modules/{moduleID}/topic/{topicID}", controller.UpdateTopic).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/modules/{moduleID}/topic/{topicID}", controller.DeleteTopic).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/modules/{moduleID}/topic", controller.GetAllTopics).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/topic-list", controller.GetTopicList).Methods(http.MethodGet)

	router.Use(controller.auth.Middleware)
	controller.log.Info("Topic Routes Registered")
}

// AddTopic will add topic for specified module.
func (controller *TopicController) AddTopic(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddTopic call==============================")

	topic := course.ModuleTopic{}

	err := web.UnmarshalJSON(r, &topic)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	parser := web.NewParser(r)

	topic.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	topic.ModuleID, err = parser.GetUUID("moduleID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse module id", http.StatusBadRequest))
		return
	}

	topic.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = topic.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TopicService.AddTopic(&topic)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Topic successfully added")
}

// AddTopics will add multiple topics for specified module.
func (controller *TopicController) AddTopics(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddTopics call==============================")

	topics := &[]course.ModuleTopic{}

	err := web.UnmarshalJSON(r, topics)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	moduleID, err := parser.GetUUID("moduleID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse module id", http.StatusBadRequest))
		return
	}

	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, topic := range *topics {
		err = topic.Validate()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = controller.TopicService.AddTopics(topics, tenantID, moduleID, credentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Topics successfully added")
}

// UpdateTopic will update the topics for specified module
func (controller *TopicController) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateTopic call==============================")

	topic := course.ModuleTopic{}
	// param := mux.Vars(r)
	err := web.UnmarshalJSON(r, &topic)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	parser := web.NewParser(r)

	topic.ID, err = parser.GetUUID("topicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse session id", http.StatusBadRequest))
		return
	}

	topic.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	topic.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	topic.ModuleID, err = parser.GetUUID("moduleID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse course id", http.StatusBadRequest))
		return
	}

	err = topic.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TopicService.UpdateTopic(&topic)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Topic successfully updated")
}

// DeleteTopic will delete specified topic from the DB.
func (controller *TopicController) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteTopic call==============================")

	topic := &course.ModuleTopic{}
	var err error

	parser := web.NewParser(r)

	topic.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	topic.ModuleID, err = parser.GetUUID("moduleID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse module id", http.StatusBadRequest))
		return
	}

	topic.ID, err = parser.GetUUID("topicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic id", http.StatusBadRequest))
		return
	}

	topic.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.TopicService.DeleteTopic(topic)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Topic Successfully Deleted")
}

// GetAllTopics returns all the sessions for the specified course
func (controller *TopicController) GetAllTopics(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllTopics call==============================")
	sessions := &[]course.ModuleTopicDTO{}

	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	moduleID, err := parser.GetUUID("moduleID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse module id", http.StatusBadRequest))
		return
	}

	var totalCount int

	err = controller.TopicService.GetAllTopics(tenantID, moduleID, sessions, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, sessions)
}

// GetTopicList returns all the sessions for the specified course
func (controller *TopicController) GetTopicList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetTopicList call==============================")

	topics := &[]list.ModuleTopic{}

	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	err = controller.TopicService.GetTopicList(tenantID, topics, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, topics)
}
