package service

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingQuestionTestCaseService provide method to update, delete, add, get method for programming question test case.
type ProgrammingQuestionTestCaseService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewProgrammingQuestionTestCaseService returns new instance of ProgrammingQuestionTestCaseService.
func NewProgrammingQuestionTestCaseService(db *gorm.DB, repository repository.Repository) *ProgrammingQuestionTestCaseService {
	return &ProgrammingQuestionTestCaseService{
		DB:         db,
		Repository: repository,
	}
}

// AddProgrammingQuestionTestCase adds new programming question test case to database.
func (service *ProgrammingQuestionTestCaseService) AddProgrammingQuestionTestCase(testCase *programming.ProgrammingQuestionTestCase) error {

	// Validate tenant id.
	err := service.doesTenantExist(testCase.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(testCase.CreatedBy, testCase.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate if programming question exists.
	err = service.doesProgrammingQuestionExist(testCase.ProgrammingQuestionID, testCase.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	//  Creating unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add programming question test case to database.
	err = service.Repository.Add(uow, testCase)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateProgrammingQuestionTestCase updates programming question test case to database.
func (service *ProgrammingQuestionTestCaseService) UpdateProgrammingQuestionTestCase(testCase *programming.ProgrammingQuestionTestCase) error {

	// Validate tenant ID.
	err := service.doesTenantExist(testCase.TenantID)
	if err != nil {
		return err
	}

	// Validate programming question test case ID.
	err = service.doesProgrammingQuestionTestCaseExist(testCase.ID, testCase.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(testCase.UpdatedBy, testCase.TenantID)
	if err != nil {
		return err
	}

	// Validate if programming question exists.
	err = service.doesProgrammingQuestionExist(testCase.ProgrammingQuestionID, testCase.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update programming question test case.
	err = service.Repository.Update(uow, testCase)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteProgrammingQuestionTestCase deletes programming question test case from database.
func (service *ProgrammingQuestionTestCaseService) DeleteProgrammingQuestionTestCase(testCase *programming.ProgrammingQuestionTestCase) error {

	credentialID := testCase.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(testCase.TenantID)
	if err != nil {
		return err
	}

	// Validate programming question test case ID.
	err = service.doesProgrammingQuestionTestCaseExist(testCase.ID, testCase.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, testCase.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update programming question test case for updating deleted_by and deleted_at fields of programming question test case.
	if err := service.Repository.UpdateWithMap(uow, testCase, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", testCase.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Programming question test case could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetProgrammingQuestionTestCase returns one programming question test case.
func (service *ProgrammingQuestionTestCaseService) GetProgrammingQuestionTestCase(testCase *programming.ProgrammingQuestionTestCase,
	tenantID uuid.UUID) error {

	// Validate tenant ID.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get one programming question test case by id.
	err = service.Repository.GetForTenant(uow, tenantID, testCase.ID, testCase)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetProgrammingQuestionTestCases returns all programming question test cases.
func (service *ProgrammingQuestionTestCaseService) GetProgrammingQuestionTestCases(testCases *[]programming.ProgrammingQuestionTestCase,
	tenantID, questionID uuid.UUID) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate programming question id.
	if err := service.doesProgrammingQuestionExist(questionID, tenantID); err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get programming question test cases from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, testCases, "`created_at`",
		repository.Filter("`programming_question_id`=?", questionID)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *ProgrammingQuestionTestCaseService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesProgrammingQuestionTestCaseExist validates if programming question test case exists or not in database.
func (service *ProgrammingQuestionTestCaseService) doesProgrammingQuestionTestCaseExist(testCaseID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingQuestionTestCase{},
		repository.Filter("`id` = ?", testCaseID))
	if err := util.HandleError("Invalid programming question test case ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *ProgrammingQuestionTestCaseService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesProgrammingQuestion validates if programming question exists or not in database.
func (service *ProgrammingQuestionTestCaseService) doesProgrammingQuestionExist(questionID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingQuestion{},
		repository.Filter("`id` = ?", questionID))
	if err := util.HandleError("Invalid programming question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
