package service

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/models/talent"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	talenq "github.com/techlabs/swabhav/tsam/models/talentenquiry"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// EnquiryService Provide method to update, delete, add, geta all methods for Enquiry.
type EnquiryService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewEnquiryService returns new instance Of EnquiryServcie.
func NewEnquiryService(db *gorm.DB, repository repository.Repository) *EnquiryService {
	return &EnquiryService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Academics", "Experiences", "Technologies", "State", "Country", "Academics.Degree",
			"Experiences.Technologies", "Experiences.Designation", "Academics.Specialization", "SalesPerson",
			"EnquirySource", "MastersAbroad", "MastersAbroad.Scores", "MastersAbroad.Scores.Examination",
			"MastersAbroad.Degree", "MastersAbroad.Universities", "MastersAbroad.Countries", "Courses",
		},
	}
}

// AddEnquiry adds new Enquiry to database.
func (service *EnquiryService) AddEnquiry(enquiry *talenq.Enquiry) error {
	// Validate compulsary fields.
	err := enquiry.Validate(false)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Get credential id from CreatedBy field of enquiry(set in controller).
	credentialID := enquiry.CreatedBy

	// Extract foreign key IDs and remove the object.
	service.extractID(enquiry)

	// Give tenant id to academics and experiences.
	service.setTenantID(enquiry, enquiry.TenantID)

	// Validate tenant id.
	if err := service.doesTenantExist(enquiry.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, enquiry.TenantID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(enquiry); err != nil {
		return err
	}

	// Starting transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Assign Enquiry Code.
	enquiry.Code, err = util.GenerateUniqueCode(uow.DB, enquiry.FirstName, "`code` = ?", talenq.Enquiry{})
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Failed to generate enquiry code", http.StatusInternalServerError)
	}

	// Check if login id's role name is salesperson or not.
	isRecordNotFound := false
	tempUser := general.User{}
	err = service.Repository.GetRecord(uow, &tempUser,
		repository.Select("users.`id`"),
		repository.Join("join credentials on credentials.`sales_person_id` = users.`id`"),
		repository.Filter("credentials.`id`=? AND credentials.`tenant_id`=? AND users.`tenant_id`=?",
			credentialID, enquiry.TenantID, enquiry.TenantID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			isRecordNotFound = true
		} else {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
		}
	}
	if !isRecordNotFound {
		enquiry.SalesPersonID = &tempUser.ID
	}

	// Set college id by college name.
	if err := service.setCollegeNameAndID(uow, enquiry); err != nil {
		uow.RollBack()
		return err
	}

	// Set created_by field of academics.
	if enquiry.Academics != nil && len(enquiry.Academics) != 0 {
		for i := 0; i < len(enquiry.Academics); i++ {
			//set created_by field of academics
			enquiry.Academics[i].CreatedBy = credentialID
		}
	}

	// Give 01 to all enquiry's experiences' fromDate and toDate field and set createdBy field.
	if enquiry.Experiences != nil && len(enquiry.Experiences) != 0 {
		for i := 0; i < len(enquiry.Experiences); i++ {
			//set created_by field of experiences
			enquiry.Experiences[i].CreatedBy = credentialID
			if len(enquiry.Experiences[i].FromDate) != 0 {
				util.AddDateToMonth(&(enquiry.Experiences[i].FromDate))
			}
			if enquiry.Experiences[i].ToDate != nil && len(*enquiry.Experiences[i].ToDate) != 0 {
				util.AddDateToMonth((enquiry.Experiences[i].ToDate))
			}
		}
	}

	// Give masters abroad and its score arrays created by field.
	if enquiry.MastersAbroad != nil {
		enquiry.MastersAbroad.CreatedBy = credentialID
		for i := 0; i < len(enquiry.MastersAbroad.Scores); i++ {
			enquiry.MastersAbroad.Scores[i].CreatedBy = credentialID
		}
	}

	// Set talent id if exists.
	if err := service.setTalentIDIfExists(uow, enquiry); err != nil {
		uow.RollBack()
		return err
	}

	// Add Enquiry to database.
	err = service.Repository.Add(uow, enquiry)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Failed to add Enquiry", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// AddEnquiryForm adds new enquiry to database from registration form.
func (service *EnquiryService) AddEnquiryForm(enquiry *talenq.Enquiry) error {

	// Give enquiry date as current date.
	currentTime := time.Now()
	enquiry.EnquiryDate = currentTime.Format("2006-01-02")

	// Validate compulsary fields.
	err := enquiry.Validate(true)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Starting transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Extract foreign key IDs and remove the object.
	service.extractID(enquiry)

	// Give tenant id to academics and experiences.
	service.setTenantID(enquiry, enquiry.TenantID)

	// Validate tenant id.
	if err := service.doesTenantExist(enquiry.TenantID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(enquiry); err != nil {
		return err
	}

	// Assign Enquiry Code.
	enquiry.Code, err = util.GenerateUniqueCode(uow.DB, enquiry.FirstName, "`code` = ?", talenq.Enquiry{})
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Failed to generate enquiry code", http.StatusInternalServerError)
	}

	// Set college id by college name.
	if err := service.setCollegeNameAndID(uow, enquiry); err != nil {
		uow.RollBack()
		return err
	}

	// Give 01 to all enquiry's experiences' fromDate and toDate field.
	if enquiry.Experiences != nil && len(enquiry.Experiences) != 0 {
		for i := 0; i < len(enquiry.Experiences); i++ {
			// Set created_by field of experiences.
			if len(enquiry.Experiences[i].FromDate) != 0 {
				util.AddDateToMonth(&(enquiry.Experiences[i].FromDate))
			}
			if enquiry.Experiences[i].ToDate != nil && len(*enquiry.Experiences[i].ToDate) != 0 {
				util.AddDateToMonth((enquiry.Experiences[i].ToDate))
			}
		}
	}

	// Set talent id if exists.
	if err := service.setTalentIDIfExists(uow, enquiry); err != nil {
		uow.RollBack()
		return err
	}

	// Add Enquiry to database.
	err = service.Repository.Add(uow, enquiry)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Failed to add Enquiry", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// AddApplicationForm adds new enquiry and waiting list to database from application form.
func (service *EnquiryService) AddApplicationForm(applicationForm *talenq.ApplicationForm, tenantID uuid.UUID) error {

	// Craete enquiry from applicationForm.
	enquiry := &applicationForm.Enquiry

	// Craete waiting list from applicationForm.
	waitingList := &applicationForm.WaitingList

	// Give enquiry date as current date.
	currentTime := time.Now()
	enquiry.EnquiryDate = currentTime.Format("2006-01-02")

	// Validate compulsary fields.
	err := enquiry.Validate(true)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Starting transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Extract foreign key IDs and remove the object.
	service.extractID(enquiry)

	// Give tenant id to enquiry, academics and experiences.
	enquiry.TenantID = tenantID
	service.setTenantID(enquiry, tenantID)

	// Validate tenant id.
	if err := service.doesTenantExist(enquiry.TenantID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(enquiry); err != nil {
		return err
	}

	// Assign Enquiry Code.
	enquiry.Code, err = util.GenerateUniqueCode(uow.DB, enquiry.FirstName, "`code` = ?", talenq.Enquiry{})
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Failed to generate enquiry code", http.StatusInternalServerError)
	}

	// Set college id by college name.
	if err := service.setCollegeNameAndID(uow, enquiry); err != nil {
		uow.RollBack()
		return err
	}

	// Give 01 to all enquiry's experiences' fromDate and toDate field.
	if enquiry.Experiences != nil && len(enquiry.Experiences) != 0 {
		for i := 0; i < len(enquiry.Experiences); i++ {
			// Set created_by field of experiences.
			if len(enquiry.Experiences[i].FromDate) != 0 {
				util.AddDateToMonth(&(enquiry.Experiences[i].FromDate))
			}
			if enquiry.Experiences[i].ToDate != nil && len(*enquiry.Experiences[i].ToDate) != 0 {
				util.AddDateToMonth((enquiry.Experiences[i].ToDate))
			}
		}
	}

	// Set talent id if exists.
	if err := service.setTalentIDIfExists(uow, enquiry); err != nil {
		uow.RollBack()
		return err
	}

	// Add Enquiry to database.
	err = service.Repository.Add(uow, enquiry)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Failed to add Enquiry", http.StatusInternalServerError)
	}

	// Add waiting list.
	waitingListService := NewWaitingListService(service.DB, service.Repository)
	if err := waitingListService.AddWaitingListForApplicationForm(waitingList, enquiry.ID, tenantID, enquiry.Email, uow); err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// UpdateEnquiry updates one Enquiry in database.
func (service *EnquiryService) UpdateEnquiry(enquiry *talenq.Enquiry) error {
	// Give tenant id to academics and experiences.
	service.setTenantID(enquiry, enquiry.TenantID)

	// Get credential id from UpdatedBy field of enquiry(set in controller).
	credentialID := enquiry.UpdatedBy

	// Extract all foreign key IDs,assign to entityID field and make entity object nil.
	service.extractID(enquiry)

	// Validate tenant id.
	if err := service.doesTenantExist(enquiry.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, enquiry.TenantID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(enquiry); err != nil {
		return err
	}

	// Give 01 to all enquiry's experiences' fromDate and toDate field.
	if enquiry.Experiences != nil && len(enquiry.Experiences) != 0 {
		for i := 0; i < len(enquiry.Experiences); i++ {
			if len(enquiry.Experiences[i].FromDate) != 0 {
				util.AddDateToMonth(&(enquiry.Experiences[i].FromDate))
			}
			if enquiry.Experiences[i].ToDate != nil && len(*enquiry.Experiences[i].ToDate) != 0 {
				util.AddDateToMonth((enquiry.Experiences[i].ToDate))
			}
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Set college id by college name.
	if err := service.setCollegeNameAndID(uow, enquiry); err != nil {
		uow.RollBack()
		return err
	}

	// Create bucket for getting enquiry already present in database.
	tempEnquiry := talenq.Enquiry{}

	// Get enquiry for getting created_by field of enquiry from database.
	if err := service.Repository.GetForTenant(uow, enquiry.TenantID, enquiry.ID, &tempEnquiry,
		repository.PreloadAssociations([]string{"Academics", "Experiences", "MastersAbroad", "MastersAbroad.Scores", "Technologies"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp enquiry to enquiry to be updated.
	enquiry.CreatedBy = tempEnquiry.CreatedBy

	// Give credential id to updated_by of enquiry.
	enquiry.UpdatedBy = credentialID

	// Give code to enquiry.
	enquiry.Code = tempEnquiry.Code

	// Update academics.
	err := service.updateAcademics(uow, enquiry.Academics, enquiry.TenantID, enquiry.UpdatedBy, enquiry.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Make academics nil to avoid any inserts or updates in academics table.
	enquiry.Academics = nil

	// Update experiences.
	err = service.updateExperiences(uow, enquiry.Experiences, enquiry.TenantID, enquiry.UpdatedBy, enquiry.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Make experiences nil to avoid any inserts or updates in experiences table.
	enquiry.Experiences = nil

	// Update matsers abroad.
	err = service.updateMastersAbroad(uow, enquiry.MastersAbroad, tempEnquiry.MastersAbroad, enquiry.TenantID,
		enquiry.UpdatedBy, enquiry.ID, enquiry)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Make matsers abroad nil to avoid any inserts or updates in matsers abroad table.
	enquiry.MastersAbroad = nil

	// Update enquiry associations.
	if err := service.updateEnquiryAssociation(uow, enquiry); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Enquiry could not be updated", http.StatusInternalServerError)
	}

	// Make technologies and experience technologies nil so that it is not inserted again.
	enquiry.Technologies = nil
	enquiry.Courses = nil

	// Update enquiry.
	if err := service.Repository.Save(uow, enquiry); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Enquiry could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetTalentsByWaitingList returns all enquiries by waiting list.
func (service *EnquiryService) GetEnquiriesByWaitingList(enquiries *[]talenq.DTO, tenantID uuid.UUID, limit, offset int,
	totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult, queryParams url.Values) error {

	//********************************************ID filter***************************************************
	// Variables for all possible IDs.
	companyBranchID := uuid.Nil
	requirementID := uuid.Nil
	courseID := uuid.Nil
	batchID := uuid.Nil
	technologyID := uuid.Nil

	// Get query params for all IDs.
	// Company branch ID.
	IDArray, _ := queryParams["companyBranchID"]
	if IDArray != nil && len(IDArray) != 0 {
		var err error
		companyBranchID, err = uuid.FromString(IDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get enquiries", http.StatusInternalServerError)
		}
	}

	// Company requirement ID.
	IDArray, _ = queryParams["requirementID"]
	if IDArray != nil && len(IDArray) != 0 {
		var err error
		requirementID, err = uuid.FromString(IDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get enquiries", http.StatusInternalServerError)
		}
	}

	// Course ID.
	IDArray, _ = queryParams["courseID"]
	if IDArray != nil && len(IDArray) != 0 {
		var err error
		courseID, err = uuid.FromString(IDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get enquiries", http.StatusInternalServerError)
		}
	}

	// Batch ID.
	IDArray, _ = queryParams["batchID"]
	if IDArray != nil && len(IDArray) != 0 {
		var err error
		batchID, err = uuid.FromString(IDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get enquiries", http.StatusInternalServerError)
		}
	}

	// Technology ID.
	IDArray, _ = queryParams["technologyID"]
	if IDArray != nil && len(IDArray) != 0 {
		var err error
		technologyID, err = uuid.FromString(IDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get talents", http.StatusInternalServerError)
		}
	}

	// Create query processors according to conditions.
	var queryProcesors []repository.QueryProcessor

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate company branch ID if present.
	if companyBranchID != uuid.Nil {
		if err := service.doesCompanyBranchExist(companyBranchID, tenantID); err != nil {
			return err
		}

		// Add company branch filter.
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN waiting_list ON talent_enquiries.`id` = waiting_list.`enquiry_id`"),
			repository.Filter("waiting_list.`company_branch_id`=?", companyBranchID))
	}

	// Validate company requirement ID if present.
	if requirementID != uuid.Nil {
		if err := service.doesCompanyRequirementExist(requirementID, tenantID); err != nil {
			return err
		}

		// Add company requirement filter.
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN waiting_list ON talent_enquiries.`id` = waiting_list.`enquiry_id`"),
			repository.Filter("waiting_list.`company_requirement_id`=?", requirementID))
	}

	// Validate course ID if present.
	if courseID != uuid.Nil {
		if err := service.doesCourseExist(courseID, tenantID); err != nil {
			return err
		}

		// Add course filter.
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN waiting_list ON talent_enquiries.`id` = waiting_list.`enquiry_id`"),
			repository.Filter("waiting_list.`course_id`=?", courseID))
	}

	// Validate batch ID if present.
	if batchID != uuid.Nil {
		if err := service.doesBatchExist(tenantID, batchID); err != nil {
			return err
		}

		// Add batch filter.
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN waiting_list ON talent_enquiries.`id` = waiting_list.`enquiry_id`"),
			repository.Filter("waiting_list.`batch_id`=?", batchID))
	}

	// Validate technology ID if present.
	if technologyID != uuid.Nil {
		if err := service.doesTechnologyExist(technologyID, tenantID); err != nil {
			return err
		}

		// Add technology filter.
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN waiting_list ON talent_enquiries.`id` = waiting_list.`enquiry_id`"),
			repository.Join("LEFT JOIN company_requirements on waiting_list.`company_requirement_id` = company_requirements.`id`"),
			repository.Join("LEFT JOIN company_requirements_technologies on company_requirements.`id` = company_requirements_technologies.`requirement_id`"),
			repository.Join("LEFT JOIN courses on waiting_list.`course_id` = courses.`id`"),
			repository.Join("LEFT JOIN courses_technologies on courses.`id` = courses_technologies.`course_id`"),
			repository.Filter("company_requirements_technologies.`technology_id`=? OR courses_technologies.`technology_id`=?", technologyID, technologyID))
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	queryProcesors = append(queryProcesors,
		repository.Filter("talent_enquiries.`id` IS NOT NULL"),
		repository.Filter("waiting_list.`deleted_at` IS NULL"),
		repository.Filter("waiting_list.`tenant_id`=?", tenantID),
		repository.Filter("waiting_list.`talent_id` IS NULL"),
		repository.Filter("talent_enquiries.`tenant_id`=?", tenantID),
		repository.GroupBy("talent_enquiries.`id`"),
		repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount))

	err := service.Repository.GetAllInOrder(uow, enquiries, "`enquiry_date` desc, `first_name`, `last_name`", queryProcesors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}

	// Range enquiries to get expected ctc.
	for index := range *enquiries {
		//**************************************************ExpectedCTC(call records)***********************************************************

		// Create bucket for expected ctc.
		exepectedCTC := &tal.ExpectedCTCLatest{}

		// Get expected ctc from database.
		if err := service.Repository.Scan(uow, exepectedCTC,
			repository.Filter("talent_enquiries.`deleted_at` IS NULL AND talent_enquiry_call_records.`deleted_at` IS NULL AND talent_enquiries.`tenant_id`=? AND talent_enquiry_call_records.`tenant_id`=? AND talent_enquiries.`id`=?",
				tenantID, tenantID, (*enquiries)[index].ID),
			repository.Filter("talent_enquiry_call_records.`expected_ctc` IS NOT NULL"),
			repository.Table("talent_enquiries"),
			repository.Join("JOIN talent_enquiry_call_records on talent_enquiry_call_records.`enquiry_id` = talent_enquiries.`id`"),
			repository.Select("expected_ctc"),
			repository.GroupBy("talent_enquiries.`id`"),
			repository.OrderBy("talent_enquiry_call_records.`date_time`")); err != nil {
			uow.RollBack()
			if err != gorm.ErrRecordNotFound {
				log.NewLogger().Error(err.Error())
				return errors.NewValidationError("Internal server error")
			}
		}

		// Give expected CTC to enquiry
		if exepectedCTC != nil && exepectedCTC.ExpectedCTC > 0 {
			(*enquiries)[index].ExpectedCTC = &exepectedCTC.ExpectedCTC
		}
	}

	// Sort the child tables.
	service.sortEnquiryChildTables(enquiries)

	uow.Commit()
	return nil
}

// ConvertToTalent converts enquiry to talent.
func (service *EnquiryService) ConvertToTalent(enquiry *talenq.Enquiry) error {
	credentialID := enquiry.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(enquiry.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, enquiry.TenantID); err != nil {
		return err
	}

	// Start new transaction.
	db := service.DB
	uow := repository.NewUnitOfWork(db, false)

	// Get enquiry from database.
	if err := service.Repository.GetForTenant(uow, enquiry.TenantID, enquiry.ID, enquiry,
		repository.PreloadAssociations([]string{"Technologies", "Academics", "Experiences", "MastersAbroad", "Experiences.Technologies"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give credential id to update_by field of enquiry.
	enquiry.UpdatedBy = credentialID

	// Generate code.
	code, err := util.GenerateUniqueCode(uow.DB, enquiry.FirstName, "`code` = ?", &tal.Talent{})
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	talent := &talent.Talent{
		Code: code,
	}
	assignEnquiryFieldsToTalent(enquiry, talent)

	// Add entry in talent table.
	talentService := NewTalentService(db, service.Repository)
	err = talentService.AddTalentAndCredential(uow, talent)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Update masters abroad's talent id.
	if enquiry.MastersAbroad != nil {
		err = service.Repository.UpdateWithMap(uow, &general.MastersAbroad{},
			map[string]interface{}{
				"TalentID": talent.ID,
			}, repository.Filter("`id` = ?", enquiry.MastersAbroad.ID))
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Enquiry could not be updated", http.StatusInternalServerError)
		}
	}

	// In waiting list table, if any entry exists related to the enquiry, then give all the waiting lists talent id.
	// Update waiting list record.
	err = service.Repository.UpdateWithMap(uow, tal.WaitingList{},
		map[string]interface{}{
			"TalentID":  talent.ID,
			"UpdatedBy": enquiry.UpdatedBy,
		}, repository.Filter("`email`=?", enquiry.Email))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Waiting list could not be updated", http.StatusInternalServerError)
	}

	// Update enquiry record.
	err = service.Repository.UpdateWithMap(uow, talenq.Enquiry{},
		map[string]interface{}{
			"TalentID":  talent.ID,
			"UpdatedBy": enquiry.UpdatedBy,
		}, repository.Filter("`email`=?", enquiry.Email))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Enquiry could not be updated", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// assignEnquiryFieldsToTalent assigns all fields of enquiry to all fields of talent.
func assignEnquiryFieldsToTalent(enquiry *talenq.Enquiry, tal *talent.Talent) {
	credentialID := enquiry.UpdatedBy
	tenantID := enquiry.TenantID
	var trueFlag bool = true

	tal.IgnoreHook = trueFlag
	tal.ID = util.GenerateUUID()
	enquiry.TalentID = &tal.ID
	tal.IsActive = true
	tal.CreatedBy = credentialID
	tal.TenantID = tenantID
	tal.FirstName = enquiry.FirstName
	tal.LastName = enquiry.LastName
	tal.Email = enquiry.Email
	tal.Contact = enquiry.Contact
	tal.IsExperience = enquiry.IsExperience
	tal.SalesPersonID = enquiry.SalesPersonID
	tal.Resume = enquiry.Resume
	tal.AlternateContact = enquiry.AlternateContact
	tal.AlternateEmail = enquiry.AlternateEmail
	tal.FacebookURL = enquiry.FacebookUrl
	tal.InstagramURL = enquiry.InstagramUrl
	tal.GithubURL = enquiry.GithubUrl
	tal.LinkedInURL = enquiry.LinkedInUrl
	tal.AcademicYear = enquiry.AcademicYear
	tal.IsSwabhavTalent = false

	// Masters abroad fields.
	tal.IsMastersAbroad = enquiry.IsMastersAbroad

	// Address fields.
	tal.Address = enquiry.Address
	tal.City = enquiry.City
	tal.PINCode = enquiry.PINCode
	tal.StateID = enquiry.StateID
	tal.CountryID = enquiry.CountryID
	tal.SourceID = enquiry.SourceID
	tal.Technologies = enquiry.Technologies

	// Academics fields.
	for _, academic := range enquiry.Academics {
		talentAcademic := &talent.Academic{
			DegreeID:         academic.DegreeID,
			College:          academic.College,
			CollegeID:        academic.CollegeID,
			Percentage:       academic.Percentage,
			Passout:          academic.Passout,
			SpecializationID: academic.SpecializationID,
			TalentID:         tal.ID,
		}
		talentAcademic.CreatedBy = credentialID
		talentAcademic.TenantID = tenantID
		tal.Academics = append(tal.Academics, talentAcademic)
	}

	//Experience fields.
	for _, experience := range enquiry.Experiences {
		if len(experience.FromDate) > 10 {
			experience.FromDate = experience.FromDate[:10]
		}
		if experience.ToDate != nil {
			if len(*experience.ToDate) > 10 {
				*experience.ToDate = (*experience.ToDate)[:10]
			}
		}
		talentExperience := &talent.Experience{
			Company:       experience.Company,
			Technologies:  experience.Technologies,
			FromDate:      experience.FromDate,
			ToDate:        experience.ToDate,
			Package:       experience.Package,
			DesignationID: experience.DesignationID,
			TalentID:      tal.ID,
		}
		talentExperience.CreatedBy = credentialID
		talentExperience.TenantID = tenantID
		tal.Experiences = append(tal.Experiences, talentExperience)
	}
}

// DeleteEnquiry deletes Enquiry from database.
func (service *EnquiryService) DeleteEnquiry(enquiry *talenq.Enquiry) error {
	// Get credential id from DeletedBy field of enquiry(set in controller).
	credentialID := enquiry.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(enquiry.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, enquiry.TenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get enquiry for updating deleted_by field of enquiry.
	if err := service.Repository.GetForTenant(uow, enquiry.TenantID, enquiry.ID, enquiry,
		repository.PreloadAssociations([]string{"Technologies", "MastersAbroad"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Delete enquiry association from database.
	if err := service.deleteEnquiryAssociation(uow, enquiry, credentialID); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Failed to delete enquiry", http.StatusInternalServerError)
	}

	// Make technologies and masters abroad field nil to avoid any updates or inserts.
	enquiry.Technologies = nil
	enquiry.Courses = nil
	enquiry.MastersAbroad = nil

	// Update talent for updating deleted_by and deleted_at fields of talent.
	if err := service.Repository.UpdateWithMap(uow, enquiry, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", enquiry.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetEnquiry gets one enquiry from database.
func (service *EnquiryService) GetEnquiry(enquiry *talenq.DTO, tenantID uuid.UUID) error {
	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get enquiry.
	if err := service.Repository.GetForTenant(uow, tenantID, enquiry.ID, enquiry,
		repository.PreloadAssociations(service.associations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Remove timestamp from EnquiryDate field of enquiry.
	enquiry.EnquiryDate = enquiry.EnquiryDate[:10]

	uow.Commit()
	return nil
}

// GetEnquiries gets all Enquiries from database.
func (service *EnquiryService) GetEnquiries(enquiries *[]talenq.DTO, tenantID uuid.UUID, limit int, offset int, totalCount *int,
	queryParams url.Values) error {
	// Variables for role name and login id.
	roleName := ""
	loginID := uuid.Nil

	// Get query params for role name and login id.
	roleNameArray, _ := queryParams["roleName"]
	if roleNameArray != nil && len(roleNameArray) != 0 {
		roleName = roleNameArray[0]
	}
	loginIDArray, _ := queryParams["loginID"]
	if loginIDArray != nil && len(loginIDArray) != 0 {
		var err error
		loginID, err = uuid.FromString(loginIDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get enquiries", http.StatusInternalServerError)
		}
	}

	// Create query processors according to conditions.
	var queryProcesors []repository.QueryProcessor

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// If login is salesperson.
	if roleName == "SalesPerson" {
		// Validate salespersonid (login id).
		if err := service.doesSalespersonExist(loginID, tenantID); err != nil {
			return err
		}
		// Add salesperson filter.
		queryProcesors = append(queryProcesors, repository.Filter("talent_enquiries.`sales_person_id`=?", loginID))
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Add paginate and preload.
	queryProcesors = append(queryProcesors, repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount))

	// Get enquiries from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, enquiries, "`enquiry_date` desc, `first_name`, `last_name`",
		queryProcesors...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Range enquiries to get expected ctc.
	for index := range *enquiries {
		//**************************************************ExpectedCTC(call records)***********************************************************

		// Create bucket for expected ctc.
		exepectedCTC := &tal.ExpectedCTCLatest{}

		// Get expected ctc from database.
		if err := service.Repository.Scan(uow, exepectedCTC,
			repository.Filter("talent_enquiries.`deleted_at` IS NULL AND talent_enquiry_call_records.`deleted_at` IS NULL AND talent_enquiries.`tenant_id`=? AND talent_enquiry_call_records.`tenant_id`=? AND talent_enquiries.`id`=?",
				tenantID, tenantID, (*enquiries)[index].ID),
			repository.Filter("talent_enquiry_call_records.`expected_ctc` IS NOT NULL"),
			repository.Table("talent_enquiries"),
			repository.Join("JOIN talent_enquiry_call_records on talent_enquiry_call_records.`enquiry_id` = talent_enquiries.`id`"),
			repository.Select("expected_ctc"),
			repository.GroupBy("talent_enquiries.`id`"),
			repository.OrderBy("talent_enquiry_call_records.`date_time`")); err != nil {
			uow.RollBack()
			if err != gorm.ErrRecordNotFound {
				log.NewLogger().Error(err.Error())
				return errors.NewValidationError("Internal server error")
			}
		}

		// Give expected CTC to enquiry
		if exepectedCTC != nil && exepectedCTC.ExpectedCTC > 0 {
			(*enquiries)[index].ExpectedCTC = &exepectedCTC.ExpectedCTC
		}
	}

	// Sort the child tables.
	service.sortEnquiryChildTables(enquiries)

	uow.Commit()
	return nil
}

// updateEnquiryAssociation updates enquiry's associations.
func (service *EnquiryService) updateEnquiryAssociation(uow *repository.UnitOfWork, enquiry *talenq.Enquiry) error {
	// Replace technologies of enquiry
	if err := service.Repository.ReplaceAssociations(uow, enquiry, "Technologies",
		enquiry.Technologies); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace courses of enquiry
	if err := service.Repository.ReplaceAssociations(uow, enquiry, "Courses",
		enquiry.Courses); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// deleteEnquiryAssociation deletes enquiry's associations.
func (service *EnquiryService) deleteEnquiryAssociation(uow *repository.UnitOfWork, enquiry *talenq.Enquiry, credentialID uuid.UUID) error {
	// ******************************************Delete experiences*****************************************************
	if err := service.Repository.UpdateWithMap(uow, &talenq.Experience{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`enquiry_id`=? AND `tenant_id`=?", enquiry.ID, enquiry.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Enquiry could not be deleted", http.StatusInternalServerError)
	}

	//*******************************************Delete academics*******************************************************
	if err := service.Repository.UpdateWithMap(uow, &talenq.Academic{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`enquiry_id`=? AND `tenant_id`=?", enquiry.ID, enquiry.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Enquiry could not be deleted", http.StatusInternalServerError)
	}

	//******************************************Delete masters abroad*************************************************************
	if enquiry.IsMastersAbroad && enquiry.MastersAbroad.TalentID == nil {
		// Update score(s) for updating deleted_by field of score(s).
		err := service.Repository.UpdateWithMap(uow, &general.Score{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		},
			repository.Filter("`masters_abroad_id`=? AND `tenant_id`=?", enquiry.MastersAbroad.ID, enquiry.TenantID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}

		// Make countries and universities field nil to avoid any updates or inserts.
		enquiry.MastersAbroad.Countries = nil
		enquiry.MastersAbroad.Universities = nil

		// Update masters abroad for updating deleted_by and deleted_at fields of masters abroad.
		if err := service.Repository.UpdateWithMap(uow, &general.MastersAbroad{},
			map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			},
			repository.Filter("`enquiry_id`=? AND `tenant_id`=?", enquiry.ID, enquiry.TenantID)); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Enquiry could not be deleted", http.StatusInternalServerError)
		}
	}

	//*******************************************Deleting call records*******************************************************
	if err := service.Repository.UpdateWithMap(uow, &talenq.CallRecord{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`enquiry_id`=? AND `tenant_id`=?", enquiry.ID, enquiry.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Enquiry could not be deleted", http.StatusInternalServerError)
	}
	return nil
}

// GetAllSearchEnquiries gets all searched enquiries from database.
func (service *EnquiryService) GetAllSearchEnquiries(enquiries *[]talenq.DTO, enquirySearch *talenq.Search, tenantID uuid.UUID,
	limit int, offset int, totalCount *int, queryParams url.Values) error {
	// Variables for role name and login id.
	roleName := ""
	loginID := uuid.Nil

	// Get query params for role name and login id.
	roleNameArray, _ := queryParams["roleName"]
	if roleNameArray != nil && len(roleNameArray) != 0 {
		roleName = roleNameArray[0]
	}
	loginIDArray, _ := queryParams["loginID"]
	if loginIDArray != nil && len(loginIDArray) != 0 {
		var err error
		loginID, err = uuid.FromString(loginIDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get enquiries", http.StatusInternalServerError)
		}
	}

	// Create query processors according to conditions.
	var queryProcessors []repository.QueryProcessor

	// If search all talents is false then add filter by salesperson.
	if enquirySearch.SearchAllEnquiries != nil && !(*enquirySearch.SearchAllEnquiries) {
		// If login is salesperson.
		if roleName == "SalesPerson" {
			// Validate salespersonid (login id).
			if err := service.doesSalespersonExist(loginID, tenantID); err != nil {
				return err
			}
			// Add salesperson filter.
			queryProcessors = append(queryProcessors, repository.Filter("talent_enquiries.`sales_person_id`=?", loginID))
		}
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	//A dd preload, pagination and search queries.
	queryProcessors = append(queryProcessors, service.enquiryAddSearchQueries(tenantID, enquirySearch, roleName)...)
	queryProcessors = append(queryProcessors, repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount),
		repository.Filter("talent_enquiries.`tenant_id`=?", tenantID))

	if err := service.Repository.GetAllInOrder(uow, enquiries, "`enquiry_date` desc, `first_name`, `last_name`", queryProcessors...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Ranhe enquiries to get expected ctc.
	for index := range *enquiries {
		//**************************************************ExpectedCTC(call records)***********************************************************

		// Create bucket for expected ctc.
		exepectedCTC := &tal.ExpectedCTCLatest{}

		// Get expected ctc from database.
		if err := service.Repository.Scan(uow, exepectedCTC,
			repository.Filter("talent_enquiries.`deleted_at` IS NULL AND talent_enquiry_call_records.`deleted_at` IS NULL AND talent_enquiries.`tenant_id`=? AND talent_enquiry_call_records.`tenant_id`=? AND talent_enquiries.`id`=?",
				tenantID, tenantID, (*enquiries)[index].ID),
			repository.Filter("talent_enquiry_call_records.`expected_ctc` IS NOT NULL"),
			repository.Table("talent_enquiries"),
			repository.Join("JOIN talent_enquiry_call_records on talent_enquiry_call_records.`enquiry_id` = talent_enquiries.`id`"),
			repository.Select("expected_ctc"),
			repository.GroupBy("talent_enquiries.`id`"),
			repository.OrderBy("talent_enquiry_call_records.`date_time`")); err != nil {
			if err != gorm.ErrRecordNotFound {
				log.NewLogger().Error(err.Error())
				uow.RollBack()
				return errors.NewValidationError("Record not found")
			}
		}

		// Give expected CTC to enquiry
		if exepectedCTC != nil && exepectedCTC.ExpectedCTC > 0 {
			(*enquiries)[index].ExpectedCTC = &exepectedCTC.ExpectedCTC
		}
	}

	// Sort the child tables.
	service.sortEnquiryChildTables(enquiries)

	uow.Commit()
	return nil
}

// UpdateEnquiriesSalesperson updates multiple enquiries' salesperson id.
func (service *EnquiryService) UpdateEnquiriesSalesperson(enquiries *[]talenq.EnquiryUpdate, salepersonID uuid.UUID,
	tenantID uuid.UUID, credentialID uuid.UUID) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, tenantID); err != nil {
		return err
	}

	// Validate salespersonid (login id).
	if err := service.doesSalespersonExist(salepersonID, tenantID); err != nil {
		return err
	}

	// Collect all enquiry ids in variable.
	var enquiryIDs []uuid.UUID
	for _, enquiry := range *enquiries {
		enquiryIDs = append(enquiryIDs, enquiry.EnquiryID)
	}

	// Validate all enquiries.
	for _, enquiryID := range enquiryIDs {
		if err := service.doesEnquiryExist(enquiryID, tenantID); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update sales_person_id field of all enquiries.
	if err := service.Repository.UpdateWithMap(uow, &talenq.Enquiry{}, map[string]interface{}{
		"SalesPersonID": salepersonID,
		"UpdatedBy":     credentialID,
	},
		repository.Filter("`id` IN (?)", enquiryIDs)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Sales person could not be allocated to enquiries", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// AddEnquiries adds multiple enquiries to database.
func (service *EnquiryService) AddEnquiries(enquiries *[]talenq.EnquiryExcel, enquiryAddedCount *int,
	tenantID uuid.UUID, errorList *[]error, credentialID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get IDs for all values and add individual enquiry to database.
	if enquiries != nil && len(*enquiries) != 0 {

		// ExcelLoop is a label for outer loop
	ExcelLoop:

		// Loop the enquiries in enquiries from excel.
		for _, enquiryExcel := range *enquiries {

			// Create bucket for enquiry.
			enquiry := talenq.Enquiry{}

			// Get country ID.
			if enquiryExcel.CountryName != nil {

				// Create bucket for country ID.
				countryID := tal.IDModel{}

				// Get country ID from database.
				err := service.Repository.GetRecordForTenant(uow, tenantID, &countryID,
					repository.Select("`id`"),
					repository.Filter("`name`=?", enquiryExcel.CountryName),
					repository.Filter("`tenant_id`=?", tenantID),
					repository.Table("countries"))

				if err != nil {
					// More efficiency needed #Niranjan
					log.NewLogger().Error(err.Error())
					errorString := ""
					if err == gorm.ErrRecordNotFound {
						errorString = enquiryExcel.Email + " : Invalid country name"
						err = errors.NewHTTPError(errorString, http.StatusBadRequest)
					} else {
						errorString = enquiryExcel.Email + " : Internal Server Error"
						err = errors.NewHTTPError(errorString, http.StatusInternalServerError)
					}
					*errorList = append(*errorList, err)
					continue
				}

				// Give country ID to enquiry.
				enquiry.CountryID = &countryID.ID
			}

			// Get state ID.
			if enquiryExcel.StateName != nil {

				// Create bucket for state ID.
				stateID := tal.IDModel{}

				// Get state ID from database.
				err := service.Repository.GetRecordForTenant(uow, tenantID, &stateID,
					repository.Select("`id`"),
					repository.Filter("`name`=?", enquiryExcel.StateName),
					repository.Filter("`tenant_id`=?", tenantID),
					repository.Table("states"))

				if err != nil {
					log.NewLogger().Error(err.Error())
					errorString := ""
					if err == gorm.ErrRecordNotFound {
						errorString = enquiryExcel.Email + " : Invalid state name"
						err = errors.NewHTTPError(errorString, http.StatusBadRequest)
					} else {
						errorString = enquiryExcel.Email + " : Internal Server Error"
						err = errors.NewHTTPError(errorString, http.StatusInternalServerError)
					}
					*errorList = append(*errorList, err)
					continue
				}

				// Give state ID to enquiry.
				enquiry.StateID = &stateID.ID
			}

			// Check if academics exist or not.
			if enquiryExcel.Academics != nil && len(enquiryExcel.Academics) > 0 {

				// Create bucket for academics.
				academics := []*talenq.Academic{}

				// Loop the academics of enquiry excel.
				for _, academicExcel := range enquiryExcel.Academics {

					// Create bucket for academic.
					academic := talenq.Academic{}

					// Create bucket for degree ID.
					degreeID := talenq.IDModel{}

					// Get degree ID from database.
					err := service.Repository.GetRecordForTenant(uow, tenantID, &degreeID,
						repository.Select("`id`"),
						repository.Filter("`name`=?", academicExcel.DegreeName),
						repository.Filter("`tenant_id`=?", tenantID),
						repository.Table("degrees"))

					if err != nil {
						log.NewLogger().Error(err.Error())
						errorString := ""
						if err == gorm.ErrRecordNotFound {
							errorString = enquiryExcel.Email + " : Invalid degree name"
							err = errors.NewHTTPError(errorString, http.StatusBadRequest)
						} else {
							errorString = enquiryExcel.Email + " : Internal Server Error"
							err = errors.NewHTTPError(errorString, http.StatusInternalServerError)
						}
						*errorList = append(*errorList, err)
						continue ExcelLoop
					}

					// Give degree ID to academic.
					academic.DegreeID = degreeID.ID

					// Create bucket for specialization ID.
					specializationID := tal.IDModel{}

					// Get specialization ID from database.
					err = service.Repository.GetRecordForTenant(uow, tenantID, &specializationID,
						repository.Select("`id`"),
						repository.Filter("`branch_name`=?", academicExcel.SpecializationName),
						repository.Filter("`degree_id`=?", degreeID.ID),
						repository.Filter("`tenant_id`=?", tenantID),
						repository.Table("specializations"))

					if err != nil {
						log.NewLogger().Error(err.Error())
						errorString := ""
						if err == gorm.ErrRecordNotFound {
							errorString = enquiryExcel.Email + " : Invalid specialization name"
							err = errors.NewHTTPError(errorString, http.StatusBadRequest)
						} else {
							errorString = enquiryExcel.Email + " : Internal Server Error"
							err = errors.NewHTTPError(errorString, http.StatusInternalServerError)
						}
						*errorList = append(*errorList, err)
						continue ExcelLoop
					}

					// Give specialization ID to academic.
					academic.SpecializationID = specializationID.ID

					// Give college name to academic.
					academic.College = academicExcel.CollegeName

					// Give percentage to academic.
					academic.Percentage = academicExcel.Percentage

					// Give year of passout to academic.
					academic.Passout = academicExcel.YearOfPassout

					// Push academic into academics.
					academics = append(academics, &academic)
				}

				// Give academics to enquiry.
				enquiry.Academics = academics
			}

			// Give first name to enquiry.
			enquiry.FirstName = enquiryExcel.FirstName

			// Give last name to enquiry.
			enquiry.LastName = enquiryExcel.LastName

			// Give email to enquiry.
			enquiry.Email = enquiryExcel.Email

			// Give contact to enquiry.
			enquiry.Contact = enquiryExcel.Contact

			// Give academic year to enquiry.
			enquiry.AcademicYear = enquiryExcel.AcademicYear

			// Give address to enquiry.
			enquiry.Address.Address = enquiryExcel.Address

			// Give city to enquiry.
			enquiry.City = enquiryExcel.City

			// Give Pin code to enquiry.
			enquiry.PINCode = enquiryExcel.PINCode

			// Give tenant ID to enquiry.
			enquiry.TenantID = tenantID

			// Give created by to enquiry.
			enquiry.CreatedBy = credentialID

			// Give enquiry date as current date.
			currentTime := time.Now()
			enquiry.EnquiryDate = currentTime.Format("2006-01-02")

			// Give enquiry source as default "Training and Placement" #Will change in future******************
			enquiry.EnquiryType = "Training And Placement"

			// Add enquiry individually.
			if err := service.AddEnquiry(&enquiry); err != nil {

				// If error then push it in error list.
				errorString := enquiryExcel.Email + " : " + err.Error()
				er := errors.NewHTTPError(errorString, http.StatusBadRequest)
				*errorList = append(*errorList, er)
				continue
			}

			// Increment count of enquiries added successfully.
			*enquiryAddedCount++
		}
	}
	return nil
}

// AddEnquiryFromExcel adds one enquiry from excel to database.
func (service *EnquiryService) AddEnquiryFromExcel(enquiryExcel *talenq.EnquiryExcel, tenantID, credentialID uuid.UUID) error {

	// Start new transaction
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get IDs for all values and add enquiry to database.
	// Create bucket for enquiry.
	enquiry := talenq.Enquiry{}

	// Get country ID.
	if enquiryExcel.CountryName != nil {

		// Create bucket for country ID.
		countryID := tal.IDModel{}

		// Get country ID from database.
		err := service.Repository.GetRecordForTenant(uow, tenantID, &countryID,
			repository.Select("`id`"),
			repository.Filter("`name`=?", enquiryExcel.CountryName),
			repository.Filter("`tenant_id`=?", tenantID),
			repository.Table("countries"))

		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.Commit()
			if err == gorm.ErrRecordNotFound {
				return errors.NewValidationError("Invalid country name")
			} else {
				return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
			}
		}

		// Give country ID to enquiry.
		enquiry.CountryID = &countryID.ID
	}

	// Get state ID.
	if enquiryExcel.StateName != nil {

		// Create bucket for state ID.
		stateID := tal.IDModel{}

		// Get state ID from database.
		err := service.Repository.GetRecordForTenant(uow, tenantID, &stateID,
			repository.Select("`id`"),
			repository.Filter("`name`=?", enquiryExcel.StateName),
			repository.Filter("`tenant_id`=?", tenantID),
			repository.Table("states"))

		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.Commit()
			if err == gorm.ErrRecordNotFound {
				return errors.NewValidationError("Invalid state name")
			} else {
				return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
			}
		}

		// Give state ID to enquiry.
		enquiry.StateID = &stateID.ID
	}

	// Check if academics exist or not.
	if enquiryExcel.Academics != nil && len(enquiryExcel.Academics) > 0 {

		// Create bucket for academics.
		academics := []*talenq.Academic{}

		// Loop the academics of enquiry excel.
		for _, academicExcel := range enquiryExcel.Academics {

			// Create bucket for academic.
			academic := talenq.Academic{}

			// Create bucket for degree ID.
			degreeID := tal.IDModel{}

			// Get degree ID from database.
			err := service.Repository.GetRecordForTenant(uow, tenantID, &degreeID,
				repository.Select("`id`"),
				repository.Filter("`name`=?", academicExcel.DegreeName),
				repository.Filter("`tenant_id`=?", tenantID),
				repository.Table("degrees"))

			if err != nil {
				log.NewLogger().Error(err.Error())
				uow.Commit()
				if err == gorm.ErrRecordNotFound {
					return errors.NewValidationError("Invalid degree name")
				} else {
					return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
				}
			}

			// Give degree ID to academic.
			academic.DegreeID = degreeID.ID

			// Create bucket for specialization ID.
			specializationID := tal.IDModel{}

			// Get specialization ID from database.
			err = service.Repository.GetRecordForTenant(uow, tenantID, &specializationID,
				repository.Select("`id`"),
				repository.Filter("`branch_name`=?", academicExcel.SpecializationName),
				repository.Filter("`degree_id`=?", degreeID.ID),
				repository.Filter("`tenant_id`=?", tenantID),
				repository.Table("specializations"))

			if err != nil {
				log.NewLogger().Error(err.Error())
				uow.Commit()
				if err == gorm.ErrRecordNotFound {
					return errors.NewValidationError("Invalid specialization name")
				} else {
					return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
				}
			}

			// Give specialization ID to academic.
			academic.SpecializationID = specializationID.ID

			// Give college name to academic.
			academic.College = academicExcel.CollegeName

			// Give percentage to academic.
			academic.Percentage = academicExcel.Percentage

			// Give year of passout to academic.
			academic.Passout = academicExcel.YearOfPassout

			// Push academic into academics.
			academics = append(academics, &academic)
		}

		// Give academics to enquiry.
		enquiry.Academics = academics
	}

	// Give first name to enquiry.
	enquiry.FirstName = enquiryExcel.FirstName

	// Give last name to enquiry.
	enquiry.LastName = enquiryExcel.LastName

	// Give email to enquiry.
	enquiry.Email = enquiryExcel.Email

	// Give contact to enquiry.
	enquiry.Contact = enquiryExcel.Contact

	// Give academic year to enquiry.
	enquiry.AcademicYear = enquiryExcel.AcademicYear

	// Give address to enquiry.
	enquiry.Address.Address = enquiryExcel.Address

	// Give city to enquiry.
	enquiry.City = enquiryExcel.City

	// Give Pin code to enquiry.
	enquiry.PINCode = enquiryExcel.PINCode

	// Give tenant ID to enquiry.
	enquiry.TenantID = tenantID

	// Give created by to enquiry.
	enquiry.CreatedBy = credentialID

	// Give enquiry date as current date.
	currentTime := time.Now()
	enquiry.EnquiryDate = currentTime.Format("2006-01-02")

	// Give enquiry source as default "Training and Placement" #Will change in future******************
	enquiry.EnquiryType = "Training And Placement"

	uow.Commit()

	// Add enquiry individually.
	if err := service.AddEnquiry(&enquiry); err != nil {
		return err
	}

	return nil
}

// updateAcademics will update the academics for specified enquiry.
func (service *EnquiryService) updateAcademics(uow *repository.UnitOfWork, academics []*talenq.Academic,
	tenantID, credentialID, enquiryID uuid.UUID) error {

	// If previous academics is present and current academics is not present.
	if academics == nil {
		err := service.Repository.UpdateWithMap(uow, talenq.Academic{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		}, repository.Filter("`enquiry_id`=?", enquiryID))
		if err != nil {
			return err
		}
	}

	// Create temp academics to get presvious academics of enquiry.
	tempAcademics := []talenq.Academic{}

	err := service.Repository.GetAllForTenant(uow, tenantID, &tempAcademics,
		repository.Filter("`enquiry_id`=?", enquiryID))
	if err != nil {
		return err
	}

	// Make map to count occurences of academic id in previous and current academics.
	academicIDMap := make(map[uuid.UUID]uint)

	// Count the number of occurence of previous academic ID.
	for _, tempAcademic := range tempAcademics {
		academicIDMap[tempAcademic.ID]++
	}

	// Count the number of occurrence of current academic ID to know total count of occurrenec of each ID.
	for _, academic := range academics {

		// If ID is valid then push its occurence in ID map.
		if util.IsUUIDValid(academic.ID) {
			academicIDMap[academic.ID]++
		} else {
			// If ID is nil create new academic entry i table.
			academic.CreatedBy = credentialID
			academic.TenantID = tenantID
			academic.EnquiryID = enquiryID
			err = service.Repository.Add(uow, &academic)
			if err != nil {
				return err
			}
		}

		// If number of occurrence is more than one (present in previous and current academics) then update academic.
		if academicIDMap[academic.ID] > 1 {
			academic.UpdatedBy = credentialID
			err = service.Repository.Update(uow, &academic)
			if err != nil {
				return err
			}
			// Make the number of occurrences 0 after updating academic.
			academicIDMap[academic.ID] = 0
		}
	}

	// If the number of occurrence is one, the academic was presnt in previous academics only, delete it.
	for _, tempAcademic := range tempAcademics {
		if academicIDMap[tempAcademic.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, talenq.Academic{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempAcademic.ID))
			if err != nil {
				return err
			}
			// Make the number of occurrences 0 after deleting academic.
			academicIDMap[tempAcademic.ID] = 0
		}
	}
	return nil
}

// updateExperiences will update the experiences for specified enquiry.
func (service *EnquiryService) updateExperiences(uow *repository.UnitOfWork, experiences []*talenq.Experience,
	tenantID, credentialID, enquiryID uuid.UUID) error {

	// If previous experiences is present and current experiences is not present.
	if experiences == nil {
		err := service.Repository.UpdateWithMap(uow, talenq.Experience{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		}, repository.Filter("`enquiry_id`=?", enquiryID))
		if err != nil {
			return err
		}
	}

	// Create temp experiences to get presvious experiences of enquiry.
	tempExperiences := []talenq.Experience{}

	err := service.Repository.GetAllForTenant(uow, tenantID, &tempExperiences,
		repository.Filter("`enquiry_id`=?", enquiryID))
	if err != nil {
		return err
	}

	// Make map to count occurences of experience id in previous and current experiences.
	experienceIDMap := make(map[uuid.UUID]uint)

	// Count the number of occurence of previous experience ID.
	for _, tempExperience := range tempExperiences {
		experienceIDMap[tempExperience.ID]++
	}

	// Count the number of occurrence of current experience ID to know total count of occurrenec of each ID.
	for _, experience := range experiences {

		// If ID is valid then push its occurence in ID map.
		if util.IsUUIDValid(experience.ID) {
			experienceIDMap[experience.ID]++
		} else {
			// If ID is nil create new experience entry i table.
			experience.CreatedBy = credentialID
			experience.TenantID = tenantID
			experience.EnquiryID = enquiryID
			err = service.Repository.Add(uow, &experience)
			if err != nil {
				return err
			}
		}

		// If number of occurrence is more than one (present in previous and current experiences) then update experience.
		if experienceIDMap[experience.ID] > 1 {
			// Give created_by field of previous experiences to current experiences.
			for _, tempExperience := range tempExperiences {
				if tempExperience.ID == experience.ID {
					experience.CreatedBy = tempExperience.CreatedBy
				}
			}

			// Replace technologies of expereince.
			if err := service.Repository.ReplaceAssociations(uow, experience, "Technologies",
				experience.Technologies); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}

			// Make technologies nil to avoid insertion during update enquiry.
			experience.Technologies = nil

			// Give updated_by field to current experiences.
			experience.UpdatedBy = credentialID
			experience.EnquiryID = enquiryID
			err = service.Repository.Save(uow, &experience)
			if err != nil {
				return err
			}

			// Make the number of occurrences 0 after updating experience.
			experienceIDMap[experience.ID] = 0
		}
	}

	// If the number of occurrence is one, the experience was presnt in previous experiences only, delete it.
	for _, tempExperience := range tempExperiences {
		if experienceIDMap[tempExperience.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, talenq.Experience{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempExperience.ID))
			if err != nil {
				return err
			}
			// Make the number of occurrences 0 after deleting experience.
			experienceIDMap[tempExperience.ID] = 0
		}
	}
	return nil
}

// updateMastersAbroad will update masters abroad for specified enquiry.
func (service *EnquiryService) updateMastersAbroad(uow *repository.UnitOfWork, mastersAbroad, tempMastersAbroad *general.MastersAbroad,
	tenantID, credentialID, enquiryID uuid.UUID, enquiry *talenq.Enquiry) error {
	// If previous masters abroad does not exist and current masters abroad exists.
	if tempMastersAbroad == nil && mastersAbroad != nil {
		// Give masters abroad and its score arrays created by field.
		mastersAbroad.CreatedBy = credentialID
		for i := 0; i < len(mastersAbroad.Scores); i++ {
			mastersAbroad.Scores[i].CreatedBy = credentialID
		}
		// Add masters abroad.
		err := service.Repository.Add(uow, &mastersAbroad)
		if err != nil {
			return err
		}
	}

	// Do checks with previous data only if previous masters abroad and current masters abroad exists.
	if tempMastersAbroad != nil && mastersAbroad != nil {
		// Create temp scores to get presvious scores of masters abroad.
		tempScores := []general.Score{}

		err := service.Repository.GetAllForTenant(uow, tenantID, &tempScores,
			repository.Filter("`masters_abroad_id`=?", tempMastersAbroad.ID))
		if err != nil {
			return err
		}

		// Make map to count occurences of scores id in previous and current scores.
		scoreIDMap := make(map[uuid.UUID]uint)

		// Count the number of occurence of previous score ID.
		for _, tempScore := range tempScores {
			scoreIDMap[tempScore.ID]++
		}

		// Count the number of occurrence of current score ID to know total count of occurrence of each ID.
		for _, score := range mastersAbroad.Scores {

			// If ID is valid then push its occurence in ID map.
			if util.IsUUIDValid(score.ID) {
				scoreIDMap[score.ID]++
			} else {
				// If ID is nil create new score entry in table.
				score.CreatedBy = credentialID
				score.TenantID = tenantID
				err = service.Repository.Add(uow, &score)
				if err != nil {
					return err
				}
			}

			// If number of occurrence is more than one (present in previous and current scores) then update score.
			if scoreIDMap[score.ID] > 1 {

				// Gice updated_by field to current scores.
				score.UpdatedBy = credentialID
				err = service.Repository.Update(uow, &score)
				if err != nil {
					return err
				}

				// Make the number of occurrences 0 after updating experience.
				scoreIDMap[score.ID] = 0
			}
		}

		// If the number of occurrence is one, the score was presnt in previous scores only, delete it.
		for _, tempScore := range tempScores {
			if scoreIDMap[tempScore.ID] == 1 {
				err = service.Repository.UpdateWithMap(uow, general.Score{}, map[string]interface{}{
					"DeletedBy": credentialID,
					"DeletedAt": time.Now(),
				}, repository.Filter("`id` = ?", tempScore.ID))
				if err != nil {
					return err
				}
				// Make the number of occurrences 0 after deleting scores.
				scoreIDMap[tempScore.ID] = 0
			}
		}

		// Replace countries and universities of masters abroad of enquiry.
		// Replace countries.
		if err := service.Repository.ReplaceAssociations(uow, mastersAbroad, "Countries",
			mastersAbroad.Countries); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}
		// Replace universities.
		if err := service.Repository.ReplaceAssociations(uow, mastersAbroad, "Universities",
			mastersAbroad.Universities); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}

		// Make universities empty to avoid unnecessary updates.
		mastersAbroad.Universities = nil

		// Make countries empty to avoid unnecessary updates.
		mastersAbroad.Countries = nil

		// Make scores empty to avoid unnecessary updates.
		mastersAbroad.Scores = nil

		// Give updated_by field to masters abroad.
		mastersAbroad.UpdatedBy = credentialID

		// Update masters abroad.
		err = service.Repository.Update(uow, &mastersAbroad)
		if err != nil {
			return err
		}
	}

	// Check if current talent's masters abroad exists or not, if doesnt exist then make it nil
	// and delete masters abroad from database.
	if mastersAbroad == nil && tempMastersAbroad != nil {
		// Make talent's master abroad id as nil.
		enquiry.IsMastersAbroad = false

		// Update score(s) for updating deleted_by field of score(s).
		err := service.Repository.UpdateWithMap(uow, &general.Score{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		},
			repository.Filter("`masters_abroad_id`=? AND `tenant_id`=?", tempMastersAbroad.ID, enquiry.TenantID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}

		// Delete masters abroad form database.
		// Update masters abroad for updating deleted_by field of masters abroad.
		if err := service.Repository.UpdateWithMap(uow, &general.MastersAbroad{}, map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
			repository.Filter("`enquiry_id`=? AND `tenant_id`=?", enquiryID, enquiry.TenantID)); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Masters Abroad details could not be deleted", http.StatusInternalServerError)
		}
	}
	return nil
}

// enquiryAddSearchQueries adds all search queries by comparing with the enquiry data recieved from request body.
func (service *EnquiryService) enquiryAddSearchQueries(tenantID uuid.UUID, enquiry *talenq.Search, roleName string) []repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcesors []repository.QueryProcessor

	fmt.Println("===================================SEARCH ENQUIRY===================================================")
	fmt.Println("enquiry ->", enquiry)
	fmt.Println("=================================================================================================")

	if enquiry.FirstName != nil {
		util.AddToSlice("talent_enquiries.`first_name`", "LIKE ?", "AND", "%"+*enquiry.FirstName+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if enquiry.LastName != nil {
		util.AddToSlice("talent_enquiries.`last_name`", "LIKE ?", "AND", "%"+*enquiry.LastName+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if enquiry.Email != nil {
		util.AddToSlice("talent_enquiries.`email`", "LIKE ?", "AND", "%"+*enquiry.Email+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if enquiry.IsExperience != nil {
		util.AddToSlice("talent_enquiries.`is_experience`", "= ?", "AND", *enquiry.IsExperience,
			&columnNames, &conditions, &operators, &values)
	}
	if enquiry.AcademicYears != nil && len(enquiry.AcademicYears) != 0 {
		util.AddToSlice("talent_enquiries.`academic_year`", "IN(?)", "AND", enquiry.AcademicYears,
			&columnNames, &conditions, &operators, &values)
	}
	// remove uuid.nil element from the IDs & add an IS NULL query.
	if len(enquiry.SalesPersonIDs) > 0 {
		util.AddToSlice("talent_enquiries.`sales_person_id`", "IN(?)", "AND", enquiry.SalesPersonIDs,
			&columnNames, &conditions, &operators, &values)
	}
	if enquiry.EnquirySource != nil {
		util.AddToSlice("talent_enquiries.`source_id`", "=?", "AND", enquiry.EnquirySource,
			&columnNames, &conditions, &operators, &values)
	}
	if enquiry.EnquiryFromDate != nil {
		util.AddToSlice("`enquiry_date`", ">= ?", "AND", enquiry.EnquiryFromDate, &columnNames, &conditions, &operators, &values)
	}
	if enquiry.EnquiryToDate != nil {
		util.AddToSlice("`enquiry_date`", "<= ?", "AND", enquiry.EnquiryToDate, &columnNames, &conditions, &operators, &values)
	}
	if enquiry.IsLastThirtyDays != nil && *enquiry.IsLastThirtyDays {
		queryProcesors = append(queryProcesors, repository.Filter("`enquiry_date` BETWEEN NOW() - INTERVAL 30 DAY AND NOW()"))
	}
	if enquiry.IsMastersAbroad != nil {
		util.AddToSlice("talent_enquiries.`is_masters_abroad`", "=?", "AND", enquiry.IsMastersAbroad,
			&columnNames, &conditions, &operators, &values)
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN masters_abroad ON talent_enquiries.`id` = masters_abroad.`enquiry_id`"),
			repository.Filter("masters_abroad.`deleted_at` IS NULL"),
			repository.Filter("masters_abroad.`tenant_id`=?", tenantID))

		if enquiry.YearOfMS != nil {
			util.AddToSlice("masters_abroad.`year_of_ms`", "=?", "AND", enquiry.YearOfMS,
				&columnNames, &conditions, &operators, &values)
		} else {
			queryProcesors = append(queryProcesors, repository.Filter("masters_abroad.`year_of_ms` IS NULL"))
		}
	}
	if enquiry.City != nil {
		util.AddToSlice("talent_enquiries.`city`", "LIKE ?", "AND", "%"+*enquiry.City+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if enquiry.CountryID != nil {
		util.AddToSlice("talent_enquiries.`country_id`", "=?", "AND", enquiry.CountryID,
			&columnNames, &conditions, &operators, &values)
	}
	// If college or qualifications is present then join talent_enquiry_academics table.
	if enquiry.College != nil || len(enquiry.Qualifications) != 0 {
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN talent_enquiry_academics ON talent_enquiries.`id` = talent_enquiry_academics.`enquiry_id`"),
			repository.Filter("talent_enquiry_academics.`deleted_at` IS NULL"),
			repository.Filter("talent_enquiry_academics.`tenant_id`=?", tenantID))

		if enquiry.College != nil {
			util.AddToSlice("talent_enquiry_academics.`college`", "LIKE ?", "AND", "%"+*enquiry.College+"%",
				&columnNames, &conditions, &operators, &values)
		}

		if len(enquiry.Qualifications) > 0 {
			util.AddToSlice("talent_enquiry_academics.`degree_id`", "IN(?)", "AND", enquiry.Qualifications,
				&columnNames, &conditions, &operators, &values)
		}
	}

	if enquiry.CallRecordOutcomeID != nil || enquiry.CallRecordPurposeID != nil {
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN talent_enquiry_call_records ON talent_enquiries.`id` = talent_enquiry_call_records.`enquiry_id`"),
			repository.Filter("talent_enquiry_call_records.`deleted_at` IS NULL"),
			repository.Filter("talent_enquiry_call_records.`tenant_id`=?", tenantID))

		if enquiry.CallRecordPurposeID != nil {
			util.AddToSlice("talent_enquiry_call_records.`purpose_id`", "=?", "AND", enquiry.CallRecordPurposeID,
				&columnNames, &conditions, &operators, &values)
		}
		if enquiry.CallRecordOutcomeID != nil {
			util.AddToSlice("talent_enquiry_call_records.`outcome_id`", "=?", "AND", enquiry.CallRecordOutcomeID,
				&columnNames, &conditions, &operators, &values)
		}
	}

	// If experince technologies is present then join talent_enquiry_experiences table,
	//talent_enquiry_experiences_technologies table.
	if len(enquiry.ExperienceTechnologies) != 0 || len(enquiry.Designations) != 0 {

		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN talent_enquiry_experiences ON talent_enquiries.`id` = talent_enquiry_experiences.`enquiry_id`"),
			repository.Join("LEFT JOIN talent_enquiry_experiences_technologies ON talent_enquiry_experiences.`id` = talent_enquiry_experiences_technologies.`experience_id`"),
			repository.Filter("talent_enquiry_experiences.`deleted_at` IS NULL"),
			repository.Filter("talent_enquiry_experiences.`tenant_id`=?", tenantID))
		if len(enquiry.ExperienceTechnologies) > 0 {
			fmt.Println("====================IN EXPERIENCE tech===============================")
			util.AddToSlice("talent_enquiry_experiences_technologies.`technology_id`", "IN(?)", "AND",
				enquiry.ExperienceTechnologies, &columnNames, &conditions, &operators, &values)
		}
		if len(enquiry.Designations) > 0 {
			fmt.Println("====================IN EXPERIENCE Designation===============================")
			util.AddToSlice("talent_enquiry_experiences.`designation_id`", "IN(?)", "AND",
				enquiry.Designations, &columnNames, &conditions, &operators, &values)
		}
	}

	// If technologies is present then join talent_enquiries_technologies and technolgies table.
	if len(enquiry.Technologies) != 0 {
		fmt.Println("====================IN REGULAR TECH===============================")
		queryProcesors = append(queryProcesors, repository.Join("INNER JOIN talent_enquiries_technologies ON talent_enquiries.`id` = talent_enquiries_technologies.`enquiry_id`"))

		if len(enquiry.Technologies) > 0 {
			util.AddToSlice("talent_enquiries_technologies.`technology_id`", "IN(?)", "AND", enquiry.Technologies, &columnNames, &conditions, &operators, &values)
		}
	}

	// If one or more waiting list fields are present then join waiting_list table.
	if enquiry.WaitingFor != nil || enquiry.WaitingForCompanyBranchID != nil || enquiry.WaitingForRequirementID != nil ||
		enquiry.WaitingForCourseID != nil || enquiry.WaitingForBatchID != nil || enquiry.WaitingForIsActive != nil ||
		enquiry.WaitingForFromDate != nil || enquiry.WaitingForToDate != nil {
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN waiting_list ON talent_enquiries.`id` = waiting_list.`enquiry_id`"),
			repository.Filter("waiting_list.`deleted_at` IS NULL"),
			repository.Filter("waiting_list.`tenant_id`=?", tenantID))

		// Waiting for is 'Company' search.
		if enquiry.WaitingFor != nil && *enquiry.WaitingFor == "Company" {
			queryProcesors = append(queryProcesors,
				repository.Filter("waiting_list.`company_branch_id` IS NOT NULL OR waiting_list.`company_requirement_id` IS NOT NULL"))
		}

		// Waiting for is 'Course' search.
		if enquiry.WaitingFor != nil && *enquiry.WaitingFor == "Course" {
			queryProcesors = append(queryProcesors,
				repository.Filter("waiting_list.`course_id` IS NOT NULL OR waiting_list.`batch_id` IS NOT NULL"))
		}

		// Company branch id search.
		if enquiry.WaitingForCompanyBranchID != nil {
			util.AddToSlice("waiting_list.`company_branch_id`", "=?", "AND", enquiry.WaitingForCompanyBranchID, &columnNames, &conditions, &operators, &values)
		}

		// Company requirement id search.
		if enquiry.WaitingForRequirementID != nil {
			util.AddToSlice("waiting_list.`company_requirement_id`", "=?", "AND", enquiry.WaitingForRequirementID, &columnNames, &conditions, &operators, &values)
		}

		// Course id search.
		if enquiry.WaitingForCourseID != nil {
			util.AddToSlice("waiting_list.`course_id`", "=?", "AND", enquiry.WaitingForCourseID, &columnNames, &conditions, &operators, &values)
		}

		// Batch id search.
		if enquiry.WaitingForBatchID != nil {
			util.AddToSlice("waiting_list.`batch_id`", "=?", "AND", enquiry.WaitingForBatchID, &columnNames, &conditions, &operators, &values)
		}

		// Is active search.
		if enquiry.WaitingForIsActive != nil {
			util.AddToSlice("waiting_list.`is_active`", "=?", "AND", enquiry.WaitingForIsActive, &columnNames, &conditions, &operators, &values)
		}

		// Applied for from date.
		if enquiry.WaitingForFromDate != nil {
			util.AddToSlice("waiting_list.`created_at`", ">= ?", "AND", enquiry.WaitingForFromDate, &columnNames, &conditions, &operators, &values)
		}

		// Applied for to date.
		if enquiry.WaitingForToDate != nil {
			util.AddToSlice("waiting_list.`created_at`", "<= ?", "AND", enquiry.WaitingForToDate, &columnNames, &conditions, &operators, &values)
		}
	}

	//min max experience coming soon*****************************************************************************

	queryProcesors = append(queryProcesors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("talent_enquiries.`id`"))
	return queryProcesors
}

// GetCount returns total count of result set.
func (service *EnquiryService) GetCount(out interface{}, count *int, queryProcessor ...repository.QueryProcessor) error {
	uow := repository.NewUnitOfWork(service.DB, true)
	return service.Repository.GetCount(uow, out, count, queryProcessor...)
}

// doesCredentialExist validates if credental exists or not in database.
func (service *EnquiryService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if credential(parent credential) exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCompanyBranchExist validates if company branch exists or not in database.
func (service *EnquiryService) doesCompanyBranchExist(companyBranchID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, company.Branch{},
		repository.Filter("`id` = ?", companyBranchID))
	if err := util.HandleError("Company branch not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *EnquiryService) doesTenantExist(tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if tenant(parent tenant) exists or not.
	exists, err := repository.DoesRecordExist(uow.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExist validates all foreign keys of enquiry.
func (service *EnquiryService) doForeignKeysExist(enquiry *talenq.Enquiry) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check parent source exists or not.
	if enquiry.SourceID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, enquiry.TenantID, general.Source{},
			repository.Filter("`id` = ?", enquiry.SourceID))
		if err := util.HandleError("Invalid source ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check parent salesperson exists or not.
	if enquiry.SalesPersonID != nil {
		exists, err := repository.DoesRecordExist(uow.DB, general.User{},
			repository.Join("left join roles on users.`role_id` = roles.`id`"),
			repository.Filter("users.`id`=? AND roles.`role_name`=? AND users.`tenant_id`=? AND roles.`tenant_id`=?",
				enquiry.SalesPersonID, "salesperson", enquiry.TenantID, enquiry.TenantID))
		if err := util.HandleError("Invalid salesperson ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check if country exists or not.
	if enquiry.CountryID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, enquiry.TenantID, general.Country{},
			repository.Filter("`id` = ?", enquiry.CountryID))
		if err := util.HandleError("Invalid country ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check if state exists or not.
	if enquiry.StateID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, enquiry.TenantID, general.State{},
			repository.Filter("`id` = ?", enquiry.StateID))
		if err := util.HandleError("Invalid state ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check if technologies exist or not.
	if enquiry.Technologies != nil && len(enquiry.Technologies) != 0 {
		if err := service.doTechnologiesExist(enquiry.Technologies, enquiry.TenantID); err != nil {
			return err
		}
	}

	// Check if courses exist or not.
	if enquiry.Courses != nil && len(enquiry.Courses) != 0 {
		var courseIDs []uuid.UUID
		for _, course := range enquiry.Courses {
			courseIDs = append(courseIDs, course.ID)
		}
		var count int = 0
		err := service.Repository.GetCountForTenant(uow, enquiry.TenantID, course.Course{}, &count,
			repository.Filter("`id` IN (?)", courseIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(enquiry.Courses) {
			log.NewLogger().Error("Course ID is invalid")
			return errors.NewValidationError("Course ID is invalid")
		}
	}

	// Check if experiences' technologies and designation exists or not.
	if enquiry.Experiences != nil && len(enquiry.Experiences) != 0 {
		for _, experience := range enquiry.Experiences {
			// Check if designation exists or not.
			exists, err := repository.DoesRecordExistForTenant(uow.DB, enquiry.TenantID, general.Designation{},
				repository.Filter("`id` = ?", experience.DesignationID))
			if err := util.HandleError("Invalid designation ID", exists, err); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}
			if experience.Technologies != nil && len(experience.Technologies) != 0 {
				if err := service.doTechnologiesExist(experience.Technologies, enquiry.TenantID); err != nil {
					return err
				}
			}
		}
	}

	// Check if academic's degree and specialization exists or not.
	if enquiry.Academics != nil && len(enquiry.Academics) != 0 {
		for _, academic := range enquiry.Academics {
			// Check if degree exists or not.
			exists, err := repository.DoesRecordExistForTenant(uow.DB, enquiry.TenantID, general.Degree{},
				repository.Filter("`id` = ?", academic.DegreeID))
			if err := util.HandleError("Invalid degree ID", exists, err); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}

			// Check if specialization exists or not only if not coming from enquiry form.
			exists, err = repository.DoesRecordExistForTenant(uow.DB, enquiry.TenantID, general.Specialization{},
				repository.Filter("`id`=? AND `degree_id`=?", academic.SpecializationID, academic.DegreeID))
			if err := util.HandleError("Invalid specialization ID", exists, err); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}
		}
	}

	// Check if masters abroad degree, universities and countries exists or not.
	if enquiry.MastersAbroad != nil {
		// Check if degree exists or not.
		exists, err := repository.DoesRecordExistForTenant(uow.DB, enquiry.TenantID, general.Degree{},
			repository.Filter("`id` = ?", enquiry.MastersAbroad.DegreeID))
		if err := util.HandleError("Invalid degree ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}

		// Check if universities exist or not.
		var universityIDs []uuid.UUID
		for _, university := range enquiry.MastersAbroad.Universities {
			universityIDs = append(universityIDs, university.ID)
		}
		var count int = 0
		err = service.Repository.GetCountForTenant(uow, enquiry.TenantID, general.University{}, &count,
			repository.Filter("`id` IN (?)", universityIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(enquiry.MastersAbroad.Universities) {
			log.NewLogger().Error("University ID is invalid")
			return errors.NewValidationError("University ID is invalid")
		}

		// Check if countries exist or not.
		var countryIDs []uuid.UUID
		for _, country := range enquiry.MastersAbroad.Countries {
			countryIDs = append(countryIDs, country.ID)
		}
		err = service.Repository.GetCountForTenant(uow, enquiry.TenantID, general.Country{}, &count,
			repository.Filter("`id` IN (?)", countryIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(enquiry.MastersAbroad.Countries) {
			log.NewLogger().Error("Country ID is invalid")
			return errors.NewValidationError("Country ID is invalid")
		}

		// Check if score's examination exist or not.
		var examinationIDs []uuid.UUID
		for _, score := range enquiry.MastersAbroad.Scores {
			examinationIDs = append(examinationIDs, score.ExaminationID)
		}
		err = service.Repository.GetCountForTenant(uow, enquiry.TenantID, general.Examination{}, &count,
			repository.Filter("`id` IN (?)", examinationIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(examinationIDs) {
			log.NewLogger().Error("Examination ID is invalid")
			return errors.NewValidationError("Examination ID is invalid")
		}
	}
	return nil
}

// doTechnologiesExist validates if technolgy exists or not in database.
func (service *EnquiryService) doTechnologiesExist(technologies []*general.Technology, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)
	// Keep all technology ids in one variable.
	var technologyIds []uuid.UUID
	for _, technology := range technologies {
		technologyIds = append(technologyIds, technology.ID)
	}
	// Get count for technologyIDs.
	var count int = 0
	err := service.Repository.GetCountForTenant(uow, tenantID, general.Technology{}, &count,
		repository.Filter("`id` IN (?)", technologyIds))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}
	if count != len(technologies) {
		log.NewLogger().Error("Technology ID is invalid")
		return errors.NewValidationError("Technology ID is invalid")
	}
	return nil
}

// doesTechnologyExist validates if technology exists or not in database.
func (service *EnquiryService) doesTechnologyExist(technologyID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Technology{},
		repository.Filter("`id` = ?", technologyID))
	if err := util.HandleError("Invalid technology ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesSalespersonExist validates if salesperson exists or not in database.
func (service *EnquiryService) doesSalespersonExist(salespersonID uuid.UUID, tenantID uuid.UUID) error {
	// Check parent sales person exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.User{},
		repository.Filter("`id` = ?", salespersonID))
	if err := util.HandleError("Invalid salesperson ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBatchExist validates if batch exists or not in database.
func (service *EnquiryService) doesBatchExist(tenantID uuid.UUID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Batch{}, repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Batch not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCompanyRequirementExist validates if company requirement exists or not in database.
func (service *EnquiryService) doesCompanyRequirementExist(requirementID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, company.Requirement{},
		repository.Filter("`id` = ?", requirementID))
	if err := util.HandleError("Company requirement not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCourseExist validates if course exists or not in database.
func (service *EnquiryService) doesCourseExist(courseID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Course{},
		repository.Filter("`id` = ?", courseID))
	if err := util.HandleError("Course not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesEnquiryExist validates if enquiry exists or not in database.
func (service *EnquiryService) doesEnquiryExist(enquiryID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talenq.Enquiry{}, repository.Filter("`id` = ?", enquiryID))
	if err := util.HandleError("Enquiry not found", exists, err); err != nil {
		return err
	}
	return nil
}

//setTenantID gives tenant id to all academics, experiences and masters abroad of enquiry.
func (service *EnquiryService) setTenantID(enquiry *talenq.Enquiry, tenantID uuid.UUID) {
	// If academics is present then give all the academics tenant id.
	if len(enquiry.Academics) != 0 {
		for _, academic := range enquiry.Academics {
			academic.TenantID = tenantID
		}
	}

	// If experiences is present then give all the experiences tenant id.
	if len(enquiry.Experiences) != 0 {
		for _, experience := range enquiry.Experiences {
			experience.TenantID = tenantID
		}
	}

	// If masters abroad is present then give masters abroad and scores tenant id.
	if enquiry.MastersAbroad != nil {
		enquiry.MastersAbroad.TenantID = tenantID
		for _, score := range enquiry.MastersAbroad.Scores {
			score.TenantID = tenantID
		}
	}
}

// // sortEnquiryChildTables sorts enquiry's academics and experineces.
func (service *EnquiryService) sortEnquiryChildTables(enquiries *[]talenq.DTO) {
	if enquiries != nil && len(*enquiries) != 0 {
		for i := 0; i < len(*enquiries); i++ {

			// Sort academics by order od passout year in ascending order.
			academics := &(*enquiries)[i].Academics
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

			// Sort experiences by order od from date in ascending order.
			experiences := &(*enquiries)[i].Experiences
			for j := 0; j < len(*experiences); j++ {
				if len((*experiences)[j].FromDate) == 0 {
					return
				}
			}
			for j := 0; j < len(*experiences); j++ {
				sort.Slice(*experiences, func(p, q int) bool {
					FromDateSmallInStr := (*experiences)[p].FromDate[:4]
					FromDateLargeInStr := (*experiences)[q].FromDate[:4]
					FromDateSmallInInt, _ := strconv.Atoi(FromDateSmallInStr)
					FromDateLargeInInt, _ := strconv.Atoi(FromDateLargeInStr)
					return FromDateSmallInInt < FromDateLargeInInt
				})
			}
		}
	}
}

// extractID extracts ID from object and removes data from the object.
// this is done so that the foreign key entity records are not updated in their respective tables
// when the college branch entity is being added or updated.
func (service *EnquiryService) extractID(enquiry *talenq.Enquiry) {
	// State field.
	if enquiry.State != nil {
		enquiry.StateID = &enquiry.State.ID
	}

	// Country field.
	if enquiry.Country != nil {
		enquiry.CountryID = &enquiry.Country.ID
	}

}

// setCollegeNameAndID gets the college by college name and sets the college id.
func (service *EnquiryService) setCollegeNameAndID(uow *repository.UnitOfWork, enquiry *talenq.Enquiry) error {
	academics := enquiry.Academics
	for i := 0; i < len(academics); i++ {
		// Get college branch from database.
		collegeBranch := list.Branch{}
		if err := service.Repository.GetRecordForTenant(uow, enquiry.TenantID, &collegeBranch,
			repository.Filter("`branch_name`=?", academics[i].College)); err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
		}

		// Give college id to academics collegeID field.
		academics[i].CollegeID = &collegeBranch.ID
	}
	return nil
}

// setTalentIDIfExists checks if talent exists by email and gets its id to give to enquiry.
func (service *EnquiryService) setTalentIDIfExists(uow *repository.UnitOfWork, enquiry *talenq.Enquiry) error {
	// Check id talent with same email exists or not.
	tempTalent := tal.Talent{}
	if err := service.Repository.GetRecordForTenant(uow, enquiry.TenantID, &tempTalent,
		repository.Filter("`email`=?", enquiry.Email),
		repository.Select("id")); err != nil {
		if err != gorm.ErrRecordNotFound {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
		}
	}

	if tempTalent.ID != uuid.Nil {
		enquiry.TalentID = &tempTalent.ID
	}
	return nil
}
