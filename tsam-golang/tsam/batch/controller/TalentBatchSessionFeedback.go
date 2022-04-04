package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentSessionFeedbackController provides methods to do CRUD operations.
type TalentSessionFeedbackController struct {
	SessionFeedbackService *service.TalentSessionFeedbackService
	auth                   *security.Authentication
	log                    log.Logger
}

// NewTalentSessionFeedbackController creates new instance of Sessionfeedback type controller.
func NewTalentSessionFeedbackController(sessionFeedbackService *service.TalentSessionFeedbackService,
	log log.Logger, auth *security.Authentication) *TalentSessionFeedbackController {
	return &TalentSessionFeedbackController{
		SessionFeedbackService: sessionFeedbackService,
		log:                    log,
		auth:                   auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (con *TalentSessionFeedbackController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add one talent's batch session feedback.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/talent/feedback",
		con.AddBatchSessionFeedback).Methods(http.MethodPost)

	// Add multiple talent's batch session feedback.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/talent/feedbacks",
		con.AddBatchSessionFeedbacks).Methods(http.MethodPost)

	// Get multiple talent's batch session feedback.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/talent/feedback",
		con.GetTalentBatchSessionFeedback).Methods(http.MethodGet)

	// Get one talent's batch feedback.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}/feedback",
		con.GetSpecifiedTalentBatchFeedback).Methods(http.MethodGet)

	// Delete one talent's batch session feedback.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/faculty/talent/{talentID}/feedback",
		con.DeleteTalentBatchSessionFeedback).Methods(http.MethodDelete)

	// Get one talent's one batch session's feedback.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/talent/{talentID}/feedback",
		con.GetSpecifiedTalentBatchSessionFeedback).Methods(http.MethodGet)

	// Get talent's feedback by each batch session's.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session-feedbacks",
		con.GetTalentSessionFeedback).Methods(http.MethodGet)

	// delete
	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session/{sessionID}/faculty/{facultyID}/talent/feedback",
	// 	con.DeleteFacultyBatchSessionFeedback).Methods(http.MethodDelete)

	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/feedback",
	// 	con.GetAllTalentBatchSessionFeedback).Methods(http.MethodGet)

	log.NewLogger().Info("Talent Batch Session Feedback Routes Registered")
}

// GetTalentSessionFeedback will get session feedback from talent to faculty for all batch-session
func (con *TalentSessionFeedbackController) GetTalentSessionFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentSessionFeedback called==============================")

	parser := web.NewParser(r)

	sessionTalentFeedbacks := []bat.TalentSessionFeedbackDTO{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.SessionFeedbackService.GetTalentSessionFeedback(&sessionTalentFeedbacks, tenantID, batchID, parser.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, sessionTalentFeedbacks)

}

// AddBatchSessionFeedback will add single session feedback from talent to faculty for specified batch-session
func (con *TalentSessionFeedbackController) AddBatchSessionFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddBatchSessionFeedback called==============================")
	sessionFeedback := bat.TalentBatchSessionFeedback{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &sessionFeedback)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	sessionFeedback.CreatedBy, err = con.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	sessionFeedback.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	sessionFeedback.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	sessionFeedback.BatchSessionID, err = parser.GetUUID("batchSessionID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = sessionFeedback.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.SessionFeedbackService.AddBatchSessionFeedback(&sessionFeedback)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback successfully submitted")
}

