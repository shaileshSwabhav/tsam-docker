package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	genService "github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	colg "github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// CollegeService provides methods to do Update, Delete, Add, Get operations on College.
// associations field will contain details about the sub-structs in college for preload and other operations.
type CollegeService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewCollegeService returns a new instance of CollegeService.
func NewCollegeService(db *gorm.DB, repository repository.Repository) *CollegeService {
	return &CollegeService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"CollegeBranches", "CollegeBranches.Country", "CollegeBranches.State",
			"CollegeBranches.SalesPerson"},
	}
}

// assignCollegeFieldsToBranch will assign the college fields to branch
func assignCollegeFieldsToBranch(college *colg.College, branch *colg.Branch) {
	branch.CollegeID = college.ID
	branch.TenantID = college.TenantID
	branch.CreatedBy = college.CreatedBy
	branch.UpdatedBy = college.UpdatedBy
	branch.DeletedBy = college.DeletedBy
}

// AddCollege adds a new college to database.
func (service *CollegeService) AddCollege(college *colg.College,
	uows ...*repository.UnitOfWork) error {

	// Check unique fields are unique within JSON.
	err := service.checkFieldUniquenessInJSON(college.CollegeBranches)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	credentialID := college.CreatedBy

	// Check all foreign key records.
	err = service.doForeignKeysExist(credentialID, college)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(college)
	if err != nil {
		return err
	}

	// Starting transaction.
	// Create new unit of work, if no transaction has been passed to the function.
	length := len(uows)
	if length == 0 {
		uows = append(uows, repository.NewUnitOfWork(service.DB, false))
	}
	uow := uows[0]

	// Generate unique code for branch
	college.Code, err = util.GenerateUniqueCode(uow.DB, college.CollegeName, "`code` = ?", &colg.Branch{})
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// Extract branches out of college so that the branches are not added without proper validation.
	collegeBranches := college.CollegeBranches
	college.CollegeBranches = nil

	// Adding college without branches.
	err = service.Repository.Add(uow, college)
	if err != nil {
		// Rollback only if no transaction is passed.
		if length == 0 {
			uow.RollBack()
		}
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Adding college branches associated with it.
	branchService := NewCollegeBranchService(uow.DB, service.Repository)
	for _, branch := range collegeBranches {

		// Assign college related fields to branch
		assignCollegeFieldsToBranch(college, branch)

		// Call add branch of branch service.
		err := branchService.AddCollegeBranch(branch, uow)
		if err != nil {
			// Rollback only if no transaction is passed.
			if length == 0 {
				uow.RollBack()
			}
			log.NewLogger().Error(err.Error())
			return err
		}
	}
	// Add credential only if affiliated with swabhav*******************************************************************

	// //add credential
	// err = service.AddCredential(uow, college)
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }

	// Commit only if no transaction is passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddColleges adds multiple colleges to Database.
func (service *CollegeService) AddColleges(colleges []*colg.College) error {

	// Add one college record at a time.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, college := range colleges {
		err := service.AddCollege(college, uow)
		if err != nil {
			uow.RollBack()
			return errors.NewValidationError("name-" + college.CollegeName + ": " + err.Error())
		}
	}
	// Commit only if all colleges have been added.
	uow.Commit()
	return nil

}

