package service

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	general "github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingQuestionSolutionService provides method to update, delete, add, get all, get one for programming question solutions.
type ProgrammingQuestionSolutionService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// solutionAssociations provides preload associations array for programming question solution.
var solutionAssociations []string = []string{"ProgrammingLanguage"}

// NewProgrammingQuestionSolutionService returns new instance of ProgrammingQuestionSolutionService.
func NewProgrammingQuestionSolutionService(db *gorm.DB, repository repository.Repository) *ProgrammingQuestionSolutionService {
	return &ProgrammingQuestionSolutionService{
		DB:         db,
		Repository: repository,
	}
}

// AddProgrammingQuestionSolution adds one programming question solution to database.
func (service *ProgrammingQuestionSolutionService) AddProgrammingQuestionSolution(solution *programming.ProgrammingQuestionSolution) error {

	// Get credential id from CreatedBy field of programming question solution(set in controller).
	credentialID := solution.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(solution.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, solution.TenantID); err != nil {
		return err
	}

	// Validate programming question id.
	if err := service.doesProgrammingQuestionExist(solution.TenantID, solution.ProgrammingQuestionID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(solution); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add programming question solution to database.
	if err := service.Repository.Add(uow, solution); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Programming question solution could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetProgrammingQuestionSolutions gets all programming question solutions from database.
func (service *ProgrammingQuestionSolutionService) GetProgrammingQuestionSolutions(solutions *[]programming.ProgrammingQuestionSolutionDTO,
	tenantID uuid.UUID, questionID uuid.UUID, uows ...*repository.UnitOfWork) error {

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate programming question id.
	if err := service.doesProgrammingQuestionExist(tenantID, questionID); err != nil {
		return err
	}

	// Get programming question solutions from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, solutions, "`created_at`",
		repository.Filter("`programming_question_id`=?", questionID),
		repository.PreloadAssociations(solutionAssociations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// GetProgrammingQuestionSolution gets one programming question solution form database.
func (service *ProgrammingQuestionSolutionService) GetProgrammingQuestionSolution(solution *programming.ProgrammingQuestionSolutionDTO,
	tenantID uuid.UUID) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate programming question id.
	if err := service.doesProgrammingQuestionExist(tenantID, solution.ProgrammingQuestionID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get programming question solution.
	if err := service.Repository.GetForTenant(uow, tenantID, solution.ID, solution,
		repository.PreloadAssociations(solutionAssociations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// UpdateProgrammingQuestionSolution updates programming question solution in Database.
func (service *ProgrammingQuestionSolutionService) UpdateProgrammingQuestionSolution(solution *programming.ProgrammingQuestionSolution) error {

	// Get credential id from UpdatedBy field of programming question solution(set in controller).
	credentialID := solution.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(solution.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, solution.TenantID); err != nil {
		return err
	}

	// Validate programming question id.
	if err := service.doesProgrammingQuestionExist(solution.TenantID, solution.ProgrammingQuestionID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(solution); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting programming question solution already present in database.
	tempSolution := programming.ProgrammingQuestionSolution{}

	// Get programming question solution for getting created_by field of programming question solution from database.
	if err := service.Repository.GetForTenant(uow, solution.TenantID, solution.ID, &tempSolution); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp programming question solution to programming question solution to be updated.
	solution.CreatedBy = tempSolution.CreatedBy

	// Update programming question solution.
	if err := service.Repository.Save(uow, solution); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Programming question solution could not be updated", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// DeleteProgrammingQuestionSolution deletes one programming question solution form database.
func (service *ProgrammingQuestionSolutionService) DeleteProgrammingQuestionSolution(solution *programming.ProgrammingQuestionSolution) error {

	// Starting new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get credential id from DeletedBy field of programming question solution(set in controller).
	credentialID := solution.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(solution.TenantID); err != nil {
		return err
	}

	// Validate programming question solution id.
	if err := service.doesProgrammingQuestionSolutionExist(solution.ID, solution.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, solution.TenantID); err != nil {
		return err
	}

	// Validate programming question solution id.
	if err := service.doesProgrammingQuestionExist(solution.TenantID, solution.ProgrammingQuestionID); err != nil {
		return err
	}

	// Update programming question solution for updating deleted_by and deleted_at field of programming question solution.
	if err := service.Repository.UpdateWithMap(uow, &programming.ProgrammingQuestionSolution{}, map[string]interface{}{
		"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`tenant_id`=? AND `id`=?", solution.TenantID, solution.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Programming question solution could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *ProgrammingQuestionSolutionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *ProgrammingQuestionSolutionService) doesCredentialExist(credentialID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesProgrammingQuestionExist validates if programming question exists or not in database.
func (service *ProgrammingQuestionSolutionService) doesProgrammingQuestionExist(tenantID, programmingQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingQuestion{},
		repository.Filter("`id` = ?", programmingQuestionID))
	if err := util.HandleError("Invalid programming question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesProgrammingQuestionSolutionExist validates if programming question solution exists or not in database.
func (service *ProgrammingQuestionSolutionService) doesProgrammingQuestionSolutionExist(solutionID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingQuestionSolution{},
		repository.Filter("`id` = ?", solutionID))
	if err := util.HandleError("Invalid programming question solution ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExist validates if language id exists or not in database.
func (service *ProgrammingQuestionSolutionService) doForeignKeysExist(solution *programming.ProgrammingQuestionSolution) error {

	// Check if programming language id exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.ProgrammingLanguage{},
		repository.Filter("`id` = ?", solution.ProgrammingLanguageID))
	if err := util.HandleError("Invalid programming language ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}
