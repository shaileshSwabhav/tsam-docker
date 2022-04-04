package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// BatchSessionTalentController provides methods to do Update, Delete, Add, Get operations on batch_sessions_talents.
type BatchSessionTalentController struct {
	log     log.Logger
	service *service.BatchSessionTalentService
	auth    *security.Authentication
}

// NewBatchSessionsTalentController creates new instance of BatchSessionsTalentController.
func NewBatchSessionsTalentController(service *service.BatchSessionTalentService,
	log log.Logger, auth *security.Authentication) *BatchSessionTalentController {
	return &BatchSessionTalentController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *BatchSessionTalentController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add batch session talent.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/talent/{talentID}/attendance",
		controller.AddTalentAttendance).Methods(http.MethodPost)

	// Add batch session talent.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/talent-attendance",
		controller.AddTalentsAttendance).Methods(http.MethodPost)

	// Get all batch session talents.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/talent",
		controller.GetBatchSessionTalents).Methods(http.MethodGet)

	// Get all batch session talents for specific talent.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/talent/{talentID}",
		controller.GetAllBatchSessionTalentsForTalent).Methods(http.MethodGet)

	// Get one batch session talents for specific talent.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/talent/{talentID}",
		controller.GetOneBatchSessionTalentsForTalent).Methods(http.MethodGet)

	// GetTalentTopicDetails will return all the sessions list for specified talent.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}/talent-details",
		controller.GetTalentTopicDetails).Methods(http.MethodGet)

	// GetAverageRatingForTalent will get average rating for batch.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}/talent-average-rating",
		controller.GetAverageRatingForTalentByWeek).Methods(http.MethodGet)

	// update -> commented because attendance should not be updated
	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/topic/{batchSessionID}/talent/{talentID}/attendance/{talentDetailID}",
	// 	controller.UpdateTalentAttendance).Methods(http.MethodPut)

	controller.log.Info("Batch Session Talent Routes Registered")
}

// GetTalentTopicDetails will return all the sessions list for specified talent.
func (controller *BatchSessionTalentController) GetTalentTopicDetails(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Talent Topic Details Called==============================")

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

	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}
	talentSessionDetails := &[]batch.TalentDetailsDTO{}

	err = controller.service.GetTalentTopicDetails(tenantID, batchID, talentID, talentSessionDetails, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, talentSessionDetails)

}

// AddTalentAttendance will add talent's attendance details to the table.
func (controller *BatchSessionTalentController) AddTalentAttendance(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add Talent Attendance Called==============================")
	talentAttendance := batch.BatchSessionTalent{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &talentAttendance)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	talentAttendance.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch ID.
	talentAttendance.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Parse and set batch session ID.
	// talentAttendance.BatchSessionID, err = parser.GetUUID("batchSessionID")
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, errors.NewHTTPError("unable to parse batch session id", http.StatusBadRequest))
	// 	return
	// }
	talentAttendance.BatchSessionID, err = parser.GetUUID("batchSessionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch session id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentAttendance.TalentID, err = parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	talentAttendance.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = talentAttendance.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddTalentAttendance(&talentAttendance)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Talent attendance successfully added")
}

// AddTalentsAttendance will add multiple talent's attendance details to the table.
func (controller *BatchSessionTalentController) AddTalentsAttendance(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add Talent Attendance Called==============================")
	talentAttendance := []batch.BatchSessionTalent{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &talentAttendance)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch ID.
	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	batchSessionID, err := parser.GetUUID("batchSessionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch session id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	// talentAttendance.TalentID, err = parser.GetUUID("talentID")
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
	// 	return
	// }

	// Parse and set credentialID in CreatedBy field.
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for i := 0; i < len(talentAttendance); i++ {

		err = talentAttendance[i].Validate()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = controller.service.AddTalentsAttendance(tenantID, credentialID, batchSessionID, batchID, &talentAttendance)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Talents attendance successfully added")
}

// UpdateTalentAttendance will update attendance for specified talent.
func (controller *BatchSessionTalentController) UpdateTalentAttendance(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Update Talent Attendance Called==============================")
	talentAttendance := batch.BatchSessionTalent{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &talentAttendance)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	talentAttendance.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch ID.
	talentAttendance.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// talentAttendance.BatchSessionID, err = parser.GetUUID("batchSessionID")
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, errors.NewHTTPError("unable to parse batch session id", http.StatusBadRequest))
	// 	return
	// }
	talentAttendance.BatchSessionID, err = parser.GetUUID("batchSessionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch session id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentAttendance.TalentID, err = parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set batch session talent ID.
	talentAttendance.ID, err = parser.GetUUID("talentDetailID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent attendance id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	talentAttendance.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = talentAttendance.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdateTalentAttendance(&talentAttendance)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Talent attendance successfully updated")
}

// GetBatchSessionTalents will get all the batch-sessions talents.
func (controller *BatchSessionTalentController) GetBatchSessionTalents(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetBatchSessionTalents Called==============================")
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

	batchSessionID, err := parser.GetUUID("batchSessionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch session id", http.StatusBadRequest))
		return
	}

	sessionTalents := &[]batch.BatchSessionTalentDTO{}
	err = controller.service.GetBatchSessionTalents(tenantID, batchID, batchSessionID, sessionTalents, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, sessionTalents)
}

// GetAllBatchSessionTalentsForTalent will get all the batch sessions talents for specific talent.
func (controller *BatchSessionTalentController) GetAllBatchSessionTalentsForTalent(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllBatchSessionForTalent Called==============================")
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

	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	sessionTalents := &[]batch.BatchSessionTalent{}
	err = controller.service.GetAllBatchSessionForTalent(tenantID, batchID, talentID, sessionTalents, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, sessionTalents)
}

// GetAverageRatingForTalentByWeek will get average rating for batch for current week.
func (controller *BatchSessionTalentController) GetAverageRatingForTalentByWeek(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAverageRatingForTalentByWeek Called==============================")
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

	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	averageRatingBatch := batch.AverageRatingBatch{}
	err = controller.service.GetAverageRatingForTalentByWeek(tenantID, batchID, talentID, &averageRatingBatch, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, averageRatingBatch)
}

// GetOneBatchSessionTalentsForTalent will get one the batch sessions talents for specific talent.
func (controller *BatchSessionTalentController) GetOneBatchSessionTalentsForTalent(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetOneBatchSessionTalentsForTalent Called==============================")
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

	batchSessionID, err := parser.GetUUID("batchSessionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch session id", http.StatusBadRequest))
		return
	}

	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	sessionTalents := &batch.BatchSessionTalent{}
	err = controller.service.GetOneBatchSessionTalentsForTalent(tenantID, batchID, talentID, batchSessionID, sessionTalents, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, sessionTalents)
}
