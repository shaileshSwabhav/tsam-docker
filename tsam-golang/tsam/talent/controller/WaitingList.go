package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	service "github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// WaitingListController provides method to update, delete, add, get for waiting list.
type WaitingListController struct {
	WaitingListService *service.WaitingListService
}

// NewWaitingListController creates new instance of WaitingListController.
func NewWaitingListController(waitingListService *service.WaitingListService) *WaitingListController {
	return &WaitingListController{
		WaitingListService: waitingListService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *WaitingListController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add one waiting list.
	router.HandleFunc("/tenant/{tenantID}/waiting-list/credential/{credentialID}",
		controller.AddWaitingList).Methods(http.MethodPost)

	// Get all waiting lists by talent id.
	router.HandleFunc("/tenant/{tenantID}/waiting-list/talent/{talentID}",
		controller.GetWaitingListsByTalent).Methods(http.MethodGet)

	// Get all waiting lists by enquiry id.
	router.HandleFunc("/tenant/{tenantID}/waiting-list/enquiry/{enquiryID}",
		controller.GetWaitingListsByEnquiry).Methods(http.MethodGet)

	// Get two waiting lists.
	router.HandleFunc("/tenant/{tenantID}/waiting-list-two",
		controller.GetTwoWaitingLists).Methods(http.MethodGet)

	// Update one waiting list.
	router.HandleFunc("/tenant/{tenantID}/waiting-list/{waitingListID}/credential/{credentialID}",
		controller.UpdateWaitingList).Methods(http.MethodPut)

	// Update some fields of waiting lists.
	router.HandleFunc("/tenant/{tenantID}/waiting-list-transfer/credential/{credentialID}",
		controller.TransferWaitingList).Methods(http.MethodPut)

	// Delete one waiting list.
	router.HandleFunc("/tenant/{tenantID}/waiting-list/{waitingListID}/credential/{credentialID}",
		controller.DeleteWaitingList).Methods(http.MethodDelete)

	log.NewLogger().Info("Waiting List Routes Registered")
}

// AddWaitingList adds one waiting list.
func (controller *WaitingListController) AddWaitingList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddWaitingList called=======================================")

	// Create bucket.
	waitingList := tal.WaitingList{}

	// Declare error variable
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Fill the waitingList variable with given data.
	if err := web.UnmarshalJSON(r, &waitingList); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	waitingList.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of waitingList.
	waitingList.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.WaitingListService.AddWaitingList(&waitingList); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Waiting List added successfully")
}

// GetWaitingListsByTalent gets all waiting lists by talent id.
func (controller *WaitingListController) GetWaitingListsByTalent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetWaitingListsByTalent called=======================================")

	// Create bucket.
	waitingList := []tal.WaitingListDTO{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parsing for query params.
	r.ParseForm()

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentID, err := util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.WaitingListService.GetWaitingLists(&waitingList, tenantID, &talentID, nil, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, waitingList)
}

// GetWaitingListsByEnquiry gets all waiting lists by enquiry id.
func (controller *WaitingListController) GetWaitingListsByEnquiry(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetWaitingListsByEnquiry called=======================================")

	// Create bucket.
	waitingList := []tal.WaitingListDTO{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parsing for query params.
	r.ParseForm()

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set enquiry ID.
	enquiryID, err := util.ParseUUID(params["enquiryID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.WaitingListService.GetWaitingLists(&waitingList, tenantID, nil, &enquiryID, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, waitingList)
}

// GetTwoWaitingLists gets two waiting lists.
func (controller *WaitingListController) GetTwoWaitingLists(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTwoWaitingLists called=======================================")

	// Create bucket.
	twoWaitingLists := tal.TwoWaitingLists{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parsing for query params.
	r.ParseForm()

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.WaitingListService.GetTwoWaitingLists(&twoWaitingLists, tenantID, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, twoWaitingLists)
}

// UpdateWaitingList updates waiting list.
func (controller *WaitingListController) UpdateWaitingList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateWaitingList called=======================================")

	// Create bucket.
	waitingList := tal.WaitingList{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the waitingList variable with given data.
	err := web.UnmarshalJSON(r, &waitingList)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set waitingList ID to waitingList.
	waitingList.ID, err = util.ParseUUID(params["waitingListID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse waiting list id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to waitingList.
	waitingList.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of waitingList.
	waitingList.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.WaitingListService.UpdateWaitingList(&waitingList); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Waiting list updated successfully")
}

// TransferWaitingList updates some fields of waiting lists.
func (controller *WaitingListController) TransferWaitingList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================TransferWaitingList called=======================================")

	// Create bucket.
	updateWaitingList := tal.UpdateWaitingList{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the updateWaitingLists variable with given data.
	err := web.UnmarshalJSON(r, &updateWaitingList)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to waitingList.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of waitingList.
	credentialID, err := util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.WaitingListService.TransferWaitingList(&updateWaitingList, tenantID, credentialID); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Waiting list transfered successfully")
}

// DeleteWaitingList deletes one waiting list.
func (controller *WaitingListController) DeleteWaitingList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteWaitingList called=======================================")

	// Create bucket.
	waitingList := tal.WaitingList{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set waitingList ID.
	waitingList.ID, err = util.ParseUUID(params["waitingListID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse waiting list id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	waitingList.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to waitingList's DeletedBy field.
	waitingList.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.WaitingListService.DeleteWaitingList(&waitingList); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Waiting list deleted successfully")
}
