package service

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	college "github.com/techlabs/swabhav/tsam/models/college"
	general "github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	talentService "github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
)

// CandidateService provides method to update, delete, add, get all, get one for candidate.
type CandidateService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// campusTalentAssociations provides preload associations array for talent.
var campusTalentAssociations []string = []string{"Academics", "Academics.Degree", "Academics.Specialization", "State", "Country"}

// NewCandidateService returns new instance of CandidateService.
func NewCandidateService(db *gorm.DB, repository repository.Repository) *CandidateService {
	return &CandidateService{
		DB:         db,
		Repository: repository,
	}
}

// AddCandidateWithCredential adds one candidate to database with credential id.
func (service *CandidateService) AddCandidateWithCredential(candidate *college.Candidate) error {
	// Get credential id from CreatedBy field of candidate(set in controller).
	credentialID := candidate.CreatedBy

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, candidate.TenantID); err != nil {
		return err
	}

	// Call add candidate service method.
	if err := service.AddCandidate(candidate); err != nil {
		return err
	}

	return nil
}

// AddCandidate adds one candidate to database.
func (service *CandidateService) AddCandidate(candidate *college.Candidate) error {

	// Validate tenant id.
	if err := service.doesTenantExist(candidate.TenantID); err != nil {
		return err
	}

	// Validate campus drive id.
	if err := service.doesCampusDriveExist(candidate.TenantID, candidate.CampusDriveID); err != nil {
		return err
	}

	// Check if same email exists in campus drive.
	if err := service.doesEmailExist(candidate); err != nil {
		return err
	}

	// Give current date to registration field of candidate.
	currentTime := time.Now()
	candidate.RegistrationDate = currentTime.Format("2006-01-02")

	// Validate compulsary fields.
	err := candidate.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get source from database with name as 'cd'.
	tempSource := general.Source{}
	if err := service.Repository.GetRecordForTenant(uow, candidate.TenantID, &tempSource,
		repository.Filter("`name`=?", "cd")); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}

	// Set college id by college name.
	if err := service.setCollegeNameAndID(uow, candidate); err != nil {
		uow.RollBack()
		return err
	}

	// Create talent.
	talent := tal.Talent{}
	service.createTalent(&talent, candidate)

	// Give soucrce id to source id field of talent.
	talent.SourceID = &tempSource.ID

	// Flag to check if talent already exists.
	doesTalentExist := false

	// Add talent to database.
	talentService := talentService.NewTalentService(service.DB, service.Repository)
	exists, err := talentService.AddTalent(&talent, true, uow)
	if err != nil {
		// Check if talent already exists.
		if exists {
			doesTalentExist = true
		}
		if !doesTalentExist {
			uow.RollBack()
			return err
		}
	}

	// Create a temp talent bucket to get already present talent from database.
	tempTalent := tal.Talent{}

	// If talent already exists then get the talent for its id.
	if doesTalentExist {

		// Get temp talent from database.
		if err := service.Repository.GetRecordForTenant(uow, candidate.TenantID, &tempTalent,
			repository.Filter("`email`=?", candidate.Email),
			repository.PreloadAssociations([]string{"Academics"})); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewValidationError("Candidate could not be added")
		}

		// Assign some fields of talent to temp talent.
		tempTalent.FirstName = talent.FirstName
		tempTalent.LastName = talent.LastName
		tempTalent.Contact = talent.Contact
		tempTalent.AcademicYear = talent.AcademicYear
		tempTalent.Address = talent.Address
		tempTalent.Country = talent.Country
		tempTalent.State = talent.State
		tempTalent.City = talent.City
		tempTalent.PINCode = talent.PINCode
		tempTalent.UpdatedBy = candidate.CreatedBy

		// Store academics of temp talent in a bucket.
		tempTalentAcademics := tempTalent.Academics

		// Make temp talent academics nil to avoid any unnecessary updates.
		tempTalent.Academics = nil

		// Update temp talent with the new fields.
		if err := service.Repository.Save(uow, tempTalent); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Candidate could not be added", http.StatusInternalServerError)
		}

		// Create flag for academic present.
		academicPresent := false

		// Range the temp academic to see if talent academic is already presesnt or not.
		for i := range tempTalentAcademics {

			// Check if college is already present in temp talent.
			if tempTalentAcademics[i].College == talent.Academics[0].College {
				academicPresent = true

				// If college is already present then update the temp talent academic.
				// Give the talent academic fields to temp talent academic.
				tempTalentAcademics[i].DegreeID = talent.Academics[0].DegreeID
				tempTalentAcademics[i].SpecializationID = talent.Academics[0].SpecializationID
				tempTalentAcademics[i].Percentage = talent.Academics[0].Percentage
				tempTalentAcademics[i].Passout = talent.Academics[0].Passout
				tempTalentAcademics[i].UpdatedBy = candidate.CreatedBy

				// Update temp talent academic.
				if err := service.Repository.Save(uow, tempTalentAcademics[i]); err != nil {
					log.NewLogger().Error(err.Error())
					uow.RollBack()
					return errors.NewHTTPError("Candidate could not be added", http.StatusInternalServerError)
				}
				break
			}
		}

		// If academic is not present.
		if !academicPresent {
			// If college is not present in temp talent then create new academic.
			// Create bucket for academic.
			newAcademic := tal.Academic{}

			// Give new fields from talent to academic.
			newAcademic.DegreeID = talent.Academics[0].DegreeID
			newAcademic.SpecializationID = talent.Academics[0].SpecializationID
			newAcademic.Percentage = talent.Academics[0].Percentage
			newAcademic.Passout = talent.Academics[0].Passout
			newAcademic.College = talent.Academics[0].College
			newAcademic.CollegeID = talent.Academics[0].CollegeID
			newAcademic.TenantID = candidate.TenantID
			newAcademic.CreatedBy = candidate.CreatedBy

			// Give talent id to new academic.
			newAcademic.TalentID = tempTalent.ID

			// Add new academic.
			if err := service.Repository.Add(uow, &newAcademic); err != nil {
				log.NewLogger().Error(err.Error())
				return errors.NewHTTPError("Candidate could not be added", http.StatusInternalServerError)
			}
		}
	}

	// Create campus talent registration.
	campusTalReg := college.CampusTalentRegistration{}
	service.createCampusTalentRegistration(&campusTalReg, candidate)

	// If talent already existed then give the temp talent id, else give the new talent id.
	if doesTalentExist {
		campusTalReg.TalentID = tempTalent.ID
	} else {
		campusTalReg.TalentID = talent.ID
	}

	// Add campus talent registration to database.
	if err := service.Repository.Add(uow, &campusTalReg); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Candidate could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetCandidates gets all candidates from database.
