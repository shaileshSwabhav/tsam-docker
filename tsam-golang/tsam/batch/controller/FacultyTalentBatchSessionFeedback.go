package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// FacultySessionFeedbackController provides methods to do CRUD operations.
type FacultySessionFeedbackController struct {
	SessionFeedbackService *service.FacultySessionFeedbackService
	auth                   *security.Authentication
	log                    log.Logger
}

// FacultySessionFeedbackController creates new instance of Sessionfeedback type controller.
func NewFacultySessionFeedbackController(sessionFeedbackService *service.FacultySessionFeedbackService,
	log log.Logger, auth *security.Authentication) *FacultySessionFeedbackController {
	return &FacultySessionFeedbackController{
		SessionFeedbackService: sessionFeedbackService,
		log:                    log,
		auth:                   auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *FacultySessionFeedbackController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/topic/{batchSessionID}/faculty/feedback",
		controller.AddBatchSessionFeedback).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/topic/{batchSessionID}/faculty/feedbacks",
		controller.AddBatchSessionFeedbacks).Methods(http.MethodPost)

	// delete
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/topic/{batchSessionID}/talent/{talentID}/faculty/feedback",
		controller.DeleteFacultyTalentBatchSessionFeedback).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session/{sessionID}/faculty/feedback",
		controller.GetAllFacultyBatchSessionFeedback).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session/{sessionID}/faculty/{facultyID}/feedback",
		controller.GetFacultyBatchSessionFeeback).Methods(http.MethodGet)

	// GetTalentFeedbackDeatils will return feedback for the specified talent
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session/{sessionID}/feedbacks",
		controller.GetTalentFeedbackDetails).Methods(http.MethodGet)

	log.NewLogger().Info("Faculty Batch Session Feedback Routes Registered")

	// update
	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session/{sessionID}/feedback/{sessionFeedbackID}/credential/{credentialID}",
	// 	con.UpdateSessionFeedback).Methods(http.MethodPut)

	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session/{sessionID}/feedback/{sessionFeedbackID}/credential/{credentialID}",
	// 	con.DeleteSessionFeedback).Methods(http.MethodDelete)

	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session/faculty/feedback",
	// 	con.GetAllFeedbackForBatch).Methods(http.MethodGet)

	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session/{sessionID}/talent/{talentID}/faculty/{facultyID}/feedback",
	// 	con.GetFacultyTalentBatchSessionFeeback).Methods(http.MethodGet)
}

// GetTalentFeedbackDetails will return all the sessions list for specified talent.
func (controller *FacultySessionFeedbackController) GetTalentFeedbackDetails(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Talent Feedback Details Called==============================")
	parser := web.NewParser(r)
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	sessionID, err := parser.GetUUID("sessionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	talentSessionFeedback := &[]bat.FacultyTalentBatchSessionFeedback{}

	err = controller.SessionFeedbackService.GetTalentFeedbackDetails(tenantID, batchID, sessionID, talentSessionFeedback, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, talentSessionFeedback)

}

// AddBatchSessionFeedback adds session feedback of talent in the table
func (controller *FacultySessionFeedbackController) AddBatchSessionFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddBatchSessionFeedback called==============================")

	parser := web.NewParser(r)
	sessionFeedback := bat.FacultyTalentBatchSessionFeedback{}

	err := web.UnmarshalJSON(r, &sessionFeedback)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	sessionFeedback.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
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

	// sessionFeedback.BatchSessionID, err = util.ParseUUID(param["sessionID"])
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

	err = controller.SessionFeedbackService.AddBatchSessionFeedback(&sessionFeedback)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback successfully submitted")
}

// AddBatchSessionFeedbacks will add feedback for multiple questions
func (controller *FacultySessionFeedbackController) AddBatchSessionFeedbacks(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddBatchSessionFeedbacks called==============================")

	parser := web.NewParser(r)
	sessionFeedbacks := []bat.FacultyTalentBatchSessionFeedback{}

	err := web.UnmarshalJSON(r, &sessionFeedbacks)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
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

	err = controller.SessionFeedbackService.AddBatchSessionFeedbacks(&sessionFeedbacks, tenantID, batchID, batchSessionID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback successfully submitted")
}

// DeleteFacultyTalentBatchSessionFeedback will delete specified feedback in table
func (controller *FacultySessionFeedbackController) DeleteFacultyTalentBatchSessionFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteFacultyTalentBatchSessionFeedback called==============================")

	var err error
	parser := web.NewParser(r)
	sessionFeedback := bat.FacultyTalentBatchSessionFeedback{}

	sessionFeedback.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
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

	sessionFeedback.TalentID, err = parser.GetUUID("talentID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.SessionFeedbackService.DeleteFacultyTalentBatchSessionFeedback(&sessionFeedback)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feedback successfully deleted")
}

