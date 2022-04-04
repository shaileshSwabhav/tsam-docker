package service

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// LoginSessionService Provide method to Update, Delete, Add, Get Method For login sesion.
type LoginSessionService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewLoginSessionService returns new instance of LoginSession.
func NewLoginSessionService(db *gorm.DB, repository repository.Repository) *LoginSessionService {
	return &LoginSessionService{
		DB:         db,
		Repository: repository,
	}
}

// AddLoginSession Add New Login session to Database.
// 	transasction commit/rollback will be done by the caller.
func (service *LoginSessionService) AddLoginSession(loginSession *general.LoginSession, uow *repository.UnitOfWork) error {

	// Validate login session.
	if err := service.doesForeignKeyExist(loginSession, loginSession.CreatedBy); err != nil {
		return err
	}

	// give start time for session.
	loginSession.StartTime = time.Now()
	//add new login session
	if err := service.Repository.Add(uow, loginSession); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Login session could not be added")
	}
	return nil
}

// EndLoginSession updates end time of login session.
// 	ids should be validated by the caller.
// 	transasction commit/rollback will be done by the caller.
func (service *LoginSessionService) EndLoginSession(uow *repository.UnitOfWork, credentialID, tenantID,
	loginSessionID uuid.UUID) error {

	// assuming tenantID has been validated already.
	//Get login session
	loginSession := &general.LoginSession{}

	var queryProcessors []repository.QueryProcessor

	queryProcessors = append(queryProcessors, repository.OrderBy("start_time DESC"),
		repository.Filter("`credential_id` = ?", credentialID))

	if util.IsUUIDValid(loginSessionID) {
		queryProcessors = append(queryProcessors, repository.Filter("`id` = ?", loginSessionID))
	}

	if err := service.Repository.GetRecord(uow, loginSession, queryProcessors...); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Invalid login ID")
	}

	endTime := time.Now()
	loginSession.EndTime = &endTime
	loginSession.UpdatedBy = credentialID

	//update login session
	err := service.Repository.Update(uow, loginSession)
	// err := service.Repository.UpdateWithMap(uow, loginSession, map[interface{}]interface{}{
	// 	"EndTime": time.Now(),
	// 	"UpdatedBy": credentialID,
	// })
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Login session could not be ended")
	}
	uow.Commit()
	return nil
}

// DeleteLoginSession By Login session ID
func (service *LoginSessionService) DeleteLoginSession(loginSessionID uuid.UUID) error {
	//create bucket
	loginSession := general.LoginSession{}

	//uow for both get and delete operations
	uow := repository.NewUnitOfWork(service.DB, false)

	//get login session
	err := service.Repository.Get(uow, loginSessionID, &loginSession)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}
	//check parent tenant exists or not
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", loginSession.TenantID))
	if err != nil {
		return err
	}
	if !exists {
		return errors.NewValidationError("No tenant found")
	}
	//delete login session
	err = service.Repository.Delete(uow, &loginSession)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewValidationError("Login session could not be deleted")
	}
	uow.Commit()
	return nil
}

//doesForeignKeyExist for login session id and parent id validation
func (service *LoginSessionService) doesForeignKeyExist(loginSession *general.LoginSession,
	credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(loginSession.TenantID)
	if err != nil {
		return err
	}

	err = service.doesCredentialExist(loginSession.TenantID, credentialID)
	if err != nil {
		return err
	}

	err = service.doesCredentialExist(loginSession.TenantID, loginSession.CredentialID)
	if err != nil {
		return err
	}

	//check parent role id  is nil or not
	// if loginSession.RoleID == uuid.Nil { //parent role is nil
	// 	return errors.NewValidationError("Invalid role ID")
	// }
	// role := general.Role{}
	// role.ID = loginSession.RoleID
	// ok, err := repository.DoesRecordExistForTenant(service.DB, loginSession.TenantID, role)
	// if err != nil || !ok {
	// 	return errors.NewValidationError("No role found")
	// }
	// credentials := general.Credential{}
	// credentials.ID = loginSession.CredentialID
	// ok, err = repository.DoesRecordExistForTenant(service.DB, loginSession.TenantID, credentials)
	// if err != nil || !ok {
	// 	return errors.NewValidationError("No login found")
	// }

	//if update operation check if id is present or not
	// if loginSession.ID != uuid.Nil {
	// 	// id is nil
	// 	session := general.LoginSession{}
	// 	session.ID = loginSession.ID
	// 	ok, err = repository.DoesRecordExistForTenant(service.DB, loginSession.TenantID, session)
	// 	if err != nil || !ok {
	// 		uow.RollBack()
	// 		return errors.NewValidationError("Invalid login session ID")
	// 	}
	// }

	return nil
}

// returns error if there is no tenant record in table.
func (service *LoginSessionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *LoginSessionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
