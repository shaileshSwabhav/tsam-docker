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
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"

	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// CampusDriveService provides methods to do Update, Delete, Add, Get operations on CampusDrive.
// associations field will contain details about the sub-structs in campusDrive for preload and other operations.
type CampusDriveService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewCampusDriveService returns a new instance of CampusDrive.
func NewCampusDriveService(db *gorm.DB, repository repository.Repository) *CampusDriveService {
	return &CampusDriveService{
		DB:         db,
		Repository: repository,
		associations: []string{"CollegeBranches", "Faculties", "SalesPeople", "CompanyRequirements", "CompanyRequirements.Branch",
			"Developers"},
	}
}

// AddCampusDrive adds a new campus drive to database.
func (service *CampusDriveService) AddCampusDrive(campusDrive *college.CampusDrive) error {
	// Get credential id from CreatedBy field of campus drive(set in controller).
	credentialID := campusDrive.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(campusDrive.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, campusDrive.TenantID); err != nil {
		return err
	}

	// Validate foreign key ids.
	if err := service.doForeignKeysExist(campusDrive); err != nil {
		return err
	}

	// Check for same campus drive name.
	if err := service.doesCampusDriveNameExist(campusDrive); err != nil {
		return err
	}

	// Starting transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Assign Campus drive code.
	var codeError error
	campusDrive.Code, codeError = util.GenerateUniqueCode(uow.DB, campusDrive.CampusName, "`code` = ?", &college.CampusDrive{})
	if codeError != nil {
		log.NewLogger().Error(codeError.Error())
		uow.RollBack()
		return errors.NewHTTPError("Internal server err or", http.StatusInternalServerError)
	}

	// Create the registration link and assign it.
	tempRegLink := "https://swabhavtechlabs.com/test/tsm-forms/#/campus-drive-registration?code=" + campusDrive.Code
	campusDrive.StudentRegistrationLink = &tempRegLink

	// Add campus drive to database.
	if err := service.Repository.Add(uow, campusDrive); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Campus Drive could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// UpdateCampusDrive updates the campus drive along with it's associations in database.
func (service *CampusDriveService) UpdateCampusDrive(campusDrive *college.CampusDrive) error {
	// Get credential id from UpdatedBy field of campus drive(set in controller).
	credentialID := campusDrive.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(campusDrive.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, campusDrive.TenantID); err != nil {
		return err
	}

	// Validate foreign key ids.
	if err := service.doForeignKeysExist(campusDrive); err != nil {
		return err
	}

	// Validate campus drive id.
	if err := service.doesCampusDriveExist(campusDrive.ID, campusDrive.TenantID); err != nil {
		return err
	}

	// Check for same campus drive name.
	if err := service.doesCampusDriveNameExist(campusDrive); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting campus drive already present in database.
	tempCampusDrive := college.CampusDrive{}

	// Get campus drive for getting created_by field of campus drive from database.
	if err := service.Repository.GetForTenant(uow, campusDrive.TenantID, campusDrive.ID, &tempCampusDrive,
		repository.PreloadAssociations(service.associations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp campus drive to campus drive to be updated.
	campusDrive.CreatedBy = tempCampusDrive.CreatedBy

	// Give code to campus drive.
	campusDrive.Code = tempCampusDrive.Code

	// Update campus drive associations.
	if err := service.updateCampusDriveAssociations(uow, campusDrive); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Campus Drive could not be updated", http.StatusInternalServerError)
	}

	// Update campus drive.
	if err := service.Repository.Save(uow, campusDrive); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Campus Drive could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteCampusDrive deletes campus drive from database.
func (service *CampusDriveService) DeleteCampusDrive(campusDrive *college.CampusDrive) error {
	// Get credential id from DeletedBy field of campus drive(set in controller).
	credentialID := campusDrive.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(campusDrive.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, campusDrive.TenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get campus drive for updating deleted_by field of campus drive.
	if err := service.Repository.GetForTenant(uow, campusDrive.TenantID, campusDrive.ID, campusDrive); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//  Delete campus drive association from database.
	if err := service.deleteCampusDriveAssociations(uow, campusDrive, credentialID); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Failed to delete campus drive", http.StatusInternalServerError)
	}

	// Update campus drive for updating deleted_by and deleted_at fields of campus drive.
	if err := service.Repository.UpdateWithMap(uow, campusDrive, map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()}); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Campus Drive could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetCampusDrive returns campus drive filtering by ID.
func (service *CampusDriveService) GetCampusDrive(campusDrive *college.CampusDrive) error {
	// Validate tenant id.
	if err := service.doesTenantExist(campusDrive.TenantID); err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get campus drive.
	if err := service.Repository.GetForTenant(uow, campusDrive.TenantID, campusDrive.ID, campusDrive,
		repository.PreloadAssociations(service.associations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetCampusDriveByCode returns campus drive filtering by code.
func (service *CampusDriveService) GetCampusDriveByCode(campusDrive *college.CampusDrive) error {
	// Validate tenant id.
	if err := service.doesTenantExist(campusDrive.TenantID); err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get campus drive.
	if err := service.Repository.GetRecordForTenant(uow, campusDrive.TenantID, campusDrive,
		repository.Filter("`code` = ?", campusDrive.Code),
		repository.PreloadAssociations(service.associations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetAllCampusDrives returns all campus drives in database.
func (service *CampusDriveService) GetAllCampusDrives(campusDrives *[]college.CampusDriveDTO,
	tenantID uuid.UUID, limit int, offset int, totalCount *int, searchForm url.Values) error {
	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Query processors for search and get.
	queryProcessors := service.addSearchQueries(searchForm, tenantID)
	queryProcessors = append(queryProcessors,
		repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount), repository.Filter("campus_drives.`tenant_id`=?", tenantID))

	// Get campus drives from database.
	if err := service.Repository.GetAllInOrder(uow, campusDrives, "`campus_date`", queryProcessors...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Range campus drives for getting total requirements, total registered candidates and total apperaed candidates.
	err := service.getValuesForCampusDrives(uow, campusDrives, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// updateCampusDriveAssociations updates campus drive's associations.
func (service *CampusDriveService) updateCampusDriveAssociations(uow *repository.UnitOfWork, campusDrive *college.CampusDrive) error {
	// Replace salesPeople of campus drive.
	if err := service.Repository.ReplaceAssociations(uow, campusDrive, "SalesPeople",
		campusDrive.SalesPeople); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace faculties of campus drive.
	if err := service.Repository.ReplaceAssociations(uow, campusDrive, "Faculties",
		campusDrive.Faculties); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace developers of campus drive.
	if err := service.Repository.ReplaceAssociations(uow, campusDrive, "Developers",
		campusDrive.Developers); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace company requirements of campus drive.
	if err := service.Repository.ReplaceAssociations(uow, campusDrive, "CompanyRequirements",
		campusDrive.CompanyRequirements); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace college branches of campus drive.
	if err := service.Repository.ReplaceAssociations(uow, campusDrive, "CollegeBranches",
		campusDrive.CollegeBranches); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	campusDrive.SalesPeople = nil
	campusDrive.Faculties = nil
	campusDrive.Developers = nil
	campusDrive.CompanyRequirements = nil
	campusDrive.CollegeBranches = nil

	return nil
}

// deleteCampusDriveAssociations deletes campus drive's associations (sub-structs).
func (service *CampusDriveService) deleteCampusDriveAssociations(uow *repository.UnitOfWork, campusDrive *college.CampusDrive,
	credentialID uuid.UUID) error {

	//**********************************************Delete campus talent registrations***************************************************
	if err := service.Repository.UpdateWithMap(uow, &college.CampusTalentRegistration{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`campus_drive_id`=? AND `tenant_id`=?", campusDrive.ID, campusDrive.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Campus Drive could not be deleted", http.StatusInternalServerError)
	}

	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *CampusDriveService) doesTenantExist(tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if tenant(parent tenant) exists or not.
	exists, err := repository.DoesRecordExist(uow.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCampusDriveExist validates if campus drive exists or not in database.
func (service *CampusDriveService) doesCampusDriveExist(campusDriveID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.CampusDrive{}, repository.Filter("`id` = ?", campusDriveID))
	if err := util.HandleError("Campus drive not found", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credental exists or not in database.
func (service *CampusDriveService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExist validates if foreign keys are present or not in database.
func (service *CampusDriveService) doForeignKeysExist(campusDrive *college.CampusDrive) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if college branches exists or not.
	if campusDrive.CollegeBranches != nil && len(campusDrive.CollegeBranches) != 0 {
		var count int = 0
		var collegeBranchIDs []uuid.UUID
		for _, collegeBranch := range campusDrive.CollegeBranches {
			collegeBranchIDs = append(collegeBranchIDs, collegeBranch.ID)
		}
		err := service.Repository.GetCountForTenant(uow, campusDrive.TenantID, college.Branch{}, &count,
			repository.Filter("`id` IN (?)", collegeBranchIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(campusDrive.CollegeBranches) {
			log.NewLogger().Error("College Branch ID is invalid")
			return errors.NewValidationError("College Branch ID is invalid")
		}
	}

	// Check if company requiremnets exists or not.
	if len(campusDrive.CompanyRequirements) != 0 {
		var count int = 0
		var requirementIDs []uuid.UUID
		for _, requirement := range campusDrive.CompanyRequirements {
			requirementIDs = append(requirementIDs, requirement.ID)
		}
		err := service.Repository.GetCountForTenant(uow, campusDrive.TenantID, company.Requirement{}, &count,
			repository.Filter("`id` IN (?)", requirementIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(campusDrive.CompanyRequirements) {
			log.NewLogger().Error("Company Requirement ID is invalid")
			return errors.NewValidationError("Company Requirement ID is invalid")
		}
	}

	// Check if salespeople exist or not.
	if campusDrive.SalesPeople != nil && len(campusDrive.SalesPeople) != 0 {
		var count int = 0
		var salesPersonIDs []uuid.UUID
		for _, salesPerson := range campusDrive.SalesPeople {
			salesPersonIDs = append(salesPersonIDs, salesPerson.ID)
		}
		err := service.Repository.GetCountForTenant(uow, campusDrive.TenantID, general.User{}, &count,
			repository.Filter("`id` IN (?)", salesPersonIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(campusDrive.SalesPeople) {
			log.NewLogger().Error("SalesPerson ID is invalid")
			return errors.NewValidationError("SalesPerson ID is invalid")
		}
	}

	// Check if faculties exist or not.
	if campusDrive.Faculties != nil && len(campusDrive.Faculties) != 0 {
		var count int = 0
		var facultyIDs []uuid.UUID
		for _, faculty := range campusDrive.Faculties {
			facultyIDs = append(facultyIDs, faculty.ID)
		}
		err := service.Repository.GetCountForTenant(uow, campusDrive.TenantID, faculty.Faculty{}, &count,
			repository.Filter("`id` IN (?)", facultyIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(campusDrive.Faculties) {
			log.NewLogger().Error("Faculty ID is invalid")
			return errors.NewValidationError("Faculty ID is invalid")
		}
	}

	// Check if developers exist or not.
	if campusDrive.Developers != nil && len(campusDrive.Developers) != 0 {
		var count int = 0
		var developerIDs []uuid.UUID
		for _, developer := range campusDrive.Developers {
			developerIDs = append(developerIDs, developer.ID)
		}
		err := service.Repository.GetCountForTenant(uow, campusDrive.TenantID, general.Employee{}, &count,
			repository.Filter("`id` IN (?)", developerIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(campusDrive.Developers) {
			log.NewLogger().Error("Developer ID is invalid")
			return errors.NewValidationError("Developer ID is invalid")
		}
	}

	return nil
}

// doesCampusDriveNameExist check for same campus drive name conflict, if campus drive name exists return true.
func (service *CampusDriveService) doesCampusDriveNameExist(campusDrive *college.CampusDrive) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, campusDrive.TenantID, &college.CampusDrive{},
		repository.Filter("`campus_name`=? AND `id`!=?", campusDrive.CampusName, campusDrive.ID))
	if err := util.HandleIfExistsError("Campus Drive Name already exists", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// addSearchQueries adds all search queries if any when getAll is called.
func (service *CampusDriveService) addSearchQueries(searchForm url.Values, tenantID uuid.UUID) []repository.QueryProcessor {
	fmt.Println("=========================In campus drive search============================", searchForm)

	// Check if there is search criteria given.
	if len(searchForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcesors []repository.QueryProcessor

	// Get login id from params.
	loginID, _ := searchForm["loginID"]

	// Get role name form params.
	roleName, _ := searchForm["roleName"]

	// Campus name search.
	if _, ok := searchForm["campusName"]; ok {
		util.AddToSlice("`campus_name`", "LIKE ?", "AND", "%"+searchForm.Get("campusName")+"%", &columnNames, &conditions, &operators, &values)
	}

	// From date search.
	if fromDate, ok := searchForm["fromDate"]; ok {
		util.AddToSlice("`campus_date`", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}

	// To date search.
	if toDate, ok := searchForm["toDate"]; ok {
		util.AddToSlice("`campus_date`", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}

	// Salesperson login and search.
	// Get salesperson ids from params.
	salesPersonIDs, _ := searchForm["salesPersonIDs"]

	// If rolename exists and is salesperson or salesperson ids search is given.
	if (len(roleName) != 0 && searchForm.Get("roleName") == "salesperson") || len(salesPersonIDs) != 0 {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN campus_drives_sales_people ON campus_drives.`id` = campus_drives_sales_people.`campus_drive_id`"),
			repository.Join("JOIN users ON users.`id` = campus_drives_sales_people.`sales_person_id`"),
			repository.Filter("users.`deleted_at` IS NULL"),
			repository.Filter("users.`tenant_id`=?", tenantID))

		// If role is salesperson then give filter for salesperson id.
		if len(roleName) != 0 && searchForm.Get("roleName") == "salesperson" {
			queryProcesors = append(queryProcesors, repository.Filter("users.`id` = ?", loginID))
		}

		// If salesperson ids in search is given then give filter for those saleseprson.
		if len(salesPersonIDs) != 0 {
			util.AddToSlice("users.`id`", "IN(?)", "AND", salesPersonIDs,
				&columnNames, &conditions, &operators, &values)
		}
	}

	// Faculty login and search.
	// Get faculty ids from params.
	facultyIDs, _ := searchForm["facultyIDs"]

	// If rolename exists and is faculty or faculty ids search is given.
	if (len(roleName) != 0 && searchForm.Get("roleName") == "faculty") || len(facultyIDs) != 0 {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN campus_drives_faculty ON campus_drives.`id` = campus_drives_faculty.`campus_drive_id`"),
			repository.Join("JOIN faculties ON faculties.`id` = campus_drives_faculty.`faculty_id`"),
			repository.Filter("faculties.`deleted_at` IS NULL"),
			repository.Filter("faculties.`tenant_id`=?", tenantID))

		// If role is faculty then give filter for faculty id.
		if len(roleName) != 0 && searchForm.Get("roleName") == "faculty" {
			queryProcesors = append(queryProcesors, repository.Filter("faculties.`id` = ?", loginID))
		}

		// If faculty ids in search is given then give filter for those faculties.
		if len(facultyIDs) != 0 {
			util.AddToSlice("faculties.`id`", "IN(?)", "AND", facultyIDs,
				&columnNames, &conditions, &operators, &values)
		}
	}

	// Developer login and search.
	// Get developer ids from params.
	developerIDs, _ := searchForm["developerIDs"]

	// If rolename exists and is developer or developer ids search is given.
	if (len(roleName) != 0 && searchForm.Get("roleName") == "developer") || len(developerIDs) != 0 {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN campus_drives_developers ON campus_drives.`id` = campus_drives_developers.`campus_drive_id`"),
			repository.Join("JOIN employees ON employees.`id` = campus_drives_developers.`developer_id`"),
			repository.Filter("employees.`deleted_at` IS NULL"),
			repository.Filter("employees.`tenant_id`=?", tenantID))

		// If role is developer then give filter for developer id.
		if len(roleName) != 0 && searchForm.Get("roleName") == "developer" {
			queryProcesors = append(queryProcesors, repository.Filter("employees.`id` = ?", loginID))
		}

		// If developer ids in search is given then give filter for those developer.
		if len(developerIDs) != 0 {
			util.AddToSlice("employees.`id`", "IN(?)", "AND", developerIDs,
				&columnNames, &conditions, &operators, &values)
		}
	}

	// Colege branch search.
	if collegeIDs, ok := searchForm["collegeIDs"]; ok {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN campus_drives_college_branches ON campus_drives.`id` = campus_drives_college_branches.`campus_drive_id`"),
			repository.Join("JOIN college_branches ON college_branches.`id` = campus_drives_college_branches.`branch_id`"),
			repository.Filter("college_branches.`deleted_at` IS NULL"),
			repository.Filter("college_branches.`tenant_id`=?", tenantID))

		util.AddToSlice("campus_drives_college_branches.`branch_id`", "IN(?)", "AND", collegeIDs,
			&columnNames, &conditions, &operators, &values)
	}

	// Company requirement search.
	if companyRequirementIDs, ok := searchForm["companyRequirementIDs"]; ok {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN campus_drives_company_requirements ON campus_drives.`id` = campus_drives_company_requirements.`campus_drive_id`"),
			repository.Join("JOIN company_requirements ON company_requirements.`id` = campus_drives_company_requirements.`requirement_id`"),
			repository.Filter("company_requirements.`deleted_at` IS NULL"),
			repository.Filter("company_requirements.`tenant_id`=?", tenantID))

		util.AddToSlice("campus_drives_company_requirements.`requirement_id`", "IN(?)", "AND", companyRequirementIDs,
			&columnNames, &conditions, &operators, &values)
	}

	// Group by campus drive id and add all filters.
	queryProcesors = append(queryProcesors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("campus_drives.`id`"))
	return queryProcesors
}

// getValuesForCampusDrives gets values for campus drives by firing individual query for each campus drive.
func (service *CampusDriveService) getValuesForCampusDrives(uow *repository.UnitOfWork, campsuDrives *[]college.CampusDriveDTO,
	tenantID uuid.UUID) error {
	for index := range *campsuDrives {
		//********************************************Total Required Candidates*********************************************************************

		// Collect ids of all company requirements of campus drive.
		var companyRequirementIDs []uuid.UUID

		for _, companyRequirement := range (*campsuDrives)[index].CompanyRequirements {
			companyRequirementIDs = append(companyRequirementIDs, companyRequirement.ID)
		}

		// Calculate total requirements by adding vacancy of all company requirements.
		var totalVacancy college.TotalVacancy
		if err := service.Repository.Scan(uow, &totalVacancy,
			repository.Filter("company_requirements.`deleted_at` IS NULL"),
			repository.Filter("company_requirements.`tenant_id`=?", tenantID),
			repository.Filter("company_requirements.`id` IN (?)", companyRequirementIDs),
			repository.Table("company_requirements"),
			repository.Select("sum(company_requirements.`vacancy`) as total_vacancy")); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Assign total vacancy to total requirements of campus drive.
		tempVacancy := uint16(totalVacancy.TotalVacancy)
		(*campsuDrives)[index].TotalRequirements = &tempVacancy

		//********************************************Total Registered Candidates*********************************************************************
		// Create bucket for total registered candidates.
		var totalRegisteredCount uint16 = 0

		// Get count of all registered candidates from database.
		err := service.Repository.GetCountForTenant(uow, tenantID, college.CampusTalentRegistration{}, &totalRegisteredCount,
			repository.Filter("`campus_drive_id`=?", (*campsuDrives)[index].ID),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}

		// Assign total registered to total registered candidates of campus drive.
		(*campsuDrives)[index].TotalRegisteredCandidates = &totalRegisteredCount

		//********************************************Total Appeared Candidates*********************************************************************
		// Create bucket for total appeared candidates.
		var totalAppearedCount uint16 = 0

		// Get count of all appeared candidates from database.
		err = service.Repository.GetCountForTenant(uow, tenantID, college.CampusTalentRegistration{}, &totalAppearedCount,
			repository.Filter("`campus_drive_id`=?", (*campsuDrives)[index].ID),
			repository.Filter("`has_attempted`=?", true),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}

		// Assign total appeared to total appeared candidates of campus drive.
		(*campsuDrives)[index].TotalAppearedCandidates = &totalAppearedCount
	}
	return nil
}
