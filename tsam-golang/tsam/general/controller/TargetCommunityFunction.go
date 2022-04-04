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

// TargetCommunityFunctionController provide method to update, delete, add, get methods for targetCommunityFunction.
type TargetCommunityFunctionController struct {
	log                            log.Logger
	TargetCommunityFunctionService *service.TargetCommunityFunctionService
	auth                           *security.Authentication
}

// NewTargetCommunityFunctionController create new instance of TargetCommunityFunctionController.
func NewTargetCommunityFunctionController(targetCommunityFunctionService *service.TargetCommunityFunctionService, log log.Logger, auth *security.Authentication) *TargetCommunityFunctionController {
	return &TargetCommunityFunctionController{
		TargetCommunityFunctionService: targetCommunityFunctionService,
		log:                            log,
		auth:                           auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *TargetCommunityFunctionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	
	// Get all targetCommunityFunctions by limit and offset.
	router.HandleFunc("/tenant/{tenantID}/target-community-function",
		controller.GetTargetCommunityFunctions).Methods(http.MethodGet)

	// Get targetCommunityFunction list.
	router.HandleFunc("/tenant/{tenantID}/target-community-function-list",
		controller.GetTargetCommunityFunctionList).Methods(http.MethodGet)

	// Get all targetCommunityFunctions by department id.
	router.HandleFunc("/tenant/{tenantID}/target-community-function/department/{departmentID}",
		controller.GetTargetCommunityFunctionByDepartment).Methods(http.MethodGet)

	// Add one targetCommunityFunction.
	router.HandleFunc("/tenant/{tenantID}/target-community-function",
		controller.AddTargetCommunityFunction).Methods(http.MethodPost)

	// Add targetCommunityFunctions.
	router.HandleFunc("/tenant/{tenantID}/target-community-function",
		controller.AddTargetCommunityFunctions).Methods(http.MethodPost)

	// Get one targetCommunityFunction.
	router.HandleFunc("/tenant/{tenantID}/target-community-function/{targetCommunityFunctionID}",
		controller.GetTargetCommunityFunction).Methods(http.MethodGet)

	// Update one targetCommunityFunction.
	router.HandleFunc("/tenant/{tenantID}/target-community-function/{targetCommunityFunctionID}",
		controller.UpdateTargetCommunityFunction).Methods(http.MethodPut)

	// Delete one targetCommunityFunction.
	router.HandleFunc("/tenant/{tenantID}/target-community-function/{targetCommunityFunctionID}",
		controller.DeleteTargetCommunityFunction).Methods(http.MethodDelete)

	controller.log.Info("TargetCommunityFunction Route Registered")
}

// AddTargetCommunityFunction adds one targetCommunityFunction.
func (controller *TargetCommunityFunctionController) AddTargetCommunityFunction(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddTargetCommunityFunction call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	targetCommunityFunction := general.TargetCommunityFunction{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &targetCommunityFunction)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := targetCommunityFunction.Validate(); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	targetCommunityFunction.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	targetCommunityFunction.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.TargetCommunityFunctionService.AddTargetCommunityFunction(&targetCommunityFunction); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Target community function added successfully")
}

// GetTargetCommunityFunctionList returns all targetCommunityFunctions.
func (controller *TargetCommunityFunctionController) GetTargetCommunityFunctionList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetTargetCommunityFunctionList call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	targetCommunityFunctions := []general.TargetCommunityFunction{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get targetCommunityFunctions method.
	if err := controller.TargetCommunityFunctionService.GetTargetCommunityFunctionList(&targetCommunityFunctions, tenantID); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, targetCommunityFunctions)
}

// GetTargetCommunityFunctions returns all targetCommunityFunctions.
func (controller *TargetCommunityFunctionController) GetTargetCommunityFunctions(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetTargetCommunityFunctions call**************************************")
	parser := web.NewParser(r)
	//create bucket
	targetCommunityFunctions := []general.TargetCommunityFunctionDTO{}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// For pagination.
	var totalCount int

	// Call get targetCommunityFunctions method.
	if err := controller.TargetCommunityFunctionService.GetTargetCommunityFunctions(&targetCommunityFunctions, tenantID, parser, &totalCount); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, targetCommunityFunctions)
}

