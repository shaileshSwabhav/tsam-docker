package service

import (
	"fmt"
	"net/url"
	"runtime"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ModuleService Provide method to Update, Delete, Add, Get Method For module.
type ModuleService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewModuleService returns new instance of ModuleService.
func NewModuleService(db *gorm.DB, repository repository.Repository) *ModuleService {
	return &ModuleService{
		DB:         db,
		Repository: repository,
	}
}

// AddModule will add new module to the modules table.
func (service *ModuleService) AddModule(module *course.Module) error {

	// check if foreign keys exist.
	err := service.doForeignKeysExist(module, module.CreatedBy, false)
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

// UpdateModule will update specified module in the modules table.
func (service *ModuleService) UpdateModule(module *course.Module) error {

	// check if foreign keys exist.
	err := service.doForeignKeysExist(module, module.UpdatedBy, false)
	if err != nil {
		return err
	}

	// check if module exist.
	err = service.doesModuleExist(module.TenantID, module.ID)
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

// DeleteModule will delete specified module in the modules table.
func (service *ModuleService) DeleteModule(module *course.Module) error {

	// check if foreign keys exist.
	err := service.doForeignKeysExist(module, module.DeletedBy, true)
	if err != nil {
		return err
	}

	// check if module exist.
	err = service.doesModuleExist(module.TenantID, module.ID)
	if err != nil {
		return err
	}

	exist, err := repository.DoesRecordExistForTenant(service.DB, module.TenantID, course.CourseModule{},
		repository.Filter("`module_id` = ?", module.ID))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("Module cannot be deleted has it is assigned to course.")
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, course.Module{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": module.DeletedBy,
	}, repository.Filter("`id` = ? ", module.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.UpdateWithMap(uow, programming.ModuleProgrammingConcepts{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": module.DeletedBy,
	}, repository.Filter("`module_id` = ?", module.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllModules will get all modules.
func (service *ModuleService) GetAllModules(tenantID uuid.UUID, modules *[]course.ModuleDTO,
	parser *web.Parser, totalCount *int) error {

	fmt.Println(" ============================== parser.Form.Get(isModuleCount) ->", parser.Form.Get("isModuleCount"))

	if !util.IsEmpty(parser.Form.Get("isModuleCount")) && parser.Form.Get("isModuleCount") == "1" {
		err := service.getModules(tenantID, modules, parser, totalCount)
		if err != nil {
			return err
		}

		return nil
	}

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	limit, offset := parser.ParseLimitAndOffset()

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, modules, "`module_name`",
		service.addSearchQueries(parser.Form), repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "ModuleTopics",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("module_topics.`topic_id` IS NULL"), repository.OrderBy("module_topics.`order`"),
				repository.PreloadAssociations([]string{
					"TopicProgrammingConcept", "TopicProgrammingConcept.ProgrammingConcept",
					"TopicProgrammingQuestions", "TopicProgrammingQuestions.ProgrammingQuestion",
				}),
			},
		}, repository.Preload{
			Schema: "ModuleTopics.SubTopics",
			Queryprocessors: []repository.QueryProcessor{
				repository.OrderBy("module_topics.`order`"),
			},
		}),
		repository.PreloadAssociations([]string{"Resources"}), repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

type parameters struct {
	tenantID uuid.UUID
	moduleID uuid.UUID
	channel  chan error
	wg       *sync.WaitGroup
}

// getModules will get all modules.
func (service *ModuleService) getModules(tenantID uuid.UUID, modules *[]course.ModuleDTO,
	parser *web.Parser, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	limit, offset := parser.ParseLimitAndOffset()

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, modules, "`module_name`",
		service.addSearchQueries(parser.Form), repository.PreloadAssociations([]string{"Resources"}),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	var wg sync.WaitGroup

	params := parameters{
		tenantID: tenantID,
		wg:       &wg,
	}

	for index := range *modules {
		channel := make(chan error, 3)
		params.moduleID = (*modules)[index].ID
		params.channel = channel

		params.wg.Add(3)
		go service.getTotalTopics(uow, params, &(*modules)[index].TotalModuleTopics)
		go service.getTotalSubTopics(uow, params, &(*modules)[index].TotalSubTopics)
		go service.getTotalProgrammingQuestions(uow, params, &(*modules)[index].TotalProgrammingQuestions)

		go func() {
			defer close(params.channel)
			params.wg.Wait()
		}()

		for err := range params.channel {
			if err != nil {
				uow.RollBack()
				return err
			}
		}

	}

	defer func() {
		// fmt.Println(" ===================== duration ->", time.Since(now))
		fmt.Println("Go routine Number is --------------------------", runtime.NumGoroutine())
	}()

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

func (service *ModuleService) getTotalTopics(uow *repository.UnitOfWork, param parameters,
	totalCount *int) {

	defer func() {
		param.wg.Done()
	}()

	err := service.Repository.GetCountForTenant(uow, param.tenantID, &course.ModuleTopic{}, totalCount,
		repository.Filter("`module_id` = ? AND `topic_id` IS NULL", param.moduleID))
	if err != nil {
		param.channel <- err

		// return err
	}

	param.channel <- nil
	// return nil
}

func (service *ModuleService) getTotalSubTopics(uow *repository.UnitOfWork, param parameters,
	totalCount *int) {

	defer func() {
		param.wg.Done()
	}()
	// defer param.wg.Done()

	err := service.Repository.GetCountForTenant(uow, param.tenantID, &course.ModuleTopic{}, totalCount,
		repository.Filter("`module_id` = ? AND `topic_id` IS NOT NULL", param.moduleID))
	if err != nil {
		param.channel <- err
		// return err
	}

	// param.wg.Done()
	param.channel <- nil
	// return nil
}

func (service *ModuleService) getTotalProgrammingQuestions(uow *repository.UnitOfWork, param parameters,
	totalCount *int) {

	defer func() {
		param.wg.Done()
	}()

	// defer param.wg.Done()

	err := service.Repository.GetCount(uow, &course.TopicProgrammingQuestion{}, totalCount,
		repository.Join("INNER JOIN module_topics ON topic_programming_questions.`topic_id` = module_topics.`id` AND "+
			"module_topics.`tenant_id` = topic_programming_questions.`tenant_id`"),
		repository.Filter("module_topics.`id` = ? AND topic_programming_questions.`tenant_id` = ? AND"+
			" module_topics.`deleted_at` IS NULL", param.moduleID, param.tenantID))
	if err != nil {
		param.channel <- err
		// return err
	}

	// param.wg.Done()
	param.channel <- nil
	// return nil
}

// doForeignKeysExist checks if all foregin keys exist and if not returns error.
func (service *ModuleService) doForeignKeysExist(module *course.Module, credentialID uuid.UUID,
	isDeleteCall bool) error {

	// Check if tenant exists.
	err := service.doesTenantExist(module.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(module.TenantID, credentialID)
	if err != nil {
		return err
	}

	if !isDeleteCall {
		// check if module with same name exist
		err = service.doesModuleNameExist(module.TenantID, module.ID, module.ModuleName)
		if err != nil {
			return err
		}
	}

	return nil
}

// addSearchQueries adds search criteria.
func (service *ModuleService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["moduleName"]; ok {
		util.AddToSlice("`module_name`", "LIKE ?", "AND", "%"+requestForm.Get("moduleName")+"%", &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *ModuleService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *ModuleService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesModuleNameExist returns error if there is no module record with same name exist in table for the given tenant.
func (service *ModuleService) doesModuleNameExist(tenantID, moduleID uuid.UUID, moduleName string) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Module{},
		repository.Filter("`module_name`=? AND `id`!=?", moduleName, moduleID))
	if err := util.HandleIfExistsError("Module name already exists", exists, err); err != nil {
		return err
	}
	return nil
}

// doesModuleExist returns error if there is no module record in table for the given tenant.
func (service *ModuleService) doesModuleExist(tenantID, moduleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Module{},
		repository.Filter("`id`=?", moduleID))
	if err := util.HandleError("Invalid module ID", exists, err); err != nil {
		return err
	}
	return nil
}
