package service

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// OutcomeService provides methods to do different CRUD operations on outcomes table.
type OutcomeService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewOutcomeService returns a new instance Of OutcomeService.
func NewOutcomeService(db *gorm.DB, repository repository.Repository) *OutcomeService {
	return &OutcomeService{
		DB:         db,
		Repository: repository,
	}
}

// AddOutcome adds a new record in the outcomes table.
func (service *OutcomeService) AddOutcome(outcome *general.Outcome,
	uows ...*repository.UnitOfWork) error {
	tenantID := outcome.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, outcome.CreatedBy, outcome.PurposeID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.doesOutcomeNameExist(outcome)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Repository call.
	err = service.Repository.Add(uow, outcome)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return err
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// UpdateOutcome updates an existing record in the outcome table.
func (service *OutcomeService) UpdateOutcome(outcome *general.Outcome) error {

	tenantID := outcome.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, outcome.UpdatedBy, outcome.PurposeID)
	if err != nil {
		return err
	}

	// Check outcome record.
	err = service.doesOutcomeExist(tenantID, outcome.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.doesOutcomeNameExist(outcome)
	if err != nil {
		return err
	}

	// Update method of repo to update.
	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, outcome)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetOutcome returns all outcome records for the given tenant.
func (service *OutcomeService) GetOutcome(outcome *general.Outcome) error {

	tenantID := outcome.TenantID

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	err = service.doesOutcomeExist(tenantID, outcome.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetRecordForTenant(uow, tenantID, outcome)
	if err != nil {
		return err
	}
	return nil
}

// AddOutcomes adds multiple outcomes by calling the add service for each entity.
func (service *OutcomeService) AddOutcomes(purposeID uuid.UUID, outcomes *[]*general.Outcome) error {
	// Check for same name conflict.
	for i := 0; i < len(*outcomes); i++ {
		for j := 0; j < len(*outcomes); j++ {
			if i != j && (*outcomes)[i].Outcome == (*outcomes)[j].Outcome && (*outcomes)[i].PurposeID == (*outcomes)[j].PurposeID {
				log.NewLogger().Error("Outcome:" + (*outcomes)[j].Outcome + " exists")
				return errors.NewValidationError("Outcome:" + (*outcomes)[j].Outcome + " exists")
			}
		}
	}

	// Add one outcome record at a time.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, outcome := range *outcomes {
		outcome.PurposeID = purposeID
		err := service.AddOutcome(outcome, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// GetAllOutcomes gets all outcomes from the DB for the specific tenant.
func (service *OutcomeService) GetAllOutcomes(tenantID uuid.UUID, outcomes *[]general.Outcome) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Repo get all call.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, outcomes)
	if err != nil {
		return err
	}
	return nil
}

// GetAllOutcomesByPurpose returns all outcomes based on purposeID & tenant ID.
func (service *OutcomeService) GetAllOutcomesByPurpose(tenantID uuid.UUID, purposeID uuid.UUID,
	outcomes *[]general.Outcome) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Validate purpse if
	err = service.doesPurposeExist(tenantID, purposeID)
	if err != nil {
		return err
	}

	// Unit of work created.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, outcomes,
		repository.Filter("`purpose_id` = ?", purposeID))
	if err != nil {
		return err
	}
	return nil
}

// DeleteOutcome deletes the specific outcome record from the table.
func (service *OutcomeService) DeleteOutcome(outcome *general.Outcome) error {
	tenantID := outcome.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, outcome.DeletedBy, outcome.PurposeID)
	if err != nil {
		return err
	}

	// Check outcome record.
	err = service.doesOutcomeExist(tenantID, outcome.ID)
	if err != nil {
		return err
	}

	// Repository deleted_at update call.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update outcome for updating deleted_by and deleted_at fields of outcome
	err = service.Repository.UpdateWithMap(uow, outcome, map[string]interface{}{
		"DeletedBy": outcome.DeletedBy,
		"DeletedAt": time.Now(),
	})
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *OutcomeService) doForeignKeysExist(tenantID, credentialID, purposeID uuid.UUID) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Check if purpose exists.
	if err := service.doesPurposeExist(tenantID, purposeID); err != nil {
		return err
	}

	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *OutcomeService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *OutcomeService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesPurposeExist returns error if there is no purpose record in table for the given tenant.
func (service *OutcomeService) doesPurposeExist(tenantID, purposeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Purpose{},
		repository.Filter("`id` = ?", purposeID))
	if err := util.HandleError("Invalid purpose ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesOutcomeExist returns error if there is no outcome record for the given tenant.
func (service *OutcomeService) doesOutcomeExist(tenantID, outcomeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Outcome{},
		repository.Filter("`id` = ?", outcomeID))
	if err := util.HandleError("Invalid outcome ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesOutcomeNameExist returns error if there is no outcome record for the given tenant.
func (service *OutcomeService) doesOutcomeNameExist(outcome *general.Outcome) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, outcome.TenantID, general.Outcome{},
		repository.Filter("`purpose_id`=? AND `outcome`=? AND `id`!=?", outcome.PurposeID,
			outcome.Outcome, outcome.ID))
	if err := util.HandleIfExistsError("Record already exists with the same purpose id and outcome.",
		exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}
