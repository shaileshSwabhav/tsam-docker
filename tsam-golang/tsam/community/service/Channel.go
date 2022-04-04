package service

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/repository"
)

// ChannelService gives acces to CRUD operations on channel.
type ChannelService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewChannelService returns new instance of ChannelService.
func NewChannelService(db *gorm.DB, repo repository.Repository) *ChannelService {
	return &ChannelService{
		DB:         db,
		Repository: repo,
	}
}

// AddChannel adds new channel to database.
func (service *ChannelService) AddChannel(channel *community.Channel) error {

	// Check all foreign key records.
	err := service.doForeignKeysExist(channel.TenantID, channel.CreatedBy)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(channel)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, channel)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// UpdateChannel updates channel in database.
func (service *ChannelService) UpdateChannel(channel *community.Channel) error {

	// Check if channelID exist.
	err := service.doesChannelExist(channel.ID)
	if err != nil {
		return err
	}

	// Check all foreign key records.
	err = service.doForeignKeysExist(channel.TenantID, channel.UpdatedBy)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(channel)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Update(uow, channel)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("An error occured", http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// DeleteChannel Delete Channel By ChannelID
func (service *ChannelService) DeleteChannel(channel *community.Channel) error {

	// Check if channelID exist.
	err := service.doesChannelExist(channel.ID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(channel.TenantID, channel.DeletedBy); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	// First update the deleted_by and deleted_at field to soft delete.
	err = service.Repository.UpdateWithMap(uow, channel, map[string]interface{}{
		"DeletedBy": channel.DeletedBy,
		"DeletedAt": time.Now(),
	})
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetChannels gets specific channels.
func (service *ChannelService) GetChannels(tenantID uuid.UUID,
	channels *[]community.Channel, parser *web.Parser, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()

	// Get all channels and orderby channel_name.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, channels, "`channel_name`",
		service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	uow.Commit()
	fmt.Println("Go routine Number is ---------------------------------------------------------", runtime.NumGoroutine())
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// validateFieldUniqueness checks field uniqueness.
func (service *ChannelService) validateFieldUniqueness(channel *community.Channel) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, channel.TenantID, community.Channel{},
		repository.Filter("`channel_name`=? AND `id`!= ?", channel.ChannelName,
			channel.ID))
	err = util.HandleIfExistsError("Record already exists with the name "+channel.ChannelName, exists, err)
	if err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *ChannelService) doForeignKeysExist(tenantID, credentialID uuid.UUID) error {

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

// doesChannelExist will check if channelID is valid
func (service *ChannelService) doesChannelExist(channelID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, community.Channel{},
		repository.Filter("`id`=?", channelID))
	if err := util.HandleError("Channel doesn't exist", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *ChannelService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Tenant doesn't exist", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *ChannelService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Credential doesn't exist", exists, err); err != nil {
		return err
	}
	return nil
}

// addSearchQueries adds search criteria to get all channels.
func (service *ChannelService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	channelName := requestForm.Get("channelName")
	if !util.IsEmpty(channelName) {
		util.AddToSlice("`channel_name`", "LIKE ?", "AND", "%"+channelName+"%", &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
