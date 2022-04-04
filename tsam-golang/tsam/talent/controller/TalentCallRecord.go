package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	calrec "github.com/techlabs/swabhav/tsam/models/talent"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	service "github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentCallRecordController provides method to update, delete, add, get all, get one for talent call records.
type TalentCallRecordController struct {
	TalentCallRecordService *service.TalentCallRecordService
}

// NewTalentCallRecordController creates new instance of TalentCallRecordController.
func NewTalentCallRecordController(talentCallRecordService *service.TalentCallRecordService) *TalentCallRecordController {
	return &TalentCallRecordController{
		TalentCallRecordService: talentCallRecordService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *TalentCallRecordController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all talent call records by talent id.
	router.HandleFunc("/tenant/{tenantID}/talent-call-record/talent/{talentID}",
		controller.GetTalentCallRecords).Methods(http.MethodGet)

	// Add one talent call record.
	router.HandleFunc("/tenant/{tenantID}/talent-call-record/talent/{talentID}/credential/{credentialID}",
		controller.AddTalentCallRecord).Methods(http.MethodPost)

	// Get one talent call record.
	router.HandleFunc("/tenant/{tenantID}/talent-call-record/{talentCallRecordID}/talent/{talentID}",
		controller.GetTalentCallRecord).Methods(http.MethodGet)

	// Update one talent call record.
	router.HandleFunc("/tenant/{tenantID}/talent-call-record/{talentCallRecordID}/talent/{talentID}/credential/{credentialID}",
		controller.UpdateTalentCallRecord).Methods(http.MethodPut)

	// Delete one talent call record.
	router.HandleFunc("/tenant/{tenantID}/talent-call-record/{talentCallRecordID}/talent/{talentID}/credential/{credentialID}",
		controller.DeleteTalentCallRecord).Methods(http.MethodDelete)

	log.NewLogger().Info("TalentCallRecord Routes Registered")
}

// AddTalentCallRecord adds one talent call record.
func (controller *TalentCallRecordController) AddTalentCallRecord(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddTalentCallRecord called=======================================")

	// Create bucket.
	talentCallRecord := calrec.CallRecord{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the talentCallRecord variable with given data.
	if err := web.UnmarshalJSON(r, &talentCallRecord); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err := talentCallRecord.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	talentCallRecord.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of talentCallRecord.
	talentCallRecord.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentCallRecord.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.TalentCallRecordService.AddTalentCallRecord(&talentCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Call record added successfully")
}

//GetTalentCallRecords gets all talent call records.
func (controller *TalentCallRecordController) GetTalentCallRecords(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTalentCallRecords called=======================================")

	// Create bucket.
	talentCallRecords := []tal.CallRecordDTO{}

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting talent id from param and parsing it to uuid.
	talentID, err := util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call get talent call records method.
	if err := controller.TalentCallRecordService.GetTalentCallRecords(&talentCallRecords, tenantID, talentID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, talentCallRecords)
}

//GetTalentCallRecord gets one talent call record.
func (controller *TalentCallRecordController) GetTalentCallRecord(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTalentCallRecord called=======================================")

	// Create bucket.
	talentCallRecord := tal.CallRecord{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set talentCallRecord ID.
	talentCallRecord.ID, err = util.ParseUUID(params["talentCallRecordID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent call record id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	talentCallRecord.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentCallRecord.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.TalentCallRecordService.GetTalentCallRecord(&talentCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, talentCallRecord)
}

//UpdateTalentCallRecord updates talent call record.
func (controller *TalentCallRecordController) UpdateTalentCallRecord(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateTalentCallRecord called=======================================")

	// Create bucket.
	talentCallRecord := tal.CallRecord{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the talentCallRecord variable with given data.
	err := web.UnmarshalJSON(r, &talentCallRecord)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := talentCallRecord.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set talentCallRecord ID to talentCallRecord.
	talentCallRecord.ID, err = util.ParseUUID(params["talentCallRecordID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent call record id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to talentCallRecord.
	talentCallRecord.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of talentCallRecord.
	talentCallRecord.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set talentCallRecord ID to talentCallRecord.
	talentCallRecord.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.TalentCallRecordService.UpdateTalentCallRecord(&talentCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Call record updated successfully")
}

//DeleteTalentCallRecord deletes one talent call record.
func (controller *TalentCallRecordController) DeleteTalentCallRecord(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteTalentCallRecord called=======================================")

	// Create bucket.
	talentCallRecord := tal.CallRecord{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set talentCallRecord ID.
	talentCallRecord.ID, err = util.ParseUUID(params["talentCallRecordID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent call record id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	talentCallRecord.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to talentCallRecord's DeletedBy field.
	talentCallRecord.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set talentCallRecord ID to talentCallRecord.
	talentCallRecord.TalentID, err = util.ParseUUID(mux.Vars(r)["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.TalentCallRecordService.DeleteTalentCallRecord(&talentCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Call record deleted successfully")
}
