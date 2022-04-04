package service

import (
	"net/url"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// TalentEventRegistrationService Provide method to Update, Delete, Add, Get method for talent_registration_event.
type TalentEventRegistrationService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewTalentEventRegistrationService returns new instance of TalentEventRegistrationService.
func NewTalentEventRegistrationService(db *gorm.DB, repository repository.Repository) *TalentEventRegistrationService {
	return &TalentEventRegistrationService{
		DB:         db,
		Repository: repository,
		association: []string{
			"Talent", "Event",
		},
	}
}

// AddTalentRegistration will add registration of talent for the specified event.
func (service *TalentEventRegistrationService) AddTalentRegistration(registration *talent.TalentEventRegistration) error {

	// Check if all foregin keys exist.
	err := service.doesForeignKeyExist(registration, registration.CreatedBy, true)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, registration)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateTalentRegistration will update the hasAttended field for specified talent and event.
func (service *TalentEventRegistrationService) UpdateTalentRegistration(registration *talent.TalentEventRegistration) error {

	// Check if all foregin keys exist.
	err := service.doesForeignKeyExist(registration, registration.UpdatedBy, false)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempRegistration := talent.TalentEventRegistration{}
	err = service.Repository.GetRecordForTenant(uow, registration.TenantID, &tempRegistration,
		repository.Filter("`id` = ?", registration.ID), repository.Select("`created_by`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	registration.CreatedBy = tempRegistration.CreatedBy

	err = service.Repository.Save(uow, registration)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTalentRegistration will for specified talent and event.
func (service *TalentEventRegistrationService) GetTalentRegistration(registration *talent.TalentEventRegistrationDTO,
	tenantID uuid.UUID, form url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &talent.TalentEventRegistration{},
		service.addSearchQueries(form))
	if err != nil {
		return err
	}
	if !exist {
		registration.IsTalentRegistered = false
		return nil
	}

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, registration, "`registration_date`",
		service.addSearchQueries(form))
	if err != nil {
		uow.RollBack()
		return err
	}
	registration.IsTalentRegistered = true

	uow.Commit()
	return nil
}

// GetTalentRegistrations will return all the talents registered with limit and offset.
func (service *TalentEventRegistrationService) GetTalentRegistrations(registration *[]talent.TalentEventRegistrationDTO,
	tenantID uuid.UUID, form url.Values, limit, offset int, totalCount *int) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, registration, "`registration_date`",
		service.addSearchQueries(form), repository.PreloadAssociations(service.association),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllTalentRegistration will return all the talents registered.
// func (service *TalentEventRegistrationService) GetAllTalentRegistration(registration *[]talent.TalentEventRegistrationDTO,
// 	tenantID uuid.UUID, form url.Values) error {

// 	// Check if tenant exist.
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(service.DB, true)

// 	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, registration, "`registration_date`",
// 		service.addSearchQueries(form), repository.PreloadAssociations(service.association))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	uow.Commit()
// 	return nil
// }

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doesForeignKeyExist will check all the foreign keys of event.
func (service *TalentEventRegistrationService) doesForeignKeyExist(registration *talent.TalentEventRegistration,
	credentialID uuid.UUID, isAddOperation bool) error {

	// Check if tenant exist.
	err := service.doesTenantExists(registration.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(registration.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if talent exist.
	err = service.doesTalentExist(registration.TenantID, registration.TalentID)
	if err != nil {
		return err
	}

	// Check if event exist.
	if isAddOperation {
		err = service.doesUpcomingEventExist(registration.TenantID, registration.EventID)
		if err != nil {
			return err
		}
	} else {
		err = service.doesLiveEventExist(registration.TenantID, registration.EventID)
		if err != nil {
			return err
		}

	}

	// Check if talent is registered for the event.
	err = service.checkIfTalentIsRegistered(registration.TenantID, registration.TalentID, registration.EventID, isAddOperation)
	if err != nil {
		return err
	}

	return nil
}

// doesTenantExists validates tenantID.
func (service *TalentEventRegistrationService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates credentialID.
func (service *TalentEventRegistrationService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Credential not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTalentExist validates talentID.
func (service *TalentEventRegistrationService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Talent not found", exists, err); err != nil {
		return err
	}
	return nil
}

// check if event is active and it is an upcoming event.
func (service *TalentEventRegistrationService) doesUpcomingEventExist(tenantID, eventID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, admin.SwabhavEvent{},
		repository.Filter("`id` = ? AND `event_status` = 'Upcoming' AND `is_active` = 1", eventID))
	if err := util.HandleError("Event not found", exists, err); err != nil {
		return err
	}
	return nil
}

// check if event is active and it is an upcoming event.
func (service *TalentEventRegistrationService) doesLiveEventExist(tenantID, eventID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, admin.SwabhavEvent{},
		repository.Filter("`id` = ? AND `event_status` = 'Live' AND `is_active` = 1", eventID))
	if err := util.HandleError("Event not found", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *TalentEventRegistrationService) checkIfTalentIsRegistered(tenantID, talentID, eventID uuid.UUID,
	isAddOperation bool) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.TalentEventRegistration{},
		repository.Filter("`talent_id` = ? AND `event_id` = ?", talentID, eventID))
	if isAddOperation {
		if err := util.HandleIfExistsError("You have already registered for the event", exists, err); err != nil {
			return err
		}
		return nil
	}
	if err := util.HandleError("Talent not registered for event the found", exists, err); err != nil {
		return err
	}

	return nil
}

// addSearchQueries will append search queries from queryParams to queryProcessor
func (service *TalentEventRegistrationService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("talent_id", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}

	if eventID, ok := requestForm["eventID"]; ok {
		util.AddToSlice("event_id", "= ?", "AND", eventID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
