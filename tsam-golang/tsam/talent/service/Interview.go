package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	general "github.com/techlabs/swabhav/tsam/models/general"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// InterviewService provides method to update, delete, add, get all, get one for interview.
type InterviewService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewInterviewService returns new instance of InterviewService.
func NewInterviewService(db *gorm.DB, repository repository.Repository) *InterviewService {
	return &InterviewService{
		DB:         db,
		Repository: repository,
	}
}

// AddInterview adds one interview to database.
func (service *InterviewService) AddInterview(interview *tal.Interview) error {
	// Get credential id from CreatedBy field of interview(set in controller).
	credentialID := interview.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(interview.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, interview.TenantID); err != nil {
		return err
	}

	// Validate interview schedule id.
	if err := service.doesInterviewScheduleExist(interview.TenantID, interview.ScheduleID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(interview); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add interview to database.
	if err := service.Repository.Add(uow, interview); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Interview could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetInterviews gets all interview from database.
func (service *InterviewService) GetInterviews(interviews *[]tal.Interview,
	tenantID uuid.UUID, interviewScheduleID uuid.UUID, uows ...*repository.UnitOfWork) error {
	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
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

	// Validate interview schedule id.
	if err := service.doesInterviewScheduleExist(tenantID, interviewScheduleID); err != nil {
		return err
	}

	// Get interview from database.
	if err := service.Repository.GetAllForTenant(uow, tenantID, interviews,
		repository.Filter("`schedule_id`=?", interviewScheduleID),
		repository.PreloadAssociations([]string{"TakenBy"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// GetInterview gets one interview form database.
func (service *InterviewService) GetInterview(interview *tal.Interview, form url.Values) error {

	// Validate tenant id.
	if err := service.doesTenantExist(interview.TenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get interview.
	if err := service.Repository.GetRecordForTenant(uow, interview.TenantID, interview,
		service.addSearchQueries(form),
		repository.PreloadAssociations([]string{"TakenBy"})); err != nil {
		uow.RollBack()
		if err == gorm.ErrRecordNotFound {
			uow.Commit()
			return nil
		}
		return errors.NewValidationError("Internal server error")
	}

	uow.Commit()
	return nil
}

// UpdateInterview updates interview in Database.
func (service *InterviewService) UpdateInterview(interview *tal.Interview) error {
	// Get credential id from UpdatedBy field of interview(set in controller).
	credentialID := interview.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(interview.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, interview.TenantID); err != nil {
		return err
	}

	// Validate interview schedule id.
	if err := service.doesInterviewScheduleExist(interview.TenantID, interview.ScheduleID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(interview); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting interview already present in database.
	tempInterview := tal.Interview{}

	// Get interview for getting created_by field of interview from database.
	if err := service.Repository.GetForTenant(uow, interview.TenantID, interview.ID, &tempInterview); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//replace taken by associations
	if err := service.Repository.ReplaceAssociations(uow, interview, "TakenBy",
		interview.TakenBy); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Make taken by nil to avoid insertion during update talent.
	interview.TakenBy = nil

	// Give created_by id from temp interview to interview to be updated.
	interview.CreatedBy = tempInterview.CreatedBy

	// Update interview.
	if err := service.Repository.Save(uow, interview); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Interview could not be updated", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// DeleteInterview deletes one interview form database.
func (service *InterviewService) DeleteInterview(interview *tal.Interview) error {
	// Starting new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get credential id from DeletedBy field of interview(set in controller).
	credentialID := interview.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(interview.TenantID); err != nil {
		return err
	}

	// Validate interview id.
	if err := service.doesInterviewExist(interview.ID, interview.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, interview.TenantID); err != nil {
		return err
	}

	// Validate interview schedule id.
	if err := service.doesInterviewScheduleExist(interview.TenantID, interview.ScheduleID); err != nil {
		return err
	}

	// Update interview for updating deleted_by and deleted_at field of interview.
	if err := service.Repository.UpdateWithMap(uow, &tal.Interview{}, map[string]interface{}{
		"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`tenant_id`=? AND `id`=?", interview.TenantID, interview.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Interview could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *InterviewService) doesTenantExist(tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if tenant(parent tenant) exists or not.
	exists, err := repository.DoesRecordExist(uow.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *InterviewService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
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

// doesInterviewScheduleExist validates if interview schedule exists or not in database.
func (service *InterviewService) doesInterviewScheduleExist(tenantID uuid.UUID, interviewSchdeuleID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check parent interview schedule exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, tal.InterviewSchedule{},
		repository.Filter("`id` = ?", interviewSchdeuleID))
	if err := util.HandleError("Invalid interview schedule ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesInterviewExist validates if interview exists or not in database.
func (service *InterviewService) doesInterviewExist(interviewID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check interview exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, tal.Interview{},
		repository.Filter("`id` = ?", interviewID))
	if err := util.HandleError("Invalid interview ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExist validates if talent and interview round are present or not in database.
func (service *InterviewService) doForeignKeysExist(interview *tal.Interview) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check parent talent exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, interview.TenantID, tal.Talent{},
		repository.Filter("`id` = ?", interview.TalentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check interview round exists or not.
	exists, err = repository.DoesRecordExistForTenant(uow.DB, interview.TenantID, tal.InterviewRound{},
		repository.Filter("`id` = ?", interview.RoundID))
	if err := util.HandleError("Invalid interview round ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// GetAllInterviewersList gets all interviewers list from credential table from database.
func (service *InterviewService) GetAllInterviewersList(interviewers *[]general.Credential, tenantID uuid.UUID) error {
	// Create new unit of work, if no transaction has been passed to the function.

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get interviewers list from database.
	// if err := service.Repository.GetAllInOrder(uow, interviewers, "`first_name`",
	// 	repository.Join("INNER JOIN employees ON credentials.`employee_id`=employees.`id` AND "+
	// 		"credentials.`tenant_id` = employees.`tenant_id`"),
	// 	repository.Filter("employees.`tenant_id` = ? AND employees.`deleted_at` IS NULL AND "+
	// 		"employees.`type` = ?", tenantID, "Developer"),
	// 	repository.Filter("`sales_person_id` OR `faculty_id` OR `employee_id` IS NOT NULL")); err != nil {
	// 	uow.RollBack()
	// 	log.NewLogger().Error(err.Error())
	// 	return errors.NewValidationError("Record not found")
	// }

	err := service.Repository.GetAllInOrder(uow, interviewers, "credentials.`first_name`",
		repository.Join("LEFT JOIN faculties f ON f.`id` = credentials.`faculty_id` AND credentials.`tenant_id` = f.`tenant_id`"),
		repository.Join("LEFT JOIN employees e ON e.`id` = credentials.`employee_id` AND credentials.`tenant_id` = e.`tenant_id`"),
		repository.Join("LEFT JOIN users u ON credentials.`sales_person_id` = u.`id` AND credentials.`tenant_id` = u.`tenant_id`"),
		repository.Filter("credentials.`deleted_at` IS NULL AND f.`deleted_at` IS NULL AND e.`deleted_at` IS NULL"+
			" AND u.`deleted_at` IS NULL AND credentials.`tenant_id` = ?", tenantID),
		repository.Filter("((credentials.`faculty_id` IS NOT NULL AND f.`is_active` = ?)"+
			" OR (credentials.`employee_id` IS NOT NULL AND e.`type` = ? AND e.`is_active` = ?)"+
			"	OR (credentials.`sales_person_id` IS NOT NULL  AND u.`is_active` = ?))", 1, "Developer", 1, 1),
		repository.PreloadAssociations([]string{"Role"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds search queries.
func (service *InterviewService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	// Interview schedule ID.
	if interviewScheduleID, ok := requestForm["interviewScheduleID"]; ok {
		util.AddToSlice("`schedule_id`", "= ?", "AND", interviewScheduleID, &columnNames, &conditions, &operators, &values)
	}

	// Talent ID.
	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("`talent_id`", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}

	// Status.
	if _, ok := requestForm["status"]; ok {
		util.AddToSlice("`status`", "LIKE ?", "AND", "%"+requestForm.Get("status")+"%", &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
