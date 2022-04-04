package service

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// FeedbackOptionService provides methods to do different CRUD operations on feedback_options table.
type FeedbackOptionService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewFeedbackOptionService returns a new instance Of FeedbackOptionService.
func NewFeedbackOptionService(db *gorm.DB, repository repository.Repository) *FeedbackOptionService {
	return &FeedbackOptionService{
		DB:         db,
		Repository: repository,
	}
}

// AddFeedbackOption will add feedback options to the table
func (service *FeedbackOptionService) AddFeedbackOption(feedbackOption *general.FeedbackOption, uows ...*repository.UnitOfWork) error {

	// check if tenant exist
	err := service.doesTenantExist(feedbackOption.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExist(feedbackOption.TenantID, feedbackOption.CreatedBy)
	if err != nil {
		return err
	}

	// check if feedback question exist
	err = service.doesFeedbackQuestionExist(feedbackOption.TenantID, feedbackOption.QuestionID)
	if err != nil {
		return err
	}

	// check if fields are unique
	err = service.validateFieldUniqueness(feedbackOption)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	tempQuestion := general.FeedbackQuestion{}
	// check if question can have options
	err = service.Repository.GetRecordForTenant(uow, feedbackOption.TenantID, &tempQuestion,
		repository.Filter("`id` = ?", feedbackOption.QuestionID), repository.Select("has_options"))
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return err
	}
	if !(*tempQuestion.HasOptions) {
		return errors.NewValidationError("This question cannot have options.")
	}

	err = service.Repository.Add(uow, feedbackOption)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return err
	}
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddFeedbackOptions will add multiple feedback options to the table
func (service *FeedbackOptionService) AddFeedbackOptions(feedbackOptions *[]general.FeedbackOption, tenantID, credentialID,
	questionID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, false)
	for _, feedbackOption := range *feedbackOptions {
		feedbackOption.TenantID = tenantID
		feedbackOption.CreatedBy = credentialID
		feedbackOption.QuestionID = questionID
		err := service.AddFeedbackOption(&feedbackOption)
		if err != nil {
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// UpdateFeedbackOption will update the specified feedback option in the table
func (service *FeedbackOptionService) UpdateFeedbackOption(feedbackOption *general.FeedbackOption) error {

	// check if tenant exist
	err := service.doesTenantExist(feedbackOption.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExist(feedbackOption.TenantID, feedbackOption.UpdatedBy)
	if err != nil {
		return err
	}

	// check if feedback question exist
	err = service.doesFeedbackQuestionExist(feedbackOption.TenantID, feedbackOption.QuestionID)
	if err != nil {
		return err
	}

	// check if feedback option exist
	err = service.doesFeedbackOptionExist(feedbackOption.TenantID, feedbackOption.ID)
	if err != nil {
		return err
	}

	// check if fields are unique
	err = service.validateFieldUniqueness(feedbackOption)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempQuestion := general.FeedbackQuestion{}
	// check if question can have options
	err = service.Repository.GetRecordForTenant(uow, feedbackOption.TenantID, &tempQuestion,
		repository.Filter("`id` = ?", feedbackOption.QuestionID), repository.Select("has_options"))
	if err != nil {
		uow.RollBack()
		return err
	}
	if !(*tempQuestion.HasOptions) {
		return errors.NewValidationError("This question cannot have options.")
	}

	err = service.Repository.Update(uow, feedbackOption)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// DeleteFeedbackOption will delete the specified feedbackQuestion
func (service *FeedbackOptionService) DeleteFeedbackOption(feedbackOption *general.FeedbackOption) error {

	// check if tenant exist
	err := service.doesTenantExist(feedbackOption.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExist(feedbackOption.TenantID, feedbackOption.DeletedBy)
	if err != nil {
		return err
	}

	// check if feedback question exist
	err = service.doesFeedbackOptionExist(feedbackOption.TenantID, feedbackOption.ID)
	if err != nil {
		return err
	}

	// check is feedback question exist
	err = service.doesFeedbackOptionExist(feedbackOption.TenantID, feedbackOption.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.UpdateWithMap(uow, general.FeedbackOption{}, map[interface{}]interface{}{
		"DeletedBy": feedbackOption.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", feedbackOption.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFeedbackOption will return specified feedback option
func (service *FeedbackOptionService) GetFeedbackOption(feedbackOption *general.FeedbackOption) error {

	// check if tenant exist
	err := service.doesTenantExist(feedbackOption.TenantID)
	if err != nil {
		return err
	}

	// check is feedback question exist
	err = service.doesFeedbackOptionExist(feedbackOption.TenantID, feedbackOption.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetForTenant(uow, feedbackOption.TenantID, feedbackOption.ID, feedbackOption)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetAllFeedbackOptions will return all the feedback options
func (service *FeedbackOptionService) GetAllFeedbackOptions(feedbackOption *[]general.FeedbackOption, tenantID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feedbackOption, "`key`")
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetFeedbackOptionsForQuestion will return all the feedback options
func (service *FeedbackOptionService) GetFeedbackOptionsForQuestion(feedbackOption *[]general.FeedbackOption,
	tenantID, questionID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if feedback question exist
	err = service.doesFeedbackQuestionExist(tenantID, questionID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feedbackOption, "`key`",
		repository.Filter("question_id=?", questionID))
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

// returns error if there is no tenant record in table.
func (service *FeedbackOptionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *FeedbackOptionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no feedback option record for the given tenant.
func (service *FeedbackOptionService) doesFeedbackOptionExist(tenantID, feedbackOptionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackOption{},
		repository.Filter("`id` = ?", feedbackOptionID))
	if err := util.HandleError("Invalid feedback option ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no feedback question record for the given tenant.
func (service *FeedbackOptionService) doesFeedbackQuestionExist(tenantID, feedbackQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackQuestion{},
		repository.Filter("`id` = ?", feedbackQuestionID))
	if err := util.HandleError("Invalid feedback question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no feedback option record for the given tenant.
func (service *FeedbackOptionService) validateFieldUniqueness(feedbackOption *general.FeedbackOption) error {

	// updating the value of an already existing type would give error as the existing record(`feedbackOption`) would
	// have a record in table which matches `key` `value` `type` fields in the table. Hence, NOT IN `id` is used to
	// exempt the `feedbackOption` record
	exists, err := repository.DoesRecordExistForTenant(service.DB, feedbackOption.TenantID, general.FeedbackOption{},
		repository.Filter("`question_id`=? AND (`key`=? OR `value`=?) AND `id`!=?", feedbackOption.QuestionID,
			feedbackOption.Key, feedbackOption.Value, feedbackOption.ID))
	if err := util.HandleIfExistsError("Record already exists with the same type and key OR value.",
		exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}
