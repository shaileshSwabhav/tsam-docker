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

// PrjectRatingController provides methods to do CRUD operations.
type ProjectRatingController struct {
	ProjectRatingService *service.ProjectRatingService
	log                  log.Logger
	auth                 *security.Authentication
}

// NewPrjectRatingController creates new instance of concept controller.
func NewPrjectRatingController(ProjectRatingService *service.ProjectRatingService,
	log log.Logger, auth *security.Authentication) *ProjectRatingController {
	return &ProjectRatingController{
		ProjectRatingService: ProjectRatingService,
		log:                  log,
		auth:                 auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *ProjectRatingController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all progrject rating parameter.
	router.HandleFunc("/tenant/{tenantID}/project-rating-parameter",
		controller.GetAllProjectRatingParameters).Methods(http.MethodGet)

	log.NewLogger().Info("Project rating parameter Route Registered")
}

// GetAllProgrammingProject will get all the programming-project from the table.
func (controller *ProjectRatingController) GetAllProjectRatingParameters(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllProjectRatingParameters call==============================")

	projectRatingParameters := []batch.ProgrammingProjectRatingParameter{}
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	var totalCount int

	err = controller.ProjectRatingService.GetAllProgrammingProject(&projectRatingParameters, tenantID, &totalCount, parser)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, projectRatingParameters)
}
