package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	services "github.com/techlabs/swabhav/tsam/company/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	cmp "github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CompanyRequirementController Provide method to Update, Delete, Add, Get Method For Company Requirement.
type CompanyRequirementController struct {
	CompanyRequirementService *services.CompanyRequirementService
	log                       log.Logger
	auth                      *security.Authentication
}

// NewCompanyRequirementController Create New Instance Of CompanyRequirementController.
func NewCompanyRequirementController(requirementService *services.CompanyRequirementService, log log.Logger, auth *security.Authentication) *CompanyRequirementController {
	return &CompanyRequirementController{
		CompanyRequirementService: requirementService,
		log:                       log,
		auth:                      auth,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *CompanyRequirementController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get my opportunities.
	router.HandleFunc("/tenant/{tenantID}/company-requirement/my-opportunities",
		controller.GetMyOpportunities).Methods(http.MethodGet)

	// Get requirement list.
	requirementList := router.HandleFunc("/tenant/{tenantID}/company-requirement/list",
		controller.GetRequirementList).Methods(http.MethodGet)

	// Get one requirement by id.
	router.HandleFunc("/tenant/{tenantID}/company-requirement/{companyRequirementID}",
		controller.GetCompanyRequirement).Methods(http.MethodGet)

	// Get all requiremnets with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/company-requirement",
		controller.GetAllCompanyRequirements).Methods(http.MethodGet)

	// Get company details.
	router.HandleFunc("/tenant/{tenantID}/company-requirement/{companyRequirementID}/company-details",
		controller.GetCompanyDetails).Methods(http.MethodGet)

	// Add one requirement.
	router.HandleFunc("/tenant/{tenantID}/company-requirement",
		controller.AddCompanyRequirement).Methods(http.MethodPost)

	// Update one requirement.
	router.HandleFunc("/tenant/{tenantID}/company-requirement/{companyRequirementID}",
		controller.UpdateCompanyRequirement).Methods(http.MethodPut)

	// Update company requirement' salesperson.
	router.HandleFunc("/tenant/{tenantID}/company-requirement/salesperson/{salesPersonID}",
		controller.UpdateCompanyRequirementsSalesperson).Methods(http.MethodPut)

	// Close one requirement.
	router.HandleFunc("/tenant/{tenantID}/company-requirement/close/{companyRequirementID}",
		controller.CloseCompanyRequirement).Methods(http.MethodDelete)

	// Delete one requirement.
	router.HandleFunc("/tenant/{tenantID}/company-requirement/{companyRequirementID}",
		controller.DeleteCompanyRequirement).Methods(http.MethodDelete)

	// Add mutiple talents to requirement.
	router.HandleFunc("/tenant/{tenantID}/company-requirement/{companyRequirementID}/talent",
		controller.AddTalentsToCompanyRequirement).Methods(http.MethodPost)

	// Exculde routes.
	*exclude = append(*exclude, requirementList)

	controller.log.Info("Company Requirement routes registered..")
}

