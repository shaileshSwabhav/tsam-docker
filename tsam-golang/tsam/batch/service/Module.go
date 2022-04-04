package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

type ModuleService struct {
	db   *gorm.DB
	repo repository.Repository
}

func NewModuleService(db *gorm.DB, repo repository.Repository) *ModuleService {
	return &ModuleService{
		db:   db,
		repo: repo,
	}
}

// AddModule will add module for specified batch.
func (service *ModuleService) AddModule(batchModule *batch.Module, uows ...*repository.UnitOfWork) error {

	// check if all foreign key's are valid.
	err := service.doesForeignKeyExist(batchModule, batchModule.CreatedBy)
	if err != nil {
		return err
	}

	// check if module exist for same faculty and batch
	err = service.doesDuplicateModuleExit(batchModule)
	if err != nil {
		return err
	}

	if batchModule.ModuleTiming != nil {
		assignIDToModuleTiming(batchModule, true)
	}

	var uow *repository.UnitOfWork
	if len(uows) == 0 {
		uow = repository.NewUnitOfWork(service.db, false)
	} else {
		uow = uows[0]
	}

	err = service.repo.Add(uow, batchModule)
	if err != nil {
		if len(uows) == 0 {
			uow.RollBack()
		}
		return err
	}

	if len(uows) == 0 {
		uow.Commit()
	}
	return nil
}

