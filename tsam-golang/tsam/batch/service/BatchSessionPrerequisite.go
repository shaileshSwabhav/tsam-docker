package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

type PrerequisiteService struct {
	db   *gorm.DB
	repo repository.Repository
}

func NewPrerequisiteService(db *gorm.DB, repo repository.Repository) *PrerequisiteService {
	return &PrerequisiteService{
		db:   db,
		repo: repo,
	}
}

// AddBatchSessionPrerequisite will add pre-requisite for batch session.
func (service *PrerequisiteService) AddBatchSessionPrerequisite(prerequisite *batch.BatchSessionPrerequisite) error {

	// check if foreign exist.
	err := service.doesForeignKeyExist(prerequisite, prerequisite.CreatedBy)
	if err != nil {
		return err
	}

	err = service.doesPrerequisiteExistForSession(prerequisite.TenantID, prerequisite.BatchID, prerequisite.BatchSessionID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	err = service.repo.Add(uow, prerequisite)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateBatchSessionPrerequisite will update pre-requisite for batch session.
func (service *PrerequisiteService) UpdateBatchSessionPrerequisite(prerequisite *batch.BatchSessionPrerequisite) error {

	// check if foreign exist.
	err := service.doesForeignKeyExist(prerequisite, prerequisite.UpdatedBy)
	if err != nil {
		return err
	}

	// check if pre-requisite exist.
	err = service.doesPrerequisiteExist(prerequisite.TenantID, prerequisite.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	err = service.repo.Update(uow, prerequisite)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteBatchSessionPrerequisite will delete pre-requisite for batch session.
func (service *PrerequisiteService) DeleteBatchSessionPrerequisite(prerequisite *batch.BatchSessionPrerequisite) error {

	// check if foreign exist.
	err := service.doesForeignKeyExist(prerequisite, prerequisite.DeletedBy)
	if err != nil {
		return err
	}

	// check if pre-requisite exist.
	err = service.doesPrerequisiteExist(prerequisite.TenantID, prerequisite.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	err = service.repo.UpdateWithMap(uow, batch.BatchSessionPrerequisite{}, map[string]interface{}{
		"DeletedBy": prerequisite.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", prerequisite.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBatchSessionPrerequisite will fetch pre-requisite for specified batch.
func (service *PrerequisiteService) GetBatchSessionPrerequisite(prerequisite *[]batch.BatchSessionPrerequisiteDTO,
	tenantID, batchID uuid.UUID, totalCount *int, parser *web.Parser) error {

	// check tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()

	uow := repository.NewUnitOfWork(service.db, true)

	err = service.repo.GetAllInOrderForTenant(uow, tenantID, prerequisite, "`created_at`",
		repository.Filter("`batch_id` = ?", batchID), service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// addSearchQueries adds all search queries if any when get is called.
func (service *PrerequisiteService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if batchSessionID, ok := requestForm["batchSessionID"]; ok {
		util.AddToSlice("`batch_session_id`", "= ?", "AND", batchSessionID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesForeignKeyExist checks if all foreign keys are valid.
func (service *PrerequisiteService) doesForeignKeyExist(prerequisite *batch.BatchSessionPrerequisite, credentialID uuid.UUID) error {

	// check tenant exist
	err := service.doesTenantExists(prerequisite.TenantID)
	if err != nil {
		return err
	}

	// check credential exist
	err = service.doesCredentialExist(prerequisite.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if batch session exist.
	err = service.doesBatchSessionExist(prerequisite.TenantID, prerequisite.BatchSessionID)
	if err != nil {
		return err
	}

	return nil
}

// doesTenantExists validates tenantID.
func (service *PrerequisiteService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.db, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *PrerequisiteService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

func (service *PrerequisiteService) doesBatchSessionExist(tenantID, batchSessionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.Session{},
		repository.Filter("`id` = ?", batchSessionID))
	if err := util.HandleError("Batch Session not found", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *PrerequisiteService) doesPrerequisiteExist(tenantID, preqrequisiteID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.BatchSessionPrerequisite{},
		repository.Filter("`id` = ?", preqrequisiteID))
	if err := util.HandleError("Prerequisite not found", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *PrerequisiteService) doesPrerequisiteExistForSession(tenantID, batchID, sessionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.BatchSessionPrerequisite{},
		repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, sessionID))
	if err := util.HandleIfExistsError("Prerequisite already exist", exists, err); err != nil {
		return err
	}
	return nil
}
