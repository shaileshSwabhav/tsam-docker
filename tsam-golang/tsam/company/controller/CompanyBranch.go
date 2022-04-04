package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	services "github.com/techlabs/swabhav/tsam/company/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	model "github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CompanyBranchController provides methods to do Update, Delete, Add, Get operations on CompanyBranch.
type CompanyBranchController struct {
	CompanyBranchService *services.CompanyBranchService
	log                  log.Logger
	auth                 *security.Authentication
}

// NewCompanyBranchController creates new instance of CompanyBranchController.
func NewCompanyBranchController(companyBranchService *services.CompanyBranchService, log log.Logger, auth *security.Authentication) *CompanyBranchController {
	return &CompanyBranchController{
		CompanyBranchService: companyBranchService,
		log:                  log,
		auth:                 auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *CompanyBranchController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/company/{companyID}/branch",
		controller.AddCompanyBranch).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/company/{companyID}/branches",
		controller.AddCompanyBranches).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/company/{companyID}/branch/{branchID}",
		controller.UpdateCompanyBranch).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/company/{companyID}/branch/{branchID}",
		controller.DeleteCompanyBranch).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/company-branch",
		controller.GetAllBranches).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/company/{companyID}/branch",
		controller.GetAllBranchesOfCompany).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/company/{companyID}/branch/{branchID}",
		controller.GetCompanyBranch).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/company/branch/salesperson/{salesPersonID}",
		controller.GetAllBranchesForSalesPerson).Methods(http.MethodGet)

	// get company branches list
	companyBranchList := router.HandleFunc("/tenant/{tenantID}/company/branch/list",
		controller.GetCompanyBranchList).Methods(http.MethodGet)

	controller.log.Info("Company Branch Routes Registered")

	// Exculde routes.
	*exclude = append(*exclude, companyBranchList)
}

// AddCompanyBranch adds a new company branch
func (controller *CompanyBranchController) AddCompanyBranch(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================AddCompanyBranch Called==============================")

	companyBranch := &model.Branch{}
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the college branch variable with given data.
	err := web.UnmarshalJSON(r, companyBranch)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	// util.ParseUUID(params["tenantID"])
	companyBranch.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["companyID"])
	companyBranch.CompanyID, err = parser.GetUUID("comapnyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse company id", http.StatusBadRequest))
		return
	}

	//  util.ParseUUID(params["credentialID"])
	companyBranch.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Validate company branch
	err = companyBranch.ValidateCompanyBranch()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// call add service
	err = controller.CompanyBranchService.AddCompanyBranch(companyBranch)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "company branch added")
}

// AddCompanyBranches adds a new company branch
func (controller *CompanyBranchController) AddCompanyBranches(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================AddCompanyBranches Called==============================")

	companyBranches := &[]model.Branch{}
	companyBranchesIDs := &[]uuid.UUID{}
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the college branch variable with given data.
	err := web.UnmarshalJSON(r, companyBranches)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	// util.ParseUUID(params["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["companyID"])
	companyID, err := parser.GetUUID("companyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse company id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["credentialID"])
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Validate company branch
	for _, companyBranch := range *companyBranches {
		err = companyBranch.ValidateCompanyBranch()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	// call add service
	err = controller.CompanyBranchService.AddCompanyBranches(companyBranches, companyBranchesIDs, companyID, tenantID, credentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "company branches added")
}

// UpdateCompanyBranch updates the specific company branch
func (controller *CompanyBranchController) UpdateCompanyBranch(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================UpdateCompanyBranch Called==============================")

	companyBranch := &model.Branch{}
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Parse college branch from request.
	err := web.UnmarshalJSON(r, companyBranch)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Validate company branch
	err = companyBranch.ValidateCompanyBranch()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	//  util.ParseUUID(params["tenantID"])
	companyBranch.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["companyID"])
	companyBranch.CompanyID, err = parser.GetUUID("companyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse company id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["credentialID"])
	companyBranch.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	//  util.ParseUUID(params["branchID"])
	companyBranch.ID, err = parser.GetUUID("branchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse branch id", http.StatusBadRequest))
		return
	}

	// call update service
	err = controller.CompanyBranchService.UpdateCompanyBranch(companyBranch)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Company Branch Updated")
}

