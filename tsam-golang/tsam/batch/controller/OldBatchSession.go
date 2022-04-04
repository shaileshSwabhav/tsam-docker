package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// OBatchSessionController Provide method to Update, Delete, Add, Get Method For Batch.
type OBatchSessionController struct {
	BatchSessionService *service.OBatchSessionService
}

// OldBatchSessionController Create New Instance Of BatchController.
func OldBatchSessionController(bs *service.OBatchSessionService) *OBatchSessionController {
	return &OBatchSessionController{
		BatchSessionService: bs,
	}
}

func (con *OBatchSessionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/old-batch-session/credential/{credentialID}",
		con.AddSessionForBatch).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/old-batch-sessions/credential/{credentialID}",
		con.UpdateSessionsForBatch).Methods(http.MethodPut)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/old-batch-session/{batchSessionID}/credential/{credentialID}",
		con.UpdateSessionForBatch).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/old-batch-session/{batchSessionID}/credential/{credentialID}",
		con.DeleteSessionForBatch).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/old-session",
		con.GetSessionsForBatch).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/old-module",
		con.GetModuleSessionsForBatch).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/old-session-assignment",
		con.GetSessionsAndAssignmentsForBatch).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/old-session-list", con.GetBatchSessionList).Methods(http.MethodGet)

	log.NewLogger().Info("Batch-Session routes registered")

}

// AddSessionForBatch will add sessions to specified batch
func (con *OBatchSessionController) AddSessionForBatch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddSession call==============================")
	param := mux.Vars(r)

	batchID, err := util.ParseUUID(param["batchID"])
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

	credentialID, err := util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSessions := []bat.MappedSession{}
	err = web.UnmarshalJSON(r, &batchSessions)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, batchSession := range batchSessions {
		err = batchSession.ValidateSession()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = con.BatchSessionService.AddSessionForBatch(&batchSessions, tenantID, batchID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Sessions succesfully added to batch")
}

func (con *OBatchSessionController) UpdateSessionForBatch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateSessionForBatch call==============================")
	param := mux.Vars(r)
	batchSession := bat.MappedSession{}

	err := web.UnmarshalJSON(r, &batchSession)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSession.BatchID, err = util.ParseUUID(param["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSession.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// batchSession.CourseSessionID, err = util.ParseUUID(param["sessionID"])
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	web.RespondError(w, err)
	// 	return
	// }

	batchSession.ID, err = util.ParseUUID(param["batchSessionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSession.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = batchSession.ValidateSession()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.BatchSessionService.UpdateSessionForBatch(&batchSession)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Sessions succesfully updated to batch")
}

func (con *OBatchSessionController) UpdateSessionsForBatch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateSessionsForBatch call==============================")
	param := mux.Vars(r)

	batchID, err := util.ParseUUID(param["batchID"])
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

	credentialID, err := util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSessions := []bat.MappedSession{}
	err = web.UnmarshalJSON(r, &batchSessions)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, batchSession := range batchSessions {
		err = batchSession.ValidateSession()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = con.BatchSessionService.UpdateSessionsForBatch(&batchSessions, tenantID, batchID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Sessions succesfully updated to batch")
}

func (con *OBatchSessionController) DeleteSessionForBatch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteSessionForBatch call==============================")
	param := mux.Vars(r)
	var err error
	batchSession := &bat.MappedSession{}

	batchSession.BatchID, err = util.ParseUUID(param["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSession.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// batchSession.CourseSessionID, err = util.ParseUUID(param["sessionID"])
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	web.RespondError(w, err)
	// 	return
	// }

	batchSession.ID, err = util.ParseUUID(param["batchSessionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchSession.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = con.BatchSessionService.DeleteSessionForBatch(batchSession)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Session successfully deleted from batch")
}

// GetSessionsForBatch will return all the sessions assigned for specified batch
func (con *OBatchSessionController) GetSessionsForBatch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetSessionsForBatch call==============================")

	// param := mux.Vars(r)
	batchSessions := &[]bat.MappedSessionDTO{}

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

	err = con.BatchSessionService.GetSessionsForBatch(batchSessions, batchID, tenantID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batchSessions)
}

// GetModuleSessionsForBatch will return all the module wise sessions for specified batch.
func (con *OBatchSessionController) GetModuleSessionsForBatch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetModuleSessionsForBatch call==============================")

	// batchModules := &[]bat.CourseModuleDTO{}

	// parser := web.NewParser(r)

	// tenantID, err := parser.GetTenantID()
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	web.RespondError(w, err)
	// 	return
	// }

	// batchID, err := parser.GetUUID("batchID")
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	web.RespondError(w, err)
	// 	return
	// }

	// err = con.BatchSessionService.GetBatchModules(tenantID, batchID, batchModules, parser)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	web.RespondError(w, err)
	// 	return
	// }

	// web.RespondJSON(w, http.StatusOK, batchModules)
}

// GetSessionsAndAssignmentsForBatch will return all the sessions and assignments assigned for specified batch.
func (con *OBatchSessionController) GetSessionsAndAssignmentsForBatch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetSessionsAndAssignmentsForBatch call==============================")

	param := mux.Vars(r)
	batchSessions := &[]bat.MappedSessionDTO{}

	batchID, err := util.ParseUUID(param["batchID"])
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

	r.ParseForm()

	err = con.BatchSessionService.GetSessionsAndAssignmentsForBatch(batchSessions, batchID, tenantID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, batchSessions)
}

// GetBatchSessionList will return list of all the sessions assigned for specified batch
func (con *OBatchSessionController) GetBatchSessionList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetBatchSessionList call==============================")
	param := mux.Vars(r)
	batchSessions := &[]list.BatchSession{}
	batchID, err := util.ParseUUID(param["batchID"])
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

	r.ParseForm()

	err = con.BatchSessionService.GetBatchSessionList(batchSessions, batchID, tenantID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, batchSessions)
}
