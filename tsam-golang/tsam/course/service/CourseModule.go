package service

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CourseModuleService Provide method to Update, Delete, Add, Get Method For Course.
type CourseModuleService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewCourseModuleService returns new instance of CourseModuleService.
func NewCourseModuleService(db *gorm.DB, repository repository.Repository) *CourseModuleService {
	return &CourseModuleService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Course", "Module",
			//  "Module.ModuleTopics",
		},
	}
}

// AddCourseModule will add new course_module to the table.
func (service *CourseModuleService) AddCourseModule(module *course.CourseModule) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(module, module.CreatedBy, false)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, module)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateCourseModule will update course_module for specified course.
func (service *CourseModuleService) UpdateCourseModule(module *course.CourseModule) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(module, module.UpdatedBy, false)
	if err != nil {
		return err
	}

	// check if course programming assignment exist.
	err = service.doesCourseModuleExist(module.TenantID, module.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, module)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteCourseModule will delete course_module for specified course.
func (service *CourseModuleService) DeleteCourseModule(module *course.CourseModule) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(module, module.DeletedBy, true)
	if err != nil {
		return err
	}

	// check if course programming assignment exist.
	err = service.doesCourseModuleExist(module.TenantID, module.ID)
	if err != nil {
		return err
	}

	// exist, err := repository.DoesRecordExistForTenant(service.DB, module.TenantID, course.CourseSession{},
	// 	repository.Filter("`module_id` = ?", module.ID))
	// if err != nil {
	// 	return err
	// }
	// if exist {
	// 	return errors.NewValidationError("Module cannot be deleted as sessions are assigned.")
	// }

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, &course.CourseModule{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": module.DeletedBy,
	}, repository.Filter("`id` = ?", module.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetCourseModule will get all course modules.
func (service *CourseModuleService) GetCourseModule(courseModules *[]course.CourseModuleDTO, tenantID,
	courseID uuid.UUID, parser *web.Parser, totalCount *int) error {

	now := time.Now()

	defer func() {
		fmt.Println("=================== duration ->", time.Since(now))
	}()

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, courseModules, "`order`",
		repository.Filter("`course_id` = ?", courseID), service.addSearchQueries(parser.Form),
		repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *courseModules {
		err = service.getCourseModuleTopic(uow, tenantID, courseID, (*courseModules)[index].ModuleID,
			&(*courseModules)[index].Module.ModuleTopics)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doesForeignKeyExist will check if all foreign keys are valid.
func (service *CourseModuleService) doesForeignKeyExist(module *course.CourseModule,
	credentialID uuid.UUID, isDeleteOperation bool) error {

	// check if tenant exist.
	err := service.doesTenantExist(module.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(module.TenantID, credentialID)
	if err != nil {
		return err
	}

	if !isDeleteOperation {

		// check if course exist.
		err = service.doesCourseExist(module.TenantID, module.CourseID)
		if err != nil {
			return err
		}

		// check if module exist.
		err = service.doesModuleExist(module.TenantID, module.ModuleID)
		if err != nil {
			return err
		}

		// check if order already exist.
		err = service.doesCourseModuleOrderExist(module.TenantID, module.ID, module.CourseID, module.Order)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *CourseModuleService) getCourseModuleTopic(uow *repository.UnitOfWork, tenantID, courseID,
	moduleID uuid.UUID, moduleTopics *[]*course.ModuleTopicDTO) error {

	err := service.Repository.GetAllInOrderForTenant(uow, tenantID, moduleTopics, "`order`",
		repository.Filter("`topic_id` IS NULL AND `module_id` = ?", moduleID),
		repository.PreloadAssociations([]string{
			// "Resources",
			"TopicProgrammingConcept",
			"TopicProgrammingQuestions", "TopicProgrammingQuestions.ProgrammingQuestion",
		}),
		// repository.PreloadAssociations([]string{"Resources"}),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "SubTopics",
			Queryprocessors: []repository.QueryProcessor{
				repository.OrderBy("`order`"),
			},
			// Queryprocessors: []repository.QueryProcessor{
			// 	repository.OrderBy("`order`"), repository.PreloadAssociations([]string{
			// 		// "Resources",
			// 		"TopicProgrammingConcept",
			// 		"TopicProgrammingQuestions", "TopicProgrammingQuestions.ProgrammingQuestion",
			// 	}),
			// },
		}))
	if err != nil {
		return err

	}
	return nil
}

// addSearchQueries adds all search queries.
func (service *CourseModuleService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	// var queryProcessors []repository.QueryProcessor

	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("`course_modules`.`is_active`", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	// if moduleName, ok := requestForm["moduleName"]; ok {
	// 	util.AddToSlice("`course_modules`.`module_name`", "= ?", "AND", moduleName, &columnNames, &conditions, &operators, &values)
	// }

	// if courseID, ok := requestForm["courseID"]; ok {
	// 	util.AddToSlice("`course_modules`.`course_id`", "= ?", "AND", courseID, &columnNames, &conditions, &operators, &values)
	// }

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// returns error if there is no tenant record in table.
func (service *CourseModuleService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *CourseModuleService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no course record in table for the given tenant.
func (service *CourseModuleService) doesCourseExist(tenantID, courseID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Course{},
		repository.Filter("`id` = ?", courseID))
	if err := util.HandleError("Invalid course ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no module record in table for the given tenant.
func (service *CourseModuleService) doesModuleExist(tenantID, moduleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Module{},
		repository.Filter("`id` = ?", moduleID))
	if err := util.HandleError("Invalid module ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no course module record in table for the given tenant.
func (service *CourseModuleService) doesCourseModuleOrderExist(tenantID, moduleID, courseID uuid.UUID, order uint) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.CourseModule{},
		repository.Filter("`id` != ? AND `course_id` = ? AND `order` = ? AND `is_active` = ?",
			moduleID, courseID, order, true))
	if err := util.HandleIfExistsError("Order already exist", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no course modules record in table for the given tenant.
func (service *CourseModuleService) doesCourseModuleExist(tenantID, courseModuleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.CourseModule{},
		repository.Filter("`id` = ?", courseModuleID))
	if err := util.HandleError("Invalid course module ID", exists, err); err != nil {
		return err
	}
	return nil
}
