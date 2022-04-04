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

// CountryController provides method to update, delete, add, get all for country.
type CountryController struct {
	log            log.Logger
	auth           *security.Authentication
	CountryService *service.CountryService
}

// NewCountryController creates new instance of CountryController.
func NewCountryController(enquiryservice *service.CountryService, log log.Logger, auth *security.Authentication) *CountryController {
	return &CountryController{
		CountryService: enquiryservice,
		log:            log,
		auth:           auth,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *CountryController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all countries list.
	countryList := router.HandleFunc("/tenant/{tenantID}/country",
		controller.GetCountries).Methods(http.MethodGet)

	// Get one country by id.
	router.HandleFunc("/tenant/{tenantID}/country/{countryID}",
		controller.GetCountry).Methods(http.MethodGet)

	// Add one country.
	router.HandleFunc("/tenant/{tenantID}/country",
		controller.AddCountry).Methods(http.MethodPost)

	// Add multiple countries.
	router.HandleFunc("/tenant/{tenantID}/countries",
		controller.AddCountries).Methods(http.MethodPost)

	// Update one country.
	router.HandleFunc("/tenant/{tenantID}/country/{countryID}",
		controller.UpdateCountry).Methods(http.MethodPut)

	// Delete one country.
	router.HandleFunc("/tenant/{tenantID}/country/{countryID}",
		controller.DeleteCountry).Methods(http.MethodDelete)

	controller.log.Info("Country Routes Registered")

	// Temporary.
	router.HandleFunc("/add-id/{tableName}", controller.AddIDToEntity).Methods(http.MethodPut)
	router.HandleFunc("/changeCode", controller.ChangeTalentCode).Methods(http.MethodGet)

	// Exculde routes.
	*exclude = append(*exclude, countryList)
}

// GetCountries returns all countries.
func (controller *CountryController) GetCountries(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetCountries Call********************************")
	parser := web.NewParser(r)
	// Create bucket.
	countries := &[]general.Country{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get all countries service method.
	err = controller.CountryService.GetCountries(countries, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, countries)
}

// GetCountry return specific country by id.
func (controller *CountryController) GetCountry(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetCountry Call********************************")
	parser := web.NewParser(r)
	// Create bucket.
	country := general.Country{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	country.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting country id from param and parsing it to uuid.
	country.ID, err = parser.GetUUID("countryID")
	if err != nil {
		controller.log.Error("unable to parse country id")
		web.RespondError(w, errors.NewHTTPError("unable to parse country id", http.StatusBadRequest))
		return
	}

	// Call get country service method.
	err = controller.CountryService.GetCountry(&country)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, country)
}

// AddCountry add new country.
func (controller *CountryController) AddCountry(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("**************************************AddCountry Call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	country := general.Country{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	country.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	country.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &country)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate country fields.
	err = country.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call add service method.
	err = controller.CountryService.AddCountry(&country)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Country added successfully")
}

// AddCountries adds multiple countries.
func (controller *CountryController) AddCountries(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddCountries Call********************************")
	parser := web.NewParser(r)
	// Create bucket.
	countryIDs := []uuid.UUID{}
	countries := []general.Country{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	credentialID, err := parser.GetUUID("credentialID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &countries)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary country fields.
	for _, country := range countries {
		err = country.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	// Call add service method.
	err = controller.CountryService.AddCountries(&countries, &countryIDs, tenantID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Countries added successfully")
}

// UpdateCountry updates the specified country by id.
func (controller *CountryController) UpdateCountry(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************UpdateCountry Call********************************")
	parser := web.NewParser(r)
	// Create bucket.
	country := general.Country{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	country.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting country id from param and parsing it to uuid.
	country.ID, err = parser.GetUUID("countryID")
	if err != nil {
		controller.log.Error("unable to parse country id")
		web.RespondError(w, errors.NewHTTPError("unable to parse country id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	country.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &country)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate country fields.
	err = country.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.CountryService.UpdateCountry(&country)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Country updated successfully")
}

// DeleteCountry deletes the specified country.
func (controller *CountryController) DeleteCountry(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************DeleteCountry Call********************************")
	parser := web.NewParser(r)
	// Create bucket.
	country := general.Country{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	country.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting country id from param and parsing it to uuid.
	country.ID, err = parser.GetUUID("countryID")
	if err != nil {
		controller.log.Error("unable to parse country id")
		web.RespondError(w, errors.NewHTTPError("unable to parse country id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	country.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Delete country service method call.
	err = controller.CountryService.DeleteCountry(&country)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Country deleted successfullys")
}