// GetAllFacultyBatchSessionFeedback will return all the feedback for specified batch-sesion
func (controller *FacultySessionFeedbackController) GetAllFacultyBatchSessionFeedback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllFeedbackForBatch called==============================")

	parser := web.NewParser(r)
	sessionFeedbacks := []bat.FacultyTalentBatchSessionFeedbackDTO{}

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

	err = controller.SessionFeedbackService.GetAllFacultyBatchSessionFeedback(&sessionFeedbacks, tenantID, batchID, batchSessionID, parser.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, sessionFeedbacks)
}

// GetFacultyBatchSessionFeeback will return batch-session feedback for specified faculty
func (controller *FacultySessionFeedbackController) GetFacultyBatchSessionFeeback(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetFacultyBatchSessionFeeback called==============================")

	parser := web.NewParser(r)
	// sessionFeedbacks := []bat.FacultyTalentBatchSessionFeedback{}
	sessionFeedbacks := []bat.FacultyTalentBatchSessionFeedbackDTO{}

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

	facultyID, err := parser.GetUUID("facultyID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.SessionFeedbackService.GetFacultyBatchSessionFeeback(&sessionFeedbacks, tenantID, batchID, batchSessionID, facultyID, parser.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, sessionFeedbacks)
}

// // GetFacultyTalentBatchSessionFeeback will return batch-session feedback from faculty for specified talent
// func (con *FacultySessionFeedbackController) GetFacultyTalentBatchSessionFeeback(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetFacultyTalentBatchSessionFeeback called==============================")
// 	param := mux.Vars(r)
// 	// sessionFeedbacks := []bat.FacultyTalentBatchSessionFeedback{}
// 	sessionFeedbacks := []bat.FacultyTalentBatchSessionFeedbackDTO{}

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

// 	sessionID, err := util.ParseUUID(param["sessionID"])
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

// 	err = con.SessionFeedbackService.GetFacultyTalentBatchSessionFeeback(&sessionFeedbacks, tenantID, batchID,
// 		sessionID, facultyID, talentID)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, sessionFeedbacks)
// }

// // UpdateSessionFeedback will update specified feedback in table
// func (con *FacultySessionFeedbackController) UpdateSessionFeedback(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================UpdateSessionFeedback called==============================")
// 	param := mux.Vars(r)
// 	sessionFeedback := bat.FacultyTalentBatchSessionFeedback{}

// 	err := web.UnmarshalJSON(r, &sessionFeedback)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	sessionFeedback.UpdatedBy, err = util.ParseUUID(param["credentialID"])
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

// 	sessionFeedback.ID, err = util.ParseUUID(param["sessionFeedbackID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = sessionFeedback.Validate()
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = con.SessionFeedbackService.UpdateSessionFeedback(&sessionFeedback)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, "Feedback successfully updated")
// }

// // DeleteSessionFeedback will delete specified feedback in table
// func (con *FacultySessionFeedbackController) DeleteSessionFeedback(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================DeleteSessionFeedback called==============================")
// 	param := mux.Vars(r)
// 	var err error
// 	sessionFeedback := bat.FacultyTalentBatchSessionFeedback{}

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

// 	sessionFeedback.ID, err = util.ParseUUID(param["sessionFeedbackID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = con.SessionFeedbackService.DeleteSessionFeedback(&sessionFeedback)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, "Feedback successfully deleted")
// }

// // GetAllFeedbackForBatch will update specified feedback in table
// func (con *FacultySessionFeedbackController) GetAllFeedbackForBatch(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetAllFeedbackForBatch called==============================")
// 	param := mux.Vars(r)
// 	sessionFeedbacks := []bat.FacultyTalentBatchSessionFeedback{}

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

// 	err = con.SessionFeedbackService.GetAllFeedbackForBatch(&sessionFeedbacks, tenantID, batchID)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, sessionFeedbacks)
// }
