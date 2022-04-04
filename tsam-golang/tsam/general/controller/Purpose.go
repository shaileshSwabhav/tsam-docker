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

// PurposeController provides methods to do CRUD operations.
type PurposeController struct {
	log            log.Logger
	auth           *security.Authentication
	PurposeService *service.PurposeService
}

// NewPurposeController creates new instance of purpose controller.
func NewPurposeController(purposeService *service.PurposeService, log log.Logger, auth *security.Authentication) *PurposeController {
	return &PurposeController{
		PurposeService: purposeService,
		log:            log,
		auth:           auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *PurposeController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all purposes.
	router.HandleFunc("/tenant/{tenantID}/purpose",
		controller.GetAllPurposes).Methods(http.MethodGet)

	// Get all purposes by type.
	router.HandleFunc("/tenant/{tenantID}/purpose/type/{type}",
		controller.GetAllPurposesByType).Methods(http.MethodGet)

	// Add one purpose.
	router.HandleFunc("/tenant/{tenantID}/purpose",
		controller.AddPurpose).Methods(http.MethodPost)

	// Add multiple purposes.
	router.HandleFunc("/tenant/{tenantID}/purposes",
		controller.AddMultiplePurposes).Methods(http.MethodPost)

	// Get one purpose.
	router.HandleFunc("/tenant/{tenantID}/purpose/{purposeID}",
		controller.GetPurpose).Methods(http.MethodGet)

	// Update one purpose.
	router.HandleFunc("/tenant/{tenantID}/purpose/{purposeID}",
		controller.UpdatePurpose).Methods(http.MethodPut)

	// Delete one purpose.
	router.HandleFunc("/tenant/{tenantID}/purpose/{purposeID}",
		controller.DeletePurpose).Methods(http.MethodDelete)

	controller.log.Info("Purpose routes registered")
}

// GetAllPurposes returns all purposes table data for a specific tenant.
func (controller *PurposeController) GetAllPurposes(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllPurposes called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	purposes := &[]general.Purpose{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get all service method.
	err = controller.PurposeService.GetAllPurposes(tenantID, purposes)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, purposes)
}

// GetAllPurposesByType returns all purposes which have the specific purpose type & tenant.
func (controller *PurposeController) GetAllPurposesByType(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllPurposesByType called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	purposes := []general.Purpose{}

	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Get purpose type form quey params.
	purposeType := params["type"]

	// Call get purposes by type service method.
	err = controller.PurposeService.GetPurposesByType(tenantID, purposeType, &purposes)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, purposes)
}

// DeletePurpose deletes the specific purpose record.
func (controller *PurposeController) DeletePurpose(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeletePurpose Called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	purpose := &general.Purpose{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	purpose.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	purpose.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting purpose id from param and parsing it to uuid.
	purpose.ID, err = parser.GetUUID("purposeID")
	if err != nil {
		controller.log.Error("Invalid purpose ID")
		web.RespondError(w, errors.NewValidationError("Invalid purpose ID"))
		return
	}

	// Call delete purpose service method.
	err = controller.PurposeService.DeletePurpose(purpose)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Purpose deleted successfully")
}

// AddMultiplePurposes adds multiple purpose records.
func (controller *PurposeController) AddMultiplePurposes(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddMultiplePurposes called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	purposes := []*general.Purpose{}

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

	// Parse purpose from request.
	err = web.UnmarshalJSON(r, &purposes)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Validate every purpose entry and give them tenant id and created_by field.
	for _, purpose := range purposes {
		if err := purpose.Validate(); err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
		purpose.CreatedBy = credentialID
		purpose.TenantID = tenantID
	}

	// Call add purposes service method.
	err = controller.PurposeService.AddPurposes(&purposes)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// IDCollection will have the list of the UUIDs of the newly added purposes.
	IDCollection := []uuid.UUID{}
	for _, purpose := range purposes {
		if purpose != nil {
			IDCollection = append(IDCollection, purpose.ID)
		}
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Purposes added successfully")
}

// GetPurpose returns a specific purpose record of a specific tenant.
func (controller *PurposeController) GetPurpose(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetPurpose Called==============================")
	parser := web.NewParser(r)
	var err error
	purpose := &general.Purpose{}

	purpose.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	purpose.ID, err = parser.GetUUID("purposeID")
	if err != nil {
		controller.log.Error("Invalid purpose ID")
		web.RespondError(w, errors.NewValidationError("Invalid purpose ID"))
		return
	}

	err = controller.PurposeService.GetPurpose(purpose)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status
	web.RespondJSON(w, http.StatusOK, purpose)
}

// UpdatePurpose updates the specific purpose record.
func (controller *PurposeController) UpdatePurpose(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdatePurpose Called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	purpose := general.Purpose{}

	// Parse Purpose from request.
	err := web.UnmarshalJSON(r, &purpose)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Validate Purpose.
	err = purpose.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	purpose.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// UpdatedBy will be set with the given credentialID
	purpose.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting purpose id from param and parsing it to uuid.
	purpose.ID, err = parser.GetUUID("purposeID")
	if err != nil {
		controller.log.Error("Invalid purpose ID")
		web.RespondError(w, errors.NewValidationError("Invalid purpose ID"))
		return
	}

	// Call update purpose service method.
	err = controller.PurposeService.UpdatePurpose(&purpose)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Purpose updated successfully")
}

// AddPurpose validates and calls service to add new purpose record.
func (controller *PurposeController) AddPurpose(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddPurpose Called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	purpose := general.Purpose{}

	// Parse Purpose from request.
	err := web.UnmarshalJSON(r, &purpose)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Parse and set tenant ID.
	purpose.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	purpose.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Validate Purpose.
	err = purpose.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Service call.
	err = controller.PurposeService.AddPurpose(&purpose)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Purpose added successfully")
}
