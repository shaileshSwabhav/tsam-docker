package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TargetCommunityFunctionService provides methods to update, delete, add, get, get all and get all by department for
// target community function.
type TargetCommunityFunctionService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewTargetCommunityFunctionService returns new instance of TargetCommunityFunctionService.
func NewTargetCommunityFunctionService(db *gorm.DB, repository repository.Repository) *TargetCommunityFunctionService {
	return &TargetCommunityFunctionService{
		DB:         db,
		Repository: repository,
	}
}

// AddTargetCommunityFunction adds new targetCommunityFunction in database.
func (service *TargetCommunityFunctionService) AddTargetCommunityFunction(targetCommunityFunction *general.TargetCommunityFunction, uows ...*repository.UnitOfWork) error {

	// Validate tenant id.
	err := service.doesTenantExist(targetCommunityFunction.TenantID)
	if err != nil {
		return err
	}

	// Validate department id.
	err = service.doesDepartmentExist(targetCommunityFunction.TenantID, targetCommunityFunction.DepartmentID)
	if err != nil {
		return err
	}

	// Validate if same targetCommunityFunction function name exists.
	err = service.doesFunctionNameExist(targetCommunityFunction)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(targetCommunityFunction.CreatedBy, targetCommunityFunction.TenantID)
	if err != nil {
		return err
	}

	//  Creating unit of work.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Add targetCommunityFunction to database
	if err := service.Repository.Add(uow, targetCommunityFunction); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Target Community Function could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetTargetCommunityFunctionList gets all targetCommunityFunctions from database.
func (service *TargetCommunityFunctionService) GetTargetCommunityFunctionList(targetCommunityFunctions *[]general.TargetCommunityFunction, tenantID uuid.UUID) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get targetCommunityFunctions from database
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, targetCommunityFunctions, "`function_name`")
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTargetCommunityFunctions gets all targetCommunityFunctions from database.
func (service *TargetCommunityFunctionService) GetTargetCommunityFunctions(targetCommunityFunctions *[]general.TargetCommunityFunctionDTO,
	tenantID uuid.UUID,
	parser *web.Parser, totalCount *int) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()
	// Get targetCommunityFunctions from database
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, targetCommunityFunctions, "`function_name`",
		repository.Filter("`deleted_at` IS NULL"),
		service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{"Department"}),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTargetCommunityFunctionByDepartment gets all targetCommunityFunctions from database by specific department id.
func (service *TargetCommunityFunctionService) GetTargetCommunityFunctionByDepartment(targetCommunityFunctions *[]general.TargetCommunityFunction,
	tenantID uuid.UUID, departmentID uuid.UUID) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Validate department id.
	err = service.doesDepartmentExist(tenantID, departmentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get targetCommunityFunctions from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, targetCommunityFunctions, "`function_name`",
		repository.Filter("`department_id`=?", departmentID)); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// AddTargetCommunityFunctions adds multiple targetCommunityFunctions in database.
func (service *TargetCommunityFunctionService) AddTargetCommunityFunctions(targetCommunityFunctions *[]general.TargetCommunityFunction, targetCommunityFunctionsIDs *[]uuid.UUID,
	tenantID, credentialID uuid.UUID) error {
	// Check for same name conflict.
	for i := 0; i < len(*targetCommunityFunctions); i++ {
		for j := 0; j < len(*targetCommunityFunctions); j++ {
			if i != j && (*targetCommunityFunctions)[i].FunctionName == (*targetCommunityFunctions)[j].FunctionName &&
				(*targetCommunityFunctions)[i].DepartmentID == (*targetCommunityFunctions)[j].DepartmentID {
				log.NewLogger().Error("Function Name:" + (*targetCommunityFunctions)[j].FunctionName + " exists")
				return errors.NewValidationError("Function Name:" + (*targetCommunityFunctions)[j].FunctionName + " exists")
			}
		}
	}

	// Add individual targetCommunityFunction.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, targetCommunityFunction := range *targetCommunityFunctions {
		targetCommunityFunction.TenantID = tenantID
		targetCommunityFunction.CreatedBy = credentialID
		err := service.AddTargetCommunityFunction(&targetCommunityFunction, uow)
		if err != nil {
			return err
		}
		*targetCommunityFunctionsIDs = append(*targetCommunityFunctionsIDs, targetCommunityFunction.ID)
	}

	uow.Commit()
	return nil
}

