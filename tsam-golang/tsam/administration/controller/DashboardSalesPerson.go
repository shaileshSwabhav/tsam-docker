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

// SalesPersonDashboardController  provides methods to do Update, Delete, Add, Get operations on salesPerson dashbaord.
type SalesPersonDashboardController struct {
	SalesPersonDashboardService *service.SalesPersonDashboardService
}

// NewSalesPersonDashboardController returns new instance of SalesPersonDashboardController.
func NewSalesPersonDashboardController(service *service.SalesPersonDashboardService) *SalesPersonDashboardController {
	return &SalesPersonDashboardController{
		SalesPersonDashboardService: service,
	}
}

// RegisterRoutes registers all the routes of SalesPersonDashboard.
func (controller *SalesPersonDashboardController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("tenant/{tenantID}/dashboard/salesperson", controller.GetSalesPersonDashboardDetails).Methods(http.MethodGet)
	log.NewLogger().Info("Dashboard SalesPerson Routes Registered")
}

// GetSalesPersonDashboardDetails gets all dashboard details of salesPerson using service.
func (controller *SalesPersonDashboardController) GetSalesPersonDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetSalesPersonDashboardDetails called")
	salesPersonDashboardDetails := dashboard.SalesPersonDashboard{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	if err := controller.SalesPersonDashboardService.GetSalesPersonDashboardDetails(tenantID, &salesPersonDashboardDetails); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, salesPersonDashboardDetails)
}