// AddModules will add multiple modules to specified batch.
func (service *ModuleService) AddModules(batchModules *[]batch.Module, tenantID,
	credentialID, batchID uuid.UUID) error {

	err := service.checkModuleOrderInJSON(batchModules)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	for _, module := range *batchModules {
		module.BatchID = batchID
		module.TenantID = tenantID
		module.CreatedBy = credentialID

		err := service.AddModule(&module, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// UpdateModule will update specified batch_module.
func (service *ModuleService) UpdateModule(batchModule *batch.Module) error {

	// check if all foreign key's are valid.
	err := service.doesForeignKeyExist(batchModule, batchModule.UpdatedBy)
	if err != nil {
		return err
	}

	// check module exist
	err = service.doesBatchModuleExists(batchModule.TenantID, batchModule.ID)
	if err != nil {
		return err
	}

	if batchModule.ModuleTiming != nil {
		assignIDToModuleTiming(batchModule, false)
	}

	uow := repository.NewUnitOfWork(service.db, false)

	err = service.updateModuleTiming(uow, batchModule, batchModule.UpdatedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	tempModule := batch.Module{}

	err = service.repo.GetRecordForTenant(uow, batchModule.TenantID, &tempModule,
		repository.Filter("`id` = ?", batchModule.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	batchModule.CreatedBy = tempModule.CreatedBy

	// err = service.repo.Update(uow, module)
	err = service.repo.Save(uow, batchModule)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteModule will delete specified batch_module.
func (service *ModuleService) DeleteModule(batchModule *batch.Module) error {

	// check tenant exist
	err := service.doesTenantExists(batchModule.TenantID)
	if err != nil {
		return err
	}

	// check credential exist
	err = service.doesCredentialExist(batchModule.TenantID, batchModule.DeletedBy)
	if err != nil {
		return err
	}

	// check module exist
	err = service.doesBatchModuleExists(batchModule.TenantID, batchModule.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	err = service.repo.UpdateWithMap(uow, batch.Module{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": batchModule.DeletedBy,
	}, repository.Filter("`id` = ?", batchModule.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBatchModules will get all the modules for specified batch.
func (service *ModuleService) GetBatchModules(batchModules *[]batch.ModuleDTO, tenantID, batchID uuid.UUID,
	totalCount *int, parser *web.Parser) error {

	// check tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// check batch exist
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	limit, offset := parser.ParseLimitAndOffset()

	var queryProcessor []repository.QueryProcessor

	queryProcessor = append(queryProcessor, repository.Filter("`batch_id` = ?", batchID), service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{"Module", "Faculty"}),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "ModuleTiming",
			Queryprocessors: []repository.QueryProcessor{
				repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND" +
					" days.`tenant_id` = batch_module_timings.`tenant_id`"), repository.Filter("batch_module_timings.`tenant_id` = ?",
					tenantID), repository.Filter("batch_module_timings.`deleted_at` IS NULL"),
				repository.OrderBy("days.`order`"), repository.PreloadAssociations([]string{"Day"}),
			}}), repository.Paginate(limit, offset, totalCount))

	err = service.repo.GetAllInOrderForTenant(uow, tenantID, batchModules, "`order`", queryProcessor...)
	if err != nil {
		uow.RollBack()
		return err
	}

	if !util.IsEmpty(parser.Form.Get("field")) {
		err = service.getAdditionalFields(uow, batchModules, tenantID, parser)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

func (service *ModuleService) getAdditionalFields(uow *repository.UnitOfWork, batchModules *[]batch.ModuleDTO,
	tenantID uuid.UUID, parser *web.Parser) error {

	for index := range *batchModules {

		if parser.Form["field"][0] == "ModuleTopics" {
			err := service.getModuleTopics(uow, &(*batchModules)[index].Module.ModuleTopics, tenantID, (*batchModules)[index].Module.ID)
			if err != nil {
				return err
			}
		}

		if parser.Form["field"][1] == "TopicProgrammingQuestions" {
			for k := range (*batchModules)[index].Module.ModuleTopics {
				err := service.getTopicProgrammingQuestion(uow, &(*batchModules)[index].Module.ModuleTopics[k].TopicProgrammingQuestions,
					tenantID, (*batchModules)[index].Module.ModuleTopics[k].ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (service *ModuleService) getModuleTopics(uow *repository.UnitOfWork, moduleTopics *[]*course.ModuleTopicDTO,
	tenantID, moduleID uuid.UUID) error {

	err := service.repo.GetAll(uow, moduleTopics,
		repository.Join("INNER JOIN batch_session_topics ON batch_session_topics.`topic_id` = module_topics.`id` AND"+
			" batch_session_topics.`tenant_id` = module_topics.`tenant_id`"),
		repository.Filter("batch_session_topics.`deleted_at` IS NULL AND module_topics.`topic_id` IS NULL AND"+
			" module_topics.`tenant_id` = ? AND module_topics.`module_id` = ?", tenantID, moduleID),
		repository.PreloadAssociations([]string{"BatchTopicAssignment", "BatchTopicAssignment.ProgrammingQuestion"}),
		repository.GroupBy("module_topics.`id`"), repository.OrderBy("module_topics.`order`"))
	if err != nil {
		return err
	}
	return nil
}

func (service *ModuleService) getTopicProgrammingQuestion(uow *repository.UnitOfWork,
	topicProgrammingQuestions *[]*course.TopicProgrammingQuestionDTO, tenantID, topicID uuid.UUID) error {

	err := service.repo.GetAll(uow, topicProgrammingQuestions, repository.Join("INNER JOIN batch_topic_assignments ON"+
		" batch_topic_assignments.`programming_question_id` = topic_programming_questions.`programming_question_id` AND"+
		" batch_topic_assignments.`tenant_id` = topic_programming_questions.`tenant_id`"),
		repository.Filter("batch_topic_assignments.`deleted_at` IS NULL AND topic_programming_questions.`deleted_at` IS NULL"+
			" AND topic_programming_questions.`tenant_id` = ?", tenantID),
		repository.Filter("batch_topic_assignments.`topic_id` = ? AND topic_programming_questions.`is_active` = ?",
			topicID, 1), repository.PreloadAssociations([]string{"ProgrammingQuestion"}),
		repository.GroupBy("batch_topic_assignments.`programming_question_id`"))
	if err != nil {
		return err
	}

	return nil
}

// GetBatchModulesWithAllFields will get all the modules with all the preloads.
func (service *ModuleService) GetBatchModulesWithAllFields(batchModules *[]batch.ModuleDTO, tenantID, batchID uuid.UUID,
	totalCount *int, parser *web.Parser) error {

	// Check tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check batch exist.
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	limit, offset := parser.ParseLimitAndOffset()

	err = service.repo.GetAllInOrderForTenant(uow, tenantID, batchModules, "`order`",
		repository.Filter("`batch_id` = ?", batchID), service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{"Faculty", "Module"}),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "Module.ModuleTopics",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("`topic_id` IS NULL"), repository.OrderBy("module_topics.`order`")}}),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "Module.ModuleTopics.SubTopics",
			Queryprocessors: []repository.QueryProcessor{repository.OrderBy("module_topics.`order`")}}),
		repository.PreloadAssociations([]string{
			"Module.ModuleTopics.TopicProgrammingConcept",
		}),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "Module.ModuleTopics.TopicProgrammingQuestions",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("topic_programming_questions.`deleted_at` IS NULL")}}),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "Module.ModuleTopics.TopicProgrammingQuestions.ProgrammingQuestion",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("programming_questions.`deleted_at` IS NULL")}}),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "ModuleTiming",
			Queryprocessors: []repository.QueryProcessor{
				repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND" +
					" days.`tenant_id` = batch_module_timings.`tenant_id`"),
				repository.Filter("batch_module_timings.`tenant_id` = ?",
					tenantID),
				repository.Filter("batch_module_timings.`deleted_at` IS NULL"),
				repository.OrderBy("days.`order`"),
				repository.PreloadAssociations([]string{"Day"})}}),
		repository.Paginate(limit, offset, totalCount))
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

