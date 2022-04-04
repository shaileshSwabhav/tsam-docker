package services

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	model "github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

//CompanyBranchService :
type CompanyBranchService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewCompanyBranchService returns new instance of CompanyBranchService.
func NewCompanyBranchService(db *gorm.DB, repository repository.Repository) *CompanyBranchService {
	return &CompanyBranchService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Country", "State", "Domains", "Technologies", "SalesPerson",
		},
	}
}

// AddCompanyBranch adds new company branch to database.
func (service *CompanyBranchService) AddCompanyBranch(companyBranch *model.Branch, uows ...*repository.UnitOfWork) error {

	credentialID := companyBranch.CreatedBy

	// Extract foreign key IDs and remove the object.
	service.extractID(companyBranch)

	// check all foreign key records
	err := service.doForeignKeysExist(credentialID, companyBranch)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(companyBranch)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	length := len(uows)
	if length == 0 {
		uows = append(uows, repository.NewUnitOfWork(service.DB, false))
	}
	uow := uows[0]

	// //extract domain record
	// for i := 0; i < len(companyBranch.Domains); i++ {
	// 	domain := &model.Domain{}
	// 	uow = repository.NewUnitOfWork(companyBranchService.DB, true)

	// 	err := companyBranchService.Repository.GetRecordForTenant(uow, companyBranch.TenantID, domain,
	// 		repository.Filter("`id`=?", companyBranch.Domains[i].ID))

	// 	if err != nil {
	// 		uow.RollBack()
	// 		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	// 	}
	// 	companyBranch.Domains[i] = domain
	// }

	// // extract technologies here
	// for i := 0; i < len(companyBranch.Technologies); i++ {
	// 	technologies := &general.Technology{}
	// 	uow = repository.NewUnitOfWork(companyBranchService.DB, true)

	// 	err := companyBranchService.Repository.GetRecordForTenant(uow, companyBranch.TenantID, technologies,
	// 		repository.Filter("`id`=?", companyBranch.Technologies[i].ID))

	// 	if err != nil {
	// 		uow.RollBack()
	// 		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	// 	}
	// 	companyBranch.Technologies[i] = technologies
	// }

	// Check if main branch exist.
	if *companyBranch.MainBranch {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, companyBranch.TenantID, &model.Branch{},
			repository.Filter("`company_id`=? AND `main_branch` = '1'", companyBranch.CompanyID))
		if err := util.HandleIfExistsError("Company can have only 1 main branch.", exists, err); err != nil {
			return err
		}
	}

	// Generate unique code for branch
	companyBranch.Code, err = util.GenerateUniqueCode(uow.DB, companyBranch.BranchName, "`code` = ?", &model.Branch{})
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// Add repo call.
	err = service.Repository.Add(uow, companyBranch)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Create login for company branch
	// credentialService := genService.NewCredentialService(companyBranchService.DB, companyBranchService.Repository)
	// err = companyBranchService.addLoginForCompanyBranch(credentialService, companyBranch, uow)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	uow.RollBack()
	// 	return err
	// }

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}

	return nil
}

// AddCompanyBranches add multiple branches to the DB
func (service *CompanyBranchService) AddCompanyBranches(companyBranches *[]model.Branch, companyBranchesIDs *[]uuid.UUID,
	companyID, tenantID, credentialID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, false)
	for _, companyBranch := range *companyBranches {
		companyBranch.CompanyID = companyID
		companyBranch.TenantID = tenantID
		companyBranch.CreatedBy = credentialID
		err := service.AddCompanyBranch(&companyBranch, uow)
		if err != nil {
			return err
		}
		*companyBranchesIDs = append(*companyBranchesIDs, companyBranch.ID)

	}
	uow.Commit()
	return nil
}

