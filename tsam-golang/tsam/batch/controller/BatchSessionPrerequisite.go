package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// PrerequisiteController provide method to Update, Delete, Add, Get Method For Batch session pre-requisite controller.
type PrerequisiteController struct {
	service *service.PrerequisiteService
	log     log.Logger
	auth    *security.Authentication
}

func NewPrerequisiteController(service *service.PrerequisiteService,
	log log.Logger, auth *security.Authentication) *PrerequisiteController {
	return &PrerequisiteController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *PrerequisiteController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	//add
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/pre-requisite",
		controller.addBatchSessionPrerequisite).Methods(http.MethodPost)

	//update
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/pre-requisite/{prerequisiteID}",
		controller.updateBatchSessionPrerequisite).Methods(http.MethodPut)

	//delete
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/batch-session/{batchSessionID}/pre-requisite/{prerequisiteID}",
		controller.deleteBatchSessionPrerequisite).Methods(http.MethodDelete)

	//get
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/pre-requisite",
		controller.getBatchSessionPrerequisite).Methods(http.MethodGet)
	controller.log.Info("Batch Session Pre-requisite Routes Registered")
}

// addBatchSessionPrerequisite will add pre-requisite for batch session.
func (controller *PrerequisiteController) addBatchSessionPrerequisite(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== addBatchSessionPrerequisite call ==============================")

	prerequisite := batch.BatchSessionPrerequisite{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &prerequisite)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.BatchSessionID, err = parser.GetUUID("batchSessionID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = prerequisite.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddBatchSessionPrerequisite(&prerequisite)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Batch session pre-requisite successfully added")
}

// updateBatchSessionPrerequisite will update pre-requisite for batch session.
func (controller *PrerequisiteController) updateBatchSessionPrerequisite(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== updateBatchSessionPrerequisite call ==============================")

	prerequisite := batch.BatchSessionPrerequisite{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &prerequisite)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.BatchSessionID, err = parser.GetUUID("batchSessionID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.ID, err = parser.GetUUID("prerequisiteID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = prerequisite.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	fmt.Println("###%$$^^$%^*&^%", &prerequisite)
	err = controller.service.UpdateBatchSessionPrerequisite(&prerequisite)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Batch session pre-requisite successfully updated")
}

// deleteBatchSessionPrerequisite will delete pre-requisite for batch session.
func (controller *PrerequisiteController) deleteBatchSessionPrerequisite(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== deleteBatchSessionPrerequisite call ==============================")

	prerequisite := batch.BatchSessionPrerequisite{}
	parser := web.NewParser(r)
	var err error

	prerequisite.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.BatchSessionID, err = parser.GetUUID("batchSessionID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.ID, err = parser.GetUUID("prerequisiteID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	prerequisite.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.DeleteBatchSessionPrerequisite(&prerequisite)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Batch session pre-requisite successfully deleted")
}

// getBatchSessionPrerequisite will get all the batch session pre-requisite.
func (controller *PrerequisiteController) getBatchSessionPrerequisite(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getBatchSessionPrerequisite call ==============================")

	prerequisite := []batch.BatchSessionPrerequisiteDTO{}
	parser := web.NewParser(r)
	var totalCount int

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

	err = controller.service.GetBatchSessionPrerequisite(&prerequisite, tenantID, batchID, &totalCount, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, prerequisite)
}
