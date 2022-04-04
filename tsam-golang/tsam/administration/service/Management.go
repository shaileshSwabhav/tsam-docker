package service

import (
	"net/http"
	"net/url"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ManagementService provides methods to Add, Get operation Login
type ManagementService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewManagementService returns new instance of LoginService
func NewManagementService(db *gorm.DB, repository repository.Repository) *ManagementService {
	return &ManagementService{
		DB:         db,
		Repository: repository,
	}
}

// GetAllEmployees gets all employee details using the credential table
func (service *ManagementService) GetAllEmployees(tenantID uuid.UUID, limit, offset int, totalCount *int, form url.Values,
	allEmployees *[]admin.Employee) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrder(uow, allEmployees, "`credentials`.`first_name`",
		repository.Join("INNER JOIN `roles` ON `credentials`.`role_id` = `roles`.`id` AND "+
			"`credentials`.`tenant_id` = `roles`.`tenant_id`"),
		repository.Filter("`credentials`.`tenant_id` = ? AND `roles`.`deleted_at` IS NULL AND "+
			"`roles`.`is_employee` = ?", tenantID, true),
		service.addSearchQueries(form),
		repository.PreloadAssociations([]string{"Role", "Supervisors", "Supervisors.Role"}),
		repository.Paginate(limit, offset, totalCount))

	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetAllEmployeeList gets a list of names of all employees
func (service *ManagementService) GetAllEmployeeList(tenantID uuid.UUID, employeeList *[]list.Credential) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrder(uow, employeeList, "`credentials`.`first_name`",
		repository.Join("INNER JOIN `roles` ON `credentials`.`role_id` = `roles`.`id` AND "+
			"`credentials`.`tenant_id` = `roles`.`tenant_id`"),
		repository.Filter("`credentials`.`tenant_id` = ? AND `roles`.`deleted_at` IS NULL AND "+
			"`roles`.`is_employee` = ?", tenantID, true),
		repository.PreloadAssociations([]string{"Role"}))

	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetDirectReports gets a list of direct reports.
func (service *ManagementService) GetDirectReports(tenantID, supervisorID uuid.UUID,
	directReports *[]list.Credential) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrder(uow, directReports, "`credentials`.`first_name`",
		repository.Join("INNER JOIN `roles` ON `credentials`.`role_id` = `roles`.`id` AND "+
			"`credentials`.`tenant_id` = `roles`.`tenant_id`"),
		repository.Join("LEFT JOIN `employee_supervisors` ON "+
			"`credentials`.`id` = `employee_supervisors`.`employee_credential_id`"),
		repository.Filter("`credentials`.`tenant_id` = ? AND `roles`.`deleted_at` IS NULL AND "+
			"`roles`.`is_employee` = ? AND `employee_supervisors`.`supervisor_credential_id` = ?",
			tenantID, true, supervisorID),
		repository.PreloadAssociations([]string{"Role"}))

	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// AddSupervisor adds a supervisor to employee.
func (service *ManagementService) AddSupervisor(tenantID uuid.UUID, supervisor *admin.EmployeeSupervisor) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check for employee credential.
	err = service.doesCredentialExist(tenantID, supervisor.EmployeeCredentialID)
	if err != nil {
		return err
	}

	// Check for supervisor credential.
	err = service.doesCredentialExist(tenantID, supervisor.SupervisorCredentialID)
	if err != nil {
		return err
	}

	// validate unique fields have unique values.
	err = service.validateFieldUniqueness(supervisor)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.Add(uow, supervisor)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteSupervisor deletes supervisor from employee.
func (service *ManagementService) DeleteSupervisor(supervisor *admin.EmployeeSupervisor) error {
	uow := repository.NewUnitOfWork(service.DB, false)

	// Check if supervisor exists.
	exists, err := repository.DoesRecordExist(service.DB, supervisor,
		repository.Filter("`employee_credential_id` = ? AND `supervisor_credential_id` = ?",
			supervisor.EmployeeCredentialID, supervisor.SupervisorCredentialID))
	if err := util.HandleError("Supervisor not found", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	err = service.Repository.NewDelete(uow, supervisor,
		repository.Filter("`employee_credential_id` = ? AND `supervisor_credential_id` = ?",
			supervisor.EmployeeCredentialID, supervisor.SupervisorCredentialID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Some error occurred.", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetFacultySupervisorCount gets count of faculty supervisors.
func (service *ManagementService) GetFacultySupervisorCount(totalCount *admin.CountModel, tenantID, credentialID uuid.UUID) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	var totalCountInNumber int

	// Get count of supervisors.
	err := service.Repository.GetCount(uow, &admin.EmployeeSupervisor{}, &totalCountInNumber,
		repository.Join("JOIN credentials ON credentials.`id` = employee_supervisors.`supervisor_credential_id`"),
		repository.Filter("credentials.`deleted_at` IS NULL AND credentials.`tenant_id`=?", tenantID),
		repository.Filter("employee_supervisors.`employee_credential_id`=?", credentialID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// If there is no supervisor then it is not head faculty.
	if totalCountInNumber == 0 {
		totalCount.TotalCount = 1
		uow.Commit()
		return nil
	}

	// Get count of supervisors that are faculty.
	err = service.Repository.GetCount(uow, &admin.EmployeeSupervisor{}, &totalCountInNumber,
		repository.Join("JOIN credentials ON credentials.`id` = employee_supervisors.`supervisor_credential_id`"),
		repository.Filter("credentials.`deleted_at` IS NULL AND credentials.`tenant_id`=?", tenantID),
		repository.Filter("employee_supervisors.`employee_credential_id`=?", credentialID),
		repository.Filter("credentials.`faculty_id` IS NOT NULL"))
	if err != nil {
		uow.RollBack()
		return err
	}
	
	totalCount.TotalCount = totalCountInNumber

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doesTenantExist returns error if there is no tenant record in table.
func (service *ManagementService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *ManagementService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// validateFieldUniqueness will check the table for any repitition for unique fields.
func (service *ManagementService) validateFieldUniqueness(supervisor *admin.EmployeeSupervisor) error {
	// Return error if any record has the same email in DB.
	exists, err := repository.DoesRecordExist(service.DB, supervisor,
		repository.Filter("`employee_credential_id` = ? AND `supervisor_credential_id` = ? ",
			supervisor.EmployeeCredentialID, supervisor.SupervisorCredentialID))
	if err := util.HandleIfExistsError("Same supervisor exists.", exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

func (service *ManagementService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if firstName, ok := requestForm["firstName"]; ok {
		util.AddToSlice("`credentials`.`first_name`", "LIKE ?", "AND", "%"+firstName[0]+"%", &columnNames, &conditions, &operators, &values)
	}
	if lastName, ok := requestForm["lastName"]; ok {
		util.AddToSlice("`credentials`.`last_name`", "LIKE ?", "AND", "%"+lastName[0]+"%", &columnNames, &conditions, &operators, &values)
	}
	if email, ok := requestForm["email"]; ok {
		util.AddToSlice("`credentials`.`email`", "LIKE ?", "AND", "%"+email[0]+"%", &columnNames, &conditions, &operators, &values)
	}
	if roleID, ok := requestForm["roleID"]; ok {
		util.AddToSlice("`credentials`.`role_id`", "= ?", "AND", roleID, &columnNames, &conditions, &operators, &values)
	}
	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}


