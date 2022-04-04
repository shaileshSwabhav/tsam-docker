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

// CityService provides methods to update, delete, add, get method for city.
type CityService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewCityService returns new instance of CityService.
func NewCityService(db *gorm.DB, repository repository.Repository) *CityService {
	return &CityService{
		DB:         db,
		Repository: repository,
	}
}

// AddCity adds new city to database.
func (service *CityService) AddCity(city *general.City, uows ...*repository.UnitOfWork) error {
	// Validate tenant id.
	err := service.doesTenantExist(city.TenantID)
	if err != nil {
		return err
	}

	// Validate state id.
	err = service.doesStateExist(city.StateID, city.TenantID)
	if err != nil {
		return err
	}

	// Validate if same city name exists.
	err = service.doesCityNameExist(city)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(city.CreatedBy, city.TenantID)
	if err != nil {
		return err
	}

	//  Creating unit of work.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Add city.
	err = service.Repository.Add(uow, city)
	if err != nil {
		uow.RollBack()
		return err
	}

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddCities adds multiple cities to database.
func (service *CityService) AddCities(cities *[]general.City, cityIDs *[]uuid.UUID, tenantID, credentialID uuid.UUID) error {
	// Check for same name conflict.
	for i := 0; i < len(*cities); i++ {
		for j := 0; j < len(*cities); j++ {
			if i != j && (*cities)[i].Name == (*cities)[j].Name && (*cities)[i].StateID == (*cities)[j].StateID {
				log.NewLogger().Error("Name:" + (*cities)[j].Name + " exists")
				return errors.NewValidationError("Name:" + (*cities)[j].Name + " exists")
			}
		}
	}

	// Add individual city.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, city := range *cities {
		city.TenantID = tenantID
		city.CreatedBy = credentialID
		err := service.AddCity(&city, uow)
		if err != nil {
			return err
		}
		*cityIDs = append(*cityIDs, city.ID)
	}

	uow.Commit()
	return nil
}

// UpdateCity updates city to database.
func (service *CityService) UpdateCity(city *general.City) error {

	// Validate tenant ID.
	err := service.doesTenantExist(city.TenantID)
	if err != nil {
		return err
	}

	// Validate city ID.
	err = service.doesCityExist(city.ID, city.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(city.UpdatedBy, city.TenantID)
	if err != nil {
		return err
	}

	// Validate if same city name exists.
	err = service.doesCityNameExist(city)
	if err != nil {
		return err
	}

	// Validate state id.
	err = service.doesStateExist(city.StateID, city.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update city.
	err = service.Repository.Update(uow, city)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteCity deletes city from database.
func (service *CityService) DeleteCity(city *general.City) error {
	credentialID := city.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(city.TenantID)
	if err != nil {
		return err
	}

	// Validate city ID.
	err = service.doesCityExist(city.ID, city.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, city.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update city for updating deleted_by and deleted_at fields of city
	if err := service.Repository.UpdateWithMap(uow, city, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", city.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("City could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetCity returns one city by specific id.
func (service *CityService) GetCity(city *general.City) error {
	// Validate tenant id.
	err := service.doesTenantExist(city.TenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get one city by city id from database.
	err = service.Repository.GetForTenant(uow, city.TenantID, city.ID, city)
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetCitiesByStateID returns all cities by state id.
func (service *CityService) GetCitiesByStateID(cities *[]general.City, tenantID, stateID uuid.UUID) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Validate state id.
	err = service.doesStateExist(stateID, tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all cities by state id from database.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, cities, "`name`",
		repository.Filter("`state_id` = ?", stateID))
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetCities returns all cities.
func (service *CityService) GetCities(cities *[]general.City, tenantID uuid.UUID) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all cities from database.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, cities, "`name`")
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// doesCityNameExist returns true if the city name already exists for the state in database.
func (service *CityService) doesCityNameExist(city *general.City) error {
	// Check for same city name conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, city.TenantID, &general.City{},
		repository.Filter("`name`=? AND `state_id` =? AND `id`!=?", city.Name, city.StateID, city.ID))
	if err := util.HandleIfExistsError("Name:"+city.Name+" exists", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *CityService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCityExist validates if city exists or not in database.
func (service *CityService) doesCityExist(cityID uuid.UUID, tenantID uuid.UUID) error {
	// Check city exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.City{},
		repository.Filter("`id` = ?", cityID))
	if err := util.HandleError("Invalid city ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesStateExist validates if state exists or not in database.
func (service *CityService) doesStateExist(stateID uuid.UUID, tenantID uuid.UUID) error {
	// Check state exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.State{},
		repository.Filter("`id` = ?", stateID))
	if err := util.HandleError("Invalid state ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *CityService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}
