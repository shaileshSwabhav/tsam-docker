package callreportcntlr

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/callreport/callreportsvc"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/callreport"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CallingReportController provides methods to do Update, Delete, Add, Get operations on callingReport.
type CallingReportController struct {
	CallingReportService *callreportsvc.CallingReportService
}

// New creates new instance of CallingReportController.
func New(service *callreportsvc.CallingReportService) *CallingReportController {
	return &CallingReportController{
		CallingReportService: service,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *CallingReportController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	router.HandleFunc("/tenant/{tenantID}/talent-calling-report/loginwise",
		controller.GetAllLoginwiseTalentCallingReports).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/talent-calling-report/daywise",
		controller.GetAllDaywiseTalentCallingReports).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/talent-calling-report/limit/{limit}/offset/{offset}",
		controller.GetTalentCallingReports).Methods(http.MethodGet)

	// talent-enquiry
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry-calling-report/loginwise",
		controller.GetLoginwiseTalentEnquiryCallingReports).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/talent-enquiry-calling-report/daywise",
		controller.GetDaywiseTalentEnquiryCallingReports).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/talent-enquiry-calling-report/limit/{limit}/offset/{offset}",
		controller.GetTalentEnquiryCallingReports).Methods(http.MethodGet)

	// router.HandleFunc("/tenant/{tenantID}/talent-daywise-calling-report",
	// callingReportController.GetAllDaywiseTalentCallingReports).Methods(http.MethodGet)
	// router.HandleFunc("/tenant/{tenantID}/talent-daywise-calling-report/limit{limit}/offset/{offset}",
	// callingReportController.GetDaywiseTalentCallingReports).Methods(http.MethodGet)
	log.NewLogger().Info("CallingReport Routes Registered")
}

// GetAllLoginwiseTalentCallingReports gets all loginwise callingReports in talent
func (controller *CallingReportController) GetAllLoginwiseTalentCallingReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllLoginwiseTalentCallingReports Called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	callingReports := &[]callreport.LoginwiseCallingReport{}
	r.ParseForm()

	err = controller.CallingReportService.GetLoginwiseTalentCallingReports(callingReports, tenantID, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status
	web.RespondJSON(w, http.StatusOK, callingReports)

}

// GetAllDaywiseTalentCallingReports gets all daywise callingReports in talent
func (controller *CallingReportController) GetAllDaywiseTalentCallingReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllDaywiseTalentCallingReports Called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	callingReports := &[]callreport.DaywiseCallingReport{}
	err = controller.CallingReportService.GetDaywiseTalentCallingReports(callingReports, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status
	web.RespondJSON(w, http.StatusOK, callingReports)

}

// GetTalentCallingReports gets talent callingReports with limit and offset based on search.
func (controller *CallingReportController) GetTalentCallingReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetLoginwiseTalentCallingReports Called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)
	allCallingReports := &[]callreport.TalentCallingReportDTO{}

	// Fills the form.
	r.ParseForm()
	err = controller.CallingReportService.GetTalentCallingReports(tenantID, allCallingReports, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allCallingReports)
}

// =======================================================TALENT-ENQUIRY=======================================================

// GetLoginwiseTalentEnquiryCallingReports gets all loginwise callingReports in talent-enqury
func (controller *CallingReportController) GetLoginwiseTalentEnquiryCallingReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetLoginwiseTalentEnquiryCallingReports Called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	callingReports := &[]callreport.LoginwiseTalentEnquiryCallingReport{}
	r.ParseForm()

	err = controller.CallingReportService.GetLoginwiseTalentEnquiryCallingReports(callingReports, tenantID, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status
	web.RespondJSON(w, http.StatusOK, callingReports)
}

// GetDaywiseTalentEnquiryCallingReports gets all daywise callingReports in talent-enquiry.
func (controller *CallingReportController) GetDaywiseTalentEnquiryCallingReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetDaywiseTalentEnquiryCallingReports Called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	callingReports := &[]callreport.DaywiseTalentEnquiryCallingReport{}

	err = controller.CallingReportService.GetDaywiseTalentEnquiryCallingReports(callingReports, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status
	web.RespondJSON(w, http.StatusOK, callingReports)
}

// GetTalentEnquiryCallingReports gets talent callingReports with limit and offset based on search.
func (controller *CallingReportController) GetTalentEnquiryCallingReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetLoginwiseTalentCallingReports Called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)
	allCallingReports := &[]callreport.TalentEnquiryCallingReportDTO{}

	// Fills the form.
	r.ParseForm()
	err = controller.CallingReportService.GetTalentEnquiryCallingReports(tenantID, allCallingReports, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allCallingReports)
}
