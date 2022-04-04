package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/report"
	"github.com/techlabs/swabhav/tsam/report/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// PackageSummaryController provides methods to do Update, Delete, Add, Get operations on package-summary.
type PackageSummaryController struct {
	PackageSummaryService *service.PackageSummaryService
}

// NewPackageSummaryController creates new instance of PackageSummaryController.
func NewPackageSummaryController(service *service.PackageSummaryService) *PackageSummaryController {
	return &PackageSummaryController{
		PackageSummaryService: service,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *PackageSummaryController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// get
	router.HandleFunc("/tenant/{tenantID}/package-summary-report", controller.GetPackageSummary).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/package-summary-report/technology",
		controller.GetTechnologyPackageSummary).Methods(http.MethodGet)

	log.NewLogger().Info("Package Summary Routes Registered")
}

// GetPackageSummary will return package summary for different experiences.
func (controller *PackageSummaryController) GetPackageSummary(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetPackageSummary Called==============================")
	packageSummary := []report.PackageSummary{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	err = controller.PackageSummaryService.GetPackageSummary(tenantID, &packageSummary, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, packageSummary)
}

// GetTechnologyPackageSummary will give package-summary based on technology.
func (controller *PackageSummaryController) GetTechnologyPackageSummary(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTechnologyPackageSummary Called==============================")
	technologySummary := []report.ExperienceTechnologySummary{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	err = controller.PackageSummaryService.GetTechnologyPackageSummary(tenantID, &technologySummary, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, technologySummary)
}
