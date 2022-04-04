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

// NextActionReportController provides methods to Get nextaction report.
type NextActionReportController struct {
	NextActionReportService *service.NextActionReportService
}

// New creates new instance of NextActionReportController.
func New(ser *service.NextActionReportService) *NextActionReportController {
	return &NextActionReportController{
		NextActionReportService: ser,
	}
}

// RegisterRoutes registers all endpoints To router.
func (con *NextActionReportController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	router.HandleFunc("/tenant/{tenantID}/talent-next-action-report/limit/{limit}/offset/{offset}",
		con.GetTalentNextActionReports).Methods(http.MethodGet)
	// router.HandleFunc("/tenant/{tenantID}/talent-next-action-report/search/limit/{limit}/offset/{offset}",
	// 	con.GetTalentNextActionSearchReports).Methods(http.MethodPost)

	log.NewLogger().Info("Next Action Report Routes Registered")
}

// GetTalentNextActionReports gets talent nextAction Reports with limit and offset based on search.
func (con *NextActionReportController) GetTalentNextActionReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentNextActionReports Called==============================")

	nextActionReports := &[]report.TalentNextActionReportDTO{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Fills the form.
	r.ParseForm()

	err = con.NextActionReportService.GetTalentNextActionReports(tenantID, nextActionReports, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, nextActionReports)
}

// // GetTalentNextActionSearchReports gets searched talent nextAction Reports with limit and offset based on search.
// func (con *NextActionReportController) GetTalentNextActionSearchReports(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetTalentNextActionSearchReports Called==============================")

// 	nextActionReports := &[]report.TalentNextActionReportDTO{}
// 	search := report.NextActionSearch{}

// 	err := web.UnmarshalJSON(r, &search)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}

// 	// limit,offset & totalCount for pagination
// 	var totalCount int
// 	limit, offset := web.GetLimitAndOffset(r)

// 	err = con.NextActionReportService.GetTalentNextActionSearchReports(tenantID, nextActionReports, &search, limit, offset, &totalCount)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, nextActionReports)
// }
