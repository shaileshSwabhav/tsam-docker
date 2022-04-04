package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// BatchTopicAssignmentService provides method to Update, Delete, Add, Get Method For batch_topic_assignments.
type BatchTopicAssignmentService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// SessionProgrammingAssignmentService returns a new instance of SessionProgrammingAssignmentService.
func NewBatchTopicAssignmentService(db *gorm.DB, repository repository.Repository) *BatchTopicAssignmentService {
	return &BatchTopicAssignmentService{
		DB:         db,
		Repository: repository,
		association: []string{
			"ProgrammingQuestion", "ProgrammingQuestion.ProgrammingQuestionTypes",
		},
	}
}

// AddBatchTopicAssignment will add new assignment to topic assigned in batch.
func (service *BatchTopicAssignmentService) AddBatchTopicAssignment(topicAssignment *batch.TopicAssignment) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(topicAssignment, topicAssignment.CreatedBy)
	if err != nil {
		return err
	}

	// check if duplicate record exist.
	err = service.doesAssignmentExistForBatch(topicAssignment.TenantID, topicAssignment.BatchID,
		topicAssignment.ProgrammingQuestionID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	defer uow.Commit()
	if topicAssignment.DueDate != nil {
		now := time.Now()
		topicAssignment.AssignedDate = &now
	}

	err = service.Repository.Add(uow, topicAssignment)
	if err != nil {
		uow.RollBack()
		return err
	}
	return nil
}

