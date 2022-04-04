package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// DashboardController  provides methods to do Update, Delete, Add, Get operations on talent dashbaord.
type DashboardController struct {
	TalentDashboardService *service.DashboardService
}

// NewTalentDashboardController returns new instance of TalentDashboardController.
func NewTalentDashboardController(service *service.DashboardService) *DashboardController {
	return &DashboardController{
		TalentDashboardService: service,
	}
}

// RegisterRoutes registers all the routes of TalentDashboard.
func (controller *DashboardController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/dashboard/talent", controller.GetTalentDashboardDetails).Methods(http.MethodGet)

	// Get faculty feedbacks for talent for current week and previous week.
	router.HandleFunc("/tenant/{tenantID}/talent-dashboard/faculty-feedback-week-wise/talent/{talentID}/batch/{batchID}",
		controller.GetFacultyFeedbackToTalentWeekWiseDashboardDetails).Methods(http.MethodGet)

	// Get faculty feedbacks for talent.
	router.HandleFunc("/tenant/{tenantID}/talent-dashboard/faculty-feedback/talent/{talentID}/batch/{batchID}",
		controller.GetFacultyFeedbackToTalentDashboardDetails).Methods(http.MethodGet)

	// Get faculty feedbacks leader board for talent.
	router.HandleFunc("/tenant/{tenantID}/talent-dashboard/faculty-feedback-leader-board/batch/{batchID}",
		controller.GetFacultyFeedbackRatingLeaderBoard).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent-weekly-rating",
		controller.GetWeeklyRating).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent-feedback-rating",
		controller.GetFeedbackRating).Methods(http.MethodGet)

	// GetTalentConceptRatingWithBatchTopicAssignment will get talent concept ratings for each batch topic assignemnt for all talents.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent-concept-rating-with-assignment",
		controller.GetTalentConceptRatingWithBatchTopicAssignment).Methods(http.MethodGet)

	log.NewLogger().Info("Talent Dashboard Routes Registered")
}

// GetTalentDashboardDetails gets all dashboard details of talent using service.
func (controller *DashboardController) GetTalentDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetTalentDashboardDetails called")
	talentDashboardDetails := dashboard.TalentDashboard{}
	if err := controller.TalentDashboardService.GetTalentDashboardDetails(&talentDashboardDetails); err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, talentDashboardDetails)
}

// GetFacultyFeedbackToTalentWeekWiseDashboardDetails will return minimum details for faculty feedback to be
// displayed on talent dashboard for current week and previous week.
func (controller *DashboardController) GetFacultyFeedbackToTalentWeekWiseDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetFacultyFeedbackToTalentWeekWiseDashboardDetails called=======================================")

	// Create bucket.
	twoFeedbacks := dashboard.ThisAndPreviousWeekFacultyToTalentFeedackDashboard{}

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting batch id from param and parsing it to uuid.
	batchID, err := util.ParseUUID(params["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Getting talent id from param and parsing it to uuid.
	talentID, err := util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call get interviews method.
	if err := controller.TalentDashboardService.GetFacultyFeedbackToTalentWeekWiseDashboardDetails(&twoFeedbacks, tenantID, batchID, talentID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, twoFeedbacks)
}

// GetFacultyFeedbackToTalentDashboardDetails will return minimum details for faculty feedback to be
// displayed on talent dashboard.
func (controller *DashboardController) GetFacultyFeedbackToTalentDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetFacultyFeedbackToTalentDashboardDetails called=======================================")

	// Create bucket.
	feedbacks := []dashboard.FacultyFeedbackToTalentDashboard{}

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting batch id from param and parsing it to uuid.
	batchID, err := util.ParseUUID(params["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Getting talent id from param and parsing it to uuid.
	talentID, err := util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call get interviews method.
	if err := controller.TalentDashboardService.GetFacultyFeedbackToTalentDashboardDetails(&feedbacks, tenantID, batchID, talentID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, feedbacks)
}

// GetFacultyFeedbackRatingLeaderBoard will return minimum details for faculty feedback to be
// displayed on talent dashboard leader board.
func (controller *DashboardController) GetFacultyFeedbackRatingLeaderBoard(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetFacultyFeedbackRatingLeaderBoard called=======================================")

	// Crete parser for query params.
	parser := web.NewParser(r)

	// Create bucket.
	feedbacks := []dashboard.FacultyFeedbackRatingLeaderBoard{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting batch id from param and parsing it to uuid.
	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Call get interviews method.
	if err := controller.TalentDashboardService.GetFacultyFeedbackRatingLeaderBoard(&feedbacks, tenantID, batchID, parser); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, feedbacks)
}

// GetWeeklyRating will return weekly rating for talents
func (controller *DashboardController) GetWeeklyRating(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info(" ========================= GetWeeklyRating called ========================= ")

	parser := web.NewParser(r)

	var performanceDetails []talent.PerformanceDetails

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	if err := controller.TalentDashboardService.GetWeeklyRating(&performanceDetails, tenantID, batchID, parser); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, performanceDetails)
}

// GetFeedbackRating will fetch feedback rating of all talents
func (controller *DashboardController) GetFeedbackRating(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info(" ========================= GetFeedbackRating called ========================= ")

	parser := web.NewParser(r)

	var performanceDetails []talent.PerformanceDetails

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	if err := controller.TalentDashboardService.GetFeedbackRating(&performanceDetails, tenantID, batchID, parser); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, performanceDetails)
}

// GetTalentConceptRatingWithBatchTopicAssignment will get talent concept ratings for each batch topic assignemnt for all talents.
func (controller *DashboardController) GetTalentConceptRatingWithBatchTopicAssignment(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info(" ========================= GetTalentConceptRatingWithBatchTopicAssignment called ========================= ")

	parser := web.NewParser(r)

	var batchTalents []talent.ConceptRatingWithAssignment

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	if err := controller.TalentDashboardService.GetTalentConceptRatingWithBatchTopicAssignment(&batchTalents, tenantID, batchID, parser); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batchTalents)
}