// UpdateCollege updates the college along with it's associations in database.
func (service *CollegeService) UpdateCollege(college *colg.College) error {

	credentialID := college.UpdatedBy

	// check all foreign key records.
	err := service.doForeignKeysExist(credentialID, college)
	if err != nil {
		return err
	}

	// Check if college exists for college.
	err = service.doesCollegeExist(college.TenantID, college.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(college)
	if err != nil {
		return err
	}

	// Transaction for update
	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, college, map[string]interface{}{
		"CollegeName":     college.CollegeName,
		"ChairmanName":    college.ChairmanName,
		"ChairmanContact": college.ChairmanContact,
		"UpdatedBy":       college.UpdatedBy,
	})
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// DeleteCollege deletes college from database.
func (service *CollegeService) DeleteCollege(college *colg.College) error {

	tenantID := college.TenantID
	credentialID := college.DeletedBy

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
	err = service.doesCollegeExist(tenantID, college.ID)
	if err != nil {
		return err
	}

	// Starting transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update the deleted_by and deleted_at field of the record.
	// Deleting the branch
	err = service.Repository.UpdateWithMap(uow, college, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	})
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Unable to delete college", http.StatusBadRequest)
	}

	// Delete college branches.
	err = service.deleteCollegeBranches(uow, credentialID, college.ID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// GetCollege returns a specific college.
func (service *CollegeService) GetCollege(college *colg.College) error {

	tenantID := college.TenantID
	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if college exists
	err = service.doesCollegeExist(tenantID, college.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetRecordForTenant(uow, tenantID, college,
		repository.PreloadAssociations(service.associations))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil
}

// GetAllColleges returns all colleges in database.
func (service *CollegeService) GetAllColleges(tenantID uuid.UUID, colleges *[]colg.College, form url.Values, limit, offset int, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, colleges, "college_name",
		service.addSearchQueries(form),
		repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil
}

// ===========================================================================================================================================================
// Below are private methods used for various operations related to college/branch
// ===========================================================================================================================================================

// // updateCollegeAssociations updates college's associations (sub-structs)
// func (service *CollegeService) updateCollegeAssociations(uow *repository.UnitOfWork, college *colg.College) error {

// 	if err := service.Repository.ReplaceAssociations(uow, college, "CollegeBranches",
// 		college.CollegeBranches); err != nil {
// 		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
// 	}
// 	if err := service.Repository.BatchUpdate(uow, colg.Branch{}, "college_id IS NULL", map[string]interface{}{
// 		"DeletedBy": college.UpdatedBy,
// 	}); err != nil {
// 		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
// 	}

// 	if err := service.Repository.Delete(uow, colg.Branch{}, "college_id IS NULL"); err != nil {
// 		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
// 	}
// 	return nil
// }

// deleteCollegeBranches soft deletes all associations of college
func (service *CollegeService) deleteCollegeBranches(uow *repository.UnitOfWork,
	deletedBy, collegeID uuid.UUID) error {

	// Soft delete college branches
	// First update the deleted_by field of record.
	err := service.Repository.UpdateWithMap(uow, colg.Branch{}, map[string]interface{}{
		"DeletedBy": deletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`college_id`=?", collegeID))
	if err != nil {
		return err
	}
	return nil
}

// returns error if fields are not unique within the JSON.
func (service *CollegeService) checkFieldUniquenessInJSON(collegeBranches []*colg.Branch) error {

	totalBranches := len(collegeBranches)

	// Map to store all TPO emails in JSON
	tpoEmailMap := make(map[string]uint, totalBranches)

	// Map to store all college branchemails in JSON
	emailMap := make(map[string]uint, totalBranches)

	// Map to store all college branchemails in JSON
	branchNameMap := make(map[string]uint, totalBranches)

	// check to see no values of unique fields are repeated in JSON
	for _, branch := range collegeBranches {
		branchNameMap[branch.BranchName] = branchNameMap[branch.BranchName] + 1
		if branchNameMap[branch.BranchName] > 1 {
			return errors.NewHTTPError("Same branch name given for more than 1 branch", http.StatusBadRequest)
		}
		if branch.Email == nil && branch.TPOEmail == nil {
			continue
		}
		if branch.TPOEmail != nil {
			tpoEmailMap[*branch.TPOEmail] = tpoEmailMap[*branch.TPOEmail] + 1
			if tpoEmailMap[*branch.TPOEmail] > 1 {
				return errors.NewHTTPError("Same TPO email given for more than 1 branch", http.StatusBadRequest)
			}
		}
		if branch.Email != nil {
			emailMap[*branch.Email] = emailMap[*branch.Email] + 1
			if emailMap[*branch.Email] > 1 {
				return errors.NewHTTPError("Same email given for more than 1 branch", http.StatusBadRequest)
			}
		}
	}
	return nil
}

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *CollegeService) doForeignKeysExist(credentialID uuid.UUID, college *colg.College) error {
	tenantID := college.TenantID

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
func (service *CollegeService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *CollegeService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no college record for the given tenant.
func (service *CollegeService) doesCollegeExist(tenantID, collegeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, colg.College{},
		repository.Filter("`id` = ?", collegeID))
	if err := util.HandleError("Invalid college ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// validateFieldUniqueness will check the table for any repitition for unique fields
func (service *CollegeService) validateFieldUniqueness(college *colg.College) error {
	// return error if any record has the same college name in DB.
	exists, err := repository.DoesRecordExistForTenant(service.DB, college.TenantID, colg.College{},
		repository.Filter("`college_name`=? AND `id`!= ?", college.CollegeName, college.ID))
	if err := util.HandleIfExistsError("Record already exists with the same college name.", exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

// Need to test properly.
func (service *CollegeService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	if collegeName, ok := requestForm["collegeName"]; ok {
		return repository.Filter("`college_name` LIKE ?", "%"+collegeName[0]+"%")
	}
	return nil
}

// GetCollegeList returns listing of all the colleges
func (service *CollegeService) GetCollegeList(colleges *[]list.College, tenantID uuid.UUID) error {

	// Check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	uow := repository.NewUnitOfWork(service.DB, true)
	// get all colleges ordered by name order
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, colleges, "college_name",
		repository.Filter("`deleted_at` IS NULL"))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}

	return nil
}

//AddCredential adds credential to credential table according to colege details and generates password
func (ser *CollegeService) AddCredential(uow *repository.UnitOfWork, college *colg.College) error {
	//create bucket for credential
	credential := general.Credential{}

	//create bucket for role
	role := general.Role{}

	//get role id by role name as 'talent'
	if err := ser.Repository.GetAllForTenant(uow, college.TenantID, &role, repository.Filter("role_name=?", "college")); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Record not found", http.StatusInternalServerError)
	}

	//generate random password
	password := util.GeneratePassword()

	//give college details to credential
	credential.FirstName = college.CollegeName
	// credential.Email = talent.Email
	credential.Password = password
	// credential.Contact = *college.ChairmanContact
	credential.RoleID = role.ID
	credential.TenantID = college.TenantID
	credential.TalentID = &college.ID
	credential.CreatedBy = college.CreatedBy

	//add credential to database
	loginService := genService.NewCredentialService(ser.DB, ser.Repository)
	if err := loginService.AddCredential(&credential, uow); err != nil {
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError(err.Error())
		}
		return err
	}
	return nil
}