// DeleteCompanyBranch deletes the specific Company Branch
func (controller *CompanyBranchController) DeleteCompanyBranch(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteCompanyBranch Called==============================")

	// params := mux.Vars(r)
	parser := web.NewParser(r)
	companyBranch := &model.Branch{}
	var err error

	// Parse and set tenant ID.
	// util.ParseUUID(params["tenantID"])
	companyBranch.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse college ID and assign to branch
	// util.ParseUUID(params["companyID"])
	companyBranch.CompanyID, err = util.ParseUUID("companyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse company id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to branch's DeletedBy field.
	// util.ParseUUID(mux.Vars(r)["credentialID"])
	companyBranch.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse branch ID and assign it to branch
	// util.ParseUUID(mux.Vars(r)["branchID"])
	companyBranch.ID, err = parser.GetUUID("branchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse branch id", http.StatusBadRequest))
		return
	}

	err = controller.CompanyBranchService.DeleteCompanyBranch(companyBranch)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Company Branch Deleted")
}

// GetAllBranches returns all branches in database as response
func (controller *CompanyBranchController) GetAllBranches(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================GetAllBranches Called==============================")

	parser := web.NewParser(r)

	//  util.ParseUUID(mux.Vars(r)["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form
	r.ParseForm()

	// limit,offset & totalCount for pagination
	var totalCount int
	// limit, offset := web.GetLimitAndOffset(r)

	allBranches := &[]*model.Branch{}
	err = controller.CompanyBranchService.GetAllBranches(tenantID, allBranches, parser, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allBranches)
}

// GetAllBranchesForSalesPerson returns all branches in where the given sales person is assigned.
func (controller *CompanyBranchController) GetAllBranchesForSalesPerson(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllBranchesForSalesPerson Called==============================")

	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// util.ParseUUID(params["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["salesPersonID"])
	salesPersonID, err := parser.GetUUID("salesPersonID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse sales person id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form
	// r.ParseForm()

	// limit,offset & totalCount for pagination
	var totalCount int
	// limit, offset := web.GetLimitAndOffset(r)
	limit, offset := parser.ParseLimitAndOffset()
	Branches := &[]*model.Branch{}
	fmt.Println("================================limit", limit, offset)
	err = controller.CompanyBranchService.GetAllBranchesForSalesPerson(tenantID, salesPersonID, Branches, parser.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, Branches)
}

// GetCompanyBranch returns a specific company branch as response
func (controller *CompanyBranchController) GetCompanyBranch(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================GetCompanyBranch Called==============================")
	var err error
	companyBranch := &model.Branch{}
	params := mux.Vars(r)

	// Assign Tenant ID
	companyBranch.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Assign Company ID
	companyBranch.CompanyID, err = util.ParseUUID(params["companyID"])
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse company id", http.StatusBadRequest))
		return
	}

	// Assign Company branch ID
	companyBranch.ID, err = util.ParseUUID(params["branchID"])
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse branch id", http.StatusBadRequest))
		return
	}

	err = controller.CompanyBranchService.GetCompanyBranch(companyBranch)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status
	web.RespondJSON(w, http.StatusOK, companyBranch)
}

// GetAllBranchesOfCompany returns all branches of a speicifc company as response
func (controller *CompanyBranchController) GetAllBranchesOfCompany(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllBranchesOfCompany Called==============================")

	params := mux.Vars(r)

	// Parse tenant ID
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	companyID, err := util.ParseUUID(params["companyID"])
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse company id", http.StatusBadRequest))
		return
	}

	Branches := &[]model.Branch{}
	err = controller.CompanyBranchService.GetAllBranchesOfCompany(tenantID, companyID, Branches)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, Branches)
}

