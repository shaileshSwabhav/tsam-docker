package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// DepartmentService provides methods to do different CRUD operations on department table.
type DepartmentService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewDepartmentService returns a new instance Of DepartmentService.
func NewDepartmentService(db *gorm.DB, repository repository.Repository) *DepartmentService {
	return &DepartmentService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Role",
		},
	}
}

// AddDepartment will add new department to the table.
func (service *DepartmentService) AddDepartment(department *general.Department) error {

	// Check if all foreign keys exist.
	err := service.doForeignKeysExist(department, department.CreatedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, department)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateDepartment will update the specified department record in the table.
func (service *DepartmentService) UpdateDepartment(department *general.Department) error {

	// Checks if all foreign key exist.
	err := service.doForeignKeysExist(department, department.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if department record exist.
	err = service.doesDepartmentExist(department.TenantID, department.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, department)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// DeleteDepartment will delete the specified department record from the table.
func (service *DepartmentService) DeleteDepartment(department *general.Department) error {
	credentialID := department.DeletedBy

	// Check if tenant exists.
	err := service.doesTenantExist(department.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(department.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if department record exist.
	err = service.doesDepartmentExist(department.TenantID, department.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update department for updating deleted_by and deleted_at fields of department
	if err := service.Repository.UpdateWithMap(uow, department, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", department.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Department could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetAllDepartments will return all the records from department table.
func (service *DepartmentService) GetAllDepartments(departments *[]general.DepartmentDTO, tenantID uuid.UUID,
	parser *web.Parser, totalCount *int) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, departments, "`name`",
		service.addSearchQueries(parser.Form), repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetDepartmentList will return all the records by role from department table (no limit, offset).
func (service *DepartmentService) GetDepartmentList(departments *[]general.DepartmentDTO, tenantID uuid.UUID, form url.Values) error {

	// Get query params for role name and login id.
	roleNames := form["roleNames"]

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Query processors gettomg withut role name filter.
	var queryProcessors []repository.QueryProcessor
	queryProcessors = append(queryProcessors,
		repository.PreloadAssociations(service.associations),
		repository.Filter("departments.`tenant_id`=?", tenantID))

	// If role names filter exists in params then add filter to query processors.
	if len(roleNames) != 0 {
		queryProcessors = append(queryProcessors,
			repository.Join("INNER JOIN roles ON roles.`id` = departments.`role_id`"),
			repository.Filter("roles.`tenant_id`=? AND roles.`deleted_at` IS NULL", tenantID),
			repository.Filter("`role_name` IN(?)", roleNames))
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrder(uow, departments, "departments.`name`", queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetDepartment will return specified record from department table.
func (service *DepartmentService) GetDepartment(department *general.DepartmentDTO, tenantID, departmentID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if department record exist.
	err = service.doesDepartmentExist(tenantID, departmentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetRecordForTenant(uow, tenantID, department,
		repository.Filter("`id`=?", departmentID), repository.PreloadAssociations(service.associations))
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

// doForeignKeysExist checks if all foregin key exist and if not returns error.
func (service *DepartmentService) doForeignKeysExist(department *general.Department, credentialID uuid.UUID) error {

	// Check if tenant exists.
	err := service.doesTenantExist(department.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(department.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if department with same name exist
	err = service.doesDepartmentNameExist(department.TenantID, department.ID, department.Name)
	if err != nil {
		return err
	}

	// check if role exist
	err = service.doesRoleExist(department.TenantID, department.RoleID)
	if err != nil {
		return err
	}

	return nil
}

// addSearchQueries adds search queries.
func (service *DepartmentService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["name"]; ok {
		util.AddToSlice("`name`", "LIKE ?", "AND", "%"+requestForm.Get("name")+"%", &columnNames, &conditions, &operators, &values)
	}
	if roleID, ok := requestForm["roleID"]; ok {
		util.AddToSlice("`role_id`", "=?", "AND", roleID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *DepartmentService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *DepartmentService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDepartmentNameExist returns error if there is no department record with same name exist in table for the given tenant.
func (service *DepartmentService) doesDepartmentNameExist(tenantID, departmentID uuid.UUID, departmentName string) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Department{},
		repository.Filter("`name`=? AND `id`!=?", departmentName, departmentID))
	if err := util.HandleIfExistsError("Department name already exist", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDepartmentExist returns error if there is no department record in table for the given tenant.
func (service *DepartmentService) doesDepartmentExist(tenantID, departmentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Department{},
		repository.Filter("`id` = ?", departmentID))
	if err := util.HandleError("Invalid department ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesRoleExist returns error if there is no role record in table for the given tenant.
func (service *DepartmentService) doesRoleExist(tenantID, roleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Role{},
		repository.Filter("`id` = ?", roleID))
	if err := util.HandleError("Invalid role ID", exists, err); err != nil {
		return err
	}
	return nil
}
