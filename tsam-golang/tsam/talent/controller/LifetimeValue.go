package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	service "github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentLifetimeValueController provides method to update, delete, add, get all, get one for talent lifetime value.
type TalentLifetimeValueController struct {
	TalentLifetimeValueService *service.TalentLifetimeValueService
}

// NewTalentCallRecordController creates new instance of TalentLifetimeValueController.
func NewTalentLifetimeValueController(talentCallRecordService *service.TalentLifetimeValueService) *TalentLifetimeValueController {
	return &TalentLifetimeValueController{
		TalentLifetimeValueService: talentCallRecordService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *TalentLifetimeValueController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add one talent lifetime value.
	router.HandleFunc("/tenant/{tenantID}/talent-lifetime-value/talent/{talentID}/credential/{credentialID}",
		controller.AddTalentLifetimeValue).Methods(http.MethodPost)

	// Get one talent lifetime value.
	router.HandleFunc("/tenant/{tenantID}/talent-lifetime-value/talent/{talentID}",
		controller.GetTalentLifetimeValue).Methods(http.MethodGet)

	// Get all talent lifetime values.
	router.HandleFunc("/tenant/{tenantID}/talent-lifetime-value/login/{lognID}/limit/{limit}/offset/{offset}",
		controller.GetTalentLifetimeValueReports).Methods(http.MethodGet)

	// Update one talent lifetime value.
	router.HandleFunc("/tenant/{tenantID}/talent-lifetime-value/{talentLifetimeValueID}/talent/{talentID}/credential/{credentialID}",
		controller.UpdateTalentLifetimeValue).Methods(http.MethodPut)

	// Delete one talent lifetime value.
	router.HandleFunc("/tenant/{tenantID}/talent-lifetime-value/{talentLifetimeValueID}/talent/{talentID}/credential/{credentialID}",
		controller.DeleteTalentLifetimeValue).Methods(http.MethodDelete)

	log.NewLogger().Info("TalentCallRecord Routes Registered")
}

// AddTalentLifetimeValue adds one talent lifetime value.
func (controller *TalentLifetimeValueController) AddTalentLifetimeValue(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddTalentLifetimeValue called=======================================")

	// Create bucket.
	lifetimeValue := tal.LifetimeValue{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the lifetimeValue variable with given data.
	if err := web.UnmarshalJSON(r, &lifetimeValue); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err := lifetimeValue.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	lifetimeValue.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of lifetimeValue.
	lifetimeValue.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	lifetimeValue.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.TalentLifetimeValueService.AddTalentLifetimeValue(&lifetimeValue); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Lifetime value added successfully")
}

//GetTalentLifetimeValue gets one talent lifetime value.
func (controller *TalentLifetimeValueController) GetTalentLifetimeValue(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTalentLifetimeValue called=======================================")

	// Create bucket.
	lifetimeValue := tal.LifetimeValue{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set tenant ID.
	lifetimeValue.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	lifetimeValue.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.TalentLifetimeValueService.GetTalentLifetimeValue(&lifetimeValue); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, lifetimeValue)
}

//GetTalentLifetimeValueReports gets all talent lifetime values.
func (controller *TalentLifetimeValueController) GetTalentLifetimeValueReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTalentLifetimeValueReports called=======================================")

	// Create bucket.
	lifetimeValues := []tal.LifetimeValueReport{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the r.Form.
	r.ParseForm()

	// Create bucket for total talents count.
	var totalCount int

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Parse and get tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and get login ID.
	lognID, err := util.ParseUUID(params["lognID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Get limit and offset from param and convert it to int.
	limit, offset := web.GetLimitAndOffset(r)

	// Call get service method.
	if err := controller.TalentLifetimeValueService.GetTalentLifetimeValueReports(&lifetimeValues, lognID, tenantID, limit, offset, &totalCount, &totalLifetimeValue, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, lifetimeValues)
}

//UpdateTalentLifetimeValue updates talent lifetime value.
func (controller *TalentLifetimeValueController) UpdateTalentLifetimeValue(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateTalentLifetimeValue called=======================================")

	// Create bucket.
	lifetimeValue := tal.LifetimeValue{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the lifetimeValue variable with given data.
	err := web.UnmarshalJSON(r, &lifetimeValue)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := lifetimeValue.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set lifetimeValue ID to lifetimeValue.
	lifetimeValue.ID, err = util.ParseUUID(params["talentLifetimeValueID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent lifetime value id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to lifetimeValue.
	lifetimeValue.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of lifetimeValue.
	lifetimeValue.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set lifetimeValue ID to lifetimeValue.
	lifetimeValue.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.TalentLifetimeValueService.UpdateTalentLifetimeValue(&lifetimeValue); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Lifetime value updated successfully")
}

//DeleteTalentCallRecord deletes one talent lifetime value.
func (controller *TalentLifetimeValueController) DeleteTalentLifetimeValue(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteTalentCallRecord called=======================================")

	// Create bucket.
	lifetimeValue := tal.LifetimeValue{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set lifetimeValue ID.
	lifetimeValue.ID, err = util.ParseUUID(params["talentLifetimeValueID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent lifetime value id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	lifetimeValue.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to lifetimeValue's DeletedBy field.
	lifetimeValue.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set lifetimeValue ID to lifetimeValue.
	lifetimeValue.TalentID, err = util.ParseUUID(mux.Vars(r)["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.TalentLifetimeValueService.DeleteTalentLifetimeValue(&lifetimeValue); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Lifetime value deleted successfully")
}
