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

// UniversityController provide method to update, delete, add, get method for University.
type UniversityController struct {
	log               log.Logger
	auth              *security.Authentication
	UniversityService *service.UniversityService
}

// NewUniversityController creates new Iistance of UniversityController.
func NewUniversityController(service *service.UniversityService, log log.Logger, auth *security.Authentication) *UniversityController {
	return &UniversityController{
		UniversityService: service,
		log:               log,
		auth:              auth,
	}
}

// RegisterRoutes register all endpoint to Router.
func (controller *UniversityController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get university list.
	router.HandleFunc("/tenant/{tenantID}/university-list",
		controller.GetUniversityList).Methods(http.MethodGet)

	// Get all universitiesw with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/university",
		controller.GetUniversities).Methods(http.MethodGet)

	// Get university list by country.
	universityList := router.HandleFunc("/tenant/{tenantID}/country/{countryID}/university",
		controller.GetUniversitiesByCountryList).Methods(http.MethodGet)

	// Add university.
	router.HandleFunc("/tenant/{tenantID}/university",
		controller.AddUniversity).Methods(http.MethodPost)

	// Add multiple universities.
	router.HandleFunc("/tenant/{tenantID}/universities",
		controller.AddUniversities).Methods(http.MethodPost)

	// Get one university.
	router.HandleFunc("/tenant/{tenantID}/university/{universityID}",
		controller.GetUniversity).Methods(http.MethodGet)

	// Update university.
	router.HandleFunc("/tenant/{tenantID}/university/{universityID}",
		controller.UpdateUniversity).Methods(http.MethodPut)

	// Delete university.
	router.HandleFunc("/tenant/{tenantID}/university/{universityID}",
		controller.DeleteUniversity).Methods(http.MethodDelete)

	// Exculde routes.
	*exclude = append(*exclude, universityList)

	controller.log.Info("University routes registered.")
}

// AddUniversity will call add service to add new university.
func (controller *UniversityController) AddUniversity(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddUniversity Called==============================")
	university := &general.University{}
	parser := web.NewParser(r)
	// Parse university from request.
	err := web.UnmarshalJSON(r, university)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Parse and set tenant ID.
	university.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	university.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Validate university.
	err = university.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Add service call.
	err = controller.UniversityService.AddUniversity(university)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "University added successfully")
}

// AddUniversities if used to add multiple universities.
func (controller *UniversityController) AddUniversities(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddUniversities called==============================")
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

	universities := []*general.University{}
	// Parse university from request.
	err = web.UnmarshalJSON(r, &universities)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewValidationError("unable to parse requested data"))
		return
	}

	// Validate every university entry.
	for _, university := range universities {
		if err := university.Validate(); err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
		university.CreatedBy = credentialID
		university.TenantID = tenantID
	}

	err = controller.UniversityService.AddUniversities(universities)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// IDCollection will have the list of the UUIDs of the newly added universities.
	IDCollection := []uuid.UUID{}
	for _, university := range universities {
		if university != nil {
			IDCollection = append(IDCollection, university.ID)
		}
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Universities added successfully")
}

// GetUniversityList fetches all universities.
func (controller *UniversityController) GetUniversityList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllUniversities called==============================")
	parser := web.NewParser(r)
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	universities := &[]general.University{}
	err = controller.UniversityService.GetUniversityList(tenantID, universities)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, universities)
}

// GetUniversities fetches all universities.
func (controller *UniversityController) GetUniversities(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetUniversities called==============================")
	parser := web.NewParser(r)
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	var totalCount int
	universities := &[]general.University{}
	err = controller.UniversityService.GetUniversities(tenantID, universities, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, universities)
}

// GetUniversitiesByCountryList fetches all universities.
func (controller *UniversityController) GetUniversitiesByCountryList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllUniversitiesOfCountry called==============================")
	parser := web.NewParser(r)
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	countryID, err := parser.GetUUID("countryID")
	if err != nil {
		controller.log.Error("unable to parse country id")
		web.RespondError(w, errors.NewHTTPError("unable to parse country id", http.StatusBadRequest))
		return
	}

	universities := &[]general.University{}
	err = controller.UniversityService.GetUniversitiesByCountryList(tenantID, countryID, universities)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, universities)
}

// UpdateUniversity updates the given university.
func (controller *UniversityController) UpdateUniversity(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateUniversity Called==============================")
	parser := web.NewParser(r)
	university := &general.University{}

	// Parse university from request.
	err := web.UnmarshalJSON(r, university)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	university.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// UpdatedBy will be set with the given credentialID.
	university.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	university.ID, err = parser.GetUUID("universityID")
	if err != nil {
		controller.log.Error("Invalid university ID")
		web.RespondError(w, errors.NewValidationError("Invalid university ID"))
		return
	}

	// Validate university.
	err = university.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.UniversityService.UpdateUniversity(university)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "University updated successfully")
}

// GetUniversity returns a specific university.
func (controller *UniversityController) GetUniversity(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetUniversity called==============================")
	parser := web.NewParser(r)
	var err error
	university := &general.University{}

	university.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	university.ID, err = parser.GetUUID("universityID")
	if err != nil {
		controller.log.Error("Invalid university ID")
		web.RespondError(w, errors.NewValidationError("Invalid university ID"))
		return
	}

	err = controller.UniversityService.GetUniversity(university)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, university)

}

// DeleteUniversity deletes the given university.
func (controller *UniversityController) DeleteUniversity(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteUniversity==============================")
	parser := web.NewParser(r)
	var err error
	university := &general.University{}

	university.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	university.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	university.ID, err = parser.GetUUID("universityID")
	if err != nil {
		controller.log.Error("Invalid university ID")
		web.RespondError(w, errors.NewValidationError("Invalid university ID"))
		return
	}

	err = controller.UniversityService.DeleteUniversity(university)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "University deleted successfully")
}