// GetTargetCommunityFunction gets one targetCommunityFunction by specific targetCommunityFunction id from database.
func (service *TargetCommunityFunctionService) GetTargetCommunityFunction(targetCommunityFunction *general.TargetCommunityFunction) error {
	// Validate tenant id.
	err := service.doesTenantExist(targetCommunityFunction.TenantID)
	if err != nil {
		return err
	}

	// Validate targetCommunityFunction id.
	err = service.doesTargetCommunityFunctionExist(targetCommunityFunction.ID, targetCommunityFunction.TenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get targetCommunityFunction,
	if err := service.Repository.GetForTenant(uow, targetCommunityFunction.TenantID, targetCommunityFunction.ID, targetCommunityFunction); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// UpdateTargetCommunityFunction updates one targetCommunityFunction by specific targetCommunityFunction id in database.
func (service *TargetCommunityFunctionService) UpdateTargetCommunityFunction(targetCommunityFunction *general.TargetCommunityFunction) error {

	// Validate tenant ID.
	err := service.doesTenantExist(targetCommunityFunction.TenantID)
	if err != nil {
		return err
	}

	// Validate targetCommunityFunction ID.
	err = service.doesTargetCommunityFunctionExist(targetCommunityFunction.ID, targetCommunityFunction.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(targetCommunityFunction.UpdatedBy, targetCommunityFunction.TenantID)
	if err != nil {
		return err
	}

	// Validate if same function name exists.
	err = service.doesFunctionNameExist(targetCommunityFunction)
	if err != nil {
		return err
	}

	// Validate department id.
	err = service.doesDepartmentExist(targetCommunityFunction.TenantID, targetCommunityFunction.DepartmentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update targetCommunityFunction.
	if err := service.Repository.Update(uow, targetCommunityFunction); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Target Community Function could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteTargetCommunityFunction deletes one targetCommunityFunction by specific targetCommunityFunction id from database.
func (service *TargetCommunityFunctionService) DeleteTargetCommunityFunction(targetCommunityFunction *general.TargetCommunityFunction) error {
	credentialID := targetCommunityFunction.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(targetCommunityFunction.TenantID)
	if err != nil {
		return err
	}

	// Validate targetCommunityFunction ID.
	err = service.doesTargetCommunityFunctionExist(targetCommunityFunction.ID, targetCommunityFunction.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, targetCommunityFunction.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update targetCommunityFunction for updating deleted_by and deleted_at fields of targetCommunityFunction.
	if err := service.Repository.UpdateWithMap(uow, targetCommunityFunction, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", targetCommunityFunction.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Target Community Function could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *TargetCommunityFunctionService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant(parent tenant) exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDepartmentExist validates if department exists or not in database.
func (service *TargetCommunityFunctionService) doesDepartmentExist(tenantID uuid.UUID, departmentID uuid.UUID) error {
	// Check parent department exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Department{},
		repository.Filter("`id` = ?", departmentID))
	if err := util.HandleError("Invalid department ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTargetCommunityFunctionExist if targetCommunityFunction exists or not in database.
func (service *TargetCommunityFunctionService) doesTargetCommunityFunctionExist(targetCommunityFunctionsID uuid.UUID, tenantID uuid.UUID) error {
	// Check targetCommunityFunction exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.TargetCommunityFunction{},
		repository.Filter("`id`=?", targetCommunityFunctionsID))
	if err := util.HandleError("Invalid Target Community Function ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if credental exists or not in database.
func (service *TargetCommunityFunctionService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	// Check if credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{}, repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesFunctionNameExist checks if targetCommunityFunction's function name already exists or not for the
// speicifc department in database.
func (service *TargetCommunityFunctionService) doesFunctionNameExist(targetCommunityFunction *general.TargetCommunityFunction) error {
	// Check for same function name conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, targetCommunityFunction.TenantID, &general.TargetCommunityFunction{},
		repository.Filter("`function_name`=? AND `department_id` = ? AND `id`!=?",
			targetCommunityFunction.FunctionName, targetCommunityFunction.DepartmentID, targetCommunityFunction.ID))
	if err := util.HandleIfExistsError("Function name exists", exists, err); err != nil {
		return err
	}
	return nil
}

// addSearchQueries adds search critera.
func (service *TargetCommunityFunctionService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if _, ok := requestForm["functionName"]; ok {
		util.AddToSlice("`function_name`", "LIKE ?", "AND", "%"+requestForm.Get("functionName")+"%", &columnNames, &conditions, &operators, &values)
	}
	if departmentID, ok := requestForm["departmentID"]; ok {
		util.AddToSlice("`department_id`", "= ?", "AND", departmentID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
