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

// OutcomeController provides methods to do CRUD operations.
type OutcomeController struct {
	log            log.Logger
	auth           *security.Authentication
	OutcomeService *service.OutcomeService
}

// NewOutcomeController creates new instance of outcome controller.
func NewOutcomeController(outcomeService *service.OutcomeService, log log.Logger, auth *security.Authentication) *OutcomeController {
	return &OutcomeController{
		OutcomeService: outcomeService,
		log:            log,
		auth:           auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *OutcomeController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all outcomes.
	router.HandleFunc("/tenant/{tenantID}/outcome",
		controller.GetAllOutcomes).Methods(http.MethodGet)

	// get outcomes by purpose id.
	router.HandleFunc("/tenant/{tenantID}/purpose/{purposeID}/outcome",
		controller.GetAllOutcomesByPurpose).Methods(http.MethodGet)

	// Add one outcome.
	router.HandleFunc("/tenant/{tenantID}/purpose/{purposeID}/outcome",
		controller.AddOutcome).Methods(http.MethodPost)

	// Add multiple outcomes.
	router.HandleFunc("/tenant/{tenantID}/purpose/{purposeID}/outcomes",
		controller.AddOutcomes).Methods(http.MethodPost)

	// Get one outcome.
	router.HandleFunc("/tenant/{tenantID}/purpose/{purposeID}/outcome/{outcomeID}",
		controller.GetOutcome).Methods(http.MethodGet)

	// Update outcome.
	router.HandleFunc("/tenant/{tenantID}/purpose/{purposeID}/outcome/{outcomeID}",
		controller.UpdateOutcome).Methods(http.MethodPut)

	// Delete outcome.
	router.HandleFunc("/tenant/{tenantID}/purpose/{purposeID}/outcome/{outcomeID}",
		controller.DeleteOutcome).Methods(http.MethodDelete)

	controller.log.Info("Outcome routes registered")
}

// GetAllOutcomes returns all outcomes table data for a specific tenant.
func (controller *OutcomeController) GetAllOutcomes(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllOutcomes called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	outcomes := []general.Outcome{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get all service method.
	err = controller.OutcomeService.GetAllOutcomes(tenantID, &outcomes)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, outcomes)
}

// GetAllOutcomesByPurpose returns all outcomes which have the specific purpose type & tenant.
func (controller *OutcomeController) GetAllOutcomesByPurpose(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllOutcomesByPurpose called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	outcomes := []general.Outcome{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting purpose id from param and parsing it to uuid.
	purposeID, err := parser.GetUUID("purposeID")
	if err != nil {
		controller.log.Error("unable to parse purpose id")
		web.RespondError(w, errors.NewHTTPError("unable to parse purpose id", http.StatusBadRequest))
		return
	}

	// Call get by purpose service method.
	err = controller.OutcomeService.GetAllOutcomesByPurpose(tenantID, purposeID, &outcomes)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, outcomes)
}

// DeleteOutcome deletes the specific outcome record.
func (controller *OutcomeController) DeleteOutcome(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteOutcome Called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	outcome := &general.Outcome{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	outcome.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	outcome.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting purpose id from param and parsing it to uuid.
	outcome.PurposeID, err = parser.GetUUID("purposeID")
	if err != nil {
		controller.log.Error("unable to parse purpose id")
		web.RespondError(w, errors.NewHTTPError("unable to parse purpose id", http.StatusBadRequest))
		return
	}

	// Getting outcome id from param and parsing it to uuid.
	outcome.ID, err = parser.GetUUID("outcomeID")
	if err != nil {
		controller.log.Error("Invalid outcome ID")
		web.RespondError(w, errors.NewValidationError("Invalid outcome ID"))
		return
	}

	// Call delete service method.
	err = controller.OutcomeService.DeleteOutcome(outcome)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Outcome deleted successfully")
}

// AddOutcomes adds multiple outcome records.
func (controller *OutcomeController) AddOutcomes(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddMultipleOutcomes called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	outcomes := []*general.Outcome{}

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

	// Getting purpose id from param and parsing it to uuid.
	purposeID, err := parser.GetUUID("purposeID")
	if err != nil {
		controller.log.Error("unable to parse purpose id")
		web.RespondError(w, errors.NewHTTPError("unable to parse purpose id", http.StatusBadRequest))
		return
	}

	// Parse outcome from request.
	err = web.UnmarshalJSON(r, &outcomes)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Validate every outcome entry.
	for _, outcome := range outcomes {
		if err := outcome.Validate(); err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
		outcome.CreatedBy = credentialID
		outcome.TenantID = tenantID
	}

	// Call add multiple service method.
	err = controller.OutcomeService.AddOutcomes(purposeID, &outcomes)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// IDCollection will have the list of the UUIDs of the newly added outcome.
	IDCollection := []uuid.UUID{}
	for _, outcome := range outcomes {
		if outcome != nil {
			IDCollection = append(IDCollection, outcome.ID)
		}
	}

	// Writing response with OK status
	web.RespondJSON(w, http.StatusOK, "Outcomes added successfully")
}

// GetOutcome returns a specific outcome record of a specific tenant.
func (controller *OutcomeController) GetOutcome(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetOutcome Called==============================")
	parser := web.NewParser(r)
	var err error
	outcome := &general.Outcome{}

	outcome.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	outcome.PurposeID, err = parser.GetUUID("purposeID")
	if err != nil {
		controller.log.Error("unable to parse purpose id")
		web.RespondError(w, errors.NewHTTPError("unable to parse purpose id", http.StatusBadRequest))
		return
	}

	outcome.ID, err = parser.GetUUID("outcomeID")
	if err != nil {
		controller.log.Error("Invalid outcome ID")
		web.RespondError(w, errors.NewValidationError("Invalid outcome ID"))
		return
	}

	err = controller.OutcomeService.GetOutcome(outcome)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status
	web.RespondJSON(w, http.StatusOK, outcome)
}

// UpdateOutcome updates the specific outcome record.
func (controller *OutcomeController) UpdateOutcome(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateOutcome Called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	outcome := general.Outcome{}

	// Parse Outcome from request.
	err := web.UnmarshalJSON(r, &outcome)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Validate Outcome.
	err = outcome.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	outcome.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// UpdatedBy will be set with the given credentialID
	outcome.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting purpose id from param and parsing it to uuid.
	outcome.PurposeID, err = parser.GetUUID("purposeID")
	if err != nil {
		controller.log.Error("unable to parse purpose id")
		web.RespondError(w, errors.NewHTTPError("unable to parse purpose id", http.StatusBadRequest))
		return
	}

	// Getting outcome id from param and parsing it to uuid.
	outcome.ID, err = parser.GetUUID("outcomeID")
	if err != nil {
		controller.log.Error("Invalid outcome ID")
		web.RespondError(w, errors.NewValidationError("Invalid outcome ID"))
		return
	}

	// Call update service method.
	err = controller.OutcomeService.UpdateOutcome(&outcome)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Outcome updated successfully")
}

// AddOutcome validates and calls service to add new outcome record.
func (controller *OutcomeController) AddOutcome(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddOutcome Called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	outcome := general.Outcome{}

	// Parse Outcome from request.
	err := web.UnmarshalJSON(r, &outcome)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Parse and set tenant ID.
	outcome.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	outcome.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set purposeID.
	outcome.PurposeID, err = parser.GetUUID("purposeID")
	if err != nil {
		controller.log.Error("unable to parse purpose id")
		web.RespondError(w, errors.NewHTTPError("unable to parse purpose id", http.StatusBadRequest))
		return
	}

	// Validate Outcome.
	err = outcome.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call add service method.
	err = controller.OutcomeService.AddOutcome(&outcome)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Outcome added successfully")
}
