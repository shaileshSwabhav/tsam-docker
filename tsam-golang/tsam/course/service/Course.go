package service

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	crs "github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// Preload Associations
// var courseAssociations []string = []string{
// 	"Technologies", "Eligibility", "Eligibility.Technologies",
// 	// "Sessions",
// }

// "Sessions.Resources"

// CourseService Provide method to Update, Delete, Add, Get Method For Course.
type CourseService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewCourseService returns new instance of CourseService.
func NewCourseService(db *gorm.DB, repository repository.Repository) *CourseService {
	return &CourseService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Technologies", "Eligibility", "Eligibility.Technologies",
		},
	}
}

// AddCourse Add New Course status to Database.
func (service *CourseService) AddCourse(course *crs.Course) error {

	// check if tenantID exists
	err := service.doesTenantExists(course.TenantID)
	if err != nil {
		return err
	}

	// assign tenantID to eligibility
	if course.Eligibility != nil {
		course.Eligibility.CreatedBy = course.CreatedBy
		course.Eligibility.TenantID = course.TenantID
	}

	// credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// // check if credentialID exists and has permission to update the course
	// err = credentialService.ValidatePermission(course.TenantID, course.CreatedBy, "/course/master", "add")
	// if err != nil {
	// 	return err
	// }

	uow := repository.NewUnitOfWork(service.DB, false)

	// Assign Course Code
	course.Code, err = util.GenerateUniqueCode(uow.DB, course.Name,
		"`code` = ?", &crs.Course{})
	if err != nil {
		return errors.NewHTTPError("Fail to generate Code", http.StatusInternalServerError)
	}

	err = service.Repository.Add(uow, course)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateCourse Update Course to Database.
func (service *CourseService) UpdateCourse(course *crs.Course) error {

	// check if tenantID exists
	err := service.doesTenantExists(course.TenantID)
	if err != nil {
		return err
	}

	// check if tenantID exists
	err = service.doesCourseExists(course.TenantID, course.ID)
	if err != nil {
		return err
	}

	// assign tenantID to eligibility
	if course.Eligibility != nil {
		// course.Eligibility.CreatedBy = course.CreatedBy
		course.Eligibility.TenantID = course.TenantID
	}

	// credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// // check if credentialID exists and has permission to update the course
	// err = credentialService.ValidatePermission(course.TenantID, course.UpdatedBy, "/course/master", "update")
	// if err != nil {
	// 	// log.NewLogger().Error(err.Error())
	// 	return err
	// }

	// Start transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	tempCourse := crs.Course{}
	err = service.Repository.GetRecordForTenant(uow, course.TenantID, &tempCourse,
		repository.Filter("`id` = ?", course.ID), repository.Select([]string{"`created_by`", "`code`"}))
	if err != nil {
		uow.RollBack()
		return err
	}
	course.CreatedBy = tempCourse.CreatedBy
	course.Code = tempCourse.Code

	err = service.updateCourseEligibility(uow, course)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.updateCourseAssociation(uow, course)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.Save(uow, course)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteCourse Delete Course From Database.
func (ser *CourseService) DeleteCourse(course *crs.Course) error {

	// check if tenantID exists
	err := ser.doesTenantExists(course.TenantID)
	if err != nil {
		return err
	}

	// check if tenantID exists
	err = ser.doesCourseExists(course.TenantID, course.ID)
	if err != nil {
		return err
	}

	// credentialService := genService.NewCredentialService(ser.DB, ser.Repository)

	// // check if credentialID exists and has permission to update the course
	// err = credentialService.ValidatePermission(course.TenantID, course.DeletedBy, "/course/master", "delete")
	// if err != nil {
	// 	return err
	// }

	// Start transaction
	deletedBy := course.DeletedBy

	uow := repository.NewUnitOfWork(ser.DB, false)

	err = ser.Repository.GetRecordForTenant(uow, course.TenantID, course,
		repository.Filter("`id`=?", course.ID), repository.PreloadAssociations(ser.associations))
	if err != nil {
		return err
	}
	course.DeletedBy = deletedBy

	// deletes associations of course
	err = ser.deleteCourseAssociation(uow, course)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = ser.deleteSessions(uow, course)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = ser.Repository.UpdateWithMap(uow, &crs.Course{}, map[interface{}]interface{}{
		"DeletedBy": course.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id`=?", course.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetCoursesList Return Course List
func (service *CourseService) GetCoursesList(courses *[]list.Course, tenantID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, courses, "name",
		repository.Filter("`deleted_at` IS NULL"))
	if err != nil {
		uow.RollBack()
		return err
	}
	return nil
}

// GetCourse Return Course By ID.
func (service *CourseService) GetCourse(tenantID uuid.UUID, course *crs.CourseDTO) error {

	// Check if tenantID exists
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if course exists
	err = service.doesCourseExists(tenantID, course.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetForTenant(uow, tenantID, course.ID, course,
		repository.PreloadAssociations(service.associations))
	if err != nil {
		uow.RollBack()
		return err
	}

	var totalCount int
	err = service.Repository.GetCountForTenant(uow, tenantID, &crs.CourseSession{},
		&totalCount, repository.Filter("session_id IS NULL"))
	if err != nil {
		uow.RollBack()
		return err
	}

	course.TotalSessions = uint(totalCount)
	uow.Commit()
	return nil
}

// GetCourses Return All Course.
func (service *CourseService) GetCourses(courses *[]crs.CourseDTO, tenantID uuid.UUID,
	parser *web.Parser, totalCount *int) error {

	// Check if tenantID exists.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()

	var queryProcessors []repository.QueryProcessor

	queryProcessors = append(queryProcessors, service.addCourseSearchQueriesParams(parser.Form)...)
	queryProcessors = append(queryProcessors, repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount))

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, courses, "name",
		queryProcessors...)
	if err != nil {
		return err
	}

	for index := range *courses {
		// var totalCount int
		// err = service.Repository.GetCountForTenant(uow, tenantID, &crs.CourseSession{},
		// 	&totalCount, repository.Filter("session_id IS NULL"),
		// 	repository.Filter("course_id=?", (*courses)[index].ID))
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }
		// (*courses)[index].TotalSessions = uint(totalCount)

		err = service.getTotalModules(uow, &(*courses)[index], tenantID)
		if err != nil {
			uow.RollBack()
			return err
		}

		err = service.getTotalTopics(uow, &(*courses)[index], tenantID)
		if err != nil {
			uow.RollBack()
			return err
		}

		err = service.getTotalQuestions(uow, &(*courses)[index], tenantID)
		if err != nil {
			uow.RollBack()
			return err
		}

		err = service.getTotalConcepts(uow, &(*courses)[index], tenantID)
		if err != nil {
			uow.RollBack()
			return err
		}

	}

	uow.Commit()
	return nil
}

// GetCourseDetails returns specific details of single course by id.
func (service *CourseService) GetCourseDetails(course *crs.CourseDetails, tenantID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if course exist.
	err = service.doesCourseExists(tenantID, course.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetForTenant(uow, tenantID, course.ID, course)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Count total sessions of course.
	var totalCount int

	err = service.Repository.GetCount(uow, &crs.Course{}, &totalCount,
		repository.Join("INNER JOIN module_topics ON module_topics.`course_id` = courses.`id`"),
		repository.Join("INNER JOIN course_modules ON course_modules.`course_id` = courses.`id` AND"+
			" course_modules.`tenant_id = courses.`tenant_id`"),
		repository.Filter("courses.`id`=? AND module_topics.`topic_id` IS NULL", course.ID),
		repository.Filter("module_topics.`deleted_at` IS NULL AND courses.`deleted_at` IS NULL"),
		repository.Filter("module_topics.`tenant_id` = ? AND courses.`tenant_id`=?", tenantID, tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}
	course.SessionsCount = uint(totalCount)

	uow.Commit()
	return nil
}

// GetCourseMinimumDetails returns minimum specific details of single course by id.
func (service *CourseService) GetCourseMinimumDetails(courses *[]crs.CourseMinimumDetails, talentID, tenantID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if talent exist.
	err = service.doesTalentExist(tenantID, talentID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors.
	var queryProcessorsForSubQuery []repository.QueryProcessor
	queryProcessorsForSubQuery = append(queryProcessorsForSubQuery,
		repository.Select("IF(batch_talents.`talent_id`=?,1, 0) as enrolled, courses.*", talentID),
		repository.Table("courses"),
		repository.Join("LEFT JOIN batches ON batches.`course_id` = courses.`id`"),
		repository.Join("LEFT JOIN batch_talents ON batch_talents.`batch_id` = batches.`id`"),
		repository.Filter("courses.`deleted_at` IS NULL AND batches.`deleted_at` IS NULL AND batch_talents.`deleted_at` IS NULL"),
		repository.Filter("courses.`tenant_id`=?", tenantID))

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, crs.Course{}, queryProcessorsForSubQuery...)
	if err != nil {
		uow.RollBack()
		// log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get courses.
	if err := service.Repository.GetAll(uow, courses, repository.RawQuery("SELECT SUM(enrolled) as is_enrolled, subquery.* FROM ? as subquery group by courses.`id`", subQuery)); err != nil {
		uow.RollBack()
		// log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

func (service *CourseService) updateCourseEligibility(uow *repository.UnitOfWork, course *crs.Course) error {

	tempCourse := &crs.Course{}

	// get batch record
	err := service.Repository.GetRecordForTenant(uow, course.TenantID, tempCourse,
		repository.Filter("`id`=?", course.ID),
		repository.PreloadAssociations([]string{"Eligibility", "Eligibility.Technologies"}))
	if err != nil {
		return err
	}

	// no eligibilty exists
	if tempCourse.EligibilityID == nil && course.Eligibility == nil {
		return nil
	}

	// previously eligibility exists but now eligibility is removed
	if course.Eligibility == nil {

		err := service.Repository.UpdateWithMap(uow, general.Eligibility{}, map[interface{}]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": course.UpdatedBy,
		}, repository.Filter("`id`=?", tempCourse.EligibilityID))
		if err != nil {
			return err
		}

		err = service.Repository.RemoveAssociations(uow, tempCourse.Eligibility, "Technologies", tempCourse.Eligibility.Technologies)
		if err != nil {
			return err
		}

		// set course.eligibilityID to null
		err = service.Repository.UpdateWithMap(uow, &crs.Course{}, map[interface{}]interface{}{
			"EligibilityID": nil,
		}, repository.Filter("`id`=?", course.ID))
		if err != nil {
			return err
		}
		return nil
	}

	if tempCourse.EligibilityID != nil && course.Eligibility != nil {
		course.Eligibility.TenantID = course.TenantID
		course.Eligibility.UpdatedBy = course.UpdatedBy
		err = service.Repository.Update(uow, course.Eligibility)
		if err != nil {
			// log.NewLogger().Error(err.Error())
			return err
		}
		course.EligibilityID = &course.Eligibility.ID
		course.Eligibility = nil
		return nil
	}

	course.Eligibility.TenantID = course.TenantID
	course.Eligibility.CreatedBy = course.UpdatedBy
	err = service.Repository.Add(uow, course.Eligibility)
	if err != nil {
		// log.NewLogger().Error(err.Error())
		return err
	}
	course.EligibilityID = &course.Eligibility.ID
	course.Eligibility = nil

	return nil
}

// updateCourseAssociation Update Course Dependencies
func (service *CourseService) updateCourseAssociation(uow *repository.UnitOfWork, course *crs.Course) error {

	err := service.Repository.ReplaceAssociations(uow, course, "Technologies", course.Technologies)
	if err != nil {
		return err
	}
	course.Technologies = nil

	return nil
}

// deleteCourseAssociation Delete Course Dependencies
func (service *CourseService) deleteCourseAssociation(uow *repository.UnitOfWork, course *crs.Course) error {

	if course.EligibilityID != nil {

		err := service.Repository.UpdateWithMap(uow, general.Eligibility{}, map[interface{}]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": course.DeletedBy,
		}, repository.Filter("`id`=?", course.EligibilityID))
		if err != nil {
			return err
		}

		// if err := ser.Repository.RemoveAssociations(uow, course.Eligibility, "Technologies", course.Eligibility.Technologies); err != nil {
		// 	log.NewLogger().Error(err.Error())
		// 	return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		// }
	}

	// if err := ser.Repository.RemoveAssociations(uow, course, "Technologies", course.Technologies); err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }
	return nil
}

func (service *CourseService) getTotalModules(uow *repository.UnitOfWork, course *crs.CourseDTO,
	tenantID uuid.UUID) error {

	err := service.Repository.GetCountForTenant(uow, tenantID, &crs.CourseModule{}, &course.TotalModules,
		repository.Filter("course_modules.`course_id` = ?", course.ID))
	if err != nil {
		return err
	}

	return nil
}

func (service *CourseService) getTotalTopics(uow *repository.UnitOfWork, course *crs.CourseDTO,
	tenantID uuid.UUID) error {

	err := service.Repository.GetCount(uow, &crs.Course{}, &course.TotalTopics,
		repository.Join("INNER JOIN course_modules ON course_modules.`course_id` = courses.`id` AND "+
			" course_modules.`tenant_id` = courses.`tenant_id`"),
		repository.Join("INNER JOIN module_topics ON module_topics.`module_id` = course_modules.`module_id` "+
			" AND module_topics.`tenant_id` = course_modules.`tenant_id`"),
		repository.Filter("course_modules.`deleted_at` IS NULL AND courses.`id` = ?", course.ID),
		repository.Filter("module_topics.`deleted_at` IS NULL AND courses.`tenant_id` = ?", tenantID))
	if err != nil {
		return err
	}

	return nil
}

func (service *CourseService) getTotalConcepts(uow *repository.UnitOfWork, course *crs.CourseDTO,
	tenantID uuid.UUID) error {

	err := service.Repository.GetCount(uow, &crs.Course{}, &course.TotalConcepts,
		repository.Join("INNER JOIN course_modules ON course_modules.`course_id` = courses.`id` AND "+
			" course_modules.`tenant_id` = courses.`tenant_id`"),
		repository.Join("INNER JOIN module_topics ON module_topics.`module_id` = course_modules.`module_id` "+
			" AND module_topics.`tenant_id` = course_modules.`tenant_id`"),
		repository.Join("INNER JOIN topic_programming_concepts ON topic_programming_concepts.`topic_id` = module_topics.`id` "+
			" AND topic_programming_concepts.`tenant_id` = module_topics.`tenant_id`"),
		repository.Filter("course_modules.`deleted_at` IS NULL"),
		repository.Filter("module_topics.`deleted_at` IS NULL AND courses.`tenant_id` = ?", tenantID),
		repository.Filter("topic_programming_concepts.`deleted_at` IS NULL AND courses.`id` = ?", course.ID),
		repository.GroupBy("topic_programming_concepts.programming_concept_id"))
	if err != nil {
		return err
	}
	fmt.Println("=====================", course.TotalConcepts)
	return nil
}

func (service *CourseService) getTotalQuestions(uow *repository.UnitOfWork, course *crs.CourseDTO,
	tenantID uuid.UUID) error {

	err := service.Repository.GetCount(uow, &crs.Course{}, &course.TotalQuestions,
		repository.Join("INNER JOIN course_modules ON course_modules.`course_id` = courses.`id` AND "+
			" course_modules.`tenant_id` = courses.`tenant_id`"),
		repository.Join("INNER JOIN module_topics ON module_topics.`module_id` = course_modules.`module_id` "+
			" AND module_topics.`tenant_id` = course_modules.`tenant_id`"),
		repository.Join("INNER JOIN topic_programming_questions ON topic_programming_questions.`topic_id` = module_topics.`id` "+
			" AND topic_programming_questions.`tenant_id` = module_topics.`tenant_id`"),
		repository.Filter("course_modules.`deleted_at` IS NULL"),
		repository.Filter("module_topics.`deleted_at` IS NULL AND courses.`tenant_id` = ?", tenantID),
		repository.Filter("topic_programming_questions.`deleted_at` IS NULL AND courses.`id` = ?", course.ID))
	if err != nil {
		return err
	}

	return nil
}

// deleteSessions will delete all the sessions, sub-sessions their resources for the specified course
func (service *CourseService) deleteSessions(uow *repository.UnitOfWork, course *crs.Course) error {

	sessions := []crs.CourseSession{}
	err := service.Repository.GetAllForTenant(uow, course.TenantID, &sessions,
		repository.Filter("`course_id`=?", course.ID), repository.PreloadAssociations([]string{"Resources"}))
	if err != nil {
		return err
	}

	for _, session := range sessions {
		// delete resource
		// err = ser.Repository.UpdateWithMap(uow, crs.Resource{}, map[interface{}]interface{}{
		// 	"DeletedBy": course.DeletedBy,
		// 	"DeletedAt": time.Now(),
		// }, repository.Filter("`session_id`=?", session.ID))
		// if err != nil {
		// 	return err
		// }
		if len(session.Resources) > 0 {
			err = service.Repository.RemoveAssociations(uow, session, "Resources", session.Resources)
			if err != nil {
				uow.RollBack()
				return err
			}
		}
	}

	// delete sessions
	err = service.Repository.UpdateWithMap(uow, crs.CourseSession{}, map[string]interface{}{
		"DeletedBy": course.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`course_id`=?", course.ID))
	if err != nil {
		return err
	}

	// err = service.deleteCourseSessionFeedback(uow, course, sessions)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// deleteCourseSessionFeedback will delete all the feedback for courses
// func (service *CourseService) deleteCourseSessionFeedback(uow *repository.UnitOfWork, course *crs.Course,
// 	batchTopics []crs.CourseSession) error {

// 	for _, topic := range batchTopics {
// 		exist, err := repository.DoesRecordExistForTenant(service.DB, course.TenantID, batch.BatchTopic{},
// 			repository.Filter("`module_topic_id` = ?", topic.ID))
// 		if err != nil {
// 			return err
// 		}
// 		if exist {
// 			tempBatchTopic := batch.BatchTopic{}
// 			err = service.Repository.GetRecordForTenant(uow, course.TenantID, &tempBatchTopic,
// 				repository.Select("`id`, `batch_id`"), repository.Filter("`module_topic_id` = ?", topic.ID))
// 			if err != nil {
// 				return err
// 			}

// 			// delete session
// 			err = service.Repository.UpdateWithMap(uow, batch.BatchTopic{}, map[string]interface{}{
// 				"DeletedBy": course.DeletedBy,
// 				"DeletedAt": time.Now(),
// 			}, repository.Filter("`id`=?", tempBatchTopic.ID))
// 			if err != nil {
// 				uow.RollBack()
// 				return err
// 			}

// 			// delete specified session's feedback (talent and faculty feedback)
// 			err = service.Repository.UpdateWithMap(uow, batch.TalentBatchSessionFeedback{}, map[string]interface{}{
// 				"DeletedBy": course.DeletedBy,
// 				"DeletedAt": time.Now(),
// 			}, repository.Filter("`batch_id`=? AND `batch_topic_id`=?", tempBatchTopic.BatchID, tempBatchTopic.ID))
// 			if err != nil {
// 				uow.RollBack()
// 				return err
// 			}

// 			err = service.Repository.UpdateWithMap(uow, batch.FacultyTalentBatchSessionFeedback{}, map[string]interface{}{
// 				"DeletedBy": course.DeletedBy,
// 				"DeletedAt": time.Now(),
// 			}, repository.Filter("`batch_id`=? AND `batch_topic_id`=?", tempBatchTopic.BatchID, tempBatchTopic.ID))
// 			if err != nil {
// 				uow.RollBack()
// 				return err
// 			}

// 			// delete specified batch's feedback (talent and faculty feedback)
// 			err = service.Repository.UpdateWithMap(uow, batch.TalentFeedback{}, map[string]interface{}{
// 				"DeletedBy": course.DeletedBy,
// 				"DeletedAt": time.Now(),
// 			}, repository.Filter("`batch_id`=?", tempBatchTopic.BatchID, tempBatchTopic.ID))
// 			if err != nil {
// 				uow.RollBack()
// 				return err
// 			}

// 			err = service.Repository.UpdateWithMap(uow, batch.FacultyTalentFeedback{}, map[string]interface{}{
// 				"DeletedBy": course.DeletedBy,
// 				"DeletedAt": time.Now(),
// 			}, repository.Filter("`batch_id`=?", tempBatchTopic.BatchID, tempBatchTopic.ID))
// 			if err != nil {
// 				uow.RollBack()
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// doesTenantExists validates tenantID
func (service *CourseService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCourseExists validates courseID
func (service *CourseService) doesCourseExists(tenantID, courseID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, crs.Course{}, repository.Filter("`id` = ?", courseID))
	if err := util.HandleError("Course not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTalentExist validates talentID
func (service *CourseService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tal.Talent{}, repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Talent not found", exists, err); err != nil {
		return err
	}
	return nil
}

// addCourseSearchQueriesParams adds all search queries if any when getAll is called
func (service *CourseService) addCourseSearchQueriesParams(requestForm url.Values) []repository.QueryProcessor {

	fmt.Println("==================================================requestForm ->", requestForm)

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	var queryProcessors []repository.QueryProcessor

	if createdAt, ok := requestForm["createdAt"]; ok {
		util.AddToSlice("`created_at`", ">= ?", "AND", createdAt, &columnNames, &conditions, &operators, &values)
	}
	if _, ok := requestForm["courseName"]; ok {
		util.AddToSlice("`name`", "LIKE ?", "AND", "%"+requestForm.Get("courseName")+"%", &columnNames, &conditions, &operators, &values)
	}
	if courseType, ok := requestForm["courseType"]; ok {
		util.AddToSlice("`course_type`", "LIKE ?", "AND", courseType, &columnNames, &conditions, &operators, &values)
	}
	if courseLevel, ok := requestForm["courseLevel"]; ok {
		util.AddToSlice("`course_level`", "LIKE ?", "AND", courseLevel, &columnNames, &conditions, &operators, &values)
	}
	//if technologies is present then join courses and technolgies table
	if technologies, ok := requestForm["technologies"]; ok {
		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN courses_technologies ON courses.`id` = courses_technologies.`course_id`"))
		if len(technologies) > 0 {
			util.AddToSlice("courses_technologies.`technology_id`", "IN (?)", "AND", technologies, &columnNames, &conditions, &operators, &values)
		}
	}

	queryProcessors = append(queryProcessors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("courses.`id`"))

	return queryProcessors
}

// // GetSearchedCourses returns all the searched courses from DB.
// func (ser *CourseService) GetSearchedCourses(courses *[]crs.Course, search *crs.Search, tenantID uuid.UUID, limit, offset int, totalCount *int) error {

// 	// Check if tenantID exists
// 	err := ser.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	var queryProcessors []repository.QueryProcessor
// 	queryProcessors = append(queryProcessors, ser.addCourseSearchQueries(search)...)
// 	queryProcessors = append(queryProcessors, repository.PreloadAssociations(ser.associations),
// 		repository.Paginate(limit, offset, totalCount))

// 	uow := repository.NewUnitOfWork(ser.DB, true)
// 	err = ser.Repository.GetAllInOrder(uow, courses, "`name`", queryProcessors...)
// 	if err != nil {
// 		return err
// 	}

// 	for index := range *courses {
// 		var totalCount int
// 		err = ser.Repository.GetCountForTenant(uow, (*courses)[index].TenantID, &crs.Session{},
// 			&totalCount, repository.Filter("`session_id` IS NULL"),
// 			repository.Filter("`course_id`=?", (*courses)[index].ID))
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 		(*courses)[index].TotalSessions = uint(totalCount)
// 	}
// 	uow.Commit()
// 	return nil
// }

// // doesCredentialExist validates credentialID
// func (ser *CourseService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.Credential{}, repository.Filter("`id` = ?", credentialID))
// 	if err := util.HandleError("Course not found", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// addFacultySearchQueries adds all search queries if any when getAll is called
// func (ser *CourseService) addCourseSearchQueries(courseSearch *crs.Search) []repository.QueryProcessor {
// 	var columnNames []string
// 	var conditions []string
// 	var operators []string
// 	var values []interface{}
// 	var queryProcessors []repository.QueryProcessor

// 	if !util.IsEmpty(courseSearch.CreatedAt) {
// 		util.AddToSlice("created_at", ">= ?", "AND", courseSearch.CreatedAt, &columnNames, &conditions, &operators, &values)
// 	}
// 	if !util.IsEmpty(courseSearch.CourseType) {
// 		util.AddToSlice("course_type", "LIKE ?", "AND", courseSearch.CourseType, &columnNames, &conditions, &operators, &values)
// 	}
// 	if !util.IsEmpty(courseSearch.CourseLevel) {
// 		util.AddToSlice("course_level", "LIKE ?", "AND", courseSearch.CourseLevel, &columnNames, &conditions, &operators, &values)
// 	}
// 	//if technologies is present then join courses and technolgies table
// 	if technologies := courseSearch.Technologies; len(technologies) > 0 {
// 		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN course_technologies ON courses.`id` = course_technologies.`course_id`"))
// 		if len(technologies) > 0 {
// 			util.AddToSlice("course_technologies.`technology_id`", "IN(?)", "AND", technologies, &columnNames, &conditions, &operators, &values)
// 		}
// 	}
// 	// if technologies := courseSearch.Technologies; len(technologies) > 0 {
// 	// 	var technologyIDs, courseIDs []string
// 	// 	for _, technology := range technologies {
// 	// 		technologyIDs = append(technologyIDs, technology.ID.String())
// 	// 	}
// 	// 	repository.PluckColumn(ser.DB, "course_technologies", "course_id",
// 	// 		&courseIDs, repository.Filter("technology_id IN(?)", technologyIDs))
// 	// 	util.AddToSlice("id", "IN(?)", "AND", courseIDs, &columnNames, &conditions, &operators, &values)
// 	// }

// 	return queryProcessors
// }
