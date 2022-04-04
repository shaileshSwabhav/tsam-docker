package service

import (
	"net/http"
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

// InterviewScheduleService provides method to update, delete, add, get all, get one for interview schedule.
type InterviewScheduleService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewInterviewScheduleService returns new instance of InterviewScheduleService.
func NewInterviewScheduleService(db *gorm.DB, repository repository.Repository) *InterviewScheduleService {
	return &InterviewScheduleService{
		DB:         db,
		Repository: repository,
	}
}

// AddInterviewSchedule adds one interview schedule to database.
func (service *InterviewScheduleService) AddInterviewSchedule(interviewSchedule *tal.InterviewSchedule) error {
	// Get credential id from CreatedBy field of interviewSchedule(set in controller).
	credentialID := interviewSchedule.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(interviewSchedule.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, interviewSchedule.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(interviewSchedule.TenantID, interviewSchedule.TalentID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add interview schedule to database.
	if err := service.Repository.Add(uow, interviewSchedule); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Interview Schedule could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetInterviewSchedules gets all interview schedules from database.
func (service *InterviewScheduleService) GetInterviewSchedules(interviewSchedules *[]tal.InterviewSchedule,
	tenantID uuid.UUID, talentID uuid.UUID, uows ...*repository.UnitOfWork) error {
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

	// Validate talent id.
	if err := service.doesTalentExist(tenantID, talentID); err != nil {
		return err
	}

	// Get interview schedules from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, interviewSchedules, "`scheduled_date`",
		repository.Filter("`talent_id`=?", talentID)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Remove timestamp from ScheduledDate field of interview schedules.
	if interviewSchedules != nil && len(*interviewSchedules) != 0 {
		for i := 0; i < len(*interviewSchedules); i++ {
			(*interviewSchedules)[i].ScheduledDate = (*interviewSchedules)[i].ScheduledDate[:10]
		}
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// GetInterviewSchedule gets one interview schedule form database.
func (service *InterviewScheduleService) GetInterviewSchedule(interviewSchedule *tal.InterviewSchedule) error {
	// Validate tenant id.
	if err := service.doesTenantExist(interviewSchedule.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(interviewSchedule.TenantID, interviewSchedule.TalentID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get interview schedule.
	if err := service.Repository.GetForTenant(uow, interviewSchedule.TenantID, interviewSchedule.ID, interviewSchedule); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Remove timestamp from ScheduledDate field of interview schedules.
	interviewSchedule.ScheduledDate = interviewSchedule.ScheduledDate[:10]

	uow.Commit()
	return nil
}

// UpdateInterviewSchedule updates interview schedule in Database.
func (service *InterviewScheduleService) UpdateInterviewSchedule(interviewSchedule *tal.InterviewSchedule) error {
	// Get credential id from UpdatedBy field of interviewSchedule(set in controller).
	credentialID := interviewSchedule.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(interviewSchedule.TenantID); err != nil {
		return err
	}

	// Validate interview schedule id.
	if err := service.doesInterviewScheduleExist(interviewSchedule.ID, interviewSchedule.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, interviewSchedule.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(interviewSchedule.TenantID, interviewSchedule.TalentID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting interview schedule already present in database.
	tempInterviewSchedule := tal.InterviewSchedule{}

	// Get interview schedule for getting created_by field of interview schedule from database.
	if err := service.Repository.GetForTenant(uow, interviewSchedule.TenantID, interviewSchedule.ID, &tempInterviewSchedule); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp interview schedule to interview schedule to be updated.
	interviewSchedule.CreatedBy = tempInterviewSchedule.CreatedBy

	// Update interview schedule.
	if err := service.Repository.Save(uow, interviewSchedule); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Interview Schedule could not be updated", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// DeleteInterviewSchedule deletes one interview schedule form database.
func (service *InterviewScheduleService) DeleteInterviewSchedule(interviewSchedule *tal.InterviewSchedule, uows ...*repository.UnitOfWork) error {
	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Get credential id from DeletedBy field of interviewSchedule(set in controller).
	credentialID := interviewSchedule.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(interviewSchedule.TenantID); err != nil {
		return err
	}

	// Validate interview schedule id.
	if err := service.doesInterviewScheduleExist(interviewSchedule.ID, interviewSchedule.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, interviewSchedule.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(interviewSchedule.TenantID, interviewSchedule.TalentID); err != nil {
		return err
	}

	//***********************************************Delete interviews********************************************
	if err := service.Repository.UpdateWithMap(uow, &tal.Interview{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`schedule_id`=? AND `tenant_id`=?", interviewSchedule.ID, interviewSchedule.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Interview Schedule could not be deleted", http.StatusInternalServerError)
	}

	// Get interview schedule for updating deleted_by field of interview schedule.
	if err := service.Repository.GetForTenant(uow, interviewSchedule.TenantID, interviewSchedule.ID, interviewSchedule); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Update interview schedule for updating deleted_by and deleted_at field of interview schedule.
	if err := service.Repository.UpdateWithMap(uow, interviewSchedule, map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()}); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Interview Schedule could not be deleted", http.StatusInternalServerError)
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *InterviewScheduleService) doesTenantExist(tenantID uuid.UUID) error {
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
func (service *InterviewScheduleService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
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

// doesTalentExist validates if talent exists or not in database.
func (service *InterviewScheduleService) doesTalentExist(tenantID uuid.UUID, talentID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check parent talent exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, tal.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesInterviewScheduleExist validates if interview schedule exists or not in database.
func (service *InterviewScheduleService) doesInterviewScheduleExist(interviewScheduleID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check interview schedule exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, tal.InterviewSchedule{},
		repository.Filter("`id` = ?", interviewScheduleID))
	if err := util.HandleError("Invalid interview schedule ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
