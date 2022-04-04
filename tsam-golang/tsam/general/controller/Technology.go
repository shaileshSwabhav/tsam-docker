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

// TechnologyController provide method to update, delete, add, get for Technology.
type TechnologyController struct {
	log               log.Logger
	TechnologyService *service.TechnologyService
	auth              *security.Authentication
}

// NewTechnologyController creates new instance of TechnologyController.
func NewTechnologyController(enquiryservice *service.TechnologyService, log log.Logger, auth *security.Authentication) *TechnologyController {
	return &TechnologyController{
		TechnologyService: enquiryservice,
		log:               log,
		auth:              auth,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *TechnologyController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all technologies by limit and offset.
	router.HandleFunc("/tenant/{tenantID}/technology",
		controller.GetTechnologies).Methods(http.MethodGet)

	// Get technology list.
	technoligyList := router.HandleFunc("/tenant/{tenantID}/technology-list",
		controller.GetTechnologyList).Methods(http.MethodGet)

	// Get one technology.
	router.HandleFunc("/tenant/{tenantID}/technology/{technologyID}",
		controller.GetTechnology).Methods(http.MethodGet)

	// Add one technology.
	router.HandleFunc("/tenant/{tenantID}/technology",
		controller.AddTechnology).Methods(http.MethodPost)

	// Add multiple technologies.
	router.HandleFunc("/tenant/{tenantID}/technologies",
		controller.AddTechnologies).Methods(http.MethodPost)

	// Update technology.
	router.HandleFunc("/tenant/{tenantID}/technology/{technologyID}",
		controller.UpdateTechnology).Methods(http.MethodPut)

	// Delete technology.
	router.HandleFunc("/tenant/{tenantID}/technology/{technologyID}",
		controller.DeleteTechnology).Methods(http.MethodDelete)

	// Exculde routes.
	*exclude = append(*exclude, technoligyList)

	controller.log.Info("Technology Route Registered")
}

// GetTechnologies returns all the technologies.
func (controller *TechnologyController) GetTechnologies(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetTechnologies called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	technologies := []general.Technology{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// For pagination.
	var totalCount int

	// Call get all service method.
	err = controller.TechnologyService.GetTechnologies(&technologies, parser, tenantID, &totalCount)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, technologies)
}

// GetTechnologyList returns all the technologies (without pagination).
func (controller *TechnologyController) GetTechnologyList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetTechnologyList call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	technologies := []general.Technology{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	r.ParseForm()

	// Call get all service method.
	err = controller.TechnologyService.GetTechnologyList(&technologies, tenantID, r.Form)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, technologies)
}

// GetTechnology returns the specified technology by id.
func (controller *TechnologyController) GetTechnology(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("**********************************GetTechnology Call*******************************")
	parser := web.NewParser(r)
	technology := general.Technology{}
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	technologyID, err := parser.GetUUID("technologyID")
	if err != nil {
		controller.log.Error("unable to parse technology id")
		web.RespondError(w, errors.NewHTTPError("unable to parse technology id", http.StatusBadRequest))
		return
	}
	err = controller.TechnologyService.GetTechnology(&technology, tenantID, technologyID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, technology)
}

// AddTechnology adds new technology.
func (controller *TechnologyController) AddTechnology(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddTechnology call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	technology := general.Technology{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &technology)
	if err != nil {
		controller.log.Error("Invalid Request")
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = technology.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	technology.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("Unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	technology.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.TechnologyService.AddTechnology(&technology)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Technology added successfully")
}

// AddTechnologies adds multiple technologies.
func (controller *TechnologyController) AddTechnologies(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddTechnologies call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	technologiesIDs := []uuid.UUID{}
	technologies := []general.Technology{}

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
	err = web.UnmarshalJSON(r, &technologies)
	if err != nil {
		controller.log.Error("Unable to parse data")
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	for _, technology := range technologies {
		err = technology.Validate()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	// Call add multiple service method.
	err = controller.TechnologyService.AddTechnologies(&technologies, &technologiesIDs, tenantID, credentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Technologies added successfully")
}

// UpdateTechnology updates the specified technology by id.
func (controller *TechnologyController) UpdateTechnology(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************UpdateTechnology call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	technology := general.Technology{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	technology.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("Unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting technlogy id from param and parsing it to uuid.
	technology.ID, err = parser.GetUUID("technologyID")
	if err != nil {
		controller.log.Error("Unable to parse technology id")
		web.RespondError(w, errors.NewHTTPError("Unable to parse technology id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	technology.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &technology)
	if err != nil {
		controller.log.Error("Unable to parse data")
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate technology fields.
	err = technology.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.TechnologyService.UpdateTechnology(&technology)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Technology updated successfully")
}

// DeleteTechnology deletes specific technology by id.
func (controller *TechnologyController) DeleteTechnology(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************DeleteTechnology call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	technology := general.Technology{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	technology.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("Unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	technology.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting technology id from param and parsing it to uuid.
	technology.ID, err = parser.GetUUID("technologyID")
	if err != nil {
		controller.log.Error("Unable to parse technology id")
		web.RespondError(w, errors.NewHTTPError("Unable to parse technology id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.TechnologyService.DeleteTechnology(&technology)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Technology deleted successfully")
}
