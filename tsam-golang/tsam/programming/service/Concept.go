package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingConceptService provides methods to update, delete, add, get for programming concept.
type ProgrammingConceptService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewProgrammingConceptService returns new instance of ProgrammingConceptService.
func NewProgrammingConceptService(db *gorm.DB, repository repository.Repository) *ProgrammingConceptService {
	return &ProgrammingConceptService{
		DB:         db,
		Repository: repository,
		association: []string{
			"ProgrammingQuestions",
		},
	}
}

// AddProgrammingConcept will add new programming concept to the table.
func (service *ProgrammingConceptService) AddProgrammingConcept(concept *programming.ProgrammingConcept) error {

	// Checks if all foreign keys exist.
	err := service.doForeignKeysExist(concept, concept.CreatedBy)
	if err != nil {
		return err
	}

	// // Set createdBy field and tenant id to all modules.
	// if concept.ConceptModules != nil && len(concept.ConceptModules) != 0 {
	// 	for i := 0; i < len(concept.ConceptModules); i++ {
	// 		concept.ConceptModules[i].CreatedBy = concept.CreatedBy
	// 		concept.ConceptModules[i].TenantID = concept.TenantID
	// 		err := service.doesConceptModuleExist(concept.ConceptModules[i])
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add programming concept.
	err = service.Repository.Add(uow, concept)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateProgrammingConcept will update the specified programming concept record in the table.
func (service *ProgrammingConceptService) UpdateProgrammingConcept(concept *programming.ProgrammingConcept) error {

	// Checks if all foreign key exist.
	err := service.doForeignKeysExist(concept, concept.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if programming concept record exist.
	err = service.doesProgrammingConceptExist(concept.TenantID, concept.ID)
	if err != nil {
		return err
	}

	// // Valudate concept module.
	// if concept.ConceptModules != nil && len(concept.ConceptModules) != 0 {
	// 	for i := 0; i < len(concept.ConceptModules); i++ {
	// 		err := service.doesConceptModuleExist(concept.ConceptModules[i])
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// // Update modules.
	// err = service.updateModules(uow, concept.ConceptModules, concept.TenantID, concept.UpdatedBy, concept.ID)
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }

	// // Make modules nil to avoid any inserts or updates in concepts modules table.
	// concept.ConceptModules = nil

	// Create bucket for getting programming concept already present in database.
	tempConcept := programming.ProgrammingConcept{}

	// Get programming concept round for getting created_by field of programming concept from database.
	if err := service.Repository.GetForTenant(uow, concept.TenantID, concept.ID, &tempConcept); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp programming concept to programming concept to be updated.
	concept.CreatedBy = tempConcept.CreatedBy

	// Update programming concept.
	err = service.Repository.Save(uow, concept)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteProgrammingConcept will delete the specified programming concept record from the table.
func (service *ProgrammingConceptService) DeleteProgrammingConcept(concept *programming.ProgrammingConcept) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(concept.TenantID); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(concept.TenantID, concept.DeletedBy); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if programming concept record exist.
	err := service.doesProgrammingConceptExist(concept.TenantID, concept.ID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// // Delete concept modules.
	// err = service.Repository.UpdateWithMap(uow, &programming.ModuleProgrammingConcept{}, map[interface{}]interface{}{
	// 	"DeletedBy": concept.DeletedBy,
	// 	"DeletedAt": time.Now(),
	// }, repository.Filter("`programming_concept_id`=?", concept.ID))
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }

	// Delete programming concept.
	err = service.Repository.UpdateWithMap(uow, &programming.ProgrammingConcept{}, map[interface{}]interface{}{
		"DeletedBy": concept.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id`=?", concept.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllProgrammingConcepts will return all the records from programming concept table.
func (service *ProgrammingConceptService) GetAllProgrammingConcepts(concepts *[]programming.ProgrammingConceptDTO, tenantID uuid.UUID,
	form url.Values, pagination repository.Paginator) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	queryProcessors := service.addSearchQueries(form)
	queryProcessors = append(queryProcessors,
		repository.Filter("`programming_concepts`.`tenant_id` = ?", tenantID),
		pagination.Paginate(),
		repository.PreloadAssociations([]string{"ConceptModules"}))

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)
	
	// Get all programming concepts.
	err = service.Repository.GetAllInOrder(uow, concepts, "`name`", queryProcessors...)
	// err = service.Repository.GetAllInOrderForTenant(uow, tenantID, concepts, "`name`",
	// 	service.addQueryProcessors(form,pagination)...)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetProgrammingConcept will return specified record from programming concept table.
func (service *ProgrammingConceptService) GetProgrammingConcept(concept *programming.ProgrammingConceptDTO,
	tenantID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if programming concept record exist.
	err = service.doesProgrammingConceptExist(tenantID, concept.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetForTenant(uow, tenantID, concept.ID, concept)
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

// doForeignKeysExist checks if all foregin keys exist and if not returns error.
func (service *ProgrammingConceptService) doForeignKeysExist(conecpt *programming.ProgrammingConcept, credentialID uuid.UUID) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(conecpt.TenantID); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(conecpt.TenantID, credentialID); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if programming concept with same name exist.
	if err := service.doesProgrammingConceptNameExist(conecpt); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// #Niranjan
// func (service *ProgrammingConceptService) addQueryProcessors(requestForm url.Values,
// 	pagination repository.Paginator) []repository.QueryProcessor {
// 	temp := []repository.QueryProcessor{service.addSearchQueries(requestForm), pagination.Paginate()}
// 		temp = append(temp, service.addOptionalData(requestForm)...)
// 	return temp
// }

// // addSearchQueries adds search criteria.
// func (service *ProgrammingConceptService) addOptionalData(requestForm url.Values) []repository.QueryProcessor {
// 	var qps []repository.QueryProcessor
// 	if batchTopicAssignmentID := requestForm.Get("batchTopicAssignmentID"); !util.IsEmpty(batchTopicAssignmentID) {
// 		qps = append(qps, repository.Join("INNER JOIN"))
// 	}
// 	if programmingQuestionID := requestForm.Get("programmingQuestionID"); !util.IsEmpty(programmingQuestionID) {
// 		qps = append(qps, repository.Join("INNER JOIN `programming_questions_programming_concepts` ON "+
// 		"`programmming_concept`.`id` = `programming_questions_programming_concepts`.`programmming_concept_id`",
// 		))
// 	}

// 	return nil
// }

// addSearchQueries adds search criteria.
func (service *ProgrammingConceptService) addSearchQueries(requestForm url.Values) []repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcesors []repository.QueryProcessor

	if _, ok := requestForm["name"]; ok {
		util.AddToSlice("`name`", "LIKE ?", "AND", "%"+requestForm.Get("name")+"%",
			&columnNames, &conditions, &operators, &values)
	}

	if topicID, ok := requestForm["topicID"]; ok {
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN `topic_programming_concepts` ON `topic_programming_concepts`.`programming_concept_id` = "+
				"`programming_concepts`.`id` AND `topic_programming_concepts`.`tenant_id` = `programming_concepts`.`tenant_id`"),
			repository.Filter("topic_programming_concepts.`deleted_at` IS NULL"))

		util.AddToSlice("`topic_programming_concepts`.`topic_id`", "= ?", "AND", topicID,
		&columnNames, &conditions, &operators, &values)
	}

	// if topicID := requestForm.Get("topicID"); !util.IsEmpty(topicID) {
	// 	util.AddToSlice("`topic_programming_concepts`.`topic_id`", "= ?", "AND", topicID,
	// 		&columnNames, &conditions, &operators, &values)
	// }

	queryProcesors = append(queryProcesors,
		repository.FilterWithOperator(columnNames, conditions, operators, values))
	return queryProcesors
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *ProgrammingConceptService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *ProgrammingConceptService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesProgrammingConceptNameExist returns error if there is no programming concept record with same name exist in table for the given tenant.
func (service *ProgrammingConceptService) doesProgrammingConceptNameExist(concept *programming.ProgrammingConcept) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, concept.TenantID, programming.ProgrammingConcept{},
		repository.Filter("`name`=? AND `id`!=?", concept.Name, concept.ID))
	if err := util.HandleIfExistsError("Programming concept name already exists", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesProgrammingConceptExist returns error if there is no programming concept record in table for the given tenant.
func (service *ProgrammingConceptService) doesProgrammingConceptExist(tenantID, conceptID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingConcept{},
		repository.Filter("`id`=?", conceptID))
	if err := util.HandleError("Invalid programming concept ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesConceptModuleExist returns error if there is module id for the concept id.
func (service *ProgrammingConceptService) doesConceptModuleExist(conceptModule *programming.ModuleProgrammingConcepts) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, conceptModule.TenantID, programming.ModuleProgrammingConcepts{},
		repository.Filter("`programming_concept_id`=? AND module_id=?", conceptModule.ProgrammingConceptID, conceptModule.ModuleID))
	if err := util.HandleIfExistsError("Module already exists for the concept", exists, err); err != nil {
		return err
	}
	return nil
}

// updateModules will update the concept modules for specified concept.
func (service *ProgrammingConceptService) updateModules(uow *repository.UnitOfWork, conceptModules []*programming.ModuleProgrammingConcepts,
	tenantID, credentialID, conceptID uuid.UUID) error {

	// If previous concept modules is present and current concept modules is not present.
	if conceptModules == nil {
		err := service.Repository.UpdateWithMap(uow, programming.ModuleProgrammingConcepts{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		}, repository.Filter("`programming_concept_id`=?", conceptID))
		if err != nil {
			return err
		}
	}

	// Create temp concept modules to get presvious concept modules of concept.
	tempConceptModules := []programming.ModuleProgrammingConcepts{}

	err := service.Repository.GetAllForTenant(uow, tenantID, &tempConceptModules,
		repository.Filter("`programming_concept_id`=?", conceptID))
	if err != nil {
		return err
	}

	// Make map to count occurences of concept module id in previous and current concept modules.
	conceptModulesIDMap := make(map[uuid.UUID]uint)

	// Count the number of occurence of previous concept module ID.
	for _, tempConceptModule := range tempConceptModules {
		conceptModulesIDMap[tempConceptModule.ID]++
	}

	// Count the number of occurrence of current concept module ID to know total count of occurrenec of each ID.
	for _, conceptModule := range conceptModules {

		// If ID is valid then push its occurence in ID map.
		if util.IsUUIDValid(conceptModule.ID) {
			conceptModulesIDMap[conceptModule.ID]++
		} else {
			// If ID is nil create new concept module entry in table.
			conceptModule.CreatedBy = credentialID
			conceptModule.TenantID = tenantID
			conceptModule.ProgrammingConceptID = conceptID
			err = service.Repository.Add(uow, &conceptModule)
			if err != nil {
				return err
			}
		}

		// If number of occurrence is more than one (present in previous and current concept modules)
		// then update concept module.
		if conceptModulesIDMap[conceptModule.ID] > 1 {
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
	return nil
}
