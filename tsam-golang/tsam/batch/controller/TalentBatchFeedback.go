package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentFeedbackController provides methods to do CRUD operations.
type TalentFeedbackController struct {
	FeedbackService *service.TalentFeedbackService
}

// NewTalentFeedbackController creates new instance of feedback type controller.
func NewTalentFeedbackController(feedbackService *service.TalentFeedbackService) *TalentFeedbackController {
	return &TalentFeedbackController{
		FeedbackService: feedbackService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (con *TalentFeedbackController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/feedback/credential/{credentialID}",
		con.AddBatchFeedback).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/feedbacks/credential/{credentialID}",
		con.AddBatchFeedbacks).Methods(http.MethodPost)

	// delete
	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/faculty/{facultyID}/talent/feedback/credential/{credentialID}",
	// 	con.DeleteFacultyBatchFeedback).Methods(http.MethodDelete)
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/faculty/{facultyID}/talent/{talentID}/feedback/credential/{credentialID}",
		con.DeleteTalentBatchFeedback).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/feedback",
		con.GetAllTalentBatchFeeback).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}/feedback",
		con.GetTalentBatchFeedback).Methods(http.MethodGet)

	log.NewLogger().Info("Talent Batch Feedback Routes Registered")
}

// AddBatchFeedback will add talent_feedback to the table
func (con *TalentFeedbackController) AddBatchFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddBatchFeedback called==============================")
	param := mux.Vars(r)
	feedback := bat.TalentFeedback{}

	err := web.UnmarshalJSON(r, &feedback)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedback.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedback.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedback.BatchID, err = util.ParseUUID(param["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = feedback.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.FeedbackService.AddBatchFeedback(&feedback)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback successfully submitted")
}

// AddBatchFeedbacks will add feedback for multiple questions
func (con *TalentFeedbackController) AddBatchFeedbacks(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddBatchFeedbacks called==============================")
	param := mux.Vars(r)
	feedbacks := []bat.TalentFeedback{}

	err := web.UnmarshalJSON(r, &feedbacks)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := util.ParseUUID(param["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, feedback := range feedbacks {
		err = feedback.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = con.FeedbackService.AddBatchFeedbacks(&feedbacks, tenantID, batchID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback successfully submitted")
}

// DeleteTalentBatchFeedback will delete specified feedback in table
func (con *TalentFeedbackController) DeleteTalentBatchFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteTalentBatchFeedback called==============================")
	param := mux.Vars(r)
	var err error
	feedback := bat.TalentFeedback{}

	feedback.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedback.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedback.BatchID, err = util.ParseUUID(param["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedback.FacultyID, err = util.ParseUUID(param["facultyID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	feedback.TalentID, err = util.ParseUUID(param["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.FeedbackService.DeleteTalentBatchFeedback(&feedback)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback successfully deleted")
}

// GetAllTalentBatchFeeback will return specified talent's batch feedback in table (admin and faculty login)
func (con *TalentFeedbackController) GetAllTalentBatchFeeback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllTalentBatchFeeback called==============================")
	param := mux.Vars(r)
	feedbacks := []bat.TalentBatchFeedbackDTO{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := util.ParseUUID(param["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse Form
	r.ParseForm()

	err = con.FeedbackService.GetAllTalentBatchFeeback(&feedbacks, tenantID, batchID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, feedbacks)
}

// GetTalentBatchFeedback will return specified talent's batch feedback in table (talent login)
func (con *TalentFeedbackController) GetTalentBatchFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentBatchFeedback called==============================")
	param := mux.Vars(r)
	feedbacks := []bat.TalentBatchFeedbackDTO{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := util.ParseUUID(param["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	talentID, err := util.ParseUUID(param["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.FeedbackService.GetTalentBatchFeedback(&feedbacks, tenantID, batchID, talentID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, feedbacks)
}

// DeleteFacultyBatchFeedback will delete specified feedback in table
// func (con *TalentFeedbackController) DeleteFacultyBatchFeedback(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================DeleteBatchFeedbackForFaculty called==============================")
// 	param := mux.Vars(r)
// 	var err error
// 	feedback := bat.TalentFeedback{}

// 	feedback.DeletedBy, err = util.ParseUUID(param["credentialID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	feedback.TenantID, err = util.ParseUUID(param["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	feedback.BatchID, err = util.ParseUUID(param["batchID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	feedback.FacultyID, err = util.ParseUUID(param["facultyID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = con.FeedbackService.DeleteBatchFeedbackForFaculty(&feedback)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, "Feedback successfully deleted")
// }
