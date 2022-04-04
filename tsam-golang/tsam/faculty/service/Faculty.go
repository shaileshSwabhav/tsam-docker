package service

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	genService "github.com/techlabs/swabhav/tsam/general/service"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/course"
	fct "github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// FacultyService Provide method to Update, Delete, Add, Get Method For Faculty.
type FacultyService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// FacultyAssociationNames Preload Association Array For Faculty
var FacultyAssociationNames []string = []string{
	"Academics", "Experiences", "Technologies", "Experiences.Technologies", "Academics.Specialization",
	"Experiences.Designation", "Academics.Degree", "Country", "State",
}

// NewFacultyService creates a new instance of FacultyService
func NewFacultyService(db *gorm.DB, repository repository.Repository) *FacultyService {
	return &FacultyService{
		DB:         db,
		Repository: repository,
	}
}

// AddFaculty Add New Faculty to Database.
func (service *FacultyService) AddFaculty(faculty *fct.Faculty, uows ...*repository.UnitOfWork) error {

	// Validates tenantID for faculty
	err := service.doesTenantExists(faculty.TenantID)
	if err != nil {
		//
		return err
	}

	// check if credential exist
	// err = facultyService.doesCredentialExist(faculty.TenantID, faculty.CreatedBy)
	// if err != nil {
	// 	return err
	// }

	// creating login service
	credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// // Checks if credentialID has permission to add new faculty
	// err = credentialService.ValidatePermission(faculty.TenantID, faculty.CreatedBy, "/admin/employee/faculty", "add")
	// if err != nil {
	// 	//
	// 	return errors.NewValidationError(err.Error())
	// }

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Check if same email exists
	exists, err := repository.DoesRecordExistForTenant(service.DB, faculty.TenantID, faculty,
		repository.Filter("email=?", faculty.Email))
	if err != nil {
		//
		return errors.NewValidationError(err.Error())
	}
	if exists {
		// create credential for faculty and then prompt error for email already exists
		// get faculty existing with same email
		existingFaculty := fct.Faculty{}
		err := service.Repository.GetRecordForTenant(uow, faculty.TenantID, &existingFaculty,
			repository.Filter("email=?", faculty.Email))
		if err != nil {
			//
			uow.RollBack()
			return err
		}
		// existingFaculty.Password = util.GeneratePassword()
		// existingFaculty.Password = faculty.FirstName
		existingFaculty.Password = strings.ToLower(existingFaculty.FirstName)

		// Create login for faculty
		err = service.addFacultyCredential(uow, credentialService, &existingFaculty)
		if err != nil {
			//
			uow.RollBack()
			return err
		}
		// update password for existing record
		err = service.Repository.UpdateWithMap(uow, &existingFaculty, map[interface{}]interface{}{
			"password": existingFaculty.Password,
		})
		if err != nil {
			//
			uow.RollBack()
			return err
		}
		if length == 0 {
			uow.Commit()
		}

		// log.NewLogger().Error("Email already exists.")
		return errors.NewValidationError("Email already exists.")
	}

	// extract id from faculty
	err = service.extractAndCheckAllID(faculty)
	if err != nil {
		//
		return err
	}

	// calucaltes the years of experience using fromDate and toDate
	// if faculty.Experiences != nil { //  && len(faculty.Experiences) != 0
	// 	err = facultyService.caluclateDate(faculty)
	// 	if err != nil {
	//
	// 		return err
	// 	}
	// }

	service.addDateToExperience(faculty)

	// Assign faculty code
	faculty.Code, err = util.GenerateUniqueCode(uow.DB, faculty.FirstName,
		"`code` = ?", faculty)
	if err != nil {
		//
		return errors.NewHTTPError("Fail to generate Faculty Code", http.StatusInternalServerError)
	}

	// generate random password
	faculty.Password = strings.ToLower(faculty.FirstName)

	// Add faculty to DB
	err = service.Repository.Add(uow, faculty)
	if err != nil {
		//
		uow.RollBack()
		return err
	}

	// Create login for faculty
	err = service.addFacultyCredential(uow, credentialService, faculty)
	if err != nil {
		//
		uow.RollBack()
		return err
	}
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddFaculties Add Multiple Faculty to Database.
func (service *FacultyService) AddFaculties(faculties *[]fct.Faculty, facultiesIDs *[]uuid.UUID, tenantID,
	credentialID uuid.UUID) error {

	// Add individual Faculty To Database
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, faculty := range *faculties {
		faculty.TenantID = tenantID
		faculty.CreatedBy = credentialID
		err := service.AddFaculty(&faculty, uow)
		if err != nil {
			//
			return err
		}
		*facultiesIDs = append(*facultiesIDs, faculty.ID)
	}
	uow.Commit()
	return nil
}

// UpdateFaculty Add New Faculty to Database.
func (service *FacultyService) UpdateFaculty(faculty *fct.Faculty) error {

	// checks if tenant exists in DB
	err := service.doesTenantExists(faculty.TenantID)
	if err != nil {
		//
		return err
	}

	// check if credential exist
	// err = facultyService.doesCredentialExist(faculty.TenantID, faculty.UpdatedBy)
	// if err != nil {
	// 	return err
	// }

	// extract id from faculty
	err = service.extractAndCheckAllID(faculty)
	if err != nil {
		//
		return err
	}

	// Check if faculty exists in DB
	err = service.doesFacultyExists(faculty.TenantID, faculty.ID)
	if err != nil {
		//
		return err
	}

	// calculates years of experience using fromDate and toDate
	// if faculty.Experiences != nil { //&& len(faculty.Experiences) != 0
	// 	err = facultyService.caluclateDate(faculty)
	// 	if err != nil {
	//
	// 		return err
	// 	}
	// }

	// Create service for credential
	// credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// // Checks if loginID has permission to update faculty
	// err = credentialService.ValidatePermission(faculty.TenantID, faculty.UpdatedBy, "/admin/employee/faculty", "update")
	// if err != nil {
	// 	//
	// 	return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	// }

	// trim time paramters from date of birth and date of joining
	// if faculty.DateOfBirth != nil {
	// 	faculty.DateOfBirth = util.RemoveTimeStampFromDate(*faculty.DateOfBirth)
	// 	// fmt.Println("*****dob ->", *faculty.DateOfBirth)
	// }
	// if faculty.DateOfJoining != nil {
	// 	faculty.DateOfJoining = util.RemoveTimeStampFromDate(*faculty.DateOfJoining)
	// 	// fmt.Println("*****dob ->", *faculty.DateOfJoining)
	// }

	service.addDateToExperience(faculty)

	// Start transaction
	uow := repository.NewUnitOfWork(service.DB, false)
	// get faculty password
	tempFaculty := &fct.Faculty{}
	err = service.Repository.GetRecordForTenant(uow, faculty.TenantID, tempFaculty,
		repository.Select([]string{"`created_by`", "`password`, `code`"}), repository.Filter("`id` = ?", faculty.ID))
	if err != nil {
		//
		uow.RollBack()
		return err
	}

	// fmt.Println("*****password ->", tempFaculty.Password)
	faculty.CreatedBy = tempFaculty.CreatedBy
	faculty.Code = tempFaculty.Code
	faculty.Password = tempFaculty.Password

	// update academics of faculty
	err = service.updateFacultyAcadmecis(faculty, uow)
	if err != nil {
		//
		uow.RollBack()
		return err
	}
	// update experince of faculty
	err = service.updateFacultyExperince(faculty, uow)
	if err != nil {
		//
		uow.RollBack()
		return err
	}

	// Update faculty asssociations(technology)
	err = service.updateFacultyAssociation(uow, faculty)
	if err != nil {

		uow.RollBack()
		return err
	}

	// Update Faculty in Database
	err = service.Repository.Save(uow, faculty)
	if err != nil {

		uow.RollBack()
		return err
	}

	// update is_active field of faculty credential.
	err = service.updateFacultyCredential(uow, faculty)
	if err != nil {

		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteFaculty soft deletes record from database
func (service *FacultyService) DeleteFaculty(faculty *fct.Faculty) error {

	// check if tenant exist
	err := service.doesTenantExists(faculty.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	// err = facultyService.doesCredentialExist(faculty.TenantID, faculty.DeletedBy)
	// if err != nil {
	// 	return err
	// }

	// check if faculty exist
	err = service.doesFacultyExists(faculty.TenantID, faculty.ID)
	if err != nil {
		return err
	}

	// Create service for credential
	// credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// checks if loginID has permission to delete the faculty
	// err = credentialService.ValidatePermission(faculty.TenantID, faculty.DeletedBy, "/admin/employee/faculty", "delete")
	// if err != nil {
	// 	return err
	// }

	// allow faculty delete only if batches are not assigned
	exists, err := repository.DoesRecordExistForTenant(service.DB, faculty.TenantID, &bat.Batch{},
		repository.Filter("`faculty_id`=?", faculty.ID))
	if err != nil {
		return err
	}
	if exists {
		return errors.NewValidationError("Faculty cannot be deleted as they are assigned to a batch")
	}

	credential := general.Credential{
		FacultyID: &faculty.ID,
	}
	credential.DeletedBy = faculty.DeletedBy

	// Start transaction
	uow := repository.NewUnitOfWork(service.DB, false)

	// Delete login of faculty
	login := genService.NewCredentialService(service.DB, service.Repository)
	err = login.DeleteCredential(&credential, faculty.TenantID, faculty.DeletedBy, *credential.FacultyID, "`faculty_id`=?", uow)
	if err != nil {
		uow.RollBack()
		return err
	}

	// get faculty record
	tempFaculty := &fct.Faculty{}
	err = service.Repository.GetRecordForTenant(uow, faculty.TenantID, tempFaculty,
		repository.Filter("`id`=?", faculty.ID), repository.PreloadAssociations(FacultyAssociationNames))
	if err != nil {
		uow.RollBack()
		return err
	}
	tempFaculty.DeletedBy = faculty.DeletedBy
	// Delete faculty association from DB
	err = service.deleteFacultyAssociation(uow, tempFaculty)
	if err != nil {
		uow.RollBack()
		return err
	}
	// Delete faculty from DB
	// err = facultyService.Repository.Delete(uow, &faculty)
	err = service.Repository.UpdateWithMap(uow, &fct.Faculty{}, map[string]interface{}{
		"DeletedBy": faculty.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", faculty.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Delete faculty credential.
	err = service.Repository.UpdateWithMap(uow, &general.Credential{}, map[string]interface{}{
		"DeletedBy": faculty.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`faculty_id` = ?", faculty.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFacultyList Return Faculty List
func (service *FacultyService) GetFacultyList(faculties *[]*list.Faculty, tenantID uuid.UUID) error {

	// Validates foreign keys of faculty
	err := service.doesTenantExists(tenantID)
	if err != nil {

		return err
	}

	// Start transaction.
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllForTenant(uow, tenantID, faculties,
		repository.Filter("`deleted_at` IS NULL"))
	if err != nil {

		uow.RollBack()
		errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// GetFaculty gets specified faculty from database.
func (service *FacultyService) GetFaculty(faculty *fct.Faculty) error {

	// Validates foreign keys of faculty
	err := service.doesTenantExists(faculty.TenantID)
	if err != nil {

		return err
	}

	// Checks if specified faculty exists
	err = service.doesFacultyExists(faculty.TenantID, faculty.ID)
	if err != nil {

		return err
	}

	// Start transaction.
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetForTenant(uow, faculty.TenantID, faculty.ID, faculty,
		repository.PreloadAssociations(FacultyAssociationNames))
	if err != nil {
		uow.RollBack()

		return err
	}
	uow.Commit()
	return nil
}

// GetFacultyForBatch Return Specific faculty for batch.
func (service *FacultyService) GetFacultyForBatch(faculty *list.Faculty, tenantID, batchID uuid.UUID) error {

	// Check if tenant exists
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Checks if specified batch exists.
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	// Start transaction.
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetRecord(uow, faculty,
		repository.Join("INNER JOIN batch_modules ON batch_modules.`faculty_id` = faculties.`id`"),
		repository.Filter("batch_modules.`tenant_id`=? AND batch_modules.`deleted_at` IS NULL", tenantID),
		repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_modules.`batch_id`=?", batchID))

	if err != nil {
		uow.RollBack()

		return err
	}
	uow.Commit()
	return nil
}

// GetFaculties Return All Faculty From Database.
func (service *FacultyService) GetFaculties(faculty *[]fct.FacultyDTO, tenantID uuid.UUID, parser *web.Parser, totalCount *int) error {

	// check if tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}
	limit, offset := parser.ParseLimitAndOffset()

	var queryProcessors []repository.QueryProcessor
	queryProcessors = append(queryProcessors, service.addFacultyQueriesFromParams(parser.Form)...)
	queryProcessors = append(queryProcessors, repository.OrderBy("faculties.`is_active` DESC"),
		repository.Filter("faculties.`tenant_id` = ?", tenantID), repository.PreloadAssociations(FacultyAssociationNames),
		repository.Paginate(limit, offset, totalCount))

	uow := repository.NewUnitOfWork(service.DB, true)
	// Get Facultys From Database
	err = service.Repository.GetAllInOrder(uow, faculty, "faculties.`first_name`",
		queryProcessors...)
	if err != nil {
		uow.RollBack()

		return err
	}

	err = service.calculateAverageAssessmentScore(faculty, tenantID, uow)
	if err != nil {
		uow.RollBack()
		return err
	}

	// service.removeDateField(faculty)

	uow.Commit()
	return nil
}

// GetFacultyListFromCredential returns faculty list from credential
func (service *FacultyService) GetFacultyListFromCredential(faculty *[]list.FacultyCredentialDTO,
	tenantID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {

		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	// err = service.Repository.GetAllInOrderForTenant(uow, tenantID, faculty, "`first_name`",
	// 	repository.Filter("`faculty_id` IS NOT NULL"))
	err = service.Repository.Scan(uow, faculty,
		repository.Filter("`faculty_id` IS NOT NULL AND `deleted_at` IS NULL"),
		repository.Table("credentials"))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetSearchedFaculties Return All Faculty From Database.
// func (service *FacultyService) GetSearchedFaculties(faculty *[]fct.Faculty, searchFaculty *fct.Search, tenantID uuid.UUID, limit, offset int, totalCount *int) error {

// 	// check if tenant exist
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(service.DB, true)
// 	// Get Facultys From Database
// 	err = service.Repository.GetAllForTenant(uow, tenantID, faculty, service.addFacultySearchQueries(searchFaculty),
// 		repository.PreloadAssociations(FacultyAssociationNames), repository.Paginate(limit, offset, totalCount))
// 	if err != nil {
// 		uow.RollBack()
//
// 		return err
// 	}

// 	err = service.calculateAverageAssessmentScore(faculty, tenantID, uow)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	// service.removeDateField(faculty)

// 	uow.Commit()
// 	return nil
// }

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// updateFacultyAssociation Update Faculty's Dependencies
func (service *FacultyService) updateFacultyAssociation(uow *repository.UnitOfWork, faculty *fct.Faculty) error {
	// associations := []interface{}{}
	// associations = append(associations, faculty.Academics)
	// for _, experience := range faculty.Experiences {
	// 	if err := service.Repository.ReplaceAssociations(uow, *experience, FacultyAssociationNames[2], experience.Technologies); err != nil {
	//
	// 		return err
	// 	}
	// }
	// associations = append(associations, faculty.Experiences)
	// associations = append(associations, faculty.Technologies)
	// for index, association := range associations {
	// 	if err := service.Repository.ReplaceAssociations(uow, faculty, FacultyAssociationNames[index], association); err != nil {
	//
	// 		return err
	// 	}
	// }

	// fmt.Println("===============================FACULTY TECHNOLOGIES======================================")

	// replace faculty experience technologies
	for _, experience := range faculty.Experiences {
		// fmt.Println("***updating experience technologies")
		if err := service.Repository.ReplaceAssociations(uow, experience, "Technologies",
			experience.Technologies); err != nil {

			return err
		}
		experience.Technologies = nil
	}
	//replace technologies of faculty
	// fmt.Println("***updating technologies")
	if err := service.Repository.ReplaceAssociations(uow, faculty, "Technologies",
		faculty.Technologies); err != nil {

		return err
	}
	faculty.Technologies = nil
	// fmt.Println("===============================FACULTY TECHNOLOGIES END======================================")
	return nil
}

// deleteFacultyAssociation delete faculty's dependencies from intermediate table.
func (service *FacultyService) deleteFacultyAssociation(uow *repository.UnitOfWork, faculty *fct.Faculty) error {

	err := service.extractAndCheckAllID(faculty)
	if err != nil {
		return err
	}

	// if len(faculty.Academics) > 0 {
	// 	for _, academic := range faculty.Academics {
	err = service.Repository.UpdateWithMap(uow, fct.Academic{}, map[interface{}]interface{}{
		"DeletedBy": faculty.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`faculty_id`=?", faculty.ID))
	if err != nil {
		return err
	}
	// 	}
	// }

	// remove technology associations
	// if len(faculty.Technologies) > 0 {
	// 	if err := service.Repository.RemoveAssociations(uow, faculty, FacultyAssociationNames[2], faculty.Technologies); err != nil {
	// 		return err
	// 	}
	// }

	err = service.Repository.UpdateWithMap(uow, fct.Experience{}, map[interface{}]interface{}{
		"DeletedBy": faculty.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`faculty_id`=?", faculty.ID))
	if err != nil {
		return err
	}

	// remove experience technology assosications
	// for _, experience := range faculty.Experiences {

	// 	err = service.Repository.RemoveAssociations(uow, experience, FacultyAssociationNames[2], experience.Technologies)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// addFacultyCredential creates a record in credentials table.
func (service *FacultyService) addFacultyCredential(uow *repository.UnitOfWork,
	credentialService *genService.CredentialService, faculty *fct.Faculty) error {

	// get roleID for faculty
	role := general.Role{}
	err := service.Repository.GetRecordForTenant(uow, faculty.TenantID, &role,
		repository.Filter("`role_name`=?", "Faculty"))
	if err != nil {
		return err
	}

	// password := "faculty" // change to util.GeneratePassword()
	credentials := general.Credential{
		FirstName: faculty.FirstName,
		LastName:  &faculty.LastName,
		Email:     faculty.Email,
		Password:  faculty.Password,
		FacultyID: &faculty.ID,
		Contact:   faculty.Contact,
		RoleID:    role.ID,
	}
	credentials.TenantID = faculty.TenantID
	credentials.CreatedBy = faculty.CreatedBy
	err = credentialService.AddCredential(&credentials, uow)
	if err != nil {
		return err
	}
	return nil
}

// updateFacultyCredential updates specified faculty record in credentials table.
func (service *FacultyService) updateFacultyCredential(uow *repository.UnitOfWork, faculty *fct.Faculty) error {

	err := service.Repository.UpdateWithMap(uow, &general.Credential{}, map[string]interface{}{
		"UpdatedBy": faculty.UpdatedBy,
		"IsActive":  faculty.IsActive,
	}, repository.Filter("`faculty_id` = ?", faculty.ID))
	if err != nil {
		return err
	}
	return nil
}

// updateFacultyAcadmecis updates the academic records of the faculty
func (service *FacultyService) updateFacultyAcadmecis(faculty *fct.Faculty, uow *repository.UnitOfWork) error {

	tempAcademics := &[]fct.Academic{}
	academicMap := make(map[uuid.UUID]uint)

	// get all academics of current faculty
	err := service.Repository.GetAllForTenant(uow, faculty.TenantID, tempAcademics,
		repository.Filter("`faculty_id` = ?", faculty.ID))
	if err != nil {

		return err
	}

	// populating academicMap
	for _, tempAcademic := range *tempAcademics {
		academicMap[tempAcademic.ID] = 1
	}
	for _, academic := range faculty.Academics {

		if util.IsUUIDValid(academic.ID) {
			academicMap[academic.ID]++
		}

		// check if experience already exists in the DB
		if academicMap[academic.ID] > 1 {
			academic.UpdatedBy = faculty.UpdatedBy
			err = service.Repository.Update(uow, &academic)
			if err != nil {

				return err
			}
			academicMap[academic.ID] = 0
		}

		if !util.IsUUIDValid(academic.ID) {
			// add experience when uuid is nil

			academic.FacultyID = faculty.ID
			academic.TenantID = faculty.TenantID
			academic.CreatedBy = faculty.UpdatedBy
			err := service.Repository.Add(uow, &academic)
			if err != nil {

				return err
			}
		}
	}

	// deleting all records where count is 1 as they have been removed from the experience
	for _, academic := range *tempAcademics {
		if academicMap[academic.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, &academic, map[interface{}]interface{}{
				"DeletedBy": faculty.UpdatedBy,
				"DeletedAt": time.Now(),
			})
			if err != nil {

				return err
			}
		}
		academicMap[academic.ID] = 0
	}

	faculty.Academics = nil

	return nil
}

// updateFacultyExperince updates experiences of the faculty in the DB
func (service *FacultyService) updateFacultyExperince(faculty *fct.Faculty, uow *repository.UnitOfWork) error {

	tempExperiences := &[]fct.Experience{}
	// get all academics of current faculty
	err := service.Repository.GetAllForTenant(uow, faculty.TenantID, tempExperiences,
		repository.Filter("faculty_id=?", faculty.ID))
	if err != nil {

		return err
	}
	experienceMap := make(map[uuid.UUID]uint)

	// populating experinceMap
	for _, tempExperience := range *tempExperiences {
		experienceMap[tempExperience.ID] = 1
	}
	for _, experience := range faculty.Experiences {

		if util.IsUUIDValid(experience.ID) {
			experienceMap[experience.ID]++
		}

		if experienceMap[experience.ID] > 1 && util.IsUUIDValid(experience.ID) {

			experience.UpdatedBy = faculty.UpdatedBy
			experience.FacultyID = faculty.ID

			// replace faculty experience technologies
			if err := service.Repository.ReplaceAssociations(uow, experience, "Technologies",
				experience.Technologies); err != nil {
				return err
			}

			experience.Technologies = nil

			err = service.Repository.Save(uow, &experience)
			if err != nil {

				return err
			}
			experienceMap[experience.ID] = 0
		}

		if !util.IsUUIDValid(experience.ID) {
			experience.FacultyID = faculty.ID
			experience.TenantID = faculty.TenantID
			experience.CreatedBy = faculty.UpdatedBy
			// experince.UpdatedBy = faculty.UpdatedBy
			err := service.Repository.Add(uow, &experience)
			if err != nil {

				return err
			}
		}
	}

	// deleting all records where count is 1 as they have been removed from the experience
	for _, experience := range *tempExperiences {
		if experienceMap[experience.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, experience, map[interface{}]interface{}{
				"DeletedBy": faculty.UpdatedBy,
				"DeletedAt": time.Now(),
			})
			if err != nil {

				return err
			}
		}
		experienceMap[experience.ID] = 0
	}

	// setting experience to nil as all data as been added to db
	faculty.Experiences = nil

	return nil
}

// extractAndCheckAllID extracts ID from the json and make the model nil, also checks if the record exist in the DB
func (service *FacultyService) extractAndCheckAllID(faculty *fct.Faculty) error {

	faculty.CountryID = &faculty.Country.ID
	faculty.StateID = &faculty.State.ID

	// assign tenantID to academic
	for i := range faculty.Academics {
		faculty.Academics[i].TenantID = faculty.TenantID
		faculty.Academics[i].CreatedBy = faculty.CreatedBy
		// extract degreeID
		faculty.Academics[i].DegreeID = faculty.Academics[i].Degree.ID
		// faculty.Academics[i].Degree = nil
		err := service.doesDegreeExists(faculty.TenantID, faculty.Academics[i].DegreeID)
		if err != nil {

			return err
		}
		// extract specializationID
		faculty.Academics[i].SpecializationID = faculty.Academics[i].Specialization.ID
		// faculty.Academics[i].Specialization = nil
		err = service.doesSpecializationExists(faculty.TenantID, faculty.Academics[i].DegreeID,
			faculty.Academics[i].SpecializationID)
		if err != nil {

			return err
		}
	}

	// assign tenantID to experiences
	if faculty.Experiences != nil {
		for i := range faculty.Experiences {

			faculty.Experiences[i].TenantID = faculty.TenantID
			faculty.Experiences[i].CreatedBy = faculty.CreatedBy
			// extract designationID
			faculty.Experiences[i].DesignationID = faculty.Experiences[i].Designation.ID
			// faculty.Experiences[i].Designation = nil
			err := service.doesDesignationExists(faculty.TenantID, faculty.Experiences[i].DesignationID)
			if err != nil {

				return err
			}
			// check if technology exists in DB
			for _, technolgy := range faculty.Experiences[i].Technologies {
				err = service.doesTechnologyExists(faculty.TenantID, technolgy.ID)
				if err != nil {

					return err
				}
			}
		}
	}
	// check if technology exists in DB
	if faculty.Technologies != nil {
		for _, technolgy := range faculty.Technologies {
			err := service.doesTechnologyExists(faculty.TenantID, technolgy.ID)
			if err != nil {

				return err
			}
		}
	}

	return nil
}

// caluclateDate will caluclate years of experience using fromDate and toDate
// func (service *FacultyService) caluclateDate(faculty *fct.Faculty) error {
// 	var time1 time.Time
// 	var time2 time.Time
// 	var year float64
// 	var month float64
// 	var err error
// 	for i := 0; i < len(faculty.Experiences); i++ {

// 		if faculty.Experiences[i].FromDate != nil && len(*faculty.Experiences[i].FromDate) != 0 {
// 			util.AddDateToMonth(faculty.Experiences[i].FromDate)
// 		}
// 		if faculty.Experiences[i].ToDate != nil && len(*faculty.Experiences[i].ToDate) != 0 {
// 			util.AddDateToMonth((faculty.Experiences[i].ToDate))
// 		}
// 		if faculty.Experiences[i].FromDate != nil {
// 			layout1 := "2006-01-02"
// 			time1, err = time.Parse(layout1, *faculty.Experiences[i].FromDate)
// 			if err != nil {
//
// 				return err
// 			}
// 		}
// 		if faculty.Experiences[i].ToDate != nil {
// 			layout1 := "2006-01-02"
// 			time2, err = time.Parse(layout1, *faculty.Experiences[i].ToDate)
// 			if err != nil {
//
// 				return err
// 			}
// 		}
// 		y1, M1, _ := time1.Date()
// 		y2, M2, _ := time2.Date()
// 		year = float64(y2 - y1)
// 		month = float64(M2 - M1)

// 		if year > 0 {
// 			month += 12 * year
// 			year--
// 		}

// 		faculty.Experiences[i].YearOfExperience = float32(math.Round(((month / 12) * 100)) / 100)
// 	}

// 	return nil
// }

// calculateAverageAssessmentScore will calculate average score for the assessment
func (service *FacultyService) calculateAverageAssessmentScore(faculty *[]fct.FacultyDTO, tenantID uuid.UUID,
	uow *repository.UnitOfWork) error {

	for index := range *faculty {

		exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &fct.FacultyAssessment{},
			repository.Filter("`faculty_id`=?", (*faculty)[index].ID))
		if err != nil {
			return err
		}
		if exist {
			err = service.Repository.Scan(uow, &(*faculty)[index],
				repository.Model(&fct.Faculty{}),
				repository.Select("((SUM(feedback_options.`key`) * 10) / SUM(feedback_questions.`max_score`)) AS average_assessment_score"),
				repository.Join("INNER JOIN faculty_assessments ON faculty_assessments.`faculty_id` = faculties.`id`"+
					" AND faculty_assessments.`tenant_id` = faculties.`tenant_id`"),
				repository.Join("INNER JOIN feedback_options ON faculty_assessments.`option_id` = feedback_options.`id`"+
					" AND faculty_assessments.`tenant_id` = feedback_options.`tenant_id`"),
				repository.Join("INNER JOIN feedback_questions ON faculty_assessments.`question_id` = feedback_questions.`id`"+
					" AND faculty_assessments.`tenant_id` = feedback_questions.`tenant_id`"),
				repository.Filter("faculties.`id` = ?", (*faculty)[index].ID), repository.Filter("feedback_questions.`is_active` = true"),
				repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL AND faculty_assessments.`deleted_at` IS NULL", tenantID),
				repository.Filter("feedback_options.`deleted_at` IS NULL AND feedback_questions.`deleted_at` IS NULL"),
				repository.GroupBy("faculty_assessments.`faculty_id`"))
			if err != nil {
				return err
			}
		}

		exist, err = repository.DoesRecordExistForTenant(service.DB, tenantID, &course.CourseTechnicalAssessment{},
			repository.Filter("`faculty_id`=?", (*faculty)[index].ID))
		if err != nil {
			return err
		}
		if exist {
			err = service.Repository.Scan(uow, &(*faculty)[index],
				repository.Model(&fct.Faculty{}), repository.Select("AVG(course_technical_assessments.`rating`) AS average_technical_score"),
				repository.Join("INNER JOIN course_technical_assessments ON course_technical_assessments.`faculty_id`=faculties.`id` AND "+
					"faculties.`tenant_id`=? AND course_technical_assessments.`deleted_at` IS NULL", tenantID),
				repository.Filter("faculties.`id` = ?", (*faculty)[index].ID), repository.GroupBy("course_technical_assessments.`faculty_id`"))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (service *FacultyService) addDateToExperience(faculty *fct.Faculty) error {

	//give 01 to all faculty's experiences' fromDate and toDate field and set createdBy field
	if faculty.Experiences != nil && len(faculty.Experiences) != 0 {
		for i := 0; i < len(faculty.Experiences); i++ {
			//set created_by field of experiences
			faculty.Experiences[i].CreatedBy = faculty.CreatedBy
			if faculty.Experiences[i].FromDate == nil {
				return errors.NewValidationError("Working from date must be specified")
			}

			if len(*(faculty.Experiences[i].FromDate)) != 0 {
				util.AddDateToMonth((faculty.Experiences[i].FromDate))
			}
			if faculty.Experiences[i].ToDate != nil && len(*faculty.Experiences[i].ToDate) != 0 {
				util.AddDateToMonth((faculty.Experiences[i].ToDate))
			}
		}
	}
	return nil
}

// RemoveDateFromFromDateToDate removes date(-01) from FromDate and ToDate fields of talents' experiences
// func (service *FacultyService) removeDateField(faculties *[]fct.Faculty) {
// 	if faculties != nil && len(*faculties) != 0 {
// 		for _, faculty := range *faculties {
// 			if faculty.Experiences != nil {
// 				for _, experience := range faculty.Experiences {
// 					if experience.FromDate != nil {
// 						*experience.FromDate = (*experience.FromDate)[:7]
// 					}
// 					if experience.ToDate != nil {
// 						*experience.ToDate = (*experience.ToDate)[:7]
// 					}

// 				}
// 			}

// 		}
// 	}

// }

// doesTenantExists validates tenantID
func (service *FacultyService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
// func (service *FacultyService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
// 		repository.Filter("`id`=?", credentialID))
// 	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
//
// 		return err
// 	}
// 	return nil
// }

// doesEmailExists checks if duplicate email exists
// func (service *FacultyService) doesEmailExists(tenantID uuid.UUID, email string) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, fct.Faculty{},
// 		repository.Filter("email=?", email))
// 	if err := util.HandleIfExistsError("Email already exists", exists, err); err != nil {
//
// 		return errors.NewValidationError(err.Error())
// 	}
// 	return nil
// }

// doesFacultyExists checks if faculty exists in DB
func (service *FacultyService) doesFacultyExists(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, fct.Faculty{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Faculty not found", exists, err); err != nil {

		return errors.NewValidationError(err.Error())
	}
	return nil
}

// doesDegreeExists checks if degree exists in DB
func (service *FacultyService) doesDegreeExists(tenantID, degreeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Degree{},
		repository.Filter("`id` = ?", degreeID))
	if err := util.HandleError("Degree not found", exists, err); err != nil {

		return errors.NewValidationError(err.Error())
	}
	return nil
}

// doesSpecializationExists checks if specialization exists for specified degree in DB
func (service *FacultyService) doesSpecializationExists(tenantID, degreeID, specializationID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Specialization{},
		repository.Filter("id=? AND degree_id=?", specializationID, degreeID))
	if err := util.HandleError("Specialization not found", exists, err); err != nil {

		return errors.NewValidationError(err.Error())
	}
	return nil
}

// doesDesignationExists checks if degree exists in DB
func (service *FacultyService) doesDesignationExists(tenantID, designationID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Designation{},
		repository.Filter("`id` = ?", designationID))
	if err := util.HandleError("Designation not found", exists, err); err != nil {

		return errors.NewValidationError(err.Error())
	}
	return nil
}

func (service *FacultyService) doesTechnologyExists(tenantID, technologyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Technology{},
		repository.Filter("`id` = ?", technologyID))
	if err := util.HandleError("Technology not found", exists, err); err != nil {

		return errors.NewValidationError(err.Error())
	}
	return nil
}

// doesBatchExists validates if batch exists.
func (service *FacultyService) doesBatchExists(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Batch not found", exists, err); err != nil {

		return errors.NewValidationError(err.Error())
	}
	return nil
}

// addFacultyQueriesFromParams adds all search queries if any when getAll is called
func (service *FacultyService) addFacultyQueriesFromParams(requestForm url.Values) []repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	var queryProcessors []repository.QueryProcessor

	if _, ok := requestForm["city"]; ok {
		util.AddToSlice("city", "LIKE ?", "AND", "%"+requestForm.Get("city")+"%", &columnNames, &conditions, &operators, &values)
	}
	if _, ok := requestForm["contact"]; ok {
		util.AddToSlice("contact", "LIKE ?", "AND", "%"+requestForm.Get("contact")+"%", &columnNames, &conditions, &operators, &values)
	}
	if _, ok := requestForm["email"]; ok {
		util.AddToSlice("email", "LIKE ?", "AND", "%"+requestForm.Get("email")+"%", &columnNames, &conditions, &operators, &values)
	}
	if _, ok := requestForm["firstName"]; ok {
		util.AddToSlice("first_name", "LIKE ?", "AND", "%"+requestForm.Get("firstName")+"%", &columnNames, &conditions, &operators, &values)
	}
	if _, ok := requestForm["lastName"]; ok {
		util.AddToSlice("last_name", "LIKE ?", "AND", "%"+requestForm.Get("lastName")+"%", &columnNames, &conditions, &operators, &values)
	}
	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("is_active", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	//if technologies is present then join talent_next_actions_technologies and technolgies table
	if technologies, ok := requestForm["technologies"]; ok {
		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN faculties_technologies ON faculties.`id` = faculties_technologies.`faculty_id`"))
		if len(technologies) > 0 {
			util.AddToSlice("faculties_technologies.`technology_id`", "IN(?)", "AND", technologies, &columnNames, &conditions, &operators, &values)
		}
	}

	queryProcessors = append(queryProcessors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("faculties.`id`"))

	// return repository.FilterWithOperator(columnNames, conditions, operators, values)
	return queryProcessors
}

// // GetCount Return Total Count Of Result Set.
// func (service *FacultyService) GetCount(faculty interface{}, count *int, queryProcessor ...repository.QueryProcessor) error {
// 	uow := repository.NewUnitOfWork(service.DB, true)
// 	return service.Repository.GetCount(uow, faculty, count, queryProcessor...)
// }

// addFacultySearchQueries adds all search queries if any when getAll is called
// func (service *service) addFacultySearchQueries(searchFaculty *fct.Search) repository.QueryProcessor {
// 	var columnNames []string
// 	var conditions []string
// 	var operators []string
// 	var values []interface{}
// 	if !util.IsEmpty(searchFaculty.City) {
// 		util.AddToSlice("city", "LIKE ?", "AND", "%"+searchFaculty.City+"%", &columnNames,
// 			&conditions, &operators, &values)
// 	}
// 	if !util.IsEmpty(searchFaculty.Contact) {
// 		util.AddToSlice("contact", "LIKE ?", "AND", "%"+searchFaculty.Contact+"%", &columnNames,
// 			&conditions, &operators, &values)
// 	}
// 	if !util.IsEmpty(searchFaculty.Email) {
// 		util.AddToSlice("email", "LIKE ?", "AND", "%"+searchFaculty.Email+"%", &columnNames,
// 			&conditions, &operators, &values)
// 	}
// 	if !util.IsEmpty(searchFaculty.FirstName) {
// 		util.AddToSlice("first_name", "LIKE ?", "AND", "%"+searchFaculty.FirstName+"%",
// 			&columnNames, &conditions, &operators, &values)
// 	}
// 	if !util.IsEmpty(searchFaculty.FirstName) {
// 		util.AddToSlice("last_name", "LIKE ?", "AND", "%"+searchFaculty.LastName+"%",
// 			&columnNames, &conditions, &operators, &values)
// 	}
// 	return repository.FilterWithOperator(columnNames, conditions, operators, values)
// }
