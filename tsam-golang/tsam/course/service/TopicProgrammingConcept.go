package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CourseTopicConceptService provides method to Update, Delete, Add, Get Method For course_topic_programming_concepts.
type CourseTopicConceptService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewCourseTopicConceptService returns a new instance of CourseTopicConceptService.
func NewCourseTopicConceptService(db *gorm.DB, repository repository.Repository) *CourseTopicConceptService {
	return &CourseTopicConceptService{
		DB:          db,
		Repository:  repository,
		association: []string{},
	}
}

// AddTopicConcept will add new course topic programming concept to the table.
func (service *CourseTopicConceptService) AddTopicConcept(courseConcept *course.TopicProgrammingConcept) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(courseConcept, courseConcept.CreatedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, courseConcept)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateCourseProgrammingAssignment will update programming_assignment for specified course.
func (service *CourseTopicConceptService) UpdateCourseProgrammingAssignment(courseConcept *course.TopicProgrammingConcept) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(courseConcept, courseConcept.UpdatedBy)
	if err != nil {
		return err
	}

	// check if course programming assignment exist.
	err = service.doesCourseProgrammingConceptExist(courseConcept.TenantID, courseConcept.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, courseConcept)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteCourseProgrammingConcept will delete programming_concept for specified course.
func (service *CourseTopicConceptService) DeleteCourseProgrammingConcept(courseConcept *course.TopicProgrammingConcept) error {

	// check if tenant exist.
	err := service.doesTenantExist(courseConcept.TenantID)
	if err != nil {
		return err
	}

	// check if course exist.
	// err = service.doesCourseExist(courseConcept.TenantID, courseConcept.CourseID)
	// if err != nil {
	// 	return err
	// }

	// check if credential exist.
	err = service.doesCredentialExist(courseConcept.TenantID, courseConcept.DeletedBy)
	if err != nil {
		return err
	}

	// check if course programming assignment exist.
	err = service.doesCourseProgrammingConceptExist(courseConcept.TenantID, courseConcept.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, &course.TopicProgrammingConcept{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": courseConcept.DeletedBy,
	}, repository.Filter("`id` = ?", courseConcept.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTopicProgrammingConcepts will fetch the programming concepts with limit and offset.
func (service *CourseTopicConceptService) GetTopicProgrammingConcepts(tenantID, moduleTopicID uuid.UUID,
	topicProgrammingConcepts *[]course.TopicProgrammingConceptDTO, parser *web.Parser, totalCount *int) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// err = service.doesCourseExist(tenantID, courseID)
	// if err != nil {
	// 	return err
	// }

	uow := repository.NewUnitOfWork(service.DB, true)

	limit, offset := parser.ParseLimitAndOffset()

	// err = service.Repository.GetAll(uow, topicProgrammingConcepts,
	// 	repository.Join("INNER JOIN `programming_concepts` ON `topic_programming_concepts`.`programming_concept_id` = "+
	// 		"`programming_concepts`.`id` AND `topic_programming_concepts`.`tenant_id` = `programming_concepts`.`tenant_id`"),
	// 	repository.Filter("`topic_programming_concepts`.`tenant_id` =?", tenantID),
	// 	repository.Filter("`topic_programming_concepts`.`topic_id` = ?", moduleTopicID),
	// 	repository.Filter("`programming_concepts`.`deleted_at` IS NULL"), service.addSearchQueries(parser.Form),
	// 	// repository.OrderBy("topic_programming_concepts.`is_active` DESC, topic_programming_concepts.`order`"),
	// 	repository.PreloadAssociations(service.association), repository.Paginate(limit, offset, totalCount))
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }

	err = service.Repository.GetAll(uow, topicProgrammingConcepts,
		repository.Filter("`topic_programming_concepts`.`topic_id` = ?", moduleTopicID),
		repository.PreloadAssociations([]string{"ProgrammingConcept"}), service.addSearchQueries(parser.Form),
		// repository.OrderBy("topic_programming_concepts.`is_active` DESC, topic_programming_concepts.`order`"),
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

// addSearchQueries adds search criteria to get all topic_programming_concepts.
func (service *CourseTopicConceptService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	if len(requestForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	// if isActive, ok := requestForm["isActive"]; ok {
	// 	util.AddToSlice("`topic_programming_questions`.`is_active`", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	// }

	// if _, ok := requestForm["name"]; ok {
	// 	util.AddToSlice("`programming_assignments`.`title`", "= ?", "AND", "%"+requestForm.Get("name")+"%", &columnNames, &conditions, &operators, &values)
	// }

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesForeignKeyExist will check if all foreign keys are valid.
func (service *CourseTopicConceptService) doesForeignKeyExist(courseAssignment *course.TopicProgrammingConcept,
	credentialID uuid.UUID) error {

	// check if tenant exist.
	err := service.doesTenantExist(courseAssignment.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(courseAssignment.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if concept exist.
	err = service.doesConceptExist(courseAssignment.TenantID, courseAssignment.ProgrammingConceptID)
	if err != nil {
		return err
	}

	// check if sub topic exist.
	err = service.doesSubTopicExist(courseAssignment.TenantID, courseAssignment.TopicID)
	if err != nil {
		return err
	}

	// check if course exist.
	// err = service.doesCourseExist(courseAssignment.TenantID, courseAssignment.CourseID)
	// if err != nil {
	// 	return err
	// }

	// check if order already exist.
	// err = service.doesOrderExist(courseAssignment.TenantID, courseAssignment.ID,
	// 	courseAssignment.CourseID, courseAssignment.Order)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// returns error if there is no tenant record in table.
func (service *CourseTopicConceptService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *CourseTopicConceptService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no course record in table for the given tenant.
// func (service *CourseTopicConceptService) doesCourseExist(tenantID, courseID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Course{},
// 		repository.Filter("`id` = ?", courseID))
// 	if err := util.HandleError("Invalid course ID", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// returns error if there is no concept record in table for the given tenant.
func (service *CourseTopicConceptService) doesConceptExist(tenantID, conceptID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingConcept{},
		repository.Filter("`id` = ?", conceptID))
	if err := util.HandleError("Invalid concept ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no concept record in table for the given tenant.
func (service *CourseTopicConceptService) doesSubTopicExist(tenantID, subTopicID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.ModuleTopic{},
		repository.Filter("`id` = ? AND `topic_id` IS NOT NULL", subTopicID))
	if err := util.HandleError("Invalid sub topic ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no course programming assignment record in table for the given tenant.
// func (service *CourseTopicConceptService) doesOrderExist(tenantID, courseAssignmentID, courseID uuid.UUID, order uint) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.TopicProgrammingConcept{},
// 		repository.Filter("`id` != ? AND `course_id` = ? AND `order` = ? AND `is_active` = ?",
// 			courseAssignmentID, courseID, order, true))
// 	if err := util.HandleIfExistsError("Order already exist", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// returns error if there is no course programming assignment record in table for the given tenant.
func (service *CourseTopicConceptService) doesCourseProgrammingConceptExist(tenantID, courseAssignmentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.TopicProgrammingConcept{},
		repository.Filter("`id` = ?", courseAssignmentID))
	if err := util.HandleError("Invalid course concept ID", exists, err); err != nil {
		return err
	}
	return nil
}
