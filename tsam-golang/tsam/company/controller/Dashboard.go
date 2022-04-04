package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	services "github.com/techlabs/swabhav/tsam/company/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/web"
)

// CompanyDashboardController  provides methods to do Update, Delete, Add, Get operations on company dashbaord.
type CompanyDashboardController struct {
	CompanyDashboardService *services.CompanyDashboardService
}

// NewCompanyDashboardController returns new instance of CompanyDashboardController.
func NewCompanyDashboardController(services *services.CompanyDashboardService) *CompanyDashboardController {
	return &CompanyDashboardController{
		CompanyDashboardService: services,
	}
}

// RegisterRoutes registers all the routes of CompanyDashboard.
func (controller *CompanyDashboardController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/dashboard/company", controller.GetCompanyDashboardDetails).Methods(http.MethodGet)
	log.NewLogger().Info("Dashboard Company Routes Registered")
}

// GetCompanyDashboardDetails gets all dashboard details of company using service.
func (controller *CompanyDashboardController) GetCompanyDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetCompanyDashboardDetails called")
	companyDashboardDetails := dashboard.CompanyDashboard{}
	if err := controller.CompanyDashboardService.GetCompanyDashboardDetails(&companyDashboardDetails); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, companyDashboardDetails)
}
