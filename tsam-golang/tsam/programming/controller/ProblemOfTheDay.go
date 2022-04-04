package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProblemOfTheDayController provides methods to do CRUD operations.
type ProblemOfTheDayController struct {
	ProblemOfTheDayService *service.ProblemOfTheDayService
}

// NewProblemOfTheDayController creates new instance of outcome controller.
func NewProblemOfTheDayController(problemOfTheDayService *service.ProblemOfTheDayService) *ProblemOfTheDayController {
	return &ProblemOfTheDayController{
		ProblemOfTheDayService: problemOfTheDayService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *ProblemOfTheDayController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all problems of the day questions.
	router.HandleFunc("/tenant/{tenantID}/problem-of-the-day",
		controller.GetProblemsOfTheDay).Methods(http.MethodPost)

	// Get all problems of the previous days questions.
	router.HandleFunc("/tenant/{tenantID}/problem-of-the-day-previous",
		controller.GetProblemsOfThePreviousDays).Methods(http.MethodPost)

	// Get leader board for problem of the day.
	router.HandleFunc("/tenant/{tenantID}/leader-board-potd/talent/{talentID}",
		controller.GetLeaderBoardPotd).Methods(http.MethodPost)

	// Get leader board for all questions.
	router.HandleFunc("/tenant/{tenantID}/leader-board/talent/{talentID}",
		controller.GetLeaderBoard).Methods(http.MethodPost)

	log.NewLogger().Info("Outcome routes registered")
}

// GetProblemsOfTheDay gets problems of the day by date.
func (controller *ProblemOfTheDayController) GetProblemsOfTheDay(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetProblemsOfTheDay called==============================")

	// Create bucket.
	questions := []programming.QuestionProblemOfTheDayDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get all service method.
	err = controller.ProblemOfTheDayService.GetProblemsOfTheDay(tenantID, &questions)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, questions)
}

// GetProblemsOfThePreviousDays gets problems of the previous days by date.
func (controller *ProblemOfTheDayController) GetProblemsOfThePreviousDays(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetProblemsOfThePreviousDays called==============================")

	// Create bucket.
	questions := []programming.QuestionProblemOfTheDayDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Call get all service method.
	err = controller.ProblemOfTheDayService.GetProblemsOfThePreviousDays(tenantID, &questions, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, questions)
}

// GetLeaderBoardPotd gets performers for problem of the day by ranking and score.
func (controller *ProblemOfTheDayController) GetLeaderBoardPotd(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetLeaderBoardPotd called==============================")

	// Create bucket.
	leaderBoard := general.LeaderBoradDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting talent id from param and parsing it to uuid.
	talentID, err := util.ParseUUID(mux.Vars(r)["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Call get leader board service method.
	err = controller.ProblemOfTheDayService.GetLeaderBoardPotd(tenantID, talentID, &leaderBoard, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, leaderBoard)
}

// GetLeaderBoard gets performers for all questions by ranking and score.
func (controller *ProblemOfTheDayController) GetLeaderBoard(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetLeaderBoard called==============================")

	// Create bucket.
	leaderBoard := general.LeaderBoradDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting talent id from param and parsing it to uuid.
	talentID, err := util.ParseUUID(mux.Vars(r)["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call get leader board service method.
	err = controller.ProblemOfTheDayService.GetLeaderBoard(tenantID, talentID, &leaderBoard)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, leaderBoard)
}