// UpdatedTopicAssignment will add update assignment assigned to topic.
func (service *BatchTopicAssignmentService) UpdatedTopicAssignment(topicAssignment *batch.TopicAssignment) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(topicAssignment, topicAssignment.UpdatedBy)
	if err != nil {
		return err
	}

	err = service.doesBatchTopicAssignmentExist(topicAssignment.TenantID, topicAssignment.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	defer uow.RollBack()

	tempTopicAssignment := &batch.TopicAssignment{}
	err = service.Repository.GetRecordForTenant(uow, topicAssignment.TenantID, &tempTopicAssignment,
		repository.Filter("`id` = ?", topicAssignment.ID),
		repository.Select("`id`,`tenant_id`,`due_date`"))
	if err != nil {
		return err
	}

	if tempTopicAssignment.DueDate == nil && topicAssignment.DueDate != nil {
		now := time.Now()
		topicAssignment.AssignedDate = &now
	}

	err = service.Repository.Update(uow, topicAssignment)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// DeleteTopicAssignment will add delete assignment assigned to batch topic.
func (service *BatchTopicAssignmentService) DeleteTopicAssignment(topicAssignment *batch.TopicAssignment) error {

	// check if tenant exist.
	err := service.doesTenantExist(topicAssignment.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(topicAssignment.TenantID, topicAssignment.DeletedBy)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesBatchTopicAssignmentExist(topicAssignment.TenantID, topicAssignment.ID)
	if err != nil {
		return err
	}

	// check if assignment is answered by talent.
	exist, err := repository.DoesRecordExistForTenant(service.DB, topicAssignment.TenantID, talent.AssignmentSubmission{},
		repository.Filter("`batch_topic_assignment_id` = ?", topicAssignment.ID))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("Assignment cannot be deleted as talent have already answered it.")
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, &batch.TopicAssignment{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": topicAssignment.DeletedBy,
	}, repository.Filter("`id` = ?", topicAssignment.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBatchTopicAssignments will get all assignments for specified batch.
func (service *BatchTopicAssignmentService) GetAllTopicAssignments(tenantID, batchID uuid.UUID,
	batchTopicAssignment *[]batch.TopicAssignmentDTO, parser *web.Parser) error {
	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllForTenant(uow, tenantID, batchTopicAssignment,
		repository.Filter("`batch_id` = ?", batchID),
		service.addSearchQueries(parser.Form), repository.PreloadAssociations(service.association), repository.PreloadAssociations([]string{"Topic"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()

	return nil
}

// GetLatestBatchProgrammingAssignment will fetch all the programming assignments for the given batch session.
func (service *BatchTopicAssignmentService) GetLatestBatchProgrammingAssignment(tenantID, batchID uuid.UUID,
	sessionAssignments *[]batch.TopicAssignmentDTO, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, sessionAssignments,
		repository.Filter("`batch_id` = ?", batchID), repository.OrderBy("`created_at` DESC, `order`"),
		service.addSearchQueries(parser.Form), repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetSessionProgrammingAssignment will fetch all the programming assignments for the given batch session.
func (service *BatchTopicAssignmentService) GetSessionProgrammingAssignment(tenantID, batchTopicID uuid.UUID,
	sessionAssignments *[]batch.TopicAssignmentDTO, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch session exist.
	// err = service.doesTopicExist(tenantID, batchTopicID)
	// if err != nil {
	// 	return err
	// }

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, sessionAssignments, "`order`",
		repository.Filter("`batch_topic_id` = ?", batchTopicID), service.addSearchQueries(parser.Form),
		repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetSessionAssignmentList will return session_assignment list.
func (service *BatchTopicAssignmentService) GetSessionAssignmentList(tenantID, batchID uuid.UUID,
	sessionAssignmentList *[]batch.TopicAssignmentDTO, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch session exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// get all session assignments.
	err = service.Repository.GetAllForTenant(uow, tenantID, sessionAssignmentList,
		repository.Filter("batch_topic_assignments.`batch_id` = ?", batchID),
		service.addSearchQueries(parser.Form),
		// repository.OrderBy("module_topics.`order`, batch_topic_assignments.`order`"),
		repository.PreloadAssociations([]string{
			"ProgrammingQuestion", "Submissions",
		}))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllAssignmentsWithSubmissions will return batch_topic_assignments with scores.
func (service *BatchTopicAssignmentService) GetAllAssignmentsWithSubmissions(tenantID, batchID uuid.UUID,
	topicAssignment *[]batch.TopicAssignmentDTO, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch session exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// get all batch topic assignments.
	// err = service.Repository.GetAllInOrder(uow, topicAssignment, "`batch_topic_assignments`.`assigned_date`",
	// 	repository.Filter("`batch_topic_assignments`.`batch_id` = ?", batchID),
	// 	repository.Filter("`batch_topic_assignments`.`tenant_id` = ?", tenantID),
	// 	repository.Filter("`batch_topic_assignments`.`assigned_date` IS NOT NULL"),
	// 	// Need to move below join & optimize code #Niranjan.
	// 	repository.Join("INNER JOIN `batch_modules` ON "+
	// 		"`batch_modules`.`module_id` = `batch_topic_assignments`.`module_id` AND "+
	// 		"`batch_modules`.`tenant_id` = `batch_topic_assignments`.`tenant_id`"),
	// 	service.addSearchQueries(parser.Form),
	// 	repository.PreloadWithCustomCondition(repository.Preload{Schema: "ProgrammingQuestion"},
	// 		repository.Preload{Schema: "Submissions",
	// 			Queryprocessors: []repository.QueryProcessor{
	// 				repository.Join("INNER JOIN `talents` ON `talents`.`id` = `talent_assignment_submissions`.`talent_id`"),
	// 				repository.GroupBy("`talent_assignment_submissions`.`talent_id`," +
	// 					"`talent_assignment_submissions`.`batch_topic_assignment_id`"),
	// 				repository.OrderBy("`talents`.`first_name`,`talent_assignment_submissions`.`submitted_on` DESC"),
	// 				repository.PreloadAssociations([]string{"Talent"}),
	// 			},
	// 		},
	// 	),
	// )
	err = service.Repository.GetAllInOrder(uow, topicAssignment, "`batch_topic_assignments`.`assigned_date`",
		repository.Filter("`batch_topic_assignments`.`batch_id` = ?", batchID),
		repository.Filter("`batch_topic_assignments`.`tenant_id` = ?", tenantID),
		repository.Filter("`batch_topic_assignments`.`assigned_date` IS NOT NULL"),
		// Need to move below join & optimize code #Niranjan.
		repository.Join("INNER JOIN `batch_modules` ON "+
			"`batch_modules`.`batch_id` = `batch_topic_assignments`.`batch_id` AND "+
			"`batch_modules`.`module_id` = `batch_topic_assignments`.`module_id` AND "+
			"`batch_modules`.`tenant_id` = `batch_topic_assignments`.`tenant_id`"),
		service.addSearchQueries(parser.Form),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "ProgrammingQuestion"},
			repository.Preload{Schema: "Submissions",
				Queryprocessors: []repository.QueryProcessor{
					repository.Join("INNER JOIN `talents` ON `talents`.`id` = `talent_assignment_submissions`.`talent_id`"),
					repository.OrderBy("`talents`.`first_name`,`talent_assignment_submissions`.`submitted_on` DESC"),
					repository.PreloadAssociations([]string{"Talent","TalentConceptRatings","TalentConceptRatings.ModuleProgrammingConcept"}),
				},
			},
		),
	)
	if err != nil {
		uow.RollBack()
		return err
	}

	// if topicAssignment != nil {
	// 	for i,assign := range *topicAssignment {
	// 		for j,sub := range assign.Submissions {
	// 			if sub.SubmittedOn. {

	// 			}
	// 		}
	// 	}
	// }

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doesForeignKeyExist will check if all foreign keys are valid.
func (service *BatchTopicAssignmentService) doesForeignKeyExist(topicAssignment *batch.TopicAssignment, credentialID uuid.UUID) error {

	// check if tenant exist.
	err := service.doesTenantExist(topicAssignment.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(topicAssignment.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if question exist.
	err = service.doesProgramingQuestionExist(topicAssignment.TenantID, topicAssignment.ProgrammingQuestionID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExist(topicAssignment.TenantID, topicAssignment.BatchID)
	if err != nil {
		return err
	}

	// check if topic exist.
	err = service.doesModuleTopicExist(topicAssignment.TenantID, topicAssignment.TopicID)
	if err != nil {
		return err
	}

	// check if batch session topic exist.
	err = service.doesBatchSessionTopicExist(topicAssignment.TenantID, topicAssignment.BatchID, topicAssignment.TopicID)
	if err != nil {
		return err
	}

	// check if order already exist.
	// err = service.doesOrderExist(sessionAssignment.TenantID, sessionAssignment.ID,
	// 	sessionAssignment.TopicID, sessionAssignment.Order)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// addSearchQueries adds search criteria to get all batch_topic_assignments.
func (service *BatchTopicAssignmentService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if dueDate, ok := requestForm["dueDate"]; ok {
		util.AddToSlice("`batch_topic_assignments`.`days`", "<= ?", "AND", dueDate, &columnNames, &conditions, &operators, &values)
	}

	// String value should be enough, this should be moved #niranjan.
	facultyID := requestForm.Get("facultyID")
	if !util.IsEmpty(facultyID) {
		util.AddToSlice("`batch_modules`.`faculty_id`", "= ?", "AND",
			facultyID, &columnNames, &conditions, &operators, &values)
	}

	if topicID, ok := requestForm["topicID"]; ok {
		util.AddToSlice("`batch_topic_assignments`.`topic_id`", "IN(?)", "AND", topicID,
			&columnNames, &conditions, &operators, &values)
	}

	if assignedDate := requestForm.Get("assignedDate"); assignedDate == "0" {
		util.AddToSlice("`batch_topic_assignments`.`assigned_date`", "IS NULL", "AND", nil,
			&columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// returns error if there is no tenant record in table.
func (service *BatchTopicAssignmentService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *BatchTopicAssignmentService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch record in table for the given tenant.
func (service *BatchTopicAssignmentService) doesProgramingQuestionExist(tenantID, questionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingQuestion{},
		repository.Filter("`id` = ?", questionID))
	if err := util.HandleError("Invalid programming question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch record in table for the given tenant.
func (service *BatchTopicAssignmentService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no topic record in table for the given tenant.
func (service *BatchTopicAssignmentService) doesModuleTopicExist(tenantID, topicID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.ModuleTopic{},
		repository.Filter("`id` = ?", topicID))
	if err := util.HandleError("Invalid topic ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch session record in table for the given tenant.
func (service *BatchTopicAssignmentService) doesBatchSessionTopicExist(tenantID, batchID, topicID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.SessionTopic{},
		repository.Filter("`topic_id` = ? AND `batch_id` = ?", topicID, batchID))
	if err := util.HandleError("session for specfied topic not found", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no session record in table for the given tenant.
// func (service *SessionProgrammingAssignmentService) doesTopicExist(tenantID, batchTopicID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.BatchTopic{},
// 		repository.Filter("`id` = ?", batchTopicID))
// 	if err := util.HandleError("Invalid batch topic ID", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// returns error if there is no project assignment record in table for the given tenant.
// func (service *BatchTopicAssignmentService) doesProgrammingAssignmentExist(tenantID, assignmentID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingAssignment{},
// 		repository.Filter("`id` = ?", assignmentID))
// 	if err := util.HandleError("Invalid assignment ID", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// returns error if there is no project assignment in batch sessions record in table for the given tenant.
func (service *BatchTopicAssignmentService) doesBatchTopicAssignmentExist(tenantID, topicAssignmentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.TopicAssignment{},
		repository.Filter("`id` = ?", topicAssignmentID))
	if err := util.HandleError("Invalid session assignment ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is project assignment is already assigned to batch session for the given tenant.
func (service *BatchTopicAssignmentService) doesAssignmentExistForBatch(tenantID, batchID, programmingQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.TopicAssignment{},
		repository.Filter("`batch_id` = ? AND `programming_question_id` = ?", batchID, programmingQuestionID))
	if err := util.HandleIfExistsError("Coding question already assigned to a session in this batch", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is order exist for batch-session-programming-assignment in table for the given tenant.
// func (service *SessionProgrammingAssignmentService) doesOrderExist(tenantID, sessionAssignmentID, batchTopicID uuid.UUID, order uint) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.BatchTopicsProgrammingAssignment{},
// 		repository.Filter("`id` != ? AND `batch_topic_id` = ? AND `order` = ?", sessionAssignmentID, batchTopicID, order))
// 	if err := util.HandleIfExistsError("Order already exist", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }
