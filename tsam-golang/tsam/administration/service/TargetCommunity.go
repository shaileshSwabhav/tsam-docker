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
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// TargetCommunityService provides methods to do different CRUD operations on target_community table.
type TargetCommunityService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewTargetCommunityService returns a new instance Of TargetCommunityService.
func NewTargetCommunityService(db *gorm.DB, repository repository.Repository) *TargetCommunityService {
	return &TargetCommunityService{
		DB:         db,
		Repository: repository,
	}
}

// TargetCommunityAssociationNames provides preload associations array for target community.
var TargetCommunityAssociationNames []string = []string{
	"Department", "Department.Role", "SalesPerson", "Function", "Courses", "Colleges", "Companies", "Faculty",
}

// AddTargetCommunity will add new target community record to the table.
func (service *TargetCommunityService) AddTargetCommunity(community *admin.TargetCommunity) error {

	// Checks if all foreign keys exist.
	err := service.doForeignKeysExist(community.TenantID, community.CreatedBy, community)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, community)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// UpdateTargetCommunity will update the specified tagert community record in the table.
func (service *TargetCommunityService) UpdateTargetCommunity(community *admin.TargetCommunity) error {

	// Checks if all foreign key exist.
	err := service.doForeignKeysExist(community.TenantID, community.UpdatedBy, community)
	if err != nil {
		return err
	}

	// Check if target community record exist.
	err = service.doesTargetCommunityExist(community.TenantID, community.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update target community associations.
	if err := service.updateTargetCommunityAssociation(uow, community); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Target Community could not be updated", http.StatusInternalServerError)
	}

	err = service.Repository.Update(uow, community)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// UpdateIsAchievedTargetCommunity will update the is target achieved field of specified tagert community record in the table.
func (service *TargetCommunityService) UpdateIsAchievedTargetCommunity(community *admin.TargetCommunityUpdate, tenantID,
	credentialID uuid.UUID) error {

	// Check if tenant record exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}
	// Check if credential record exist.
	err = service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if target community record exist.
	err = service.doesTargetCommunityExist(tenantID, community.TargetCommunityID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update only is target achieved field.
	err = service.Repository.UpdateWithMap(uow, &admin.TargetCommunity{}, map[interface{}]interface{}{
		"IsTargetAchieved": true,
		"UpdatedBy":        credentialID,
	}, repository.Filter("`id`=?", community.TargetCommunityID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteTargetCommunity will delete the specified target community record from the table.
func (service *TargetCommunityService) DeleteTargetCommunity(community *admin.TargetCommunity) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(community.TenantID); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(community.TenantID, community.DeletedBy); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if target community record exists.
	err := service.doesTargetCommunityExist(community.TenantID, community.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, &admin.TargetCommunity{}, map[interface{}]interface{}{
		"DeletedBy": community.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id`=?", community.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTargetCommunities will return all the records from target_community table.
func (service *TargetCommunityService) GetTargetCommunities(communities *[]admin.TargetCommunityDTO, tenantID uuid.UUID,
	form url.Values, limit, offset int, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Query processors for search and get.
	queryProcessors := service.addSearchQueries(form, tenantID)
	queryProcessors = append(queryProcessors, repository.PreloadAssociations(TargetCommunityAssociationNames),
		repository.Paginate(limit, offset, totalCount),
		repository.Filter("target_communities.`tenant_id`=?", tenantID))

	err = service.Repository.GetAllInOrder(uow, communities, "`target_start_date`", queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTargetCommunityList will return all the records from target_community table.
func (service *TargetCommunityService) GetTargetCommunityList(communities *[]admin.TargetCommunity, tenantID uuid.UUID) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, communities, "`target_start_date`")
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTargetCommunity will return specified record from target_community table.
func (service *TargetCommunityService) GetTargetCommunity(community *admin.TargetCommunity, tenantID, communityID uuid.UUID) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if community record exists.
	err = service.doesTargetCommunityExist(tenantID, communityID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetRecordForTenant(uow, tenantID, community,
		repository.Filter("`id`=?", communityID),
		repository.PreloadAssociations(TargetCommunityAssociationNames))
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

// updateTargetCommunityAssociation updates target community's associations.
func (service *TargetCommunityService) updateTargetCommunityAssociation(uow *repository.UnitOfWork, community *admin.TargetCommunity) error {
	// Replace college branches of target community.
	if err := service.Repository.ReplaceAssociations(uow, community, "Colleges",
		community.Colleges); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace company branches of target community.
	if err := service.Repository.ReplaceAssociations(uow, community, "Companies",
		community.Companies); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace courses of target community.
	if err := service.Repository.ReplaceAssociations(uow, community, "Courses",
		community.Courses); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	// Make all association map nil to avoid inserts or updates while updating target community.
	community.Colleges = nil
	community.Companies = nil
	community.Courses = nil

	return nil
}

// doForeignKeysExist checks if all foregin keys exist and if not returns error.
func (service *TargetCommunityService) doForeignKeysExist(tenantID, credentialID uuid.UUID, community *admin.TargetCommunity) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if department exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Department{},
		repository.Filter("`id` = ?", community.DepartmentID))
	if err := util.HandleError("Invalid department ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if function exists or not.
	exists, err = repository.DoesRecordExist(service.DB, general.TargetCommunityFunction{},
		repository.Filter("`id` = ?", community.FunctionID))
	if err := util.HandleError("Invalid function ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if salesperson exists or not.
	exists, err = repository.DoesRecordExist(service.DB, general.User{},
		repository.Filter("`id` = ?", community.SalesPersonID))
	if err := util.HandleError("Invalid salesperson ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if faculty exists or not.
	exists, err = repository.DoesRecordExist(service.DB, faculty.Faculty{},
		repository.Filter("`id` = ?", community.FacultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if colleges exist or not.
	if community.Colleges != nil && len(community.Colleges) > 0 {
		var collegeIDs []uuid.UUID
		for _, tempCollege := range community.Colleges {
			collegeIDs = append(collegeIDs, tempCollege.ID)
		}
		// Get count for collegeIDs.
		var count int = 0
		err := service.Repository.GetCountForTenant(uow, tenantID, college.Branch{}, &count,
			repository.Filter("`id` IN (?)", collegeIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(community.Colleges) {
			log.NewLogger().Error("College ID is invalid")
			return errors.NewValidationError("College ID is invalid")
		}
	}

	// Check if companies exist or not.
	if community.Companies != nil && len(community.Companies) > 0 {
		var companyIDs []uuid.UUID
		for _, tempCompany := range community.Companies {
			companyIDs = append(companyIDs, tempCompany.ID)
		}
		// Get count for companyIDs.
		var count int = 0
		err := service.Repository.GetCountForTenant(uow, tenantID, company.Branch{}, &count,
			repository.Filter("`id` IN (?)", companyIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(community.Companies) {
			log.NewLogger().Error("Company ID is invalid")
			return errors.NewValidationError("Company ID is invalid")
		}
	}

	// Check if courses exist or not.
	if community.Courses != nil && len(community.Courses) > 0 {
		var courseIDs []uuid.UUID
		for _, tempCourse := range community.Courses {
			courseIDs = append(courseIDs, tempCourse.ID)
		}
		fmt.Println(courseIDs)
		// Get count for courseIDs.
		var count int = 0
		err := service.Repository.GetCountForTenant(uow, tenantID, course.Course{}, &count,
			repository.Filter("`id` IN (?)", courseIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(community.Courses) {
			log.NewLogger().Error("Course ID is invalid")
			return errors.NewValidationError("Course ID is invalid")
		}
	}

	return nil
}

// addSearchQueries adds all search queries if any when getAll is called.
func (service *TargetCommunityService) addSearchQueries(searchForm url.Values, tenantID uuid.UUID) []repository.QueryProcessor {
	fmt.Println("=========================In target community search============================", searchForm)

	// Check if there is search criteria given.
	if len(searchForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcesors []repository.QueryProcessor

	// // Get login id from params.
	// loginID, _ := searchForm["loginID"]

	// //get role name form params
	// roleName, _ := searchForm["roleName"]

	// Target type search.
	if _, ok := searchForm["targetType"]; ok {
		util.AddToSlice("`target_type`", "LIKE ?", "AND", "%"+searchForm.Get("targetType")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Student type search.
	if _, ok := searchForm["studentType"]; ok {
		util.AddToSlice("`student_type`", "LIKE ?", "AND", "%"+searchForm.Get("studentType")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Department id search.
	if departmentID, ok := searchForm["departmentID"]; ok {
		util.AddToSlice("`department_id`", "=?", "AND", departmentID, &columnNames, &conditions, &operators, &values)
	}

	// SalesPerson id search.
	if salesPersonID, ok := searchForm["salesPersonID"]; ok {
		util.AddToSlice("`sales_person_id`", "=?", "AND", salesPersonID, &columnNames, &conditions, &operators, &values)
	}

	// Function id search.
	if functionID, ok := searchForm["functionID"]; ok {
		util.AddToSlice("`function_id`", "=?", "AND", functionID, &columnNames, &conditions, &operators, &values)
	}

	// Faculty id search.
	if facultyID, ok := searchForm["facultyID"]; ok {
		util.AddToSlice("`faculty_id`", "=?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
	}

	// Number of batches(greater than or equal to) search.
	if numberOfBatches, ok := searchForm["numberOfBatches"]; ok {
		util.AddToSlice("`number_of_batches`", ">=?", "AND", numberOfBatches, &columnNames, &conditions, &operators, &values)
	}

	// Target Student Count(greater than or equal to) search.
	if targetStudentCount, ok := searchForm["targetStudentCount"]; ok {
		util.AddToSlice("`target_student_count`", ">=?", "AND", targetStudentCount, &columnNames, &conditions, &operators, &values)
	}

	// Hours(greater than or equal to) search.
	if hours, ok := searchForm["hours"]; ok {
		util.AddToSlice("`hours`", ">=?", "AND", hours, &columnNames, &conditions, &operators, &values)
	}

	// Fees(greater than or equal to) search.
	if fees, ok := searchForm["fees"]; ok {
		util.AddToSlice("`fees`", ">=?", "AND", fees, &columnNames, &conditions, &operators, &values)
	}

	// Is Target Achieved  search.
	if isTargetAchieved, ok := searchForm["isTargetAchieved"]; ok {
		util.AddToSlice("`is_target_achieved`", "=?", "AND", isTargetAchieved, &columnNames, &conditions, &operators, &values)
	}

	// Target Start Date from date.
	if startDateFromDate, ok := searchForm["startDateFromDate"]; ok {
		util.AddToSlice("`target_start_date`", ">= ?", "AND", startDateFromDate, &columnNames, &conditions, &operators, &values)
	}

	// Target Start Date to date.
	if startDateToDate, ok := searchForm["startDateToDate"]; ok {
		util.AddToSlice("`target_start_date`", "<= ?", "AND", startDateToDate, &columnNames, &conditions, &operators, &values)
	}

	// Target End Date from date.
	if endDateFromDate, ok := searchForm["endDateFromDate"]; ok {
		util.AddToSlice("`target_end_date`", ">= ?", "AND", endDateFromDate, &columnNames, &conditions, &operators, &values)
	}

	// Target End Date to date.
	if endDateToDate, ok := searchForm["endDateToDate"]; ok {
		util.AddToSlice("`target_end_date`", "<= ?", "AND", endDateToDate, &columnNames, &conditions, &operators, &values)
	}

	// Talent type(greater than or equal to) search.
	if talentType, ok := searchForm["talentType"]; ok {
		util.AddToSlice("`talent_type`", ">=?", "AND", talentType, &columnNames, &conditions, &operators, &values)
	}

	// Minimum Experience Years(greater than or equal to) search.
	if minExperienceYears, ok := searchForm["minExperienceYears"]; ok {
		util.AddToSlice("`min_experience_years`", ">=?", "AND", minExperienceYears, &columnNames, &conditions, &operators, &values)
	}

	// Maximum Experience Years(greater than or equal to) search.
	if maxExperienceYears, ok := searchForm["maxExperienceYears"]; ok {
		util.AddToSlice("`max_experience_years`", ">=?", "AND", maxExperienceYears, &columnNames, &conditions, &operators, &values)
	}

	// Salary(greater than or equal to) search.
	if salary, ok := searchForm["salary"]; ok {
		util.AddToSlice("`salary`", ">=?", "AND", salary, &columnNames, &conditions, &operators, &values)
	}

	// UpSell(greater than or equal to) search.
	if upSell, ok := searchForm["upSell"]; ok {
		util.AddToSlice("`up_sell`", ">=?", "AND", upSell, &columnNames, &conditions, &operators, &values)
	}

	// CrossSell(greater than or equal to) search.
	if crossSell, ok := searchForm["crossSell"]; ok {
		util.AddToSlice("`cross_sell`", ">=?", "AND", crossSell, &columnNames, &conditions, &operators, &values)
	}

	// Referral(greater than or equal to) search.
	if referral, ok := searchForm["referral"]; ok {
		util.AddToSlice("`referral`", ">=?", "AND", referral, &columnNames, &conditions, &operators, &values)
	}

	// Required Talent Rating(greater than or equal to) search.
	if requiredTalentRating, ok := searchForm["requiredTalentRating"]; ok {
		util.AddToSlice("`required_talent_rating`", ">=?", "AND", requiredTalentRating, &columnNames, &conditions, &operators, &values)
	}

	// Required Talent Rating(greater than or equal to) search.
	if rating, ok := searchForm["rating"]; ok {
		util.AddToSlice("`rating`", ">=?", "AND", rating, &columnNames, &conditions, &operators, &values)
	}

	// Colege branch search.
	if collegeIDs, ok := searchForm["collegeIDs"]; ok {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN target_communities_colleges ON target_communities.`id` = target_communities_colleges.`target_community_id`"),
			repository.Join("JOIN college_branches ON college_branches.`id` = target_communities_colleges.`branch_id`"),
			repository.Filter("college_branches.`deleted_at` IS NULL"),
			repository.Filter("college_branches.`tenant_id`=?", tenantID))

		util.AddToSlice("target_communities_colleges.`branch_id`", "IN(?)", "AND", collegeIDs,
			&columnNames, &conditions, &operators, &values)
	}

	// Company branch search.
	if companyIDs, ok := searchForm["companyIDs"]; ok {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN target_communities_companies ON target_communities.`id` = target_communities_companies.`target_community_id`"),
			repository.Join("JOIN company_branches ON company_branches.`id` = target_communities_companies.`company_branch_id`"),
			repository.Filter("company_branches.`deleted_at` IS NULL"),
			repository.Filter("company_branches.`tenant_id`=?", tenantID))

		util.AddToSlice("target_communities_companies.`company_branch_id`", "IN(?)", "AND", companyIDs,
			&columnNames, &conditions, &operators, &values)
	}

	// Company branch search.
	if courseIDs, ok := searchForm["courseIDs"]; ok {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN target_communities_courses ON target_communities.`id` = target_communities_courses.`target_community_id`"),
			repository.Join("JOIN courses ON courses.`id` = target_communities_courses.`course_id`"),
			repository.Filter("courses.`deleted_at` IS NULL"),
			repository.Filter("courses.`tenant_id`=?", tenantID))

		util.AddToSlice("target_communities_courses.`course_id`", "IN(?)", "AND", courseIDs,
			&columnNames, &conditions, &operators, &values)
	}

	//group by target community id and add all filters
	queryProcesors = append(queryProcesors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("target_communities.`id`"))
	return queryProcesors
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *TargetCommunityService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *TargetCommunityService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTargetCommunityExist returns error if there is no target community record in table for the given tenant.
func (service *TargetCommunityService) doesTargetCommunityExist(tenantID, communityID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, admin.TargetCommunity{},
		repository.Filter("`id` = ?", communityID))
	if err := util.HandleError("Invalid community ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
