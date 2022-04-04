package service

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// SwabhavProjectService Provide method to Update, Delete, Add, Get Method For Project.
type SwabhavProjectService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewSwabhavProjectService returns new instance of ProjectService.
func NewSwabhavProjectService(db *gorm.DB, repository repository.Repository) *SwabhavProjectService {
	return &SwabhavProjectService{
		DB:         db,
		Repository: repository,
		association: []string{
			"Project",
		},
	}
}

// AddProject will add new project to the table.
func (service *SwabhavProjectService) AddProject(project *general.SwabhavProject) error {

	// Checks if all foreign keys exist.
	err := service.doForeignKeysExist(project, project.CreatedBy)
	if err != nil {
		return err
	}

	if project.SubProjects != nil {
		for index := range project.SubProjects {
			(project.SubProjects)[index].TenantID = project.TenantID
			(project.SubProjects)[index].CreatedBy = project.CreatedBy
		}
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

// UpdateProject will update the specified project record in the table.
func (service *SwabhavProjectService) UpdateProject(project *general.SwabhavProject) error {

	// Checks if all foreign key exist.
	err := service.doForeignKeysExist(project, project.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if project record exist.
	err = service.doesProjectExist(project.TenantID, project.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update sub-projects.
	err = service.updateSubProject(uow, project)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.Update(uow, project)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// DeleteProject will delete the specified project record from the table.
func (service *SwabhavProjectService) DeleteProject(project *general.SwabhavProject) error {

	// Checks if all foreign keys exist.
	err := service.doForeignKeysExist(project, project.DeletedBy)
	if err != nil {
		return err
	}

	// Check if project record exist.
	err = service.doesProjectExist(project.TenantID, project.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.deleteSubProject(uow, project.ID, project.DeletedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.UpdateWithMap(uow, &general.SwabhavProject{}, map[interface{}]interface{}{
		"DeletedBy": project.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id`=?", project.ID))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetAllProject will return all the records from project table.
func (service *SwabhavProjectService) GetAllProject(projects *[]general.SwabhavProject, tenantID uuid.UUID,
	parser *web.Parser, totalCount *int) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}
	// if parser.Form.Get("isPaginate") == "" {
	// 	limit, offset := parser.ParseLimitAndOffset()
	// 	queryProcessor = repository.Paginate(limit, offset, totalCount)
	//   }
	limit, offset := parser.ParseLimitAndOffset()
	fmt.Println("limit ", limit, "offset ", offset)
	uow := repository.NewUnitOfWork(service.DB, true)
	// repository.PreloadAssociations(service.association)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, projects, "`name`",
		service.addSearchQueries(parser.Form), repository.Filter("`project_id` IS NULL"),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *projects {
		err = service.getSubProjects(&((*projects)[index].SubProjects), tenantID, (*projects)[index].ID, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetProjects will return all the records from project table.
func (service *SwabhavProjectService) GetProjects(projects *[]general.SwabhavProject, tenantID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// repository.PreloadAssociations(service.association)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, projects, "`name`",
		repository.Filter("`project_id` IS NULL"))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *projects {
		err = service.getSubProjects(&((*projects)[index].SubProjects), tenantID, (*projects)[index].ID, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetProject will return specified record from project table.
func (service *SwabhavProjectService) GetProject(project *general.SwabhavProject, tenantID, projectID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if project record exist.
	err = service.doesProjectExist(tenantID, projectID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	// repository.PreloadAssociations(service.association)
	err = service.Repository.GetRecordForTenant(uow, tenantID, project,
		repository.Filter("`id`=?", projectID), repository.Filter("`project_id` IS NULL"))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.getSubProjects(&(project.SubProjects), tenantID, project.ID, uow)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// Niranjan

// GetProjectList will return list of projects.
func (service *SwabhavProjectService) GetProjectList(tenantID uuid.UUID, projects *[]general.SwabhavProject) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	// repository.PreloadAssociations(service.association)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, projects, "`name`",
		repository.Filter("`swabhav_projects`.`project_id` IS NULL"),
		repository.Select([]string{"`swabhav_projects`.`id`", "`swabhav_projects`.`name`"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetSubProjectList will return list of subproject with respect to Project.
func (service *SwabhavProjectService) GetSubProjectList(tenantID, projectID uuid.UUID,
	subProjects *[]general.SwabhavProject) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	// repository.PreloadAssociations(service.association)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, subProjects, "`name`",
		repository.Filter("`swabhav_projects`.`project_id` = ? ", projectID),
		repository.Select([]string{"`swabhav_projects`.`id`", "`swabhav_projects`.`name`"}))
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

// getSubProjects gets called recursively to get all sub projects.
func (service *SwabhavProjectService) getSubProjects(projects *[]general.SwabhavProject, tenantID, projectID uuid.UUID,
	uow *repository.UnitOfWork) error {

	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &general.SwabhavProject{},
		repository.Filter("`project_id`=?", projectID))
	if err != nil {
		return err
	}
	if exist {
		err = service.Repository.GetAllInOrderForTenant(uow, tenantID, projects, "`name`",
			repository.Filter("`project_id`=?", projectID))
		if err != nil {
			return err
		}

		// Rescursive call.
		for index := range *projects {
			err = service.getSubProjects(&((*projects)[index].SubProjects), tenantID, (*projects)[index].ID, uow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// updateSubProject will update all the sub-projects.
func (service *SwabhavProjectService) updateSubProject(uow *repository.UnitOfWork, project *general.SwabhavProject) error {

	tempSubProjects := []general.SwabhavProject{}
	projectMap := make(map[uuid.UUID]uint)

	err := service.Repository.GetAllForTenant(uow, project.TenantID, &tempSubProjects,
		repository.Filter("`project_id`=?", project.ID))
	if err != nil {
		return err
	}

	for _, tempSubProject := range tempSubProjects {
		projectMap[tempSubProject.ID]++
	}

	for _, subProject := range project.SubProjects {

		// If it is an existing sub-project.
		if util.IsUUIDValid(subProject.ID) {
			projectMap[subProject.ID]++
		} else {
			// Add sub-project for non-existing sub-project.
			subProject.TenantID = project.TenantID
			subProject.CreatedBy = project.UpdatedBy
			subProject.ProjectID = &project.ID

			err = service.Repository.Add(uow, &subProject)
			if err != nil {
				return err
			}
		}

		// Update if existing sub-project.
		if projectMap[subProject.ID] > 1 {
			subProject.UpdatedBy = project.UpdatedBy

			err := service.Repository.Update(uow, &subProject)
			if err != nil {
				return err
			}

			projectMap[subProject.ID] = 0
		}
	}

	// Delete sub-projects.
	for _, tempSubProject := range tempSubProjects {
		if projectMap[tempSubProject.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, &general.SwabhavProject{}, map[interface{}]interface{}{
				"DeletedAt": time.Now(),
				"DeletedBy": project.UpdatedBy,
			}, repository.Filter("`id`=?", tempSubProject.ID))
			if err != nil {
				return err
			}
			projectMap[tempSubProject.ID] = 0
		}
	}

	project.SubProjects = nil

	return nil
}

// deleteSubProject will delete all the sub-projects if parent project is deleted.
func (service *SwabhavProjectService) deleteSubProject(uow *repository.UnitOfWork, projectID, credentialID uuid.UUID) error {

	err := service.Repository.UpdateWithMap(uow, &general.SwabhavProject{}, map[interface{}]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": credentialID,
	}, repository.Filter("`project_id`=?", projectID))
	if err != nil {
		return err
	}
	return nil
}

// doForeignKeysExist checks if all foregin keys exist and if not returns error.
func (service *SwabhavProjectService) doForeignKeysExist(project *general.SwabhavProject, credentialID uuid.UUID) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(project.TenantID); err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(project.TenantID, credentialID); err != nil {
		return err
	}

	// check if project with same name exist
	if err := service.doesProjectNameExist(project.TenantID, project.ID, project.Name); err != nil {
		return err
	}

	return nil
}

// addSearchQueries adds search criteria.
func (service *SwabhavProjectService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["projectName"]; ok {
		util.AddToSlice("`name`", "LIKE ?", "AND", "%"+requestForm.Get("projectName")+"%", &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *SwabhavProjectService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *SwabhavProjectService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesProjectNameExist returns error if there is no project record with same name exist in table for the given tenant.
func (service *SwabhavProjectService) doesProjectNameExist(tenantID, projectID uuid.UUID, projectName string) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.SwabhavProject{},
		repository.Filter("`name`=? AND `id`!=? AND `project_id` IS NULL", projectName, projectID))
	if err := util.HandleIfExistsError("Project name already exists", exists, err); err != nil {
		return err
	}
	return nil
}

// doesProjectExist returns error if there is no project record in table for the given tenant.
func (service *SwabhavProjectService) doesProjectExist(tenantID, projectID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.SwabhavProject{},
		repository.Filter("`id`=?", projectID))
	if err := util.HandleError("Invalid project ID", exists, err); err != nil {
		return err
	}
	return nil
}
