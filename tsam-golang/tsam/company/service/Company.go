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
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CompanyService provide method like Add, Update, Delete, GetByID, GetAll for company
// companyAssociationNames field will contain details about the sub-structs in company for preload and other operations.
type CompanyService struct {
	DB                      *gorm.DB
	Repository              repository.Repository
	companyAssociationNames []string
}

// NewCompanyService returns the new instance of CompanyService
func NewCompanyService(db *gorm.DB, repo repository.Repository) *CompanyService {
	return &CompanyService{
		DB:         db,
		Repository: repo,
		companyAssociationNames: []string{
			"Branches",
			// Preloading of branch is not required has this data is currently not used.
			// "Branches.State", "Branches.Country", "Branches.Domains", "Branches.Technologies", "Branches.SalesPerson",
		},
	}
}

// assignCompanyFieldsToBranch will assign the company fields to branch
func assignCompanyFieldsToBranch(company *model.Company, branch *model.Branch) {
	branch.CompanyID = company.ID
	branch.TenantID = company.TenantID
	branch.CreatedBy = company.CreatedBy
	branch.UpdatedBy = company.UpdatedBy
	branch.DeletedBy = company.DeletedBy
}

// AddCompany add to new company to database
func (service *CompanyService) AddCompany(company *model.Company) error {

	// Check unique fields are unique within JSON.
	err := service.checkFieldUniquenessInJSON(company.Branches)
	if err != nil {
		return err
	}

	credentialID := company.CreatedBy

	// check all foreign key records
	err = service.doForeignKeysExist(credentialID, company)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(company)
	if err != nil {
		return err
	}

	// Starting transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Generate unique code for company
	company.Code, err = util.GenerateUniqueCode(uow.DB, company.CompanyName, "`code` = ?", &model.Company{})
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// Extract branches out of company so that the branches are not added without proper validation.
	companyBranches := company.Branches
	company.Branches = nil

	// Adding company without branches.
	err = service.Repository.Add(uow, company)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Adding company branches associated with it.
	companyBranchService := NewCompanyBranchService(uow.DB, service.Repository)
	for _, branch := range companyBranches {
		// Assign company related fields to branch
		assignCompanyFieldsToBranch(company, branch)

		// Call add branch of branch service.
		err := companyBranchService.AddCompanyBranch(branch, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// UpdateCompany update the data of company
func (service *CompanyService) UpdateCompany(company *model.Company) error {

	credentialID := company.UpdatedBy

	// check all foreign key records
	err := service.doForeignKeysExist(credentialID, company)
	if err != nil {
		return err
	}

	// Check if company exists or not
	err = service.doesCompanyExist(company.TenantID, company.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(company)
	if err != nil {
		return err
	}

	// Starting transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// err = service.Repository.UpdateWithMap(uow, company, map[string]interface{}{
	// 	"CompanyName": company.CompanyName,
	// 	"Code":        company.Code,
	// 	"Website":     company.Website,
	// 	"UpdatedBy":   company.UpdatedBy,
	// })
	err = service.Repository.Update(uow, company)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// DeleteCompany delete the data of company
func (service *CompanyService) DeleteCompany(company *model.Company) error {

	tenantID := company.TenantID
	credentialID := company.DeletedBy

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Check if college exists.
	err = service.doesCompanyExist(tenantID, company.ID)
	if err != nil {
		return err
	}

	// Starting transaction
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update the deleted_by and deleted_at field of the record.
	err = service.Repository.UpdateWithMap(uow, company, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	})
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Unable to delete company", http.StatusBadRequest)
	}

	// Delete company branch associations
	err = service.deleteCompanyBranchAssociations(uow, tenantID, company.ID)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// Delete company branches.
	err = service.deleteCompanyBranches(uow, credentialID, company.ID)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetAllCompanies returns all company from database
func (service *CompanyService) GetAllCompanies(tenantID uuid.UUID, companies *[]model.CompanyDTO, parser *web.Parser, totalCount *int) error {

	// check if tenant exists or not
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, companies, "company_name",
		service.addSearchQueries(parser.Form), repository.PreloadAssociations(service.companyAssociationNames),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	return nil
}

// GetCompany returns particular company by ID
func (service *CompanyService) GetCompany(company *model.Company) error {

	tenantID := company.TenantID

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if company exists
	err = service.doesCompanyExist(tenantID, company.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetRecordForTenant(uow, tenantID, company,
		repository.PreloadAssociations(service.companyAssociationNames))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil
}

// deleteCompanyBranches soft deletes all associations of company
func (service *CompanyService) deleteCompanyBranches(uow *repository.UnitOfWork,
	deletedBy, companyID uuid.UUID) error {

	// Soft delete company branches
	// First update the deleted_by field of record.
	err := service.Repository.UpdateWithMap(uow, model.Branch{}, map[string]interface{}{
		"DeletedBy": deletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`company_id`=?", companyID))
	if err != nil {
		return err
	}
	return nil
}

// deleteCompanyBranchAssociations deletes companyBranch's dependencies
func (service *CompanyService) deleteCompanyBranchAssociations(uow *repository.UnitOfWork, tenantID, companyID uuid.UUID) error {

	tempBranches := []model.Branch{}
	err := service.Repository.GetAllForTenant(uow, tenantID, &tempBranches,
		repository.Filter("company_id=?", companyID), repository.PreloadAssociations([]string{"Technologies", "Domains"}))
	if err != nil {
		return err
	}

	// fmt.Println("===========================================================================================================")
	// fmt.Println("***companyBranch ->", tempBranches)
	// fmt.Println("===========================================================================================================")

	// fmt.Println("===========================================COMPANY BRANCH ASSOCIATIONS================================================")
	for _, branch := range tempBranches {
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
	}
	// fmt.Println("===========================================================================================================")
	return nil
}

// validateFieldUniqueness will check the table for any repitition for unique fields
func (service *CompanyService) validateFieldUniqueness(company *model.Company) error {
	// return error if any record has the same company name in DB.
	exists, err := repository.DoesRecordExistForTenant(service.DB, company.TenantID, model.Company{},
		repository.Filter("`company_name`=? AND `id`!= ?", company.CompanyName, company.ID))
	if err := util.HandleIfExistsError("Record already exists with the same company name.", exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

// returns error if fields are not unique within the JSON.
func (service *CompanyService) checkFieldUniquenessInJSON(companyBranches []*model.Branch) error {

	totalBranches := len(companyBranches)

	// Map to store all HR Head emails in JSON
	hrHeadEmailMap := make(map[string]uint, totalBranches)

	// Map to store all Technology Head emails in JSON
	technologyHeadEmailMap := make(map[string]uint, totalBranches)

	// Map to store all Unit Head emails in JSON
	unitHeadEmailMap := make(map[string]uint, totalBranches)

	// Map to store all Finance emails in JSON
	financeHeadEmailMap := make(map[string]uint, totalBranches)

	// Map to store all Recruitment emails in JSON
	recruitmentHeadEmailMap := make(map[string]uint, totalBranches)

	// Map to store all company branch in JSON
	branchNameMap := make(map[string]uint, totalBranches)

	// check to see no values of unique fields are repeated in JSON
	for _, branch := range companyBranches {
		branchNameMap[branch.BranchName] = branchNameMap[branch.BranchName] + 1
		if branchNameMap[branch.BranchName] > 1 {
			return errors.NewHTTPError("Same branch name given for more than 1 branch", http.StatusBadRequest)
		}

		if branch.HRHeadEmail == nil && branch.TechnologyHeadEmail == nil && branch.UnitHeadEmail == nil && branch.FinanceHeadEmail == nil && branch.RecruitmentHeadEmail == nil {
			continue
		}
		if branch.HRHeadEmail != nil {
			hrHeadEmailMap[*branch.HRHeadEmail] = hrHeadEmailMap[*branch.HRHeadEmail] + 1
			if hrHeadEmailMap[*branch.HRHeadEmail] > 1 {
				return errors.NewHTTPError("Same HR Head email given for more than 1 branch", http.StatusBadRequest)
			}
		}
		if branch.TechnologyHeadEmail != nil {
			technologyHeadEmailMap[*branch.TechnologyHeadEmail] = technologyHeadEmailMap[*branch.TechnologyHeadEmail] + 1
			if hrHeadEmailMap[*branch.HRHeadEmail] > 1 {
				return errors.NewHTTPError("Same Technology Head email given for more than 1 branch", http.StatusBadRequest)
			}
		}
		if branch.UnitHeadEmail != nil {
			unitHeadEmailMap[*branch.UnitHeadEmail] = unitHeadEmailMap[*branch.UnitHeadEmail] + 1
			if hrHeadEmailMap[*branch.HRHeadEmail] > 1 {
				return errors.NewHTTPError("Same Unit Head email given for more than 1 branch", http.StatusBadRequest)
			}
		}
		if branch.FinanceHeadEmail != nil {
			financeHeadEmailMap[*branch.FinanceHeadEmail] = financeHeadEmailMap[*branch.FinanceHeadEmail] + 1
			if hrHeadEmailMap[*branch.HRHeadEmail] > 1 {
				return errors.NewHTTPError("Same Finance Head email given for more than 1 branch", http.StatusBadRequest)
			}
		}
		if branch.RecruitmentHeadEmail != nil {
			recruitmentHeadEmailMap[*branch.RecruitmentHeadEmail] = recruitmentHeadEmailMap[*branch.RecruitmentHeadEmail] + 1
			if hrHeadEmailMap[*branch.HRHeadEmail] > 1 {
				return errors.NewHTTPError("Same Recruitment Head email given for more than 1 branch", http.StatusBadRequest)
			}
		}
	}
	return nil
}

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *CompanyService) doForeignKeysExist(credentialID uuid.UUID, company *model.Company) error {

	tenantID := company.TenantID

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	return nil
}

// returns error if there is no tenant record in table.
func (service *CompanyService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *CompanyService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no company record for the given tenant.
func (service *CompanyService) doesCompanyExist(tenantID, companyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, model.Company{},
		repository.Filter("`id` = ?", companyID))
	if err := util.HandleError("Invalid company ID", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *CompanyService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	if companyName, ok := requestForm["companyName"]; ok {
		return repository.Filter("`company_name` LIKE ?", "%"+companyName[0]+"%")
	}
	return nil
}
