package services

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	model "github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// DomainService provide method like Add, Update, Delete, GetByID, GetAll for Domain
type DomainService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewDomainService returns the new instance of DomainService
func NewDomainService(db *gorm.DB, repo repository.Repository) *DomainService {
	return &DomainService{
		DB:         db,
		Repository: repo,
	}
}

// AddDomain adds new Domain to database
func (service *DomainService) AddDomain(domain *model.Domain, uows ...*repository.UnitOfWork) error {

	credentialID := domain.CreatedBy

	// check all foreign key records
	err := service.doForeignKeysExist(credentialID, domain)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(domain)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	length := len(uows)
	if length == 0 {
		uows = append(uows, repository.NewUnitOfWork(service.DB, false))
	}
	uow := uows[0]

	// Generate ID.
	domain.ID = util.GenerateUUID()

	// Add repo call.
	err = service.Repository.Add(uow, domain)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// UpdateDomain update the data of Domain
func (service *DomainService) UpdateDomain(domain *model.Domain) error {

	credentialID := domain.UpdatedBy

	// check all foreign key records
	err := service.doForeignKeysExist(credentialID, domain)
	if err != nil {
		return err
	}

	// check if domain exist
	err = service.doesDomainExist(domain.TenantID, domain.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(domain)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, domain)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// DeleteDomain delete the data of Domain
func (service *DomainService) DeleteDomain(domain *model.Domain) error {

	credentialID := domain.DeletedBy

	// check all foreign key records
	err := service.doForeignKeysExist(credentialID, domain)
	if err != nil {
		return err
	}

	// check if domain exist
	err = service.doesDomainExist(domain.TenantID, domain.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update the deleted_by and deleted_at field of the record.
	// Deleting the domain
	err = service.Repository.UpdateWithMap(uow, domain, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	})
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Unable to delete domain", http.StatusBadRequest)
	}

	uow.Commit()
	return nil
}

// GetAllDomains returns all Domain from database
func (service *DomainService) GetAllDomains(tenantID uuid.UUID, domains *[]*model.Domain,
	form url.Values, limit, offset int, totalCount *int) error {

	// check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Get after adding paging limit and offset.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, domains, "`domain_name`",
		service.addSearchQueries(form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetDomainList returns all Domain from database
func (service *DomainService) GetDomainList(tenantID uuid.UUID, domains *[]*model.Domain) error {

	// check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Get after adding paging limit and offset.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, domains, "`domain_name`")
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetDomain returns particular Domain by ID
func (service *DomainService) GetDomain(domain *model.Domain) error {

	tenantID := domain.TenantID

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if domain exist
	err = service.doesDomainExist(tenantID, domain.ID)
	if err != nil {
		return err
	}

	// Get called.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetRecordForTenant(uow, tenantID, domain)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

func (service *DomainService) validateFieldUniqueness(domain *model.Domain) error {

	exists, err := repository.DoesRecordExistForTenant(service.DB, domain.TenantID, model.Domain{},
		repository.Filter("`domain_name`=? AND `id`!= ?", domain.DomainName, domain.ID))
	if err := util.HandleIfExistsError("Record already exist with same domain name", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError(err.Error())
	}
	return nil
}

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *DomainService) doForeignKeysExist(credentialID uuid.UUID, domain *model.Domain) error {
	tenantID := domain.TenantID

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

// returns error if there is no tenant record in table.
func (service *DomainService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *DomainService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no domain record for the given tenant.
func (service *DomainService) doesDomainExist(tenantID, domainID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, model.Domain{},
		repository.Filter("`id`=?", domainID))
	if err := util.HandleError("Invalid domain ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// adds all search queries if any when getAll is called
// Need to test properly.
func (service *DomainService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if domainName, ok := requestForm["domainName"]; ok {
		util.AddToSlice("`domain_name`", "LIKE ?", "AND", "%"+domainName[0]+"%", &columnNames, &conditions, &operators, &values)
	}
	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
