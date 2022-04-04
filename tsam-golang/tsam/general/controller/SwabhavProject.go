package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// SwabhavProjectController provides methods to do CRUD operations.
type SwabhavProjectController struct {
	log            log.Logger
	auth           *security.Authentication
	ProjectService *service.SwabhavProjectService
}

// NewSwabhavProjectController creates new instance of project controller.
func NewSwabhavProjectController(generalService *service.SwabhavProjectService, log log.Logger, auth *security.Authentication) *SwabhavProjectController {
	return &SwabhavProjectController{
		ProjectService: generalService,
		log:            log,
		auth:           auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *SwabhavProjectController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add project.
	router.HandleFunc("/tenant/{tenantID}/project",
		controller.AddProject).Methods(http.MethodPost)

	// Update project.
	router.HandleFunc("/tenant/{tenantID}/project/{projectID}",
		controller.UpdateProject).Methods(http.MethodPut)

	// Delete project.
	router.HandleFunc("/tenant/{tenantID}/project/{projectID}",
		controller.DeleteProject).Methods(http.MethodDelete)

	// Get all projects with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/project",
		controller.GetAllProjects).Methods(http.MethodGet)

	// Get one project.
	router.HandleFunc("/tenant/{tenantID}/project/{projectID}",
		controller.GetProject).Methods(http.MethodGet)

	// Get project list.
	router.HandleFunc("/tenant/{tenantID}/project-list",
		controller.GetListOfProjects).Methods(http.MethodGet)

	// Get sub-project list.
	router.HandleFunc("/tenant/{tenantID}/project/{projectID}/sub-project-list",
		controller.GetListOfSubProjects).Methods(http.MethodGet)

	// // Get project list (along with all sub projects).
	// router.HandleFunc("/tenant/{tenantID}/projectl",
	// 	controller.GetProjects).Methods(http.MethodGet)

	controller.log.Info("Project Route Registered")
}

// AddProject will add the new dpeartment record in the table.
func (controller *SwabhavProjectController) AddProject(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddProject called==============================")
	parser := web.NewParser(r)
	project := general.SwabhavProject{}

	err := web.UnmarshalJSON(r, &project)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	project.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	project.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	err = project.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProjectService.AddProject(&project)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Project added successfully")
}

// UpdateProject will update the specified project record in the table.
func (controller *SwabhavProjectController) UpdateProject(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateProject called==============================")
	parser := web.NewParser(r)
	project := general.SwabhavProject{}

	err := web.UnmarshalJSON(r, &project)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	project.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	project.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}
	project.ID, err = parser.GetUUID("projectID")
	if err != nil {
		controller.log.Error("unable to parse project id")
		web.RespondError(w, errors.NewHTTPError("unable to parse project id", http.StatusBadRequest))
		return
	}

	err = project.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProjectService.UpdateProject(&project)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Project updated successfully")
}

// DeleteProject will delete the specified department record from the table.
func (controller *SwabhavProjectController) DeleteProject(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteProject called==============================")
	parser := web.NewParser(r)
	project := general.SwabhavProject{}
	var err error

	project.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	project.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}
	project.ID, err = parser.GetUUID("projectID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse project id", http.StatusBadRequest))
		return
	}

	err = controller.ProjectService.DeleteProject(&project)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Project deleted successfully")
}

// GetAllProjects will return all the records from project table.
func (controller *SwabhavProjectController) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllProjects called==============================")
	parser := web.NewParser(r)
	projects := []general.SwabhavProject{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	var totalCount int
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProjectService.GetAllProject(&projects, tenantID, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, projects)
}

// GetProjects will return all the records from project table.
func (controller *SwabhavProjectController) GetProjects(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetProjects called==============================")
	parser := web.NewParser(r)
	projects := []general.SwabhavProject{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	err = controller.ProjectService.GetProjects(&projects, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, projects)
}

// GetProject will return specified record from department table.
func (controller *SwabhavProjectController) GetProject(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetProject called==============================")
	parser := web.NewParser(r)
	department := general.SwabhavProject{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	projectID, err := parser.GetUUID("projectID")
	if err != nil {
		controller.log.Error("unable to parse project id")
		web.RespondError(w, errors.NewHTTPError("unable to parse project id", http.StatusBadRequest))
		return
	}

	err = controller.ProjectService.GetProject(&department, tenantID, projectID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, department)
}

// Niranjan

// GetListOfProjects will return a list of projects.
func (controller *SwabhavProjectController) GetListOfProjects(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetListOfProjects called==============================")
	parser := web.NewParser(r)
	projects := []general.SwabhavProject{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse project id")
		web.RespondError(w, errors.NewHTTPError("unable to parse project id", http.StatusBadRequest))
		return
	}

	err = controller.ProjectService.GetProjectList(tenantID, &projects)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, projects)
}

// GetListOfSubProjects will return a list of sub-projects w.r.t. project.
func (controller *SwabhavProjectController) GetListOfSubProjects(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetListOfSubProjects called==============================")
	parser := web.NewParser(r)
	projects := []general.SwabhavProject{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	projectID, err := parser.GetUUID("projectID")
	if err != nil {
		controller.log.Error("unable to parse project id")
		web.RespondError(w, errors.NewHTTPError("unable to parse project id", http.StatusBadRequest))
		return
	}

	err = controller.ProjectService.GetSubProjectList(tenantID, projectID, &projects)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, projects)
}
