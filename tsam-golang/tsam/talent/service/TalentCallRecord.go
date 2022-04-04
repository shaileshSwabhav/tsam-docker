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

// TalentCallRecordService provides method to update, delete, add, get all, get one for talent call records.
type TalentCallRecordService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// talentCallRecordAssociations provides preload associations array for enquiry call record.
var talentCallRecordAssociations []string = []string{"Purpose", "Outcome"}

// NewTalentCallRecordService returns new instance of TalentCallRecordService.
func NewTalentCallRecordService(db *gorm.DB, repository repository.Repository) *TalentCallRecordService {
	return &TalentCallRecordService{
		DB:         db,
		Repository: repository,
	}
}

// AddTalentCallRecord adds one talent call record to database.
func (service *TalentCallRecordService) AddTalentCallRecord(talentCallRecord *tal.CallRecord) error {
	// Get credential id from CreatedBy field of talentCallRecord(set in controller).
	credentialID := talentCallRecord.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(talentCallRecord.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, talentCallRecord.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(talentCallRecord.TenantID, talentCallRecord.TalentID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(talentCallRecord); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add talent call record to database.
	if err := service.Repository.Add(uow, talentCallRecord); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent call record could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetTalentCallRecords gets all talent call records from database.
func (service *TalentCallRecordService) GetTalentCallRecords(talentCallRecords *[]tal.CallRecordDTO,
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

	// Get talent call records from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, talentCallRecords, "`date_time`",
		repository.Filter("`talent_id`=?", talentID),
		repository.PreloadAssociations(talentCallRecordAssociations)); err != nil {
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

// GetTalentCallRecord gets one talent call record form database.
func (service *TalentCallRecordService) GetTalentCallRecord(talentCallRecord *tal.CallRecord) error {
	// Validate tenant id.
	if err := service.doesTenantExist(talentCallRecord.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(talentCallRecord.TenantID, talentCallRecord.TalentID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get talent call record.
	if err := service.Repository.GetForTenant(uow, talentCallRecord.TenantID, talentCallRecord.ID, talentCallRecord,
		repository.PreloadAssociations(talentCallRecordAssociations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// UpdateTalentCallRecord updates talent call record in Database.
func (service *TalentCallRecordService) UpdateTalentCallRecord(talentCallRecord *tal.CallRecord) error {
	// Get credential id from UpdatedBy field of talentCallRecord(set in controller).
	credentialID := talentCallRecord.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(talentCallRecord.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, talentCallRecord.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(talentCallRecord.TenantID, talentCallRecord.TalentID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(talentCallRecord); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting talent call record already present in database.
	tempTalentCallRecord := tal.CallRecord{}

	// Get talent call record for getting created_by field of talent call record from database.
	if err := service.Repository.GetForTenant(uow, talentCallRecord.TenantID, talentCallRecord.ID, &tempTalentCallRecord); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp talent call record to talent call record to be updated.
	talentCallRecord.CreatedBy = tempTalentCallRecord.CreatedBy

	// Update Talent call record.
	if err := service.Repository.Save(uow, talentCallRecord); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent call reecord could not be updated", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// DeleteTalentCallRecord deletes one talent call record form database.
func (service *TalentCallRecordService) DeleteTalentCallRecord(talentCallRecord *tal.CallRecord) error {
	// Starting new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get credential id from DeletedBy field of talentCallRecord(set in controller).
	credentialID := talentCallRecord.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(talentCallRecord.TenantID); err != nil {
		return err
	}

	// Validate talent call record id.
	if err := service.doesTalentCallRecordExist(talentCallRecord.ID, talentCallRecord.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, talentCallRecord.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(talentCallRecord.TenantID, talentCallRecord.TalentID); err != nil {
		return err
	}

	// Update talent call record for updating deleted_by and deleted_at field of talent call record.
	if err := service.Repository.UpdateWithMap(uow, &tal.CallRecord{}, map[string]interface{}{
		"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`tenant_id`=? AND `id`=?", talentCallRecord.TenantID, talentCallRecord.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent call record could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *TalentCallRecordService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *TalentCallRecordService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentExist validates if talent exists or not in database.
func (service *TalentCallRecordService) doesTalentExist(tenantID uuid.UUID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tal.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentCallRecordExist validates if talent call record exists or not in database.
func (service *TalentCallRecordService) doesTalentCallRecordExist(talentCallRecordID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tal.CallRecord{},
		repository.Filter("`id` = ?", talentCallRecordID))
	if err := util.HandleError("Invalid talent call record ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExist validates if purpose ad outcome exists or not in database.
func (service *TalentCallRecordService) doForeignKeysExist(talentCallRecord *tal.CallRecord) error {
	// Check if purpose exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Purpose{}, repository.Filter("`id` = ?", talentCallRecord.PurposeID))
	if err := util.HandleError("Invalid purpose ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if outcome exists or not.
	exists, err = repository.DoesRecordExist(service.DB, general.Outcome{}, repository.Filter("`id` = ?", talentCallRecord.OutcomeID))
	if err := util.HandleError("Invalid outcome ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
