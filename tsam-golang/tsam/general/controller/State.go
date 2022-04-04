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

// StateController Provide method to Update, Delete, Add, Get Method For State.
type StateController struct {
	log          log.Logger
	StateService *service.StateService
	auth         *security.Authentication
}

// NewStateController Create New Instance Of StateController.
func NewStateController(enquiryservice *service.StateService, log log.Logger, auth *security.Authentication) *StateController {
	return &StateController{
		StateService: enquiryservice,
		log:          log,
		auth:         auth,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *StateController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all states.
	router.HandleFunc("/tenant/{tenantID}/state",
		controller.GetStates).Methods(http.MethodGet)

	// Get all states by country id.
	stateList := router.HandleFunc("/tenant/{tenantID}/country/{countryID}/state",
		controller.GetStatesByCountryID).Methods(http.MethodGet)

	// Get one state by id.
	router.HandleFunc("/tenant/{tenantID}/state/{stateID}",
		controller.GetState).Methods(http.MethodGet)

	// Add one state.
	router.HandleFunc("/tenant/{tenantID}/state",
		controller.AddState).Methods(http.MethodPost)

	// Add multiple states.
	router.HandleFunc("/tenant/{tenantID}/states",
		controller.AddStates).Methods(http.MethodPost)

	// Update one state.
	router.HandleFunc("/tenant/{tenantID}/state/{stateID}",
		controller.UpdateState).Methods(http.MethodPut)

	// Delete one state.
	router.HandleFunc("/tenant/{tenantID}/state/{stateID}",
		controller.DeleteState).Methods(http.MethodDelete)

	// Exculde routes.
	*exclude = append(*exclude, stateList)

	log.NewLogger().Info("State Route Registered")
}

// GetStates return all states.
func (controller *StateController) GetStates(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetStates call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	states := []general.State{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewValidationError("unable to parse tenant id"))
		return
	}

	// Call get states method.
	err = controller.StateService.GetStates(&states, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, states)
}

// GetStatesByCountryID returns all states by country ID.
func (controller *StateController) GetStatesByCountryID(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetStatesByCountryID call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	states := []general.State{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewValidationError("unable to parse tenant id"))
		return
	}

	// Getting country id from param and parsing it to uuid.
	countryID, err := parser.GetUUID("countryID")
	if err != nil {
		controller.log.Error("unable to parse Country id")
		web.RespondError(w, errors.NewValidationError("unable to parse Country id"))
		return
	}

	// Call get states by country service method.
	err = controller.StateService.GetStateByCountryID(&states, tenantID, countryID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, states)
}

// GetState returns specific state.
func (controller *StateController) GetState(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetState call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	state := general.State{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	state.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewValidationError("unable to parse tenant id"))
		return
	}

	// Getting state id from param and parsing it to uuid.
	state.ID, err = parser.GetUUID("stateID")
	if err != nil {
		controller.log.Error("unable to parse State id")
		web.RespondError(w, errors.NewValidationError("unable to parse State id"))
		return
	}

	// Call get state service method.
	err = controller.StateService.GetState(&state)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, state)
}

// AddState adds new state.
func (controller *StateController) AddState(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddState call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	state := general.State{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &state)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := state.ValidateState(); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	state.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewValidationError("unable to parse tenant id"))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	state.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add state service method.
	err = controller.StateService.AddState(&state)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "State added successfully")
}

// AddStates add multiple states.
func (controller *StateController) AddStates(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddStates call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	statesIDs := []uuid.UUID{}
	states := []general.State{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewValidationError("unable to parse tenant id"))
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
	err = web.UnmarshalJSON(r, &states)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Validate compusalry fields of all states.
	for _, state := range states {
		err = state.ValidateState()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	// Call add multiple states service method.
	err = controller.StateService.AddStates(&states, &statesIDs, tenantID, credentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "States added successfully")
}

// UpdateState updates the state.
func (controller *StateController) UpdateState(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateState call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	state := general.State{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	state.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewValidationError("unable to parse tenant id"))
		return
	}

	// Getting state id from param and parsing it to uuid.
	state.ID, err = parser.GetUUID("stateID")
	if err != nil {
		controller.log.Error("unable to parse State id")
		web.RespondError(w, errors.NewValidationError("unable to parse State id"))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	state.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &state)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary state fields.
	err = state.ValidateState()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.StateService.UpdateState(&state)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "State updated successfully")
}

// DeleteState deletes state.
func (controller *StateController) DeleteState(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteState call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	state := general.State{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	state.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewValidationError("unable to parse tenant id"))
		return
	}

	// Getting state id from param and parsing it to uuid.
	state.ID, err = parser.GetUUID("stateID")
	if err != nil {
		controller.log.Error("unable to parse State id")
		web.RespondError(w, errors.NewValidationError("unable to parse State id"))
		return
	}

	// Getting state id from param and parsing it to uuid.
	state.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.StateService.DeleteState(&state)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "State deleted successfully")
}
