package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingAssignmentController provides methods to do Update, Delete, Add, Get operations on programming_assignment.
type ProgrammingAssignmentController struct {
	log     log.Logger
	service *service.ProgrammingAssignmentService
	auth    *security.Authentication
}

// NewProgrammingAssignmentController creates new instance of ProgrammingAssignmentController.
func NewProgrammingAssignmentController(service *service.ProgrammingAssignmentService,
	log log.Logger, auth *security.Authentication) *ProgrammingAssignmentController {
	return &ProgrammingAssignmentController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *ProgrammingAssignmentController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/programming-assignment",
		controller.AddProgrammingAssignment).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/programming-assignment/{assignmentID}",
		controller.UpdateProgrammingAssignment).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/programming-assignment/{assignmentID}",
		controller.DeleteProgrammingAssignment).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/programming-assignment",
		controller.GetProgrammingAssignment).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/programming-assignment-list",
		controller.GetProgrammingAssignmentList).Methods(http.MethodGet)

	controller.log.Info("Programming Assignment Routes Registered")
}

// AddProgrammingAssignment will add new assignment question to the table.
func (controller *ProgrammingAssignmentController) AddProgrammingAssignment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add Programming Assignmnet Called==============================")
	assignment := programming.ProgrammingAssignment{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &assignment)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	assignment.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	assignment.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = assignment.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddProgrammingAssignment(&assignment)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Programming assignment successfully added")
}

// UpdateProgrammingAssignment will add new assignment question to the table.
func (controller *ProgrammingAssignmentController) UpdateProgrammingAssignment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Update Programming Assignmnet Called==============================")
	assignment := programming.ProgrammingAssignment{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &assignment)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	assignment.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	assignment.ID, err = parser.GetUUID("assignmentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming-assignment ID", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field.
	assignment.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = assignment.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdateProgrammingAssignment(&assignment)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Programming assignment successfully updated")
}

// DeleteProgrammingAssignment will add new assignment question to the table.
func (controller *ProgrammingAssignmentController) DeleteProgrammingAssignment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Delete Programming Assignmnet Called==============================")
	assignment := programming.ProgrammingAssignment{}
	parser := web.NewParser(r)
	var err error

	// Parse and set tenant ID.
	assignment.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	assignment.ID, err = parser.GetUUID("assignmentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("invalid Programming Assignment ID", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in DeletedBy field.
	assignment.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.DeleteProgrammingAssignment(&assignment)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Programming assignment successfully delete")
}

// GetProgrammingAssignment will return list of programming assignments.
func (controller *ProgrammingAssignmentController) GetProgrammingAssignment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Programming Assignment Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	var totalCount int
	assignments := &[]programming.ProgrammingAssignmentDTO{}

	err = controller.service.GetProgrammingAssignment(tenantID, assignments, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, assignments)
}

// GetProgrammingAssignmentList will return list of programming assignments.
func (controller *ProgrammingAssignmentController) GetProgrammingAssignmentList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Programming Assignment List Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	var totalCount int
	assignments := &[]list.ProgrammingAssignment{}
	err = controller.service.GetProgrammingAssignmentList(tenantID, assignments, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, assignments)
}