func (service *ModuleService) checkModuleOrderInJSON(batchModules *[]batch.Module) error {

	orderMap := make(map[uint]uint)

	for _, module := range *batchModules {
		if orderMap[module.Order] == 1 {
			return errors.NewValidationError("Multiple modules cannot have same order")
		}
		orderMap[module.Order] = 1
	}

	return nil
}

func assignIDToModuleTiming(module *batch.Module, isAddOperation bool) {
	for index := range module.ModuleTiming {
		module.ModuleTiming[index].TenantID = module.TenantID
		module.ModuleTiming[index].BatchID = module.BatchID
		module.ModuleTiming[index].ModuleID = module.ModuleID
		module.ModuleTiming[index].FacultyID = module.FacultyID

		if isAddOperation {
			module.ModuleTiming[index].CreatedBy = module.CreatedBy
			continue
		}
		module.ModuleTiming[index].UpdatedBy = module.UpdatedBy
	}
}

func (service *ModuleService) updateModuleTiming(uow *repository.UnitOfWork, batchModule *batch.Module,
	credentialID uuid.UUID) error {

	moduleTimingMap := make(map[uuid.UUID]uint)
	moduleTimings := batchModule.ModuleTiming
	tempModuleTimings := []batch.ModuleTiming{}

	err := service.repo.GetAllForTenant(uow, batchModule.TenantID, &tempModuleTimings,
		repository.Filter("`batch_module_id` = ? AND `faculty_id` = ?", batchModule.ID, batchModule.FacultyID))
	if err != nil {
		return err
	}

	// populate all entries for existing batch timing (existing)
	for _, tempModuleTime := range tempModuleTimings {
		moduleTimingMap[tempModuleTime.ID] = moduleTimingMap[tempModuleTime.ID] + 1
	}

	for i := 0; i < len(moduleTimings); i++ {

		if util.IsUUIDValid(moduleTimings[i].ID) {
			moduleTimingMap[moduleTimings[i].ID] = moduleTimingMap[moduleTimings[i].ID] + 1

			// update existing records
			if moduleTimingMap[moduleTimings[i].ID] > 1 {
				moduleTimings[i].UpdatedBy = credentialID
				moduleTimings[i].ModuleID = batchModule.ModuleID
				moduleTimings[i].BatchModuleID = batchModule.ID
				err = service.repo.Update(uow, &moduleTimings[i])
				if err != nil {
					return err
				}
				moduleTimingMap[moduleTimings[i].ID] = 0
			}
			continue
		}

		// add new module-timing
		moduleTimings[i].ModuleID = batchModule.ModuleID
		moduleTimings[i].CreatedBy = credentialID
		moduleTimings[i].BatchID = batchModule.BatchID
		moduleTimings[i].TenantID = batchModule.TenantID
		moduleTimings[i].BatchModuleID = batchModule.ID
		err = service.repo.Add(uow, &moduleTimings[i])
		if err != nil {
			return err
		}
		// (*moduleTimings) = append((*moduleTimings)[:i], (*moduleTimings)[i+1:]...)
		// i = i - 1
	}

	for _, tempModuleTime := range tempModuleTimings {
		if moduleTimingMap[tempModuleTime.ID] == 1 {
			err = service.repo.UpdateWithMap(uow, batch.ModuleTiming{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempModuleTime.ID))
			if err != nil {
				return err
			}
		}
	}

	// module.ModuleTiming = nil
	return nil
}