// UpdateCompanyBranch updates company branch in database.
func (service *CompanyBranchService) UpdateCompanyBranch(companyBranch *model.Branch) error {

	credentialID := companyBranch.UpdatedBy

	// Extract all foreign key IDs,assign to entityID field and make entity object nil.
	service.extractID(companyBranch)

	// check all foreign key records
	err := service.doForeignKeysExist(credentialID, companyBranch)
	if err != nil {
		return err
	}

	// check if branch exists for company
	err = service.doesBranchExistInCompany(companyBranch.TenantID, companyBranch.CompanyID, companyBranch.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(companyBranch)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Check if main branch exist.
	if *companyBranch.MainBranch {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, companyBranch.TenantID, &model.Branch{},
			repository.Filter("`company_id`=? AND `main_branch` = '1'", companyBranch.CompanyID))
		if err := util.HandleIfExistsError("Company can have only 1 main branch.", exists, err); err != nil {
			return err
		}
	}

	// get created_by field
	tempBranch := model.Branch{}
	err = service.Repository.GetRecordForTenant(uow, companyBranch.TenantID, &tempBranch,
		repository.Filter("`id` = ?", companyBranch.ID), repository.Select([]string{"created_by"}))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	companyBranch.CreatedBy = tempBranch.CreatedBy

	// updating associations of company branch
	err = service.updateBranchAssociations(uow, companyBranch)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// update branch
	err = service.Repository.Save(uow, companyBranch)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// DeleteCompanyBranch deletes company branch from database.
func (service *CompanyBranchService) DeleteCompanyBranch(companyBranch *model.Branch) error {

	tenantID := companyBranch.TenantID
	companyID := companyBranch.CompanyID
	credentialID := companyBranch.DeletedBy

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Check if company exists
	err = service.doesCompanyExist(tenantID, companyID)
	if err != nil {
		return err
	}

	// Check if branch exists for company.
	err = service.doesBranchExistInCompany(tenantID, companyID, companyBranch.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// // Delete login of company
	// credential := general.Credential{
	// 	CompanyID: &companyBranch.ID,
	// }
	// credentialService := genService.NewCredentialService(companyBranchService.DB, companyBranchService.Repository)
	// err = credentialService.DeleteCredential(&credential, tenantID, companyBranch.DeletedBy, *credential.CompanyID, "company_id=?", uow)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	uow.RollBack()
	// 	return err
	// }

	// get technologie and domains for current branch
	tempBranch := model.Branch{}
	err = service.Repository.GetRecordForTenant(uow, tenantID, &tempBranch,
		repository.Filter("`id` = ?", companyBranch.ID), repository.PreloadAssociations([]string{"Technologies", "Domains"}))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// delete branch associations
	err = service.deleteBranchAssociations(uow, &tempBranch)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Update the deleted_by and deleted_at field of the record.
	// Deleting the branch
	err = service.Repository.UpdateWithMap(uow, &model.Branch{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", companyBranch.ID))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Unable to delete company branch", http.StatusBadRequest)
	}

	// Get count of the total branches associated with the parent company.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, model.Branch{},
		repository.Filter("`company_id`=?", companyID))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Unable to get count of company branches.", http.StatusInternalServerError)
	}

	// If no branches exist for the parent company, delete the company.
	if !exists {
		company := model.Company{}
		company.ID = companyID
		// Update the deleted_by and deleted_at field of the record.
		// Deleting the company
		// Note:- This will also change the updated_at field value
		err = service.Repository.UpdateWithMap(uow, company, map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		})
		if err != nil {
			uow.RollBack()
			return errors.NewHTTPError("Unable to delete parent company", http.StatusBadRequest)
		}
	}

	uow.Commit()
	return nil
}

// GetAllBranches returns all branches in database.
func (service *CompanyBranchService) GetAllBranches(tenantID uuid.UUID, companyBranches *[]*model.Branch,
	parser *web.Parser, totalCount *int) error {

	// check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()

	// Get after preloading and adding paging limit and offset.
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, companyBranches, "`branch_name`",
		repository.PreloadAssociations(service.associations),
		service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil

}

