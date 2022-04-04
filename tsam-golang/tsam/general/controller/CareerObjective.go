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

// CareerObjectiveController provides methods to update, delete, add, get methods for career obejective.
type CareerObjectiveController struct {
	log                    log.Logger
	auth                   *security.Authentication
	CareerObjectiveService *service.CareerObjectiveService
}

// NewCareerObjectiveController Create New Instance Of CareerObjectiveController.
func NewCareerObjectiveController(careerObjectiveService *service.CareerObjectiveService, log log.Logger, auth *security.Authentication) *CareerObjectiveController {
	return &CareerObjectiveController{
		CareerObjectiveService: careerObjectiveService,
		log:                    log,
		auth:                   auth,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *CareerObjectiveController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	
	// Get all career objectives by limit and offset.
	router.HandleFunc("/tenant/{tenantID}/career-objective",
		controller.GetCareerObjectives).Methods(http.MethodGet)

	// Add one career objective.
	router.HandleFunc("/tenant/{tenantID}/career-objective",
		controller.AddCareerObjective).Methods(http.MethodPost)

	// Update one career objective.
	router.HandleFunc("/tenant/{tenantID}/career-objective/{careerObjectiveID}",
		controller.UpdateCareerObjective).Methods(http.MethodPut)

	// Delete one career objective.
	router.HandleFunc("/tenant/{tenantID}/career-objective/{careerObjectiveID}",
		controller.DeleteCareerObjective).Methods(http.MethodDelete)

	controller.log.Info("Career Objective Route Registered")
}

// AddCareerObjective adds one career objective.
func (controller *CareerObjectiveController) AddCareerObjective(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===============================AddCareerObjective called=======================================")
	parser := web.NewParser(r)
	// Create bucket.
	careerObjective := general.CareerObjective{}

	// Get params from api.
	//params := mux.Vars(r)

	// Fill the talent variable with given data.
	err := web.UnmarshalJSON(r, &careerObjective)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := careerObjective.Validate(); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	careerObjective.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of talent.
	careerObjective.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.CareerObjectiveService.AddCareerObjective(&careerObjective); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Career objective added successfully")
}

// GetCareerObjectives returns all career objectives.
func (controller *CareerObjectiveController) GetCareerObjectives(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===============================GetCareerObjectives called=======================================")
	parser := web.NewParser(r)
	// Create bucket.
	careerObjectives := []general.CareerObjective{}

	// Create bucket for total talents count.
	var totalCount int

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get career objectives method.
	err = controller.CareerObjectiveService.GetCareerObjectives(&careerObjectives, tenantID, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, careerObjectives)
}

// UpdateCareerObjective updates one career objective by specific career objective id.
func (controller *CareerObjectiveController) UpdateCareerObjective(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===============================UpdateCareerObjective called=======================================")
	parser := web.NewParser(r)
	// Create bucket.
	careerObjective := general.CareerObjective{}

	// Get params from api.
	//params := mux.Vars(r)

	// Fill the talent variable with given data.
	err := web.UnmarshalJSON(r, &careerObjective)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = careerObjective.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set talent ID to careerObjective.
	careerObjective.ID, err = parser.GetUUID("careerObjectiveID")
	if err != nil {
		controller.log.Error("unable to parse Career Objective id")
		web.RespondError(w, errors.NewHTTPError("unable to parse Career Objective id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to careerObjective.
	careerObjective.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of careerObjective.
	careerObjective.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.CareerObjectiveService.UpdateCareerObjective(&careerObjective); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Career objective updated successfully")
}

// DeleteCareerObjective deletes one career objective by specific career objective id.
func (controller *CareerObjectiveController) DeleteCareerObjective(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===============================DeleteCareerObjective called=======================================")
	parser := web.NewParser(r)
	// Create bucket.
	careerObjective := general.CareerObjective{}

	// Get params from api.
	//params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set talent ID.
	careerObjective.ID, err = parser.GetUUID("careerObjectiveID")
	if err != nil {
		controller.log.Error("unable to parse career objective id")
		web.RespondError(w, errors.NewHTTPError("unable to parse career objective id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	careerObjective.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to talent's DeletedBy field.
	careerObjective.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.CareerObjectiveService.DeleteCareerObjective(&careerObjective); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Career objective deleted successfully")
}
