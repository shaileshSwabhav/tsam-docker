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

// StudentService provides method to update, delete, add, get all, get one for student.
type StudentService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// seminarTalentAssociations provides preload associations array for talent.
var seminarTalentAssociations []string = []string{"Academics", "Academics.Degree", "Academics.Specialization", "State", "Country"}

// NewStudentService returns new instance of StudentService.
func NewStudentService(db *gorm.DB, repository repository.Repository) *StudentService {
	return &StudentService{
		DB:         db,
		Repository: repository,
	}
}

// AddStudentWithCredential adds one student to database with credential id.
func (service *StudentService) AddStudentWithCredential(student *college.Student) error {
	// Get credential id from CreatedBy field of student(set in controller).
	credentialID := student.CreatedBy

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, student.TenantID); err != nil {
		return err
	}

	// Call add student service method.
	if err := service.AddStudent(student); err != nil {
		return err
	}

	return nil
}

// AddStudent adds one student to database.
func (service *StudentService) AddStudent(student *college.Student) error {

	// Validate tenant id.
	if err := service.doesTenantExist(student.TenantID); err != nil {
		return err
	}

	// Validate seminar id.
	if err := service.doesSeminarExist(student.TenantID, student.SeminarID); err != nil {
		return err
	}

	// Check if same email exists in seminar.
	if err := service.doesEmailExist(student); err != nil {
		return err
	}

	// Give current date to registration field of student.
	currentTime := time.Now()
	student.RegistrationDate = currentTime.Format("2006-01-02")

	// Validate compulsary fields.
	err := student.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get source from database with name as 'semi'.
	tempSource := general.Source{}
	if err := service.Repository.GetRecordForTenant(uow, student.TenantID, &tempSource,
		repository.Filter("`name`=?", "semi")); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}

	// Set college id by college name.
	if err := service.setCollegeNameAndID(uow, student); err != nil {
		uow.RollBack()
		return err
	}

	// Create talent.
	talent := tal.Talent{}
	service.createTalent(&talent, student)

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
		if err := service.Repository.GetRecordForTenant(uow, student.TenantID, &tempTalent,
			repository.Filter("`email`=?", student.Email),
			repository.PreloadAssociations([]string{"Academics"})); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewValidationError("Student could not be added")
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
		tempTalent.UpdatedBy = student.CreatedBy

		// Store academics of temp talent in a bucket.
		tempTalentAcademics := tempTalent.Academics

		// Make temp talent academics nil to avoid any unnecessary updates.
		tempTalent.Academics = nil

		// Update temp talent with the new fields.
		if err := service.Repository.Save(uow, tempTalent); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Student could not be added", http.StatusInternalServerError)
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
				tempTalentAcademics[i].UpdatedBy = student.CreatedBy

				// Update temp talent academic.
				if err := service.Repository.Save(uow, tempTalentAcademics[i]); err != nil {
					log.NewLogger().Error(err.Error())
					uow.RollBack()
					return errors.NewHTTPError("Student could not be added", http.StatusInternalServerError)
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
			newAcademic.TenantID = student.TenantID
			newAcademic.CreatedBy = student.CreatedBy

			// Give talent id to new academic.
			newAcademic.TalentID = tempTalent.ID

			// Add new academic.
			if err := service.Repository.Add(uow, &newAcademic); err != nil {
				log.NewLogger().Error(err.Error())
				return errors.NewHTTPError("Student could not be added", http.StatusInternalServerError)
			}
		}
	}

	// Create seminar talent registration.
	seminarTalReg := college.SeminarTalentRegistration{}
	service.createSeminarTalentRegistration(&seminarTalReg, student)

	// If talent already existed then give the temp talent id, else give the new talent id.
	if doesTalentExist {
		seminarTalReg.TalentID = tempTalent.ID
	} else {
		seminarTalReg.TalentID = talent.ID
	}

	// Add seminar talent registration to database.
	if err := service.Repository.Add(uow, &seminarTalReg); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Student could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetStudents gets all students from database.
