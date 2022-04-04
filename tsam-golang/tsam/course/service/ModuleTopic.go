package service

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TopicService Provide method to Update, Delete, Add, Get Method For topics.
type TopicService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewTopicService returns new instance of TopicService.
func NewTopicService(db *gorm.DB, repository repository.Repository) *TopicService {
	return &TopicService{
		DB:          db,
		Repository:  repository,
		association: []string{
			// "Resources",
			// "SubTopics",
		},
	}
}

// AddTopic will add topic for specified module.
func (service *TopicService) AddTopic(topic *course.ModuleTopic, uows ...*repository.UnitOfWork) error {

	// check all foreign key exists
	err := service.doesForeignKeyExist(topic, topic.CreatedBy)
	if err != nil {
		return err
	}

	// order validation for topic
	if topic.TopicID == nil {
		err = service.doesTopicOrderExist(topic)
		if err != nil {
			return err
		}

		orderMap := map[uint]bool{}
		for index := range topic.SubTopics {

			if orderMap[topic.SubTopics[index].Order] {
				return errors.NewValidationError("duplicate order given for sub-topic")
			}
			orderMap[topic.SubTopics[index].Order] = true

			topic.SubTopics[index].TenantID = topic.TenantID
			topic.SubTopics[index].CreatedBy = topic.CreatedBy
			topic.SubTopics[index].ModuleID = topic.ModuleID
		}
	}

	if topic.TopicID != nil {
		err = service.doesSubTopicOrderExist(topic)
		if err != nil {
			return err
		}
	}

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)

	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	err = service.Repository.Add(uow, topic)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return err
	}
	fmt.Println("**************************", topic.TopicProgrammingConcept)

	fmt.Println("**************************", topic.SubTopics[0].TopicProgrammingConcept)

	// Add programming concepts of all sub topics.
	// for _, subTopic := range topic.SubTopics {
	for _, concept := range topic.TopicProgrammingConcept {

		concept.CreatedBy = topic.CreatedBy
		concept.TopicID = topic.ID
		concept.TenantID = topic.TenantID

		err = service.Repository.Add(uow, concept)
		if err != nil {
			if length == 0 {
				uow.RollBack()
			}
			return err
		}
	}
	// }

	// for _, subTopic := range topic.SubTopics {
	// 	err = service.Repository.Add(uow, &subTopic.TopicProgrammingConcept)
	// 	if err != nil {
	// 		if length == 0 {
	// 			uow.RollBack()
	// 		}
	// 		return err
	// 	}
	// }

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddTopics will add multiple topics for specified module.
func (service *TopicService) AddTopics(topics *[]course.ModuleTopic, tenantID,
	moduleID, credentialID uuid.UUID) error {

	err := service.checkTopicInJSON(topics)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	for _, topic := range *topics {
		// assign foreign keys
		topic.TenantID = tenantID
		topic.ModuleID = moduleID
		topic.CreatedBy = credentialID

		// assign foreign keys to sub-topics
		for index := range topic.SubTopics {
			topic.SubTopics[index].TenantID = tenantID
			topic.SubTopics[index].ModuleID = moduleID
			topic.SubTopics[index].CreatedBy = credentialID
		}

		err := service.AddTopic(&topic, uow)
		if err != nil {
			uow.RollBack()
			return err
		}

	}
	uow.Commit()
	return nil
}

