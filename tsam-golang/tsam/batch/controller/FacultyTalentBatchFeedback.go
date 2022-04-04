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

// FacultyFeedbackController provides methods to do CRUD operations.
type FacultyFeedbackController struct {
	FeedbackService *service.FacultyFeedbackService
}

// FeedbackController creates new instance of feedback type controller.
func NewFacultyFeedbackController(feedbackService *service.FacultyFeedbackService) *FacultyFeedbackController {
	return &FacultyFeedbackController{
		FeedbackService: feedbackService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (con *FacultyFeedbackController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// take confirmation for the endpoint from sir

	// add
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/faculty/feedback/credential/{credentialID}",
		con.AddBatchFeedback).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/faculty/feedbacks/credential/{credentialID}",
		con.AddBatchFeedbacks).Methods(http.MethodPost)

	// delete
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}/faculty/feedback/credential/{credentialID}",
		con.DeleteFacultyTalentBatchFeedback).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/faculty/{facultyID}/feedback",
		con.GetFacultyBatchFeedback).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/faculty/feedback",
		con.GetAllFacultyBatchFeedback).Methods(http.MethodGet)

	log.NewLogger().Info("Faculty Batch Feedback Routes Registered")

	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/feedback/{feedbackID}/credential/{credentialID}",
	// 	con.UpdateTalentFeedback).Methods(http.MethodPut)

	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}/feedback",
	// 	con.GetAllFeedbackForTalent).Methods(http.MethodGet)

	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/faculty/feedback/{feedbackID}/credential/{credentialID}",
	// 	con.DeleteTalentFeedback).Methods(http.MethodDelete)

	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}/faculty/{facultyID}/feedback",
	// 	con.GetFacultyTalentBatchFeedback).Methods(http.MethodGet)
}

// AddBatchFeedback will add talent_feedback to the table
func (con *FacultyFeedbackController) AddBatchFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddBatchFeedback called==============================")
	param := mux.Vars(r)
	feedback := bat.FacultyTalentFeedback{}

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
func (con *FacultyFeedbackController) AddBatchFeedbacks(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddBatchFeedbacks called==============================")
	param := mux.Vars(r)
	feedbacks := []bat.FacultyTalentFeedback{}

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

// DeleteFacultyTalentBatchFeedback will delete specified feedback in table
func (con *FacultyFeedbackController) DeleteFacultyTalentBatchFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteFacultyTalentBatchFeedback called==============================")
	param := mux.Vars(r)
	var err error
	feedback := bat.FacultyTalentFeedback{}

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

	feedback.TalentID, err = util.ParseUUID(param["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.FeedbackService.DeleteFacultyTalentBatchFeedback(&feedback)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback successfully deleted")
}

// GetAllFacultyBatchFeedback will update specified feedback in table (admin login)
func (con *FacultyFeedbackController) GetAllFacultyBatchFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllFacultyBatchFeedback called==============================")
	param := mux.Vars(r)
	feedbacks := []bat.FacultyTalentBatchFeedbackDTO{}

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

	// Parse form
	r.ParseForm()

	err = con.FeedbackService.GetAllFacultyBatchFeedback(&feedbacks, tenantID, batchID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, feedbacks)
}

// GetFacultyBatchFeedback will update specified feedback in table (faculty login)
func (con *FacultyFeedbackController) GetFacultyBatchFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetFacultyBatchFeedback called==============================")
	param := mux.Vars(r)
	// feedbacks := []bat.FacultyTalentFeedback{}
	feedbacks := []bat.FacultyTalentBatchFeedbackDTO{}
	// feedbacks := make([]bat.FacultyBatchFeedbackDTO, 10)

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

	facultyID, err := util.ParseUUID(param["facultyID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse form
	r.ParseForm()

	err = con.FeedbackService.GetFacultyBatchFeedback(&feedbacks, tenantID, batchID, facultyID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, feedbacks)
}

// // GetFacultyTalentBatchFeedback will return faculty feedback for specified talent
// func (con *FacultyFeedbackController) GetFacultyTalentBatchFeedback(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetFacultyTalentBatchFeedback called==============================")
// 	param := mux.Vars(r)
// 	// feedbacks := []bat.FacultyTalentFeedback{}
// 	feedbacks := []bat.FacultyTalentBatchFeedbackDTO{}
// 	// feedbacks := make([]bat.FacultyBatchFeedbackDTO, 10)

// 	tenantID, err := util.ParseUUID(param["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	batchID, err := util.ParseUUID(param["batchID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	facultyID, err := util.ParseUUID(param["facultyID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	talentID, err := util.ParseUUID(param["talentID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = con.FeedbackService.GetFacultyTalentBatchFeedback(&feedbacks, tenantID, batchID, facultyID, talentID)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, feedbacks)
// }

// // UpdateTalentFeedback will update specified feedback in table
// func (con *FacultyFeedbackController) UpdateTalentFeedback(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================UpdateTalentFeedback called==============================")
// 	param := mux.Vars(r)
// 	feedback := bat.FacultyTalentFeedback{}

// 	err := web.UnmarshalJSON(r, &feedback)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	feedback.UpdatedBy, err = util.ParseUUID(param["credentialID"])
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

// 	feedback.ID, err = util.ParseUUID(param["feedbackID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = feedback.Validate()
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = con.FeedbackService.UpdateBatchFeedback(&feedback)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, "Feedback successfully updated")
// }

// // DeleteTalentFeedback will delete specified feedback in table
// func (con *FacultyFeedbackController) DeleteTalentFeedback(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================DeleteTalentFeedback called==============================")
// 	param := mux.Vars(r)
// 	var err error
// 	feedback := bat.FacultyTalentFeedback{}

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

// 	feedback.ID, err = util.ParseUUID(param["feedbackID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = con.FeedbackService.DeleteBatchFeedback(&feedback)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, "Feedback successfully deleted")
// }

// // GetAllFeedbackForTalent will update specified feedback in table
// func (con *FacultyFeedbackController) GetAllFeedbackForTalent(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetAllFeedbackForTalent called==============================")
// 	param := mux.Vars(r)
// 	feedbacks := []bat.FacultyTalentFeedback{}

// 	tenantID, err := util.ParseUUID(param["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	batchID, err := util.ParseUUID(param["batchID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	talentID, err := util.ParseUUID(param["talentID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = con.FeedbackService.GetAllFeedbackForTalent(&feedbacks, tenantID, batchID, talentID)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, feedbacks)
// }
