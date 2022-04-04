package service

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/blog"
	"github.com/techlabs/swabhav/tsam/models/general"
	notificationdemo "github.com/techlabs/swabhav/tsam/notification_demo"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

type NotificationService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

func NewNotificationService(db *gorm.DB, repo repository.Repository) *NotificationService {
	return &NotificationService{
		DB:         db,
		Repository: repo,
	}
}

func (service *NotificationService) AddNotification(childTable *[]byte, notification *general.Notification_Test, uow *repository.UnitOfWork) error {
	// uow := repository.NewUnitOfWork(service.DB, false)

	var err error
	child := blog.BlogNotification{}
	err = json.Unmarshal(*childTable, &child)
	if err != nil {
		uow.DB.Rollback()
		return err

	}
	err = service.Repository.Add(uow, &child)
	if err != nil {
		uow.DB.Rollback()
		return err

	}

	notification.BlogNotificationID = &child.ID
	// notification.TypeID=
	err = service.Repository.Add(uow, &notification)
	if err != nil {
		uow.DB.Rollback()
		return err

	}
	uow.Commit()
	return nil
}

func (service *NotificationService) GetAllNotifications(tenantID uuid.UUID, userID uuid.UUID,
	notifications *[]notificationdemo.Notification_Test_DTO) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	err = service.doesCredentialExist(tenantID, userID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllForTenant(uow, tenantID, notifications, repository.Filter("notifier_id=?", userID),
		repository.PreloadAssociations([]string{"BlogNotifications"}))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *NotificationService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *NotificationService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