func (service *CandidateService) GetCandidates(candidates *[]college.CandidateDTO, tenantID uuid.UUID, campusDriveID uuid.UUID,
	limit int, offset int, totalCount *int, searchForm url.Values, uows ...*repository.UnitOfWork) error {
	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork

	// Check if uow is passed in argument or not.
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate campus drive id.
	if err := service.doesCampusDriveExist(tenantID, campusDriveID); err != nil {
		return err
	}

	// Create bucket for talents and campus talent registrations.
	talents := []tal.DTO{}
	campusTalRegs := []college.CampusTalentRegistration{}

	// Query processors for search and get for campus talent registrations.
	var queryProcessorsForCampusTalReg []repository.QueryProcessor
	queryProcessorsForCampusTalReg = append(queryProcessorsForCampusTalReg,
		repository.Join("JOIN talents ON talents.`id` = campus_talent_registrations.`talent_id`"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("campus_talent_registrations.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Filter("`campus_drive_id`=?", campusDriveID),
		repository.Paginate(limit, offset, totalCount),
	)
	queryProcessorsForCampusTalReg = append(queryProcessorsForCampusTalReg, service.addSearchQueries(searchForm, tenantID)...)

	// Get all campus talent registrations from database.
	if err := service.Repository.GetAllInOrder(uow, &campusTalRegs, "`first_name`, `last_name`", queryProcessorsForCampusTalReg...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for search and get for talents.
	queryProcessorsForTalent := service.addSearchQueries(searchForm, tenantID)
	queryProcessorsForTalent = append(queryProcessorsForTalent,
		repository.Join("JOIN campus_talent_registrations ON talents.`id` = campus_talent_registrations.`talent_id`"),
		repository.Filter("campus_talent_registrations.`deleted_at` IS NULL"),
		repository.Filter("campus_talent_registrations.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Filter("`campus_drive_id`=?", campusDriveID),
		repository.PreloadAssociations(campusTalentAssociations),
		repository.Paginate(limit, offset, totalCount),
	)

	// Get talents from database.
	if err := service.Repository.GetAllInOrder(uow, &talents, "`first_name`, `last_name`", queryProcessorsForTalent...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give campus talent registrations to candidates.
	for _, campusTalReg := range campusTalRegs {
		candidate := college.CandidateDTO{}
		service.createCampusTalentRegistrationPartOfCandidate(&campusTalReg, &candidate)
		*candidates = append(*candidates, candidate)
	}

	// Sort talent academics.
	service.sortTalentAcademicsByOrder(&talents)

	// Give talents to candidates.
	for _, talent := range talents {
		for i := 0; i < len(*candidates); i++ {
			if talent.ID == (*candidates)[i].TalentID {
				service.createTalentPartOfCandidate(&talent, &(*candidates)[i])
			}
		}
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}

	return nil
}

// UpdateCandidate updates candidate in Database.
func (service *CandidateService) UpdateCandidate(candidate *college.Candidate) error {
	// Get credential id from UpdatedBy field of candidate(set in controller).
	credentialID := candidate.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(candidate.TenantID); err != nil {
		return err
	}

	// Validate candidate id.
	if err := service.doesCampusTalentRegistrationExist(candidate.CampusTalentRegistrationID, candidate.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, candidate.TenantID); err != nil {
		return err
	}

	// Validate campus drive id.
	if err := service.doesCampusDriveExist(candidate.TenantID, candidate.CampusDriveID); err != nil {
		return err
	}

	// Check if same email exists in campus drive.
	if err := service.doesEmailExist(candidate); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Set college id by college name.
	if err := service.setCollegeNameAndID(uow, candidate); err != nil {
		uow.RollBack()
		return err
	}

	// Create bucket for getting candidate already present in database.
	tempCampusTalReg := college.CampusTalentRegistration{}

	// Get candidate for getting created_by field of candidate from database.
	if err := service.Repository.GetForTenant(uow, candidate.TenantID, candidate.CampusTalentRegistrationID, &tempCampusTalReg); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp candidate to candidate to be updated.
	candidate.CreatedBy = tempCampusTalReg.CreatedBy

	// Create talent.
	talent := tal.Talent{}
	service.createTalent(&talent, candidate)

	// Check if talent with email already exists.
	isTalentRecordNotFound := false
	tempTalent := tal.Talent{}
	if err := service.Repository.GetRecordForTenant(uow, candidate.TenantID, &tempTalent,
		repository.Filter("`email`=? AND `id`!=?", candidate.Email, talent.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			isTalentRecordNotFound = true
		} else {
			uow.RollBack()
			return errors.NewValidationError("Internal Server Error")
		}
	}

	if !isTalentRecordNotFound {
		// If talent already exists then send error as email already exists.
		if util.IsUUIDValid(tempTalent.ID) {
			return errors.NewValidationError("Email already exists")
		}
	}

	// Update talent to database.
	talentService := talentService.NewTalentService(service.DB, service.Repository)
	err := talentService.UpdateTalent(&talent, uow)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Create campus talent registration.
	campusTalReg := college.CampusTalentRegistration{}
	service.createCampusTalentRegistration(&campusTalReg, candidate)

	// Add campus talent registration to database.
	if err := service.Repository.Save(uow, &campusTalReg); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Campus talent registration could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteCandidate deletes one candidate form database.
func (service *CandidateService) DeleteCandidate(candidate *college.Candidate, uows ...*repository.UnitOfWork) error {
	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Get credential id from DeletedBy field of candidate(set in controller).
	credentialID := candidate.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(candidate.TenantID); err != nil {
		return err
	}

	// Validate candidate id.
	if err := service.doesCampusTalentRegistrationExist(candidate.CampusTalentRegistrationID, candidate.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, candidate.TenantID); err != nil {
		return err
	}

	// Validate campus drive id.
	if err := service.doesCampusDriveExist(candidate.TenantID, candidate.CampusDriveID); err != nil {
		return err
	}

	// Update candidate for updating deleted_by and deleted_at field of candidate.
	if err := service.Repository.UpdateWithMap(uow, &college.CampusTalentRegistration{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`id`=? AND `tenant_id`=?", candidate.CampusTalentRegistrationID, candidate.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Campus Talent Registration could not be deleted", http.StatusInternalServerError)
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// UpdateMultipleCandidate updates multiple fields of campus talent registrations.
func (service *CandidateService) UpdateMultipleCandidate(updateMultipleCandidate *college.UpdateMultipleCandidate) error {
	// Validate tenant id.
	if err := service.doesTenantExist(updateMultipleCandidate.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(updateMultipleCandidate.UpdatedBy, updateMultipleCandidate.TenantID); err != nil {
		return err
	}

	// Validate campus drive.
	if err := service.doesCampusDriveExist(updateMultipleCandidate.TenantID, updateMultipleCandidate.CampusDriveID); err != nil {
		return err
	}

	// Validate all campus talent registrations.
	for _, campusTalentRegistrationID := range updateMultipleCandidate.CampusTalentRegistrationIDs {
		if err := service.doesCampusTalentRegistrationExist(campusTalentRegistrationID, updateMultipleCandidate.TenantID); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create map for upadting values.
	updateMap := make(map[string]interface{})

	if updateMultipleCandidate.HasAttempted != nil {
		updateMap["HasAttempted"] = *updateMultipleCandidate.HasAttempted
	}

	if updateMultipleCandidate.IsTestLinkSent != nil {
		updateMap["IsTestLinkSent"] = *updateMultipleCandidate.IsTestLinkSent
	}

	if updateMultipleCandidate.Result != nil {
		updateMap["Result"] = *updateMultipleCandidate.Result
	}

	// Update campus talent registrayion field of candidates.
	if err := service.Repository.UpdateWithMap(uow, &college.CampusTalentRegistration{}, updateMap,
		repository.Filter("`id` IN (?)", updateMultipleCandidate.CampusTalentRegistrationIDs)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Candidate(s) could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds all search queries if any when getAll is called.
func (service *CandidateService) addSearchQueries(searchForm url.Values, tenantID uuid.UUID) []repository.QueryProcessor {
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

	// First name search.
	if _, ok := searchForm["firstName"]; ok {
		util.AddToSlice("`first_name`", "LIKE ?", "AND", "%"+searchForm.Get("firstName")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Last name search.
	if _, ok := searchForm["lastName"]; ok {
		util.AddToSlice("`last_name`", "LIKE ?", "AND", "%"+searchForm.Get("lastName")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Email search.
	if _, ok := searchForm["email"]; ok {
		util.AddToSlice("`email`", "LIKE ?", "AND", "%"+searchForm.Get("email")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Contact search.
	if _, ok := searchForm["contact"]; ok {
		util.AddToSlice("`contact`", "LIKE ?", "AND", "%"+searchForm.Get("contact")+"%", &columnNames, &conditions, &operators, &values)
	}

	// From date search.
	if fromDate, ok := searchForm["fromDate"]; ok {
		util.AddToSlice("`registration_date`", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}

	// To date search.
	if toDate, ok := searchForm["toDate"]; ok {
		util.AddToSlice("`registration_date`", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}

	// Academic year search.
	if academicYear, ok := searchForm["academicYear"]; ok {
		util.AddToSlice("`academic_year`", "=?", "AND", academicYear, &columnNames, &conditions, &operators, &values)
	}

	// Is swabhav talent search.
	if isSwabhavTalent, ok := searchForm["isSwabhavTalent"]; ok {
		util.AddToSlice("`is_swabhav_talent`", "=?", "AND", isSwabhavTalent, &columnNames, &conditions, &operators, &values)
	}

	// Has attempted search.
	if hasAttempted, ok := searchForm["hasAttempted"]; ok {
		util.AddToSlice("`has_attempted`", "=?", "AND", hasAttempted, &columnNames, &conditions, &operators, &values)
	}

	// Is test link sent search.
	if isTestLinkSent, ok := searchForm["isTestLinkSent"]; ok {
		util.AddToSlice("`is_test_link_sent`", "=?", "AND", isTestLinkSent, &columnNames, &conditions, &operators, &values)
	}

	// Result search.
	if result, ok := searchForm["result"]; ok {
		util.AddToSlice("`result`", "=?", "AND", result, &columnNames, &conditions, &operators, &values)
	}

	// Academics related search.
	degreeID, degreeIDok := searchForm["degreeID"]
	collegeBranchID, collegeBranchIDok := searchForm["collegeBranchID"]
	specializationID, specializationIDok := searchForm["specializationID"]
	if degreeIDok || collegeBranchIDok || specializationIDok {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN talent_academics ON talents.`id` = talent_academics.`talent_id`"),
			repository.Filter("talent_academics.`deleted_at` IS NULL"),
			repository.Filter("talent_academics.`tenant_id`=?", tenantID),
			repository.GroupBy("talents.`id`"))

		// Degree ID related search.
		if degreeIDok {
			util.AddToSlice("talent_academics.`degree_id`", "IN(?)", "AND", degreeID,
				&columnNames, &conditions, &operators, &values)
		}

		// College branch ID related search.
		if collegeBranchIDok {
			util.AddToSlice("talent_academics.`college_branch_id`", "IN(?)", "AND", collegeBranchID,
				&columnNames, &conditions, &operators, &values)
		}

		// Specialization ID related search.
		if specializationIDok {
			util.AddToSlice("talent_academics.`specialization_id`", "IN(?)", "AND", specializationID,
				&columnNames, &conditions, &operators, &values)
		}
	}

	// Group by campus drive id.
	queryProcesors = append(queryProcesors,
		repository.FilterWithOperator(columnNames, conditions, operators, values))
	return queryProcesors
}

// doForeignKeysExist validates if country, state, degree and specialization present or not in database.
func (service *CandidateService) doForeignKeysExist(candidate *college.Candidate) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if country exists or not.
	if candidate.CountryID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, candidate.TenantID, general.Country{},
			repository.Filter("`id` = ?", candidate.CountryID))
		if err := util.HandleError("Invalid country ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check if state exists or not.
	if candidate.StateID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, candidate.TenantID, general.State{},
			repository.Filter("`id` = ?", candidate.StateID))
		if err := util.HandleError("Invalid state ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check if degree exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, candidate.TenantID, general.Degree{},
		repository.Filter("`id` = ?", candidate.DegreeID))
	if err := util.HandleError("Invalid degree ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if specialization exists or not.
	exists, err = repository.DoesRecordExistForTenant(uow.DB, candidate.TenantID, general.Specialization{},
		repository.Filter("`id`=? AND `degree_id`=?", candidate.SpecializationID, candidate.DegreeID))
	if err := util.HandleError("Invalid specialization ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *CandidateService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *CandidateService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCampusDriveExist validates if campus drive exists or not in database.
func (service *CandidateService) doesCampusDriveExist(tenantID uuid.UUID, campusDriveID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.CampusDrive{},
		repository.Filter("`id` = ?", campusDriveID))
	if err := util.HandleError("Invalid campus drive ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCampusTalentRegistrationExist validates if candidate exists or not in database.
func (service *CandidateService) doesCampusTalentRegistrationExist(campusTalRegID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.CampusTalentRegistration{},
		repository.Filter("`id` = ?", campusTalRegID))
	if err := util.HandleError("Invalid Canmpus Talent Registration ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesEmailExist checks if candidate is already present in campus drive.
func (service *CandidateService) doesEmailExist(candidate *college.Candidate) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check email exists or not.
	exists, err := repository.DoesRecordExist(uow.DB, &college.CampusTalentRegistration{},
		repository.Join("JOIN talents ON talents.`id` = campus_talent_registrations.`talent_id`"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("campus_talent_registrations.`tenant_id`=?", candidate.TenantID),
		repository.Filter("talents.`tenant_id`=?", candidate.TenantID),
		repository.Filter("`campus_drive_id`=?", candidate.CampusDriveID),
		repository.Filter("`email`=? AND campus_talent_registrations.`talent_id`!=?", candidate.Email, candidate.TalentID))

	errorString := "Candidate already exists with email id " + candidate.Email
	if err := util.HandleIfExistsError(errorString, exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// createTalent creates talent struct.
func (service *CandidateService) createTalent(talent *tal.Talent, candidate *college.Candidate) {
	// Create talent struct.
	talent.TenantID = candidate.TenantID
	talent.ID = candidate.TalentID
	talent.CreatedBy = candidate.CreatedBy
	talent.UpdatedBy = candidate.UpdatedBy
	talent.DeletedBy = candidate.DeletedBy
	talent.FirstName = candidate.FirstName
	talent.LastName = candidate.LastName
	talent.Email = candidate.Email
	talent.Contact = candidate.Contact
	talent.AcademicYear = candidate.AcademicYear
	talent.IsSwabhavTalent = candidate.IsSwabhavTalent
	talent.CountryID = candidate.CountryID
	talent.StateID = candidate.StateID
	talent.Address = candidate.Address
	talent.PINCode = candidate.PINCode
	talent.City = candidate.City
	talent.Resume = candidate.Resume
	// Set is_active field of talent as true.
	talent.IsActive = true

	// Set is_experience field of talent as false.
	talent.IsExperience = false

	// Create talent academic struct.
	talentAcademic := &tal.Academic{}
	talentAcademic.CollegeID = candidate.CollegeID
	talentAcademic.College = candidate.College
	talentAcademic.DegreeID = candidate.DegreeID
	talentAcademic.SpecializationID = candidate.SpecializationID
	talentAcademic.Passout = candidate.Passout
	talentAcademic.Percentage = candidate.Percentage
	talentAcademic.ID = candidate.TalentAcademicID

	// Give talent academic to talent.
	talent.Academics = []*tal.Academic{}
	talent.Academics = append(talent.Academics, talentAcademic)
}

// createCampusTalentRegistration creates campus talent registration struct.
func (service *CandidateService) createCampusTalentRegistration(campusTalReg *college.CampusTalentRegistration,
	candidate *college.Candidate) {
	campusTalReg.ID = candidate.CampusTalentRegistrationID
	campusTalReg.TenantID = candidate.TenantID
	campusTalReg.CreatedBy = candidate.CreatedBy
	campusTalReg.UpdatedBy = candidate.UpdatedBy
	campusTalReg.DeletedBy = candidate.DeletedBy
	campusTalReg.RegistrationDate = candidate.RegistrationDate
	campusTalReg.IsTestLinkSent = candidate.IsTestLinkSent
	campusTalReg.HasAttempted = candidate.HasAttempted
	campusTalReg.Result = candidate.Result
	campusTalReg.CampusDriveID = candidate.CampusDriveID
	campusTalReg.ID = candidate.CampusTalentRegistrationID
	campusTalReg.TalentID = candidate.TalentID
}

// createCampusTalentRegistrationPartOfCandidate creates campus talent registration part of candidate.
func (service *CandidateService) createCampusTalentRegistrationPartOfCandidate(campusTalReg *college.CampusTalentRegistration,
	candidate *college.CandidateDTO) {
	candidate.RegistrationDate = campusTalReg.RegistrationDate
	candidate.IsTestLinkSent = campusTalReg.IsTestLinkSent
	candidate.HasAttempted = campusTalReg.HasAttempted
	candidate.Result = campusTalReg.Result
	candidate.CampusDriveID = campusTalReg.CampusDriveID
	candidate.CampusTalentRegistrationID = campusTalReg.ID
	candidate.TalentID = campusTalReg.TalentID
}

// createTalentPartOfCandidate creates talent part of candidate.
func (service *CandidateService) createTalentPartOfCandidate(talent *tal.DTO, candidate *college.CandidateDTO) {
	// Create personal details part talent for candidate.
	candidate.TalentID = talent.ID
	candidate.FirstName = talent.FirstName
	candidate.LastName = talent.LastName
	candidate.Email = talent.Email
	candidate.Contact = talent.Contact
	candidate.AcademicYear = talent.AcademicYear
	candidate.IsSwabhavTalent = talent.IsSwabhavTalent
	candidate.CountryID = talent.CountryID
	candidate.Country = talent.Country
	candidate.StateID = talent.StateID
	candidate.State = talent.State
	candidate.Address = talent.Address
	candidate.PINCode = talent.PINCode
	candidate.City = talent.City
	candidate.Resume = talent.Resume

	// Create talent academic struct.
	talentAcademic := talent.Academics[len(talent.Academics)-1]
	candidate.College = talentAcademic.College
	candidate.Degree = talentAcademic.Degree
	candidate.Specialization = talentAcademic.Specialization
	candidate.Passout = talentAcademic.Passout
	candidate.Percentage = talentAcademic.Percentage
	candidate.TalentAcademicID = talentAcademic.ID
}

// sortTalentAcademicsByOrder sorts all the talent academics by order of Passout field.
func (service *CandidateService) sortTalentAcademicsByOrder(talents *[]tal.DTO) {
	if talents != nil && len(*talents) != 0 {
		for i := 0; i < len(*talents); i++ {
			academics := &(*talents)[i].Academics
			for j := 0; j < len(*academics); j++ {
				if (*academics)[j].Passout == 0 {
					return
				}
			}
			for j := 0; j < len(*academics); j++ {
				sort.Slice(*academics, func(p, q int) bool {
					return (*academics)[p].Passout < (*academics)[q].Passout
				})
			}
		}
	}
}

// setCollegeNameAndID gets the college by college name and sets the college id.
func (service *CandidateService) setCollegeNameAndID(uow *repository.UnitOfWork, candidate *college.Candidate) error {
	// Get college branch from database.
	collegeBranch := list.Branch{}
	if err := service.Repository.GetRecordForTenant(uow, candidate.TenantID, &collegeBranch,
		repository.Filter("`branch_name`=?", candidate.College)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}

	// Give college id to academics collegeID field.
	candidate.CollegeID = &collegeBranch.ID
	return nil
}
