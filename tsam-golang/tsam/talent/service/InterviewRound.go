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

// InterviewRoundService provides method to update, delete, add, get all, get one for interview rounds.
type InterviewRoundService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewInterviewRoundService returns new instance of InterviewRoundService.
func NewInterviewRoundService(db *gorm.DB, repository repository.Repository) *InterviewRoundService {
	return &InterviewRoundService{
		DB:         db,
		Repository: repository,
	}
}

// AddInterviewRound adds one interview round to database.
func (service *InterviewRoundService) AddInterviewRound(interviewRound *tal.InterviewRound, uows ...*repository.UnitOfWork) error {

	// Creating unit of work.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Get credential id from CreatedBy field of interviewRound.
	credentialID := interviewRound.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(interviewRound.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, interviewRound.TenantID); err != nil {
		return err
	}

	// Add interview round to database.
	if err := service.Repository.Add(uow, interviewRound); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Interview round could not be added", http.StatusInternalServerError)
	}

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// GetInterviewRounds gets all interview rounds from database.
func (service *InterviewRoundService) GetInterviewRounds(interviewRounds *[]tal.InterviewRound, tenantID uuid.UUID) error {
	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	//Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get interview rounds from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, interviewRounds, "`name`"); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetInterviewRound gets one interview round form database.
func (service *InterviewRoundService) GetInterviewRound(interviewRound *tal.InterviewRound) error {
	// Validate tenant id.
	if err := service.doesTenantExist(interviewRound.TenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get interview round.
	if err := service.Repository.GetForTenant(uow, interviewRound.TenantID, interviewRound.ID, interviewRound); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// UpdateInterviewRound updates interview round in Database.
func (service *InterviewRoundService) UpdateInterviewRound(interviewRound *tal.InterviewRound) error {
	// Get credential id from UpdatedBy field of interviewRound(set in controller).
	credentialID := interviewRound.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(interviewRound.TenantID); err != nil {
		return err
	}

	// validate credential id.
	if err := service.doesCredentialExist(credentialID, interviewRound.TenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting interview round already present in database.
	tempInterviewRound := tal.InterviewRound{}

	// Get interview round for getting created_by field of interview round from database.
	if err := service.Repository.GetForTenant(uow, interviewRound.TenantID, interviewRound.ID, &tempInterviewRound); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp interview round to interview round to be updated.
	interviewRound.CreatedBy = tempInterviewRound.CreatedBy

	// Update Interview round.
	if err := service.Repository.Save(uow, interviewRound); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Interview round could not be updated", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// DeleteInterviewRound deletes one interview round form database.
func (service *InterviewRoundService) DeleteInterviewRound(interviewRound *tal.InterviewRound) error {

	// Get credential id from DeletedBy field of interviewRound(set in controller).
	credentialID := interviewRound.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(interviewRound.TenantID); err != nil {
		return err
	}

	// Validate interview round id.
	if err := service.doesInterviewRoundExist(interviewRound.ID, interviewRound.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, interviewRound.TenantID); err != nil {
		return err
	}

	//Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update interview round for updating deleted_by and deleted_at field of interview round.
	if err := service.Repository.UpdateWithMap(uow, &tal.InterviewRound{}, map[string]interface{}{
		"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`tenant_id`=? AND `id`=?", interviewRound.TenantID, interviewRound.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Interview round could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// AddInterviewRounds Adds multiple city to database.
func (service *InterviewRoundService) AddInterviewRounds(interviewRounds *[]tal.InterviewRound, interviewRoundIDs *[]uuid.UUID,
	tenantID, credentialID uuid.UUID) error {

	// Add individual interview round to database.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, interviewRound := range *interviewRounds {
		interviewRound.TenantID = tenantID
		interviewRound.CreatedBy = credentialID
		err := service.AddInterviewRound(&interviewRound, uow)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
		*interviewRoundIDs = append(*interviewRoundIDs, interviewRound.ID)
	}
	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *InterviewRoundService) doesTenantExist(tenantID uuid.UUID) error {
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
func (service *InterviewRoundService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
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

// doesInterviewRoundExist validates if interview round exists or not in database.
func (service *InterviewRoundService) doesInterviewRoundExist(interviewRoundID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check interview round exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, tal.InterviewRound{},
		repository.Filter("`id` = ?", interviewRoundID))
	if err := util.HandleError("Invalid interview round ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
