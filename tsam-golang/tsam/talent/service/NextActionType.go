package service

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// NextActionTypeService provides methods to update, delete, add, get method for nextActionType.
type NextActionTypeService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewNextActionTypeService returns new instance of NextActionTypeService.
func NewNextActionTypeService(db *gorm.DB, repository repository.Repository) *NextActionTypeService {
	return &NextActionTypeService{
		DB:         db,
		Repository: repository,
	}
}

// AddNextActionType will add new record to the table.
func (service *NextActionTypeService) AddNextActionType(nextActionType *talent.NextActionType, uows ...*repository.UnitOfWork) error {

	// Check if tenant exists.
	err := service.doesTenantExists(nextActionType.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(nextActionType.TenantID, nextActionType.CreatedBy)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(nextActionType)
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
	err = service.Repository.Add(uow, nextActionType)
	if err != nil {
		log.NewLogger().Error(err.Error())
		if length == 0 {
			uow.RollBack()
		}
		return err
	}
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddNextActionTypes adds nextActionTypes.
func (service *NextActionTypeService) AddNextActionTypes(nextActionTypes *[]talent.NextActionType, tenantID, credentialID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, false)
	for _, nextActionType := range *nextActionTypes {
		nextActionType.TenantID = tenantID
		nextActionType.CreatedBy = credentialID
		err := service.AddNextActionType(&nextActionType, uow)
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// UpdateNextActionType will update the specified record in the table.
func (service *NextActionTypeService) UpdateNextActionType(nextActionType *talent.NextActionType) error {

	// Check tenant exists.
	err := service.doesTenantExists(nextActionType.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(nextActionType.TenantID, nextActionType.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if nextActionType exist.
	err = service.doesNextActionTypeExists(nextActionType.TenantID, nextActionType.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(nextActionType)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Update(uow, nextActionType)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// DeleteNextActionType deletes nextActionType.
func (service *NextActionTypeService) DeleteNextActionType(nextActionType *talent.NextActionType) error {

	// Check tenant exists.
	err := service.doesTenantExists(nextActionType.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(nextActionType.TenantID, nextActionType.DeletedBy)
	if err != nil {
		return err
	}

	// Check if nextActionType exists.
	err = service.doesNextActionTypeExists(nextActionType.TenantID, nextActionType.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.UpdateWithMap(uow, talent.NextActionType{}, map[interface{}]interface{}{
		"DeletedBy": nextActionType.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", nextActionType.ID))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetNextActionTypes will return all the nextActionTypes present in the table.
func (service *NextActionTypeService) GetNextActionTypes(nextActionTypes *[]talent.NextActionType, tenantID uuid.UUID) error {

	// Check tenant exists.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, nextActionTypes, "`type`")
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

// doesTenantExists checks whether tenantID is valid.
func (service *NextActionTypeService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *NextActionTypeService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesNextActionTypeExists checks whether nextActionTypeID is valid.
func (service *NextActionTypeService) doesNextActionTypeExists(tenantID, nextActionTypeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.NextActionType{},
		repository.Filter("`id` = ?", nextActionTypeID))
	if err := util.HandleError("Invalid next action ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// validateFieldUniqueness returns error if a record with same type exists.
func (service *NextActionTypeService) validateFieldUniqueness(nextActionType *talent.NextActionType) error {

	exists, err := repository.DoesRecordExistForTenant(service.DB, nextActionType.TenantID, talent.NextActionType{},
		repository.Filter("`type`=? AND `id`!=?", nextActionType.Type, nextActionType.ID))
	if err := util.HandleIfExistsError("Record already exists with same type.",
		exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}