func (service *StudentService) GetStudents(students *[]college.StudentDTO, tenantID uuid.UUID, seminarID uuid.UUID,
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

	// Validate seminar id.
	if err := service.doesSeminarExist(tenantID, seminarID); err != nil {
		return err
	}

	// Create bucket for talents and seminar talent registrations.
	talents := []tal.DTO{}
	seminarTalRegs := []college.SeminarTalentRegistration{}

	// Query processors for search and get for seminar talent registrations.
	var queryProcessorsForSeminarTalReg []repository.QueryProcessor
	queryProcessorsForSeminarTalReg = append(queryProcessorsForSeminarTalReg,
		repository.Join("JOIN talents ON talents.`id` = seminar_talent_registrations.`talent_id`"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("seminar_talent_registrations.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Filter("`seminar_id`=?", seminarID),
		repository.Paginate(limit, offset, totalCount),
	)
	queryProcessorsForSeminarTalReg = append(queryProcessorsForSeminarTalReg, service.addSearchQueries(searchForm, tenantID)...)

	// Get all seminar talent registrations from database.
	if err := service.Repository.GetAllInOrder(uow, &seminarTalRegs, "`first_name`, `last_name`", queryProcessorsForSeminarTalReg...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for search and get for talents.
	queryProcessorsForTalent := service.addSearchQueries(searchForm, tenantID)
	queryProcessorsForTalent = append(queryProcessorsForTalent,
		repository.Join("JOIN seminar_talent_registrations ON talents.`id` = seminar_talent_registrations.`talent_id`"),
		repository.Filter("seminar_talent_registrations.`deleted_at` IS NULL"),
		repository.Filter("seminar_talent_registrations.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Filter("`seminar_id`=?", seminarID),
		repository.PreloadAssociations(seminarTalentAssociations),
		repository.Paginate(limit, offset, totalCount),
	)

	// Get talents from database.
	if err := service.Repository.GetAllInOrder(uow, &talents, "`first_name`, `last_name`", queryProcessorsForTalent...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give seminar talent registrations to students.
	for _, seminarTalReg := range seminarTalRegs {
		student := college.StudentDTO{}
		service.createSeminarTalentRegistrationPartOfStudent(&seminarTalReg, &student)
		*students = append(*students, student)
	}

	// Sort talent academics.
	service.sortTalentAcademicsByOrder(&talents)

	// Give talents to students.
	for _, talent := range talents {
		for i := 0; i < len(*students); i++ {
			if talent.ID == (*students)[i].TalentID {
				service.createTalentPartOfStudent(&talent, &(*students)[i])
			}
		}
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}

	return nil
}

// UpdateStudent updates student in Database.
func (service *StudentService) UpdateStudent(student *college.Student) error {
	// Get credential id from UpdatedBy field of student(set in controller).
	credentialID := student.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(student.TenantID); err != nil {
		return err
	}

	// Validate student id.
	if err := service.doesSeminarTalentRegistrationExist(student.SeminarTalentRegistrationID, student.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, student.TenantID); err != nil {
		return err
	}

	// Validate seminar id.
	if err := service.doesSeminarExist(student.TenantID, student.SeminarID); err != nil {
		return err
	}

	// Check if same email exists in seminar.
	if err := service.doesEmailExist(student); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Set college id by college name.
	if err := service.setCollegeNameAndID(uow, student); err != nil {
		uow.RollBack()
		return err
	}

	// Create bucket for getting student already present in database.
	tempSeminarTalReg := college.SeminarTalentRegistration{}

	// Get student for getting created_by field of student from database.
	if err := service.Repository.GetForTenant(uow, student.TenantID, student.SeminarTalentRegistrationID, &tempSeminarTalReg); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp student to student to be updated.
	student.CreatedBy = tempSeminarTalReg.CreatedBy

	// Create talent.
	talent := tal.Talent{}
	service.createTalent(&talent, student)

	// Check if talent with email already exists.
	isTalentRecordNotFound := false
	tempTalent := tal.Talent{}
	if err := service.Repository.GetRecordForTenant(uow, student.TenantID, &tempTalent,
		repository.Filter("`email`=? AND `id`!=?", student.Email, talent.ID)); err != nil {
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

	// Create seminar talent registration.
	seminarTalReg := college.SeminarTalentRegistration{}
	service.createSeminarTalentRegistration(&seminarTalReg, student)

	// Add seminar talent registration to database.
	if err := service.Repository.Save(uow, &seminarTalReg); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Seminar talent registration could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteStudent deletes one student form database.
func (service *StudentService) DeleteStudent(student *college.Student, uows ...*repository.UnitOfWork) error {
	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Get credential id from DeletedBy field of student(set in controller).
	credentialID := student.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(student.TenantID); err != nil {
		return err
	}

	// Validate student id.
	if err := service.doesSeminarTalentRegistrationExist(student.SeminarTalentRegistrationID, student.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, student.TenantID); err != nil {
		return err
	}

	// Validate seminar id.
	if err := service.doesSeminarExist(student.TenantID, student.SeminarID); err != nil {
		return err
	}

	// Update student for updating deleted_by and deleted_at field of student.
	if err := service.Repository.UpdateWithMap(uow, &college.SeminarTalentRegistration{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`id`=? AND `tenant_id`=?", student.SeminarTalentRegistrationID, student.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Seminar Talent Registration could not be deleted", http.StatusInternalServerError)
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// UpdateMultipleStudent updates multiple fields of seminar talent registrations.
func (service *StudentService) UpdateMultipleStudent(updateMultipleStudent *college.UpdateMultipleStudent) error {
	// Validate tenant id.
	if err := service.doesTenantExist(updateMultipleStudent.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(updateMultipleStudent.UpdatedBy, updateMultipleStudent.TenantID); err != nil {
		return err
	}

	// Validate seminar.
	if err := service.doesSeminarExist(updateMultipleStudent.TenantID, updateMultipleStudent.SeminarID); err != nil {
		return err
	}

	// Validate all seminar talent registrations.
	for _, seminarTalentRegistrationID := range updateMultipleStudent.SeminarTalentRegistrationIDs {
		if err := service.doesSeminarTalentRegistrationExist(seminarTalentRegistrationID, updateMultipleStudent.TenantID); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create map for upadting values.
	updateMap := make(map[string]interface{})

	if updateMultipleStudent.HasVisited != nil {
		updateMap["HasVisited"] = *updateMultipleStudent.HasVisited
	}

	// Update seminar talent registrayion field of students.
	if err := service.Repository.UpdateWithMap(uow, &college.SeminarTalentRegistration{}, updateMap,
		repository.Filter("`id` IN (?)", updateMultipleStudent.SeminarTalentRegistrationIDs)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Student(s) could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds all search queries if any when getAll is called.
func (service *StudentService) addSearchQueries(searchForm url.Values, tenantID uuid.UUID) []repository.QueryProcessor {
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

	// Has visited search.
	if hasVisited, ok := searchForm["hasVisited"]; ok {
		util.AddToSlice("`has_visited`", "=?", "AND", hasVisited, &columnNames, &conditions, &operators, &values)
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

	// Group by seminar id.
	queryProcesors = append(queryProcesors,
		repository.FilterWithOperator(columnNames, conditions, operators, values))
	return queryProcesors
}

// doForeignKeysExist validates if country, state, degree and specialization present or not in database.
func (service *StudentService) doForeignKeysExist(student *college.Student) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if country exists or not.
	if student.CountryID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, student.TenantID, general.Country{},
			repository.Filter("`id` = ?", student.CountryID))
		if err := util.HandleError("Invalid country ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check if state exists or not.
	if student.StateID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, student.TenantID, general.State{},
			repository.Filter("`id` = ?", student.StateID))
		if err := util.HandleError("Invalid state ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check if degree exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, student.TenantID, general.Degree{},
		repository.Filter("`id` = ?", student.DegreeID))
	if err := util.HandleError("Invalid degree ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if specialization exists or not.
	exists, err = repository.DoesRecordExistForTenant(uow.DB, student.TenantID, general.Specialization{},
		repository.Filter("`id`=? AND `degree_id`=?", student.SpecializationID, student.DegreeID))
	if err := util.HandleError("Invalid specialization ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *StudentService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *StudentService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesSeminarExist validates if seminar exists or not in database.
func (service *StudentService) doesSeminarExist(tenantID uuid.UUID, seminarID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.Seminar{},
		repository.Filter("`id` = ?", seminarID))
	if err := util.HandleError("Invalid seminar ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesSeminarTalentRegistrationExist validates if student exists or not in database.
func (service *StudentService) doesSeminarTalentRegistrationExist(seminarTalRegID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.SeminarTalentRegistration{},
		repository.Filter("`id` = ?", seminarTalRegID))
	if err := util.HandleError("Invalid Seminar Talent Registration ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesEmailExist checks if student is already present in seminar.
func (service *StudentService) doesEmailExist(student *college.Student) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check email exists or not.
	exists, err := repository.DoesRecordExist(uow.DB, &college.SeminarTalentRegistration{},
		repository.Join("JOIN talents ON talents.`id` = seminar_talent_registrations.`talent_id`"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("seminar_talent_registrations.`tenant_id`=?", student.TenantID),
		repository.Filter("talents.`tenant_id`=?", student.TenantID),
		repository.Filter("`seminar_id`=?", student.SeminarID),
		repository.Filter("`email`=? AND seminar_talent_registrations.`talent_id`!=?", student.Email, student.TalentID))

	errorString := "Student already exists with email id " + student.Email
	if err := util.HandleIfExistsError(errorString, exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// createTalent creates talent struct.
func (service *StudentService) createTalent(talent *tal.Talent, student *college.Student) {
	// Create talent struct.
	talent.TenantID = student.TenantID
	talent.ID = student.TalentID
	talent.CreatedBy = student.CreatedBy
	talent.UpdatedBy = student.UpdatedBy
	talent.DeletedBy = student.DeletedBy
	talent.FirstName = student.FirstName
	talent.LastName = student.LastName
	talent.Email = student.Email
	talent.Contact = student.Contact
	talent.AcademicYear = student.AcademicYear
	talent.IsSwabhavTalent = student.IsSwabhavTalent
	talent.CountryID = student.CountryID
	talent.StateID = student.StateID
	talent.Address = student.Address
	talent.PINCode = student.PINCode
	talent.City = student.City
	talent.Resume = student.Resume
	// Set is_active field of talent as true.
	talent.IsActive = true

	// Set is_experience field of talent as false.
	talent.IsExperience = false

	// Create talent academic struct.
	talentAcademic := &tal.Academic{}
	talentAcademic.CollegeID = student.CollegeID
	talentAcademic.College = student.College
	talentAcademic.DegreeID = student.DegreeID
	talentAcademic.SpecializationID = student.SpecializationID
	talentAcademic.Passout = student.Passout
	talentAcademic.Percentage = student.Percentage
	talentAcademic.ID = student.TalentAcademicID

	// Give talent academic to talent.
	talent.Academics = []*tal.Academic{}
	talent.Academics = append(talent.Academics, talentAcademic)
}

// createSeminarTalentRegistration creates seminar talent registration struct.
func (service *StudentService) createSeminarTalentRegistration(seminarTalReg *college.SeminarTalentRegistration,
	student *college.Student) {
	seminarTalReg.ID = student.SeminarTalentRegistrationID
	seminarTalReg.TenantID = student.TenantID
	seminarTalReg.CreatedBy = student.CreatedBy
	seminarTalReg.UpdatedBy = student.UpdatedBy
	seminarTalReg.DeletedBy = student.DeletedBy
	seminarTalReg.RegistrationDate = student.RegistrationDate
	seminarTalReg.HasVisited = student.HasVisited
	seminarTalReg.SeminarID = student.SeminarID
	seminarTalReg.ID = student.SeminarTalentRegistrationID
	seminarTalReg.TalentID = student.TalentID
}

// createSeminarTalentRegistrationPartOfStudent creates seminar talent registration part of student.
func (service *StudentService) createSeminarTalentRegistrationPartOfStudent(seminarTalReg *college.SeminarTalentRegistration,
	student *college.StudentDTO) {
	student.RegistrationDate = seminarTalReg.RegistrationDate
	student.HasVisited = seminarTalReg.HasVisited
	student.SeminarID = seminarTalReg.SeminarID
	student.SeminarTalentRegistrationID = seminarTalReg.ID
	student.TalentID = seminarTalReg.TalentID
}

// createTalentPartOfStudent creates talent part of student.
func (service *StudentService) createTalentPartOfStudent(talent *tal.DTO, student *college.StudentDTO) {
	// Create personal details part talent for student.
	student.TalentID = talent.ID
	student.FirstName = talent.FirstName
	student.LastName = talent.LastName
	student.Email = talent.Email
	student.Contact = talent.Contact
	student.AcademicYear = talent.AcademicYear
	student.IsSwabhavTalent = talent.IsSwabhavTalent
	student.CountryID = talent.CountryID
	student.Country = talent.Country
	student.StateID = talent.StateID
	student.State = talent.State
	student.Address = talent.Address
	student.PINCode = talent.PINCode
	student.City = talent.City
	student.Resume = talent.Resume

	// Create talent academic struct.
	talentAcademic := talent.Academics[len(talent.Academics)-1]
	student.College = talentAcademic.College
	student.Degree = talentAcademic.Degree
	student.Specialization = talentAcademic.Specialization
	student.Passout = talentAcademic.Passout
	student.Percentage = talentAcademic.Percentage
	student.TalentAcademicID = talentAcademic.ID
}

// sortTalentAcademicsByOrder sorts all the talent academics by order of Passout field.
func (service *StudentService) sortTalentAcademicsByOrder(talents *[]tal.DTO) {
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
func (service *StudentService) setCollegeNameAndID(uow *repository.UnitOfWork, student *college.Student) error {
	// Get college branch from database.
	collegeBranch := list.Branch{}
	if err := service.Repository.GetRecordForTenant(uow, student.TenantID, &collegeBranch,
		repository.Filter("`branch_name`=?", student.College)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}

	// Give college id to academics collegeID field.
	student.CollegeID = &collegeBranch.ID
	return nil
}
