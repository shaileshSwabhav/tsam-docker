package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

type ProjectService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

func NewProjectService(db *gorm.DB, repo repository.Repository) *ProjectService {
	return &ProjectService{
		DB:         db,
		Repository: repo,
	}
}

// AddProject will add programming project to batch.
func (service *ProjectService) AddProject(project *batch.Project) error {

	// check if tenant exist.
	err := service.doesForeignKeyExist(project, project.CreatedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, project)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateProject will update programming project to batch.
func (service *ProjectService) UpdateProject(project *batch.Project) error {

	// check if tenant exist.
	err := service.doesForeignKeyExist(project, project.UpdatedBy)
	if err != nil {
		return err
	}

	err = service.doesProjectExist(project.TenantID, project.ID)
	if err != nil {
		return err
	}

	if project.DueDate != nil {
		now := time.Now()
		project.AssignedDate = &now
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, project)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteProject will delete programming project to batch.
func (service *ProjectService) DeleteProject(batchProject *batch.Project) error {

	// check tenant exist
	err := service.doesTenantExists(batchProject.TenantID)
	if err != nil {
		return err
	}

	// check credential exist
	err = service.doesCredentialExist(batchProject.TenantID, batchProject.DeletedBy)
	if err != nil {
		return err
	}

	// check batch exist
	err = service.doesBatchExists(batchProject.TenantID, batchProject.BatchID)
	if err != nil {
		return err
	}
	err = service.doesProjectExist(batchProject.TenantID, batchProject.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, batch.Project{}, map[string]interface{}{
		"DeletedBy": batchProject.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", batchProject.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetProjects will get all the batch projects.
func (service *ProjectService) GetProjects(batchProjects *[]batch.ProjectDTO,
	tenantID, batchID uuid.UUID, totalCount *int, parser *web.Parser) error {

	// check tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, batchProjects, "`created_at`",
		repository.Filter("`batch_id` = ?", batchID), service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{"ProgrammingProject", "ProgrammingProject.Technologies", "ProgrammingProject.Resources"}), repository.Paginate(limit, offset, totalCount))
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

// doesForeignKeyExist checks if all foreign keys are valid.
func (service *ProjectService) doesForeignKeyExist(batchProject *batch.Project, credentialID uuid.UUID) error {

	// check tenant exist
	err := service.doesTenantExists(batchProject.TenantID)
	if err != nil {
		return err
	}

	// check credential exist
	err = service.doesCredentialExist(batchProject.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check batch exist
	err = service.doesBatchExists(batchProject.TenantID, batchProject.BatchID)
	if err != nil {
		return err
	}

	// check programming-project exist.
	err = service.doesProgrammingProjectExist(batchProject.TenantID, batchProject.ProgrammingProjectID)
	if err != nil {
		return err
	}

	return nil
}

// addSearchQueries adds all search queries if any when get is called.
func (service *ProjectService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	// if batchID, ok := requestForm["batchID"]; ok {
	// 	util.AddToSlice("`batch_id`", "= ?", "AND", batchID, &columnNames, &conditions, &operators, &values)
	// }

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesTenantExists validates tenantID.
func (service *ProjectService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *ProjectService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBatchExists validates batchID.
func (service *ProjectService) doesBatchExists(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Batch not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesProgrammingProjectExist validates projectID.
func (service *ProjectService) doesProgrammingProjectExist(tenantID, projectID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingProject{},
		repository.Filter("`id` = ?", projectID))
	if err := util.HandleError("Project not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesProjectExist validates batchProjectID.
func (service *ProjectService) doesProjectExist(tenantID, batchProjectID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Project{},
		repository.Filter("`id` = ?", batchProjectID))
	if err := util.HandleError("Batch Project not found", exists, err); err != nil {
		return err
	}
	return nil
}
