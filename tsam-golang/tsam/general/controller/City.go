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

// CityController provides method to update, delete, add, get method for city.
type CityController struct {
	log         log.Logger
	auth        *security.Authentication
	CityService *service.CityService
}

// NewCityController creates new instance of CityController.
func NewCityController(cityService *service.CityService, log log.Logger, auth *security.Authentication) *CityController {
	return &CityController{
		CityService: cityService,
		log:         log,
		auth:        auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *CityController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all cities.
	router.HandleFunc("/tenant/{tenantID}/city",
		controller.GetCities).Methods(http.MethodGet)

	// Get all cities by state id.
	router.HandleFunc("/tenant/{tenantID}/city/state/{stateID}",
		controller.GetCitiesByStateID).Methods(http.MethodGet)

	// Get one city by city id.
	router.HandleFunc("/tenant/{tenantID}/city/{cityID}",
		controller.GetCity).Methods(http.MethodGet)

	// Add one city.
	router.HandleFunc("/tenant/{tenantID}/city",
		controller.AddCity).Methods(http.MethodPost)

	// Add multiple cities.
	router.HandleFunc("/tenant/{tenantID}/cities",
		controller.AddCities).Methods(http.MethodPost)

	// Update one city.
	router.HandleFunc("/tenant/{tenantID}/city/{cityID}",
		controller.UpdateCity).Methods(http.MethodPut)

	// Delete one city.
	router.HandleFunc("/tenant/{tenantID}/city/{cityID}",
		controller.DeleteCity).Methods(http.MethodDelete)

	controller.log.Info("City Routes Registered")
}

// GetCities returns all cities.
func (controller *CityController) GetCities(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetCities call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	cities := []general.City{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get cities method.
	err = controller.CityService.GetCities(&cities, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResonseWriter,
	web.RespondJSON(w, http.StatusOK, cities)
}

// GetCitiesByStateID returns all cities by state id.
func (controller *CityController) GetCitiesByStateID(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetCitiesByStateID call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	cities := []general.City{}

	//param := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting state id from param and parsing it to uuid.
	stateID, err := parser.GetUUID("stateID")
	if err != nil {
		controller.log.Error("unable to parse state id")
		web.RespondError(w, errors.NewHTTPError("unable to parse state id", http.StatusBadRequest))
		return
	}

	// Call get cities by state id method.
	err = controller.CityService.GetCitiesByStateID(&cities, tenantID, stateID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResonseWriter.
	web.RespondJSON(w, http.StatusOK, cities)
}

// GetCity returns specific city by id.
func (controller *CityController) GetCity(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetCity call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	city := general.City{}

	//param := mux.Vars(r)

	var err error

	// Getting tenant id from param and parsing it to uuid.
	city.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting city id from param and parsing it to uuid.
	city.ID, err = parser.GetUUID("cityID")
	if err != nil {
		controller.log.Error("unable to parse city id")
		web.RespondError(w, errors.NewHTTPError("unable to parse city id", http.StatusBadRequest))
		return
	}

	// Call get city method.
	err = controller.CityService.GetCity(&city)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ReponseWriter.
	web.RespondJSON(w, http.StatusOK, city)
}

// AddCity adds new city.
func (controller *CityController) AddCity(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddCity call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	city := general.City{}

	//param := mux.Vars(r)

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &city)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := city.Validate(); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	city.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	city.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add city method.
	err = controller.CityService.AddCity(&city)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResposeWriter.
	web.RespondJSON(w, http.StatusOK, "City added successfully")
}

// AddCities adds new cities.
func (controller *CityController) AddCities(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddCities call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	cityIDs := []uuid.UUID{}
	cities := []general.City{}

	//param := mux.Vars(r)

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
	err = web.UnmarshalJSON(r, &cities)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate all compulsary fields of cities.
	for _, city := range cities {
		err = city.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}

	}

	// Call add cities method.
	err = controller.CityService.AddCities(&cities, &cityIDs, tenantID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponeWriter.
	web.RespondJSON(w, http.StatusOK, "Cities added successfully")
}

// UpdateCity updates city.
func (controller *CityController) UpdateCity(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************UpdateCity call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	city := general.City{}

	//param := mux.Vars(r)

	var err error

	// Getting tenant id from param and parsing it to uuid.
	city.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	city.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting city id from param and parsing it to uuid.
	city.ID, err = parser.GetUUID("cityID")
	if err != nil {
		controller.log.Error("unable to parse city id")
		web.RespondError(w, errors.NewHTTPError("unable to parse city id", http.StatusBadRequest))
		return
	}

	// Unmarshal JSON.
	err = web.UnmarshalJSON(r, &city)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate city.
	err = city.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update city method.
	err = controller.CityService.UpdateCity(&city)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "City updated successfully")
}

// DeleteCity deletes city.
func (controller *CityController) DeleteCity(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************DeleteCity call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	city := general.City{}

	//param := mux.Vars(r)

	var err error

	// Getting tenant id from param and parsing it to uuid.
	city.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	city.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting city id from param and parsing it to uuid.
	city.ID, err = parser.GetUUID("cityID")
	if err != nil {
		controller.log.Error("unable to parse city id")
		web.RespondError(w, errors.NewHTTPError("unable to parse city id", http.StatusBadRequest))
		return
	}

	// Call delete city method.
	err = controller.CityService.DeleteCity(&city)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "City deleted successfully")
}
