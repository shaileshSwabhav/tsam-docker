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

// DayService provides methods to update, delete, add, get method for day.
type DayService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewDayService returns new instance of DayService.
func NewDayService(db *gorm.DB, repository repository.Repository) *DayService {
	return &DayService{
		DB:         db,
		Repository: repository,
	}
}

// AddDay will add new record to the table.
func (service *DayService) AddDay(day *general.Day, uows ...*repository.UnitOfWork) error {

	// Check tenant exists.
	err := service.doesTenantExists(day.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(day.TenantID, day.CreatedBy)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.doesFieldUniquenessExist(day)
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
	err = service.Repository.Add(uow, day)
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

// AddDays adds days.
func (service *DayService) AddDays(days *[]general.Day, tenantID, credentialID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, false)
	for _, day := range *days {
		day.TenantID = tenantID
		day.CreatedBy = credentialID
		err := service.AddDay(&day, uow)
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// UpdateDay will update the specified record in the table.
func (service *DayService) UpdateDay(day *general.Day) error {

	// Check tenant exists.
	err := service.doesTenantExists(day.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(day.TenantID, day.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if day exist.
	err = service.doesDayExists(day.TenantID, day.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.doesFieldUniquenessExist(day)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Update(uow, day)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// DeleteDay deletes day.
func (service *DayService) DeleteDay(day *general.Day) error {

	// Check tenant exists.
	err := service.doesTenantExists(day.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(day.TenantID, day.DeletedBy)
	if err != nil {
		return err
	}

	// Check if day exist.
	err = service.doesDayExists(day.TenantID, day.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.UpdateWithMap(uow, general.Day{}, map[interface{}]interface{}{
		"DeletedBy": day.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", day.ID))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetDays will return all the days present in the table.
func (service *DayService) GetDays(days *[]general.Day, tenantID uuid.UUID) error {

	// Check tenant exists.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, days, "`order`")
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

// doesTenantExists validates tenant ID.
func (service *DayService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *DayService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesDayExists validates day ID.
func (service *DayService) doesDayExists(tenantID, dayID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Day{},
		repository.Filter("`id` = ?", dayID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesFieldUniquenessExist returns error if there is no general type record for the given tenant.
func (service *DayService) doesFieldUniquenessExist(day *general.Day) error {

	exists, err := repository.DoesRecordExistForTenant(service.DB, day.TenantID, general.Day{},
		repository.Filter("`order`=? AND `id`!=?", day.Order, day.ID))
	if err := util.HandleIfExistsError("Record already exists with same order.",
		exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}