// GetCompanyBranch returns all Company Branches.
func (service *CompanyBranchService) GetCompanyBranch(companyBranch *model.Branch) error {

	tenantID := companyBranch.TenantID

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if company id exists
	err = service.doesCompanyExist(tenantID, companyBranch.CompanyID)
	if err != nil {
		return err
	}

	// Check if branch exists for company.
	err = service.doesBranchExistInCompany(tenantID, companyBranch.CompanyID, companyBranch.ID)
	if err != nil {
		return err
	}

	// Get called with preload.
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetRecordForTenant(uow, tenantID, companyBranch,
		repository.PreloadAssociations(service.associations))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetAllBranchesForSalesPerson returns all branches in database where the specific sales person has been assigned.
func (service *CompanyBranchService) GetAllBranchesForSalesPerson(tenantID, salesPersonID uuid.UUID, companyBranches *[]*model.Branch,
	form url.Values, limit, offset int, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if sales-person exists.
	err = service.doesSalesPersonExist(tenantID, salesPersonID)
	if err != nil {
		return err
	}

	// Get after preloading and adding paging limit and offset.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, companyBranches, "`branch_name`",
		repository.PreloadAssociations(service.associations),
		repository.Filter("`sales_person_id`=?", salesPersonID),
		service.addSearchQueries(form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetAllBranchesOfCompany returns all branches of a specific company.
func (service *CompanyBranchService) GetAllBranchesOfCompany(tenantID uuid.UUID, companyID uuid.UUID,
	companyBranches *[]model.Branch) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if company id exists
	err = service.doesCompanyExist(tenantID, companyID)
	if err != nil {
		return err
	}

	// Get all with preload and order by name.
	uow := repository.NewUnitOfWork(service.DB, true)
	// repository.Filter("`main_branch` = false"),
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, companyBranches, "`branch_name`",
		repository.Filter("company_id=?", companyID),
		repository.PreloadAssociations(service.associations))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil

}

// GetAllCompanyBranchList returns listing of all the company branches
func (service *CompanyBranchService) GetAllCompanyBranchList(branches *[]list.CompanyBranch,
	tenantID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, branches, "`branch_name`",
		repository.Filter("`deleted_at` IS NULL"))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// updateBranchAssociations update companyBranch's dependencies
func (service *CompanyBranchService) updateBranchAssociations(uow *repository.UnitOfWork, branch *model.Branch) error {

	if branch.Technologies != nil {
		err := service.Repository.ReplaceAssociations(uow, branch, "Technologies", branch.Technologies)
		if err != nil {
			return err
		}
	}
	if branch.Domains != nil {
		err := service.Repository.ReplaceAssociations(uow, branch, "Domains", branch.Domains)
		if err != nil {
			return err
		}
	}
	branch.Technologies = nil
	branch.Domains = nil
	return nil
}

// deleteBranchAssociations deletes companyBranch's dependencies
func (service *CompanyBranchService) deleteBranchAssociations(uow *repository.UnitOfWork, branch *model.Branch) error {
	if len(branch.Technologies) > 0 {
		err := service.Repository.RemoveAssociations(uow, branch, "Technologies", branch.Technologies)
		if err != nil {
			return err
		}
	}
	if len(branch.Domains) > 0 {
		err := service.Repository.RemoveAssociations(uow, branch, "Domains", branch.Domains)
		if err != nil {
			return err
		}
	}
	return nil
}

// validateFieldUniqueness checks for duplicate fields in the json
func (service *CompanyBranchService) validateFieldUniqueness(branch *model.Branch) error {

	// err := service.doesHRHeadEmailExist(branch.TenantID, branch.ID, branch.HRHeadEmail)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }

	// err = service.doesTechnologyHeadEmailExist(branch.TenantID, branch.ID, branch.TechnologyHeadEmail)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }

	// err = service.doesUnitHeadEmailExist(branch.TenantID, branch.ID, branch.UnitHeadEmail)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }

	// err = service.doesFinanceHeadEmailExist(branch.TenantID, branch.ID, branch.FinanceHeadEmail)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }

	// err = service.doesRecruitmentHeadEmailExist(branch.TenantID, branch.ID, branch.RecruitmentHeadEmail)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }

	err := service.doesBranchNameExist(branch.TenantID, branch.ID, &branch.BranchName)
	if err != nil {
		return err
	}
	return nil
}

// Extracts ID from object and removes data from the object.
// this is done so that the foreign key entity records are not updated in their respective tables
// when the company branch entity is being added or updated.
func (service *CompanyBranchService) extractID(branch *model.Branch) error {
	if branch.SalesPerson != nil {
		branch.SalesPersonID = &branch.SalesPerson.ID
		branch.SalesPerson = nil
	}
	if branch.State != nil {
		branch.StateID = &branch.State.ID
		branch.State = nil
	}
	if branch.Country != nil {
		branch.CountryID = &branch.Country.ID
		branch.Country = nil
	}

	// // Domain fields
	// if branch.Domains != nil && len(branch.Domains) != 0 {
	// 	for i := 0; i < len(branch.Domains); i++ {
	// 		if branch.Domains[i].DomainName != nil {
	// 			branch.Domains[i].DomainID = &branch.Domains[i].ID
	// 			branch.Domains[i].DomainName = nil
	// 		}
	// 	}
	// }

	// // Technology fields
	// if branch.Technologies != nil && len(branch.Technologies) != 0 {
	// 	for i := 0; i < len(branch.Technologies); i++ {
	// 		branch.Technologies[i].TechnologyID = &branch.Technologies[i].ID
	// 		branch.Technologies[i].Language = ""
	// 	}
	// }

	return nil
}

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *CompanyBranchService) doForeignKeysExist(credentialID uuid.UUID, branch *model.Branch) error {
	tenantID := branch.TenantID

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Check if company exists.
	if err := service.doesCompanyExist(tenantID, branch.CompanyID); err != nil {
		return err
	}

	// Check if country exists.
	if branch.CountryID != nil {
		if err := service.doesCountryExist(tenantID, *branch.CountryID); err != nil {
			return err
		}
	}

	// Check if state exists.
	if branch.CountryID != nil && branch.StateID != nil {
		if err := service.doesStateExist(tenantID, *branch.CountryID, *branch.StateID); err != nil {
			return err
		}
	}

	// Check if sales person exists.
	if branch.SalesPersonID != nil {
		if err := service.doesSalesPersonExist(tenantID, *branch.SalesPersonID); err != nil {
			return err
		}
	}

	return nil
}

