package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	services "github.com/techlabs/swabhav/tsam/company/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	model "github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// CompanyController Provide method to Update, Delete, Add, Get Method For Company.
type CompanyController struct {
	CompanyService *services.CompanyService
	log            log.Logger
	auth           *security.Authentication
}

// NewCompanyController Create New Instance Of CompanyController.
func NewCompanyController(comapnyService *services.CompanyService, log log.Logger, auth *security.Authentication) *CompanyController {
	return &CompanyController{
		CompanyService: comapnyService,
		log:            log,
		auth:           auth,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *CompanyController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/company", controller.AddCompany).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/company/{companyID}", controller.UpdateCompany).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/company/{companyID}", controller.DeleteCompany).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/company/{companyID}", controller.GetCompanyByID).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/company", controller.GetAllCompanies).Methods(http.MethodGet)

	controller.log.Info("Company Routes Registered.")
}

// AddCompany Add New Company
func (controller *CompanyController) AddCompany(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddCompany API Called==============================")

	company := &model.Company{}

	// Fill the company variable with given data.
	err := web.UnmarshalJSON(r, company)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	parser := web.NewParser(r)
	// Parse and set tenant ID.
	//  util.ParseUUID(mux.Vars(r)["tenantID"])
	company.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	// util.ParseUUID(mux.Vars(r)["credentialID"])
	company.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Validate company
	err = company.ValidateCompany()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// call Service
	err = controller.CompanyService.AddCompany(company)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "company added")
}

//UpdateCompany Update The company
func (controller *CompanyController) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateCompany API Called==============================")
	company := &model.Company{}
	// params := mux.Vars(r)

	// Parse company from request.
	err := web.UnmarshalJSON(r, company)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	parser := web.NewParser(r)
	// Assign company ID before sending it to validation.
	// util.ParseUUID(params["companyID"])
	company.ID, err = parser.GetUUID("companyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse company id", http.StatusBadRequest))
		return
	}

	// Validate company
	err = company.ValidateCompany()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// util.ParseUUID(params["tenantID"])
	company.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["credentialID"])
	company.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.CompanyService.UpdateCompany(company)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Company updated")
}

//DeleteCompany delete company
func (controller *CompanyController) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteCompany API Called==============================")
	company := &model.Company{}
	// params := mux.Vars(r)
	parser := web.NewParser(r)
	var err error

	// Parse and set tenant ID.
	// util.ParseUUID(params["tenantID"])
	company.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to company's DeletedBy field.
	// util.ParseUUID(params["credentialID"])
	company.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// util.ParseUUID(params["companyID"])
	company.ID, err = parser.GetUUID("companyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewValidationError("unable to parse company id"))
		return
	}

	err = controller.CompanyService.DeleteCompany(company)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Company deleted")
}

// GetCompanyByID returns one company
func (controller *CompanyController) GetCompanyByID(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================GetCompanyByID API Called==============================")
	company := &model.Company{}
	// params := mux.Vars(r)
	parser := web.NewParser(r)
	var err error

	// Assign Tenant ID
	// util.ParseUUID(params["tenantID"])
	company.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Assign college ID
	// util.ParseUUID(params["companyID"])
	company.ID, err = parser.GetUUID("companyID")
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse company ID", http.StatusBadRequest))
		return
	}

	err = controller.CompanyService.GetCompany(company)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, company)
}

// GetAllCompanies returns all company
func (controller *CompanyController) GetAllCompanies(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================GetAllCompany API Called==============================")

	parser := web.NewParser(r)
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	var totalCount int
	companies := &[]model.CompanyDTO{}

	err = controller.CompanyService.GetAllCompanies(tenantID, companies, parser, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, companies)
}