// AddBatchSessionFeedbacks will add feedback by talent for multiple questions
func (con *TalentSessionFeedbackController) AddBatchSessionFeedbacks(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddBatchSessionFeedbacks called==============================")

	parser := web.NewParser(r)
	sessionFeedbacks := []bat.TalentBatchSessionFeedback{}

	err := web.UnmarshalJSON(r, &sessionFeedbacks)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := con.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSessionID, err := parser.GetUUID("batchSessionID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, feedback := range sessionFeedbacks {
		err = feedback.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = con.SessionFeedbackService.AddBatchSessionFeedbacks(&sessionFeedbacks, tenantID, credentialID, batchID, batchSessionID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback successfully submitted")
}

// DeleteTalentBatchSessionFeedback will delete specified feedback in table
func (con *TalentSessionFeedbackController) DeleteTalentBatchSessionFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteTalentBatchSessionFeedback called==============================")

	parser := web.NewParser(r)
	var err error
	sessionFeedback := bat.TalentBatchSessionFeedback{}

	sessionFeedback.DeletedBy, err = con.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	sessionFeedback.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	sessionFeedback.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	sessionFeedback.BatchSessionID, err = parser.GetUUID("batchSessionID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// sessionFeedback.FacultyID, err = util.ParseUUID(param["facultyID"])
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	web.RespondError(w, err)
	// 	return
	// }

	sessionFeedback.TalentID, err = parser.GetUUID("talentID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.SessionFeedbackService.DeleteTalentBatchSessionFeedback(&sessionFeedback)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback successfully deleted")
}

// GetTalentBatchSessionFeedback will return all the feedback for specified batch
func (con *TalentSessionFeedbackController) GetTalentBatchSessionFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentBatchSessionFeedback called==============================")

	// sessionFeedbacks := []bat.FacultyTalentBatchSessionFeedback{}

	parser := web.NewParser(r)
	sessionFeedbacks := []bat.TalentBatchSessionFeedbackDTO{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSessionID, err := parser.GetUUID("batchSessionID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.SessionFeedbackService.GetTalentBatchSessionFeedback(&sessionFeedbacks, tenantID, batchID, batchSessionID, parser.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, sessionFeedbacks)
}

// GetSpecifiedTalentBatchFeedback will return batch feedback for specified talent.
func (con *TalentSessionFeedbackController) GetSpecifiedTalentBatchFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetSpecifiedTalentBatchFeedback called==============================")

	// sessionFeedbacks := []bat.FacultyTalentBatchSessionFeedback{}
	sessionFeedbacks := []bat.SingleTalentBatchFeedbackDTO{}
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.SessionFeedbackService.GetSpecifiedTalentBatchFeedback(&sessionFeedbacks, tenantID, batchID, talentID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, sessionFeedbacks)
}

// GetSpecifiedTalentBatchSessionFeedback will return batch feedback for specified talent got one
// batch session.
func (con *TalentSessionFeedbackController) GetSpecifiedTalentBatchSessionFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetSpecifiedTalentBatchSessionFeedback called==============================")

	// sessionFeedbacks := []bat.FacultyTalentBatchSessionFeedback{}
	sessionFeedbacks := []bat.TalentBatchSessionFeedback{}
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSessionID, err := parser.GetUUID("batchSessionID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.SessionFeedbackService.GetSpecifiedTalentBatchSessionFeedback(&sessionFeedbacks, tenantID, batchID, talentID, batchSessionID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, sessionFeedbacks)
}

// GetAllTalentBatchSessionFeedback will return all the feedbacks for specified batch.
// func (con *TalentSessionFeedbackController) GetAllTalentBatchSessionFeedback(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetAllTalentBatchSessionFeedback called==============================")

// 	parser := web.NewParser(r)
// 	sessionFeedbacks := []bat.BatchTopicDTO{}

// 	tenantID, err := parser.GetTenantID()
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	batchID, err := parser.GetUUID("batchID")
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = con.SessionFeedbackService.GetAllTalentBatchSessionFeedback(tenantID, batchID, &sessionFeedbacks)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, sessionFeedbacks)
// }

// // DeleteFacultyBatchSessionFeedback will delete specified feedback in table
// func (con *TalentSessionFeedbackController) DeleteFacultyBatchSessionFeedback(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================DeleteFacultyBatchSessionFeedback called==============================")
// 	param := mux.Vars(r)
// 	var err error
// 	sessionFeedback := bat.TalentBatchSessionFeedback{}

// 	sessionFeedback.DeletedBy, err = util.ParseUUID(param["credentialID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	sessionFeedback.TenantID, err = util.ParseUUID(param["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	sessionFeedback.BatchID, err = util.ParseUUID(param["batchID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	sessionFeedback.SessionID, err = util.ParseUUID(param["sessionID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	sessionFeedback.FacultyID, err = util.ParseUUID(param["facultyID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = con.SessionFeedbackService.DeleteFacultyBatchSessionFeedback(&sessionFeedback)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, "Feedback successfully deleted")
// }
