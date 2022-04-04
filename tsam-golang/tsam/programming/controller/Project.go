package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingProjectController Provide method to Update, Delete, Add, Get Method For programming-project.
type ProgrammingProjectController struct {
	ProgrammingProjectService *service.ProgrammingProjectService
	log                       log.Logger
	auth                      *security.Authentication
}

// NewProgrammingProjectController Create New Instance Of ProgrammingProjectController.
func NewProgrammingProjectController(programmingProjectService *service.ProgrammingProjectService,
	log log.Logger, auth *security.Authentication) *ProgrammingProjectController {
	return &ProgrammingProjectController{
		ProgrammingProjectService: programmingProjectService,
		log:                       log,
		auth:                      auth,
	}
}

// RegisterRoutes Register All Endpoints To Router excluding a few endpoints from token check.
func (controller *ProgrammingProjectController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/programming-project",
		controller.AddProgrammingProject).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/programming-project/{programmingProjectID}",
		controller.UpdateProgrammingProject).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/programming-project/{programmingProjectID}",
		controller.DeleteProgrammingProject).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/programming-project",
		controller.GetAllProgrammingProject).Methods(http.MethodGet)

	log.NewLogger().Info("Programming Project Route Registered")
}

// AddProgrammingProject will add new batch-project in the table.
func (controller *ProgrammingProjectController) AddProgrammingProject(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddProgrammingProject call==============================")

	project := programming.ProgrammingProject{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &project)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	project.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	project.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = project.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProgrammingProjectService.AddProgrammingProject(&project)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, project.ID)
}

// UpdateProgrammingProject will update specified programming-project in the table.
func (controller *ProgrammingProjectController) UpdateProgrammingProject(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateProgrammingProject call==============================")

	project := programming.ProgrammingProject{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &project)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	project.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	project.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	project.ID, err = parser.GetUUID("programmingProjectID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = project.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProgrammingProjectService.UpdateProgrammingProject(&project)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Programming project successfully udpated")
}

// DeleteProgrammingProject will soft delete specified programming-project from the table.
func (controller *ProgrammingProjectController) DeleteProgrammingProject(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteProgrammingProject call==============================")

	project := programming.ProgrammingProject{}
	parser := web.NewParser(r)
	var err error

	project.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	project.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	project.ID, err = parser.GetUUID("programmingProjectID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProgrammingProjectService.DeleteProgrammingProject(&project)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Programming project successfully deleted")
}

// GetAllProgrammingProject will get all the programming-project from the table.
func (controller *ProgrammingProjectController) GetAllProgrammingProject(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllProgrammingProject call==============================")

	projects := []programming.ProgrammingProjectDTO{}
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	var totalCount int

	err = controller.ProgrammingProjectService.GetAllProgrammingProject(&projects, tenantID, &totalCount, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, projects)
}
