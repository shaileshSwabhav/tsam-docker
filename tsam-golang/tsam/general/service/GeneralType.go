package service

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// GeneralTypeService provides methods to do different CRUD operations on general_types table.
type GeneralTypeService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewGeneralTypeService returns a new instance Of GeneralTypeService.
func NewGeneralTypeService(db *gorm.DB, repository repository.Repository) *GeneralTypeService {
	return &GeneralTypeService{
		DB:         db,
		Repository: repository,
	}
}

// AddGeneralType adds a new record in the general_type table.
func (service *GeneralTypeService) AddGeneralType(generalType *general.CommonType,
	uows ...*repository.UnitOfWork) error {
	tenantID := generalType.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, generalType.CreatedBy)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.doesFieldUniquenessExist(generalType)
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
	err = service.Repository.Add(uow, generalType)
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

// UpdateGeneralType updates an existing record in the general_type table.
func (service *GeneralTypeService) UpdateGeneralType(generalType *general.CommonType) error {

	tenantID := generalType.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, generalType.UpdatedBy)
	if err != nil {
		return err
	}

	// Check general type record.
	err = service.doesGeneralTypeExist(tenantID, generalType.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.doesFieldUniquenessExist(generalType)
	if err != nil {
		return err
	}

	// Update method of repo to update.
	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Update(uow, generalType)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetGeneralType returns all general_type records for the given tenant.
func (service *GeneralTypeService) GetGeneralType(generalType *general.CommonType) error {

	tenantID := generalType.TenantID

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	err = service.doesGeneralTypeExist(tenantID, generalType.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetRecordForTenant(uow, tenantID, generalType)
	if err != nil {
		return err
	}
	return nil
}

// AddGeneralTypes adds multiple general types by calling the add service for each entity.
func (service *GeneralTypeService) AddGeneralTypes(generalTypes []*general.CommonType) error {
	// Add one general type record at a time.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, generalType := range generalTypes {
		err := service.AddGeneralType(generalType, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// GetAllGeneralTypes gets all general types from the DB for the specific tenant.
func (service *GeneralTypeService) GetAllGeneralTypes(tenantID uuid.UUID, generalTypes *[]general.CommonType) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Repo get all call.
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllForTenant(uow, tenantID, generalTypes)
	if err != nil {
		return err
	}
	return nil
}

// GetGeneralTypesByType returns all general types based on type name & tenant ID.
func (service *GeneralTypeService) GetGeneralTypesByType(tenantID uuid.UUID, typeName string,
	generalTypes *[]general.CommonType) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Unit of work created.
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, generalTypes,
		"key", repository.Filter("`type` = ?", typeName))
	if err != nil {
		return err
	}
	return nil
}

// DeleteGeneralType deletes the specific general type record from the table.
func (service *GeneralTypeService) DeleteGeneralType(generalType *general.CommonType) error {
	tenantID := generalType.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, generalType.DeletedBy)
	if err != nil {
		return err
	}

	// Check general type record.
	err = service.doesGeneralTypeExist(tenantID, generalType.ID)
	if err != nil {
		return err
	}

	// Repository deleted_at update call.
	uow := repository.NewUnitOfWork(service.DB, false)

	// First update the deleted_by field of record and then soft delete.
	err = service.Repository.UpdateWithMap(uow, generalType, map[string]interface{}{
		"DeletedBy": generalType.DeletedBy,
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
func (service *GeneralTypeService) doForeignKeysExist(tenantID, credentialID uuid.UUID) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *GeneralTypeService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

//doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *GeneralTypeService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesGeneralTypeExist returns error if there is no general type record for the given tenant.
func (service *GeneralTypeService) doesGeneralTypeExist(tenantID, generalTypeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.CommonType{},
		repository.Filter("`id` = ?", generalTypeID))
	if err := util.HandleError("Invalid general type ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesFieldUniquenessExist returns error if there is no general type record for the given tenant.
func (service *GeneralTypeService) doesFieldUniquenessExist(generalType *general.CommonType) error {

	// updating the value of an already existing type would give error as the existing record(`generalType`) would
	// have a record in table which matches `key` `value` `type` fields in the table. Hence, NOT IN `id` is used to
	// exempt the `generalType` record
	exists, err := repository.DoesRecordExistForTenant(service.DB, generalType.TenantID, general.CommonType{},
		repository.Filter("`type`=? AND (`key`=? OR `value`=?) AND `id`!=?", generalType.Type,
			generalType.Key, generalType.Value, generalType.ID))
	if err := util.HandleIfExistsError("Record already exists with the same type and key OR value.",
		exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}
