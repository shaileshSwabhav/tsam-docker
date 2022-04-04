package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// SwabhavEventService Provide method to Update, Delete, Add, Get method for event.
type SwabhavEventService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewSwabhavEventService returns new instance of EventService.
func NewSwabhavEventService(db *gorm.DB, repository repository.Repository) *SwabhavEventService {
	return &SwabhavEventService{
		DB:         db,
		Repository: repository,
		association: []string{
			"Country", "State",
		},
	}
}

// AddEvent will add new event in the table.
func (service *SwabhavEventService) AddEvent(event *admin.SwabhavEvent) error {

	// Checks if all foreign keys are present.
	err := service.doesForeignKeyExist(event, event.CreatedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, event)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateEvent will update specified event in the table.
func (service *SwabhavEventService) UpdateEvent(event *admin.SwabhavEvent) error {

	// Checks if all foreign keys are present.
	err := service.doesForeignKeyExist(event, event.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if event exist.
	err = service.doesEventExist(event.TenantID, event.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	var tempEvent admin.SwabhavEvent

	err = service.Repository.GetRecordForTenant(uow, event.TenantID, &tempEvent,
		repository.Select([]string{"`created_by`"}), repository.Filter("`id` = ?", event.ID))
	if err != nil {
		return err
	}

	event.CreatedBy = tempEvent.CreatedBy

	err = service.Repository.Save(uow, event)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteEvent will delete specified event in the table.
func (service *SwabhavEventService) DeleteEvent(event *admin.SwabhavEvent) error {

	// Checks if all foreign keys are present.
	err := service.doesForeignKeyExist(event, event.DeletedBy)
	if err != nil {
		return err
	}

	// Check if event exist.
	err = service.doesEventExist(event.TenantID, event.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, admin.SwabhavEvent{}, map[string]interface{}{
		"DeletedBy": event.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", event.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Delete talent registrations for the specified event.
	err = service.Repository.UpdateWithMap(uow, talent.TalentEventRegistration{}, map[string]interface{}{
		"DeletedBy": event.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`event_id` = ?", event.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetEvent will return specified event.
func (service *SwabhavEventService) GetEvent(event *admin.SwabhavEventDTO, tenantID, eventID uuid.UUID, form url.Values) error {

	//  Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetRecordForTenant(uow, tenantID, event, service.addSearchQueries(form),
		repository.Filter("`id` = ?", eventID), repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.GetCountForTenant(uow, tenantID, talent.TalentEventRegistration{},
		&event.TotalRegistrations, repository.Filter("`event_id` = ?", event.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetEvents will return events with limit and offset.
func (service *SwabhavEventService) GetEvents(events *[]admin.SwabhavEventDTO, tenantID uuid.UUID, form url.Values,
	limit, offset int, totalCount *int) error {

	//  Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, events, service.addSearchQueries(form),
		repository.OrderBy("CAST(swabhav_events.`from_date` AS DATE)"), repository.PreloadAssociations(service.association),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *events {
		err = service.Repository.GetCountForTenant(uow, tenantID, talent.TalentEventRegistration{},
			&(*events)[index].TotalRegistrations, repository.Filter("`swabhav_event_id` = ?", (*events)[index].ID))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetAllEvents will return all events.
func (service *SwabhavEventService) GetAllEvents(events *[]admin.SwabhavEventDTO, tenantID uuid.UUID, form url.Values) error {

	//  Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, events, service.addSearchQueries(form),
		repository.PreloadAssociations(service.association), repository.OrderBy("CAST(swabhav_events.`from_date` AS DATE)"))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *events {
		err = service.Repository.GetCountForTenant(uow, tenantID, talent.TalentEventRegistration{},
			&(*events)[index].TotalRegistrations, repository.Filter("`event_id` = ?", (*events)[index].ID))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doesForeignKeyExist will check all the foreign keys of event.
func (service *SwabhavEventService) doesForeignKeyExist(event *admin.SwabhavEvent, credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExists(event.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(event.TenantID, credentialID)
	if err != nil {
		return err
	}

	return nil
}

// doesTenantExists validates tenantID.
func (service *SwabhavEventService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates credentialID.
func (service *SwabhavEventService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Credential not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesEventExist validates eventID.
func (service *SwabhavEventService) doesEventExist(tenantID, eventID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, admin.SwabhavEvent{},
		repository.Filter("`id` = ?", eventID))
	if err := util.HandleError("Event not found", exists, err); err != nil {
		return err
	}
	return nil
}

// addSearchQueries will append search queries from queryParams to queryProcessor
func (service *SwabhavEventService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["title"]; ok {
		util.AddToSlice("`title`", "LIKE ?", "AND", "%"+requestForm.Get("title")+"%", &columnNames, &conditions, &operators, &values)
	}

	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("CAST(swabhav_events.`from_date` AS DATE)", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}

	if toDate, ok := requestForm["toDate"]; ok {
		util.AddToSlice("CAST(swabhav_events.`to_date` AS DATE)", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}

	if lastRegistrationDate, ok := requestForm["lastRegistrationDate"]; ok {
		util.AddToSlice("CAST(swabhav_events.`last_registration_date` AS DATE)", ">= ?", "AND", lastRegistrationDate, &columnNames, &conditions, &operators, &values)
	}

	if isOnline, ok := requestForm["isOnline"]; ok {
		util.AddToSlice("is_online", "= ?", "AND", isOnline, &columnNames, &conditions, &operators, &values)
	}

	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("is_active", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	if eventStatus, ok := requestForm["eventStatus"]; ok {
		util.AddToSlice("event_status", "= ?", "AND", eventStatus, &columnNames, &conditions, &operators, &values)
	}

	if eventID, ok := requestForm["eventID"]; ok {
		util.AddToSlice("id", "= ?", "AND", eventID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
