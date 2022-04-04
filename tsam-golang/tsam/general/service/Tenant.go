package service

import (
	"net/http"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// TenantService Provide method to Update, Delete, Add, Get Method For Tenant.
type TenantService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewTenantService returns new instance of TenantService.
func NewTenantService(db *gorm.DB, repo repository.Repository) *TenantService {
	return &TenantService{
		DB:         db,
		Repository: repo,
	}
}

// GetTenant By Tenant ID
func (service *TenantService) GetTenant(tenant *general.Tenant, tenantID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, true)
	var err error
	err = service.Repository.Get(uow, tenantID, tenant)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetAllTenants Returns all the tenants
func (service *TenantService) GetAllTenants(tenant *[]general.Tenant) error {
	uow := repository.NewUnitOfWork(service.DB, true)
	var err error
	err = service.Repository.GetAll(uow, tenant)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// SearchTenants returns all the searched tenants
func (service *TenantService) SearchTenants(tenant *[]general.Tenant, searchTenant *general.Tenant) error {
	uow := repository.NewUnitOfWork(service.DB, true)
	err := service.Repository.GetAll(uow, tenant, service.tenantAddSearchQueries(searchTenant))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// AddTenant will add new tenant
func (service *TenantService) AddTenant(tenant *general.Tenant) error {
	var err error

	err = service.validateDuplicateValues(tenant)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	uow := repository.NewUnitOfWork(service.DB, false)
	tenant.ID = util.GenerateUUID()
	err = service.Repository.Add(uow, tenant)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// UpdateTenant will update tenant
func (service *TenantService) UpdateTenant(tenant *general.Tenant) error {
	var err error
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenant.ID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	if !exists {
		log.NewLogger().Error("Tenant not found")
		return errors.NewValidationError("Tenant not found")
	}
	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Get(uow, tenant.ID, &general.Tenant{})
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	err = service.Repository.Update(uow, tenant)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// DeleteTenant will soft delete the tenant
func (service *TenantService) DeleteTenant(tenantID uuid.UUID) error {
	var err error
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	if !exists {
		log.NewLogger().Error("Tenant not found")
		return errors.NewValidationError("Tenant not found")
	}
	tenant := general.Tenant{}
	tenant.ID = tenantID
	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Delete(uow, tenant)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// validateDuplicateValues checks for duplciate field values.
func (service *TenantService) validateDuplicateValues(tenant *general.Tenant) error {

	uow := repository.NewUnitOfWork(service.DB, true)
	var err error
	var count int
	// Check for duplicate name
	if err = service.Repository.GetCount(uow, &general.Tenant{}, &count,
		repository.Filter("tenant_name=?", tenant.TenantName)); err != nil {
		uow.RollBack()
		return err
	}
	if count > 0 {
		return errors.NewValidationError("Tenant Name already exists")
	}
	// Check for duplicate email
	if err = service.Repository.GetCount(uow, &general.Tenant{}, &count,
		repository.Filter("tenant_email=?", tenant.TenantEmail)); err != nil {
		uow.RollBack()
		return err
	}
	if count > 0 {
		return errors.NewValidationError("Tenant Email already exists")
	}
	// Check for duplicate contact
	if err = service.Repository.GetCount(uow, &general.Tenant{}, &count,
		repository.Filter("tenant_contact=?", tenant.TenantContact)); err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	if count > 0 {
		return errors.NewValidationError("Tenant Contact already exists")
	}
	return nil
}

// tenantAddSearchQueries adds all search queries by comparing with the Tenant data received from
func (service *TenantService) tenantAddSearchQueries(tenant *general.Tenant) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if tenant.TenantName != "" {
		util.AddToSlice("tenant_name", "LIKE ?", "AND", "%"+tenant.TenantName+"%", &columnNames, &conditions, &operators, &values)
	}
	if tenant.TenantEmail != "" {
		util.AddToSlice("tenant_name", "LIKE ?", "AND", "%"+tenant.TenantEmail+"%", &columnNames, &conditions, &operators, &values)
	}
	if tenant.TenantContact != "" {
		util.AddToSlice("tenant_name", "LIKE ?", "AND", "%"+tenant.TenantContact+"%", &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
