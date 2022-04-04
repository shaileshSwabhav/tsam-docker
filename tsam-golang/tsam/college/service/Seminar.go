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
	"github.com/techlabs/swabhav/tsam/models/general"

	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// SeminarService provides methods to do Update, Delete, Add, Get operations on Seminar.
// associations field will contain details about the sub-structs in seminar for preload and other operations.
type SeminarService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewSeminarService returns a new instance of Seminar.
func NewSeminarService(db *gorm.DB, repository repository.Repository) *SeminarService {
	return &SeminarService{
		DB:           db,
		Repository:   repository,
		associations: []string{"CollegeBranches", "SalesPeople", "Speakers"},
	}
}

// AddSeminar adds a new seminar to database.
func (service *SeminarService) AddSeminar(seminar *college.Seminar) error {
	// Get credential id from CreatedBy field of seminar(set in controller).
	credentialID := seminar.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(seminar.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, seminar.TenantID); err != nil {
		return err
	}

	// Validate foreign key ids.
	if err := service.doForeignKeysExist(seminar); err != nil {
		return err
	}

	// Check for same seminar name.
	if err := service.doesSeminarNameExist(seminar); err != nil {
		return err
	}

	// Starting transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Assign seminar code.
	var codeError error
	seminar.Code, codeError = util.GenerateUniqueCode(uow.DB, seminar.SeminarName, "`code` = ?", &college.Seminar{})
	if codeError != nil {
		log.NewLogger().Error(codeError.Error())
		uow.RollBack()
		return errors.NewHTTPError("Internal server err or", http.StatusInternalServerError)
	}

	// Create the registration link and assign it.
	tempRegLink := "https://swabhavtechlabs.com/test/tsm-forms/#/seminar-registration?code=" + seminar.Code
	seminar.StudentRegistrationLink = &tempRegLink

	// Add seminar to database.
	if err := service.Repository.Add(uow, seminar); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Seminar could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// UpdateSeminar updates the seminar along with it's associations in database.
func (service *SeminarService) UpdateSeminar(seminar *college.Seminar) error {
	// Get credential id from UpdatedBy field of seminar(set in controller).
	credentialID := seminar.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(seminar.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, seminar.TenantID); err != nil {
		return err
	}

	// Validate foreign key ids.
	if err := service.doForeignKeysExist(seminar); err != nil {
		return err
	}

	// Validate seminar id.
	if err := service.doesSeminarExist(seminar.ID, seminar.TenantID); err != nil {
		return err
	}

	// Check for same seminar name.
	if err := service.doesSeminarNameExist(seminar); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting seminar already present in database.
	tempSeminar := college.Seminar{}

	// Get seminar for getting created_by field of seminar from database.
	if err := service.Repository.GetForTenant(uow, seminar.TenantID, seminar.ID, &tempSeminar,
		repository.PreloadAssociations(service.associations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp seminar to seminar to be updated.
	seminar.CreatedBy = tempSeminar.CreatedBy

	// Give code to seminar.
	seminar.Code = tempSeminar.Code

	// Update seminar associations.
	if err := service.updateSeminarAssociations(uow, seminar); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Seminar could not be updated", http.StatusInternalServerError)
	}

	// Update seminar.
	if err := service.Repository.Save(uow, seminar); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Seminar could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteSeminar deletes seminar from database.
func (service *SeminarService) DeleteSeminar(seminar *college.Seminar) error {
	// Get credential id from DeletedBy field of seminar(set in controller).
	credentialID := seminar.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(seminar.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, seminar.TenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get seminar for updating deleted_by field of seminar.
	if err := service.Repository.GetForTenant(uow, seminar.TenantID, seminar.ID, seminar); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//  Delete seminar association from database.
	if err := service.deleteSeminarAssociations(uow, seminar, credentialID); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Failed to delete seminar", http.StatusInternalServerError)
	}

	// Update seminar for updating deleted_by and deleted_at fields of seminar.
	if err := service.Repository.UpdateWithMap(uow, seminar, map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()}); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Seminar could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetSeminar returns seminar filtering by ID.
func (service *SeminarService) GetSeminar(seminar *college.Seminar) error {
	// Validate tenant id.
	if err := service.doesTenantExist(seminar.TenantID); err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get seminar.
	if err := service.Repository.GetForTenant(uow, seminar.TenantID, seminar.ID, seminar,
		repository.PreloadAssociations(service.associations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetSeminarByCode returns seminar filtering by code.
func (service *SeminarService) GetSeminarByCode(seminar *college.Seminar) error {
	// Validate tenant id.
	if err := service.doesTenantExist(seminar.TenantID); err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get seminar.
	if err := service.Repository.GetRecordForTenant(uow, seminar.TenantID, seminar,
		repository.Filter("`code` = ?", seminar.Code),
		repository.PreloadAssociations(service.associations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetAllSeminars returns all seminars in database.
func (service *SeminarService) GetAllSeminars(seminars *[]college.SeminarDTO,
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
		repository.Paginate(limit, offset, totalCount),
		repository.Filter("seminars.`tenant_id`=?", tenantID))

	// Get seminars from database.
	if err := service.Repository.GetAllInOrder(uow, seminars, "`seminar_date`", queryProcessors...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Range seminars for getting total registered students and total apperaed students.
	err := service.getValuesForSeminars(uow, seminars, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// updateSeminarAssociations updates seminar's associations.
func (service *SeminarService) updateSeminarAssociations(uow *repository.UnitOfWork, seminar *college.Seminar) error {

	// Replace college branches of seminar.
	if err := service.Repository.ReplaceAssociations(uow, seminar, "CollegeBranches",
		seminar.CollegeBranches); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace salesPeople of seminar.
	if err := service.Repository.ReplaceAssociations(uow, seminar, "SalesPeople",
		seminar.SalesPeople); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace speakers of seminar.
	if err := service.Repository.ReplaceAssociations(uow, seminar, "Speakers",
		seminar.Speakers); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	seminar.CollegeBranches = nil
	seminar.SalesPeople = nil
	seminar.Speakers = nil

	return nil
}

// deleteSeminarAssociations deletes seminar's associations (sub-structs).
func (service *SeminarService) deleteSeminarAssociations(uow *repository.UnitOfWork, seminar *college.Seminar,
	credentialID uuid.UUID) error {

	//**********************************************Delete seminar talent registrations***************************************************
	if err := service.Repository.UpdateWithMap(uow, &college.SeminarTalentRegistration{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`seminar_id`=? AND `tenant_id`=?", seminar.ID, seminar.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Seminar could not be deleted", http.StatusInternalServerError)
	}

	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *SeminarService) doesTenantExist(tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if tenant(parent tenant) exists or not.
	exists, err := repository.DoesRecordExist(uow.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesSeminarExist validates if seminar exists or not in database.
func (service *SeminarService) doesSeminarExist(seminarID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.Seminar{},
		repository.Filter("`id` = ?", seminarID))
	if err := util.HandleError("Seminar not found", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credental exists or not in database.
func (service *SeminarService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExist validates if foreign keys are present or not in database.
func (service *SeminarService) doForeignKeysExist(seminar *college.Seminar) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if college branches exists or not.
	if seminar.CollegeBranches != nil && len(seminar.CollegeBranches) != 0 {
		var count int = 0
		var collegeBranchIDs []uuid.UUID
		for _, collegeBranch := range seminar.CollegeBranches {
			collegeBranchIDs = append(collegeBranchIDs, collegeBranch.ID)
		}
		err := service.Repository.GetCountForTenant(uow, seminar.TenantID, college.Branch{}, &count,
			repository.Filter("`id` IN (?)", collegeBranchIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(seminar.CollegeBranches) {
			log.NewLogger().Error("College Branch ID is invalid")
			return errors.NewValidationError("College Branch ID is invalid")
		}
	}

	// Check if salespeople exist or not.
	if seminar.SalesPeople != nil && len(seminar.SalesPeople) != 0 {
		var count int = 0
		var salesPersonIDs []uuid.UUID
		for _, salesPerson := range seminar.SalesPeople {
			salesPersonIDs = append(salesPersonIDs, salesPerson.ID)
		}
		err := service.Repository.GetCountForTenant(uow, seminar.TenantID, general.User{}, &count,
			repository.Filter("`id` IN (?)", salesPersonIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(seminar.SalesPeople) {
			log.NewLogger().Error("SalesPerson ID is invalid")
			return errors.NewValidationError("SalesPerson ID is invalid")
		}
	}

	// Check if speakers exist or not.
	if seminar.Speakers != nil && len(seminar.Speakers) != 0 {
		var count int = 0
		var speakerIDs []uuid.UUID
		for _, speaker := range seminar.Speakers {
			speakerIDs = append(speakerIDs, speaker.ID)
		}
		err := service.Repository.GetCountForTenant(uow, seminar.TenantID, college.Speaker{}, &count,
			repository.Filter("`id` IN (?)", speakerIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(seminar.Speakers) {
			log.NewLogger().Error("Speaker ID is invalid")
			return errors.NewValidationError("Speaker ID is invalid")
		}
	}

	return nil
}

// doesSeminarNameExist check for same seminar name conflict, if seminar name exists return true.
func (service *SeminarService) doesSeminarNameExist(seminar *college.Seminar) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, seminar.TenantID, &college.Seminar{},
		repository.Filter("`seminar_name`=? AND `id`!=?", seminar.SeminarName, seminar.ID))
	if err := util.HandleIfExistsError("Seminar Name already exists", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// addSearchQueries adds all search queries if any when getAll is called.
func (service *SeminarService) addSearchQueries(searchForm url.Values, tenantID uuid.UUID) []repository.QueryProcessor {
	fmt.Println("=========================In seminar search============================", searchForm)

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

	// Salesperson login and search.
	// Get salesperson ids from params.
	salesPersonIDs, _ := searchForm["salesPersonIDs"]

	// If rolename exists and is salesperson or salesperson ids search is given.
	if (len(roleName) != 0 && searchForm.Get("roleName") == "salesperson") || len(salesPersonIDs) != 0 {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN seminars_sales_people ON seminars.`id` = seminars_sales_people.`seminar_id`"),
			repository.Join("JOIN users ON users.`id` = seminars_sales_people.`sales_person_id`"),
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

	// Seminar name search.
	if _, ok := searchForm["seminarName"]; ok {
		util.AddToSlice("`seminar_name`", "LIKE ?", "AND", "%"+searchForm.Get("seminarName")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Seminar from date search.
	if fromDate, ok := searchForm["fromDate"]; ok {
		util.AddToSlice("`seminar_date`", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}

	// Seminar to date search.
	if toDate, ok := searchForm["toDate"]; ok {
		util.AddToSlice("`seminar_date`", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}

	// Colege branch search.
	if collegeIDs, ok := searchForm["collegeIDs"]; ok {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN seminars_college_branches ON seminars.`id` = seminars_college_branches.`seminar_id`"),
			repository.Join("JOIN college_branches ON college_branches.`id` = seminars_college_branches.`branch_id`"),
			repository.Filter("college_branches.`deleted_at` IS NULL"),
			repository.Filter("college_branches.`tenant_id`=?", tenantID))

		util.AddToSlice("seminars_college_branches.`branch_id`", "IN(?)", "AND", collegeIDs,
			&columnNames, &conditions, &operators, &values)
	}

	// Speaker search.
	if speakerIDs, ok := searchForm["speakerIDs"]; ok {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN seminars_speakers ON seminars.`id` = seminars_speakers.`seminar_id`"),
			repository.Join("JOIN speakers ON speakers.`id` = seminars_speakers.`speaker_id`"),
			repository.Filter("speakers.`deleted_at` IS NULL"),
			repository.Filter("speakers.`tenant_id`=?", tenantID))

		util.AddToSlice("seminars_speakers.`speaker_id`", "IN(?)", "AND", speakerIDs,
			&columnNames, &conditions, &operators, &values)
	}

	// Group by seminar id and add all filters.
	queryProcesors = append(queryProcesors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("seminars.`id`"))
	return queryProcesors
}

// getValuesForSeminars gets values for seminars by firing individual query for each seminar.
func (service *SeminarService) getValuesForSeminars(uow *repository.UnitOfWork, semianrs *[]college.SeminarDTO,
	tenantID uuid.UUID) error {
	for index := range *semianrs {
		//********************************************Total Registered Students*********************************************************************
		// Create bucket for total registered students.
		var totalRegisteredCount uint16 = 0

		// Get count of all registered students from database.
		err := service.Repository.GetCountForTenant(uow, tenantID, college.SeminarTalentRegistration{}, &totalRegisteredCount,
			repository.Filter("`seminar_id`=?", (*semianrs)[index].ID),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}

		// Assign total registered to total registered students of seminar.
		(*semianrs)[index].TotalRegisteredStudents = &totalRegisteredCount

		//********************************************Total Visited Candidates*********************************************************************
		// Create bucket for total appeared candidates.
		var totalVisitedCount uint16 = 0

		// Get count of all appeared candidates from database.
		err = service.Repository.GetCountForTenant(uow, tenantID, college.SeminarTalentRegistration{}, &totalVisitedCount,
			repository.Filter("`seminar_id`=?", (*semianrs)[index].ID),
			repository.Filter("`has_visited`=?", true),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}

		// Assign total appeared to total appeared candidates of seminar.
		(*semianrs)[index].TotalVisitedStudents = &totalVisitedCount
	}
	return nil
}
