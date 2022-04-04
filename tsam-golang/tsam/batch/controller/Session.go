package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// SessionController provide method to Update, Delete, Add, Get Method For Batch session controller.
type SessionController struct {
	service *service.SessionService
	log     log.Logger
	auth    *security.Authentication
}

// NewSessionController Create New Instance Of BatchSessionController.
func NewSessionController(service *service.SessionService,
	log log.Logger, auth *security.Authentication) *SessionController {
	return &SessionController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *SessionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/generate-session",
		controller.createBatchSessionPlan).Methods(http.MethodPost)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/sub-topic/{subTopicID}/mark-as-complete",
		controller.markSubTopicAsComplete).Methods(http.MethodPut)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/skip-session",
		controller.skipPendingSession).Methods(http.MethodPut)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/sessions",
		controller.updateBatchSession).Methods(http.MethodPut)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/faculty/{facultyID}/sessions",
		controller.deleteSessionPlan).Methods(http.MethodDelete)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session-plan",
		controller.getBatchSessionPlan).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session",
		controller.getBatchSessions).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/all-session-plan",
		controller.getAllBatchSessions).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session-plan-counts",
		controller.getBatchSessionsCounts).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session-topic-name-list",
		controller.getBatchSessionWithTopicNameList).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/module-topics",
		controller.getSessionModuleTopics).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/module/{moduleID}/topics",
		controller.getModuleTopics).Methods(http.MethodGet)

	controller.log.Info("Batch Session Routes Registered")
}

// getModuleTopics will return all Module Topics and from batch session topics.
func (controller *SessionController) getModuleTopics(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getModuleTopics call ==============================")
	parser := web.NewParser(r)
	batchSessionTopics := []batch.ModuleTopics{}

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
	moduleID, err := parser.GetUUID("moduleID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.GetModuleTopics(tenantID, batchID, moduleID, &batchSessionTopics, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batchSessionTopics)

}

// createBatchSessionPlan will create a plan for batch session.
func (controller *SessionController) createBatchSessionPlan(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== createBatchSessionPlan call ==============================")

	parser := web.NewParser(r)

	batchSession := []batch.SessionTopic{}

	err := web.UnmarshalJSON(r, &batchSession)
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

	tenantID, err := parser.GetTenantID()
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

	err = controller.service.NewCreateBatchSessionPlan(tenantID, credentialID, batchID, &batchSession, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, nil)
}

// updateBatchSession will add/update topics to session.
func (controller *SessionController) updateBatchSession(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== updateBatchSession call ==============================")

	parser := web.NewParser(r)
	batchSession := []batch.SessionTopic{}

	err := web.UnmarshalJSON(r, &batchSession)
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

	tenantID, err := parser.GetTenantID()
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

	err = controller.service.NewUpdateBatchSession(tenantID, credentialID, batchID, &batchSession, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, nil)
}

// markSubTopicAsComplete will mark specified session as completed.
func (controller *SessionController) markSubTopicAsComplete(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== markSubTopicAsComplete call ==============================")

	parser := web.NewParser(r)

	batchSession := batch.SessionTopic{}

	err := web.UnmarshalJSON(r, &batchSession)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSession.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSession.SubTopicID, err = parser.GetUUID("subTopicID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSession.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSession.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.MarkSubTopicAsComplete(&batchSession)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, nil)
}

// SkipPreviousSession will skip the pending sessions to next date.
func (controller *SessionController) skipPendingSession(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== skipPreviousSession call ==============================")

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

	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.NewSkipPendingSession(tenantID, credentialID, batchID, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, nil)
}

// deleteSessionPlan will delete all the session plan assigned to specified faculty.
func (controller *SessionController) deleteSessionPlan(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== deleteSessionPlan call ==============================")

	parser := web.NewParser(r)
	session := batch.Session{}
	var err error

	session.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	session.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	session.FacultyID, err = parser.GetUUID("facultyID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	session.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.DeleteSessionPlan(&session)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, nil)
}

// getBatchSessionPlan will get all the sessions for specified date.
func (controller *SessionController) getBatchSessionPlan(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getBatchSessionPlan call ==============================")

	parser := web.NewParser(r)

	// modules := []course.ModuleDTO{}
	batchSession := batch.SessionDTO{}

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

	err = controller.service.NewGetBatchSessionPlan(&batchSession, tenantID, batchID, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batchSession)
}

// getAllBatchSessions will return all batch sessions for specified batch.
func (controller *SessionController) getAllBatchSessions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getAllBatchSessions call ==============================")

	parser := web.NewParser(r)

	batchSessions := []batch.SessionDTO{}

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

	err = controller.service.NewGetAllBatchSessions(&batchSessions, tenantID, batchID, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batchSessions)
}

// getBatchSessionsCounts will return all counts for batch session plan for specified batch.
func (controller *SessionController) getBatchSessionsCounts(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getBatchSessionsCounts call ==============================")

	parser := web.NewParser(r)

	batchSessionCounts := batch.SessionCounts{}

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

	err = controller.service.NewGetBatchSessionsCounts(&batchSessionCounts, tenantID, batchID, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batchSessionCounts)
}

// getBatchSessions will get all the sessions for specified date.
func (controller *SessionController) getBatchSessions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getBatchSessions call ==============================")

	parser := web.NewParser(r)

	batchSessions := []batch.Session{}

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

	err = controller.service.GetBatchSessions(&batchSessions, tenantID, batchID, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batchSessions)
}

// getBatchSessionWithTopicNameList will return list of topic names and sub-topic names for specified batch.
func (controller *SessionController) getBatchSessionWithTopicNameList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getBatchSessionWithTopicNameList call ==============================")

	parser := web.NewParser(r)

	batchSessions := []batch.SessionWithTopicNameDTO{}

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

	err = controller.service.GetBatchSessionWithTopicNameList(tenantID, batchID, &batchSessions, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batchSessions)
}

// getSessionModuleTopics will return topics, subTopics for specified module.
func (controller *SessionController) getSessionModuleTopics(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getBatchModuleTopics call ==============================")

	moduleTopics := []course.ModuleTopicDTO{}
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

	// moduleID, err := parser.GetUUID("moduleID")
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	web.RespondError(w, err)
	// 	return
	// }

	err = controller.service.GetSessionModuleTopics(tenantID, batchID, &moduleTopics, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, moduleTopics)

}
