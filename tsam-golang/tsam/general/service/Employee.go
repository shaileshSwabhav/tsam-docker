package service

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// EmployeeService provides methods to do different CRUD operations on employee table.
type EmployeeService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewEmployeeService returns a new instance Of EmployeeService.
func NewEmployeeService(db *gorm.DB, repository repository.Repository) *EmployeeService {
	return &EmployeeService{
		DB:         db,
		Repository: repository,
		associations: []string{
			// "Academics", "Academics.Specialization", "Academics.Degree",
			// "Experiences", "Experiences.Technologies", "Experiences.Designation",
			"Technologies",
			"Role",
			"Country", "State",
		},
	}
}

// AddEmployee will add new employee to the table and also create a login for the employee
func (service *EmployeeService) AddEmployee(employee *general.Employee) error {

	// extract all foreign key id's
	service.extractID(employee)

	// check if foreign key exist
	err := service.doForeignKeysExist(employee, employee.CreatedBy)
	if err != nil {
		return err
	}

	tenantID := employee.TenantID

	// creating login service
	credentialService := NewCredentialService(service.DB, service.Repository)

	// // Checks if credentialID has permission to add new employee
	// err = credentialService.ValidatePermission(tenantID, employee.CreatedBy, "/admin/employee/other", "add")
	// if err != nil {
	// 	return err
	// }

	uow := repository.NewUnitOfWork(service.DB, false)

	// check if email exist
	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &general.Employee{},
		repository.Filter("`email`=?", employee.Email))
	if err != nil {
		return err
	}
	if exist {
		// create credential for employee and then prompt error for email already exists
		// get employee existing with email
		existingEmployee := general.Employee{}
		err = service.Repository.GetRecordForTenant(uow, tenantID, &existingEmployee,
			repository.Filter("`email`=?", employee.Email))
		if err != nil {
			uow.RollBack()
			return err
		}

		// assign firstName as password
		existingEmployee.Password = strings.ToLower(existingEmployee.FirstName)

		err = service.addEmployeeCredential(uow, credentialService, &existingEmployee)
		if err != nil {
			uow.RollBack()
			return err
		}

		// update password for existing record
		err = service.Repository.UpdateWithMap(uow, &general.Employee{}, map[interface{}]interface{}{
			"password": existingEmployee.Password,
		}, repository.Filter("`id`=?", existingEmployee.ID))
		if err != nil {
			uow.RollBack()
			return err
		}

		uow.Commit()
		return errors.NewValidationError("Employee already exist similar email. New Login created.")
	}

	// assign employee code
	employee.Code, err = util.GenerateUniqueCode(uow.DB, employee.FirstName,
		"`code` = ?", employee)
	if err != nil {
		return err
	}

	employee.Password = strings.ToLower(employee.FirstName)
	// technologies := employee.Technologies
	// employee.Technologies = nil

	// Add employee to DB
	err = service.Repository.Add(uow, employee)
	if err != nil {
		uow.RollBack()
		return err
	}

	// adding technologies to mapped table
	// for _, technology := range technologies {

	// 	err = service.Repository.Add(uow, &general.EmployeeTechnologies{
	// 		EmployeeID:   employee.ID,
	// 		TechnologyID: technology.ID,
	// 	})
	// 	if err != nil {
	// 		uow.RollBack()
	// 		return err
	// 	}
	// }

	// Create login for employee
	err = service.addEmployeeCredential(uow, credentialService, employee)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateEmployee will update existing employee in the table