// GetCompanyBranchList will return list of company branches
func (controller *CompanyBranchController) GetCompanyBranchList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetCompanyBranchList Called==============================")
	params := mux.Vars(r)

	// Parse tenant ID
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	branches := &[]list.CompanyBranch{}
	err = controller.CompanyBranchService.GetAllCompanyBranchList(branches, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, branches)
}

// adds all search queries by comparing with the CompanyBranch data recieved from
//
// POST: GetSearchedBranches()
// func (controller *CompanyBranchController) addSearchQueries(branch *model.Branch) repository.QueryProcessor {
// 	var columnNames []string
// 	var conditions []string
// 	var operators []string
// 	var values []interface{}
// 	if state := branch.State; state != nil {
// 		util.AddToSlice("state_id", "= ?", "AND", state.ID, &columnNames, &conditions, &operators, &values)
// 	}
// 	if country := branch.Country; country != nil {
// 		util.AddToSlice("country_id", "= ?", "AND", country.ID, &columnNames, &conditions, &operators, &values)
// 	}
// 	if salesPersonID := branch.SalesPersonID; salesPersonID != nil {
// 		util.AddToSlice("sales_person_id", "= ?", "AND", salesPersonID, &columnNames, &conditions, &operators, &values)
// 	}

// 	if companyRating := branch.CompanyRating; companyRating != nil {
// 		util.AddToSlice("company_rating", ">= ?", "AND", companyRating, &columnNames, &conditions, &operators, &values)
// 	}
// 	if numberOfEmployees := branch.NumberOfEmployees; numberOfEmployees != nil {
// 		util.AddToSlice("number_of_employees", ">= ?", "AND", numberOfEmployees, &columnNames, &conditions, &operators, &values)
// 	}
// 	if city := branch.City; city != nil && util.IsEmpty(*city) {
// 		util.AddToSlice("city", "LIKE ?", "AND", "%"+*city+"%", &columnNames, &conditions, &operators, &values)
// 	}
// 	if companyName := branch.CompanyName; companyName != nil {
// 		util.AddToSlice("company_name", "LIKE ?", "AND", "%"+*companyName+"%", &columnNames, &conditions, &operators, &values)
// 	}
// 	if hrHeadName := branch.HRHeadName; hrHeadName != nil {
// 		util.AddToSlice("hr_head_name", "LIKE ?", "AND", "%"+*hrHeadName+"%", &columnNames, &conditions, &operators, &values)
// 	}
// 	if companyCode := branch.Code; len(companyCode) > 0 {
// 		util.AddToSlice("company_code", "LIKE ?", "AND", "%"+companyCode+"%", &columnNames, &conditions, &operators, &values)
// 	}
// 	if technologies := branch.Technologies; len(technologies) > 0 {
// 		var technologyIDs, branchIDs []string
// 		for _, technology := range technologies {
// 			technologyIDs = append(technologyIDs, technology.ID.String())
// 		}
// 		repository.PluckColumn(controller.CompanyBranchService.DB, "company_branch_technologies", "company_branch_id",
// 			&branchIDs, repository.Filter("technology_id IN(?)", technologyIDs))
// 		util.AddToSlice("id", "IN(?)", "AND", branchIDs, &columnNames, &conditions, &operators, &values)
// 	}
// 	if domains := branch.Domains; len(domains) > 0 {
// 		var domainIDs, branchIDs []string
// 		for _, domain := range domains {
// 			domainIDs = append(domainIDs, domain.ID.String())
// 		}
// 		repository.PluckColumn(controller.CompanyBranchService.DB, "company_branch_domains", "company_branch_id", &branchIDs,
// 			repository.Filter("domain_id IN(?)", domainIDs))
// 		util.AddToSlice("id", "IN(?)", "AND", branchIDs, &columnNames, &conditions, &operators, &values)

// 	}
// 	return repository.FilterWithOperator(columnNames, conditions, operators, values)
// }
