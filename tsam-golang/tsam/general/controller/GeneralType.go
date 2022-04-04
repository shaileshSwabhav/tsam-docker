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

// GeneralTypeController provides methods to do CRUD operations.
type GeneralTypeController struct {
	log                log.Logger
	auth               *security.Authentication
	GeneralTypeService *service.GeneralTypeService
}

// NewGeneralTypeController creates new instance of general type controller.
func NewGeneralTypeController(generalService *service.GeneralTypeService, log log.Logger, auth *security.Authentication) *GeneralTypeController {
	return &GeneralTypeController{
		GeneralTypeService: generalService,
		log:                log,
		auth:               auth,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *GeneralTypeController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all general types.
	router.HandleFunc("/tenant/{tenantID}/general-type",
		controller.GetAllGeneralTypes).Methods(http.MethodGet)

	// Get general type list by type.
	generalTypeList := router.HandleFunc("/tenant/{tenantID}/general-type/type/{type}",
		controller.GetAllGeneralTypesByType).Methods(http.MethodGet)

	// Add one general type.
	router.HandleFunc("/tenant/{tenantID}/general-type",
		controller.AddGeneralType).Methods(http.MethodPost)

	// Add multiple general type.
	router.HandleFunc("/tenant/{tenantID}/general-types",
		controller.AddMultipleGeneralTypes).Methods(http.MethodPost)

	// Get one general type.
	router.HandleFunc("/tenant/{tenantID}/general-type/{generalTypeID}",
		controller.GetGeneralType).Methods(http.MethodGet)

	// Update general type.
	router.HandleFunc("/tenant/{tenantID}/general-type/{generalTypeID}",
		controller.UpdateGeneralType).Methods(http.MethodPut)

	// Delete general type.
	router.HandleFunc("/tenant/{tenantID}/general-type/{generalTypeID}",
		controller.DeleteGeneralType).Methods(http.MethodDelete)

	// Exculde routes.
	*exclude = append(*exclude, generalTypeList)

	controller.log.Info("General Types routes registered")
}

// GetAllGeneralTypes returns all general_types table data for a specific tenant.
func (controller *GeneralTypeController) GetAllGeneralTypes(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllGeneralTypes called==============================")
	parser := web.NewParser(r)
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	generalTypes := &[]general.CommonType{}
	err = controller.GeneralTypeService.GetAllGeneralTypes(tenantID, generalTypes)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, generalTypes)
}

// GetAllGeneralTypesByType returns all general type which have the specific type name & tenant.
func (controller *GeneralTypeController) GetAllGeneralTypesByType(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllGeneralTypesByType called==============================")

	parser := web.NewParser(r)
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	params := mux.Vars(r)
	typeName := params["type"]
	generalTypes := &[]general.CommonType{}
	err = controller.GeneralTypeService.GetGeneralTypesByType(tenantID, typeName, generalTypes)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, generalTypes)
}

// DeleteGeneralType deletes the specific general type record.
func (controller *GeneralTypeController) DeleteGeneralType(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteGeneralType Called==============================")

	var err error
	parser := web.NewParser(r)
	generalType := &general.CommonType{}

	generalType.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	generalType.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	generalType.ID, err = parser.GetUUID("generalTypeID")
	if err != nil {
		controller.log.Error("Invalid general type ID")
		web.RespondError(w, errors.NewValidationError("Invalid general type ID"))
		return
	}

	err = controller.GeneralTypeService.DeleteGeneralType(generalType)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "General type deleted successfully")
}

// AddMultipleGeneralTypes adds multiple general_type records.
func (controller *GeneralTypeController) AddMultipleGeneralTypes(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddMultipleGeneralTypes called==============================")

	parser := web.NewParser(r)
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	generalTypes := []*general.CommonType{}
	// Parse GeneralType from request.
	err = web.UnmarshalJSON(r, &generalTypes)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate every general type entry.
	for _, generalType := range generalTypes {
		if err := generalType.Validate(); err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
		generalType.CreatedBy = credentialID
		generalType.TenantID = tenantID
	}

	err = controller.GeneralTypeService.AddGeneralTypes(generalTypes)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// IDCollection will have the list of the UUIDs of the newly added general type.
	IDCollection := []uuid.UUID{}
	for _, generalType := range generalTypes {
		if generalType != nil {
			IDCollection = append(IDCollection, generalType.ID)
		}
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "General types added successfully")
}

// GetGeneralType returns a specific general_type record of a specific tenant.
func (controller *GeneralTypeController) GetGeneralType(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetGeneralType Called==============================")

	var err error
	generalType := &general.CommonType{}
	parser := web.NewParser(r)

	generalType.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	generalType.ID, err = parser.GetUUID("generalTypeID")
	if err != nil {
		controller.log.Error("Invalid general type ID")
		web.RespondError(w, errors.NewValidationError("Invalid general type ID"))
		return
	}

	err = controller.GeneralTypeService.GetGeneralType(generalType)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, generalType)
}

// UpdateGeneralType updates the specific general_type record.
func (controller *GeneralTypeController) UpdateGeneralType(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateGeneralType Called==============================")
	generalType := &general.CommonType{}
	parser := web.NewParser(r)

	// Parse GeneralType from request.
	err := web.UnmarshalJSON(r, generalType)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate GeneralType.
	err = generalType.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	generalType.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// UpdatedBy will be set with the given credentialID.
	generalType.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	generalType.ID, err = parser.GetUUID("generalTypeID")
	if err != nil {
		controller.log.Error("Invalid general type ID")
		web.RespondError(w, errors.NewValidationError("Invalid general type ID"))
		return
	}

	err = controller.GeneralTypeService.UpdateGeneralType(generalType)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "General type updated successfully")
}

// AddGeneralType validates and calls service to add new general type record.
func (controller *GeneralTypeController) AddGeneralType(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddGeneralType Called==============================")
	generalType := &general.CommonType{}
	parser := web.NewParser(r)

	// Parse GeneralType from request.
	err := web.UnmarshalJSON(r, generalType)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Parse and set tenant ID.
	generalType.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	generalType.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Validate GeneralType.
	err = generalType.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Service call.
	err = controller.GeneralTypeService.AddGeneralType(generalType)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "General type added successfully")
}
