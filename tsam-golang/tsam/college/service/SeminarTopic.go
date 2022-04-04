package service

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/college"
	general "github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// SeminarTopicService provides method to update, delete, add, get all, get one for semiar topics.
type SeminarTopicService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// seminarTopicAssociations provides preload associations array for seminar topic.
var seminarTopicAssociations []string = []string{"Speaker"}

// NewSeminarTopicService returns new instance of SeminarTopicService.
func NewSeminarTopicService(db *gorm.DB, repository repository.Repository) *SeminarTopicService {
	return &SeminarTopicService{
		DB:         db,
		Repository: repository,
	}
}

// AddSeminarTopic adds one semiar topic to database.
func (service *SeminarTopicService) AddSeminarTopic(topic *college.Topic) error {
	// Get credential id from CreatedBy field of topic(set in controller).
	credentialID := topic.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(topic.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, topic.TenantID); err != nil {
		return err
	}

	// Validate seminar id.
	if err := service.doesSeminarExist(topic.TenantID, topic.SeminarID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(topic); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add semiar topic to database.
	if err := service.Repository.Add(uow, topic); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("semiar topic could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetSeminarTopics gets all semiar topics from database.
func (service *SeminarTopicService) GetSeminarTopics(topics *[]college.TopicDTO,
	tenantID uuid.UUID, seminarID uuid.UUID, uows ...*repository.UnitOfWork) error {
	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate seminar id.
	if err := service.doesSeminarExist(tenantID, seminarID); err != nil {
		return err
	}

	// Get semiar topics from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, topics, "`from_time`",
		repository.Filter("`seminar_id`=?", seminarID),
		repository.PreloadAssociations(seminarTopicAssociations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// GetSeminarTopic gets one semiar topic form database.
func (service *SeminarTopicService) GetSeminarTopic(topic *college.Topic) error {
	// Validate tenant id.
	if err := service.doesTenantExist(topic.TenantID); err != nil {
		return err
	}

	// Validate seminar id.
	if err := service.doesSeminarExist(topic.TenantID, topic.SeminarID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get semiar topic.
	if err := service.Repository.GetForTenant(uow, topic.TenantID, topic.ID, topic,
		repository.PreloadAssociations(seminarTopicAssociations)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// UpdateSeminarTopic updates semiar topic in Database.
func (service *SeminarTopicService) UpdateSeminarTopic(topic *college.Topic) error {
	// Get credential id from UpdatedBy field of topic(set in controller).
	credentialID := topic.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(topic.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, topic.TenantID); err != nil {
		return err
	}

	// Validate seminar id.
	if err := service.doesSeminarExist(topic.TenantID, topic.SeminarID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(topic); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting semiar topic already present in database.
	tempTopic := college.Topic{}

	// Get semiar topic for getting created_by field of semiar topic from database.
	if err := service.Repository.GetForTenant(uow, topic.TenantID, topic.ID, &tempTopic); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp semiar topic to semiar topic to be updated.
	topic.CreatedBy = tempTopic.CreatedBy

	// Update semiar topic.
	if err := service.Repository.Save(uow, topic); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Seminar topic could not be updated", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// DeleteSeminarTopic deletes one semiar topic form database.
func (service *SeminarTopicService) DeleteSeminarTopic(topic *college.Topic) error {
	// Starting new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get credential id from DeletedBy field of topic(set in controller).
	credentialID := topic.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(topic.TenantID); err != nil {
		return err
	}

	// Validate semiar topic id.
	if err := service.doesSeminarTopicExist(topic.ID, topic.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, topic.TenantID); err != nil {
		return err
	}

	// Validate seminar id.
	if err := service.doesSeminarExist(topic.TenantID, topic.SeminarID); err != nil {
		return err
	}

	// Update semiar topic for updating deleted_by and deleted_at field of semiar topic.
	if err := service.Repository.UpdateWithMap(uow, &college.Topic{}, map[string]interface{}{
		"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`tenant_id`=? AND `id`=?", topic.TenantID, topic.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("semiar topic could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *SeminarTopicService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *SeminarTopicService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesSeminarExist validates if seminar exists or not in database.
func (service *SeminarTopicService) doesSeminarExist(tenantID uuid.UUID, seminarID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.Seminar{},
		repository.Filter("`id` = ?", seminarID))
	if err := util.HandleError("Invalid seminar ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesSeminarTopicExist validates if semiar topic exists or not in database.
func (service *SeminarTopicService) doesSeminarTopicExist(topicID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.Topic{},
		repository.Filter("`id` = ?", topicID))
	if err := util.HandleError("Invalid semiar topic ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExist validates if speaker exists or not in database.
func (service *SeminarTopicService) doForeignKeysExist(topic *college.Topic) error {

	// Check if speaker exists or not.
	if topic.SpeakerID != nil {
		exists, err := repository.DoesRecordExist(service.DB, college.Speaker{},
			repository.Filter("`id` = ?", topic.SpeakerID))
		if err := util.HandleError("Invalid speaker ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}
	return nil
}
