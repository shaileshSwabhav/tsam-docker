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
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TopicQuestionService provides method to Update, Delete, Add, Get Method For topic_programming_questions.
type TopicQuestionService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// CourseProgrammingAssignmentService returns a new instance of CourseProgrammingQuestionService.
func NewCourseTopicQuestionService(db *gorm.DB, repository repository.Repository) *TopicQuestionService {
	return &TopicQuestionService{
		DB:         db,
		Repository: repository,
		association: []string{
			// "SubTopic",
			"ProgrammingQuestion", "ProgrammingQuestion.ProgrammingQuestionTypes",
		},
	}
}

// AddTopicProgrammingQuestion will add programming_assignment for specified course.
func (service *TopicQuestionService) AddTopicProgrammingQuestion(topicQuestion *course.TopicProgrammingQuestion) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(topicQuestion, topicQuestion.CreatedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, topicQuestion)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// TopicCourseProgrammingQuestion will update programming_assignment for specified course.
func (service *TopicQuestionService) TopicCourseProgrammingQuestion(topicQuestion *course.TopicProgrammingQuestion) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(topicQuestion, topicQuestion.UpdatedBy)
	if err != nil {
		return err
	}

	// check if course programming assignment exist.
	err = service.doesCourseProgrammingQuestionExist(topicQuestion.TenantID, topicQuestion.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, topicQuestion)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteCourseProgrammingQuestion will delete programming_question for specified topic.
