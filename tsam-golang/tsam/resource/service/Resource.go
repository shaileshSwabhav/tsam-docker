package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	res "github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ResourceService Provide method to Update, Delete, Add, Get Method For Resource.
type ResourceService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewResourceService returns new instance of ResourceService.
func NewResourceService(db *gorm.DB, repository repository.Repository) *ResourceService {
	return &ResourceService{
		DB:         db,
		Repository: repository,
		association: []string{
			// "Book",
			"Technology",
		},
	}
}

// AddResource will add new resource to the table.
func (service *ResourceService) AddResource(resource *res.Resource, uows ...*repository.UnitOfWork) error {

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {

		// Check foreign key exist.
		err := service.doesForeignKeyExist(resource, resource.CreatedBy)
		if err != nil {
			return err
		}

		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// resourceBook := resource.Book
	// resource.Book = nil

	// if *resource.IsBook {
	// 	resourceBook.CreatedBy = resource.CreatedBy
	// 	resourceBook.TenantID = resource.TenantID
	// }

	// Add resource to resources table.
	err := service.Repository.Add(uow, resource)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return err
	}

	// resourceBook.ID = resource.ID
	// resourceBook.IgnoreHook = true

	// // Add resource to resources table.
	// err = service.Repository.Add(uow, resourceBook)
	// if err != nil {
	// 	if length == 0 {
	// 		uow.RollBack()
	// 	}
	// 	return err
	// }

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddResources will add multiple resource to the table.
func (service *ResourceService) AddResources(resources *[]res.Resource, tenantID,
	credentialID uuid.UUID) error {

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

	// Create new unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	for index := range *resources {

		// Check if unique resource is added.
		err = service.doesResourceExist(&(*resources)[index])
		if err != nil {
			return err
		}

		(*resources)[index].TenantID = tenantID
		(*resources)[index].CreatedBy = credentialID

		err := service.AddResource(&(*resources)[index], uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// UpdateResource will update the specified resource.
func (service *ResourceService) UpdateResource(resource *res.Resource) error {

	// Check foreign key exist.
	err := service.doesForeignKeyExist(resource, resource.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if resource exist
	err = service.doesResourceRecordExist(resource.TenantID, resource.ID)
	if err != nil {
		return err
	}

	// Create new unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	tempResource := res.Resource{}
	err = service.Repository.GetRecordForTenant(uow, resource.TenantID, &tempResource,
		repository.Filter("`id`=?", resource.ID), repository.Select("`created_by`"))
	if err != nil {
		return err
	}

	resource.CreatedBy = tempResource.CreatedBy

	// Add resource to resources table.
	err = service.Repository.Save(uow, resource)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteResource will delete the specified resource.
func (service *ResourceService) DeleteResource(resource *res.Resource) error {

	// Check foreign key exist.
	err := service.doesForeignKeyExist(resource, resource.DeletedBy)
	if err != nil {
		return err
	}

	// Check if resource exist
	err = service.doesResourceRecordExist(resource.TenantID, resource.ID)
	if err != nil {
		return err
	}

	// Create new unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.deleteResourceAssociation(uow, resource)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Add resource to resources table.
	err = service.Repository.UpdateWithMap(uow, res.Resource{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": resource.DeletedBy,
	}, repository.Filter("`id` = ?", resource.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// delete records from course_session_resource table not handled

	uow.Commit()
	return nil
}

// GetAllResources will return all the resources.
func (service *ResourceService) GetAllResources(resources *[]res.ResourceDTO, tenantID uuid.UUID,
	requestForm url.Values, limit, offset int, totalCount *int) error {

	// check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// extract credentialID from queryParams.
	credentialID, err := uuid.FromString(requestForm.Get("credentialID"))
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExists(tenantID, credentialID)
	if err != nil {
		return err
	}

	// Start transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, resources, "`resource_name`",
		service.addSearchQueries(requestForm), repository.PreloadAssociations(service.association),
		repository.PreloadWithCondition(map[string][]interface{}{
			"ResourceLike": {
				"`credential_id` = ? AND `is_liked` = '1'", credentialID,
			},
		}), repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *resources {
		// Get total count for resource count.
		err = service.Repository.GetCountForTenant(uow, tenantID, &res.Download{}, &(*resources)[index].TotalDownload,
			repository.Filter("`resource_id` = ?", (*resources)[index].ID), repository.GroupBy("`credential_id`"))
		if err != nil {
			uow.RollBack()
			return err
		}

		// Get total coutn for resource like.
		err = service.Repository.GetCount(uow, &res.Like{}, &(*resources)[index].TotalLike,
			repository.Filter("`resource_id` = ? AND is_liked = '1'", (*resources)[index].ID))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetResourceCount will return the count for the sprcified file-type
func (service *ResourceService) GetResourceCount(resourceCount *[]res.Count,
	tenantID uuid.UUID, requestForm url.Values) error {

	// check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Start transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	tempFileType := []general.CommonType{}

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &tempFileType, "`key`",
		repository.Filter("`type` = 'file_type'"))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range tempFileType {
		var tempResourceCount res.Count

		err = service.getFileTypeCount(uow, tempFileType[index].Value, &tempResourceCount,
			tenantID, requestForm)
		if err != nil {
			uow.RollBack()
			return err
		}
		*resourceCount = append(*resourceCount, tempResourceCount)
	}

	uow.Commit()

	return nil
}

// GetResourcesList will return all the resources for specified resource and file type.
func (service *ResourceService) GetResourcesList(resource *[]res.Resource, tenantID uuid.UUID,
	requestForm url.Values) error {

	// check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Start transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, resource, "`resource_name`",
		service.addSearchQueries(requestForm))
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

// getFileTypeCount will return count of resources based on the file-type
func (service *ResourceService) getFileTypeCount(uow *repository.UnitOfWork, fileType string, tempResourceCount *res.Count,
	tenantID uuid.UUID, requestForm url.Values) error {

	tempResourceCount.FileType = fileType

	err := service.Repository.GetCountForTenant(uow, tenantID, res.Resource{},
		&tempResourceCount.TotalCount, repository.Filter("`file_type` = ?", fileType),
		service.addSearchQueries(requestForm))
	if err != nil {
		uow.RollBack()
		return err
	}

	return nil
}

func (service *ResourceService) doesForeignKeyExist(resource *res.Resource, credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExists(resource.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExists(resource.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if unique resource is added.
	err = service.doesResourceExist(resource)
	if err != nil {
		return err
	}

	// Check if technology exist.
	if resource.TechnologyID != nil {
		err = service.doesTechnologyExists(resource.TenantID, *resource.TechnologyID)
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteResourceAssociation will remove the association of session with resource
func (service *ResourceService) deleteResourceAssociation(uow *repository.UnitOfWork, resource *res.Resource) error {

	tempCourseSession := []course.CourseSession{}
	err := service.Repository.GetAll(uow, &tempCourseSession,
		repository.Join("INNER JOIN course_sessions_resources ON course_sessions_resources.`course_session_id` = course_sessions.`id`"),
		repository.Join("INNER JOIN resources ON course_sessions_resources.`resource_id` = resources.`id`"),
		repository.Filter("course_sessions.`tenant_id`=?", resource.TenantID),
		repository.Filter("resources.`id`=? AND resources.`deleted_at` IS NULL AND course_sessions.`deleted_at` IS NULL", resource.ID),
		repository.PreloadAssociations([]string{"Resources"}))
	if err != nil {
		return err
	}

	// Remove association from map table.
	for index := range tempCourseSession {
		err = service.Repository.RemoveAssociations(uow, &tempCourseSession[index], "Resources", resource)
		if err != nil {
			return err
		}
	}
	return nil
}

// doesTenantExists validates tenantID.
func (service *ResourceService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExists validates credentialID.
func (service *ResourceService) doesCredentialExists(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Credential not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesResourceExist checks if resource already exist.
func (service *ResourceService) doesResourceExist(resource *res.Resource) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, resource.TenantID, res.Resource{},
		repository.Filter("`resource_type`=? AND `file_type`=? AND `resource_name`=? AND `id`!=?",
			resource.ResourceType, resource.FileType, resource.ResourceName, resource.ID))
	if err := util.HandleIfExistsError("Resource already exist", exists, err); err != nil {
		return err
	}
	return nil
}

// doesResourceNameExist checks if resource name already exsit.
// func (ser *ResourceService) doesResourceNameExist(tenantID uuid.UUID, resourceName string) error {
// 	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, resource.Resource{},
// 		repository.Filter("`resource_name`=? AND `id`!=?", resourceName))
// 	if err := util.HandleIfExistsError("Resource name already exist", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (service *ResourceService) doesResourceRecordExist(tenantID, resourceID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, res.Resource{},
		repository.Filter("`id` = ?", resourceID))
	if err := util.HandleError("Resource not found", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no technology record in table for the given tenant.
func (service *ResourceService) doesTechnologyExists(tenantID, technologyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Technology{},
		repository.Filter("`id`=?", technologyID))
	if err := util.HandleError("Technology not found", exists, err); err != nil {
		return err
	}
	return nil
}

// addSearchQueries will search queries to queryprocessor
func (service *ResourceService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["resourceName"]; ok {
		util.AddToSlice("`resource_name`", "LIKE ?", "AND", "%"+requestForm.Get("resourceName")+"%", &columnNames, &conditions, &operators, &values)
	}

	if resourceType, ok := requestForm["resourceType"]; ok {
		util.AddToSlice("`resource_type`", "= ?", "AND", resourceType, &columnNames, &conditions, &operators, &values)
	}

	if resourceSubType, ok := requestForm["resourceSubType"]; ok {
		util.AddToSlice("`resource_sub_type`", "= ?", "AND", resourceSubType, &columnNames, &conditions, &operators, &values)
	}

	if fileType, ok := requestForm["fileType"]; ok {
		util.AddToSlice("`file_type`", "= ?", "AND", fileType, &columnNames, &conditions, &operators, &values)
	}

	if technologyID, ok := requestForm["technologyID"]; ok {
		util.AddToSlice("`technology_id`", "= ?", "AND", technologyID, &columnNames, &conditions, &operators, &values)
	}

	if isBook, ok := requestForm["isBook"]; ok {
		util.AddToSlice("`is_book`", "= ?", "AND", isBook, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
