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

// FacultyReportController provides methods to do Update, Delete, Add, Get operations on faculty-report.
type FacultyReportController struct {
	FacultyReportService *service.FacultyReportService
}

// NewFacultyReportController creates new instance of FacultyReportController.
func NewFacultyReportController(service *service.FacultyReportService) *FacultyReportController {
	return &FacultyReportController{
		FacultyReportService: service,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *FacultyReportController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// get
	router.HandleFunc("/tenant/{tenantID}/faculty-report", controller.GetFacultyReport).Methods(http.MethodGet)

	log.NewLogger().Info("Faculty Report Routes Registered")
}

// GetFacultyReport will return fresher-summary
func (controller *FacultyReportController) GetFacultyReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetFacultyReport Called==============================")
	facultyReport := []report.FacultyReport{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Fills the form.
	r.ParseForm()

	err = controller.FacultyReportService.GetFacultyReport(tenantID, &facultyReport, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, facultyReport)
}