// UpdateTopic will update the topics for specified module
func (service *TopicService) UpdateTopic(topic *course.ModuleTopic) error {

	// check all foreign key exists
	err := service.doesForeignKeyExist(topic, topic.UpdatedBy)
	if err != nil {
		return err
	}

	// order validation for topic
	if topic.TopicID == nil {
		err = service.doesTopicOrderExist(topic)
		if err != nil {
			return err
		}
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempTopic := course.ModuleTopic{}
	err = service.Repository.GetRecordForTenant(uow, topic.TenantID, &tempTopic,
		repository.Select("created_by"), repository.Filter("`id` = ?", topic.ID))
	if err != nil {
		uow.RollBack()
		return err
	}
	topic.CreatedBy = tempTopic.CreatedBy

	// update associations
	err = service.updateSubTopics(uow, topic, topic.ID, topic.ModuleID,
		topic.TenantID, topic.UpdatedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.updateTopicProgrammingConcept(uow, topic, topic.TenantID, topic.UpdatedBy)
	if err != nil {
		return err
	}

	err = service.Repository.Save(uow, topic)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteTopic will delete specified topic from the DB.
func (service *TopicService) DeleteTopic(topic *course.ModuleTopic) error {

	// check if tenant exist
	err := service.doesTenantExists(topic.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExists(topic.TenantID, topic.DeletedBy)
	if err != nil {
		return err
	}

	// check if course exist
	err = service.doesModuleExist(topic.TenantID, topic.ModuleID)
	if err != nil {
		return err
	}

	// check if topic exist for tenant
	err = service.doesTopicExist(topic.TenantID, topic.ID)
	if err != nil {
		return err
	}

	// check if topic is assigned to batch
	exist, err := repository.DoesRecordExistForTenant(service.DB, topic.TenantID, batch.SessionTopic{},
		repository.Filter("`topic_id` = ? OR `sub_topic_id` = ?", topic.ID, topic.ID))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("topic cannot be deleted as it has been assigned in session plan")
	}

	// Start transaction
	uow := repository.NewUnitOfWork(service.DB, false)

	// delete sub-topics for the current topic
	err = service.deleteSubTopic(uow, topic)
	if err != nil {
		uow.RollBack()
		return err
	}

	// delete the specified topic
	err = service.Repository.UpdateWithMap(uow, &course.ModuleTopic{}, map[string]interface{}{
		"DeletedBy": topic.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", topic.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllTopics will return all the topics.
func (service *TopicService) GetAllTopics(tenantID, moduleID uuid.UUID, topics *[]course.ModuleTopicDTO,
	parser *web.Parser, totalCount *int) error {

	// check if tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()

	uow := repository.NewUnitOfWork(service.DB, true)

	var queryProcessors []repository.QueryProcessor
	queryProcessors = append(queryProcessors, repository.Filter("module_topics.`topic_id` IS NULL AND"+
		" module_topics.`module_id` = ? AND module_topics.`tenant_id` = ?", moduleID, tenantID),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "SubTopics",
			Queryprocessors: []repository.QueryProcessor{
				repository.OrderBy("`order`"), repository.PreloadAssociations([]string{
					"Module", "TopicProgrammingConcept", "TopicProgrammingConcept.ProgrammingConcept",
					"TopicProgrammingQuestions", "TopicProgrammingQuestions.ProgrammingQuestion",
				}),
			},
		}), repository.PreloadAssociations([]string{
			"Module", "TopicProgrammingConcept", "TopicProgrammingConcept.ProgrammingConcept",
			"TopicProgrammingQuestions", "TopicProgrammingQuestions.ProgrammingQuestion",
		}), repository.Paginate(limit, offset, totalCount))
	queryProcessors = append(queryProcessors, service.addSearchQueries(parser.Form)...)

	err = service.Repository.GetAllInOrder(uow, topics, "module_topics.`order`", queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTopicList will return list of all the topics.
func (service *TopicService) GetTopicList(tenantID uuid.UUID, topics *[]list.ModuleTopic,
	parser *web.Parser) error {

	// check if tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	var queryProcessors []repository.QueryProcessor

	queryProcessors = append(queryProcessors, repository.Filter("module_topics.`topic_id` IS NULL AND module_topics.`tenant_id` = ?", tenantID))
	queryProcessors = append(queryProcessors, service.addSearchQueries(parser.Form)...)

	err = service.Repository.GetAllInOrder(uow, topics, "module_topics.`order`", queryProcessors...)
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

// doesForeignKeyExist checks if all the foreign key are valid
func (service *TopicService) doesForeignKeyExist(topic *course.ModuleTopic, credentialID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExists(topic.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExists(topic.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if course exist
	err = service.doesModuleExist(topic.TenantID, topic.ModuleID)
	if err != nil {
		return err
	}

	if topic.TopicID != nil {
		err = service.doesTopicExist(topic.TenantID, *topic.TopicID)
		if err != nil {
			return err
		}
	}

	return nil
}

// addSearchQueries adds search criteria.
func (service *TopicService) addSearchQueries(requestForm url.Values) []repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcessors []repository.QueryProcessor

	if _, ok := requestForm["topicName"]; ok {
		util.AddToSlice("`topic_name`", "LIKE ?", "AND", "%"+requestForm.Get("topicName")+"%", &columnNames, &conditions, &operators, &values)
	}

	if moduleID, ok := requestForm["moduleID"]; ok {
		util.AddToSlice("`module_id`", "= ?", "AND", moduleID, &columnNames, &conditions, &operators, &values)
	}

	if batchID, ok := requestForm["batchID"]; ok {
		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN batch_session_topics ON"+
			" batch_session_topics.`topic_id` = module_topics.`id` AND module_topics.`tenant_id` = batch_session_topics.`tenant_id`"),
			repository.Filter("batch_session_topics.`deleted_at` IS NULL AND batch_session_topics.`batch_id` = ?", batchID))
	}

	queryProcessors = append(queryProcessors, repository.FilterWithOperator(columnNames, conditions, operators, values))
	queryProcessors = append(queryProcessors, repository.GroupBy("module_topics.`id`"))
	return queryProcessors
}

// deleteSubTopic will delete all the sub-topics for the current topic.
func (service *TopicService) deleteSubTopic(uow *repository.UnitOfWork, topic *course.ModuleTopic) error {

	err := service.Repository.UpdateWithMap(uow, &course.ModuleTopic{}, map[string]interface{}{
		"DeletedBy": topic.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`topic_id`=?", topic.ID))
	if err != nil {
		return err
	}

	return nil
}

// updateSubTopics will update sub-topics
func (service *TopicService) updateSubTopics(uow *repository.UnitOfWork, topic *course.ModuleTopic,
	topicID, moduleID, tenantID, credentialID uuid.UUID) error {

	subTopics := &topic.SubTopics
	tempSubTopics := &[]course.ModuleTopic{}
	topicMap := make(map[uuid.UUID]uint)

	err := service.checkTopicInJSON(subTopics)
	if err != nil {
		return err
	}

	err = service.Repository.GetAllForTenant(uow, tenantID, tempSubTopics,
		repository.Filter("`module_id`=? AND `topic_id`=?", moduleID, topicID))
	if err != nil {
		return err
	}

	// populating topicMap
	for _, tempSubTopic := range *tempSubTopics {
		topicMap[tempSubTopic.ID] = 1
	}

	// checking with new topics on new, existing and deleted sub-topics
	for index, subTopic := range *subTopics {

		// if it is not an existing sub-topic
		if util.IsUUIDValid(subTopic.ID) {
			topicMap[subTopic.ID]++
		}

		if topicMap[subTopic.ID] > 1 {

			// check foreign keys of sub-topic
			assignForeignKeys(&subTopic, moduleID, tenantID, topicID, index+1)
			subTopic.UpdatedBy = credentialID

			err = service.Repository.Update(uow, &subTopic)
			if err != nil {
				return nil
			}

			topicMap[subTopic.ID] = 0
		}

		// new entry of sub-topic
		if !util.IsUUIDValid(subTopic.ID) {

			// check if order exists
			exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.ModuleTopic{},
				repository.Filter("module_id=? AND topic_id=? AND `order`=?", moduleID, topicID, subTopic.Order))
			if err != nil {
				return err
			}
			if exists {
				return errors.NewValidationError("Similar order exist")
			}

			assignForeignKeys(&subTopic, moduleID, tenantID, topicID, index+1)
			subTopic.CreatedBy = credentialID

			// adding new record to the DB
			err = service.Repository.Add(uow, &subTopic)
			if err != nil {
				return err
			}
		}

		// err = service.updateTopicProgrammingConcept(uow, &subTopic, tenantID, credentialID)
		// if err != nil {
		// 	return err
		// }

	}

	// deleting records where value is 1 as they are the deleted entries
	for _, tempSubTopic := range *tempSubTopics {
		if topicMap[tempSubTopic.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, &tempSubTopic, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			})
			if err != nil {
				return err
			}

			err = service.Repository.UpdateWithMap(uow, course.TopicProgrammingConcept{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`topic_id` = ?", tempSubTopic.ID))
			if err != nil {
				return err
			}

			topicMap[tempSubTopic.ID] = 0
		}

	}

	topic.SubTopics = nil
	return nil
}

func (service *TopicService) updateTopicProgrammingConcept(uow *repository.UnitOfWork, topic *course.ModuleTopic,
	tenantID, credentialID uuid.UUID) error {

	tempTopicProgrammingAssignments := []course.TopicProgrammingConcept{}
	topicProgrammingConceptMap := make(map[uuid.UUID]uint)

	err := service.Repository.GetAllForTenant(uow, tenantID, &tempTopicProgrammingAssignments,
		repository.Filter("`topic_id` = ?", topic.ID))
	if err != nil {
		return err
	}

	// populating topicProgrammingConceptMap
	for _, tempTopicProgrammingAssignment := range tempTopicProgrammingAssignments {
		topicProgrammingConceptMap[tempTopicProgrammingAssignment.ID] = 1
	}

	for _, topicProgrammingConcept := range topic.TopicProgrammingConcept {

		// if it is not an existing programmingConcept
		if util.IsUUIDValid(topicProgrammingConcept.ID) {
			topicProgrammingConceptMap[topicProgrammingConcept.ID]++
		}

		// new entry of programmingConcept
		if !util.IsUUIDValid(topicProgrammingConcept.ID) {

			topicProgrammingConcept.CreatedBy = credentialID
			topicProgrammingConcept.TopicID = topic.ID
			topicProgrammingConcept.TenantID = tenantID

			// adding new record to the DB
			err = service.Repository.Add(uow, &topicProgrammingConcept)
			if err != nil {
				return err
			}
		}

		if topicProgrammingConceptMap[topicProgrammingConcept.ID] > 1 {

			// check foreign keys of programmingConcept
			topicProgrammingConcept.TenantID = tenantID
			topicProgrammingConcept.TopicID = topic.ID
			topicProgrammingConcept.UpdatedBy = credentialID

			err = service.Repository.Update(uow, &topicProgrammingConcept)
			if err != nil {
				return nil
			}

			topicProgrammingConceptMap[topicProgrammingConcept.ID] = 0
			continue
		}

	}

	// deleting records where value is 1 as they are the deleted entries
	for _, tempTopicProgrammingAssignment := range tempTopicProgrammingAssignments {

		if topicProgrammingConceptMap[tempTopicProgrammingAssignment.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, course.TopicProgrammingConcept{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempTopicProgrammingAssignment.ID))
			if err != nil {
				return err
			}

			topicProgrammingConceptMap[tempTopicProgrammingAssignment.ID] = 0
		}
	}
	topic.TopicProgrammingConcept = nil
	return nil
}

func assignForeignKeys(topic *course.ModuleTopic, moduleID, tenantID,
	topicID uuid.UUID, index int) {
	topic.ModuleID = moduleID
	topic.TenantID = tenantID
	topic.TopicID = &topicID
	topic.Order = uint(index)
}

// checkTopicInJSON will check order of topics in JSON.
func (service *TopicService) checkTopicInJSON(topics *[]course.ModuleTopic) error {

	topicMap := make(map[uint]uint)

	for _, topic := range *topics {
		topicMap[topic.Order]++
		if topicMap[topic.Order] > 1 {
			return errors.NewValidationError("Same order given for two topics")
		}
	}

	return nil
}

func (service *TopicService) doesSubTopicOrderExist(topic *course.ModuleTopic) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, topic.TenantID, &course.ModuleTopic{},
		repository.Filter("`module_id`=? AND `order`=? AND `id`!=? AND `topic_id` =?",
			topic.ModuleID, topic.Order, topic.ID, topic.TopicID))
	if err != nil {
		return err
	}
	if exists {
		return errors.NewValidationError("Sub Topic with similar order already present")
	}
	return nil
}

func (service *TopicService) doesTopicOrderExist(topic *course.ModuleTopic) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, topic.TenantID, &course.ModuleTopic{},
		repository.Filter("`module_id`=? AND `order`=? AND `id`!=? AND `topic_id` IS NULL",
			topic.ModuleID, topic.Order, topic.ID))
	if err := util.HandleIfExistsError("Topic with similar order already present", exists, err); err != nil {
		return err
	}
	// if err != nil {
	// 	return err
	// }
	// if exists {
	// 	return errors.NewValidationError("Topic with similar order already present")
	// }
	return nil
}

// doesTenantExists validates tenantID
func (service *TopicService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExists validates courseID
func (service *TopicService) doesCredentialExists(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Credential not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesModuleExist validates courseID
func (service *TopicService) doesModuleExist(tenantID, moduleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Module{},
		repository.Filter("`id` = ?", moduleID))
	if err := util.HandleError("Module not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTopicExist validates courseID
func (service *TopicService) doesTopicExist(tenantID, topicID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.ModuleTopic{},
		repository.Filter("`id` = ?", topicID))
	if err := util.HandleError("Topic not found", exists, err); err != nil {
		return err
	}
	return nil
}
