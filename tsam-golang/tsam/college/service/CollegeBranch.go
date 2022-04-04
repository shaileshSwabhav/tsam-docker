package service

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/college"
	colg "github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// CollegeBranchService provides methods to do Update, Delete, Add, Get operations on CollegeBranch.
// associations field will contain details about the sub-structs in branch for preload and other operations.
type CollegeBranchService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewCollegeBranchService returns new instance of service.
func NewCollegeBranchService(db *gorm.DB, repository repository.Repository) *CollegeBranchService {
	return &CollegeBranchService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Country", "State", "SalesPerson", "University",
		},
	}
}

// AddCollegeBranch adds new college branch to database.
func (service *CollegeBranchService) AddCollegeBranch(branch *college.Branch, uows ...*repository.UnitOfWork) error {

	credentialID := branch.CreatedBy
	err := service.getIDFromName(branch)
	if err != nil {
		return err
	}

	// Extract foreign key IDs and remove the object.
	err = service.extractID(branch)
	if err != nil {
		return err
	}

	// Check all foreign key records.
	err = service.doForeignKeysExist(credentialID, branch)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(branch)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	length := len(uows)
	if length == 0 {
		uows = append(uows, repository.NewUnitOfWork(service.DB, false))
	}
	uow := uows[0]

	// Generate unique code for branch
	branch.Code, err = util.GenerateUniqueCode(uow.DB, branch.BranchName, "`code` = ?", &colg.Branch{})
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// Add repo call.
	err = service.Repository.Add(uow, branch)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// UpdateCollegeBranch updates college branch in database.
func (service *CollegeBranchService) UpdateCollegeBranch(branch *college.Branch) error {

	credentialID := branch.UpdatedBy

	// Extract all foreign key IDs,assign to entityID field and make entity object nil.
	err := service.extractID(branch)
	if err != nil {
		return err
	}

	// check all foreign key records.
	err = service.doForeignKeysExist(credentialID, branch)
	if err != nil {
		return err
	}

	// Check if branch exists for college.
	err = service.doesBranchExistInCollege(branch.TenantID, branch.CollegeID, branch.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(branch)
	if err != nil {
		return err
	}

	// Transaction for update
	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Save(uow, branch)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// DeleteCollegeBranch deletes college branch from database.
func (service *CollegeBranchService) DeleteCollegeBranch(branch *colg.Branch) error {

	tenantID := branch.TenantID
	collegeID := branch.CollegeID
	credentialID := branch.DeletedBy

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Check if college exists
	err = service.doesCollegeExist(tenantID, collegeID)
	if err != nil {
		return err
	}

	// Check if branch exists for college.
	err = service.doesBranchExistInCollege(tenantID, collegeID, branch.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update the deleted_by and deleted_at field of the record.
	// Deleting the branch
	err = service.Repository.UpdateWithMap(uow, branch, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	})
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Unable to delete college branch", http.StatusBadRequest)
	}

	// Get count of the total branches associated with the parent college.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, colg.Branch{},
		repository.Filter("`college_id`=?", collegeID))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Unable to get count of college branches.", http.StatusInternalServerError)
	}

	// If no branches exist for the parent college, delete the college.
	if !exists {
		college := colg.College{}
		college.ID = collegeID
		// Update the deleted_by and deleted_at field of the record.
		// Deleting the college
		// Note:- This will also change the updated_at field value
		err = service.Repository.UpdateWithMap(uow, college, map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		})
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to delete parent college", http.StatusBadRequest)
		}
	}

	uow.Commit()
	return nil
}

// GetCollegeBranch returns all CollegeBranches.
func (service *CollegeBranchService) GetCollegeBranch(branch *college.Branch) error {

	tenantID := branch.TenantID
	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if college exists
	err = service.doesCollegeExist(tenantID, branch.CollegeID)
	if err != nil {
		return err
	}

	// Check if branch exists for college.
	err = service.doesBranchExistInCollege(tenantID, branch.CollegeID, branch.ID)
	if err != nil {
		return err
	}

	// Get called with preload.
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetRecordForTenant(uow, tenantID, branch,
		repository.PreloadAssociations(service.associations))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetCollegeBranchList returns all college branch names from database.
func (service *CollegeBranchService) GetCollegeBranchList(tenantID uuid.UUID,
	collegeBranches *[]*list.Branch) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Getting college_id and tenant_id as they will have the UUID.Nil value and won't be omitted.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, collegeBranches, "`branch_name`",
		repository.Filter("`deleted_at` IS NULL"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil
}

// GetCollegeBranchList returns all college branch names from database.
func (service *CollegeBranchService) GetCollegeBranchListWithLimit(tenantID uuid.UUID,
	collegeBranches *[]*list.Branch, requestForm url.Values, limit, offset int, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Getting college_id and tenant_id as they will have the UUID.Nil value and won't be omitted.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, collegeBranches, "`branch_name`",
		repository.Filter("`deleted_at` IS NULL"), service.addSearchQueries(requestForm), repository.Paginate(limit, offset, totalCount))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil
}

