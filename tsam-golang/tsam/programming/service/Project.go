package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingProjectService Provide method to Update, Delete, Add, Get Method For programming-project.
type ProgrammingProjectService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewProgrammingProjectService returns new instance of ProgrammingProjectService.
func NewProgrammingProjectService(db *gorm.DB, repository repository.Repository) *ProgrammingProjectService {
	return &ProgrammingProjectService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Technologies", "Resources",
		},
	}
}

// AddProgrammingProject will add new programming-project in the table
func (service *ProgrammingProjectService) AddProgrammingProject(project *programming.ProgrammingProject) error {

	// Checks if all foreign keys are present.
	err := service.doesForeignKeyExist(project, project.CreatedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Assign programming project code
	project.Code, err = util.GenerateUniqueCode(uow.DB, project.ProjectName,
		"`code` = ?", &programming.ProgrammingProject{})
	if err != nil {
		return errors.NewHTTPError("Fail to generate Code", http.StatusInternalServerError)
	}

	err = service.Repository.Add(uow, project)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateProgrammingProject will update specified programming-project in the table
func (service *ProgrammingProjectService) UpdateProgrammingProject(project *programming.ProgrammingProject) error {

	// Checks if all foreign keys are present.
	err := service.doesForeignKeyExist(project, project.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if specified programming-project exist.
	err = service.doesProgrammingProjectExist(project.TenantID, project.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempProgrammingProject := programming.ProgrammingProject{}

	// Get createdBy and code for specified programming-project.
	err = service.Repository.GetRecordForTenant(uow, project.TenantID, &tempProgrammingProject,
		repository.Filter("`id` = ?", project.ID), repository.Select([]string{"`created_by`, `code`"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	project.Code = tempProgrammingProject.Code
	project.CreatedBy = tempProgrammingProject.CreatedBy

	// replace m2m table assocations.
	err = service.updateProgrammingProjectAssociation(uow, project)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Update programming-project.
	err = service.Repository.Save(uow, project)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteProgrammingProject will soft delete specified programming-project from the table
func (service *ProgrammingProjectService) DeleteProgrammingProject(project *programming.ProgrammingProject) error {

	// Checks if all foreign keys are present.
	err := service.doesForeignKeyExist(project, project.DeletedBy)
	if err != nil {
		return err
	}

	// Check if specified programming-project exist.
	err = service.doesProgrammingProjectExist(project.TenantID, project.ID)
	if err != nil {
		return err
	}

	// Check if project is assigned to btahc.
	err = service.doesProgrammingPtojectAssigned(project.TenantID, project.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Soft Delete programming-project.
	err = service.Repository.UpdateWithMap(uow, programming.ProgrammingProject{}, map[string]interface{}{
		"DeletedBy": project.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", project.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllProgrammingProject will get all the programming-project from the table.
func (service *ProgrammingProjectService) GetAllProgrammingProject(projects *[]programming.ProgrammingProjectDTO,
	tenantID uuid.UUID, totalCount *int, parser *web.Parser) error {

	// Check if tenant exists.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()

	var queryProcessors []repository.QueryProcessor
	queryProcessors = append(queryProcessors, service.addSearchQueriesParams(parser.Form)...)
	queryProcessors = append(queryProcessors, repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount))

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, projects, "`project_name`",
		queryProcessors...)
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

// updateProgrammingProjectAssociation will update technologies and resources m2m table entries.
func (ser *ProgrammingProjectService) updateProgrammingProjectAssociation(uow *repository.UnitOfWork, project *programming.ProgrammingProject) error {

	// Replace technologies associations.
	err := ser.Repository.ReplaceAssociations(uow, project, "Technologies", project.Technologies)
	if err != nil {
		return err
	}

	// Replace resources associations.
	err = ser.Repository.ReplaceAssociations(uow, project, "Resources", project.Resources)
	if err != nil {
		return err
	}

	project.Technologies = nil
	project.Resources = nil

	return nil
}

// doesForeignKeyExist will check all the foreign keys of programming-project.
func (ser *ProgrammingProjectService) doesForeignKeyExist(programmingProject *programming.ProgrammingProject, credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := ser.doesTenantExists(programmingProject.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = ser.doesCredentialExist(programmingProject.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if project name already exist
	err = ser.doesProjectNameExist(programmingProject.TenantID, programmingProject.ID, programmingProject.ProjectName)
	if err != nil {
		return err
	}

	return nil
}

// doesTenantExists validates tenantID.
func (ser *ProgrammingProjectService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(ser.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates credentialID.
func (ser *ProgrammingProjectService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Credential not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesProgrammingProjectExist validates programmingProjectID.
func (ser *ProgrammingProjectService) doesProgrammingProjectExist(tenantID, projectID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, programming.ProgrammingProject{},
		repository.Filter("`id` = ?", projectID))
	if err := util.HandleError("Programming project not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesProjectNameExist validates programming-project name.
func (ser *ProgrammingProjectService) doesProjectNameExist(tenantID, projectID uuid.UUID, projectName string) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, programming.ProgrammingProject{},
		repository.Filter("`project_name` = ? AND `is_active` = true AND `id` != ?", projectName, projectID))
	if err := util.HandleIfExistsError("Project name already exist", exists, err); err != nil {
		return err
	}
	return nil
}

// addSearchQueriesParams adds all search queries if any when getAll is called
func (ser *ProgrammingProjectService) addSearchQueriesParams(requestForm url.Values) []repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	var queryProcessors []repository.QueryProcessor

	if _, ok := requestForm["projectName"]; ok {
		util.AddToSlice("`project_name`", "LIKE ?", "AND", "%"+requestForm.Get("projectName")+"%", &columnNames, &conditions, &operators, &values)
	}

	if complexityLevel, ok := requestForm["complexityLevel"]; ok {
		util.AddToSlice("`complexity_level`", "= ?", "AND", complexityLevel, &columnNames, &conditions, &operators, &values)
	}

	if resourceType, ok := requestForm["resourceType"]; ok {
		util.AddToSlice("`resource_type`", "= ?", "AND", resourceType, &columnNames, &conditions, &operators, &values)
	}

	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("`is_active`", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	//if technologies is present then join programming_projects and technolgies table
	if technologies, ok := requestForm["technologies"]; ok {
		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN programming_projects_technologies ON programming_projects.`id` = programming_projects_technologies.`programming_project_id`"))
		if len(technologies) > 0 {
			util.AddToSlice("programming_projects_technologies.`technology_id`", "IN(?)", "AND", technologies, &columnNames, &conditions, &operators, &values)
		}
	}

	//if resources is present then join programming_projects and resources table
	if resources, ok := requestForm["resources"]; ok {
		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN programming_projects_resources ON programming_projects.`id` = programming_projects_resources.`programming_project_id`"))
		if len(resources) > 0 {
			util.AddToSlice("programming_projects_resources.`resource_id`", "IN(?)", "AND", resources, &columnNames, &conditions, &operators, &values)
		}
	}

	queryProcessors = append(queryProcessors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("programming_projects.`id`"))

	return queryProcessors
}

// doesProgrammingPtojectAssigned validates project is assigned to batch
func (ser *ProgrammingProjectService) doesProgrammingPtojectAssigned(tenantID, projectID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, batch.Project{},
		repository.Filter("`programming_project_id` = ? AND deleted_at IS NULL", projectID))

	if err := util.HandleIfExistsError("Project is Assigned for batch, So can't be Deleted", exists, err); err != nil {
		return err
	}
	return nil
}
