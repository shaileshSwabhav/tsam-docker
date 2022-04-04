package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/administration/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// AdminDashboardController  provides methods to do Update, Delete, Add, Get operations on admin dashbaord.
type AdminDashboardController struct {
	AdminDashboardService *service.AdminDashboardService
}

// NewAdminDashboardController returns new instance of AdminDashboardController.
func NewAdminDashboardController(service *service.AdminDashboardService) *AdminDashboardController {
	return &AdminDashboardController{
		AdminDashboardService: service,
	}
}

// RegisterRoutes registers all the routes of AdminDashboard.
func (controller *AdminDashboardController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/tenant/{tenantID}/dashboard/admin",
		controller.GetAdminDashboardDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/admin/sales-people",
		controller.getSalesPeopleDashboardDetails).
		Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/admin/talent",
		controller.getTalentDashboardDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/admin/talent-enquiry",
		controller.getTalentEnquiryDashboardDetails).
		Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/admin/faculty",
		controller.getFacultyDashboardDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/admin/college",
		controller.getCollegeDashboardDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/admin/company",
		controller.getCompanyDashboardDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/admin/course",
		controller.getCourseDashboardDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/admin/batch",
		controller.getBatchDashboardDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/admin/technology",
		controller.getTechnologyDashboardDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/batch/score", controller.GetBatchPerformances).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/batch/{batchID}/talent/{talentID}/session/score",
		controller.GetSessionWiseTalentScore).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/batch/status", controller.GetBatchStatusDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/batch/{batchID}/status", controller.GetBatchTalentFeedbackScore).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/talent-enquiry-source", controller.GetEnquirySourceCount).Methods(http.MethodGet)

	log.NewLogger().Info("Dashboard Admin Routes Registered")
}

// getCourseDashboardDetails gets all dashboard details of admin using service.
func (controller *AdminDashboardController) getCourseDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("getCourseBatchDashboardDetails called")
	dashboard := dashboard.CourseDashboard{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	if err := controller.AdminDashboardService.GetCourseDashboardDetails(tenantID, &dashboard); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, dashboard)
}

// getBatchDashboardDetails gets all dashboard details of admin using service.
func (controller *AdminDashboardController) getBatchDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("getBatchBatchDashboardDetails called")
	dashboard := dashboard.BatchDashboard{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	if err := controller.AdminDashboardService.GetBatchDashboardDetails(tenantID, &dashboard); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, dashboard)
}

// getCompanyDashboardDetails gets all dashboard details of admin using service.
func (controller *AdminDashboardController) getCompanyDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("getCompanyDashboardDetails called")
	dashboard := dashboard.CompanyDashboard{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	if err := controller.AdminDashboardService.GetCompanyDashboardDetails(tenantID, &dashboard); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, dashboard)
}

// getCollegeDashboardDetails gets all dashboard details of admin using service.
func (controller *AdminDashboardController) getCollegeDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("getCollegeDashboardDetails called")
	dashboard := dashboard.CollegeDashboard{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	if err := controller.AdminDashboardService.GetCollegeDashboardDetails(tenantID, &dashboard); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, dashboard)
}

// getFacultyDashboardDetails gets all dashboard details of admin using service.
func (controller *AdminDashboardController) getFacultyDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("getFacultyDashboardDetails called")
	dashboard := dashboard.FacultyDashboard{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	if err := controller.AdminDashboardService.GetFacultyDashboardDetails(tenantID, &dashboard); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, dashboard)
}

// getTalentDashboardDetails gets all dashboard details of admin using service.
func (controller *AdminDashboardController) getTalentDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("getTalentDashboardDetails called")
	dashboard := dashboard.TalentDashboard{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	if err := controller.AdminDashboardService.GetTalentDashboardDetails(tenantID, &dashboard); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, dashboard)
}

// getTalentEnquiryDashboardDetails gets all dashboard details of admin using service.
func (controller *AdminDashboardController) getTalentEnquiryDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================getTalentEnquiryDashboardDetails call==============================")
	dashboard := dashboard.EnquiryDashboard{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	if err := controller.AdminDashboardService.GetTalentEnquiryDashboardDetails(tenantID, &dashboard, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, dashboard)
}

// getSalesPeopleDashboardDetails gets all dashboard details of admin using service.
func (controller *AdminDashboardController) getSalesPeopleDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("getSalesPeopleDashboardDetails called")
	dashboard := dashboard.SalesPersonDashboard{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Parse for salesperson ID.
	r.ParseForm()
	if err := controller.AdminDashboardService.GetSalesPeopleDashboardDetails(tenantID, r.Form, &dashboard); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, dashboard)
}

// getTechnologyDashboardDetails gets all dashboard details of admin using service.
func (controller *AdminDashboardController) getTechnologyDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("getTechnologyDashboardDetails called")
	dashboard := dashboard.TechnologyDashboard{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	if err := controller.AdminDashboardService.GetTechnologyDashboardDetails(tenantID, &dashboard); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, dashboard)
}

// GetAdminDashboardDetails gets total values of entites that are not called at init eg:Course,tech.
func (controller *AdminDashboardController) GetAdminDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetAdminDashboardDetails called")
	dashboard := dashboard.AdminDashboard{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	if err := controller.AdminDashboardService.GetAdminDashboardDetails(tenantID, &dashboard); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, dashboard)
}

// GetBatchPerformances will return talents having outstanding, good or average score
func (controller *AdminDashboardController) GetBatchPerformances(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetBatchPerformances call==============================")
	batchPerformance := dashboard.BatchPerformance{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// parse form
	r.ParseForm()

	err = controller.AdminDashboardService.GetBatchPerformances(&batchPerformance, tenantID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batchPerformance)
}

// GetSessionWiseTalentScore returns specified talent's session-wise feedback score
func (controller *AdminDashboardController) GetSessionWiseTalentScore(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetSessionWiseTalentScore call==============================")

	talentSessionFeedbackScore := dashboard.TalentSessionFeedbackScore{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	talentID, err := util.ParseUUID(mux.Vars(r)["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := util.ParseUUID(mux.Vars(r)["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	err = controller.AdminDashboardService.GetSessionWiseTalentScore(&talentSessionFeedbackScore, tenantID, talentID, batchID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, talentSessionFeedbackScore)
}

// GetBatchStatusDetails returns details of batch for specified batch status
func (controller *AdminDashboardController) GetBatchStatusDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetBatchStatusDetails call==============================")

	batchDetails := []dashboard.BatchDetails{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	err = controller.AdminDashboardService.GetBatchStatusDetails(&batchDetails, tenantID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batchDetails)
}

// GetBatchTalentFeedbackScore returns talents feedback score for specified batch
func (controller *AdminDashboardController) GetBatchTalentFeedbackScore(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetBatchTalentFeedbackScore call==============================")

	talentFeedbackScore := []dashboard.TalentFeedbackScore{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := util.ParseUUID(mux.Vars(r)["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.AdminDashboardService.GetBatchTalentFeedbackScore(&talentFeedbackScore, tenantID, batchID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, talentFeedbackScore)
}

// GetEnquirySourceCount will return count of all the enquiries from a specific source
func (controller *AdminDashboardController) GetEnquirySourceCount(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetEnquirySourceCount call==============================")

	sources := []dashboard.EnquirySource{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	err = controller.AdminDashboardService.GetEnquirySourceCount(tenantID, &sources, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, sources)
}