// AddCompanyRequirement adds new company requirement.
func (controller *CompanyRequirementController) AddCompanyRequirement(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================AddCompanyRequirement call==============================")

	// Create bucket.
	companyRequirement := &cmp.Requirement{}

	// Declare error variable.
	var err error

	// Get params from api.
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the talent variable with given data.
	err = web.UnmarshalJSON(r, companyRequirement)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = companyRequirement.Validate()
	if err != nil {
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Parse and set tenant ID.
	// util.ParseUUID(param[paramTenantID])
	companyRequirement.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credential ID.
	// util.ParseUUID(param[paramCredentialID])
	companyRequirement.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call add service method.
	err = controller.CompanyRequirementService.AddCompanyRequirement(companyRequirement)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Company requirement added successfully")
}

// UpdateCompanyRequirement updates the company requirement.
func (controller *CompanyRequirementController) UpdateCompanyRequirement(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateCompanyRequirement call==============================")

	// Create bucket.
	companyRequirement := &cmp.Requirement{}

	// Get params from api.
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	// Declare erroe variable.
	var err error

	// Parse company requirement.
	err = web.UnmarshalJSON(r, companyRequirement)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = companyRequirement.Validate()
	if err != nil {
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Parse tenant ID.
	// util.ParseUUID(param[paramTenantID])
	companyRequirement.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	// Parse company requirement ID.
	// util.ParseUUID(param[paramCompanyRequirementID])
	companyRequirement.ID, err = parser.GetUUID(paramCompanyRequirementID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requirement id", http.StatusBadRequest))
		return
	}
	// Parse credential ID.
	// util.ParseUUID(param[paramCredentialID])
	companyRequirement.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.CompanyRequirementService.UpdateCompanyRequirement(companyRequirement)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Company requirement updated successfully")
}

// CloseCompanyRequirement updates the company requirement's is active field to false.
func (controller *CompanyRequirementController) CloseCompanyRequirement(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================CloseCompanyRequirement call==============================")

	// Create bucket.
	companyRequirement := cmp.Requirement{}

	// Get params from api.
	// param := mux.Vars(r)/
	parser := web.NewParser(r)

	// Declare erroe variable.
	var err error

	// Parse tenant ID.
	// util.ParseUUID(param[paramTenantID])
	companyRequirement.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	// Parse company requirement ID.
	// util.ParseUUID(param[paramCompanyRequirementID])
	companyRequirement.ID, err = parser.GetUUID(paramCompanyRequirementID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requirement id", http.StatusBadRequest))
		return
	}

	// Parse credential ID.
	// util.ParseUUID(param[paramCredentialID])
	companyRequirement.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.CompanyRequirementService.CloseCompanyRequirement(&companyRequirement)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Company requirement closed successfully")
}

// UpdateCompanyRequirementsSalesperson update one or more company requirement' salesperson.
func (controller *CompanyRequirementController) UpdateCompanyRequirementsSalesperson(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("===============================UpdateCompanyRequirementsSalesperson called=======================================")

	// Create bucket.
	requirements := []cmp.RequirementUpdate{}

	// Get params from api.
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the requirements variable with given data.
	err := web.UnmarshalJSON(r, &requirements)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Getting salesperson id from param and parsing it to uuid.
	// util.ParseUUID(params["salesPersonID"])
	salesPersonID, err := parser.GetUUID("salesPersonID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse salesPerson id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	// util.ParseUUID(params["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	// util.ParseUUID(params["credentialID"])
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.CompanyRequirementService.UpdateCompanyRequirementsSalesperson(&requirements, salesPersonID, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Company requirements' salesperson updated successfully")
}

// DeleteCompanyRequirement deletes company Requirement by id.
func (controller *CompanyRequirementController) DeleteCompanyRequirement(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteCompanyRequirement call==============================")

	// Get params from api.
	// param := mux.Vars(r)
	parser := web.NewParser(r)
	// Decalre error variable.
	var err error

	// Create bucket.
	companyRequirement := &cmp.Requirement{}

	// Parse and set requirement ID.
	// util.ParseUUID(param[paramCompanyRequirementID])
	companyRequirement.ID, err = parser.GetUUID(paramCompanyRequirementID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requirement id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	// util.ParseUUID(param[paramTenantID])
	companyRequirement.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credential ID.
	// util.ParseUUID(param[paramCredentialID])
	companyRequirement.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)

	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Add delete service method.
	err = controller.CompanyRequirementService.DeleteCompanyRequirement(companyRequirement)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Company requirement deleted successfully")
}

// AddTalentsToCompanyRequirement add talents to the company requirement.
func (controller *CompanyRequirementController) AddTalentsToCompanyRequirement(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddTalentsToCompanyRequirement call==============================")

	// Create bucket.
	requirementTalents := []cmp.CompanyRequirementTalents{}

	parser := web.NewParser(r)
	// Getting tenant id from param and parsing it to uuid.
	err := web.UnmarshalJSON(r, &requirementTalents)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	// util.ParseUUID(mux.Vars(r)[paramTenantID])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting requirement id from param and parsing it to uuid.
	// util.ParseUUID(mux.Vars(r)[paramCompanyRequirementID])
	requirementID, err := parser.GetUUID(paramCompanyRequirementID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call add talent to requirement service method.
	err = controller.CompanyRequirementService.AddTalentsToCompanyRequirement(&requirementTalents, tenantID, requirementID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Talents added to company requirement successfully")
}

// func (controller *CompanyRequirementController) UpdateRequirementRating(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================UpdateRequirementRating call==============================")

// 	// Get params from api.
// 	param := mux.Vars(r)

// 	// Create bucket.
// 	companyRequirement := cmp.Requirement{}

// 	err := web.UnmarshalJSON(r, &companyRequirement)
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse requirement", http.StatusBadRequest))
// 		return
// 	}

// 	// Parse and set requirement ID.
// 	companyRequirement.ID, err = util.ParseUUID(param[paramCompanyRequirementID])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse requirement id", http.StatusBadRequest))
// 		return
// 	}

// 	// Parse and set tenant ID.
// 	companyRequirement.TenantID, err = util.ParseUUID(param[paramTenantID])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}

// 	// Parse and set credential ID.
// 	companyRequirement.DeletedBy, err = util.ParseUUID(param[paramCredentialID])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
// 		return
// 	}

// 	// Add delete service method.
// 	err = controller.CompanyRequirementService.UpdateRequirementRating(&companyRequirement)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, "Company requirement rating updated")
// }

// GetAllCompanyRequirements returns all company requirements.
func (controller *CompanyRequirementController) GetAllCompanyRequirements(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllCompanyRequirements call==============================")

	// Craete bucket.
	companyRequirements := &[]cmp.RequirementDTO{}

	parser := web.NewParser(r)
	// Getting tenant id from param and parsing it to uuid.
	// util.ParseUUID(mux.Vars(r)[paramTenantID])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Get total, limit and offset from param and convert it to int.
	var totalCount int
	// limit, offset := web.GetLimitAndOffset(r)

	// Parse form
	// r.ParseForm()

	// Call get all service method.
	err = controller.CompanyRequirementService.GetAllCompanyRequirements(companyRequirements, parser, tenantID, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// web.RespondJSON(w, http.StatusOK, companyRequirements)
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, companyRequirements)
}

// GetMyOpportunities returns specific details of all company requirements.
func (controller *CompanyRequirementController) GetMyOpportunities(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================GetMyOpportunities call==============================")

	// Craete bucket.
	companyRequirements := &[]cmp.MyOpportunityDTO{}

	parser := web.NewParser(r)

	// Getting tenant id from param and parsing it to uuid.
	// util.ParseUUID(mux.Vars(r)[paramTenantID])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Get total, limit and offset from param and convert it to int.
	var totalCount int
	// limit, offset := web.GetLimitAndOffset(r)

	// Parse form
	// r.ParseForm()

	// Call get all service method.
	err = controller.CompanyRequirementService.GetMyOpportunities(companyRequirements, parser, tenantID, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// web.RespondJSON(w, http.StatusOK, companyRequirements)
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, companyRequirements)
}

// GetCompanyDetails will return requirement details for the specified requirement.
func (controller *CompanyRequirementController) GetCompanyDetails(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================GetCompanyDetails call==============================")

	opportunity := cmp.CompanyDetails{}
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	// Parse and set tenant ID.
	// util.ParseUUID(param[paramTenantID])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set requirement ID.
	// util.ParseUUID(param[paramCompanyRequirementID])
	requirementID, err := util.ParseUUID(paramCompanyRequirementID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requirement id", http.StatusBadRequest))
		return
	}

	err = controller.CompanyRequirementService.GetCompanyDetails(tenantID, requirementID, &opportunity)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	web.RespondJSON(w, http.StatusOK, opportunity)
}

// GetCompanyRequirement returns one company requirement by id.
func (controller *CompanyRequirementController) GetCompanyRequirement(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================GetCompanyRequirement call==============================")

	// Create bucket.
	companyRequirement := &cmp.RequirementDTO{}

	// Declare error variable.
	var err error

	// Get params from api.
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	// Parse and set tenant ID.
	// util.ParseUUID(mux.Vars(r)[paramTenantID])
	companyRequirement.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set requirement ID.
	// util.ParseUUID(param[paramCompanyRequirementID])
	companyRequirement.ID, err = parser.GetUUID(paramCompanyRequirementID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}
	// Get service method call.
	err = controller.CompanyRequirementService.GetCompanyRequirement(companyRequirement)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, companyRequirement)
}

// GetRequirementList will return list of all active requirements.
func (controller *CompanyRequirementController) GetRequirementList(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================GetRequirementList call==============================")

	// Create bucket for requirements.
	requirementList := []list.Requirement{}

	// Parse Form
	parser := web.NewParser(r)
	// r.ParseForm()

	// Get tenant id and parse it to uuid.
	// util.ParseUUID(mux.Vars(r)[paramTenantID])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Get requirement list service call.
	err = controller.CompanyRequirementService.GetRequirementList(&requirementList, tenantID, parser)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, requirementList)

}