func (service *EmployeeService) UpdateEmployee(employee *general.Employee) error {

	// extract all foreign key id's
	service.extractID(employee)

	// check if foreign key's exist
	err := service.doForeignKeysExist(employee, employee.UpdatedBy)
	if err != nil {
		return err
	}

	// check if employee exist
	err = service.doesEmployeeExist(employee.TenantID, employee.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempEmployee := general.Employee{}
	err = service.Repository.GetRecordForTenant(uow, employee.TenantID, &tempEmployee,
		repository.Filter("`id`=?", employee.ID), repository.Select([]string{"`code`", "`created_by`", "`password`"}))
	if err != nil {
		return err
	}

	employee.Code = tempEmployee.Code
	employee.CreatedBy = tempEmployee.CreatedBy
	employee.Password = tempEmployee.Password

	// replace technology associations and make it nil
	err = service.replaceEmployeeAssociation(uow, employee, employee.UpdatedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.Save(uow, employee)
	if err != nil {
		uow.RollBack()
		return err
	}

	// updates is_active field of employee credential.
	err = service.updateEmployeeCredential(uow, employee)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteEmployee will delete specified employee
func (service *EmployeeService) DeleteEmployee(employee *general.Employee) error {

	// check is tenant exist
	err := service.doesTenantExist(employee.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExist(employee.TenantID, employee.DeletedBy)
	if err != nil {
		return err
	}

	// check if employee exist
	err = service.doesEmployeeExist(employee.TenantID, employee.ID)
	if err != nil {
		return err
	}

	credential := general.Credential{
		EmployeeID: &employee.ID,
	}
	// Start transaction
	uow := repository.NewUnitOfWork(service.DB, false)

	// Delete credentialService of employee
	credentialService := NewCredentialService(service.DB, service.Repository)
	err = credentialService.DeleteCredential(&credential, employee.TenantID, employee.DeletedBy, *credential.EmployeeID, "employee_id=?", uow)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.UpdateWithMap(uow, &general.Employee{}, map[string]interface{}{
		"DeletedBy": employee.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id`=?", employee.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Delete employee credential.
	err = service.Repository.UpdateWithMap(uow, &general.Credential{}, map[string]interface{}{
		"DeletedBy": employee.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`employee_id` = ?", employee.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllEmployee will return all the employees from the table with limit and offset
func (service *EmployeeService) GetAllEmployee(employees *[]general.EmployeeDTO, tenantID uuid.UUID, form url.Values,
	limit, offset int, totalCount *int) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return gorm.ErrCantStartTransaction
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, employees, "employees.`first_name`",
		repository.OrderBy("employees.`is_active` DESC"), service.addSearchQueries(form),
		repository.PreloadAssociations(service.associations), repository.Paginate(limit, offset, totalCount))
	if err != nil {
		return err
	}

	return nil
}

// GetAllEmployeeList returns an employee list.
func (service *EmployeeService) GetAllEmployeeList(employees *[]college.Developer, tenantID uuid.UUID) error {

	// Validates foreign keys of faculty
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start transaction.
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllForTenant(uow, tenantID, employees,
		repository.Filter("`is_active` = ?", true))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

func (service *EmployeeService) doForeignKeysExist(employee *general.Employee, credentialID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExist(employee.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	// err = service.doesCredentialExist(employee.TenantID, credentialID)
	// if err != nil {
	// 	return err
	// }

	// check if role exist
	err = service.doesRoleExist(employee.TenantID, employee.RoleID)
	if err != nil {
		return err
	}

	// check if type exist
	// err = service.doesTypeExist(employee.TenantID, employee.Type)
	// if err != nil {
	// 	return err
	// }

	// check if academics exist
	// if employee.Academics != nil {
	// 	for _, academic := range employee.Acadmeics {
	// 		err := service.doesDegreeExists(employee.TenantID, academic.DegreeID)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		err = service.doesSpecializationExists(employee.TenantID, academic.DegreeID,
	// 			academic.SpecializationID)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// check if experience exist
	// if employee.Experiences != nil {
	// 	for _, experience := range employee.Experiences {
	// 		err := service.doesDesignationExists(employee.TenantID, experience.DesignationID)
	// 		if err != nil {
	// 		// 			return err
	// 		}
	// 		// check if technology exists in DB
	// 		for _, technolgy := range experience.Technologies {
	// 			err = service.doesTechnologyExists(employee.TenantID, technolgy.ID)
	// 			if err != nil {
	// 			// 				return err
	// 			}
	// 		}
	// 	}
	// }

	// check if technology exists in DB
	if employee.Technologies != nil {
		for _, technolgy := range employee.Technologies {
			err := service.doesTechnologyExists(employee.TenantID, technolgy.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// replaceEmployeeAssociation will update technologies of employee
func (service *EmployeeService) replaceEmployeeAssociation(uow *repository.UnitOfWork, employee *general.Employee, credentialID uuid.UUID) error {

	if employee.Technologies != nil {
		if err := service.Repository.ReplaceAssociations(uow, employee, "Technologies", employee.Technologies); err != nil {
			return err
		}
	}

	employee.Technologies = nil
	return nil
}

// extractID extracts ID from the json and make the model nil, also checks if the record exist in the DB
func (service *EmployeeService) extractID(employee *general.Employee) {

	// extract countryID
	if employee.Country != nil {
		employee.CountryID = &employee.Country.ID
		// employee.Country = nil
	}

	// extract stateID
	if employee.State != nil {
		employee.StateID = &employee.State.ID
		// employee.State = nil
	}

	// extract role
	employee.RoleID = employee.Role.ID
	employee.Type = employee.Role.RoleName
	// employee.Role = nil

	// assign tenantID to academic
	// for i := range employee.Academics {

	// 	employee.Academics[i].TenantID = employee.TenantID
	// 	employee.Academics[i].CreatedBy = employee.CreatedBy

	// 	// extract degreeID
	// 	employee.Academics[i].DegreeID = employee.Academics[i].Degree.ID
	// 	employee.Academics[i].Degree = nil

	// 	// extract specializationID
	// 	employee.Academics[i].SpecializationID = academic.Specialization.ID
	// 	employee.Academics[i].Specialization = nil

	// }

	// assign tenantID to experiences
	// if employee.Experiences != nil {
	// 	for i := range employee.Experiences {

	// 		employee.Experiences[i].TenantID = employee.TenantID
	// 		employee.Experiences[i].CreatedBy = employee.CreatedBy

	// 		// extract designationID
	// 		employee.Experiences[i].DesignationID = employee.Experiences[i].Designation.ID
	// 		employee.Experiences[i].Designation = nil

	// 	}
	// }

}

// addEmployeeCredential creates a record in credentials table.
func (service *EmployeeService) addEmployeeCredential(uow *repository.UnitOfWork, credentialService *CredentialService, employee *general.Employee) error {

	// get roleID for employee
	// role := general.Role{}
	// err := service.Repository.GetRecordForTenant(uow, employee.TenantID, &role,
	// 	repository.Filter("`role_name`=?", employee.Type), repository.Select([]string{"`id`"}))
	// if err != nil {
	// 	return err
	// }

	credentials := general.Credential{
		FirstName:  employee.FirstName,
		LastName:   &employee.LastName,
		Email:      employee.Email,
		Contact:    employee.Contact,
		Password:   employee.Password,
		EmployeeID: &employee.ID,
		RoleID:     employee.RoleID,
	}
	credentials.TenantID = employee.TenantID
	credentials.CreatedBy = employee.CreatedBy
	err := credentialService.AddCredential(&credentials, uow)
	if err != nil {
		return err
	}
	return nil
}

// updateEmployeeCredential updates specified employee record in credentials table.
func (service *EmployeeService) updateEmployeeCredential(uow *repository.UnitOfWork, employee *general.Employee) error {

	err := service.Repository.UpdateWithMap(uow, &general.Credential{}, map[string]interface{}{
		"UpdatedBy": employee.UpdatedBy,
		"IsActive":  employee.IsActive,
	}, repository.Filter("`employee_id` = ?", employee.ID))
	if err != nil {
		return err
	}
	return nil
}

// returns error if there is no tenant record in table.
func (service *EmployeeService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *EmployeeService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no department record in table for the given tenant.
func (service *EmployeeService) doesEmployeeExist(tenantID, employeeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Employee{},
		repository.Filter("`id`=?", employeeID))
	if err := util.HandleError("Invalid employee ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no technology record in table for the given tenant.
func (service *EmployeeService) doesTechnologyExists(tenantID, technologyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Technology{},
		repository.Filter("`id`=?", technologyID))
	if err := util.HandleError("Technology not found", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no role record in table for the given tenant.
func (service *EmployeeService) doesRoleExist(tenantID, roleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Role{},
		repository.Filter("`id`=?", roleID))
	if err := util.HandleError("Invalid role", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *EmployeeService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["firstName"]; ok {
		util.AddToSlice("`first_name`", "LIKE ?", "AND", "%"+requestForm.Get("firstName")+"%", &columnNames, &conditions, &operators, &values)
	}

	if _, ok := requestForm["lastName"]; ok {
		util.AddToSlice("`last_name`", "LIKE ?", "AND", "%"+requestForm.Get("lastName")+"%", &columnNames, &conditions, &operators, &values)
	}

	if _, ok := requestForm["email"]; ok {
		util.AddToSlice("`email`", "LIKE ?", "AND", "%"+requestForm.Get("email")+"%", &columnNames, &conditions, &operators, &values)
	}

	if _, ok := requestForm["contact"]; ok {
		util.AddToSlice("`contact`", "LIKE ?", "AND", "%"+requestForm.Get("contact")+"%", &columnNames, &conditions, &operators, &values)
	}
	if isActive := requestForm.Get("isActive"); util.IsEmpty(isActive) {
		util.AddToSlice("`is_active`", "= ?", "AND", true, &columnNames, &conditions, &operators, &values)
	} else {
		util.AddToSlice("`is_active`", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// returns error if there is no email record in table for the given tenant.
// func (service *EmployeeService) doesEmailExists(tenantID uuid.UUID, email string) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Employee{},
// 		repository.Filter("`email`=?", email))
// 	if err := util.HandleIfExistsError("Email already exists", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // returns error if there is no degree record in table for the given tenant.
// func (service *EmployeeService) doesDegreeExists(tenantID, degreeID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Degree{},
// 		repository.Filter("`id`=?", degreeID))
// 	if err := util.HandleError("Degree not found", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // returns error if there is no specialization record in table for the given tenant.
// func (service *EmployeeService) doesSpecializationExists(tenantID, degreeID, specializationID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Specialization{},
// 		repository.Filter("`id`=? AND `degree_id`=?", specializationID, degreeID))
// 	if err := util.HandleError("Specialization not found", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // returns error if there is no designation record in table for the given tenant.
// func (service *EmployeeService) doesDesignationExists(tenantID, designationID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Designation{},
// 		repository.Filter("`id`=?", designationID))
// 	if err := util.HandleError("Designation not found", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// returns error if there is no role record in table for the given tenant.
// func (service *EmployeeService) doesTypeExist(tenantID uuid.UUID, typeName string) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Role{},
// 		repository.Filter("`role_name`=?", typeName))
// 	if err := util.HandleError("Invalid type", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// updateDeveloperAcademics will update academics of specified employee (NOT TESTED)
// func (service *DeveloperService) updateDeveloperAcademics(uow *repository.UnitOfWork, employee *general.Developer) error {

// tempAcademics := &[]fct.Academic{}
// academicMap := make(map[uuid.UUID]uint)

// get all academics of current employee
// err := service.Repository.GetAllForTenant(uow, employee.TenantID, tempAcademics,
// 	repository.Filter("`developer_id`=?", employee.ID))
// if err != nil {
// 	return err
// }

// populating academicMap
// for _, tempAcademic := range *tempAcademics {
// 	academicMap[tempAcademic.ID] = 1
// }
// for _, academic := range employee.Academics {

// 	if util.IsUUIDValid(academic.ID) {
// 		academicMap[academic.ID]++
// 	}

// 	// check if experience already exists in the DB
// 	if academicMap[academic.ID] > 1 {
// 		academic.UpdatedBy = employee.UpdatedBy
// 		err = service.Repository.Update(uow, &academic)
// 		if err != nil {
// 			return err
// 		}
// 		academicMap[academic.ID] = 0
// 	}

// 	if !util.IsUUIDValid(academic.ID) {

// 		academic.EmployeeID = employee.ID
// 		academic.TenantID = employee.TenantID
// 		academic.CreatedBy = employee.UpdatedBy
// 		err := service.Repository.Add(uow, &academic)
// 		if err != nil {
// 			return err
// 		}
// 	}
// }

// deleting all records where count is 1 as they have been removed from the experience
// for _, academic := range *tempAcademics {
// 	if academicMap[academic.ID] == 1 {
// 		// fmt.Println("******deleting academic")
// 		err = service.Repository.UpdateWithMap(uow, &academic, map[interface{}]interface{}{
// 			"DeletedBy": employee.UpdatedBy,
// 			"DeletedAt": time.Now(),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	academicMap[academic.ID] = 0
// }

// employee.Academics = nil

// return nil
// }

// updateDeveloperExperince updates experiences of the employee in the DB
// func (service *DeveloperService) updateDeveloperExperince(uow *repository.UnitOfWork, employee *general.Developer) error {

// tempExperiences := &[]fct.Experience{}
// get all experince of current employee
// err := service.Repository.GetAllForTenant(uow, employee.TenantID, tempExperiences,
// 	repository.Filter("`developer_id`=?", employee.ID))
// if err != nil {
// 	return err
// }

// experienceMap := make(map[uuid.UUID]uint)

// populating experinceMap
// for _, tempExperience := range *tempExperiences {
// 	experienceMap[tempExperience.ID] = 1
// }

// for _, experience := range employee.Experiences {

// 	if util.IsUUIDValid(experience.ID) {
// 		experienceMap[experience.ID]++
// 	}

// 	// check if experience already exists in the DB
// 	if experienceMap[experience.ID] > 1 && util.IsUUIDValid(experience.ID) {

// 		experience.UpdatedBy = employee.UpdatedBy
// 		experience.DevepolerID = employee.ID

// 		// replace employee experience technologies
// 		// fmt.Println("***updating experience technologies")
// 		if err := service.Repository.ReplaceAssociations(uow, experience, "Technologies",
// 			experience.Technologies); err != nil {
// 			return err
// 		}

// 		experience.Technologies = nil

// 		err = service.Repository.Save(uow, &experience)
// 		if err != nil {
// 			return err
// 		}
// 		experienceMap[experience.ID] = 0
// 	}

// 	if !util.IsUUIDValid(experience.ID) {

// 		experience.DevepolerID = employee.ID
// 		experience.TenantID = employee.TenantID
// 		experience.CreatedBy = employee.UpdatedBy
// 		err := service.Repository.Add(uow, &experience)
// 		if err != nil {
// 			return err
// 		}
// 	}
// }

// deleting all records where count is 1 as they have been removed from the experience
// for _, experience := range *tempExperiences {
// 	if experienceMap[experience.ID] == 1 {
// 		err = service.Repository.UpdateWithMap(uow, experience, map[interface{}]interface{}{
// 			"DeletedBy": employee.UpdatedBy,
// 			"DeletedAt": time.Now(),
// 		})
// 		if err != nil {
// 			log.NewLogger().Error(err.Error())
// 			return err
// 		}
// 	}
// 	experienceMap[experience.ID] = 0
// }

// setting experience to nil as all data as been added to db
// employee.Experiences = nil

// 	return nil
// }
