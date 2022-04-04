package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// SpecializationController provide method to update, delete, add, get methods for specialization.
type SpecializationController struct {
	log                   log.Logger
	SpecializationService *service.SpecializationService
	auth                  *security.Authentication
}

// NewSpecializationController create new instance of SpecializationController.
func NewSpecializationController(specializationService *service.SpecializationService, log log.Logger, auth *security.Authentication) *SpecializationController {
	return &SpecializationController{
		SpecializationService: specializationService,
		log:                   log,
		auth:                  auth,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *SpecializationController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all specializations by limit and offset.
	router.HandleFunc("/tenant/{tenantID}/specialization",
		controller.GetSpecializations).Methods(http.MethodGet)

	// Get specialization list.
	specializationList := router.HandleFunc("/tenant/{tenantID}/specialization-list",
		controller.GetSpecializationList).Methods(http.MethodGet)

	// Get all specializations by degree id.
	degreeSpecializationList := router.HandleFunc("/tenant/{tenantID}/specialization/degree/{degreeID}",
		controller.GetSpecializationsByDegree).Methods(http.MethodGet)

	// Add one specialization.
	router.HandleFunc("/tenant/{tenantID}/specialization",
		controller.AddSpecialization).Methods(http.MethodPost)

	// Add specializations.
	router.HandleFunc("/tenant/{tenantID}/specializations",
		controller.AddSpecializations).Methods(http.MethodPost)

	// Get one specialization.
	router.HandleFunc("/tenant/{tenantID}/specialization/{specializationID}",
		controller.GetSpecialization).Methods(http.MethodGet)

	// Update one specialization.
	router.HandleFunc("/tenant/{tenantID}/specialization/{specializationID}",
		controller.UpdateSpecialization).Methods(http.MethodPut)

	// Delete one specialization.
	router.HandleFunc("/tenant/{tenantID}/specialization/{specializationID}",
		controller.DeleteSpecialization).Methods(http.MethodDelete)

	// Exculde routes.
	*exclude = append(*exclude, specializationList, degreeSpecializationList)

	controller.log.Info("Specialization Route Registered")
}

// AddSpecialization adds one specialization.
func (controller *SpecializationController) AddSpecialization(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddSpecialization call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	specialization := general.Specialization{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &specialization)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := specialization.Validate(); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	specialization.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	specialization.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.SpecializationService.AddSpecialization(&specialization); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Specialization added successfully")
}

// GetSpecializationList returns all specializations.
func (controller *SpecializationController) GetSpecializationList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetSpecializationList call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	specializations := []list.Specialization{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get specializations method.
	if err := controller.SpecializationService.GetSpecializationList(&specializations, tenantID); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, specializations)
}

// GetSpecializations returns all specializations.
func (controller *SpecializationController) GetSpecializations(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetSpecializations call**************************************")
	parser := web.NewParser(r)
	//create bucket
	specializations := []general.SpecializationDTO{}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// For pagination.
	var totalCount int

	// Call get specializations method.
	if err := controller.SpecializationService.GetSpecializations(&specializations, tenantID, parser, &totalCount); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, specializations)
}

// GetSpecializationsByDegree returns all specializations by specific degree id.
func (controller *SpecializationController) GetSpecializationsByDegree(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetSpecializationsByDegree call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	specializations := []general.Specialization{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting degree id from param and parsing it to uuid.
	degreeID, err := parser.GetUUID("degreeID")
	if err != nil {
		controller.log.Error("unable to parse degree id")
		web.RespondError(w, errors.NewHTTPError("unable to parse degree id", http.StatusBadRequest))
		return
	}

	// Call get specialization by degree method.
	if err = controller.SpecializationService.GetSpecializationsByDegree(&specializations, tenantID, degreeID); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, specializations)
}

// AddSpecializations adds multiple specializations.
func (controller *SpecializationController) AddSpecializations(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddSpecializations call**************************************")
	parser := web.NewParser(r)
	// Create bucket for specialization ids to be added.
	specializationIDs := []uuid.UUID{}

	// Create bucket for specializations to be added.
	specializations := []general.Specialization{}

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

	// Unmarshal json.
	if err := web.UnmarshalJSON(r, &specializations); err != nil {
		controller.log.Error("Unable to parse data")
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate all compulsary fields of cities.
	for _, specialization := range specializations {
		err = specialization.Validate()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}

	}

	// Call add specializations service method.
	if err := controller.SpecializationService.AddSpecializations(&specializations, &specializationIDs, tenantID, credentialID); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Specializations added successfully")
}

// GetSpecialization returns one specialization by specific specialization id.
func (controller *SpecializationController) GetSpecialization(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetSpecialization call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	specialization := general.Specialization{}

	var err error

	// Getting specialization id from param and parsing it to uuid.
	specialization.ID, err = parser.GetUUID("specializationID")
	if err != nil {
		controller.log.Error("unable to parse specialization id")
		web.RespondError(w, errors.NewHTTPError("unable to parse specialization id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	specialization.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.SpecializationService.GetSpecialization(&specialization); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, specialization)
}

// UpdateSpecialization updates one specialization by specific specialization id.
func (controller *SpecializationController) UpdateSpecialization(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************UpdateSpecialization call**************************************")
	parser := web.NewParser(r)
	// Create bucket
	specialization := general.Specialization{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &specialization)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := specialization.Validate(); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting specialization id from param and parsing it to uuid.
	specialization.ID, err = parser.GetUUID("specializationID")
	if err != nil {
		controller.log.Error("unable to parse specialization id")
		web.RespondError(w, errors.NewHTTPError("unable to parse specialization id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	specialization.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	specialization.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.SpecializationService.UpdateSpecialization(&specialization); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Specialization updated successfully")
}

// DeleteSpecialization deletes one specialization by specific specialization id.
func (controller *SpecializationController) DeleteSpecialization(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************DeleteSpecialization call**************************************")
	parser := web.NewParser(r)
	//create bucket.
	specialization := general.Specialization{}

	var err error

	// Getting id from param and parsing it to uuid.
	specialization.ID, err = parser.GetUUID("specializationID")
	if err != nil {
		controller.log.Error("unable to parse specialization id")
		web.RespondError(w, errors.NewHTTPError("unable to parse specialization id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	specialization.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	specialization.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.SpecializationService.DeleteSpecialization(&specialization); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Specialization deleted successfully")
}
