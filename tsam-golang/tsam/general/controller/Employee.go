package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// EmployeeController provides methods to do CRUD operations.
type EmployeeController struct {
	EmployeeService *service.EmployeeService
}

// NewEmployeeController creates new instance of employee controller.
func NewEmployeeController(generalService *service.EmployeeService) *EmployeeController {
	return &EmployeeController{
		EmployeeService: generalService,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *EmployeeController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/other-employee/credential/{credentialID}",
		controller.AddEmployee).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/other-employee/{employeeID}/credential/{credentialID}",
		controller.UpdateEmployee).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/other-employee/{employeeID}/credential/{credentialID}",
		controller.DeleteEmployee).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/other-employee/limit/{limit}/offset/{offset}",
		controller.GetAllEmployee).Methods(http.MethodGet)

	// get list
	router.HandleFunc("/tenant/{tenantID}/other-employee-list",
		controller.GetAllEmployeeList).Methods(http.MethodGet)

	log.NewLogger().Info("Employee Route Registered")
}

// AddEmployee will add new employee to the table and also create a login for the employee
func (controller *EmployeeController) AddEmployee(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddEmployee call==============================")
	employee := general.Employee{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &employee)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	employee.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	employee.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = employee.ValidateEmployee()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.EmployeeService.AddEmployee(&employee)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Employee added successfully")
}

// UpdateEmployee will update existing employee in the table
func (controller *EmployeeController) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateEmployee call==============================")
	employee := general.Employee{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &employee)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	employee.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	employee.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	employee.ID, err = util.ParseUUID(param["employeeID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = employee.ValidateEmployee()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.EmployeeService.UpdateEmployee(&employee)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Employee updated successfully")
}

// DeleteEmployee deletes an employee from database.
func (controller *EmployeeController) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteEmployee call==============================")
	employee := general.Employee{}
	param := mux.Vars(r)
	var err error

	employee.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	employee.ID, err = util.ParseUUID(param["employeeID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	employee.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.EmployeeService.DeleteEmployee(&employee)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Employee deleted successfully")
}

// GetAllEmployee will return all the employees from the table with limit and offset
func (controller *EmployeeController) GetAllEmployee(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllEmployee call==============================")
	employees := []general.EmployeeDTO{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Fill the r.Form
	r.ParseForm()

	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	err = controller.EmployeeService.GetAllEmployee(&employees, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, employees)
}

// GetAllEmployeeList will return an employee list with is_active true.
func (controller *EmployeeController) GetAllEmployeeList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllEmployeeList call==============================")
	employees := []college.Developer{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	err = controller.EmployeeService.GetAllEmployeeList(&employees, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, employees)
}
