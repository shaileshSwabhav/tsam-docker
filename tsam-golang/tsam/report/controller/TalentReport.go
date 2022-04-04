package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/report"
	"github.com/techlabs/swabhav/tsam/report/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentReportController provides methods to do Update, Delete, Add, Get operations on faculty-report.
type TalentReportController struct {
	TalentReportService *service.TalentReportService
}

// NewTalentReportController creates new instance of TalentReportController.
func NewTalentReportController(service *service.TalentReportService) *TalentReportController {
	return &TalentReportController{
		TalentReportService: service,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *TalentReportController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get talent report.
	router.HandleFunc("/tenant/{tenantID}/talent-report/{talentID}",
		controller.GetTalentReport).Methods(http.MethodGet)

	log.NewLogger().Info("Faculty Report Routes Registered")
}

// GetTalentReport will return talent report.
func (controller *TalentReportController) GetTalentReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentReport Called==============================")

	// Create bucket for talent report.
	talentReport := report.TalentReport{}

	// Parse and get tennat id.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and get talnet id.
	talentID, err := util.ParseUUID(mux.Vars(r)["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Fills the form.
	r.ParseForm()

	// Call the get talent report service method.
	err = controller.TalentReportService.GetTalentReport(&talentReport, talentID, tenantID, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, talentReport)
}
