package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/web"
)

// BatchDashboardController  provides methods to do Update, Delete, Add, Get operations on course-batch dashbaord.
type BatchDashboardController struct {
	CourseBatchDashboardService *service.CourseBatchDashboardService
}

// NewCourseBatchDashboardController returns new instance of CourseBatchDashboardController.
func NewCourseBatchDashboardController(service *service.CourseBatchDashboardService) *BatchDashboardController {
	return &BatchDashboardController{
		CourseBatchDashboardService: service,
	}
}

// RegisterRoutes registers all the routes of Course-Batch Dashboard.
func (con *BatchDashboardController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/dashboard/course-batch", con.GetCourseBatchDashboardDetails).Methods(http.MethodGet)

	log.NewLogger().Info("Dashboard Batch-Course Routes Registered")
}

// GetCourseBatchDashboardDetails gets all dashboard details of course using service.
func (con *BatchDashboardController) GetCourseBatchDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetCourseBatchDashboardDetails called")
	batchDashboardDetails := dashboard.CourseDashboard{}
	if err := con.CourseBatchDashboardService.GetCourseBatchDashboardDetails(&batchDashboardDetails); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, batchDashboardDetails)
}
