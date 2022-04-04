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

// PurposeService provides methods to do different CRUD operations on purposes table.
type PurposeService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewPurposeService returns a new instance Of PurposeService.
func NewPurposeService(db *gorm.DB, repository repository.Repository) *PurposeService {
	return &PurposeService{
		DB:         db,
		Repository: repository,
	}
}

// AddPurpose adds a new record in the purposes table.
func (service *PurposeService) AddPurpose(purpose *general.Purpose,
	uows ...*repository.UnitOfWork) error {
	tenantID := purpose.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, purpose.CreatedBy)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.doesPurposeNameExist(purpose)
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
	err = service.Repository.Add(uow, purpose)
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

// UpdatePurpose updates an existing record in the purpose table.
func (service *PurposeService) UpdatePurpose(purpose *general.Purpose) error {

	tenantID := purpose.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, purpose.UpdatedBy)
	if err != nil {
		return err
	}

	// Check purpose record.
	err = service.doesPurposeExist(tenantID, purpose.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.doesPurposeNameExist(purpose)
	if err != nil {
		return err
	}

	// Update method of repo to update.
	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Update(uow, purpose)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetPurpose returns all purpose records for the given tenant.
func (service *PurposeService) GetPurpose(purpose *general.Purpose) error {

	tenantID := purpose.TenantID

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	err = service.doesPurposeExist(tenantID, purpose.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetRecordForTenant(uow, tenantID, purpose)
	if err != nil {
		return err
	}
	return nil
}

// AddPurposes adds multiple purposes by calling the add service for each entity.
func (service *PurposeService) AddPurposes(purposes *[]*general.Purpose) error {
	// Check for same name conflict.
	for i := 0; i < len(*purposes); i++ {
		for j := 0; j < len(*purposes); j++ {
			if i != j && (*purposes)[i].Purpose == (*purposes)[j].Purpose && (*purposes)[i].PurposeType == (*purposes)[j].PurposeType {
				log.NewLogger().Error("Purpose:" + (*purposes)[j].Purpose + " exists")
				return errors.NewValidationError("Purpose:" + (*purposes)[j].Purpose + " exists")
			}
		}
	}
	// Add one purpose record at a time.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, purpose := range *purposes {
		err := service.AddPurpose(purpose, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// GetAllPurposes gets all purposes from the DB for the specific tenant.
func (service *PurposeService) GetAllPurposes(tenantID uuid.UUID, purposes *[]general.Purpose) error {
	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, purposes)
	if err != nil {
		return err
	}
	return nil
}

// GetPurposesByType returns all purposes based on purpose type  & tenant ID.
func (service *PurposeService) GetPurposesByType(tenantID uuid.UUID, purposeType string,
	purposes *[]general.Purpose) error {
	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all purposes by type form database.
	err = service.Repository.GetAllForTenant(uow, tenantID, purposes,
		repository.Filter("`purpose_type` = ?", purposeType))
	if err != nil {
		return err
	}
	return nil
}

// DeletePurpose deletes the specific purpose record from the table.
func (service *PurposeService) DeletePurpose(purpose *general.Purpose) error {
	tenantID := purpose.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, purpose.DeletedBy)
	if err != nil {
		return err
	}

	// Check purpose record.
	err = service.doesPurposeExist(tenantID, purpose.ID)
	if err != nil {
		return err
	}

	// Repository deleted_at update call.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update deleted_by field of purpose.
	err = service.Repository.UpdateWithMap(uow, purpose, map[string]interface{}{
		"DeletedBy": purpose.DeletedBy,
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
func (service *PurposeService) doForeignKeysExist(tenantID, credentialID uuid.UUID) error {
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
func (service *PurposeService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *PurposeService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesPurposeExist returns error if there is no purpose record for the given tenant.
func (service *PurposeService) doesPurposeExist(tenantID, purposeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Purpose{},
		repository.Filter("`id` = ?", purposeID))
	if err := util.HandleError("Invalid purpose ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesPurposeNameExist returns error if there is no purpose record for the given tenant.
func (service *PurposeService) doesPurposeNameExist(purpose *general.Purpose) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, purpose.TenantID, general.Purpose{},
		repository.Filter("`purpose_type`=? AND `purpose`=? AND `id`!=?", purpose.PurposeType,
			purpose.Purpose, purpose.ID))
	if err := util.HandleIfExistsError("Record already exists with the same type and purpose.",
		exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}
