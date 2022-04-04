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

// ProfessionalSummaryReportController provides methods to get professional summary report of talents.
type ProfessionalSummaryReportController struct {
	ProfessionalSummaryReportService *service.ProfessionalSummaryReportService
}

// NewProfessionalSummaryReportController creates new instance of ProfessionalSummaryReportController.
func NewProfessionalSummaryReportController(professionalSummaryReportService *service.ProfessionalSummaryReportService) *ProfessionalSummaryReportController {
	return &ProfessionalSummaryReportController{
		ProfessionalSummaryReportService: professionalSummaryReportService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *ProfessionalSummaryReportController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get professional summary report with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/professional-summary-report",
		controller.GetProfessionalSummaryReport).Methods(http.MethodGet)

	// Get talent counts by company name and technology.
	router.HandleFunc("/tenant/{tenantID}/professional-summary-report-tech-count",
		controller.GetProfessionalSummaryReportByTechCount).Methods(http.MethodGet)

	log.NewLogger().Info("Talent Routes Registered")
}

// GetProfessionalSummaryReport gets professional summary report of talents.
func (controller *ProfessionalSummaryReportController) GetProfessionalSummaryReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetProfessionalSummaryReport called=======================================")

	// Create bucket.
	proSummaryReport := []tal.ProfessionalSummaryReport{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parsing for query params.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	limit, offset := web.GetLimitAndOffset(r)

	// Create bucket for total counts of all categories.
	proSummaryReportCounts := tal.ProfessionalSummaryReportCounts{}

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.ProfessionalSummaryReportService.GetProfessionalSummaryReport(&proSummaryReport, tenantID, r.Form, limit, offset, &proSummaryReportCounts); err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total category counts value in header.
	web.SetNewHeader(w, "fisrtTotalCount", strconv.Itoa(int(*proSummaryReportCounts.FirstCountTotal)))
	web.SetNewHeader(w, "secondTotalCount", strconv.Itoa(int(*proSummaryReportCounts.SecondCountTotal)))
	web.SetNewHeader(w, "thirdTotalCount", strconv.Itoa(int(*proSummaryReportCounts.ThirdCountTotal)))
	web.SetNewHeader(w, "fourthTotalCount", strconv.Itoa(int(*proSummaryReportCounts.FourthCountTotal)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, proSummaryReportCounts.TotalCount, proSummaryReport)
}

// GetProfessionalSummaryReportByTechCount gets talent counts by company name and technology.
func (controller *ProfessionalSummaryReportController) GetProfessionalSummaryReportByTechCount(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetProfessionalSummaryReportByTechCount called=======================================")

	// Create bucket.
	companyTechTalents := []tal.CompanyTechnologyTalent{}

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
	if err := controller.ProfessionalSummaryReportService.GetProfessionalSummaryReportByTechCount(&companyTechTalents, tenantID, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, companyTechTalents)
}
