package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/faculty/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// DashboardController  provides methods to do Update, Delete, Add, Get operations on faculty dashbaord.
type DashboardController struct {
	dashboardService *service.DashboardService
}

// NewFacultyDashboardController returns new instance of FacultyDashboardController.
func NewFacultyDashboardController(service *service.DashboardService) *DashboardController {
	return &DashboardController{
		dashboardService: service,
	}
}

// RegisterRoutes registers all the routes of FacultyDashboard.
func (controller *DashboardController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/dashboard/faculty", controller.GetFacultyDashboardDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/faculty/{facultyID}/faculty-batch-details",
		controller.GetFacultyBatchDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/faculty/{facultyID}/ongoing-batch-details",
		controller.GetOngoingBatchDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/faculty/{facultyID}/barchart",
		controller.GetBarchartData).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/faculty/{facultyID}/piechart",
		controller.GetPiechartData).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/faculty/credential/{credentialID}/task-list",
		controller.GetTaskList).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/dashboard/batch/{batchID}/talent-feedbacks",
		controller.getTalentFeedbackScore).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/faculty-weekly-rating",
		controller.getWeeklyAvgRating).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/weekly-feedback-rating",
		controller.getSessionFeedbackRating).Methods(http.MethodGet)

	log.NewLogger().Info("Faculty Dashboard Routes Registered")
}

// GetFacultyDashboardDetails gets all dashboard details of faculty using service.
func (controller *DashboardController) GetFacultyDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetFacultyDashboardDetails called")
	facultyDashboardDetails := dashboard.FacultyDashboard{}
	if err := controller.dashboardService.GetFacultyDashboardDetails(&facultyDashboardDetails); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, facultyDashboardDetails)
}

// GetFacultyBatchDetails will return total count of all batches for specified faculty.
func (controller *DashboardController) GetFacultyBatchDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetFacultyBatchDetails Call==============================")
	facultyBatch := new(faculty.FacultyBatch)

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	facultyID, err := util.ParseUUID(mux.Vars(r)["facultyID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse form.
	r.ParseForm()

	err = controller.dashboardService.GetFacultyBatchDetails(tenantID, facultyID, facultyBatch, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, facultyBatch)
}

// GetOngoingBatchDetails will return ongoing batch details for specified faculty.
func (controller *DashboardController) GetOngoingBatchDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetOngoingBatchDetails Call==============================")
	batchDetails := new([]faculty.OngoingBatchDetails)

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	facultyID, err := util.ParseUUID(mux.Vars(r)["facultyID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse form.
	r.ParseForm()

	err = controller.dashboardService.GetOngoingBatchDetails(tenantID, facultyID, batchDetails, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, batchDetails)
}

// GetBarchartData will return data for barchart.
func (controller *DashboardController) GetBarchartData(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetBarchartData call==============================")

	barchartData := faculty.BarchartData{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	facultyID, err := util.ParseUUID(mux.Vars(r)["facultyID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	r.ParseForm()

	err = controller.dashboardService.GetBarchartData(tenantID, facultyID, &barchartData, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, barchartData)
}

// GetPiechartData will return project and its count in a week.
func (controller *DashboardController) GetPiechartData(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetPiechartData call==============================")

	piechartData := []faculty.PiechartData{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	facultyID, err := util.ParseUUID(mux.Vars(r)["facultyID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	r.ParseForm()

	err = controller.dashboardService.GetPiechartData(tenantID, facultyID, &piechartData, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, piechartData)
}

// GetTaskList will get activites from timesheet for specified credentialID.
func (controller *DashboardController) GetTaskList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTaskList call==============================")

	activities := []admin.TimesheetActivity{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	r.ParseForm()

	err = controller.dashboardService.GetTaskList(tenantID, credentialID, &activities, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, activities)
}

func (controller *DashboardController) getTalentFeedbackScore(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getTalentFeedbackScore call ==============================")

	feedbacks := faculty.Feedback{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := util.ParseUUID(mux.Vars(r)["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse batch id", http.StatusBadRequest))
		return
	}

	r.ParseForm()

	err = controller.dashboardService.GetTalentFeedbackScore(&feedbacks, tenantID, batchID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, feedbacks)
}

// getWeeklyAvgRating will return average rating of faculty for specified batch.
func (controller *DashboardController) getWeeklyAvgRating(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getWeeklyAvgRating call ==============================")

	rating := faculty.WeeklyAvgRating{}
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse batch id", http.StatusBadRequest))
		return
	}

	err = controller.dashboardService.GetWeeklyAvgRating(&rating, tenantID, batchID, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		if gorm.ErrRecordNotFound == err {
			web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusNoContent))
			return
		}
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, rating)
}

func (controller *DashboardController) getSessionFeedbackRating(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getSessionFeedbackRating call ==============================")

	rating := []faculty.WeeklyAvgRating{}
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse batch id", http.StatusBadRequest))
		return
	}

	err = controller.dashboardService.GetSessionFeedbackRating(&rating, tenantID, batchID, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		// if gorm.ErrRecordNotFound == err {
		// 	web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusNoContent))
		// 	return
		// }
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, rating)
}
