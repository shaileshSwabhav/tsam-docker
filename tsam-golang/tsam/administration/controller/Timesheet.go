package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/administration/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TimesheetTestController provides methods to do CRUD operations.
type TimesheetTestController struct {
	TimesheetService *service.TimesheetService
}

// TimesheetController creates new instance of Timesheet controller.
func TimesheetController(generalService *service.TimesheetService) *TimesheetTestController {
	return &TimesheetTestController{
		TimesheetService: generalService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *TimesheetTestController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// get
	router.HandleFunc("/tenant/{tenantID}/new-timesheet/limit/{limit}/offset/{offset}",
		controller.GetTimesheets).Methods(http.MethodGet)

	// add
	router.HandleFunc("/tenant/{tenantID}/new-timesheet/credential/{credentialID}",
		controller.AddTimesheet).Methods(http.MethodPost)

	// add multiple
	router.HandleFunc("/tenant/{tenantID}/new-timesheets/credential/{credentialID}",
		controller.AddTimesheets).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/new-timesheet/{timesheetID}/credential/{credentialID}",
		controller.UpdateTimesheet).Methods(http.MethodPut)

	router.HandleFunc("/tenant/{tenantID}/timesheet-activity/{timesheetActivityID}/credential/{credentialID}",
		controller.UpdateTimesheetActivity).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/new-timesheet/{timesheetID}/credential/{credentialID}",
		controller.DeleteTimesheet).Methods(http.MethodDelete)

	log.NewLogger().Info("timesheet Route Registered")
}

// Need to add pagination if below API is needed.
// // GetAllTimesheets fetches all timesheets.
// func (controller *TimesheetTestController) GetAllTimesheets(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetAllTimesheets called==============================")

// 	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	timesheets := &[]admin.Timesheet{}
// 	err = controller.TimesheetService.GetAllTimesheets(tenantID, timesheets)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}

// 	// Writing response with OK status
// 	web.RespondJSON(w, http.StatusOK, timesheets)
// }

// GetTimesheets fetches specific timesheets.
func (controller *TimesheetTestController) GetTimesheets(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTimesheets called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	var TimesheetHeader admin.TimesheetHeader
	r.ParseForm()
	TimesheetHeader.Limit, TimesheetHeader.Offset = web.GetLimitAndOffset(r)
	timesheets := &[]*admin.TimesheetDTO{}
	err = controller.TimesheetService.GetTimesheets(tenantID, timesheets, r.Form, &TimesheetHeader)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.SetNewHeader(w, "Total-Hours", fmt.Sprintf("%v", TimesheetHeader.TotalHours))
	web.SetNewHeader(w, "Free-Hours", fmt.Sprintf("%v", TimesheetHeader.FreeHours))

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSONWithXTotalCount(w, http.StatusOK, TimesheetHeader.TotalCount, timesheets)
}

// API is working but get specific can be used to serve the purpose(search)
// // GetAllTimesheetsOfCredential fetches all timesheets of credential.
// func (controller *TimesheetTestController) GetAllTimesheetsOfCredential(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetAllTimesheetsOfCredential called==============================")

// 	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	timesheets := &[]admin.Timesheet{}
// 	err = controller.TimesheetService.GetAllTimesheetsOfCredential(tenantID, credentialID, timesheets)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}

// 	// Writing response with OK status
// 	web.RespondJSON(w, http.StatusOK, timesheets)
// }

// AddTimesheet will add new Timesheet to the table
func (controller *TimesheetTestController) AddTimesheet(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddTimesheet called==============================")
	param := mux.Vars(r)
	timesheet := admin.Timesheet{}

	err := web.UnmarshalJSON(r, &timesheet)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	timesheet.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	timesheet.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = timesheet.ValidateTimesheet()
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.TimesheetService.AddTimesheet(&timesheet)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, timesheet.ID)
}

// AddTimesheets if used to add multiple timesheets.
func (controller *TimesheetTestController) AddTimesheets(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddTimesheets called==============================")

	params := mux.Vars(r)
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	credentialID, err := util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	timesheets := []*admin.Timesheet{}
	// Parse timesheet from request.
	err = web.UnmarshalJSON(r, &timesheets)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Validate every timesheet entry.
	for _, timesheet := range timesheets {
		if err := timesheet.ValidateTimesheet(); err != nil {
			web.RespondError(w, err)
			return
		}
		timesheet.CreatedBy = credentialID
		timesheet.TenantID = tenantID
	}

	err = controller.TimesheetService.AddTimesheets(timesheets)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// // IDCollection will have the list of the UUIDs of the newly added timesheets.
	// IDCollection := []uuid.UUID{}
	// for _, timesheet := range timesheets {
	// 	if timesheet != nil {
	// 		IDCollection = append(IDCollection, timesheet.ID)
	// 	}
	// }

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "All timesheets added.")
}

// UpdateTimesheet will update the specified timesheet
func (controller *TimesheetTestController) UpdateTimesheet(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateTimesheet called==============================")
	param := mux.Vars(r)
	timesheet := admin.Timesheet{}

	err := web.UnmarshalJSON(r, &timesheet)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	timesheet.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	timesheet.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	timesheet.ID, err = util.ParseUUID(param["timesheetID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = timesheet.ValidateTimesheet()
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.TimesheetService.UpdateTimesheet(&timesheet)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Timesheet updated successfully")
}

// DeleteTimesheet will update the specified timesheet
func (controller *TimesheetTestController) DeleteTimesheet(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteTimesheet called==============================")
	param := mux.Vars(r)
	timesheet := admin.Timesheet{}
	var err error

	timesheet.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	timesheet.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	timesheet.ID, err = util.ParseUUID(param["timesheetID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.TimesheetService.DeleteTimesheet(&timesheet)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Timesheet deleted successfully")
}

// UpdateTimesheetActivity will update specific timesheet activity
func (controller *TimesheetTestController) UpdateTimesheetActivity(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateTimesheetActivity called==============================")
	param := mux.Vars(r)
	timesheetActivity := admin.TimesheetActivity{}

	err := web.UnmarshalJSON(r, &timesheetActivity)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	timesheetActivity.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	timesheetActivity.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	timesheetActivity.ID, err = util.ParseUUID(param["timesheetActivityID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = timesheetActivity.ValidateTimesheetActivity()
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.TimesheetService.UpdateTimesheetActivity(&timesheetActivity)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Timesheet activity updated successfully")
}
