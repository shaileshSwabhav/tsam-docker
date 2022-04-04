package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// SalaryTrendController Provide method to Update, Delete, Add, Get Method For salary-trend.
type SalaryTrendController struct {
	SalaryTrendService *service.SalaryTrendService
}

// NewSalaryTrendController Create New Instance Of SalaryTrendController.
func NewSalaryTrendController(salaryTrendService *service.SalaryTrendService) *SalaryTrendController {
	return &SalaryTrendController{
		SalaryTrendService: salaryTrendService,
	}
}

// RegisterRoutes Register All Endpoint To Router
func (con *SalaryTrendController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/salary-trend/credential/{credentialID}", con.AddSalaryTrend).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/salary-trend/{salaryTrendID}/credential/{credentialID}", con.UpdateSalaryTrend).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/salary-trend/{salaryTrendID}/credential/{credentialID}", con.DeleteSalaryTrend).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/salary-trend/limit/{limit}/offset/{offset}", con.GetAllSalaryTrend).Methods(http.MethodGet)

	log.NewLogger().Info("Salary Trend Route Registered")

}

// AddSalaryTrend will add new salary-trend to the table.
func (con *SalaryTrendController) AddSalaryTrend(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddSalaryTrend call==============================")

	salaryTrend := general.SalaryTrend{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &salaryTrend)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	salaryTrend.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	salaryTrend.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = salaryTrend.ValidateSalaryTrend()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.SalaryTrendService.AddSalaryTrend(&salaryTrend)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Salary trend successfully added")
}

// UpdateSalaryTrend will update the specifed salary-trend in the table.
func (con *SalaryTrendController) UpdateSalaryTrend(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateSalaryTrend call==============================")

	salaryTrend := general.SalaryTrend{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &salaryTrend)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	salaryTrend.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	salaryTrend.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	salaryTrend.ID, err = util.ParseUUID(param["salaryTrendID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = salaryTrend.ValidateSalaryTrend()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.SalaryTrendService.UpdateSalaryTrend(&salaryTrend)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Salary trend successfully updated")
}

// DeleteSalaryTrend will delete the specifed salary-trend from the table.
func (con *SalaryTrendController) DeleteSalaryTrend(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteSalaryTrend call==============================")

	salaryTrend := general.SalaryTrend{}
	param := mux.Vars(r)
	var err error

	salaryTrend.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	salaryTrend.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	salaryTrend.ID, err = util.ParseUUID(param["salaryTrendID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.SalaryTrendService.DeleteSalaryTrend(&salaryTrend)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Salary trend successfully deleted")
}

// GetAllSalaryTrend will get all the salary-trend from the table.
func (con *SalaryTrendController) GetAllSalaryTrend(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllSalaryTrend call==============================")

	salaryTrends := []general.SalaryTrendDTO{}
	param := mux.Vars(r)
	var err error

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	r.ParseForm()

	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	err = con.SalaryTrendService.GetAllSalaryTrend(&salaryTrends, tenantID, limit, offset, &totalCount, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, salaryTrends)
}
