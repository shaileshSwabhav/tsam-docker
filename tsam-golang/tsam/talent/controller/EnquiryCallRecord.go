package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	calrec "github.com/techlabs/swabhav/tsam/models/talentenquiry"
	talenq "github.com/techlabs/swabhav/tsam/models/talentenquiry"
	service "github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// EnquiryCallRecordController provides method to update, delete, add, get all, get one for enquiry call records.
type EnquiryCallRecordController struct {
	EnquiryCallRecordService *service.EnquiryCallRecordService
}

// NewEnquiryCallRecordController creates new instance of EnquiryCallRecordController.
func NewEnquiryCallRecordController(enquiryCallRecordService *service.EnquiryCallRecordService) *EnquiryCallRecordController {
	return &EnquiryCallRecordController{
		EnquiryCallRecordService: enquiryCallRecordService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *EnquiryCallRecordController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all enquiry call records by enquiry id.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry-call-record/talent-enquiry/{enquiryID}",
		controller.GetEnquiryCallRecords).Methods(http.MethodGet)

	// Add one enquiry call record.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry-call-record/talent-enquiry/{enquiryID}/credential/{credentialID}",
		controller.AddEnquiryCallRecord).Methods(http.MethodPost)

	// Get one enquiry call record.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry-call-record/{enquiryCallRecordID}/talent-enquiry/{enquiryID}",
		controller.GetEnquiryCallRecord).Methods(http.MethodGet)

	// Update one enquiry call record.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry-call-record/{enquiryCallRecordID}/talent-enquiry/{enquiryID}/credential/{credentialID}",
		controller.UpdateEnquiryCallRecord).Methods(http.MethodPut)

	// Delete one enquiry call record.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry-call-record/{enquiryCallRecordID}/talent-enquiry/{enquiryID}/credential/{credentialID}",
		controller.DeleteEnquiryCallRecord).Methods(http.MethodDelete)

	log.NewLogger().Info("EnquiryCallRecord Routes Registered")
}

// AddEnquiryCallRecord adds one enquiry call record.
func (controller *EnquiryCallRecordController) AddEnquiryCallRecord(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddEnquiryCallRecord called=======================================")

	// Create bucket.
	enquiryCallRecord := calrec.CallRecord{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the enquiryCallRecord variable with given data.
	if err := web.UnmarshalJSON(r, &enquiryCallRecord); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err := enquiryCallRecord.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	enquiryCallRecord.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of enquiryCallRecord.
	enquiryCallRecord.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set enquiry ID.
	enquiryCallRecord.EnquiryID, err = util.ParseUUID(params["enquiryID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.EnquiryCallRecordService.AddEnquiryCallRecord(&enquiryCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiry call record added successfully")
}

//GetEnquiryCallRecords gets all enquiry call records.
func (controller *EnquiryCallRecordController) GetEnquiryCallRecords(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetEnquiryCallRecords called=======================================")

	// Create bucket.
	enquiryCallRecords := []talenq.CallRecordDTO{}

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting enquiry id from param and parsing it to uuid.
	enquiryID, err := util.ParseUUID(params["enquiryID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Call get enquiry call records method.
	if err := controller.EnquiryCallRecordService.GetEnquiryCallRecords(&enquiryCallRecords, tenantID, enquiryID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, enquiryCallRecords)
}

//GetEnquiryCallRecord gets one enquiry call record.
func (controller *EnquiryCallRecordController) GetEnquiryCallRecord(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetEnquiryCallRecord called=======================================")

	// Create bucket.
	enquiryCallRecord := talenq.CallRecord{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set enquiryCallRecord ID.
	enquiryCallRecord.ID, err = util.ParseUUID(params["enquiryCallRecordID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry call record id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	enquiryCallRecord.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set enquiry ID.
	enquiryCallRecord.EnquiryID, err = util.ParseUUID(params["enquiryID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.EnquiryCallRecordService.GetEnquiryCallRecord(&enquiryCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, enquiryCallRecord)
}

//UpdateEnquiryCallRecord updates enquiry call record.
func (controller *EnquiryCallRecordController) UpdateEnquiryCallRecord(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateEnquiryCallRecord called=======================================")

	// Create bucket.
	enquiryCallRecord := talenq.CallRecord{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the enquiryCallRecord variable with given data.
	err := web.UnmarshalJSON(r, &enquiryCallRecord)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := enquiryCallRecord.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set enquiryCallRecord ID to enquiryCallRecord.
	enquiryCallRecord.ID, err = util.ParseUUID(params["enquiryCallRecordID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry call record id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to enquiryCallRecord.
	enquiryCallRecord.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of enquiryCallRecord.
	enquiryCallRecord.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set enquiryCallRecord ID to enquiryCallRecord.
	enquiryCallRecord.EnquiryID, err = util.ParseUUID(params["enquiryID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.EnquiryCallRecordService.UpdateEnquiryCallRecord(&enquiryCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Enquiry call record updated successfully")
}

//DeleteEnquiryCallRecord deletes one enquiry call record.
func (controller *EnquiryCallRecordController) DeleteEnquiryCallRecord(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteEnquiryCallRecord called=======================================")

	// Create bucket.
	enquiryCallRecord := talenq.CallRecord{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set enquiryCallRecord ID.
	enquiryCallRecord.ID, err = util.ParseUUID(params["enquiryCallRecordID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry call record id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	enquiryCallRecord.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to enquiryCallRecord's DeletedBy field.
	enquiryCallRecord.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set enquiryCallRecord ID to enquiryCallRecord.
	enquiryCallRecord.EnquiryID, err = util.ParseUUID(mux.Vars(r)["enquiryID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Call delete service method
	if err := controller.EnquiryCallRecordService.DeleteEnquiryCallRecord(&enquiryCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiry call record deleted successfully")
}
