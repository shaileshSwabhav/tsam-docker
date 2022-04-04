package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	services "github.com/techlabs/swabhav/tsam/company/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// EnquiryCallRecordController provides method to update, delete, add, get all, get one for enquiry call records.
type EnquiryCallRecordController struct {
	EnquiryCallRecordService *services.EnquiryCallRecordService
	log                      log.Logger
	auth                     *security.Authentication
}

// NewEnquiryCallRecordController creates new instance of EnquiryCallRecordController.
func NewEnquiryCallRecordController(enquiryCallRecordService *services.EnquiryCallRecordService, log log.Logger, auth *security.Authentication) *EnquiryCallRecordController {
	return &EnquiryCallRecordController{
		EnquiryCallRecordService: enquiryCallRecordService,
		log:                      log,
		auth:                     auth,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *EnquiryCallRecordController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all enquiry call records by enquiry id.
	router.HandleFunc("/tenant/{tenantID}/company-enquiry-call-record/enquiry/{enquiryID}",
		controller.GetEnquiryCallRecords).Methods(http.MethodGet)

	// Add one enquiry call record.
	router.HandleFunc("/tenant/{tenantID}/company-enquiry-call-record/enquiry/{enquiryID}",
		controller.AddEnquiryCallRecord).Methods(http.MethodPost)

	// Get one enquiry call record.
	router.HandleFunc("/tenant/{tenantID}/company-enquiry-call-record/{enquiryCallRecordID}/enquiry/{enquiryID}",
		controller.GetEnquiryCallRecord).Methods(http.MethodGet)

	// Update one enquiry call record.
	router.HandleFunc("/tenant/{tenantID}/company-enquiry-call-record/{enquiryCallRecordID}/enquiry/{enquiryID}",
		controller.UpdateEnquiryCallRecord).Methods(http.MethodPut)

	// Delete one enquiry call record.
	router.HandleFunc("/tenant/{tenantID}/company-enquiry-call-record/{enquiryCallRecordID}/enquiry/{enquiryID}",
		controller.DeleteEnquiryCallRecord).Methods(http.MethodDelete)

	controller.log.Info("Enquiry Call Record Routes Registered")
}

// AddEnquiryCallRecord adds one enquiry call record.
func (controller *EnquiryCallRecordController) AddEnquiryCallRecord(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===============================AddEnquiryCallRecord called=======================================")

	// Create bucket.
	enquiryCallRecord := company.CallRecord{}

	// Get params from api.
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the enquiryCallRecord variable with given data.
	if err := web.UnmarshalJSON(r, &enquiryCallRecord); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err := enquiryCallRecord.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	// util.ParseUUID(params["tenantID"])
	enquiryCallRecord.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of enquiryCallRecord.
	// util.ParseUUID(params["credentialID"])
	enquiryCallRecord.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set enquiry ID.
	// util.ParseUUID(params["enquiryID"])
	enquiryCallRecord.EnquiryID, err = parser.GetUUID("enquiryID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.EnquiryCallRecordService.AddEnquiryCallRecord(&enquiryCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Call record added successfully")
}

//GetEnquiryCallRecords gets all enquiry call records.
func (controller *EnquiryCallRecordController) GetEnquiryCallRecords(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("===============================GetEnquiryCallRecords called=======================================")

	// Create bucket.
	enquiryCallRecords := []company.CallRecordDTO{}

	// Get params from api.
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Getting tenant id from param and parsing it to uuid.
	// util.ParseUUID(params["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting enquiry id from param and parsing it to uuid.
	// util.ParseUUID(params["enquiryID"])
	enquiryID, err := parser.GetUUID("enquiryID")
	if err != nil {
		controller.log.Error(err.Error())
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
	controller.log.Info("===============================GetEnquiryCallRecord called=======================================")

	// Create bucket.
	enquiryCallRecord := company.CallRecord{}

	// Declare err.
	var err error

	// Get params from api.
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Parse and set enquiryCallRecord ID.
	// util.ParseUUID(params["enquiryCallRecordID"])
	enquiryCallRecord.ID, err = parser.GetUUID("enquiryCallRecordID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry call record id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	// util.ParseUUID(params["tenantID"])
	enquiryCallRecord.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set enquiry ID.
	// util.ParseUUID(params["enquiryID"])
	enquiryCallRecord.EnquiryID, err = parser.GetUUID("enquiryID")
	if err != nil {
		controller.log.Error(err.Error())
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

	controller.log.Info("===============================UpdateEnquiryCallRecord called=======================================")

	// Create bucket.
	enquiryCallRecord := company.CallRecord{}

	// Get params from api.
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the enquiryCallRecord variable with given data.
	err := web.UnmarshalJSON(r, &enquiryCallRecord)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := enquiryCallRecord.Validate(); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set enquiryCallRecord ID to enquiryCallRecord.
	// util.ParseUUID(params["enquiryCallRecordID"])
	enquiryCallRecord.ID, err = parser.GetUUID("enquiryCallRecordID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry call record id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to enquiryCallRecord.
	// util.ParseUUID(params["tenantID"])
	enquiryCallRecord.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of enquiryCallRecord.
	//  util.ParseUUID(params["credentialID"])
	enquiryCallRecord.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set enquiryCallRecord ID to enquiryCallRecord.
	// util.ParseUUID(params["enquiryID"])
	enquiryCallRecord.EnquiryID, err = parser.GetUUID("enquiryID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.EnquiryCallRecordService.UpdateEnquiryCallRecord(&enquiryCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Call record updated successfully")
}

//DeleteEnquiryCallRecord deletes one enquiry call record.
func (controller *EnquiryCallRecordController) DeleteEnquiryCallRecord(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("===============================DeleteEnquiryCallRecord called=======================================")

	// Create bucket.
	enquiryCallRecord := company.CallRecord{}

	// Get params from api.
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Declare err.
	var err error

	// Parse and set enquiryCallRecord ID.
	// util.ParseUUID(params["enquiryCallRecordID"])
	enquiryCallRecord.ID, err = parser.GetUUID("enquiryCallRecordID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry call record id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	// util.ParseUUID(params["tenantID"])
	enquiryCallRecord.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to enquiryCallRecord's DeletedBy field.
	// util.ParseUUID(params["credentialID"])
	enquiryCallRecord.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set enquiryCallRecord ID to enquiryCallRecord.
	// util.ParseUUID(mux.Vars(r)["enquiryID"])
	enquiryCallRecord.EnquiryID, err = parser.GetUUID("enquiryID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.EnquiryCallRecordService.DeleteEnquiryCallRecord(&enquiryCallRecord); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Call record deleted successfully")
}
