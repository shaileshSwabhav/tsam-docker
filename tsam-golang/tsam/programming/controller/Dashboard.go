package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// ConceptDashboardController provides methods to do CRUD operations.
type ConceptDashboardController struct {
	ConceptDashboardService *service.ConceptDashboardService
	log                       log.Logger
	auth                      *security.Authentication
}

// NewConceptDashboardController creates new instance of concept dashboard controller.
func NewConceptDashboardController(conceptDashboardService *service.ConceptDashboardService,
	log log.Logger, auth *security.Authentication) *ConceptDashboardController {
	return &ConceptDashboardController{
		ConceptDashboardService: conceptDashboardService,
		log:                       log,
		auth:                      auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *ConceptDashboardController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all programming concepts with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/complex-concepts",
		controller.GetComplexConcepts).Methods(http.MethodGet)

	log.NewLogger().Info("Concept Dashboard Route Registered")
}

// GetComplexConcepts all complex concepts.
func (controller *ConceptDashboardController) GetComplexConcepts(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetComplexConcepts called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Create bucket for complex concepts.
	concepts := []programming.ComplexConcept{}

	// Call get all service method.
	err = controller.ConceptDashboardService.GetComplexConcepts(&concepts, tenantID, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, concepts)
}