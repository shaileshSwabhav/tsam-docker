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

// SourceService provides methods to do different CRUD operations on sources table.
type SourceService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewSourceService returns a new instance Of SourceService.
func NewSourceService(db *gorm.DB, repository repository.Repository) *SourceService {
	return &SourceService{
		DB:         db,
		Repository: repository,
	}
}

// AddSource adds a new record in the sources table.
func (service *SourceService) AddSource(source *general.Source, uows ...*repository.UnitOfWork) error {
	tenantID := source.TenantID

	// check all foreign key records.
	err := service.doForeignKeysExist(tenantID, source.CreatedBy)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.doesSourceNameExist(source)
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
	err = service.Repository.Add(uow, source)
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

// UpdateSource updates an existing record in the source table.
func (service *SourceService) UpdateSource(source *general.Source) error {

	tenantID := source.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, source.UpdatedBy)
	if err != nil {
		return err
	}

	// Check source record.
	err = service.doesSourceExist(tenantID, source.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.doesSourceNameExist(source)
	if err != nil {
		return err
	}

	// Update method of repo to update.
	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, source)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetSource returns all source records for the given tenant.
func (service *SourceService) GetSource(source *general.Source) error {
	tenantID := source.TenantID

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if source exist.
	err = service.doesSourceExist(tenantID, source.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetForTenant(uow, tenantID, source.ID, source)
	if err != nil {
		return err
	}

	return nil
}

// AddSources adds multiple sources by calling the add service for each entity.
func (service *SourceService) AddSources(sources *[]*general.Source) error {
	// Check for same name conflict.
	for i := 0; i < len(*sources); i++ {
		for j := 0; j < len(*sources); j++ {
			if i != j && (*sources)[i].Name == (*sources)[j].Name {
				log.NewLogger().Error("Name:" + (*sources)[j].Name + " exists")
				return errors.NewValidationError("Name:" + (*sources)[j].Name + " exists")
			}
		}
	}

	// Add individual source.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, source := range *sources {
		err := service.AddSource(source, uow)
		if err != nil {
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetAllSources gets all sources from the DB for the specific tenant.
func (service *SourceService) GetAllSources(tenantID uuid.UUID, sources *[]general.Source) error {
	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Repo get all call.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, sources)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// DeleteSource deletes the specific source record from the table.
func (service *SourceService) DeleteSource(source *general.Source) error {
	tenantID := source.TenantID

	// Check all foreign key records.
	err := service.doForeignKeysExist(tenantID, source.DeletedBy)
	if err != nil {
		return err
	}

	// Check source record.
	err = service.doesSourceExist(tenantID, source.ID)
	if err != nil {
		return err
	}

	// Repository deleted_at update call.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update source for updating deleted_by and deleted_at fields of source
	err = service.Repository.UpdateWithMap(uow, source, map[string]interface{}{
		"DeletedBy": source.DeletedBy,
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
func (service *SourceService) doForeignKeysExist(tenantID, credentialID uuid.UUID) error {
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

//  doesTenantExist returns error if there is no tenant record in table.
func (service *SourceService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *SourceService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSourceExist returns error if there is no source record for the given tenant.
func (service *SourceService) doesSourceExist(tenantID, sourceID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Source{},
		repository.Filter("`id` = ?", sourceID))
	if err := util.HandleError("Invalid source ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSourceNameExist returns true if the source name already exists in database
func (service *SourceService) doesSourceNameExist(source *general.Source) error {
	//check for same source name conflict
	exists, err := repository.DoesRecordExistForTenant(service.DB, source.TenantID, &general.Source{},
		repository.Filter("`name`=? AND `id`!=?", source.Name, source.ID))
	if err := util.HandleIfExistsError("Source name exists", exists, err); err != nil {
		return err
	}
	return nil
}
