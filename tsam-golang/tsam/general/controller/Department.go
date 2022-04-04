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

// DepartmentController provides methods to do CRUD operations.
type DepartmentController struct {
	log               log.Logger
	auth              *security.Authentication
	DepartmentService *service.DepartmentService
}

// NewDepartmentController creates new instance of DepartmentController.
func NewDepartmentController(generalService *service.DepartmentService, log log.Logger, auth *security.Authentication) *DepartmentController {
	return &DepartmentController{
		DepartmentService: generalService,
		log:               log,
		auth:              auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *DepartmentController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add one department.
	router.HandleFunc("/tenant/{tenantID}/department",
		controller.AddDepartment).Methods(http.MethodPost)

	// Update department.
	router.HandleFunc("/tenant/{tenantID}/department/{departmentID}",
		controller.UpdateDepartment).Methods(http.MethodPut)

	// Delete department.
	router.HandleFunc("/tenant/{tenantID}/department/{departmentID}",
		controller.DeleteDepartment).Methods(http.MethodDelete)

	// Get all department with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/department",
		controller.GetAllDepartments).Methods(http.MethodGet)

	// Get department list.
	router.HandleFunc("/tenant/{tenantID}/department-list",
		controller.GetDepartmentList).Methods(http.MethodGet)

	// Get one department.
	router.HandleFunc("/tenant/{tenantID}/department/{departmentID}",
		controller.GetDepartment).Methods(http.MethodGet)

	controller.log.Info("Department Route Registered")
}

// AddDepartment will add the new dpeartment record in the table.
func (controller *DepartmentController) AddDepartment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddDepartment called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	department := general.Department{}

	//param := mux.Vars(r)

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &department)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = department.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	department.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	department.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.DepartmentService.AddDepartment(&department)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Department added successfully")
}

// UpdateDepartment will update the specified department record in the table.
func (controller *DepartmentController) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateDepartment called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	department := general.Department{}

	//param := mux.Vars(r)

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &department)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate city.
	err = department.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	department.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	department.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting department id from param and parsing it to uuid.
	department.ID, err = parser.GetUUID("departmentID")
	if err != nil {
		controller.log.Error("unable to parse department id")
		web.RespondError(w, errors.NewHTTPError("unable to parse department id", http.StatusBadRequest))
		return
	}

	// Cal update service method.
	err = controller.DepartmentService.UpdateDepartment(&department)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Department updated successfully")
}

// DeleteDepartment will delete the specified department record from the table.
func (controller *DepartmentController) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteDepartment called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	department := general.Department{}

	//param := mux.Vars(r)

	var err error

	// Getting tenant id from param and parsing it to uuid.
	department.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	department.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting department id from param and parsing it to uuid.
	department.ID, err = parser.GetUUID("departmentID")
	if err != nil {
		controller.log.Error("unable to parse department id")
		web.RespondError(w, errors.NewHTTPError("unable to parse department id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.DepartmentService.DeleteDepartment(&department)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Department deleted successfully")
}

// GetAllDepartments will return all the records from department table.
func (controller *DepartmentController) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllDepartments called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	departments := []general.DepartmentDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	var totalCount int
	// Call get all service method.
	err = controller.DepartmentService.GetAllDepartments(&departments, tenantID, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, departments)
}

// GetDepartmentList returns department list bu roles.
func (controller *DepartmentController) GetDepartmentList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetDepartmentList called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	departments := []general.DepartmentDTO{}

	//param := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form
	r.ParseForm()

	// Call get all service method.
	err = controller.DepartmentService.GetDepartmentList(&departments, tenantID, r.Form)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, departments)
}

// GetDepartment will return specified record from department table.
func (controller *DepartmentController) GetDepartment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetDepartment called==============================")
	parser := web.NewParser(r)
	// Create bucket.
	department := general.DepartmentDTO{}

	//param := mux.Vars(r)

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

	// Call get one service method.
	err = controller.DepartmentService.GetDepartment(&department, tenantID, departmentID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, department)
}
