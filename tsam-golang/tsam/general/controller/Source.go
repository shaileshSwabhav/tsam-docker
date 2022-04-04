package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// SourceController provides methods to do CRUD operations.
type SourceController struct {
	log                log.Logger
	SourceService *service.SourceService
	auth               *security.Authentication
}

// NewSourceController creates new instance of SourceController.
func NewSourceController(sourceService *service.SourceService, log log.Logger, auth *security.Authentication) *SourceController {
	return &SourceController{
		SourceService: sourceService,
		log:                log,
		auth:               auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *SourceController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all sources.
	sourceList := router.HandleFunc("/tenant/{tenantID}/source",
		controller.GetAllSources).Methods(http.MethodGet)

	// Add one source.
	router.HandleFunc("/tenant/{tenantID}/source",
		controller.AddSource).Methods(http.MethodPost)

	// Add multiple sources.
	router.HandleFunc("/tenant/{tenantID}/sources",
		controller.AddMultipleSources).Methods(http.MethodPost)

	// Get one source
	router.HandleFunc("/tenant/{tenantID}/source/{sourceID}",
		controller.GetSource).Methods(http.MethodGet)

	// Update source.
	router.HandleFunc("/tenant/{tenantID}/source/{sourceID}",
		controller.UpdateSource).Methods(http.MethodPut)

	// Delete source.
	router.HandleFunc("/tenant/{tenantID}/source/{sourceID}",
		controller.DeleteSource).Methods(http.MethodDelete)

	// Exculde routes.
	*exclude = append(*exclude, sourceList)

	controller.log.Info("Purpose routes registered")
}

// GetAllSources returns all sources.
func (controller *SourceController) GetAllSources(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllSources called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	sources := []general.Source{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get all service method.
	err = controller.SourceService.GetAllSources(tenantID, &sources)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, sources)
}

// DeleteSource deletes the specific source record.
func (controller *SourceController) DeleteSource(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteSource Called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	source := general.Source{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	source.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	source.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting source id from param and parsing it to uuid.
	source.ID, err = parser.GetUUID("sourceID")
	if err != nil {
		controller.log.Error("Invalid source ID")
		web.RespondError(w, errors.NewValidationError("Invalid source ID"))
		return
	}

	// Call delete service method.
	err = controller.SourceService.DeleteSource(&source)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Source deleted successfully")
}

// AddMultipleSources adds multiple source records.
func (controller *SourceController) AddMultipleSources(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddMultipleSources called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	sources := []*general.Source{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse Purpose from request.
	err = web.UnmarshalJSON(r, &sources)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Validate every source entry.
	for _, source := range sources {
		if err := source.Validate(); err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
		source.CreatedBy = credentialID
		source.TenantID = tenantID
	}

	// Call add multiple service method.
	err = controller.SourceService.AddSources(&sources)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// IDCollection will have the list of the UUIDs of the newly added source.
	IDCollection := []uuid.UUID{}
	for _, source := range sources {
		if source != nil {
			IDCollection = append(IDCollection, source.ID)
		}
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Sources added successfully")
}

// GetSource returns a specific source record of a specific tenant.
func (controller *SourceController) GetSource(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetSource Called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	source := general.Source{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	source.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting source id from param and parsing it to uuid.
	source.ID, err = parser.GetUUID("sourceID")
	if err != nil {
		controller.log.Error("Invalid source ID")
		web.RespondError(w, errors.NewValidationError("Invalid source ID"))
		return
	}

	// Call get one service method.
	err = controller.SourceService.GetSource(&source)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, source)
}

// UpdateSource updates the specific source record.
func (controller *SourceController) UpdateSource(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateSource Called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	source := general.Source{}

	// Parse source from request.
	err := web.UnmarshalJSON(r, &source)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Validate source.
	err = source.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	source.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// UpdatedBy will be set with the given credentialID
	source.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting source id from param and parsing it to uuid.
	source.ID, err = parser.GetUUID("sourceID")
	if err != nil {
		controller.log.Error("Invalid source ID")
		web.RespondError(w, errors.NewValidationError("Invalid source ID"))
		return
	}

	// Call update service method.
	err = controller.SourceService.UpdateSource(&source)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Purpose updated successfully")
}

// AddSource validates and calls service to add new source record.
func (controller *SourceController) AddSource(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddSource Called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	source := general.Source{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &source)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Parse and set tenant ID.
	source.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	source.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Validate source.
	err = source.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call add service method.
	err = controller.SourceService.AddSource(&source)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Source added successfully")
}
