package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	service "github.com/techlabs/swabhav/tsam/college/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// SeminarTopicController provides method to update, delete, add, get all, get one for seminar topics.
type SeminarTopicController struct {
	SeminarTopicService *service.SeminarTopicService
}

// NewSeminarTopicController creates new instance of SeminarTopicController.
func NewSeminarTopicController(seminarTopicService *service.SeminarTopicService) *SeminarTopicController {
	return &SeminarTopicController{
		SeminarTopicService: seminarTopicService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *SeminarTopicController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all seminar topics by seminar id.
	router.HandleFunc("/tenant/{tenantID}/seminar-topic/seminar/{seminarID}",
		controller.GetSeminarTopics).Methods(http.MethodGet)

	// Add one seminar topic.
	router.HandleFunc("/tenant/{tenantID}/seminar-topic/seminar/{seminarID}/credential/{credentialID}",
		controller.AddSeminarTopic).Methods(http.MethodPost)

	// Get one seminar topic.
	router.HandleFunc("/tenant/{tenantID}/seminar-topic/{seminarTopicID}/seminar/{seminarID}",
		controller.GetSeminarTopic).Methods(http.MethodGet)

	// Update one seminar topic.
	router.HandleFunc("/tenant/{tenantID}/seminar-topic/{seminarTopicID}/seminar/{seminarID}/credential/{credentialID}",
		controller.UpdateSeminarTopic).Methods(http.MethodPut)

	// Delete one seminar topic.
	router.HandleFunc("/tenant/{tenantID}/seminar-topic/{seminarTopicID}/seminar/{seminarID}/credential/{credentialID}",
		controller.DeleteSeminarTopic).Methods(http.MethodDelete)

	log.NewLogger().Info("Seminar Topic Routes Registered")
}

// AddSeminarTopic adds one seminar topic.
func (controller *SeminarTopicController) AddSeminarTopic(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddSeminarTopic called=======================================")

	// Create bucket.
	seminarTopic := college.Topic{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the seminarTopic variable with given data.
	if err := web.UnmarshalJSON(r, &seminarTopic); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err := seminarTopic.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	seminarTopic.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of seminarTopic.
	seminarTopic.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set seminar ID.
	seminarTopic.SeminarID, err = util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.SeminarTopicService.AddSeminarTopic(&seminarTopic); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Seminar topic added successfully")
}

//GetSeminarTopics gets all seminar topics.
func (controller *SeminarTopicController) GetSeminarTopics(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetSeminarTopics called=======================================")

	// Create bucket.
	seminarTopics := []college.TopicDTO{}

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting seminar id from param and parsing it to uuid.
	seminarID, err := util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Call get seminar topics method.
	if err := controller.SeminarTopicService.GetSeminarTopics(&seminarTopics, tenantID, seminarID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, seminarTopics)
}

//GetSeminarTopic gets one seminar topic.
func (controller *SeminarTopicController) GetSeminarTopic(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetSeminarTopic called=======================================")

	// Create bucket.
	seminarTopic := college.Topic{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set seminarTopic ID.
	seminarTopic.ID, err = util.ParseUUID(params["seminarTopicID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar topic id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	seminarTopic.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set seminar ID.
	seminarTopic.SeminarID, err = util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.SeminarTopicService.GetSeminarTopic(&seminarTopic); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, seminarTopic)
}

//UpdateSeminarTopic updates seminar topic.
func (controller *SeminarTopicController) UpdateSeminarTopic(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateSeminarTopic called=======================================")

	// Create bucket.
	seminarTopic := college.Topic{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the seminarTopic variable with given data.
	err := web.UnmarshalJSON(r, &seminarTopic)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := seminarTopic.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set seminarTopic ID to seminarTopic.
	seminarTopic.ID, err = util.ParseUUID(params["seminarTopicID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar topic id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to seminarTopic.
	seminarTopic.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of seminarTopic.
	seminarTopic.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set seminar ID to seminarTopic.
	seminarTopic.SeminarID, err = util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.SeminarTopicService.UpdateSeminarTopic(&seminarTopic); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Seminar topic updated successfully")
}

//DeleteSeminarTopic deletes one seminar topic.
func (controller *SeminarTopicController) DeleteSeminarTopic(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteSeminarTopic called=======================================")

	// Create bucket.
	seminarTopic := college.Topic{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set seminarTopic ID.
	seminarTopic.ID, err = util.ParseUUID(params["seminarTopicID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar topic id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	seminarTopic.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to seminarTopic's DeletedBy field.
	seminarTopic.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set seminarTopic ID to seminarTopic.
	seminarTopic.SeminarID, err = util.ParseUUID(mux.Vars(r)["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.SeminarTopicService.DeleteSeminarTopic(&seminarTopic); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Seminar topic deleted successfully")
}
