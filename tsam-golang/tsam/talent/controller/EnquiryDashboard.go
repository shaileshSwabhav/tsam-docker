package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/web"
)

// EnquiryDashboardController  provides methods to do Update, Delete, Add, Get operations on enquiry dashbaord.
type EnquiryDashboardController struct {
	EnquiryDashboardService *service.EnquiryDashboardService
}

// NewEnquiryDashboardController returns new instance of EnquiryDashboardController.
func NewEnquiryDashboardController(service *service.EnquiryDashboardService) *EnquiryDashboardController {
	return &EnquiryDashboardController{
		EnquiryDashboardService: service,
	}
}

// RegisterRoutes registers all the routes of EnquiryDashboard.
func (controller *EnquiryDashboardController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/dashboard/enquiry", controller.GetEnquiryDashboardDetails).Methods(http.MethodGet)
	log.NewLogger().Info("Dashboard Enquiry Routes Registered")
}

// GetEnquiryDashboardDetails gets all dashboard details of enquiry using service.
func (controller *EnquiryDashboardController) GetEnquiryDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetEnquiryDashboardDetails called")
	enquiryDashboardDetails := dashboard.EnquiryDashboard{}
	if err := controller.EnquiryDashboardService.GetEnquiryDashboardDetails(&enquiryDashboardDetails); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, enquiryDashboardDetails)
}
