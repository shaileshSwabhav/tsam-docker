package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/report"
	"github.com/techlabs/swabhav/tsam/report/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// LoginReportController provides methods to do Update, Delete, Add, Get operations on login-report.
type LoginReportController struct {
	LoginReportService *service.LoginReportService
}

// NewLoginReportController creates new instance of LoginReportController.
func NewLoginReportController(service *service.LoginReportService) *LoginReportController {
	return &LoginReportController{
		LoginReportService: service,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *LoginReportController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// get
	router.HandleFunc("/tenant/{tenantID}/login-report/limit/{limit}/offset/{offset}",
		controller.GetLoginReports).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/credential/{credentialID}/login-report/limit/{limit}/offset/{offset}",
		controller.GetCredentialLoginReports).Methods(http.MethodGet)

	log.NewLogger().Info("Login Report Routes Registered")
}

// GetCredentialLoginReports returns details of login and logout
func (controller *LoginReportController) GetCredentialLoginReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetCredentialLoginReports Called==============================")

	// loginReports := &[]report.CredentialLoginReport{}
	loginReports := new([]report.CredentialLoginReport)

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Fills the form.
	r.ParseForm()

	err = controller.LoginReportService.GetCredentialLoginReports(tenantID, credentialID, loginReports,
		limit, offset, &totalCount, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, loginReports)
}

// GetLoginReports returns details of login and logout.
func (controller *LoginReportController) GetLoginReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetLoginReports Called==============================")

	loginReports := new([]report.LoginReport)
	// loginReports := &[]report.LoginReport{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Fills the form.
	r.ParseForm()

	err = controller.LoginReportService.GetLoginReports(tenantID, loginReports, limit, offset, &totalCount, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	fmt.Println("total count in controller ->", totalCount)

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, loginReports)
}