// GetAllBranches returns all branches in database.
func (service *CollegeBranchService) GetAllBranches(tenantID uuid.UUID, collegeBranches *[]*college.Branch, form url.Values,
	limit, offset int, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Get after preloading and adding paging limit and offset.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, collegeBranches, "`branch_name`",
		service.addSearchQueries(form),
		repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetAllBranchesForSalesPerson returns all branches in database where the specific sales person has been assigned.
func (service *CollegeBranchService) GetAllBranchesForSalesPerson(tenantID, salesPersonID uuid.UUID, collegeBranches *[]*college.Branch,
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

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, collegeBranches, "`branch_name`",
		repository.PreloadAssociations(service.associations),
		repository.Filter("`sales_person_id`=?", salesPersonID),
		service.addSearchQueries(form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetAllBranchesOfCollege returns all branches of a specific college.
func (service *CollegeBranchService) GetAllBranchesOfCollege(tenantID, collegeID uuid.UUID,
	collegeBranches *[]*college.Branch) error {
	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if college exists
	err = service.doesCollegeExist(tenantID, collegeID)
	if err != nil {
		return err
	}

	// Get all with preload and order by name.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, collegeBranches, "`branch_name`",
		repository.Filter("college_id=?", collegeID),
		repository.PreloadAssociations(service.associations))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

func (service *CollegeBranchService) validateFieldUniqueness(branch *colg.Branch) error {
	tenantID := branch.TenantID

	// return error if any record has the same college name in DB.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, colg.Branch{},
		repository.Filter("`branch_name`=? AND `college_id`=? AND `id`!= ?",
			branch.BranchName, branch.CollegeID, branch.ID))
	if err := util.HandleIfExistsError("Record already exists with the same branch name.", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError(err.Error())
	}

	// return error if any record has the same college name in DB.
	exists, err = repository.DoesRecordExistForTenant(service.DB, tenantID, colg.Branch{},
		repository.Filter("`tpo_email`=? AND `id`!= ?", branch.TPOEmail, branch.ID))
	if err := util.HandleIfExistsError("Record already exists with the same TPO email.", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError(err.Error())
	}

	// return error if any record has the same college name in DB.
	exists, err = repository.DoesRecordExistForTenant(service.DB, tenantID, colg.Branch{},
		repository.Filter("`email`=? AND `id`!= ?", branch.Email, branch.ID))
	if err := util.HandleIfExistsError("Record already exists with the same branch email.", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError(err.Error())
	}
	return nil
}

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *CollegeBranchService) doForeignKeysExist(credentialID uuid.UUID, branch *colg.Branch) error {
	tenantID := branch.TenantID

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Check if university exists.
	if err := service.doesUniversityExist(tenantID, branch.UniversityID); err != nil {
		return err
	}

	// Check if college exists.
	if err := service.doesCollegeExist(tenantID, branch.CollegeID); err != nil {
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
func (service *CollegeBranchService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *CollegeBranchService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no university record in table for the given tenant.
func (service *CollegeBranchService) doesUniversityExist(tenantID, universityID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.University{},
		repository.Filter("`id` = ?", universityID))
	if err := util.HandleError("Invalid university ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no college record for the given tenant.
func (service *CollegeBranchService) doesCollegeExist(tenantID, collegeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, colg.College{},
		repository.Filter("`id`=?", collegeID))
	if err := util.HandleError("Invalid college ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no college record for the given tenant.
func (service *CollegeBranchService) doesBranchExistInCollege(tenantID, collegeID, branchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, colg.Branch{},
		repository.Filter("`college_id`=? AND `id`=?", collegeID, branchID))
	if err := util.HandleError("Invalid college branch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no country record for the given tenant.
func (service *CollegeBranchService) doesCountryExist(tenantID, countryID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Country{},
		repository.Filter("`id`=?", countryID))
	if err := util.HandleError("Invalid country ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no state record for the given tenant.
func (service *CollegeBranchService) doesStateExist(tenantID, countryID, stateID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.State{},
		repository.Filter("`id`=? AND `country_id`=?", stateID, countryID))
	if err := util.HandleError("Invalid state ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Need to add join for roles
// returns error if there is no salesPerson record for the given tenant.
func (service *CollegeBranchService) doesSalesPersonExist(tenantID, salesPersonID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.User{},
		repository.Filter("`id`=?", salesPersonID))
	if err := util.HandleError("Invalid sales person ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

func (service *CollegeBranchService) getIDFromName(branch *colg.Branch) error {
	fmt.Println("=========================branch================================", branch)
	tenantID := branch.TenantID

	// For country
	country := &general.Country{}
	fmt.Println("=========================before error================================", branch.Country)
	countryName := branch.Country.Name
	fmt.Println("=========================after error================================")

	uow := repository.NewUnitOfWork(service.DB, true)

	fmt.Println("=========================name================================", countryName)
	// If only country name has been given, it will get record of that country and assign the ID
	err := service.Repository.GetRecordForTenant(uow, tenantID, country,
		repository.Filter("`name`=?", countryName))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Invalid country name")
	}
	fmt.Println(country.ID)
	branch.Country.ID = country.ID

	// For state
	state := &general.State{}
	stateName := branch.State.Name

	fmt.Println("=========================name================================", stateName)
	// If only state name has been given, it will get record of that state and assign the ID
	err = service.Repository.GetRecordForTenant(uow, tenantID, state,
		repository.Filter("`name`=?", stateName))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Invalid state name")
	}
	fmt.Println(state.ID)
	branch.State.ID = state.ID

	// For university
	university := &general.University{}
	universityName := branch.University.UniversityName

	fmt.Println("=========================name================================", universityName)
	// If only university name has been given, it will get record of that university and assign the ID
	err = service.Repository.GetRecordForTenant(uow, tenantID, university,
		repository.Filter("`university_name`=?", universityName))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Invalid university name")
	}
	fmt.Println(university.ID)
	branch.University.ID = university.ID
	return nil
}

// Extracts ID from object and removes data from the object.
// this is done so that the foreign key entity records are not updated in their respective tables
// when the college branch entity is being added or updated.
func (service *CollegeBranchService) extractID(branch *colg.Branch) error {
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
	if branch.University != nil {
		branch.UniversityID = branch.University.ID
		branch.University = nil
	}
	return nil
}

// adds all search queries if any when getAll is called
// Need to test properly.
func (service *CollegeBranchService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	fmt.Println("=========================In search============================", requestForm)
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if collegeID, ok := requestForm["collegeID"]; ok {
		util.AddToSlice("`college_id`", "= ?", "AND", collegeID, &columnNames, &conditions, &operators, &values)
	}
	if stateID, ok := requestForm["stateID"]; ok {
		util.AddToSlice("`state_id`", "= ?", "AND", stateID, &columnNames, &conditions, &operators, &values)
	}
	if countryID, ok := requestForm["countryID"]; ok {
		util.AddToSlice("`country_id`", "= ?", "AND", countryID, &columnNames, &conditions, &operators, &values)
	}
	if universityID, ok := requestForm["universityID"]; ok {
		util.AddToSlice("`university_id`", "= ?", "AND", universityID, &columnNames, &conditions, &operators, &values)
	}
	if collegeRating, ok := requestForm["collegeRating"]; ok {
		util.AddToSlice("`college_rating`", ">= ?", "AND", collegeRating, &columnNames, &conditions, &operators, &values)
	}
	if allIndiaRanking, ok := requestForm["allIndiaRanking"]; ok {
		util.AddToSlice("`all_india_ranking`", "<= ?", "AND", allIndiaRanking, &columnNames, &conditions, &operators, &values)
	}
	if tpoName, ok := requestForm["tpoName"]; ok {
		util.AddToSlice("`tpo_name`", "LIKE ?", "AND", "%"+tpoName[0]+"%", &columnNames, &conditions, &operators, &values)
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
	// coordinatorID has been renamed to salesPersonID.
	if salesPersonIDs, ok := requestForm["salesPersonID"]; ok {
		salesPersonID, err := util.ParseUUID(salesPersonIDs[0])
		if err != nil {
			log.NewLogger().Error(err)
			return nil
		}
		if salesPersonID == uuid.Nil {
			util.AddToSlice("`sales_person_id`", "IS NULL", "AND", nil, &columnNames, &conditions, &operators, &values)
		} else {
			util.AddToSlice("`sales_person_id`", "= ?", "AND", salesPersonID, &columnNames, &conditions, &operators, &values)
		}
	}
	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
