package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/college/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/web"
)

// CollegeDashboardController  provides methods to do Update, Delete, Add, Get operations on college dashbaord.
type CollegeDashboardController struct {
	CollegeDashboardService *service.CollegeDashboardService
}

// NewCollegeDashboardController returns new instance of CollegeDashboardController.
func NewCollegeDashboardController(service *service.CollegeDashboardService) *CollegeDashboardController {
	return &CollegeDashboardController{
		CollegeDashboardService: service,
	}
}

// RegisterRoutes registers all the routes of CollegeDashboard.
func (controller *CollegeDashboardController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/dashboard/college", controller.GetCollegeDashboardDetails).Methods(http.MethodGet)
	log.NewLogger().Info("Dashboard College Routes Registered")
}

// GetCollegeDashboardDetails gets all dashboard details of college using service.
func (controller *CollegeDashboardController) GetCollegeDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetCollegeDashboardDetails called")
	collegeDashboardDetails := dashboard.CollegeDashboard{}
	if err := controller.CollegeDashboardService.GetCollegeDashboardDetails(&collegeDashboardDetails); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, collegeDashboardDetails)
}
