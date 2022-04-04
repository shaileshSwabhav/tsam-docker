package services

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	company "github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

//CompanyEnquiryService :
type CompanyEnquiryService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewCompanyEnquiryService returns new instance of CompanyEnquiryService.
func NewCompanyEnquiryService(db *gorm.DB, repository repository.Repository) *CompanyEnquiryService {
	return &CompanyEnquiryService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Country", "State", "Technologies", "Domains", "SalesPerson",
		},
	}
}

// AddCompanyEnquiry adds new company enquiry to database.
func (service *CompanyEnquiryService) AddCompanyEnquiry(enquiry *company.Enquiry) error {

	// Get credential id from CreatedBy field of enquiry(set in controller).
	credentialID := enquiry.CreatedBy

	// Extract foreign key IDs and remove the object.
	service.extractID(enquiry)

	// Validate tenant id.
	if err := service.doesTenantExist(enquiry.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(enquiry.TenantID, credentialID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(enquiry); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Assign Company Enquiry Code.
	var codeError error
	enquiry.Code, codeError = util.GenerateUniqueCode(uow.DB, enquiry.CompanyName, "`code` = ?", &company.Enquiry{})
	if codeError != nil {
		uow.RollBack()
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}

	// Add enquiry.
	err := service.Repository.Add(uow, enquiry)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	uow.Commit()
	return nil
}

// UpdateCompanyEnquiry updates company enquiry in database.
func (service *CompanyEnquiryService) UpdateCompanyEnquiry(enquiry *company.Enquiry) error {

	// Get credential id from UpdatedBy field of equiry(set in controller).
	credentialID := enquiry.UpdatedBy

	// Extract all foreign key IDs,assign to entityID field and make entity object nil.
	service.extractID(enquiry)

	// Validate tenant id.
	if err := service.doesTenantExist(enquiry.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(enquiry.TenantID, credentialID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(enquiry); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting enquiry already present in database.
	tempEnquiry := company.Enquiry{}

	// Get enquiry for getting created_by field of enquiry from database.
	if err := service.Repository.GetForTenant(uow, enquiry.TenantID, enquiry.ID, &tempEnquiry); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp enquiry to enquiry to be updated.
	enquiry.CreatedBy = tempEnquiry.CreatedBy

	// Give code to enquiry.
	enquiry.Code = tempEnquiry.Code

	// Update enquiry associations.
	if err := service.updateCompanyEnquiryAssociations(uow, enquiry); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Enquiry could not be updated", http.StatusInternalServerError)
	}

	// Update enquiry.
	if err := service.Repository.Save(uow, enquiry); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Enquiry could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteCompanyEnquiry deletes company enquiry from database.
func (service *CompanyEnquiryService) DeleteCompanyEnquiry(enquiry *company.Enquiry) error {
	// Get credential id from DeletedBy field of enquiry(set in controller).
	credentialID := enquiry.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(enquiry.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(enquiry.TenantID, credentialID); err != nil {
		return err
	}

	// Starting new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get enquiry for updating deleted_by field of enquiry.
	if err := service.Repository.GetForTenant(uow, enquiry.TenantID, enquiry.ID, enquiry,
		repository.PreloadAssociations([]string{"Technologies", "Domains"})); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	//  Delete enquiry association from database.
	if err := service.deleteCompanyEnquiryAssociations(uow, enquiry, credentialID); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Failed to delete enquiry", http.StatusInternalServerError)
	}

	// Make technologies and domains nil to avoid any updates or inserts.
	enquiry.Technologies = nil
	enquiry.Domains = nil

	// Update enquiry for updating deleted_by and deleted_at fields of enquiry.
	if err := service.Repository.UpdateWithMap(uow, enquiry, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("tenant_id=?", enquiry.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Enquiry could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// // GetCompanyEnquiry returns all CompanyEnquiries.
// func (service *CompanyEnquiryService) GetCompanyEnquiry(companyEnquiry *company.Enquiry) error {

// 	uow := repository.NewUnitOfWork(service.DB, true)
// 	err := service.Repository.GetForTenant(uow, companyEnquiry.TenantID, companyEnquiry.ID, companyEnquiry,
// 		repository.PreloadAssociations(service.associations))
// 	if err != nil {
// 		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
// 	}
// 	companyEnquiries := []*company.Enquiry{
// 		companyEnquiry,
// 	}
// 	service.correctDateFormat(&companyEnquiries)

// 	return nil
// }

// GetAllEnquiries returns all enquiries in database.
func (service *CompanyEnquiryService) GetAllEnquiries(enquiries *[]company.EnquiryDTO,
	tenantID uuid.UUID, parser *web.Parser, totalCount *int) error {

	var queryProcessors []repository.QueryProcessor

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	queryProcessors = append(queryProcessors, service.addSearchQueries(tenantID, parser.Form)...)
	queryProcessors = append(queryProcessors, repository.Filter("company_enquiries.`tenant_id`=?", tenantID),
		repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount))

	// Get enquiries from database.
	if err := service.Repository.GetAllInOrder(uow, enquiries, "company_name", queryProcessors...); err != nil {
		uow.RollBack()
		return err
	}

	return nil
}

// UpdateCompanyEnquirysSalesperson multiple enquiry' salesperson id.
func (service *CompanyEnquiryService) UpdateCompanyEnquirysSalesperson(enquiries *[]company.EnquiryUpdate, salepersonID uuid.UUID,
	tenantID uuid.UUID, credentialID uuid.UUID) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Validate salespersonid (login id).
	if err := service.doesSalesPersonExist(salepersonID, tenantID); err != nil {
		return err
	}

	// Collect all enquiry ids in variable.
	var enquiryIDs []uuid.UUID
	for _, enquiry := range *enquiries {
		enquiryIDs = append(enquiryIDs, enquiry.EnquiryID)
	}

	// Validate all enquiries.
	for _, enquiryID := range enquiryIDs {
		if err := service.doesEnquiryExist(tenantID, enquiryID); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update sales_person_id field of all enquiries.
	if err := service.Repository.UpdateWithMap(uow, &company.Enquiry{}, map[string]interface{}{"SalesPersonID": salepersonID},
		repository.Filter("id IN (?)", enquiryIDs)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Sales person could not be allocated to enquiries", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// updateCompanyEnquiryAssociations updates enquiry's associations (sub-structs)
func (service *CompanyEnquiryService) updateCompanyEnquiryAssociations(uow *repository.UnitOfWork, enquiry *company.Enquiry) error {
	// Replace technologies of enquiry.
	if err := service.Repository.ReplaceAssociations(uow, enquiry, "Technologies",
		enquiry.Technologies); err != nil {
		return err
	}

	// Replace domains of enquiry.
	if err := service.Repository.ReplaceAssociations(uow, enquiry, "Domains",
		enquiry.Domains); err != nil {
		return err
	}

	// Make technologies and domains technologies nil so that it is not inserted again.
	enquiry.Technologies = nil
	enquiry.Domains = nil

	return nil
}

// deleteCompanyEnquiryAssociations deletes enquiry's associations (sub-structs).
func (service *CompanyEnquiryService) deleteCompanyEnquiryAssociations(uow *repository.UnitOfWork,
	enquiry *company.Enquiry, credentialID uuid.UUID) error {

	//***********************************************Deleting call records************************************************
	if err := service.Repository.UpdateWithMap(uow, &company.CallRecord{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("enquiry_id=? AND tenant_id=?", enquiry.ID, enquiry.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Enquiry could not be deleted", http.StatusInternalServerError)
	}

	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *CompanyEnquiryService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credenital record in table.
func (service *CompanyEnquiryService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credenital ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesEnquiryExist returns error if there is no enquiry record in table.
func (service *CompanyEnquiryService) doesEnquiryExist(tenantID, enquiryID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, company.Enquiry{},
		repository.Filter("`id` = ?", enquiryID))
	if err := util.HandleError("Invalid enquiry ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSalesPersonExist returns error if there is no salesperson record in table.
func (service *CompanyEnquiryService) doesSalesPersonExist(salesPersonID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.User{},
		repository.Join("left join roles ON users.role_id = roles.id"),
		repository.Filter("users.id=? AND roles.role_name=? AND users.tenant_id=? AND roles.tenant_id=?",
			salesPersonID, "salesperson", tenantID, tenantID))
	if err := util.HandleError("Invalid salesperson ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doForeignKeysExist validates all foreign keys of enquiry.
func (service *CompanyEnquiryService) doForeignKeysExist(enquiry *company.Enquiry) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if salesperson exists or not.
	if enquiry.SalesPersonID != nil {
		if err := service.doesSalesPersonExist(*enquiry.SalesPersonID, enquiry.TenantID); err != nil {
			return err
		}
	}

	// Check if company branch exists or not.
	if enquiry.CompanyBranchID != nil {
		exists, err := repository.DoesRecordExistForTenant(service.DB, enquiry.TenantID, company.Branch{},
			repository.Filter("`id`=?", enquiry.CompanyBranchID))
		if err := util.HandleError("Invalid company branch ID", exists, err); err != nil {
			return err
		}
	}

	// Check if technologies exist or not.
	if enquiry.Technologies != nil && len(enquiry.Technologies) != 0 {
		var technologyIds []uuid.UUID
		for _, technology := range enquiry.Technologies {
			technologyIds = append(technologyIds, technology.ID)
		}
		var count int = 0
		err := service.Repository.GetCountForTenant(uow, enquiry.TenantID, general.Technology{}, &count,
			repository.Filter("id IN (?)", technologyIds))
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(enquiry.Technologies) {
			return errors.NewValidationError("Technology ID is invalid")
		}
	}

	// Check if domains exist or not.
	if enquiry.Domains != nil && len(enquiry.Domains) != 0 {
		var domainIDs []uuid.UUID
		for _, domain := range enquiry.Domains {
			domainIDs = append(domainIDs, domain.ID)
		}
		var count int = 0
		err := service.Repository.GetCountForTenant(uow, enquiry.TenantID, company.Domain{}, &count,
			repository.Filter("id IN (?)", domainIDs))
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(enquiry.Domains) {
			return errors.NewValidationError("Domain ID is invalid")
		}
	}

	return nil
}

// Extracts ID from object and removes data from the object.
// this is done so that the foreign key entity records are not updated in their respective tables
// when the college branch entity is being added or updated.
func (service *CompanyEnquiryService) extractID(enquiry *company.Enquiry) error {
	// State field.
	if enquiry.State != nil {
		enquiry.StateID = &enquiry.State.ID
	}

	// Country field.
	if enquiry.Country != nil {
		enquiry.CountryID = &enquiry.Country.ID
	}

	return nil
}

// addSearchQueries will add search queries from query params
func (service *CompanyEnquiryService) addSearchQueries(tenantID uuid.UUID, requestForm url.Values) []repository.QueryProcessor {

	fmt.Println("========================================requestForm ->", requestForm)

	if len(requestForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	var queryProcessors []repository.QueryProcessor

	// Company Name.
	if _, ok := requestForm["companyName"]; ok {
		util.AddToSlice("company_name", "LIKE ?", "AND", "%"+requestForm.Get("companyName")+"%", &columnNames, &conditions, &operators, &values)
	}

	// State ID.
	if stateID, ok := requestForm["stateID"]; ok {
		util.AddToSlice("state_id", "= ?", "AND", stateID, &columnNames, &conditions, &operators, &values)
	}

	// Country ID.
	if countryID, ok := requestForm["countryID"]; ok {
		util.AddToSlice("country_id", "= ?", "AND", countryID, &columnNames, &conditions, &operators, &values)
	}

	// SalesPerson ID.
	if salesPersonID, ok := requestForm["salesPersonID"]; ok {
		util.AddToSlice("sales_person_id", "= ?", "AND", salesPersonID, &columnNames, &conditions, &operators, &values)
	}

	// Enquiry Type.
	if enquiryType, ok := requestForm["enquiryType"]; ok {
		util.AddToSlice("enquiry_type", "= ?", "AND", enquiryType, &columnNames, &conditions, &operators, &values)
	}

	// Enquiry Source.
	if enquirySource, ok := requestForm["enquirySource"]; ok {
		util.AddToSlice("enquiry_source", "= ?", "AND", enquirySource, &columnNames, &conditions, &operators, &values)
	}

	// Enquiry From Date.
	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("enquiry_date", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}

	// Enquiry Till Date.
	if tillDate, ok := requestForm["tillDate"]; ok {
		util.AddToSlice("enquiry_date", "<= ?", "AND", tillDate, &columnNames, &conditions, &operators, &values)
	}

	// City.
	if _, ok := requestForm["city"]; ok {
		util.AddToSlice("city", "LIKE ?", "AND", "%"+requestForm.Get("city")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Technologies.
	if technologyIDs, ok := requestForm["technologyIDs"]; ok {
		if len(technologyIDs) > 0 {
			queryProcessors = append(queryProcessors,
				repository.Join("INNER JOIN company_enquiries_technologies ON company_enquiries.`id` = company_enquiries_technologies.`enquiry_id`"))
			util.AddToSlice("company_enquiries_technologies.`technology_id`", "IN(?)", "AND", technologyIDs, &columnNames, &conditions, &operators, &values)
		}
	}

	// Domains.
	if domainIDs, ok := requestForm["domainIDs"]; ok {
		if len(domainIDs) > 0 {
			queryProcessors = append(queryProcessors,
				repository.Join("INNER JOIN company_enquiries_domains ON company_enquiries.`id` = company_enquiries_domains.`enquiry_id`"))
			util.AddToSlice("company_enquiries_domains.`domain_id`", "IN(?)", "AND", domainIDs, &columnNames, &conditions, &operators, &values)
		}
	}

	// Call records.
	outcomeID, outcomeOk := requestForm["outcomeID"]
	purposeID, purposeOk := requestForm["outcomeID"]
	if outcomeOk || purposeOk {
		queryProcessors = append(queryProcessors,
			repository.Join("INNER JOIN company_enquiry_call_records ON company_enquiries.`id` = company_enquiry_call_records.`company_enquiry_id`"),
			repository.Filter("company_enquiry_call_records.`deleted_at` IS NULL"),
			repository.Filter("company_enquiry_call_records.`tenant_id`=?", tenantID))
		if outcomeOk {
			util.AddToSlice("company_enquiry_call_records.`outcome_id`", "=?", "AND", outcomeID,
				&columnNames, &conditions, &operators, &values)
		}
		if purposeOk {
			util.AddToSlice("company_enquiry_call_records.`purpose_id`", "=?", "AND", purposeID,
				&columnNames, &conditions, &operators, &values)
		}
	}

	//
	// if salesPersonID, ok := requestForm["salesPersonID"]; ok {
	// 	id, _ := uuid.FromString(salesPersonID[0])
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	if util.IsUUIDValid(id) {
	// 		util.AddToSlice("sales_person_id", "= ?", "AND", id, &columnNames, &conditions, &operators, &values)
	// 	} else {
	// 		util.AddToSlice("sales_person_id", "IS NULL", "AND", nil, &columnNames, &conditions, &operators, &values)
	// 	}
	// }

	queryProcessors = append(queryProcessors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("company_enquiries.`id`"))

	return queryProcessors
}
