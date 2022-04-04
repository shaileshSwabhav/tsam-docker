package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/report"
	"github.com/techlabs/swabhav/tsam/report/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// FresherSummaryController provides methods to do Update, Delete, Add, Get operations on fresher-summary.
type FresherSummaryController struct {
	FresherSummaryService *service.FresherSummaryService
}

// NewFresherSummaryController creates new instance of FresherSummaryController.
func NewFresherSummaryController(service *service.FresherSummaryService) *FresherSummaryController {
	return &FresherSummaryController{
		FresherSummaryService: service,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *FresherSummaryController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// get
	router.HandleFunc("/tenant/{tenantID}/fresher-summary-report", controller.GetAllFresherSummary).Methods(http.MethodGet)
	// router.HandleFunc("/tenant/{tenantID}/fresher-summary-report/technology", con.GetSummaryForTechnology).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/fresher-summary-report/technology", controller.GetSummaryForAcademicTechnology).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/fresher-summary-report/technology-list",
		controller.GetSummaryTechnologyList).Methods(http.MethodGet)

	log.NewLogger().Info("Fresher Summary Routes Registered")
}

// GetAllFresherSummary will return fresher-summary
func (controller *FresherSummaryController) GetAllFresherSummary(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllFresherSummary Called==============================")
	fresherSummary := []report.FresherSummary{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Fills the form.
	r.ParseForm()

	err = controller.FresherSummaryService.GetAllFresherSummary(tenantID, &fresherSummary, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, fresherSummary)
}

// // GetSummaryForTechnology will get technology wise fresher-summary for all academic year.
// func (con *FresherSummaryController) GetSummaryForTechnology(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetSummaryForTechnology Called==============================")
// 	technologySummary := []report.TechnologySummary{}

// 	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	// Fills the form.
// 	r.ParseForm()

// 	err = con.FresherSummaryService.GetSummaryForTechnology(tenantID, &technologySummary, r.Form)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, technologySummary)
// }

// GetSummaryTechnologyList will return list of technologies used in fresher summary
func (controller *FresherSummaryController) GetSummaryTechnologyList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetSummaryTechnologyList Called==============================")
	technologies := []general.Technology{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FresherSummaryService.GetSummaryTechnologyList(tenantID, &technologies)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, technologies)
}

// GetSummaryForAcademicTechnology will get technology wise fresher-summary for specified academic year.
func (controller *FresherSummaryController) GetSummaryForAcademicTechnology(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetSummaryForAcademicTechnology Called==============================")
	technologySummary := []report.AcademicTechnologySummary{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Fills the form.
	r.ParseForm()

	err = controller.FresherSummaryService.GetSummaryForAcademicTechnology(tenantID, &technologySummary, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, technologySummary)
}
