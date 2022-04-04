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
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ModuleProgrammingConceptService provides methods to update, delete, add, get for programming concept module.
type ModuleProgrammingConceptService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewModuleProgrammingConceptService returns new instance of ModuleProgrammingConceptService.
func NewModuleProgrammingConceptService(db *gorm.DB, repository repository.Repository) *ModuleProgrammingConceptService {
	return &ModuleProgrammingConceptService{
		DB:         db,
		Repository: repository,
		// association: []string{
		// 	"ProgrammingQuestions",
		// },
	}
}

// AddModuleProgrammingConcepts will add new programming concept modules to the table.
func (service *ModuleProgrammingConceptService) AddModuleProgrammingConcepts(tenantID, credentialID uuid.UUID,
	conceptModules *[]programming.ModuleProgrammingConcepts) error {

	// Checks if all foreign keys exist.
	err := service.doForeignKeysExist(tenantID, credentialID, conceptModules)
	if err != nil {
		return err
	}

	// Check if concept module exists.
	for i := range *conceptModules {
		if err := service.doesConceptModuleExist(tenantID, &(*conceptModules)[i]); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Set createdBy field and tenant id to all concept modules.
	if conceptModules != nil && len(*conceptModules) != 0 {
		for i := 0; i < len(*conceptModules); i++ {
			(*conceptModules)[i].CreatedBy = credentialID
			(*conceptModules)[i].TenantID = tenantID
		}
	}

	// Check if all concept modules are having unique concept ids.
	conceptModuleMap := make(map[uuid.UUID]uint)
	for _, conceptModule := range *conceptModules {
		conceptModuleMap[conceptModule.ProgrammingConceptID]++
		if conceptModuleMap[conceptModule.ProgrammingConceptID] > 1 {
			return errors.NewValidationError("Modules cannot have same concepts")
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add all programming concept modules.
	for i := range *conceptModules {

		// If there are parent concept modules get id of its parent concept modules.
		if len((*conceptModules)[i].ParentModuleProgrammingConcepts) > 0 {
			if err := service.getParentConceptModuleID(uow, tenantID, &(*conceptModules)[i]); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}
		}

		// Add concept module.
		err = service.Repository.Add(uow, &(*conceptModules)[i])
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// UpdateModuleProgrammingConcepts will update the specified programming concept module record in the table.
func (service *ModuleProgrammingConceptService) UpdateModuleProgrammingConcepts(moduleID, tenantID, credentialID uuid.UUID,
	conceptModules *[]programming.ModuleProgrammingConcepts) error {

	// Checks if all foreign keys exist.
	err := service.doForeignKeysExist(tenantID, credentialID, conceptModules)
	if err != nil {
		return err
	}

	// Check if all concept modules are having unique concept ids.
	conceptModuleMap := make(map[uuid.UUID]uint)
	for _, conceptModule := range *conceptModules {
		conceptModuleMap[conceptModule.ProgrammingConceptID]++
		if conceptModuleMap[conceptModule.ProgrammingConceptID] > 1 {
			return errors.NewValidationError("Modules cannot have same concepts")
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// If previous concept modules is present and current concept modules is not present.
	if conceptModules == nil {
		err := service.Repository.UpdateWithMap(uow, programming.ModuleProgrammingConcepts{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		}, repository.Filter("`module_id`=?"))
		if err != nil {
			return err
		}

		uow.Commit()
		return nil
	}

	// Create temp concept modules to get presvious concept modules.
	tempConceptModules := []programming.ModuleProgrammingConcepts{}

	err = service.Repository.GetAllForTenant(uow, tenantID, &tempConceptModules,
		repository.Filter("`module_id`=?", moduleID))
	if err != nil {
		return err
	}

	// Make map to count occurences of concept module id in previous and current concept modules.
	conceptModulesIDMap := make(map[uuid.UUID]uint)

	// Count the number of occurence of previous concept module ID.
	for _, tempConceptModule := range tempConceptModules {
		conceptModulesIDMap[tempConceptModule.ID]++
	}

	// Count the number of occurrence of current concept module ID to know total count of occurrence of each ID.
	for _, conceptModule := range *conceptModules {

		// If ID is valid then push its occurence in ID map.
		if util.IsUUIDValid(conceptModule.ID) {
			conceptModulesIDMap[conceptModule.ID]++
		} else {

			// If ID is nil create new concept module entry in table.
			conceptModule.CreatedBy = credentialID
			conceptModule.TenantID = tenantID

			// Check if concept module exists.
			if err := service.doesConceptModuleExist(tenantID, &conceptModule); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}

			// If there are parent concept modules get id of its parent concept modules.
			if err := service.getParentConceptModuleID(uow, tenantID, &conceptModule); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}

			err = service.Repository.Add(uow, &conceptModule)
			if err != nil {
				return err
			}
		}

		// If number of occurrence is more than one (present in previous and current concept modules)
		// then update concept module.
		if conceptModulesIDMap[conceptModule.ID] > 1 {

			// If there are parent concept modules get id of its parent concept modules.
			if err := service.getParentConceptModuleID(uow, tenantID, &conceptModule); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}

			conceptModule.UpdatedBy = credentialID
			err = service.Repository.Update(uow, &conceptModule)
			if err != nil {
				return err
			}

			// Make the number of occurrences 0 after updating concept module.
			conceptModulesIDMap[conceptModule.ID] = 0
		}
	}

	// If the number of occurrence is one, the concept module was presnt in previous concept modules only, delete it.
	for _, tempConceptModule := range tempConceptModules {
		if conceptModulesIDMap[tempConceptModule.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, programming.ModuleProgrammingConcepts{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempConceptModule.ID))
			if err != nil {
				return err
			}

			// Make the number of occurrences 0 after deleting concept module.
			conceptModulesIDMap[tempConceptModule.ID] = 0
		}
	}

	uow.Commit()
	return nil
}

// DeleteModuleProgrammingConcept will delete the specified programming concept module record from the table.
func (service *ModuleProgrammingConceptService) DeleteModuleProgrammingConcept(conceptModule *programming.ModuleProgrammingConcepts) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(conceptModule.TenantID); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(conceptModule.TenantID, conceptModule.DeletedBy); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if programming concept module record exist.
	err := service.doesConceptModuleExist(conceptModule.TenantID, conceptModule)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Delete programming concept module.
	err = service.Repository.UpdateWithMap(uow, &programming.ModuleProgrammingConcepts{}, map[interface{}]interface{}{
		"DeletedBy": conceptModule.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id`=?", conceptModule.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllModuleProgrammingConceptsForConceptTree will return all the records from programming concept modules table for concept tree.
func (service *ModuleProgrammingConceptService) GetAllModuleProgrammingConceptsForConceptTree(conceptModules *[]programming.ModuleProgrammingConcepts,
	tenantID uuid.UUID, parser *web.Parser, totalCount *int) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	limit, offset := parser.ParseLimitAndOffset()

	// Get all programming concepts.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, conceptModules, "`level`",
		service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount),
		// repository.PreloadWithCustomCondition(repository.Preload{Schema: "ParentModuleProgrammingConcepts",
		// 	Queryprocessors: []repository.QueryProcessor{
		// 		repository.OrderBy("`parent_module_programming_concept_id`"),
		// 	}}))
		repository.PreloadAssociations([]string{"ParentModuleProgrammingConcepts"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllModuleProgrammingConcepts will return all the records from programming concept modules table.
func (service *ModuleProgrammingConceptService) GetAllModuleProgrammingConcepts(conceptModules *[]programming.ModuleProgrammingConceptsDTO,
	tenantID uuid.UUID, parser *web.Parser, totalCount *int) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	limit, offset := parser.ParseLimitAndOffset()

	// Get all programming concepts.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, conceptModules, "`level`",
		service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount),
		repository.PreloadAssociations([]string{"ProgrammingConcept"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetModuleConceptsForAssignment will return all concept_modules for a specific assignment.
func (service *ModuleProgrammingConceptService) GetModuleConceptsForAssignment(conceptModules *[]programming.ModuleProgrammingConceptsDTO,
	tenantID, assignmentID uuid.UUID, parser *web.Parser) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// try to shift this to search for get all, also add module_id in Batchtopic assignment. #Niranjan
	// Needs a final check.

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	bta := &batch.TopicAssignment{}
	bta.ID = assignmentID
	err = service.Repository.GetRecordForTenant(uow, tenantID, bta,
		repository.Filter("`id` = ?", bta.ID),
		repository.Select("`module_id`"))
	if err != nil {
		return err
	}

	// Get all programming concept modules.
	// err = service.Repository.GetAllInOrder(uow, conceptModules, "`modules_programming_concepts`.`level`",
	// 	service.addSearchQueries(parser.Form),
	// 	repository.PreloadAssociations([]string{"ProgrammingConcept"}),
	// 	repository.Join("INNER JOIN `module_topics` ON "+
	// 		"`modules_programming_concepts`.`module_id` = `module_topics`.`module_id` AND "+
	// 		"`modules_programming_concepts`.`tenant_id` = `module_topics`.`tenant_id`"),
	// 	repository.Join("INNER JOIN `batch_topic_assignments` ON "+
	// 		"`batch_topic_assignments`.`topic_id` = `module_topics`.`id` AND "+
	// 		"`batch_topic_assignments`.`tenant_id` = `module_topics`.`tenant_id`"),
	// 	repository.Join("INNER JOIN `programming_questions_programming_concepts` ON "+
	// 		"`programming_questions_programming_concepts`.`programming_concept_id` = `modules_programming_concepts`.`programming_concept_id`"),
	// 	repository.Filter("`modules_programming_concepts`.`tenant_id` = ?", tenantID),
	// 	repository.Filter("`module_topics`.`deleted_at` IS NULL"),
	// 	repository.Filter("`batch_topic_assignments`.`deleted_at` IS NULL"),
	// 	repository.Filter("`batch_topic_assignments`.`id` = ?", assignmentID))

	// SELECT * FROM `batch_topic_assignments` bta
	// INNER JOIN `programming_questions_programming_concepts` pqc ON
	// bta.`programming_question_id` = pqc.`programming_question_id` INNER JOIN
	// `module_topics`mt ON bta.`topic_id` = mt.`id` INNER JOIN `modules_programming_concepts` pcm ON
	// pcm.`programming_concept_id` = pqc.`programming_concept_id` WHERE
	// bta.`id` = "5163e50b-78d1-4467-b415-d56d0c194b9e"
	// Get all programming concept modules.
	err = service.Repository.GetAllInOrder(uow, conceptModules, "`modules_programming_concepts`.`level`",
		service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{"ProgrammingConcept"}),
		repository.Join("INNER JOIN `programming_questions_programming_concepts` pqpc ON "+
			"`modules_programming_concepts`.`programming_concept_id` = `pqpc`.`programming_concept_id`"),
		repository.Join("INNER JOIN `batch_topic_assignments` ON "+
			"`batch_topic_assignments`.`programming_question_id` = `pqpc`.`programming_question_id` AND "+
			"`modules_programming_concepts`.`tenant_id` = `batch_topic_assignments`.`tenant_id`"),
		repository.Join("INNER JOIN `topic_programming_concepts` ON "+
			"`topic_programming_concepts`.`topic_id` = `batch_topic_assignments`.`topic_id` AND "+
			"`topic_programming_concepts`.`programming_concept_id`=`pqpc`.`programming_concept_id` AND "+
			"`topic_programming_concepts`.`tenant_id` = `batch_topic_assignments`.`tenant_id`"),
		repository.Filter("`modules_programming_concepts`.`tenant_id` = ?", tenantID),
		repository.Filter("`modules_programming_concepts`.`module_id` = ?", bta.ModuleID),
		repository.Filter("`batch_topic_assignments`.`tenant_id` = ?", tenantID),
		repository.Filter("`topic_programming_concepts`.`deleted_at` IS NULL"),
		repository.Filter("`batch_topic_assignments`.`deleted_at` IS NULL"),
		repository.Filter("`batch_topic_assignments`.`id` = ?", assignmentID),
		repository.GroupBy("`topic_programming_concepts`.`programming_concept_id`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetModuleProgrammingConcept will return specified record from programming concept module table.
func (service *ModuleProgrammingConceptService) GetModuleProgrammingConcept(conceptModule *programming.ModuleProgrammingConcepts,
	tenantID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if programming concept record exist.
	err = service.doesConceptModuleExist(tenantID, conceptModule)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetForTenant(uow, tenantID, conceptModule.ID, conceptModule)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllModuleProgrammingConceptsForTalentScore will return all the records from programming concept modules table for talent score.
func (service *ModuleProgrammingConceptService) GetAllModuleProgrammingConceptsForTalentScore(conceptModules *[]programming.ModuleProgrammingConceptsForTalentScore,
	tenantID, talentID uuid.UUID, parser *web.Parser) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all programming concepts.
	err = service.Repository.GetAllInOrder(uow, conceptModules, "`level`",
		repository.Select("modules_programming_concepts.*, AVG(talent_concept_ratings.`score`) AS `average_score`"),
		repository.Join("LEFT JOIN talent_concept_ratings on talent_concept_ratings.`module_programming_concept_id` = modules_programming_concepts.`id`"),
		service.addSearchQueries(parser.Form),
		repository.Filter("modules_programming_concepts.`tenant_id` = ? AND modules_programming_concepts.`deleted_at` IS NULL", tenantID),
		repository.Filter("talent_concept_ratings.`deleted_at` IS NULL"),
		repository.Filter("talent_concept_ratings.`talent_id` = ? OR talent_concept_ratings.`talent_id` IS NULL", talentID),
		repository.GroupBy("modules_programming_concepts.`id`"),
		repository.PreloadAssociations([]string{"ParentModuleProgrammingConcepts", "ProgrammingConcept"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	// err = service.Repository.Scan(uow, conceptModules,
	// 	repository.Table("modules_programming_concepts"),
	// 	repository.Select("modules_programming_concepts.*, AVG(talent_concept_ratings.`score`) AS `average_score`"),
	// 	repository.Join("LEFT JOIN talent_concept_ratings on talent_concept_ratings.`module_programming_concept_id` = modules_programming_concepts.`id`"),
	// 	service.addSearchQueries(parser.Form),
	// 	repository.Filter("modules_programming_concepts.`tenant_id` = ? AND modules_programming_concepts.`deleted_at` IS NULL", tenantID),
	// 	repository.Filter("talent_concept_ratings.`deleted_at` IS NULL"),
	// 	repository.Filter("talent_concept_ratings.`talent_id` = ? OR talent_concept_ratings.`talent_id` IS NULL", talentID),
	// 	repository.GroupBy("modules_programming_concepts.`id`"))
	// 	// repository.PreloadAssociations([]string{"ParentModuleProgrammingConcepts"}))
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	// fmt.Println((*conceptModules)[0].AverageScore)
	// fmt.Println((*conceptModules)[1].AverageScore)
	// fmt.Println((*conceptModules)[2].AverageScore)
	// fmt.Println((*conceptModules)[3].AverageScore)
	// fmt.Println((*conceptModules)[4].AverageScore)

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// addSearchQueries adds search criteria.
func (service *ModuleProgrammingConceptService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if moduleID, ok := requestForm["moduleID"]; ok {
		util.AddToSlice("modules_programming_concepts.`module_id`", "=?", "AND", moduleID,
			&columnNames, &conditions, &operators, &values)
	}

	if conceptID, ok := requestForm["conceptID"]; ok {
		util.AddToSlice("modules_programming_concepts.`concept_id`", "=?", "AND", conceptID,
			&columnNames, &conditions, &operators, &values)
	}

	if level, ok := requestForm["level"]; ok {
		util.AddToSlice("`level`", "=?", "AND", level, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doForeignKeysExist checks if all foregin keys exist and if not returns error.
func (service *ModuleProgrammingConceptService) doForeignKeysExist(tenantID, credentialID uuid.UUID,
	conceptModules *[]programming.ModuleProgrammingConcepts) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	for i := range *conceptModules {

		// Check if programming concept exists.
		if err := service.doesProgrammingConceptExist(tenantID, (*conceptModules)[i].ProgrammingConceptID); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}

		// Check if module exists.
		if err := service.doesModuleExist(tenantID, (*conceptModules)[i].ModuleID); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}

		// // Check if concept module exists.
		// if err := service.doesConceptModuleExist(tenantID, &(*conceptModules)[i]); err != nil {
		// 	log.NewLogger().Error(err.Error())
		// 	return err
		// }
	}
	return nil
}

// doesConceptModuleExist returns error if there is module id for the concept id.
func (service *ModuleProgrammingConceptService) doesConceptModuleExist(tenantID uuid.UUID, conceptModule *programming.ModuleProgrammingConcepts) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ModuleProgrammingConcepts{},
		repository.Filter("`programming_concept_id`=? AND module_id=?", conceptModule.ProgrammingConceptID, conceptModule.ModuleID))
	if err := util.HandleIfExistsError("Concept already exists for the module", exists, err); err != nil {
		return err
	}
	return nil
}

// doesProgrammingConceptExist returns error if there is no programming concept record in table for the given tenant.
func (service *ModuleProgrammingConceptService) doesProgrammingConceptExist(tenantID, conceptID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingConcept{},
		repository.Filter("`id`=?", conceptID))
	if err := util.HandleError("Invalid programming concept ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesModuleExist returns error if there is no module record in table for the given tenant.
func (service *ModuleProgrammingConceptService) doesModuleExist(tenantID, moduleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Module{},
		repository.Filter("`id`=?", moduleID))
	if err := util.HandleError("Invalid module ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *ModuleProgrammingConceptService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *ModuleProgrammingConceptService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// getParentConceptModuleID gets the parent concept module id for given concept id and module id.
func (service *ModuleProgrammingConceptService) getParentConceptModuleID(uow *repository.UnitOfWork, tenantID uuid.UUID,
	conceptModule *programming.ModuleProgrammingConcepts) error {

	for j := range conceptModule.ParentModuleProgrammingConcepts {

		// Create bucket for concept module.
		tempConceptModule := programming.ModuleProgrammingConcepts{}

		// Get concept module.
		err := service.Repository.GetRecordForTenant(uow, tenantID, &tempConceptModule,
			repository.Filter("programming_concept_id=?", conceptModule.ParentModuleProgrammingConcepts[j].ProgrammingConceptID),
			repository.Filter("module_id=?", conceptModule.ParentModuleProgrammingConcepts[j].ModuleID))
		if err != nil {
			uow.RollBack()
			return err
		}

		// Give concept module id to parent concept module id.
		conceptModule.ParentModuleProgrammingConcepts[j].ID = tempConceptModule.ID
	}

	return nil
}
