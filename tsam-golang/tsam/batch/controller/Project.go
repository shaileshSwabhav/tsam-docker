package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// BatchProjectController provide method to Update, Delete, Add, Get Method For Batch project controller.
type BatchProjectController struct {
	service *service.ProjectService
	log     log.Logger
	auth    *security.Authentication
}

func NewBatchProjectController(service *service.ProjectService,
	log log.Logger, auth *security.Authentication) *BatchProjectController {
	return &BatchProjectController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *BatchProjectController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/project", controller.addBatchProject).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/project/{batchProjectID}",
		controller.updateBatchProject).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/project/{batchProjectID}",
		controller.deleteBatchProject).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/project", controller.getBatchProjects).Methods(http.MethodGet)

	controller.log.Info("Batch Project Routes Registered")
}

// addBatchProject will add programming project to batch.
func (controller *BatchProjectController) addBatchProject(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== addBatchProject call ==============================")

	batchProject := batch.Project{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &batchProject)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchProject.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchProject.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchProject.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = batchProject.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddProject(&batchProject)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, nil)
}

// updateBatchProject will update programming project to batch.
func (controller *BatchProjectController) updateBatchProject(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== updateBatchProject call ==============================")

	batchProject := batch.Project{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &batchProject)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchProject.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchProject.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchProject.ID, err = parser.GetUUID("batchProjectID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchProject.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = batchProject.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdateProject(&batchProject)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Batch project successfully updated")
}

// deleteBatchProject will delete programming project to batch.
func (controller *BatchProjectController) deleteBatchProject(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== deleteBatchProject call ==============================")

	batchProject := batch.Project{}
	parser := web.NewParser(r)
	var err error

	batchProject.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchProject.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchProject.ID, err = parser.GetUUID("batchProjectID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	batchProject.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.DeleteProject(&batchProject)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Batch project successfully deleted")
}

// getBatchProjects will get all the batch projects.
func (controller *BatchProjectController) getBatchProjects(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("============================== getBatchProjects call ==============================")

	batchProjects := []batch.ProjectDTO{}
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

	err = controller.service.GetProjects(&batchProjects, tenantID, batchID, &totalCount, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, batchProjects)
}
