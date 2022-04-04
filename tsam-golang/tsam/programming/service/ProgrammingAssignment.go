package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingAssignmentService provides method to Update, Delete, Add, Get Method For programming_assignments.
type ProgrammingAssignmentService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// ProgrammingAssignmentService returns a new instance of ProgrammingAssignmentService.
func NewProgrammingAssignmentService(db *gorm.DB, repository repository.Repository) *ProgrammingAssignmentService {
	return &ProgrammingAssignmentService{
		DB:         db,
		Repository: repository,
		association: []string{
			"ProgrammingAssignmentSubTask", "ProgrammingAssignmentSubTask.Resource",
			"ProgrammingQuestion", "ProgrammingQuestion.ProgrammingQuestionTypes",
		},
	}
}

// AddProgrammingAssignment will add new assignment question to the table.
func (service *ProgrammingAssignmentService) AddProgrammingAssignment(assignment *programming.ProgrammingAssignment) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(assignment, assignment.CreatedBy)
	if err != nil {
		return err
	}

	if assignment.ProgrammingAssignmentSubTask != nil {
		for index := range assignment.ProgrammingAssignmentSubTask {
			assignment.ProgrammingAssignmentSubTask[index].TenantID = assignment.TenantID
			assignment.ProgrammingAssignmentSubTask[index].CreatedBy = assignment.CreatedBy
		}
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, assignment)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateProgrammingAssignment will update assignment question to the table.
func (service *ProgrammingAssignmentService) UpdateProgrammingAssignment(assignment *programming.ProgrammingAssignment) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(assignment, assignment.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if programming assignment exist.
	err = service.doesProgrammingAssignmentExist(assignment.TenantID, assignment.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempAssignment := programming.ProgrammingAssignment{}
	// Get createdby
	err = service.Repository.GetRecordForTenant(uow, assignment.TenantID, &tempAssignment,
		repository.Filter("`id` = ?", assignment.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	assignment.CreatedBy = tempAssignment.CreatedBy

	// Update programming assignment sub tasks.
	err = service.updateAssignmentSubTask(uow, assignment, assignment.UpdatedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.updateProgrammingAssignmentAssociations(uow, assignment)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.Save(uow, assignment)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteProgrammingAssignment will update assignment question to the table.
func (service *ProgrammingAssignmentService) DeleteProgrammingAssignment(assignment *programming.ProgrammingAssignment) error {

	// check if tenant exist.
	err := service.doesTenantExist(assignment.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(assignment.TenantID, assignment.DeletedBy)
	if err != nil {
		return err
	}

	// Check if programming assignment exist.
	err = service.doesProgrammingAssignmentExist(assignment.TenantID, assignment.ID)
	if err != nil {
		return err
	}

	// Check is assignment is used for course.
	exist, err := repository.DoesRecordExistForTenant(service.DB, assignment.TenantID, course.TopicProgrammingQuestion{},
		repository.Filter("`programming_assignment_id` = ?", assignment.ID))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("Assignment cannot be deleted as it is assigned to course.")
	}

	// Check if assignment is used for batch_session.
	exist, err = repository.DoesRecordExistForTenant(service.DB, assignment.TenantID, batch.TopicAssignment{},
		repository.Filter("`programming_assignment_id` = ?", assignment.ID))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("Assignment cannot be deleted as it is assigned to batch_session.")
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, &programming.ProgrammingAssignment{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": assignment.DeletedBy,
	}, repository.Filter("`id` = ?", assignment.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// // Delete programming_assignment from courses.
	// err = service.Repository.UpdateWithMap(uow, &course.CoursesProgrammingAssignment{}, map[string]interface{}{
	// 	"DeletedAt": time.Now(),
	// 	"DeletedBy": assignment.DeletedBy,
	// }, repository.Filter("`programming_assignment_id` = ?", assignment.ID))
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }

	// // Delete programming_assignment from batch_sessions.
	// err = service.Repository.UpdateWithMap(uow, &batch.BatchSessionsProgrammingAssignment{}, map[string]interface{}{
	// 	"DeletedAt": time.Now(),
	// 	"DeletedBy": assignment.DeletedBy,
	// }, repository.Filter("`programming_assignment_id` = ?", assignment.ID))
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }

	uow.Commit()
	return nil
}

// GetProgrammingAssignment will return all the programming assignments with limit and offset.
func (service *ProgrammingAssignmentService) GetProgrammingAssignment(tenantID uuid.UUID,
	assignments *[]programming.ProgrammingAssignmentDTO, parser *web.Parser, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	var queryProcessor repository.QueryProcessor

	// if util.IsEmpty(parser.Form.Get("isPaginate")) {
	limit, offset := parser.ParseLimitAndOffset()
	queryProcessor = repository.Paginate(limit, offset, totalCount)
	// }

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, assignments, "`created_at`",
		service.addSearchQueries(parser.Form), repository.PreloadAssociations(service.association),
		queryProcessor)
	if err != nil {
		uow.RollBack()
		return err
	}

	// if parser.Form.Get("isPaginate") == "0" {
	// 	*totalCount = len(*assignments)
	// }

	uow.Commit()
	return nil
}

// GetProgrammingAssignmentList will return list of programming assignments.
func (service *ProgrammingAssignmentService) GetProgrammingAssignmentList(tenantID uuid.UUID,
	assignments *[]list.ProgrammingAssignment, parser *web.Parser, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, assignments, "`created_at`",
		service.addSearchQueries(parser.Form))
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

// doesForeignKeyExist will check if all foreign keys are valid.
func (service *ProgrammingAssignmentService) doesForeignKeyExist(assignment *programming.ProgrammingAssignment, credentialID uuid.UUID) error {

	// check if tenant exist.
	err := service.doesTenantExist(assignment.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(assignment.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesAssignmentTitleExist(assignment.TenantID, assignment.ID, assignment.Title)
	if err != nil {
		return err
	}

	if assignment.ProgrammingQuestion != nil {
		for _, question := range assignment.ProgrammingQuestion {
			err = service.doesProgrammingQuestionExist(assignment.TenantID, question.ID)
			if err != nil {
				return err
			}
		}
	}

	if assignment.ProgrammingAssignmentSubTask != nil {
		for _, subTask := range assignment.ProgrammingAssignmentSubTask {
			err = service.doesResourceExist(assignment.TenantID, subTask.ResourceID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (service *ProgrammingAssignmentService) updateAssignmentSubTask(uow *repository.UnitOfWork,
	assignment *programming.ProgrammingAssignment, credentialID uuid.UUID) error {

	assignmentSubTaskMap := make(map[uuid.UUID]uint)
	assignmentSubTasks := &assignment.ProgrammingAssignmentSubTask
	tempAssignmentSubTasks := &[]programming.ProgrammingAssignmentSubTask{}

	err := service.Repository.GetAllForTenant(uow, assignment.TenantID, tempAssignmentSubTasks,
		repository.Filter("`programming_assignment_id`=?", assignment.ID))
	if err != nil {
		return err
	}

	for _, tempAssignmentSubTask := range *tempAssignmentSubTasks {
		assignmentSubTaskMap[tempAssignmentSubTask.ID] = assignmentSubTaskMap[tempAssignmentSubTask.ID] + 1
	}

	for _, assignmentSubTask := range *assignmentSubTasks {

		if util.IsUUIDValid(assignmentSubTask.ID) {
			assignmentSubTaskMap[assignmentSubTask.ID] = assignmentSubTaskMap[assignmentSubTask.ID] + 1
		} else {
			assignmentSubTask.ProgrammingAssignmentID = assignment.ID
			assignmentSubTask.CreatedBy = credentialID
			assignmentSubTask.TenantID = assignment.TenantID
			err = service.Repository.Add(uow, &assignmentSubTask)
			if err != nil {
				return err
			}
		}

		// update existing records
		if assignmentSubTaskMap[assignmentSubTask.ID] > 1 {
			assignmentSubTask.UpdatedBy = credentialID
			err = service.Repository.Update(uow, &assignmentSubTask)
			if err != nil {
				return err
			}
			assignmentSubTaskMap[assignmentSubTask.ID] = 0
		}
	}

	for _, tempBatchTime := range *tempAssignmentSubTasks {
		if assignmentSubTaskMap[tempBatchTime.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, programming.ProgrammingAssignmentSubTask{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempBatchTime.ID))
			if err != nil {
				return err
			}
		}
	}

	assignment.ProgrammingAssignmentSubTask = nil
	return nil
}

func (service *ProgrammingAssignmentService) updateProgrammingAssignmentAssociations(uow *repository.UnitOfWork,
	assignment *programming.ProgrammingAssignment) error {

	err := service.Repository.ReplaceAssociations(uow, assignment, "ProgrammingQuestion", assignment.ProgrammingQuestion)
	if err != nil {
		return err
	}
	assignment.ProgrammingQuestion = nil
	return nil
}

// addSearchQueries adds search criteria to get all programming_assignments.
func (service *ProgrammingAssignmentService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["title"]; ok {
		util.AddToSlice("`title`", "LIKE ?", "AND", "%"+requestForm.Get("title")+"%", &columnNames, &conditions, &operators, &values)
	}

	if programmingAssignmentType, ok := requestForm["programmingAssignmentType"]; ok {
		util.AddToSlice("`programming_assignment_type`", "= ?", "AND", programmingAssignmentType, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// returns error if there is no tenant record in table.
func (service *ProgrammingAssignmentService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *ProgrammingAssignmentService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no programming question record in table for the given tenant.
func (service *ProgrammingAssignmentService) doesProgrammingQuestionExist(tenantID, questionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingQuestion{},
		repository.Filter("`id` = ?", questionID))
	if err := util.HandleError("Invalid programming question ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no resource record in table for the given tenant.
func (service *ProgrammingAssignmentService) doesResourceExist(tenantID, resourceID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, resource.Resource{},
		repository.Filter("`id` = ?", resourceID))
	if err := util.HandleError("Invalid resource ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no programming assignment record in table for the given tenant.
func (service *ProgrammingAssignmentService) doesProgrammingAssignmentExist(tenantID, assignmentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingAssignment{},
		repository.Filter("`id` = ?", assignmentID))
	if err := util.HandleError("Invalid programming assignmnet ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is same title for programming assignment record in table for the given tenant.
func (service *ProgrammingAssignmentService) doesAssignmentTitleExist(tenantID, assignmentID uuid.UUID, title string) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingAssignment{},
		repository.Filter("`id` != ? AND `title` = ?", assignmentID, title))
	if err := util.HandleIfExistsError("Assignment title already exist", exists, err); err != nil {
		return err
	}
	return nil
}
