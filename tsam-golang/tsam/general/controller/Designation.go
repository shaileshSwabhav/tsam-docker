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

// DesignationController provides method to update, delete, add, get method for Designation.
type DesignationController struct {
	log                log.Logger
	DesignationService *service.DesignationService
	auth               *security.Authentication
}

// NewDesignationController creates new instance of DesignationController.
func NewDesignationController(designationservice *service.DesignationService, log log.Logger, auth *security.Authentication) *DesignationController {
	return &DesignationController{
		DesignationService: designationservice,
		log:                log,
		auth:               auth,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *DesignationController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all designations with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/designation",
		controller.GetDesignations).Methods(http.MethodGet)

	// Get all designation list.
	designationList := router.HandleFunc("/tenant/{tenantID}/designation-list",
		controller.GetDesignationList).Methods(http.MethodGet)

	// Get one designation.
	router.HandleFunc("/tenant/{tenantID}/designation/{designationID}",
		controller.GetDesignation).Methods(http.MethodGet)

	// Add one designation.
	router.HandleFunc("/tenant/{tenantID}/designation",
		controller.AddDesignation).Methods(http.MethodPost)

	// Add multiple designations.
	router.HandleFunc("/tenant/{tenantID}/designations",
		controller.AddDesignations).Methods(http.MethodPost)

	// Update designation.
	router.HandleFunc("/tenant/{tenantID}/designation/{designationID}",
		controller.UpdateDesignation).Methods(http.MethodPut)

	// Delete designation.
	router.HandleFunc("/tenant/{tenantID}/designation/{designationID}",
		controller.DeleteDesignation).Methods(http.MethodDelete)

	// Exculde routes.
	*exclude = append(*exclude, designationList)

	controller.log.Info("Designation Route Registered")
}

// GetDesignations returns all designations.
func (controller *DesignationController) GetDesignations(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetDesignations call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	designations := []general.Designation{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	var totalCount int
	// Call get all designations service method.
	err = controller.DesignationService.GetDesignations(&designations, tenantID, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, designations)
}

// GetDesignationList returns designation list.
func (controller *DesignationController) GetDesignationList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetDesignationList call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	designations := []general.Designation{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get designation list method.
	err = controller.DesignationService.GetDesignationList(&designations, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, designations)
}

// GetDesignation return the specifed designation by id.
func (controller *DesignationController) GetDesignation(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetDesignation Call********************************")
	parser := web.NewParser(r)
	designation := general.Designation{}
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	designationID, err := parser.GetUUID("designationID")
	if err != nil {
		controller.log.Error("unable to parse designation id")
		web.RespondError(w, errors.NewHTTPError("unable to parse designation id", http.StatusBadRequest))
		return
	}
	err = controller.DesignationService.GetDesignation(&designation, tenantID, designationID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, designation)
}

// AddDesignation adds new designation.
func (controller *DesignationController) AddDesignation(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddDesignation call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	designation := general.Designation{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &designation)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = designation.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	designation.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	designation.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.

	err = controller.DesignationService.AddDesignation(&designation)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Designation added successfully")
}

// AddDesignations adds multiple designations.
func (controller *DesignationController) AddDesignations(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddDesignations call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	designations := []general.Designation{}
	designationIDs := []uuid.UUID{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &designations)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate designation fields.
	for _, designation := range designations {
		err = designation.Validate()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

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

	// Call add multiple service method.
	err = controller.DesignationService.AddDesignations(&designations, &designationIDs, tenantID, credentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Designations added successfully")
}

// UpdateDesignation updates the specified designation by id.
func (controller *DesignationController) UpdateDesignation(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************UpdateDesignation call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	designation := general.Designation{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &designation)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate designation fields.
	err = designation.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	designation.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting designation id from param and parsing it to uuid.
	designation.ID, err = parser.GetUUID("designationID")
	if err != nil {
		controller.log.Error("unable to parse designation id")
		web.RespondError(w, errors.NewHTTPError("unable to parse designation id", http.StatusBadRequest))
		return
	}

	// Getting cresential id from param and parsing it to uuid.
	designation.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.DesignationService.UpdateDesignation(&designation)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Designation updated successfully")
}

// DeleteDesignation delete specofic designation by id.
func (controller *DesignationController) DeleteDesignation(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************DeleteDesignation call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	designation := general.Designation{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	designation.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting designation id from param and parsing it to uuid.
	designation.ID, err = parser.GetUUID("designationID")
	if err != nil {
		controller.log.Error("unable to parse designation id")
		web.RespondError(w, errors.NewHTTPError("unable to parse designation id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	designation.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.DesignationService.DeleteDesignation(&designation)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Designation deleted successfully")
}
