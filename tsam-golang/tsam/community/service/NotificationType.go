package service

import (
	"net/http"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// NotificationTypeService gives access to CRUD operations.
type NotificationTypeService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewNotificationTypeService returns an instance of NotificationTypeService.
func NewNotificationTypeService(db *gorm.DB, repo repository.Repository) *NotificationTypeService {
	return &NotificationTypeService{
		DB:         db,
		Repository: repo,
	}
}

// AddNotificationType adds new notification type.
func (service *NotificationTypeService) AddNotificationType(notificationType *community.NotificationType) error {

	// check if all foreign keys exist.
	err := service.doForeignKeysExist(notificationType.TenantID, notificationType.CreatedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, notificationType)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Unable to add notification type", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *NotificationTypeService) doForeignKeysExist(tenantID, credentialID uuid.UUID) error {

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
func (service *NotificationTypeService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *NotificationTypeService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential id", exists, err); err != nil {
		return err
	}
	return nil
}