// returns error if there is no tenant record in table.
func (service *CompanyBranchService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no company record for the given tenant.
func (companyBranchService *CompanyBranchService) doesBranchExistInCompany(tenantID, companyID, branchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(companyBranchService.DB, tenantID, model.Branch{},
		repository.Filter("`company_id`=? AND `id`=?", companyID, branchID))
	if err := util.HandleError("Invalid company branch ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *CompanyBranchService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no country record for the given tenant.
func (service *CompanyBranchService) doesCountryExist(tenantID, countryID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Country{},
		repository.Filter("`id`=?", countryID))
	if err := util.HandleError("Invalid country ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no company record for the given tenant.
func (service *CompanyBranchService) doesCompanyExist(tenantID, companyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, model.Company{},
		repository.Filter("`id`=?", companyID))
	if err := util.HandleError("Invalid company ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no state record for the given tenant.
func (service *CompanyBranchService) doesStateExist(tenantID, countryID, stateID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.State{},
		repository.Filter("`id`=? AND `country_id`=?", stateID, countryID))
	if err := util.HandleError("Invalid state ID", exists, err); err != nil {
		return err
	}
	return nil
}

// Need to add join for roles
// returns error if there is no salesPerson record for the given tenant.
func (service *CompanyBranchService) doesSalesPersonExist(tenantID, salesPersonID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.User{},
		repository.Filter("`id`=?", salesPersonID))
	if err := util.HandleError("Invalid sales person ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if BranchName exists for the given tenant.
func (service *CompanyBranchService) doesBranchNameExist(tenantID, branchID uuid.UUID, branchName *string) error {
	// return error if any record has the same branch name in DB.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, model.Branch{},
		repository.Filter("`branch_name`=? AND `id`!= ?", branchName, branchID))
	if err := util.HandleIfExistsError("Record already exists with the same branch name.", exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

// adds all search queries if any when getAll is called
// Need to test properly.
func (service *CompanyBranchService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if companyID, ok := requestForm["companyID"]; ok {
		util.AddToSlice("company_id", "= ?", "AND", companyID, &columnNames, &conditions, &operators, &values)
	}
	if stateID, ok := requestForm["stateID"]; ok {
		util.AddToSlice("`state_id`", "= ?", "AND", stateID, &columnNames, &conditions, &operators, &values)
	}
	if countryID, ok := requestForm["countryID"]; ok {
		util.AddToSlice("`country_id`", "= ?", "AND", countryID, &columnNames, &conditions, &operators, &values)
	}
	// coordinatorID has been renaed to salesPersonID.
	if salesPersonID, ok := requestForm["salesPersonID"]; ok {
		util.AddToSlice("`sales_person_id`", "= ?", "AND", salesPersonID, &columnNames, &conditions, &operators, &values)
		// if util.IsUUIDValid(uuid.FromStringOrNil(salesPersonID[0])) {
		// 	util.AddToSlice("sales_person_id", "= ?", "AND", salesPersonID, &columnNames, &conditions, &operators, &values)
		// } else {
		// 	util.AddToSlice("sales_person_id", "IS NULL", "AND", nil, &columnNames, &conditions, &operators, &values)
		// }
	}
	if city, ok := requestForm["city"]; ok {
		util.AddToSlice("`city`", "LIKE ?", "AND", "%"+city[0]+"%", &columnNames, &conditions, &operators, &values)
	}
	if branchName, ok := requestForm["branchName"]; ok {
		util.AddToSlice("`branch_name`", "LIKE ?", "AND", "%"+branchName[0]+"%", &columnNames, &conditions, &operators, &values)
	}
	if code, ok := requestForm["code"]; ok {
		util.AddToSlice("`code`", "LIKE ?", "AND", "%"+code[0]+"%", &columnNames, &conditions, &operators, &values)
	}
	if companyRating, ok := requestForm["companyRating"]; ok {
		util.AddToSlice("company_rating", "= ?", "AND", companyRating, &columnNames, &conditions, &operators, &values)
	}
	if numberOfEmployees, ok := requestForm["numberOfEmployees"]; ok {
		util.AddToSlice("number_of_employees", ">= ?", "AND", numberOfEmployees, &columnNames, &conditions, &operators, &values)
	}
	if hrHeadName, ok := requestForm["hrHeadName"]; ok {
		util.AddToSlice("hr_head_name", "LIKE ?", "AND", "%"+hrHeadName[0]+"%", &columnNames, &conditions, &operators, &values)
	}
	if financeHeadEmail, ok := requestForm["financeHeadEmail"]; ok {
		util.AddToSlice("finance_head_email", "= ?", "OR", financeHeadEmail, &columnNames, &conditions, &operators, &values)
	}
	if technologyHeadEmail, ok := requestForm["technologyHeadEmail"]; ok {
		util.AddToSlice("technology_head_email", "= ?", "OR", technologyHeadEmail, &columnNames, &conditions, &operators, &values)
	}
	if unitHeadEmail, ok := requestForm["unitHeadEmail"]; ok {
		util.AddToSlice("unit_head_email", "= ?", "OR", unitHeadEmail, &columnNames, &conditions, &operators, &values)
	}
	if hrHeadEmail, ok := requestForm["hRHeadEmail"]; ok {
		util.AddToSlice("hr_head_email", "= ?", "OR", hrHeadEmail, &columnNames, &conditions, &operators, &values)
	}
	if recruitmentHeadEmail, ok := requestForm["recruitmentHeadEmail"]; ok {
		util.AddToSlice("recruitment_head_email", "= ?", "OR", recruitmentHeadEmail, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// // returns error if HRHeadEmail record exists for the given tenant.
// func (companyBranchService *CompanyBranchService) doesHRHeadEmailExist(tenantID, branchID uuid.UUID, HRHeadEmail *string) error {
// 	exists, err := repository.DoesRecordExistForTenant(companyBranchService.DB, tenantID, model.Branch{},
// 		repository.Filter("`hr_head_email`=? AND `id`!= ?", HRHeadEmail, branchID))
// 	if err := util.HandleIfExistsError("HR Head Email already exists", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return errors.NewValidationError(err.Error())
// 	}
// 	return nil
// }

// // returns error if TechnologyHeadEmail record exists for the given tenant.
// func (companyBranchService *CompanyBranchService) doesTechnologyHeadEmailExist(tenantID, branchID uuid.UUID, technologyHeadEmail *string) error {
// 	exists, err := repository.DoesRecordExistForTenant(companyBranchService.DB, tenantID, model.Branch{},
// 		repository.Filter("`technology_head_email`=? AND `id`!= ?", technologyHeadEmail, branchID))
// 	if err := util.HandleIfExistsError("Technology Head Email already exists", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return errors.NewValidationError(err.Error())
// 	}
// 	return nil
// }

// // returns error if UnitHeadEmail record exists for the given tenant.
// func (companyBranchService *CompanyBranchService) doesUnitHeadEmailExist(tenantID, branchID uuid.UUID, unitHeadEmail *string) error {
// 	exists, err := repository.DoesRecordExistForTenant(companyBranchService.DB, tenantID, model.Branch{},
// 		repository.Filter("`unit_head_email`=? AND `id`!= ?", unitHeadEmail, branchID))
// 	if err := util.HandleIfExistsError("Unit Head Email already exists", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return errors.NewValidationError(err.Error())
// 	}
// 	return nil
// }

// // returns error if FinanceHeadEmail record exists for the given tenant.
// func (companyBranchService *CompanyBranchService) doesFinanceHeadEmailExist(tenantID, branchID uuid.UUID, financeHeadEmail *string) error {
// 	exists, err := repository.DoesRecordExistForTenant(companyBranchService.DB, tenantID, model.Branch{},
// 		repository.Filter("`finance_head_email`=? AND `id`!= ?", financeHeadEmail, branchID))
// 	if err := util.HandleIfExistsError("Finance Head Email already exists", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return errors.NewValidationError(err.Error())
// 	}
// 	return nil
// }

// // returns error if RecruitmentHeadEmail record exists for the given tenant.
// func (companyBranchService *CompanyBranchService) doesRecruitmentHeadEmailExist(tenantID, branchID uuid.UUID, recruitmentHeadEmail *string) error {
// 	exists, err := repository.DoesRecordExistForTenant(companyBranchService.DB, tenantID, model.Branch{},
// 		repository.Filter("`recruitment_head_email`=? AND `id`!= ?", recruitmentHeadEmail, branchID))
// 	if err := util.HandleIfExistsError("Recruitment Head already exists", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return errors.NewValidationError(err.Error())
// 	}
// 	return nil
// }

// func multipleBranchCheck(branches *[]model.Branch) error {
// 	for _, branch := range *branches {
// 		if *branch.MainBranch {
// 			return errors.NewHTTPError("Company can only have 1 main branch", http.StatusBadRequest)
// 		}
// 	}
// 	return nil
// }

// // addLoginForCompanyBranch creates a record in credentials tables
// func (companyBranchService *CompanyBranchService) addLoginForCompanyBranch(credentialService *genService.CredentialService, companyBranch *model.Branch,
// 	uow *repository.UnitOfWork) error {
// 	// get roleID for faculty
// 	role := general.Role{}
// 	err := companyBranchService.Repository.GetRecordForTenant(uow, companyBranch.TenantID, &role,
// 		repository.Filter("role_name=?", "Company"))
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		uow.RollBack()
// 		return err
// 	}
// 	// err = loginService.ValidatePermission(faculty.TenantID, faculty.CreatedBy, "add")
// 	// if err != nil {
// 	// 	log.NewLogger().Error(err.Error())
// 	// 	return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
// 	// }

// 	// password := "company" // change to util.GeneratePassword()
// 	credentials := general.Credential{
// 		FirstName: companyBranch.BranchName,
// 		Email:     *companyBranch.HRHeadContact,
// 		Password:  "company",
// 		CompanyID: &companyBranch.ID,
// 		Contact:   *companyBranch.HRHeadContact,
// 		RoleID:    role.ID,
// 	}
// 	credentials.TenantID = companyBranch.TenantID
// 	credentials.CreatedBy = companyBranch.CreatedBy
// 	err = credentialService.AddCredential(&credentials, uow)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }
