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

// StateService Provide method to Update, Delete, Add, Get Method For State.
type StateService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// var stateError []string = []string{
// 	"invalid request", "Invalid State ID", "an error occurred", "State Doesn't Exist",
// }

// NewStateService returns new instance of StateService.
func NewStateService(db *gorm.DB, repository repository.Repository) *StateService {
	return &StateService{
		DB:         db,
		Repository: repository,
	}
}

// AddState adds new state to database.
func (service *StateService) AddState(state *general.State, uows ...*repository.UnitOfWork) error {
	// Validate tenant id.
	err := service.doesTenantExist(state.TenantID)
	if err != nil {
		return err
	}

	// Validate country id.
	err = service.doesCountryExist(state.CountryID, state.TenantID)
	if err != nil {
		return err
	}

	// Validate if same state name exists.
	err = service.doesStateNameExist(state)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(state.CreatedBy, state.TenantID)
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

	// Add State.
	err = service.Repository.Add(uow, state)
	if err != nil {
		uow.RollBack()
		return err
	}

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddStates adds multiple states to database.
func (service *StateService) AddStates(states *[]general.State, statesIDs *[]uuid.UUID, tenantID, credentialID uuid.UUID) error {
	// Check for same name conflict.
	for i := 0; i < len(*states); i++ {
		for j := 0; j < len(*states); j++ {
			if i != j && (*states)[i].Name == (*states)[j].Name && (*states)[i].CountryID == (*states)[j].CountryID {
				log.NewLogger().Error("Name:" + (*states)[j].Name + " exists")
				return errors.NewValidationError("Name:" + (*states)[j].Name + " exists")
			}
		}
	}

	// Add individual state.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, state := range *states {
		state.TenantID = tenantID
		state.CreatedBy = credentialID
		err := service.AddState(&state, uow)
		if err != nil {
			return err
		}
		*statesIDs = append(*statesIDs, state.ID)
	}

	uow.Commit()
	return nil
}

// UpdateState updates state to database.
func (service *StateService) UpdateState(state *general.State) error {
	// Validate tenant ID.
	err := service.doesTenantExist(state.TenantID)
	if err != nil {
		return err
	}

	// Validate state ID.
	err = service.doesStateExist(state.ID, state.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(state.UpdatedBy, state.TenantID)
	if err != nil {
		return err
	}

	// Validate if same state name exists.
	err = service.doesStateNameExist(state)
	if err != nil {
		return err
	}

	// Validate country id.
	err = service.doesCountryExist(state.CountryID, state.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update state.
	err = service.Repository.Update(uow, state)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteState deletes state from database.
func (service *StateService) DeleteState(state *general.State) error {
	credentialID := state.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(state.TenantID)
	if err != nil {
		return err
	}

	// Validate state ID.
	err = service.doesStateExist(state.ID, state.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, state.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update state for updating deleted_by and deleted_at fields of state
	if err := service.Repository.UpdateWithMap(uow, state, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", state.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("State could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetState returns one state.
func (service *StateService) GetState(state *general.State) error {
	// Validate tenant id.
	err := service.doesTenantExist(state.TenantID)
	if err != nil {
		return err
	}

	// Validate city id.
	err = service.doesStateExist(state.ID, state.TenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get state by id.
	err = service.Repository.GetForTenant(uow, state.TenantID, state.ID, state)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetStateByCountryID returns all state by country id.
func (service *StateService) GetStateByCountryID(state *[]general.State, tenantID, countryID uuid.UUID) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Validate state id.
	err = service.doesCountryExist(countryID, tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all states by country id from database.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, state, "`name`",
		repository.Filter("`country_id` = ?", countryID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetStates return all state.
func (service *StateService) GetStates(states *[]general.State, tenantID uuid.UUID) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all states from database.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, states, "`name`")
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// doesStateNameExist returns true if the state name already exists for the country in database.
func (service *StateService) doesStateNameExist(state *general.State) error {
	// Check for same state name conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, state.TenantID, &general.State{},
		repository.Filter("`name`=? AND `country_id` =? AND `id`!=?", state.Name, state.CountryID, state.ID))
	if err := util.HandleIfExistsError("Name:"+state.Name+" exists", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *StateService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesStateExist validates if state exists or not in database.
func (service *StateService) doesStateExist(stateID, tenantID uuid.UUID) error {
	// Check state exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.State{},
		repository.Filter("`id` = ?", stateID))
	if err := util.HandleError("Invalid state ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCountryExist validates if country exists or not in database.
func (service *StateService) doesCountryExist(countryID, tenantID uuid.UUID) error {
	// Check country exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Country{},
		repository.Filter("`id` = ?", countryID))
	if err := util.HandleError("Invalid country ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *StateService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}
