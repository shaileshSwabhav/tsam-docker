package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/college/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/college"

	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// SeminarController provides methods to do update, delete, add, get all and get one operations on seminar.
type SeminarController struct {
	SeminarService *service.SeminarService
}

// NewSeminarController creates new instance of SeminarController.
func NewSeminarController(seminarService *service.SeminarService) *SeminarController {
	return &SeminarController{
		SeminarService: seminarService,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *SeminarController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Add one seminar.
	router.HandleFunc("/tenant/{tenantID}/seminar/credential/{credentialID}",
		controller.AddSeminar).Methods(http.MethodPost)

	// Get all seminar with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/seminar/limit/{limit}/offset/{offset}",
		controller.GetAllSeminars).Methods(http.MethodGet)

	// Get one seminar by id.
	router.HandleFunc("/tenant/{tenantID}/seminar/{seminarID}",
		controller.GetSeminar).Methods(http.MethodGet)

	// Update one seminar.
	router.HandleFunc("/tenant/{tenantID}/seminar/{seminarID}/credential/{credentialID}",
		controller.UpdateSeminar).Methods(http.MethodPut)

	// Delete one seminar.
	router.HandleFunc("/tenant/{tenantID}/seminar/{seminarID}/credential/{credentialID}",
		controller.DeleteSeminar).Methods(http.MethodDelete)

	// Get one seminar by code.
	seminarByCode := router.HandleFunc("/tenant/{tenantID}/seminar/code/{seminarCode}",
		controller.GetSeminarByCode).Methods(http.MethodGet)

	// Exculde routes.
	*exclude = append(*exclude, seminarByCode)

	log.NewLogger().Info("Seminar Routes Registered")
}

// GetAllSeminars gets all seminars.
func (controller *SeminarController) GetAllSeminars(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetAllSeminars call**************************************")

	// Create bucket.
	seminars := []college.SeminarDTO{}

	// Create bucket for total seminar count.
	var totalCount int

	// Fill the r.Form.
	r.ParseForm()

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Get limit and offset from param and convert it to int.
	limit, offset := web.GetLimitAndOffset(r)

	// Call get seminars method.
	err = controller.SeminarService.GetAllSeminars(&seminars, tenantID, limit, offset, &totalCount, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, seminars)
}

// GetSeminar gets seminar by calling the get service.
func (controller *SeminarController) GetSeminar(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetSeminar called=======================================")

	// Create bucket.
	seminar := college.Seminar{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set seminar ID.
	seminar.ID, err = util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	seminar.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.SeminarService.GetSeminar(&seminar); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, seminar)
}

// UpdateSeminar updates the seminar by calling the update service.
func (controller *SeminarController) UpdateSeminar(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateSeminar called=======================================")

	// Create bucket.
	seminar := college.Seminar{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the seminar variable with given data.
	err := web.UnmarshalJSON(r, &seminar)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = seminar.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set seminar ID to seminar.
	seminar.ID, err = util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to seminar.
	seminar.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of seminar.
	seminar.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.SeminarService.UpdateSeminar(&seminar)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Seminar updated successfully")
}

// AddSeminar adds new seminar by calling the add service.
func (controller *SeminarController) AddSeminar(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddSeminar called=======================================")

	// Create bucket.
	seminar := college.Seminar{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the seminar variable with given data.
	err := web.UnmarshalJSON(r, &seminar)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = seminar.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	seminar.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of seminar.
	seminar.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.SeminarService.AddSeminar(&seminar); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Seminar added successfully")
}

// DeleteSeminar deletes seminar by calling the delete service.
func (controller *SeminarController) DeleteSeminar(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteSeminar called=======================================")

	// Create bucket.
	seminar := college.Seminar{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set seminar ID.
	seminar.ID, err = util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	seminar.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to seminar's DeletedBy field.
	seminar.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.SeminarService.DeleteSeminar(&seminar)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Seminar deleted successfully")

}

// GetSeminarByCode gets seminar by its code.
func (controller *SeminarController) GetSeminarByCode(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetSeminarByCode called=======================================")

	// Create bucket.
	seminar := college.Seminar{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set seminar ID.
	seminar.Code = (params["seminarCode"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar code", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	seminar.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.SeminarService.GetSeminarByCode(&seminar); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, seminar)
}
