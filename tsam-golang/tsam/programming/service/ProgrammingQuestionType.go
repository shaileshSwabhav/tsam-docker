package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingQuestionTypeService provides methods to update, delete, add, get method for ProgrammingQuestionType.
type ProgrammingQuestionTypeService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewProgrammingQuestionTypeService returns new instance of ProgrammingQuestionTypeService.
func NewProgrammingQuestionTypeService(db *gorm.DB, repository repository.Repository) *ProgrammingQuestionTypeService {
	return &ProgrammingQuestionTypeService{
		DB:         db,
		Repository: repository,
	}
}

// AddProgrammingQuestionType adds new city to database.
func (service *ProgrammingQuestionTypeService) AddProgrammingQuestionType(programmingType *programming.ProgrammingQuestionType,
	uows ...*repository.UnitOfWork) error {

	// Check if foregin keys exist.
	err := service.doesForeignKeyExist(programmingType, programmingType.CreatedBy)
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

	// Add city.
	err = service.Repository.Add(uow, programmingType)
	if err != nil {
		uow.RollBack()
		return err
	}

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddProgrammingQuestionTypes adds multiple programming question types to database.
func (service *ProgrammingQuestionTypeService) AddProgrammingQuestionTypes(programmingTypes *[]programming.ProgrammingQuestionType,
	tenantID, credentialID uuid.UUID) error {

	// json validation for unique type field.
	err := service.JSONValidation(programmingTypes)
	if err != nil {
		return err
	}

	// Add individual country to database.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, programmingType := range *programmingTypes {
		programmingType.TenantID = tenantID
		programmingType.CreatedBy = credentialID
		err := service.AddProgrammingQuestionType(&programmingType, uow)
		if err != nil {
			return err
		}
	}

	uow.Commit()
	return nil
}

// UpdateProgrammingQuestionType updates programming question type to database.
func (service *ProgrammingQuestionTypeService) UpdateProgrammingQuestionType(programmingType *programming.ProgrammingQuestionType) error {

	// Check if foregin keys exist.
	err := service.doesForeignKeyExist(programmingType, programmingType.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if programming question type exist.
	err = service.doesProgrammingTypeExist(programmingType.TenantID, programmingType.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, programmingType)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// DeleteProgrammingQuestionType updates programming question type to database.
func (service *ProgrammingQuestionTypeService) DeleteProgrammingQuestionType(programmingType *programming.ProgrammingQuestionType) error {

	// Check if tenant exist.
	err := service.doesTenantExist(programmingType.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(programmingType.TenantID, programmingType.DeletedBy)
	if err != nil {
		return err
	}

	// Check if programming question type exist.
	err = service.doesProgrammingTypeExist(programmingType.TenantID, programmingType.ID)
	if err != nil {
		return err
	}

	// Check if programming questions are assigned to specified type.
	exist, err := repository.DoesRecordExistForTenant(service.DB, programmingType.TenantID, new(programming.ProgrammingQuestion),
		repository.Filter("`Programming_question_type_id` = ?", programmingType.ID))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("Programming questions exist for the specified type.")
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, new(programming.ProgrammingQuestionType), map[string]interface{}{
		"DeletedBy": programmingType.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", programmingType.ID))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetProgrammingQuestionTypes returns programming question types with limit and offset.
func (service *ProgrammingQuestionTypeService) GetProgrammingQuestionTypes(programmingType *[]programming.ProgrammingQuestionTypeDTO,
	form url.Values, tenantID uuid.UUID, limit, offset int, totalCount *int) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, programmingType, "`programming_type`",
		service.addSearchQueries(form), repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GeProgrammingQuestionTypeList returns all programming question types.
func (service *ProgrammingQuestionTypeService) GeProgrammingQuestionTypeList(programmingType *[]programming.ProgrammingQuestionTypeDTO,
	form url.Values, tenantID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, programmingType, "`programming_type`",
		service.addSearchQueries(form))
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

func (service *ProgrammingQuestionTypeService) doesForeignKeyExist(programmingType *programming.ProgrammingQuestionType,
	credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(programmingType.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(programmingType.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if programming type is unique.
	err = service.doesTypeExist(programmingType.TenantID, programmingType.ID, programmingType.ProgrammingType)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProgrammingQuestionTypeService) JSONValidation(programmingTypes *[]programming.ProgrammingQuestionType) error {

	programmingTypeMap := make(map[string]int)
	for _, programmingType := range *programmingTypes {
		programmingTypeMap[programmingType.ProgrammingType]++

		if programmingTypeMap[programmingType.ProgrammingType] > 1 {
			return errors.NewValidationError(programmingType.ProgrammingType + " is not unique.")
		}
	}

	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *ProgrammingQuestionTypeService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *ProgrammingQuestionTypeService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesProgrammingTypeExist validates if programming question type exists or not in database.
func (service *ProgrammingQuestionTypeService) doesProgrammingTypeExist(tenantID, typeID uuid.UUID) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingQuestionType{},
		repository.Filter("`id` = ?", typeID))
	if err := util.HandleError("Invalid programmingTypeID exist.", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTypeExist validates if programming question type exists or not in database.
func (service *ProgrammingQuestionTypeService) doesTypeExist(tenantID, typeID uuid.UUID,
	programmingType string) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingQuestionType{},
		repository.Filter("`id` != ? AND `programming_type` = ?", typeID, programmingType))
	if err := util.HandleIfExistsError(programmingType+" programming type already exist.", exists, err); err != nil {
		return err
	}
	return nil
}

// addSearchQueries will append search queries from queryParams to queryProcessor
func (service *ProgrammingQuestionTypeService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["programmingType"]; ok {
		util.AddToSlice("`programming_type`", "LIKE ?", "AND", "%"+requestForm.Get("programmingType")+"%", &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
