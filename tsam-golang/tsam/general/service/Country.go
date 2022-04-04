package service

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// CountryService Provide method to Update, Delete, Add, Get Method For Country.
type CountryService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// var countryError []string = []string{
// 	"invalid request", "Invalid Country ID", "an error occurred", "Country Doesn't Exist",
// }

// NewCountryService returns new instance of CountryService.
func NewCountryService(db *gorm.DB, repository repository.Repository) *CountryService {
	return &CountryService{
		DB:         db,
		Repository: repository,
	}
}

// AddCountry adds new country to database.
func (service *CountryService) AddCountry(country *general.Country, uows ...*repository.UnitOfWork) error {
	// Get credenial id from created_by field of country set in controller.
	credentialID := country.CreatedBy

	// Validate tenant ID.
	err := service.doesTenantExist(country.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, country.TenantID)
	if err != nil {
		return err
	}

	// Validate same country name exists.
	err = service.doesCountryNameExist(country)
	if err != nil {
		return err
	}

	// Creating unit of work.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Add Country.
	err = service.Repository.Add(uow, country)
	if err != nil {
		uow.RollBack()
		return err
	}

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddCountries adds multiple countries to database.
func (service *CountryService) AddCountries(countries *[]general.Country, countryIDs *[]uuid.UUID, tenantID, credentialID uuid.UUID) error {
	// Check for same name conflict.
	for i := 0; i < len(*countries); i++ {
		for j := 0; j < len(*countries); j++ {
			if i != j && (*countries)[i].Name == (*countries)[j].Name {
				log.NewLogger().Error("Name:" + (*countries)[j].Name + " exists")
				return errors.NewValidationError("Name:" + (*countries)[j].Name + " exists")
			}
		}
	}

	// Add individual country to database.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, country := range *countries {
		country.TenantID = tenantID
		country.CreatedBy = credentialID
		err := service.AddCountry(&country, uow)
		if err != nil {
			return err
		}
		*countryIDs = append(*countryIDs, country.ID)
	}

	uow.Commit()
	return nil
}

// UpdateCountry updates country to database.
func (service *CountryService) UpdateCountry(country *general.Country) error {
	// Validate tenant ID.
	err := service.doesTenantExist(country.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(country.UpdatedBy, country.TenantID)
	if err != nil {
		return err
	}

	// Validate country ID.
	err = service.doesCountryExist(country.ID, country.TenantID)
	if err != nil {
		return err
	}

	// Validate if same country name already exists.
	err = service.doesCountryNameExist(country)
	if err != nil {
		return err
	}

	// Update country.
	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Update(uow, country)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// DeleteCountry delete country from database.
func (service *CountryService) DeleteCountry(country *general.Country) error {
	// Get credential id form deleted_by field of country set in controller.
	credentialID := country.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(country.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(country.DeletedBy, country.TenantID)
	if err != nil {
		return err
	}

	// Validate country ID.
	err = service.doesCountryExist(country.ID, country.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update country for updating deleted_by and deleted_at fields of country
	if err := service.Repository.UpdateWithMap(uow, country, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", country.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Country could not be deleted", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// GetCountry returns one country by id.
func (service *CountryService) GetCountry(country *general.Country) error {
	// Validate tenant id.
	err := service.doesTenantExist(country.TenantID)
	if err != nil {
		return err
	}

	// Validate city id.
	err = service.doesCountryExist(country.ID, country.TenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get country.
	err = service.Repository.GetForTenant(uow, country.TenantID, country.ID, country)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetCountries returns all countries.
func (service *CountryService) GetCountries(countries *[]general.Country, tenantID uuid.UUID) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all countries.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, countries, "`name`")
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}
	uow.Commit()
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *CountryService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCountryExist validates if country exists or not in database.
func (service *CountryService) doesCountryExist(countryID uuid.UUID, tenantID uuid.UUID) error {
	// Check country exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Country{},
		repository.Filter("`id` = ?", countryID))
	if err := util.HandleError("Invalid country ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *CountryService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCountryNameExist returns true if the country name already exists for the state in database.
func (service *CountryService) doesCountryNameExist(country *general.Country) error {
	// Check for same country name conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, country.TenantID, &general.Country{},
		repository.Filter("`name`=? AND `id`!=?", country.Name, country.ID))
	if err := util.HandleIfExistsError("Name:"+country.Name+" exists", exists, err); err != nil {
		return err
	}
	return nil
}
