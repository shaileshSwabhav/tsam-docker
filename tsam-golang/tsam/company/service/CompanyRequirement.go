package services

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	cmp "github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CompanyRequirementService provides method like Add, Update, Delete, GetByID, GetAll for company's requirement
// associationNames field will contain details about the sub-structs in companyRequirement for preload and other operations.
type CompanyRequirementService struct {
	DB               *gorm.DB
	Repository       repository.Repository
	associationNames []string
}

// NewCompanyRequirementService returns the new instance of CompanyRequirementService
func NewCompanyRequirementService(db *gorm.DB, repo repository.Repository) *CompanyRequirementService {
	return &CompanyRequirementService{
		DB:         db,
		Repository: repo,
		associationNames: []string{
			"Company", "Country", "State", "Technologies", "Universities",
			"Qualifications", "SalesPerson", "Designation",
			// "SelectedTalents",
			//  "SelectedTalents.Experiences", "SelectedTalents.Experiences.Technologies","Colleges"
		},
	}
}

// AddCompanyRequirement add to new companyRequirement to database
func (service *CompanyRequirementService) AddCompanyRequirement(companyRequirement *cmp.Requirement) error {

	// Check if tenant exists.
	err := service.doesTenantExist(companyRequirement.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(companyRequirement.TenantID, companyRequirement.CreatedBy)
	if err != nil {
		return err
	}

	// Check if foreign keys exist.
	err = service.doForeignKeysExist(companyRequirement)
	if err != nil {
		return err
	}

	// Extracts ID's of country and state.
	service.extractID(companyRequirement)

	// Creating credential service.
	// credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// // Check if credential ID has permission to add company requirement.
	// err = credentialService.ValidatePermission(companyRequirement.TenantID, companyRequirement.CreatedBy, "/company/requirement", "add")
	// if err != nil {
	// 	return errors.NewValidationError(err.Error())
	// }

	// Make is active field true.
	companyRequirement.IsActive = true

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create temp company branch.
	tempCompany := cmp.Company{}

	// Get company branch ny requirement's company branch id.
	err = service.Repository.GetRecordForTenant(uow, companyRequirement.TenantID, &tempCompany,
		repository.Filter("`id` = ?", companyRequirement.CompanyID),
		repository.Select("`company_name`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Genearte unique code.
	companyRequirement.Code, err = util.GenerateUniqueCode(uow.DB, tempCompany.CompanyName, "`code` = ?", companyRequirement)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// Add requirement to database.
	err = service.Repository.Add(uow, companyRequirement)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	uow.Commit()
	return nil
}

// UpdateCompanyRequirement updates the data of company requirement.
func (service *CompanyRequirementService) UpdateCompanyRequirement(companyRequirement *cmp.Requirement) error {

	// Check if tenant ID exists.
	err := service.doesTenantExist(companyRequirement.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(companyRequirement.TenantID, companyRequirement.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if company requirement ID exists.
	err = service.doesCompanyRequirementExist(companyRequirement.TenantID, companyRequirement.ID)
	if err != nil {
		return err
	}

	// Check if foreign keys exist.
	err = service.doForeignKeysExist(companyRequirement)
	if err != nil {
		return err
	}

	// Extracts ID's of country and state.
	service.extractID(companyRequirement)

	// Creating credential service.
	// credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// // Checks if credentialID has permission to update companyRequirement.
	// err = credentialService.ValidatePermission(companyRequirement.TenantID, companyRequirement.UpdatedBy,
	// 	"/company/requirement", "update")
	// if err != nil {
	// 	return errors.NewValidationError(err.Error())
	// }

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create temp company requirement.
	tempCompanyRequirement := cmp.Requirement{}

	// Get company requirement by id.
	err = service.Repository.GetRecordForTenant(uow, companyRequirement.TenantID, &tempCompanyRequirement,
		repository.Filter("`id` = ?", companyRequirement.ID), repository.Select([]string{"`created_by`", "`code`", "is_active"}))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Give code from temp to requirement.
	companyRequirement.Code = tempCompanyRequirement.Code

	// Give created_by from temp to requirement.
	companyRequirement.CreatedBy = tempCompanyRequirement.CreatedBy

	// Give is_active from temp to requirement.
	companyRequirement.IsActive = tempCompanyRequirement.IsActive

	// Update company requirement assocaitions.
	err = service.updateCompanyRequirementAssociations(uow, companyRequirement)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Update requirement in database.
	err = service.Repository.Save(uow, companyRequirement)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	uow.Commit()
	return nil
}

// CloseCompanyRequirement updates the company requirement's is active field to false.
func (service *CompanyRequirementService) CloseCompanyRequirement(companyRequirement *cmp.Requirement) error {

	// Check if tenant ID exists.
	err := service.doesTenantExist(companyRequirement.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(companyRequirement.TenantID, companyRequirement.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if company requirement ID exists.
	err = service.doesCompanyRequirementExist(companyRequirement.TenantID, companyRequirement.ID)
	if err != nil {
		return err
	}

	// Creating credential service.
	// credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// // Checks if credentialID has permission to update company requirement.
	// err = credentialService.ValidatePermission(companyRequirement.TenantID, companyRequirement.UpdatedBy,
	// 	"/company/requirement", "update")
	// if err != nil {
	// 	return errors.NewValidationError(err.Error())
	// }

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update for setting isactive and updatedBy fields for requirement.
	err = service.Repository.UpdateWithMap(uow, companyRequirement, map[interface{}]interface{}{
		"UpdatedBy": companyRequirement.UpdatedBy,
		"IsActive":  false,
	})
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Update all the waiting list entries' isActive field to false.
	err = service.Repository.UpdateWithMap(uow, tal.WaitingList{}, map[interface{}]interface{}{
		"UpdatedBy": companyRequirement.UpdatedBy,
		"IsActive":  false,
	},
		repository.Filter("`company_requirement_id`=?", companyRequirement.ID))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	uow.Commit()
	return nil
}

// UpdateCompanyRequirementsSalesperson updates multiple company requirement' salesperson id.
func (service *CompanyRequirementService) UpdateCompanyRequirementsSalesperson(requirements *[]cmp.RequirementUpdate,
	salepersonID, tenantID, credentialID uuid.UUID) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Validate salespersonid (login id).
	if err := service.doesSalesPersonExist(tenantID, salepersonID); err != nil {
		return err
	}

	// Collect all requirement ids in variable.
	var requirementIDs []uuid.UUID
	for _, requirement := range *requirements {
		requirementIDs = append(requirementIDs, requirement.RequirementID)
	}

	// Validate all requirements.
	for _, requirementID := range requirementIDs {
		if err := service.doesCompanyRequirementExist(tenantID, requirementID); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update sales_person_id field of all talents.
	if err := service.Repository.UpdateWithMap(uow, &cmp.Requirement{}, map[string]interface{}{"SalesPersonID": salepersonID},
		repository.Filter("id IN (?)", requirementIDs)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Sales person could not be allocated to company requirements", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteCompanyRequirement deletes the data of company requirement.
func (service *CompanyRequirementService) DeleteCompanyRequirement(companyRequirement *cmp.Requirement) error {

	// Checks if tenant ID exists.
	err := service.doesTenantExist(companyRequirement.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(companyRequirement.TenantID, companyRequirement.DeletedBy)
	if err != nil {
		return err
	}

	// Checks if company requirement ID exists.
	err = service.doesCompanyRequirementExist(companyRequirement.TenantID, companyRequirement.ID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update for setting deletedBy and deletedAt fields for requirement.
	err = service.Repository.UpdateWithMap(uow, companyRequirement, map[interface{}]interface{}{
		"DeletedBy": companyRequirement.DeletedBy,
		"DeletedAt": time.Now(),
	})
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Update all the waiting list entries' isActive field to false.
	err = service.Repository.UpdateWithMap(uow, tal.WaitingList{}, map[interface{}]interface{}{
		"UpdatedBy": companyRequirement.DeletedBy,
		"IsActive":  false,
	},
		repository.Filter("`company_requirement_id`=?", companyRequirement.ID))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	uow.Commit()
	return nil
}

// AddTalentsToCompanyRequirement adds requirementID and talentID to DB
func (service *CompanyRequirementService) AddTalentsToCompanyRequirement(requirementTalents *[]cmp.CompanyRequirementTalents,
	tenantID, requirementID uuid.UUID) error {

	// Check if tenant id exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if requirement id exists.
	err = service.doesCompanyRequirementExist(tenantID, requirementID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Check if talent ID's exist.
	for _, requirementTalent := range *requirementTalents {
		err = service.doesTalentExist(tenantID, requirementTalent.TalentID)
		if err != nil {
			return err
		}

		// Check if talent id and requirement id entry are unique or not.
		err = service.doesTalentExistForRequirement(uow, tenantID, requirementTalent.TalentID, requirementTalent.RequirementID)
		if err != nil {
			uow.RollBack()
			return err
		}

		// Give requirement id to all entries.
		requirementTalent.RequirementID = requirementID

		// Add entry to database.
		err = service.Repository.Add(uow, requirementTalent)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// // UpdateRequirementRating will update the criteria field of requirement table.
// func (service *CompanyRequirementService) UpdateRequirementRating(requirement *cmp.Requirement) error {

// 	// Check if tenant exist.
// 	err := service.doesTenantExist(requirement.TenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// Check if credential exist.
// 	err = service.doesCredentialExist(requirement.TenantID, requirement.UpdatedBy)
// 	if err != nil {
// 		return err
// 	}

// 	// Check if requirement exist.
// 	err = service.doesCompanyRequirementExist(requirement.TenantID, requirement.ID)
// 	if err != nil {
// 		return err
// 	}

// 	if requirement.Rating == nil {
// 		return errors.NewValidationError("Rating cannot be zero")
// 	}

// 	uow := repository.NewUnitOfWork(service.DB, false)

// 	err = service.Repository.Update(uow, requirement)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	uow.Commit()
// 	return nil
// }

// GetAllCompanyRequirements returns all company requirements from database.
func (service *CompanyRequirementService) GetAllCompanyRequirements(companyRequirements *[]cmp.RequirementDTO, parser *web.Parser,
	tenantID uuid.UUID, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()
	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// // Create query precessors fro filter, pagination and preloads.
	// var queryProcessors []repository.QueryProcessor
	// queryProcessors = append(queryProcessors, repository.OrderBy("`is_active` DESC, `required_before`"))
	// queryProcessors = append(queryProcessors, service.addSearchQueriesFromParams(requestForm)...)
	// queryProcessors = append(queryProcessors, repository.PreloadAssociations(service.associationNames),
	// 	repository.Paginate(limit, offset, totalCount))

	// Create query processors fro filter, pagination and preloads.
	var queryProcessors []repository.QueryProcessor
	searchQueryProcessor, err := service.addSearchQueriesFromParams(parser.Form)
	if err != nil {
		return err
	}
	queryProcessors = append(queryProcessors, searchQueryProcessor...)
	queryProcessors = append(queryProcessors, repository.Filter("company_requirements.tenant_id=?", tenantID))
	queryProcessors = append(queryProcessors, repository.OrderBy("company_requirements.is_active DESC, required_before"))
	queryProcessors = append(queryProcessors, repository.PreloadAssociations(service.associationNames),
		repository.Paginate(limit, offset, totalCount))

	// Get all requirements form database.
	err = service.Repository.GetAll(uow, companyRequirements, queryProcessors...)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// If requirements is not present then do not get extra fields.
	if companyRequirements != nil || len(*companyRequirements) != 0 {
		// Range requirements for getting extra fields.
		err = service.getValuesForRequirements(uow, companyRequirements, tenantID)
		if err != nil {
			uow.RollBack()
			return errors.NewValidationError("Record not found")
		}
	}

	// Get salary trends.
	for index := range *companyRequirements {

		if (*companyRequirements)[index].Designation == nil {
			continue
		}

		var averageExp float64
		var tempSalaryTrend general.SalaryTrend
		if (*companyRequirements)[index].MinimumExperience != nil {
			if (*companyRequirements)[index].MaximumExperience != nil {
				averageExp = math.Round((float64(*(*companyRequirements)[index].MinimumExperience) +
					float64(*(*companyRequirements)[index].MaximumExperience)) / 2.0)
			} else {
				averageExp = math.Round(float64(*(*companyRequirements)[index].MinimumExperience))
			}
		} else {
			averageExp = 0
		}

		for _, tech := range (*companyRequirements)[index].Technologies {

			if (*companyRequirements)[index].SalaryTrend != nil {
				break
			}

			exist, err := repository.DoesRecordExist(service.DB, &general.SalaryTrend{},
				repository.Filter(" `technology_id` = ? AND ? BETWEEN `minimum_experience` AND `maximum_experience`", tech.ID, averageExp),
				repository.Filter("`designation_id` = ?", (*companyRequirements)[index].Designation.ID),
				repository.Filter("salary_trends.`tenant_id` = ? AND salary_trends.`deleted_at` IS NULL", tenantID))
			if err != nil {
				return err
			}
			if !exist {
				continue
			}

			err = service.Repository.GetRecord(uow, &tempSalaryTrend,
				repository.Filter("`technology_id` = ? AND ? BETWEEN `minimum_experience` AND `maximum_experience`", tech.ID, averageExp),
				repository.Filter("`designation_id` = ?", (*companyRequirements)[index].Designation.ID),
				repository.Filter("salary_trends.`tenant_id` = ? AND salary_trends.`deleted_at` IS NULL", tenantID),
				repository.PreloadAssociations([]string{"Technology"}))
			if err != nil {
				uow.RollBack()
				return errors.NewValidationError("Record not found")
			}
			(*companyRequirements)[index].SalaryTrend = &tempSalaryTrend
		}
	}

	uow.Commit()
	return nil
}

// GetMyOpportunities returns specific details of all company requirements.
func (service *CompanyRequirementService) GetMyOpportunities(companyRequirements *[]cmp.MyOpportunityDTO, parser *web.Parser,
	tenantID uuid.UUID, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()

	// Create query processors fro filter, pagination and preloads.
	var queryProcessors []repository.QueryProcessor
	searchQueryProcessor, err := service.addSearchQueriesFromParams(parser.Form)
	if err != nil {
		return err
	}
	queryProcessors = append(queryProcessors, searchQueryProcessor...)
	queryProcessors = append(queryProcessors, repository.Filter("company_requirements.tenant_id=?", tenantID))
	queryProcessors = append(queryProcessors, repository.OrderBy("company_requirements.is_active DESC, required_before"))
	queryProcessors = append(queryProcessors, repository.PreloadAssociations([]string{"Company", "CompanyBranch.Company",
		"Designation"}),
		repository.Paginate(limit, offset, totalCount))

	// Get all requirements form database.
	err = service.Repository.GetAll(uow, companyRequirements, queryProcessors...)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	uow.Commit()
	return nil
}

// GetCompanyRequirement returns particular companyRequirement by ID.
func (service *CompanyRequirementService) GetCompanyRequirement(companyRequirement *cmp.RequirementDTO) error {

	// Check if tenant exists.
	err := service.doesTenantExist(companyRequirement.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get one requirement by id from database.
	err = service.Repository.GetAllForTenant(uow, companyRequirement.TenantID, companyRequirement,
		repository.Filter("`id` = ?", companyRequirement.ID),
		repository.PreloadAssociations(service.associationNames))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	return nil
}

// GetCompanyDetails will return requirement details for the specified requirement.
func (service *CompanyRequirementService) GetCompanyDetails(tenantID, requirementID uuid.UUID,
	opprutunity *cmp.CompanyDetails) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if requirement exist.
	err = service.doesCompanyRequirementExist(tenantID, requirementID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetRecordForTenant(uow, tenantID, opprutunity,
		repository.Filter("`id` = ?", requirementID),
		repository.PreloadAssociations([]string{"CompanyBranch", "CompanyBranch.Company", "Designation"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetRequirementList will return list of all active requirements.
func (service *CompanyRequirementService) GetRequirementList(requirement *[]list.Requirement, tenantID uuid.UUID,
	parser *web.Parser) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors for search critera, filter and preload.
	var queryProcessors []repository.QueryProcessor
	searchQueryProcessor, err := service.addSearchQueriesFromParams(parser.Form)
	if err != nil {
		return err
	}
	queryProcessors = append(queryProcessors, searchQueryProcessor...)
	queryProcessors = append(queryProcessors, repository.PreloadAssociations([]string{"Branch"}))

	// Get requirement list form database.
	err = service.Repository.GetAllForTenant(uow, tenantID, requirement, queryProcessors...)
	if err != nil {
		return err
	}

	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// updateCompanyRequirementAssociations update company requirement's dependencies.
func (service *CompanyRequirementService) updateCompanyRequirementAssociations(uow *repository.UnitOfWork, companyRequirement *cmp.Requirement) error {

	// Replace talents.
	err := service.Repository.ReplaceAssociations(uow, companyRequirement, "SelectedTalents",
		companyRequirement.SelectedTalents)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// Replace qulaifications.
	if err := service.Repository.ReplaceAssociations(uow, companyRequirement, "Qualifications",
		companyRequirement.Qualifications); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// Replace techbologies.
	if err := service.Repository.ReplaceAssociations(uow, companyRequirement, "Technologies",
		companyRequirement.Technologies); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// Replace universities.
	if err := service.Repository.ReplaceAssociations(uow, companyRequirement, "Universities",
		companyRequirement.Universities); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// Make all associations nil to avoid any unnecessary updates or inserts.
	companyRequirement.Qualifications = nil
	companyRequirement.Technologies = nil
	companyRequirement.SelectedTalents = nil
	companyRequirement.Universities = nil

	return nil
}

// doesTenantExist checks if tenant is present or not in database.
func (service *CompanyRequirementService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *CompanyRequirementService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSalesPersonExist checks if salesperson is present or not in database.
func (service *CompanyRequirementService) doesSalesPersonExist(tenantID, salesPersonID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.User{},
		repository.Filter("`id` = ?", salesPersonID))
	if err := util.HandleError("SalesPerson not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDesignationExist checks if designation is present or not in database.
func (service *CompanyRequirementService) doesDesignationExist(tenantID, designationID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Designation{},
		repository.Filter("`id` = ?", designationID))
	if err := util.HandleError("Designation not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCompanyRequirementExist checks if requirement is present or not in database.
func (service *CompanyRequirementService) doesCompanyRequirementExist(tenantID, companyRequirementID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, cmp.Requirement{},
		repository.Filter("`id` = ?", companyRequirementID))
	if err := util.HandleError("Company Requirement not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCompanyBranchExist checks if company branch is present or not in database.
func (service *CompanyRequirementService) doesCompanyExist(tenantID, companyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, cmp.Company{},
		repository.Filter("`id` = ?", companyID))
	if err := util.HandleError("Company not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTalentExist checks if talent is present or not in database.
func (service *CompanyRequirementService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tal.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Talent not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doForeignKeysExist validates all foreign keys of requirement.
func (service *CompanyRequirementService) doForeignKeysExist(requirement *cmp.Requirement) error {

	// Check company branch exists or not.
	err := service.doesCompanyExist(requirement.TenantID, requirement.CompanyID)
	if err != nil {
		return err
	}

	// Check salesperson exists or not.
	if requirement.SalesPersonID != nil {
		err = service.doesSalesPersonExist(requirement.TenantID, *requirement.SalesPersonID)
		if err != nil {
			return err
		}
	}

	err = service.doesDesignationExist(requirement.TenantID, requirement.DesignationID)
	if err != nil {
		return err
	}

	return nil
}

// doesTalentExistForRequirement checks if the talent is already assigned to the requirement.
func (service *CompanyRequirementService) doesTalentExistForRequirement(uow *repository.UnitOfWork,
	tenantID, talentID, requirementID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, cmp.CompanyRequirementTalents{},
		repository.Filter("`requirement_id`=? AND `talent_id`=?", requirementID, talentID))
	if err != nil {
		return err
	}
	if exists {
		tempTalent := tal.Talent{}
		err = service.Repository.GetRecordForTenant(uow, tenantID, &tempTalent,
			repository.Filter("`id`=?", talentID), repository.Select([]string{"`first_name`", "`last_name`"}))
		if err != nil {
			return err
		}
		return errors.NewValidationError(tempTalent.FirstName + " " + tempTalent.LastName + " already exist for this requirement")
	}
	// if err := util.HandleIfExistsError("Talent already assigned to company's requirement", exists, err); err != nil {
	// 	return err
	// }
	return nil
}

// extractID extracts country and state id.
func (service *CompanyRequirementService) extractID(companyRequirement *cmp.Requirement) {

	// Country.
	companyRequirement.CountryID = &companyRequirement.Country.ID
	companyRequirement.Country = nil

	// State.
	companyRequirement.StateID = &companyRequirement.State.ID
	companyRequirement.State = nil
}

// addSearchQueries will add search queries from query params.
func (service *CompanyRequirementService) addSearchQueriesFromParams(requestForm url.Values) ([]repository.QueryProcessor, error) {
	fmt.Println("========================================requestForm ->", requestForm)

	if len(requestForm) == 0 {
		return nil, nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	var queryProcessors []repository.QueryProcessor

	if companyBranchID, ok := requestForm["companyBranchID"]; ok {
		util.AddToSlice("company_branch_id", "= ?", "AND", companyBranchID[0], &columnNames, &conditions, &operators, &values)
	}

	if stateID, ok := requestForm["stateID"]; ok {
		util.AddToSlice("state_id", "= ?", "AND", stateID, &columnNames, &conditions, &operators, &values)
	}
	if country, ok := requestForm["country"]; ok {
		util.AddToSlice("country_id", "= ?", "AND", country, &columnNames, &conditions, &operators, &values)
	}
	if salesPersonID, ok := requestForm["salesPersonID"]; ok {
		id, err := uuid.FromString(salesPersonID[0])
		if err != nil {
			return nil, err
		}
		if util.IsUUIDValid(id) {
			util.AddToSlice("sales_person_id", "= ?", "AND", id, &columnNames, &conditions, &operators, &values)
		} else {
			util.AddToSlice("sales_person_id", "IS NULL", "AND", nil, &columnNames, &conditions, &operators, &values)
		}
	}
	if requirementFromDate, ok := requestForm["requirementFromDate"]; ok {
		util.AddToSlice("required_from", ">= ?", "AND", requirementFromDate, &columnNames, &conditions, &operators, &values)
	}
	if requirementTillDate, ok := requestForm["requirementTillDate"]; ok {
		util.AddToSlice("required_before", "<= ?", "AND", requirementTillDate, &columnNames, &conditions, &operators, &values)
	}
	if city, ok := requestForm["city"]; ok {
		util.AddToSlice("city", "LIKE ?", "AND", "%"+city[0]+"%", &columnNames, &conditions, &operators, &values)
	}

	if personalityType, ok := requestForm["personalityType"]; ok {
		util.AddToSlice("personality_type", "= ?", "AND", personalityType, &columnNames, &conditions, &operators, &values)
	}
	if jobRole, ok := requestForm["jobRole"]; ok {
		util.AddToSlice("job_role", "= ?", "AND", jobRole, &columnNames, &conditions, &operators, &values)
	}
	if minimumExperience, ok := requestForm["minimumExperience"]; ok {
		util.AddToSlice("minimum_experience", ">= ?", "AND", minimumExperience, &columnNames, &conditions, &operators, &values)
	}
	if maximumExperience, ok := requestForm["maximumExperience"]; ok {
		util.AddToSlice("maximum_experience", "<= ?", "AND", maximumExperience, &columnNames, &conditions, &operators, &values)
	}
	if minimumPackage, ok := requestForm["minimumPackage"]; ok {
		util.AddToSlice("minimum_package", ">= ?", "AND", minimumPackage, &columnNames, &conditions, &operators, &values)
	}
	if maximumPackage, ok := requestForm["maximumPackage"]; ok {
		util.AddToSlice("maximum_package", "<= ?", "AND", maximumPackage, &columnNames, &conditions, &operators, &values)
	}
	if packageOffered, ok := requestForm["packageOffered"]; ok {
		util.AddToSlice("package_offered", ">= ?", "AND", packageOffered, &columnNames, &conditions, &operators, &values)
	}
	if designation, ok := requestForm["designation"]; ok {
		util.AddToSlice("designation_id", "IN(?)", "AND", designation, &columnNames, &conditions, &operators, &values)
	}
	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("company_requirements.is_active", "=?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}
	if isFresher, ok := requestForm["isFresher"]; ok {
		if isFresher[0] == "true" {
			queryProcessors = append(queryProcessors,
				repository.Filter("company_requirements.minimum_experience IS NULL"))
		}
		if isFresher[0] == "false" {
			queryProcessors = append(queryProcessors,
				repository.Filter("company_requirements.minimum_experience IS NOT NULL"))
		}
	}

	if _, ok := requestForm["talentRating"]; ok {
		talentRating := requestForm.Get("talentRating")
		if talentRating == "Outstanding" {
			util.AddToSlice("talent_rating", ">= ?", "AND", 5, &columnNames, &conditions, &operators, &values)
		}
		if talentRating == "Excellent" {
			util.AddToSlice("talent_rating", "BETWEEN 3 AND 4", "AND", nil, &columnNames, &conditions, &operators, &values)
		}
		if talentRating == "Average" {
			util.AddToSlice("talent_rating", "<= ?", "AND", 2, &columnNames, &conditions, &operators, &values)
		}
		if talentRating == "Unranked" {
			util.AddToSlice("talent_rating", "IS NULL", "AND", nil, &columnNames, &conditions, &operators, &values)
		}
	}

	// If technologies is present then join company_requirements_technologies and company_requirements table.
	if technologies, ok := requestForm["technologies"]; ok {

		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN company_requirements_technologies ON company_requirements.`id` = company_requirements_technologies.`requirement_id`"))
		if technologies[0] == "Other" {
			queryProcessors = append(queryProcessors, repository.Join("INNER JOIN technologies ON technologies.`id` = company_requirements_technologies.`technology_id`"))
			techLanguage := []string{
				"Advance Java", "Dotnet", "Java", "Machine Learning", "Cloud", "Golang",
			}
			util.AddToSlice("technologies.`language`", "NOT IN(?)", "AND",
				techLanguage, &columnNames, &conditions, &operators, &values)
		} else {

			if len(technologies) > 0 {
				util.AddToSlice("company_requirements_technologies.`technology_id`", "IN(?)", "AND", technologies, &columnNames, &conditions, &operators, &values)
			}
		}

	}

	// If qualification is present then join company_requirements_qualifications and company_requirements table.
	if qualifications, ok := requestForm["qualifications"]; ok {
		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN company_requirements_qualifications ON company_requirements.`id` = company_requirements_qualifications.`requirement_id`"))
		if len(qualifications) > 0 {
			util.AddToSlice("company_requirements_qualifications.`degree_id`", "IN(?)", "AND", qualifications, &columnNames, &conditions, &operators, &values)
		}
	}

	// If talent id is present then join waiting list table.
	if talentID, ok := requestForm["talentID"]; ok {
		queryProcessors = append(queryProcessors,
			repository.Join("JOIN waiting_list ON company_requirements.id = waiting_list.company_requirement_id"),
			repository.Filter("company_requirements.deleted_at IS NULL and company_requirements.is_active = 1 "+
				"AND talent_id =? AND waiting_list.deleted_at IS NULL AND waiting_list.is_active = 1", talentID))
	}

	queryProcessors = append(queryProcessors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("company_requirements.`id`"))

	return queryProcessors, nil
}

// getValuesForRequirements gets values for requirements by firing individual query for each requirement.
func (service *CompanyRequirementService) getValuesForRequirements(uow *repository.UnitOfWork, requirements *[]cmp.RequirementDTO, tenantID uuid.UUID) error {

	for index := range *requirements {

		//********************************************APPLICANTS*********************************************************************

		// Create bucket for applicants.
		var waitingListCount uint16 = 0

		// Get all count of applicants form database.
		// If company requirement is active then get only active waiting list entries.
		if (*requirements)[index].IsActive {
			err := service.Repository.GetCountForTenant(uow, tenantID, tal.WaitingList{}, &waitingListCount,
				repository.Filter("company_requirement_id=?", (*requirements)[index].ID),
				repository.Filter("is_active=?", true),
			)
			if err != nil {
				uow.RollBack()
				return err
			}
		} else {
			err := service.Repository.GetCountForTenant(uow, tenantID, tal.WaitingList{}, &waitingListCount,
				repository.Filter("company_requirement_id=?", (*requirements)[index].ID),
			)
			if err != nil {
				uow.RollBack()
				return err
			}
		}

		// If no applicants then dont assign any value to total applicants fields.
		if waitingListCount == 0 {
			continue
		}

		// Give waiting list to requirement.
		(*requirements)[index].TotalApplicants = &waitingListCount
	}

	return nil
}

// // deleteCompanyRequirementAssociation soft deletes all associations of company Requirement
// func (service *CompanyRequirementService) deleteCompanyRequirementAssociation(uow *repository.UnitOfWork, companyRequirement *cmp.Requirement) error {

// 	// fmt.Println("======================================DELETE ASSOCIATIONS==========================================")

// 	// fmt.Println("***talents", companyRequirement.SelectedTalents)
// 	if len(companyRequirement.SelectedTalents) > 0 {
// 		// fmt.Println("Deleting selected talents")
// 		if err := service.Repository.RemoveAssociations(uow, companyRequirement, "SelectedTalents", companyRequirement.SelectedTalents); err != nil {
// 			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// 	// fmt.Println("***qualifications", companyRequirement.Qualifications)
// 	if len(companyRequirement.Qualifications) > 0 {
// 		// fmt.Println("Deleting qualifications")
// 		if err := service.Repository.RemoveAssociations(uow, companyRequirement, "Qualifications", companyRequirement.Qualifications); err != nil {
// 			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// 	// fmt.Println("***colleges", companyRequirement.Colleges)
// 	if len(companyRequirement.Colleges) > 0 {
// 		// fmt.Println("Deleting colleges")
// 		if err := service.Repository.RemoveAssociations(uow, companyRequirement, "Colleges", companyRequirement.Colleges); err != nil {
// 			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// 	// fmt.Println("***colleges", companyRequirement.Colleges)
// 	if len(companyRequirement.Universities) > 0 {
// 		// fmt.Println("Deleting universities")
// 		if err := service.Repository.RemoveAssociations(uow, companyRequirement, "Universities", companyRequirement.Universities); err != nil {
// 			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// 	// fmt.Println("***technologies", companyRequirement.Technologies)
// 	if len(companyRequirement.Technologies) > 0 {
// 		// fmt.Println("Deleting technologies")
// 		if err := service.Repository.RemoveAssociations(uow, companyRequirement, "Technologies", companyRequirement.Technologies); err != nil {
// 			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// 	// fmt.Println("======================================END==========================================")

// 	return nil
// }
