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

// DegreeController provides method to update, delete, add, get For Degree.
type DegreeController struct {
	log           log.Logger
	auth          *security.Authentication
	DegreeService *service.DegreeService
}

// NewDegreeController creates new instance of DegreeController.
func NewDegreeController(degreeService *service.DegreeService, log log.Logger, auth *security.Authentication) *DegreeController {
	return &DegreeController{
		DegreeService: degreeService,
		log:           log,
		auth:          auth,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *DegreeController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	
	// Get all degree list.
	// degreeList := router.HandleFunc("/tenant/{tenantID}/degree-list",
	// 	controller.GetDegreeList).Methods(http.MethodGet)

	// Get all degrees with limit and offset.
	degreeList := router.HandleFunc("/tenant/{tenantID}/degree",
		controller.GetDegrees).Methods(http.MethodGet)

	// Get one degree by id.
	router.HandleFunc("/tenant/{tenantID}/degree/{degreeID}",
		controller.GetDegree).Methods(http.MethodGet)

	// Add one degree.
	router.HandleFunc("/tenant/{tenantID}/degree",
		controller.AddDegree).Methods(http.MethodPost)

	// Add multiple degrees.
	router.HandleFunc("/tenant/{tenantID}/degrees",
		controller.AddDegrees).Methods(http.MethodPost)

	// Update one degree.
	router.HandleFunc("/tenant/{tenantID}/degree/{degreeID}",
		controller.UpdateDegree).Methods(http.MethodPut)

	// Delete one degree.
	router.HandleFunc("/tenant/{tenantID}/degree/{degreeID}",
		controller.DeleteDegree).Methods(http.MethodDelete)

	// Exculde routes.
	*exclude = append(*exclude, degreeList)

	controller.log.Info("Degree Route Registered")
}

// GetDegreeList returns degree list.
func (controller *DegreeController) GetDegreeList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetDegreeList call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	degrees := []general.Degree{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get degree list method.
	err = controller.DegreeService.GetDegreeList(&degrees, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, degrees)
}

// GetDegrees returns all degrees.
func (controller *DegreeController) GetDegrees(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetDegrees call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	degrees := []general.Degree{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	var totalCount int
	// Call get all degrees service method.
	err = controller.DegreeService.GetDegrees(&degrees, tenantID, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, degrees)
}

// GetDegree return specific degree by id.
func (controller *DegreeController) GetDegree(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetDegree call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	degree := general.Degree{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	degree.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant ID")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant ID", http.StatusBadRequest))
		return
	}

	// Getting degree id from param and parsing it to uuid.
	degree.ID, err = parser.GetUUID("degreeID")
	if err != nil {
		controller.log.Error("unable to parse degree ID")
		web.RespondError(w, errors.NewHTTPError("unable to parse degree ID", http.StatusBadRequest))
		return
	}

	// Call get degree by id service method.
	err = controller.DegreeService.GetDegree(&degree)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, degree)
}

// AddDegree adds new degree.
func (controller *DegreeController) AddDegree(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddDegree call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	degree := general.Degree{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &degree)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := degree.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	degree.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	degree.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.DegreeService.AddDegree(&degree)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Degree added successfully")
}

// AddDegrees adds multiple degrees.
func (controller *DegreeController) AddDegrees(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddDegrees call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	degreesIDs := []uuid.UUID{}
	degrees := []general.Degree{}

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
	err = web.UnmarshalJSON(r, &degrees)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary degree fields.
	for _, degree := range degrees {
		err = degree.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	// Call add multiple degree service method.
	err = controller.DegreeService.AddDegrees(&degrees, &degreesIDs, tenantID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Degrees added successfully")
}

// UpdateDegree updates the specified degree by id.
func (controller *DegreeController) UpdateDegree(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************UpdateDegree call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	degree := general.Degree{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	degree.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant ID")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant ID", http.StatusBadRequest))
		return
	}

	// Getting degree id from param and parsing it to uuid.
	degree.ID, err = parser.GetUUID("degreeID")
	if err != nil {
		controller.log.Error("unable to parse degree ID")
		web.RespondError(w, errors.NewHTTPError("unable to parse degree ID", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	degree.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &degree)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate country fields
	err = degree.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.DegreeService.UpdateDegree(&degree)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Degree updated successfully")
}

// DeleteDegree deletes the specified degree by id.
func (controller *DegreeController) DeleteDegree(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************DeleteDegree call**************************************")
	parser := web.NewParser(r)
	// Create bcuket.
	degree := general.Degree{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	degree.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant ID")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant ID", http.StatusBadRequest))
		return
	}

	// Getting degree id from param and parsing it to uuid.
	degree.ID, err = parser.GetUUID("degreeID")
	if err != nil {
		controller.log.Error("unable to parse degree ID")
		web.RespondError(w, errors.NewHTTPError("unable to parse degree ID", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	degree.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.DegreeService.DeleteDegree(&degree)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Degree deleted successfully")
}
