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

// WaitingListReportController provides methods to get waiting list of talents and enquiries in a report format.
type WaitingListReportController struct {
	WaitingListReportService *service.WaitingListReportService
}

// NewWaitingListReportController creates new instance of WaitingListReportController.
func NewWaitingListReportController(waitingListReportService *service.WaitingListReportService) *WaitingListReportController {
	return &WaitingListReportController{
		WaitingListReportService: waitingListReportService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *WaitingListReportController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get waiting list report by company requirement with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/waiting-list-report-company-branch/limit/{limit}/offset/{offset}",
		controller.GetCompanyBranchWaitingListReport).Methods(http.MethodGet)

	// Get and waiting list report by course with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/waiting-list-report-course/limit/{limit}/offset/{offset}",
		controller.GetCourseWaitingListReport).Methods(http.MethodGet)

	// Get waiting report list by requirement.
	router.HandleFunc("/tenant/{tenantID}/waiting-list-report-requirement/company-branch/{companyBranchID}",
		controller.GetRequirementWaitingListReport).Methods(http.MethodGet)

	// Get waiting list report by batch.
	router.HandleFunc("/tenant/{tenantID}/waiting-list-report-batch/course/{courseID}",
		controller.GetBatchWaitingListReport).Methods(http.MethodGet)

	// Get waiting list report by technology.
	router.HandleFunc("/tenant/{tenantID}/waiting-list-report-technology/limit/{limit}/offset/{offset}",
		controller.GetTechnologyWaitingListReport).Methods(http.MethodGet)

	log.NewLogger().Info("Talent Routes Registered")
}

// GetCompanyBranchWaitingListReport gets waiting list report for company branch.
func (controller *WaitingListReportController) GetCompanyBranchWaitingListReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetCompanyBranchWaitingListReport called=======================================")

	// Create bucket.
	waitingListReport := []tal.WaitingListCompanyBranchDTO{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parsing for query params.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	totalCount := tal.TotalCount{}
	limit, offset := web.GetLimitAndOffset(r)

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.WaitingListReportService.GetCompanyBranchWaitingListReport(&waitingListReport, tenantID, r.Form, limit, offset, &totalCount); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount.TotalCount, waitingListReport)
}

// GetCourseWaitingListReport gets waiting list report for course.
func (controller *WaitingListReportController) GetCourseWaitingListReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetCourseWaitingListReport called=======================================")

	// Create bucket.
	waitingListReport := []tal.WaitingListCourseDTO{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parsing for query params.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	totalCount := tal.TotalCount{}
	limit, offset := web.GetLimitAndOffset(r)

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.WaitingListReportService.GetCourseWaitingListReport(&waitingListReport, tenantID, r.Form, limit, offset, &totalCount); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount.TotalCount, waitingListReport)
}

// GetRequirementWaitingListReport gets waiting list report for requirement.
func (controller *WaitingListReportController) GetRequirementWaitingListReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetRequirementWaitingListReport called=======================================")

	// Create bucket.
	waitingListReport := []tal.WaitingListRequirementDTO{}

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

	// Parse and set company branch ID.
	companyBranchID, err := util.ParseUUID(params["companyBranchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse company branch id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.WaitingListReportService.GetRequirementWaitingListReport(&waitingListReport, tenantID, companyBranchID, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, waitingListReport)
}

// GetBatchWaitingListReport gets waiting list report for batch.
func (controller *WaitingListReportController) GetBatchWaitingListReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetBatchWaitingListReport called=======================================")

	// Create bucket.
	waitingListReport := []tal.WaitingListBatchDTO{}

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

	// Parse and set course ID.
	courseID, err := util.ParseUUID(params["courseID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse course id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.WaitingListReportService.GetBatchWaitingListReport(&waitingListReport, tenantID, courseID, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, waitingListReport)
}

// GetTechnologyWaitingListReport gets waiting list report by technology.
func (controller *WaitingListReportController) GetTechnologyWaitingListReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTechnologyWaitingListReport called=======================================")

	// Create bucket.
	waitingListReport := []tal.WaitingListTechnologyDTO{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parsing for query params.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	totalCount := tal.TotalCount{}
	limit, offset := web.GetLimitAndOffset(r)

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.WaitingListReportService.GetTechnologyWaitingListReport(&waitingListReport, tenantID, r.Form, limit, offset, &totalCount); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount.TotalCount, waitingListReport)
}

// // RegisterRoutes registers all endpoint to router.
// func (controller *WaitingListReportController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

// 	// Get all talents and waiting list by company requirement with limit and offset.
// 	router.HandleFunc("/tenant/{tenantID}/talent-waiting-list-report/company-requirement/limit/{limit}/offset/{offset}",
// 		controller.GetTalentWaitingListByCompanyRequirement).Methods(http.MethodGet)

// 	log.NewLogger().Info("Talent Routes Registered")
// }

// // GetTalentWaitingListByCompanyRequirement returns all talents and its corresponding waiting list entries by
// // company requirement id.
// func (controller *WaitingListReportController) GetTalentWaitingListByCompanyRequirement(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetTalentWaitingListByCompanyRequirement call==============================")

// 	// Create bucket.
// 	talentWaitingLists := []tal.TalentWaitingList{}

// 	// Getting tenant id from param and parsing it to uuid.
// 	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}

// 	// Limit, offset & totalCount for pagination.
// 	var totalCount int
// 	limit, offset := web.GetLimitAndOffset(r)

// 	// Call get talents and waiting lists by company requirement service method.
// 	err = controller.WaitingListReportService.GetTalentWaitingListByCompanyRequirement(&talentWaitingLists, tenantID, limit, offset, &totalCount)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}

// 	// Writing response with OK status and total count in header to ResponseWriter.
// 	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talentWaitingLists)
// }
