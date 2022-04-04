package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/administration/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ManagementController gives method for listing for all users that is salesperson and admin.
type ManagementController struct {
	ManagementService *service.ManagementService
	auth         *security.Authentication
}

// NewManagementController creates new instance  ManagementController.
func NewManagementController(mgmtService *service.ManagementService, auth *security.Authentication) *ManagementController {
	return &ManagementController{
		ManagementService: mgmtService,
		auth:         auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *ManagementController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// get all employees
	router.HandleFunc("/tenant/{tenantID}/all-employees/limit/{limit}/offset/{offset}",
		controller.GetAllEmployees).Methods(http.MethodGet)

	// get all employees list
	router.HandleFunc("/tenant/{tenantID}/all-employee-list",
		controller.GetAllEmployeeList).Methods(http.MethodGet)

	// get direct reports
	router.HandleFunc("/tenant/{tenantID}/direct-reports/supervisor/{supervisorID}",
		controller.GetDirectReports).Methods(http.MethodGet)

	// add
	router.HandleFunc("/tenant/{tenantID}/supervisor", controller.AddSupervisor).Methods(http.MethodPost)

	// delete
	router.HandleFunc("/tenant/{tenantID}/supervisor/{supervisorID}/employee/{employeeID}",
		controller.DeleteSupervisor).Methods(http.MethodDelete)

	// Get faculty supervisor count.
	router.HandleFunc("/tenant/{tenantID}/faculty-supervisor-count",
		controller.GetFacultySupervisorCount).Methods(http.MethodGet)

	log.NewLogger().Info("Management Routes Registered")

}

// GetAllEmployees will return all the employees from the table with limit and offset
func (controller *ManagementController) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllEmployee call==============================")
	allEmployees := []admin.Employee{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Fill the request form.
	r.ParseForm()

	// limit,offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	err = controller.ManagementService.GetAllEmployees(tenantID, limit, offset, &totalCount, r.Form, &allEmployees)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allEmployees)
}

// GetAllEmployeeList will return all the employees from the table with limit and offset
func (controller *ManagementController) GetAllEmployeeList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllEmployee call==============================")
	allEmployees := []list.Credential{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ManagementService.GetAllEmployeeList(tenantID, &allEmployees)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, allEmployees)
}

// AddSupervisor will add new employee to the table and also create a login for the employee
func (controller *ManagementController) AddSupervisor(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddSupervisor called==============================")
	supervisor := admin.EmployeeSupervisor{}
	param := mux.Vars(r)

	err := web.UnmarshalJSON(r, &supervisor)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = supervisor.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ManagementService.AddSupervisor(tenantID, &supervisor)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Supervisor successfully added")
}

// GetDirectReports returns a list of direct reports of the particular supervisor.
func (controller *ManagementController) GetDirectReports(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetDirectReports called==============================")
	params := mux.Vars(r)
	directReports := &[]list.Credential{}

	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
	}
	supervisorID, err := util.ParseUUID(params["supervisorID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
	}

	err = controller.ManagementService.GetDirectReports(tenantID, supervisorID, directReports)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, directReports)
}

// DeleteSupervisor deletes specific supervisor by id.
func (controller *ManagementController) DeleteSupervisor(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteSupervisor called==============================")
	param := mux.Vars(r)
	supervisor := &admin.EmployeeSupervisor{}
	var err error

	supervisor.EmployeeCredentialID, err = util.ParseUUID(param["employeeID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	supervisor.SupervisorCredentialID, err = util.ParseUUID(param["supervisorID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.ManagementService.DeleteSupervisor(supervisor)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Supervisor deleted successfully")
}

// GetFacultySupervisorCount gets count of faculty supervisors.
func (controller *ManagementController) GetFacultySupervisorCount(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetFacultySupervisorCount called=======================================")

	// Create bucket.
	totalCount := admin.CountModel{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call get talents method.
	err = controller.ManagementService.GetFacultySupervisorCount(&totalCount, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, totalCount)
}