// addSearchQueries adds search criteria.
func (service *ModuleService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if moduleID, ok := requestForm["moduleID"]; ok {
		util.AddToSlice("`module_id`", "= ?", "AND", moduleID, &columnNames, &conditions, &operators, &values)
	}

	if facultyID, ok := requestForm["facultyID"]; ok {
		util.AddToSlice("`faculty_id`", "= ?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
	}

	if batchID, ok := requestForm["batchID"]; ok {
		util.AddToSlice("`batch_id`", "= ?", "AND", batchID, &columnNames, &conditions, &operators, &values)
	}

	if order, ok := requestForm["order"]; ok {
		util.AddToSlice("`order`", "= ?", "AND", order, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesForeignKeyExist checks if all foreign keys are valid.
func (service *ModuleService) doesForeignKeyExist(batchModule *batch.Module, credentialID uuid.UUID) error {

	// check tenant exist
	err := service.doesTenantExists(batchModule.TenantID)
	if err != nil {
		return err
	}

	// check credential exist
	err = service.doesCredentialExist(batchModule.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check batch exist
	err = service.doesBatchExists(batchModule.TenantID, batchModule.BatchID)
	if err != nil {
		return err
	}

	// check if faculty exist
	err = service.doesFacultyExists(batchModule.TenantID, batchModule.FacultyID)
	if err != nil {
		return err
	}

	// check module exist
	err = service.doesModuleExists(batchModule.TenantID, batchModule.ModuleID)
	if err != nil {
		return err
	}

	// check if order exist
	err = service.doesOrderExistForBatch(batchModule)
	if err != nil {
		return err
	}

	// // check if start_date for multiple modules in same for single faculty
	// err = service.doeStartDateExistForBatch(module)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// doesTenantExists validates tenantID
func (service *ModuleService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.db, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *ModuleService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesBatchExists validates batchID
func (service *ModuleService) doesBatchExists(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Batch not found", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *ModuleService) doesFacultyExists(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, faculty.Faculty{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Faculty not found", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *ModuleService) doesModuleExists(tenantID, moduleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, course.Module{},
		repository.Filter("`id` = ?", moduleID))
	if err := util.HandleError("Module not found", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *ModuleService) doesBatchModuleExists(tenantID, batchModuleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.Module{},
		repository.Filter("`id` = ?", batchModuleID))
	if err := util.HandleError("Module for batch not found", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *ModuleService) doesDuplicateModuleExit(batchModule *batch.Module) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, batchModule.TenantID, batch.Module{},
		repository.Filter("`module_id` = ? AND `batch_id` = ?", batchModule.ModuleID, batchModule.BatchID))
	if err := util.HandleIfExistsError("Module for already added to the batch", exists, err); err != nil {
		return err
	}
	return nil
}

// func (service *ModuleService) doeStartDateExistForBatch(module *batch.Module) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.db, module.TenantID, batch.Module{},
// 		repository.Filter("`batch_id` = ? AND `start_date` = ? AND `module_id` = ? AND `faculty_id` = ? AND `id` != ?",
// 			module.BatchID, module.StartDate, module.ModuleID, module.FacultyID, module.ID))
// 	if err := util.HandleIfExistsError("Same faculty cannot have multiple modules on same date", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (service *ModuleService) doesOrderExistForBatch(module *batch.Module) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, module.TenantID, batch.Module{},
		repository.Filter("`batch_id` = ? AND `order` = ? AND `id` != ?",
			module.BatchID, module.Order, module.ID))
	if err := util.HandleIfExistsError("Order already exist for specified module in batch", exists, err); err != nil {
		return err
	}
	return nil
}