// GetTargetCommunityFunctionByDepartment returns all targetCommunityFunctions by specific department id.
func (controller *TargetCommunityFunctionController) GetTargetCommunityFunctionByDepartment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetTargetCommunityFunctionByDepartment call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	targetCommunityFunctions := []general.TargetCommunityFunction{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting department id from param and parsing it to uuid.
	departmentID, err := parser.GetUUID("departmentID")
	if err != nil {
		controller.log.Error("unable to parse department id")
		web.RespondError(w, errors.NewHTTPError("unable to parse department id", http.StatusBadRequest))
		return
	}

	// Call get targetCommunityFunction by department method.
	if err = controller.TargetCommunityFunctionService.GetTargetCommunityFunctionByDepartment(&targetCommunityFunctions, tenantID, departmentID); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, targetCommunityFunctions)
}

// AddTargetCommunityFunctions adds multiple targetCommunityFunctions.
func (controller *TargetCommunityFunctionController) AddTargetCommunityFunctions(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddTargetCommunityFunctions call**************************************")
	parser := web.NewParser(r)
	// Create bucket for targetCommunityFunction ids to be added.
	targetCommunityFunctionIDs := []uuid.UUID{}

	// Create bucket for targetCommunityFunctions to be added.
	targetCommunityFunctions := []general.TargetCommunityFunction{}

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
	if err := web.UnmarshalJSON(r, &targetCommunityFunctions); err != nil {
		controller.log.Error("Unable to parse data")
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate all compulsary fields of cities.
	for _, targetCommunityFunction := range targetCommunityFunctions {
		err = targetCommunityFunction.Validate()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}

	}

	// Call add targetCommunityFunctions service method.
	if err := controller.TargetCommunityFunctionService.AddTargetCommunityFunctions(&targetCommunityFunctions, &targetCommunityFunctionIDs, tenantID, credentialID); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Target community fucntions added successfully")
}

// GetTargetCommunityFunction returns one targetCommunityFunction by specific targetCommunityFunction id.
func (controller *TargetCommunityFunctionController) GetTargetCommunityFunction(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetTargetCommunityFunction call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	targetCommunityFunction := general.TargetCommunityFunction{}

	var err error

	// Getting targetCommunityFunction id from param and parsing it to uuid.
	targetCommunityFunction.ID, err = parser.GetUUID("targetCommunityFunctionID")
	if err != nil {
		controller.log.Error("unable to parse target community function id")
		web.RespondError(w, errors.NewHTTPError("unable to parse target community function id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	targetCommunityFunction.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.TargetCommunityFunctionService.GetTargetCommunityFunction(&targetCommunityFunction); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, targetCommunityFunction)
}

// UpdateTargetCommunityFunction updates one targetCommunityFunction by specific targetCommunityFunction id.
func (controller *TargetCommunityFunctionController) UpdateTargetCommunityFunction(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************UpdateTargetCommunityFunction call**************************************")
	parser := web.NewParser(r)
	// Create bucket
	targetCommunityFunction := general.TargetCommunityFunction{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &targetCommunityFunction)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := targetCommunityFunction.Validate(); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting targetCommunityFunction id from param and parsing it to uuid.
	targetCommunityFunction.ID, err = parser.GetUUID("targetCommunityFunctionID")
	if err != nil {
		controller.log.Error("unable to parse target community function id")
		web.RespondError(w, errors.NewHTTPError("unable to parse target community function id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	targetCommunityFunction.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	targetCommunityFunction.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.TargetCommunityFunctionService.UpdateTargetCommunityFunction(&targetCommunityFunction); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Target community function updated successfully")
}

// DeleteTargetCommunityFunction deletes one targetCommunityFunction by specific targetCommunityFunction id.
func (controller *TargetCommunityFunctionController) DeleteTargetCommunityFunction(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************DeleteTargetCommunityFunction call**************************************")
	parser := web.NewParser(r)
	//create bucket.
	targetCommunityFunction := general.TargetCommunityFunction{}

	var err error

	// Getting id from param and parsing it to uuid.
	targetCommunityFunction.ID, err = parser.GetUUID("targetCommunityFunctionID")
	if err != nil {
		controller.log.Error("unable to parse target community function id")
		web.RespondError(w, errors.NewHTTPError("unable to parse target community function id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	targetCommunityFunction.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	targetCommunityFunction.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.TargetCommunityFunctionService.DeleteTargetCommunityFunction(&targetCommunityFunction); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Target community function deleted successfully")
}
