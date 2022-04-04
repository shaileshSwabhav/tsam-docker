package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	services "github.com/techlabs/swabhav/tsam/company/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	company "github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// CompanyEnquiryController provides methods to do Update, Delete, Add, Get operations on CompanyEnquiry.
type CompanyEnquiryController struct {
	CompanyEnquiryService *services.CompanyEnquiryService
	log                   log.Logger
	auth                  *security.Authentication
}

// NewCompanyEnquiryController creates new instance of CompanyEnquiryController.
func NewCompanyEnquiryController(companyEnquiryService *services.CompanyEnquiryService, log log.Logger, auth *security.Authentication) *CompanyEnquiryController {
	return &CompanyEnquiryController{
		CompanyEnquiryService: companyEnquiryService,
		log:                   log,
		auth:                  auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *CompanyEnquiryController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Add one enquiry.
	router.HandleFunc("/tenant/{tenantID}/company-enquiry",
		controller.AddCompanyEnquiry).Methods(http.MethodPost)

	// Update one enquiry.
	router.HandleFunc("/tenant/{tenantID}/company-enquiry/{enquiryID}",
		controller.UpdateCompanyEnquiry).Methods(http.MethodPut)

	// Delete one enquiry.
	router.HandleFunc("/tenant/{tenantID}/company-enquiry/{enquiryID}",
		controller.DeleteCompanyEnquiry).Methods(http.MethodDelete)
	// router.HandleFunc("tenant/{tenantID}")

	// Get all enquiries.
	router.HandleFunc("/tenant/{tenantID}/company-enquiry",
		controller.GetAllEnquiries).Methods(http.MethodGet)

	// Update all enquiries' salesperson.
	router.HandleFunc("/tenant/{tenantID}/company-enquiry/saleperson/{salesPersonID}",
		controller.UpdateCompanyEnquirysSalesperson).Methods(http.MethodPut)

	// // Get one enquiry.
	// router.HandleFunc("/tenant/{tenantID}/company/enquiry/{enquiryID}",
	// 	controller.GetCompanyEnquiry).Methods(http.MethodGet)

	// router.HandleFunc("/tenant/{tenantID}/company/enquiry/search/limit/{limit}/offset/{offset}",
	// 	controller.GetAllSearchedEnquiries).Methods(http.MethodPost)

	controller.log.Info("Company Enquiry Routes Registered")
}

// GetAllEnquiries returns all Enquiries in database as response
func (controller *CompanyEnquiryController) GetAllEnquiries(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllEnquiries Called==============================")

	//create parser
	parser := web.NewParser(r)

	// Create bucket.
	enquiries := []company.EnquiryDTO{}

	// Limit, offset & totalCount for pagination.
	var totalCount int
	// limit, offset := web.GetLimitAndOffset(r)

	// Crete error variable.
	var err error

	// Getting tenant id from param and parsing it to uuid.
	// util.ParseUUID(mux.Vars(r)[paramTenantID])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse Form
	r.ParseForm()

	// Call get all enquiries service method.
	err = controller.CompanyEnquiryService.GetAllEnquiries(&enquiries, tenantID, parser, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// web.RespondJSON(w, http.StatusOK, companyEnquiries)
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, enquiries)

}

// // GetCompanyEnquiry returns a specific company enquiry as response.
// func (controller *CompanyEnquiryController) GetCompanyEnquiry(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================GetCompanyEnquiry Called==============================")

// 	param := mux.Vars(r)
// 	companyEnquiry := &company.Enquiry{}
// 	var err error

// 	companyEnquiry.ID, err = util.ParseUUID(param[paramCompanyEnquiryID])
// 	if err != nil {
// 		web.RespondError(w, errors.NewValidationError(err.Error()))
// 		return
// 	}

// 	companyEnquiry.TenantID, err = util.ParseUUID(param[paramTenantID])
// 	if err != nil {
// 		web.RespondError(w, errors.NewValidationError(err.Error()))
// 		return
// 	}

// 	err = controller.CompanyEnquiryService.GetCompanyEnquiry(companyEnquiry)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}
// 	// companyEnquiries := []company.Enquiry{
// 	// 	*companyEnquiry,
// 	// }
// 	web.RespondJSON(w, http.StatusOK, companyEnquiry)
// }

// UpdateCompanyEnquiry updates the specific company enquiry
func (controller *CompanyEnquiryController) UpdateCompanyEnquiry(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateCompanyEnquiry Called==============================")

	// Create bucket.
	enquiry := company.Enquiry{}

	// Get params from api.
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the enquiry variable with given data.
	err := web.UnmarshalJSON(r, &enquiry)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = enquiry.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set enquiry ID to enquiry.
	// util.ParseUUID(param[paramCompanyEnquiryID])
	enquiry.ID, err = parser.GetUUID(paramCompanyEnquiryID)
	if err != nil {
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Parse and set tenant ID to enquiry.
	// util.ParseUUID(param[paramTenantID])
	enquiry.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Parse and set update_by to enquiry.
	// util.ParseUUID(param[paramCredentialID])
	enquiry.UpdatedBy, err = parser.GetUUID(paramCredentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Call update service method.
	err = controller.CompanyEnquiryService.UpdateCompanyEnquiry(&enquiry)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Company enquiry updated successfully")
}

// AddCompanyEnquiry adds a new company enquiry.
func (controller *CompanyEnquiryController) AddCompanyEnquiry(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================AddCompanyEnquiry Called==============================")

	// Create bucket.
	enquiry := company.Enquiry{}

	// Get params from api.
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the enquiry variable with given data.
	err := web.UnmarshalJSON(r, &enquiry)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	// util.ParseUUID(param[paramTenantID])
	enquiry.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Parse and set created_by ID.
	// util.ParseUUID(param[paramCredentialID])
	enquiry.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Validates all the fields of company enquiry.
	err = enquiry.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Call add service method.
	err = controller.CompanyEnquiryService.AddCompanyEnquiry(&enquiry)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Company enquiry added successfully")
}

// DeleteCompanyEnquiry deletes the specific Company enquiry
func (controller *CompanyEnquiryController) DeleteCompanyEnquiry(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("==============================DeleteCompanyEnquiry Called==============================")

	// Create bucket.
	enquiry := company.Enquiry{}

	// Get params from api.
	// param := mux.Vars(r)
	parser := web.NewParser(r)

	// Declare err.
	var err error

	// Parse and set enquiry ID.
	// util.ParseUUID(param[paramCompanyEnquiryID])
	enquiry.ID, err = parser.GetUUID(paramCompanyEnquiryID)
	if err != nil {
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Parse and set tenant ID.
	// util.ParseUUID(param[paramTenantID])
	enquiry.TenantID, err = parser.GetTenantID()
	if err != nil {
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Parse and set deleted_by ID.
	// util.ParseUUID(param[paramCredentialID])
	enquiry.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Call delete service method.
	err = controller.CompanyEnquiryService.DeleteCompanyEnquiry(&enquiry)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Company enquiry deleted successfully")
}

// UpdateCompanyEnquirysSalesperson update one or more enquiry' salesperson.
func (controller *CompanyEnquiryController) UpdateCompanyEnquirysSalesperson(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("===============================UpdateCompanyEnquirysSalesperson called=======================================")

	// Create bucket.
	enquiries := []company.EnquiryUpdate{}

	// Get params from api.
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the enquiry variable with given data.
	err := web.UnmarshalJSON(r, &enquiries)
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
	err = controller.CompanyEnquiryService.UpdateCompanyEnquirysSalesperson(&enquiries, salesPersonID, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Company enquiries' salesperson updated successfully")
}