func (service *TopicQuestionService) DeleteCourseProgrammingQuestion(topicQuestion *course.TopicProgrammingQuestion) error {

	// check if tenant exist.
	err := service.doesTenantExist(topicQuestion.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(topicQuestion.TenantID, topicQuestion.DeletedBy)
	if err != nil {
		return err
	}

	// check if course programming assignment exist.
	err = service.doesCourseProgrammingQuestionExist(topicQuestion.TenantID, topicQuestion.ID)
	if err != nil {
		return err
	}

	exist, err := repository.DoesRecordExist(service.DB, &batch.TopicAssignment{},
		repository.Join("INNER JOIN `topic_programming_questions` ON `topic_programming_questions`.`programming_assignment_id` = "+
			"`batch_topics_programming_assignments`.`programming_assignment_id` AND `batch_topics_programming_assignments`.`tenant_id` = "+
			"`topic_programming_questions`.`tenant_id`"), repository.Filter("`topic_programming_questions`.`id` = ?", topicQuestion.ID),
		repository.Filter("`topic_programming_questions`.`deleted_at` IS NULL AND "+
			"`batch_topics_programming_assignments`.`tenant_id` = ?", topicQuestion.TenantID))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("Cannot delete assignment from course as it is assigned to batch-session")
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, &course.TopicProgrammingQuestion{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": topicQuestion.DeletedBy,
	}, repository.Filter("`id` = ?", topicQuestion.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTopicProgrammingQuestionList will fetch all the programming assignments for the specified topic.
func (service *TopicQuestionService) GetTopicProgrammingQuestionList(tenantID, topicID uuid.UUID,
	topicQuestions *[]course.TopicProgrammingQuestionDTO, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	err = service.doesSubTopicExist(tenantID, topicID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAll(uow, topicQuestions,
		repository.Join("INNER JOIN programming_questions ON topic_programming_questions.`programming_question_id` = "+
			"programming_questions.`id` AND topic_programming_questions.`tenant_id` = "+
			"programming_questions.`tenant_id`"), repository.Filter("topic_programming_questions.`tenant_id` =?", tenantID),
		repository.Filter("`topic_id` = ? AND `is_active` = ?", topicID, 1), service.addSearchQueries(parser.Form),
		repository.Filter("programming_questions.`deleted_at` IS NULL"), repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTopicProgrammingQuestions will fetch the programming questions for the specified topic with limit and offset.
func (service *TopicQuestionService) GetTopicProgrammingQuestions(tenantID, topicID uuid.UUID,
	topicAssignments *[]course.TopicProgrammingQuestionDTO, parser *web.Parser, totalCount *int) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	err = service.doesSubTopicExist(tenantID, topicID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	limit, offset := parser.ParseLimitAndOffset()

	err = service.Repository.GetAll(uow, topicAssignments,
		repository.Join("INNER JOIN programming_questions ON topic_programming_questions.`programming_question_id` = "+
			"programming_questions.`id` AND topic_programming_questions.`tenant_id` = "+
			"programming_questions.`tenant_id`"), repository.Filter("topic_programming_questions.`tenant_id` =?", tenantID),
		repository.Filter("topic_programming_questions.`topic_id` = ?", topicID),
		repository.Filter("programming_questions.`deleted_at` IS NULL"), service.addSearchQueries(parser.Form),
		repository.OrderBy("topic_programming_questions.`is_active` DESC"),
		repository.PreloadAssociations(service.association), repository.Paginate(limit, offset, totalCount))
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
func (service *TopicQuestionService) doesForeignKeyExist(topicAssignment *course.TopicProgrammingQuestion,
	credentialID uuid.UUID) error {

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

	// check if sub topic exist.
	err = service.doesSubTopicExist(topicAssignment.TenantID, topicAssignment.TopicID)
	if err != nil {
		return err
	}

	// check if programming assignment exist.
	err = service.doesProgrammingQuestionExist(topicAssignment.TenantID, topicAssignment.ProgrammingQuestionID)
	if err != nil {
		return err
	}

	// check if order already exist.
	// err = service.doesOrderExist(topicAssignment.TenantID, topicAssignment.TopicID,
	// 	topicAssignment.ID, topicAssignment.Order)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// addSearchQueries adds search criteria to get all topic_programming_questions.
func (service *TopicQuestionService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	if len(requestForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("`topic_programming_questions`.`is_active`", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	if _, ok := requestForm["name"]; ok {
		util.AddToSlice("`programming_questions`.`title`", "= ?", "AND", "%"+requestForm.Get("name")+"%", &columnNames, &conditions, &operators, &values)
	}

	if subTopicID, ok := requestForm["subTopicID"]; ok {
		util.AddToSlice("`topic_programming_questions`.`topic_id`", "= ?", "AND", subTopicID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// returns error if there is no tenant record in table.
func (service *TopicQuestionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *TopicQuestionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no concept record in table for the given tenant.
func (service *TopicQuestionService) doesConceptExist(tenantID, conceptID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingConcept{},
		repository.Filter("`id` = ?", conceptID))
	if err := util.HandleError("Invalid concept ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no concept record in table for the given tenant.
func (service *TopicQuestionService) doesSubTopicExist(tenantID, subTopicID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.ModuleTopic{},
		repository.Filter("`topic_id` = ? AND `topic_id` IS NOT NULL", subTopicID))
	if err := util.HandleError("Invalid sub topic ID", exists, err); err != nil {
		return err
	}
	return nil
}

// func (service *CourseProgrammingAssignmentService) doesCourseSessionExist(tenantID, courseID, courseSessionID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.CourseSession{},
// 		repository.Filter("`id` = ? AND `course_id` = ? AND `session_id` IS NULL", courseSessionID, courseID))
// 	if err := util.HandleError("Invalid course ID", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// returns error if there is no programming question record in table for the given tenant.
func (service *TopicQuestionService) doesProgrammingQuestionExist(tenantID, questionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingQuestion{},
		repository.Filter("`id` = ?", questionID))
	if err := util.HandleError("Invalid question ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no course programming assignment record in table for the given tenant.
func (service *TopicQuestionService) doesCourseProgrammingQuestionExist(tenantID, courseAssignmentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.TopicProgrammingQuestion{},
		repository.Filter("`id` = ?", courseAssignmentID))
	if err := util.HandleError("Invalid course assignment ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no programming assignment and course exist in table for the given tenant.
// func (service *CourseTopicAssignmentService) doesAssignmentExistInCourse(tenantID, courseAssignmentID, courseID,
// 	programmingAssignmentID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.TopicProgrammingAssignment{},
// 		repository.Filter("`id` != ? AND `course_id` = ? AND `programming_assignment_id` = ?",
// 			courseAssignmentID, courseID, programmingAssignmentID))
// 	if err := util.HandleIfExistsError("Assignment already assigned to course", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// returns error if there is no course programming assignment record in table for the given tenant.
// func (service *CourseTopicQuestionService) doesOrderExist(tenantID, topicID, topicAssignmentID uuid.UUID, order uint) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.TopicProgrammingQuestion{},
// 		repository.Filter("`id` != ? AND `order` = ? AND `is_active` = ? AND topic_id = ?",
// 			topicAssignmentID, order, true, topicID))
// 	if err := util.HandleIfExistsError("Order already exist", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (service *CourseTopicQuestionService) doesQuestionExistForTopic(tenantID, moduleID, programmingQuestionID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.TopicProgrammingQuestion{},
// 		repository.Filter("`module_id` = ? AND `programming_question_id` = ?", moduleID, programmingQuestionID))
// 	if err := util.HandleIfExistsError("Coding question already assigned to a session in this batch", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }
