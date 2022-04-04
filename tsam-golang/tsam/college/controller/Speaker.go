package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/college/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// SpeakerController provides method to update, delete, add, get for speaker.
type SpeakerController struct {
	SpeakerService *service.SpeakerService
}

// NewSpeakerController creates new instance of SpeakerController.
func NewSpeakerController(speakerService *service.SpeakerService) *SpeakerController {
	return &SpeakerController{
		SpeakerService: speakerService,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *SpeakerController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all speaker list.
	router.HandleFunc("/tenant/{tenantID}/speaker",
		controller.GetSpeakerList).Methods(http.MethodGet)

	// Get all speakers with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/speaker/limit/{limit}/offset/{offset}",
		controller.GetSpeakers).Methods(http.MethodGet)

	// Get one speaker by id.
	router.HandleFunc("/tenant/{tenantID}/speaker/{speakerID}",
		controller.GetSpeaker).Methods(http.MethodGet)

	// Add one speaker.
	router.HandleFunc("/tenant/{tenantID}/speaker/credential/{credentialID}",
		controller.AddSpeaker).Methods(http.MethodPost)

	// Add multiple speakers.
	router.HandleFunc("/tenant/{tenantID}/speakers/credential/{credentialID}",
		controller.AddSpeakers).Methods(http.MethodPost)

	// Update one speaker.
	router.HandleFunc("/tenant/{tenantID}/speaker/{speakerID}/credential/{credentialID}",
		controller.UpdateSpeaker).Methods(http.MethodPut)

	// Delete one speaker.
	router.HandleFunc("/tenant/{tenantID}/speaker/{speakerID}/credential/{credentialID}",
		controller.DeleteSpeaker).Methods(http.MethodDelete)

	log.NewLogger().Info("Speaker Route Registered")
}

// GetSpeakerList returns speaker list.
func (controller *SpeakerController) GetSpeakerList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetSpeakerList call**************************************")

	// Create bucket.
	speakers := []college.Speaker{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get speaker list method.
	err = controller.SpeakerService.GetSpeakerList(&speakers, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, speakers)
}

// GetSpeakers returns all speakers.
func (controller *SpeakerController) GetSpeakers(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetSpeakers call**************************************")

	// Create bucket.
	speakers := []college.SpeakerDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse for query params.
	r.ParseForm()

	// For pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all speakers service method.
	err = controller.SpeakerService.GetSpeakers(&speakers, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, speakers)
}

// GetSpeaker return specific speaker by id.
func (controller *SpeakerController) GetSpeaker(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetSpeaker call**************************************")

	// Create bucket.
	speaker := college.SpeakerDTO{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	speaker.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant ID", http.StatusBadRequest))
		return
	}

	// Getting speaker id from param and parsing it to uuid.
	speaker.ID, err = util.ParseUUID(mux.Vars(r)["speakerID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse speaker ID", http.StatusBadRequest))
		return
	}

	// Call get speaker by id service method.
	err = controller.SpeakerService.GetSpeaker(&speaker)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, speaker)
}

// AddSpeaker adds new speaker.
func (controller *SpeakerController) AddSpeaker(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddSpeaker call**************************************")

	// Create bucket.
	speaker := college.Speaker{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &speaker)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := speaker.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	speaker.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	speaker.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.SpeakerService.AddSpeaker(&speaker)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Speaker added successfully")
}

// AddSpeakers adds multiple speakers.
func (controller *SpeakerController) AddSpeakers(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddSpeakers call**************************************")

	// Create bucket.
	speakersIDs := []uuid.UUID{}
	speakers := []college.Speaker{}

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

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &speakers)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary speaker fields.
	for _, speaker := range speakers {
		err = speaker.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	// Call add multiple speaker service method.
	err = controller.SpeakerService.AddSpeakers(&speakers, &speakersIDs, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Speakers added successfully")
}

// UpdateSpeaker updates the specified speaker by id.
func (controller *SpeakerController) UpdateSpeaker(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateSpeaker call**************************************")

	// Create bucket.
	speaker := college.Speaker{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	speaker.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant ID", http.StatusBadRequest))
		return
	}

	// Getting speaker id from param and parsing it to uuid.
	speaker.ID, err = util.ParseUUID(mux.Vars(r)["speakerID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse speaker ID", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	speaker.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &speaker)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate country fields
	err = speaker.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.SpeakerService.UpdateSpeaker(&speaker)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Speaker updated successfully")
}

// DeleteSpeaker deletes the specified speaker by id.
func (controller *SpeakerController) DeleteSpeaker(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteSpeaker call**************************************")

	// Create bcuket.
	speaker := college.Speaker{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	speaker.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant ID", http.StatusBadRequest))
		return
	}

	// Getting speaker id from param and parsing it to uuid.
	speaker.ID, err = util.ParseUUID(mux.Vars(r)["speakerID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse speaker ID", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	speaker.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.SpeakerService.DeleteSpeaker(&speaker)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Speaker deleted successfully")
}
