package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingConceptController provides methods to do CRUD operations.
type ProgrammingConceptController struct {
	ProgrammingConceptService *service.ProgrammingConceptService
	log                       log.Logger
	auth                      *security.Authentication
}

// NewProgrammingConceptController creates new instance of concept controller.
func NewProgrammingConceptController(programmingConceptService *service.ProgrammingConceptService,
	log log.Logger, auth *security.Authentication) *ProgrammingConceptController {
	return &ProgrammingConceptController{
		ProgrammingConceptService: programmingConceptService,
		log:                       log,
		auth:                      auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *ProgrammingConceptController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add programming concept.
	router.HandleFunc("/tenant/{tenantID}/programming-concept/credential/{credentialID}",
		controller.AddProgrammingConcept).Methods(http.MethodPost)

	// Update programming concept.
	router.HandleFunc("/tenant/{tenantID}/programming-concept/{conceptID}/credential/{credentialID}",
		controller.UpdateProgrammingConcept).Methods(http.MethodPut)

	// Delete programming concept.
	router.HandleFunc("/tenant/{tenantID}/programming-concept/{conceptID}/credential/{credentialID}",
		controller.DeleteProgrammingConcept).Methods(http.MethodDelete)

	// Get all programming concepts with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/programming-concepts",
		controller.GetAllProgrammingConcepts).Methods(http.MethodGet)

	// Get one programming concept.
	router.HandleFunc("/tenant/{tenantID}/programming-concept/{conceptID}",
		controller.GetProgrammingConcept).Methods(http.MethodGet)

	log.NewLogger().Info("Programming Concept Route Registered")
}

// AddProgrammingConcept will add the new dpeartment record in the table.
func (controller *ProgrammingConceptController) AddProgrammingConcept(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddProgrammingConcept called==============================")

	// Create bucket for params.
	param := mux.Vars(r)

	// Create bucket for pragramming concept.
	concept := programming.ProgrammingConcept{}

	// Unmarshal JSON.
	err := web.UnmarshalJSON(r, &concept)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant id.
	concept.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Prade and set credential id.
	concept.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Validate programming concept compulsary fields.
	err = concept.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call add service method.
	err = controller.ProgrammingConceptService.AddProgrammingConcept(&concept)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Programming concept added successfully")
}

// UpdateProgrammingConcept will update the specified concept record in the table.
func (controller *ProgrammingConceptController) UpdateProgrammingConcept(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateProgrammingConcept called==============================")

	// Create bucket for params.
	param := mux.Vars(r)

	// Create bucket for pragramming concept.
	concept := programming.ProgrammingConcept{}

	// Unmarshal JSON.
	err := web.UnmarshalJSON(r, &concept)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant id.
	concept.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Prade and set credential id.
	concept.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Prade and set concept id.
	concept.ID, err = util.ParseUUID(param["conceptID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Validate programming concept compulsary fields.
	err = concept.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.ProgrammingConceptService.UpdateProgrammingConcept(&concept)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Programming concept updated successfully")
}

// DeleteProgrammingConcept will delete the specified department record from the table.
func (controller *ProgrammingConceptController) DeleteProgrammingConcept(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteProgrammingConcept called==============================")

	// Create bucket for params.
	param := mux.Vars(r)

	// Create bucket for pragramming concept.
	concept := programming.ProgrammingConcept{}

	// Create error varoable.
	var err error

	// Parse and set tenant id.
	concept.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Prade and set credential id.
	concept.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Prade and set concept id.
	concept.ID, err = util.ParseUUID(param["conceptID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call delete service method.
	err = controller.ProgrammingConceptService.DeleteProgrammingConcept(&concept)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Programming concept deleted successfully")
}

// GetAllProgrammingConcepts will return all the records from concept table.
func (controller *ProgrammingConceptController) GetAllProgrammingConcepts(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllProgrammingConcepts called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Create bucket for pragramming concept.
	concepts := []programming.ProgrammingConceptDTO{}
	pagination := repository.Pagination{}

	// Create variables for total, limit, offset.
	// var totalCount int
	pagination.Limit, pagination.Offset = parser.ParseLimitAndOffset()

	// Call get all service method.
	err = controller.ProgrammingConceptService.GetAllProgrammingConcepts(&concepts, tenantID, r.Form, &pagination)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, pagination.TotalCount, concepts)
}

// GetProgrammingConcept will return specified record from department table.
func (controller *ProgrammingConceptController) GetProgrammingConcept(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetProgrammingConcept called==============================")

	// Create bucket for params.
	param := mux.Vars(r)

	// Create error varoable.
	var err error

	// Create bucket for pragramming concept.
	concept := programming.ProgrammingConceptDTO{}

	// Parse and set tenant id.
	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Prade and set concept id.
	concept.ID, err = util.ParseUUID(param["conceptID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call get one service method.
	err = controller.ProgrammingConceptService.GetProgrammingConcept(&concept, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, concept)
}
