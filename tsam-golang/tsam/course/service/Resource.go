package service

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	res "github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ModuleResourceService Provide method to Update, Delete, Add, Get Method For Resource.
type ModuleResourceService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewModuleResourceService returns new instance of CourseSessionResourceService.
func NewModuleResourceService(db *gorm.DB, repository repository.Repository) *ModuleResourceService {
	return &ModuleResourceService{
		DB:         db,
		Repository: repository,
	}
}

// AddResource will add the resource for the specified session to DB
func (service *ModuleResourceService) AddResource(resource *course.ModuleResource,
	tenantID, credentialID uuid.UUID) error {

	// Check foreign key exist.
	err := service.doesForeignKeyExist(resource, tenantID, credentialID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add resource to resources table.
	err = service.Repository.Add(uow, &resource)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// AddResources will add the resource for the specified session to DB
// func (ser *ModuleResourceService) AddResources(resources *[]res.Resource, tenantID,
// 	courseSessionID, credentialID uuid.UUID) error {

// 	// Start transaction.
// 	uow := repository.NewUnitOfWork(ser.DB, false)

// 	// for _, resource := range *resources {
// 	// 	resource.TenantID = tenantID
// 	// 	resource.CreatedBy = credentialID

// 	// 	err := ser.AddResource(&resource, courseSessionID, uow)
// 	// 	if err != nil {
// 	// 		uow.RollBack()
// 	// 		return err
// 	// 	}
// 	// }

// 	uow.Commit()
// 	return nil
// }

// DeleteResources will delete the specified resource entry from the map table
func (service *ModuleResourceService) DeleteResources(resource *course.ModuleResource, tenantID, credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExists(tenantID, credentialID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempModule := course.Module{}
	tempModule.ID = resource.ModuleID

	tempResource := res.Resource{}
	tempResource.ID = resource.ResourceID

	// Remove association from map table.
	err = service.Repository.RemoveAssociations(uow, tempModule, "Resources", tempResource)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return err
}

// GetModuleResources will return all the resources for the specified module.
func (ser *ModuleResourceService) GetModuleResources(resource *[]res.Resource, tenantID, moduleID uuid.UUID) error {

	// check if tenant exist.
	err := ser.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// check if session exist.
	err = ser.doesModuleExist(tenantID, moduleID)
	if err != nil {
		return err
	}

	// Start transaction.
	uow := repository.NewUnitOfWork(ser.DB, true)

	err = ser.Repository.GetAll(uow, resource,
		repository.Join("INNER JOIN modules_resources ON modules_resources.`resource_id` = resources.`id`"),
		repository.Join("INNER JOIN modules ON modules_resources.`module_id` = modules.`id`"),
		repository.Filter("modules.`id`=? AND modules.`deleted_at` IS NULL AND resources.`tenant_id` = ?", moduleID, tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// // GetAllResources will return all the resources.
// func (ser *ModuleResourceService) GetAllResources(resource *[]res.Resource, tenantID uuid.UUID) error {

// 	// check if tenant exist.
// 	err := ser.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// Start transaction.
// 	uow := repository.NewUnitOfWork(ser.DB, true)

// 	err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, resource, "`resource_name`")
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

func (service *ModuleResourceService) doesForeignKeyExist(resource *course.ModuleResource,
	tenantID, credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExists(tenantID, credentialID)
	if err != nil {
		return err
	}

	// check if module exist.
	err = service.doesModuleExist(tenantID, resource.ModuleID)
	if err != nil {
		return err
	}

	// check if resource exist.
	err = service.doesResourceExist(tenantID, resource.ResourceID)
	if err != nil {
		return err
	}

	// check if modules_resources exist.
	err = service.doesResourceExistForModule(resource.ModuleID, resource.ResourceID)
	if err != nil {
		return err
	}

	return nil
}

// doesModuleExist validates module id.
func (service *ModuleResourceService) doesModuleExist(tenantID, moduleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Module{},
		repository.Filter("`id` = ?", moduleID))
	if err := util.HandleError("Module not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesResourceExist validates resource id.
func (service *ModuleResourceService) doesResourceExist(tenantID, resourceID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, res.Resource{},
		repository.Filter("`id` = ?", resourceID))
	if err := util.HandleError("Resource not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExists validates tenantID.
func (service *ModuleResourceService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExists validates credentialID.
func (service *ModuleResourceService) doesCredentialExists(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Credential not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesResourceExistForModule checks if map entry exist for module and resource.
func (service *ModuleResourceService) doesResourceExistForModule(moduleID, resourceID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, course.ModuleResource{},
		repository.Filter("`module_id` = ? AND `resource_id`=?", moduleID, resourceID))
	if err := util.HandleIfExistsError("Resource for current module already exist", exists, err); err != nil {
		return err
	}
	return nil
}
