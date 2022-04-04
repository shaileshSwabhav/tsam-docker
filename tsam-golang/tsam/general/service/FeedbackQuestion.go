package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// FeedbackQuestionService provides methods to do different CRUD operations on feedback_questions table.
type FeedbackQuestionService struct {
	DB              *gorm.DB
	Repository      repository.Repository
	association     []string
	maxRating       uint
	feedbackTypeMap map[string]interface{}
	feedbackType    []string
}

// NewFeedbackQuestionService returns a new instance Of FeedbackQuestionService.
func NewFeedbackQuestionService(db *gorm.DB, repository repository.Repository) *FeedbackQuestionService {
	return &FeedbackQuestionService{
		DB:         db,
		Repository: repository,
		association: []string{
			"Options", "FeedbackQuestionGroup",
		},
		maxRating: 10,
		feedbackTypeMap: map[string]interface{}{
			"Faculty_Batch_Feedback":   bat.FacultyTalentFeedback{},
			"Faculty_Session_Feedback": bat.FacultyTalentBatchSessionFeedback{},
			"Talent_Batch_Feedback":    bat.TalentFeedback{},
			"Talent_Session_Feedback":  bat.TalentBatchSessionFeedback{},
			"Faculty_Assessment":       faculty.FacultyAssessment{},
			"Aha_Moment_Feedback":      bat.AhaMomentResponse{},
		},
		feedbackType: []string{
			"Faculty_Batch_Feedback", "Faculty_Session_Feedback", "Talent_Batch_Feedback",
			"Talent_Session_Feedback", "Faculty_Assessment", "Aha_Moment_Feedback",
		},
	}
}

// AddFeedbackQuestion will add feedback questions to the table
func (service *FeedbackQuestionService) AddFeedbackQuestion(feedbackQuestion *general.FeedbackQuestion, uows ...*repository.UnitOfWork) error {

	// Check if max rating exist.
	if *feedbackQuestion.HasOptions && (feedbackQuestion.Type == service.feedbackType[0] ||
		feedbackQuestion.Type == service.feedbackType[1] || feedbackQuestion.Type == service.feedbackType[5]) {
		err := service.doesMaxRatingExist(feedbackQuestion)
		if err != nil {
			return err
		}
	}

	// Check if all foreign key exist.
	err := service.doesForeignKeyExist(feedbackQuestion, feedbackQuestion.CreatedBy)
	if err != nil {
		return err
	}

	if feedbackQuestion.Options != nil {
		for index := range feedbackQuestion.Options {
			feedbackQuestion.Options[index].TenantID = feedbackQuestion.TenantID
			feedbackQuestion.Options[index].CreatedBy = feedbackQuestion.CreatedBy
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

	err = service.Repository.Add(uow, feedbackQuestion)
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

// AddFeedbackQuestions will add multiple feedback questions to the table
func (service *FeedbackQuestionService) AddFeedbackQuestions(feedbackQuestions *[]general.FeedbackQuestion, tenantID,
	credentialID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, false)
	for _, feedbackQuestion := range *feedbackQuestions {
		feedbackQuestion.TenantID = tenantID
		feedbackQuestion.CreatedBy = credentialID
		err := service.AddFeedbackQuestion(&feedbackQuestion)
		if err != nil {
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// UpdateFeedbackQuestion will update the specified feedbackQuestion
func (service *FeedbackQuestionService) UpdateFeedbackQuestion(feedbackQuestion *general.FeedbackQuestion) error {

	if !(*feedbackQuestion.IsActive) {
		return errors.NewValidationError("Inactive questions cannot be updated.")
	}

	// Check if max rating exist.
	if *feedbackQuestion.HasOptions && (feedbackQuestion.Type == service.feedbackType[0] ||
		feedbackQuestion.Type == service.feedbackType[1] || feedbackQuestion.Type == service.feedbackType[5]) {
		err := service.doesMaxRatingExist(feedbackQuestion)
		if err != nil {
			return err
		}
	}

	// Check if all foreign key exist.
	err := service.doesForeignKeyExist(feedbackQuestion, feedbackQuestion.UpdatedBy)
	if err != nil {
		return err
	}

	// check is feedback question exist
	err = service.doesFeedbackQuestionExist(feedbackQuestion.TenantID, feedbackQuestion.ID)
	if err != nil {
		return err
	}

	// check if feedback question is already answered.
	err = service.doesFeedbackQuestionAnswerExist(feedbackQuestion.TenantID, feedbackQuestion.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.updateOptions(uow, feedbackQuestion, feedbackQuestion.UpdatedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.Update(uow, feedbackQuestion)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// UpdateFeedbackQuestionStatus will update the status of the specified feedback question.
func (service *FeedbackQuestionService) UpdateFeedbackQuestionStatus(feedbackQuestion *general.FeedbackQuestion) error {

	// Check if tenant exist.
	err := service.doesTenantExist(feedbackQuestion.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(feedbackQuestion.TenantID, feedbackQuestion.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if feedback question exist.
	err = service.doesFeedbackQuestionExist(feedbackQuestion.TenantID, feedbackQuestion.ID)
	if err != nil {
		return err
	}

	if *feedbackQuestion.IsActive {

		// Check if feedback question and keyword is unique for its type.
		err = service.validateQuestionUniqueness(feedbackQuestion)
		if err != nil {
			return err
		}

		// Check if active question with same order already exist.
		if feedbackQuestion.FeedbackQuestionGroup == nil {
			// check if same order exist.
			err = service.doesQuestionOrderExist(feedbackQuestion.TenantID, feedbackQuestion.ID,
				feedbackQuestion.Order, feedbackQuestion.Type)
			if err != nil {
				return err
			}
		} else {
			// check if same order exist.
			err = service.doesFeedbackQuestionOrderExistForGroup(feedbackQuestion.TenantID, feedbackQuestion.ID,
				feedbackQuestion.FeedbackQuestionGroup.ID, feedbackQuestion.Order, feedbackQuestion.Type)
			if err != nil {
				return err
			}
			// feedbackQuestion.GroupID = &feedbackQuestion.FeedbackQuestionGroup.ID
		}
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, general.FeedbackQuestion{}, map[string]interface{}{
		"IsActive":  feedbackQuestion.IsActive,
		"UpdatedBy": feedbackQuestion.UpdatedBy,
		"UpdatedAt": time.Now(),
	}, repository.Filter("`id` = ?", feedbackQuestion.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteFeedbackQuestion will delete the specified feedbackQuestion.
func (service *FeedbackQuestionService) DeleteFeedbackQuestion(feedbackQuestion *general.FeedbackQuestion) error {

	// check if tenant exist.
	err := service.doesTenantExist(feedbackQuestion.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(feedbackQuestion.TenantID, feedbackQuestion.DeletedBy)
	if err != nil {
		return err
	}

	// check is feedback question exist.
	err = service.doesFeedbackQuestionExist(feedbackQuestion.TenantID, feedbackQuestion.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.UpdateWithMap(uow, general.FeedbackQuestion{}, map[string]interface{}{
		"DeletedBy": feedbackQuestion.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", feedbackQuestion.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.deleteFeedbackOption(uow, feedbackQuestion.DeletedBy, feedbackQuestion.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.deletedAllFeedback(uow, feedbackQuestion.TenantID, feedbackQuestion.DeletedBy,
		feedbackQuestion.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFeedbackQuestion will return the specified feedback question.
func (service *FeedbackQuestionService) GetFeedbackQuestion(feedbackQuestion *general.FeedbackQuestion) error {

	// check if tenant exist.
	err := service.doesTenantExist(feedbackQuestion.TenantID)
	if err != nil {
		return err
	}

	// check is feedback question exist.
	err = service.doesFeedbackQuestionExist(feedbackQuestion.TenantID, feedbackQuestion.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetForTenant(uow, feedbackQuestion.TenantID, feedbackQuestion.ID, feedbackQuestion,
		repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFeedbackQuestions will return all the feedback question from the table.
func (service *FeedbackQuestionService) GetFeedbackQuestions(feedbackQuestions *[]general.FeedbackQuestionDTO, form url.Values, tenantID uuid.UUID,
	limit, offset int, totalCount *int) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// , repository.PreloadAssociations(service.association)
	err = service.Repository.GetAllForTenant(uow, tenantID, feedbackQuestions,
		service.addSearchQueries(form), repository.PreloadAssociations([]string{"FeedbackQuestionGroup"}),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "Options",
			Queryprocessors: []repository.QueryProcessor{
				repository.OrderBy("feedback_options.`order`")},
		}), repository.OrderBy("`is_active` DESC, `type`, `order`"), repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllFeedbackQuestion will return all the feedback question from the table.
func (service *FeedbackQuestionService) GetAllFeedbackQuestion(feedbackQuestions *[]general.FeedbackQuestionDTO, tenantID uuid.UUID) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feedbackQuestions, "`order`",
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "Options",
			Queryprocessors: []repository.QueryProcessor{
				repository.OrderBy("feedback_options.`order`")},
		}), repository.OrderBy("`is_active` DESC"))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFeedbackQuestionsByType will return all the feedback question from the table.
func (service *FeedbackQuestionService) GetFeedbackQuestionsByType(feedbackQuestions *[]general.FeedbackQuestionDTO, tenantID uuid.UUID,
	typeName string) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	// repository.PreloadAssociations(service.association)
	// repository.Join("LEFT JOIN feedback_options ON feedback_questions.id = feedback_options.question_id"),
	// repository.OrderBy("feedback_options.key")
	err = service.Repository.GetAll(uow, feedbackQuestions,
		repository.Filter("feedback_questions.`type`=? AND feedback_questions.`is_active` = true", typeName),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "FeedbackQuestionGroup",
			Queryprocessors: []repository.QueryProcessor{
				repository.OrderBy("feedback_question_groups.`order`")},
		}),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "Options",
			Queryprocessors: []repository.QueryProcessor{
				repository.OrderBy("feedback_options.`order`")},
		}), repository.Join("LEFT JOIN feedback_question_groups ON feedback_question_groups.`id` = feedback_questions.`group_id`"+
			" AND feedback_question_groups.`tenant_id` = feedback_questions.`tenant_id`"),
		repository.Filter("feedback_questions.`tenant_id`=?", tenantID),
		repository.Filter("feedback_question_groups.`deleted_at` IS NULL"),
		repository.OrderBy("feedback_questions.`order`, feedback_question_groups.`order`"))
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

// doesForeignKeyExist will check if all foreign keys are valid.
func (service *FeedbackQuestionService) doesForeignKeyExist(feedbackQuestion *general.FeedbackQuestion, credentialID uuid.UUID) error {

	// check if tenant exist.
	err := service.doesTenantExist(feedbackQuestion.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(feedbackQuestion.TenantID, credentialID)
	if err != nil {
		return err
	}

	if feedbackQuestion.FeedbackQuestionGroup == nil {
		// check if same order exist.
		err = service.doesQuestionOrderExist(feedbackQuestion.TenantID, feedbackQuestion.ID,
			feedbackQuestion.Order, feedbackQuestion.Type)
		if err != nil {
			return err
		}
	} else {
		// check if same order exist.
		err = service.doesFeedbackQuestionOrderExistForGroup(feedbackQuestion.TenantID, feedbackQuestion.ID,
			feedbackQuestion.FeedbackQuestionGroup.ID, feedbackQuestion.Order, feedbackQuestion.Type)
		if err != nil {
			return err
		}
		feedbackQuestion.GroupID = &feedbackQuestion.FeedbackQuestionGroup.ID
	}

	// check order field in options json.
	err = service.checkValidOptionFieldsInJSON(&feedbackQuestion.Options)
	if err != nil {
		return err
	}

	// check if fields are unique.
	err = service.validateQuestionUniqueness(feedbackQuestion)
	if err != nil {
		return err
	}

	if feedbackQuestion.Options != nil {
		for index := range feedbackQuestion.Options {
			// check if fields are unique in options.
			err = service.validateOptionsFieldUniqueness(feedbackQuestion.Options[index],
				feedbackQuestion.TenantID, feedbackQuestion.ID)
			if err != nil {
				return err
			}

			service.calculateMaxScore(feedbackQuestion)
		}
	}

	return nil
}

// doesMaxRatingExist checks if one and only key with 10 exist.
func (service *FeedbackQuestionService) doesMaxRatingExist(feedbackQuestion *general.FeedbackQuestion) error {

	if feedbackQuestion.Options != nil {
		optionsMap := make(map[int]int)

		for _, option := range feedbackQuestion.Options {
			if option.Key > int(service.maxRating) {
				return errors.NewValidationError("Score cannot be greater than 10.")
			}
			optionsMap[option.Key]++
			// if (optionsMap[option.Key] > 1) {
			// 	return errors.NewValidationError("Multiple options cannot have same score")
			// }
		}

		if optionsMap[int(service.maxRating)] == 0 {
			return errors.NewValidationError("Score with 10 must be specified.")
		}
		if optionsMap[int(service.maxRating)] > 1 {
			return errors.NewValidationError("Only 1 option with score 10 must be present.")
		}
	}

	return nil
}

// updateOptions will check the new options with the options already in table.
func (service *FeedbackQuestionService) updateOptions(uow *repository.UnitOfWork,
	feedbackQuestion *general.FeedbackQuestion, credentialID uuid.UUID) error {

	if feedbackQuestion.Options == nil {
		err := service.Repository.UpdateWithMap(uow, general.FeedbackOption{}, map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		}, repository.Filter("`question_id`=?", feedbackQuestion.ID))
		if err != nil {
			return err
		}
		return nil
	}

	feedbackOptions := feedbackQuestion.Options
	tempFeedbackOptions := []general.FeedbackOption{}

	err := service.Repository.GetAllForTenant(uow, feedbackQuestion.TenantID, &tempFeedbackOptions,
		repository.Filter("`question_id`=?", feedbackQuestion.ID))
	if err != nil {
		return err
	}

	feedbackOptionMap := make(map[uuid.UUID]uint)

	for _, tempFeedbackOption := range tempFeedbackOptions {
		feedbackOptionMap[tempFeedbackOption.ID]++
	}

	for _, feedbackOption := range feedbackOptions {

		if util.IsUUIDValid(feedbackOption.ID) {
			feedbackOptionMap[feedbackOption.ID]++
		} else {
			feedbackOption.CreatedBy = credentialID
			feedbackOption.TenantID = feedbackQuestion.TenantID
			feedbackOption.QuestionID = feedbackQuestion.ID
			err = service.Repository.Add(uow, &feedbackOption)
			if err != nil {
				return err
			}
		}

		if feedbackOptionMap[feedbackOption.ID] > 1 {
			feedbackOption.UpdatedBy = credentialID
			err = service.Repository.Update(uow, &feedbackOption)
			if err != nil {
				return err
			}
			feedbackOptionMap[feedbackOption.ID] = 0
		}
	}

	for _, tempFeedbackOption := range tempFeedbackOptions {
		if feedbackOptionMap[tempFeedbackOption.ID] == 1 {
			// tempFeedbackOption.DeletedBy =
			err = service.Repository.UpdateWithMap(uow, general.FeedbackOption{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempFeedbackOption.ID))
			if err != nil {
				return err
			}
			feedbackOptionMap[tempFeedbackOption.ID] = 0
		}
	}
	feedbackQuestion.Options = nil
	return nil
}

// doesFeedbackQuestionAnswerExist check if feedback question is already answered.
func (service *FeedbackQuestionService) doesFeedbackQuestionAnswerExist(tenantID, questionID uuid.UUID) error {

	for _, tableName := range service.feedbackTypeMap {
		exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tableName,
			repository.Filter("`question_id` = ?", questionID))
		if err != nil {
			return err
		}
		if exist {
			return errors.NewValidationError("Current feedback question cannot be updated has it is already answered.")
		}
	}

	return nil
}

// deleteFeedbackOption will delete options for specified question
func (service *FeedbackQuestionService) deleteFeedbackOption(uow *repository.UnitOfWork,
	credentialID, feedbackQuestionID uuid.UUID) error {

	err := service.Repository.UpdateWithMap(uow, general.FeedbackOption{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`question_id`=?", feedbackQuestionID))
	if err != nil {
		return err
	}
	return nil
}

// deletedAllFeedback will delete the specified question from all the feedback tables(faculty & talent)
func (service *FeedbackQuestionService) deletedAllFeedback(uow *repository.UnitOfWork,
	tenantID, credentialID, questionID uuid.UUID) error {

	// delete question from faculty_talent_feedback_batch_session_feedback
	err := service.Repository.UpdateWithMap(uow, bat.FacultyTalentBatchSessionFeedback{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`tenant_id`=? AND `question_id`=?", tenantID, questionID))
	if err != nil {
		return err
	}

	// delete question from faculty_talent_feedback_batch_feedback
	err = service.Repository.UpdateWithMap(uow, bat.FacultyTalentFeedback{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`tenant_id`=? AND `question_id`=?", tenantID, questionID))
	if err != nil {
		return err
	}

	// delete question from talent_batch_feedback_batch_feedback
	err = service.Repository.UpdateWithMap(uow, bat.TalentFeedback{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`tenant_id`=? AND `question_id`=?", tenantID, questionID))
	if err != nil {
		return err
	}

	// delete question from talent_batch_session_feedback_batch_feedback
	err = service.Repository.UpdateWithMap(uow, bat.TalentBatchSessionFeedback{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`tenant_id`=? AND `question_id`=?", tenantID, questionID))
	if err != nil {
		return err
	}
	return nil
}

// calculateMaxScore will calculate maxScore from all the options
func (service *FeedbackQuestionService) calculateMaxScore(feedbackQuestion *general.FeedbackQuestion) {

	maxScore := 0
	feedbackQuestion.MaxScore = &maxScore

	for _, option := range feedbackQuestion.Options {
		if option.Key > *feedbackQuestion.MaxScore {
			*feedbackQuestion.MaxScore = option.Key
		}
	}
}

// returns error if there is no tenant record in table.
func (service *FeedbackQuestionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *FeedbackQuestionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no feedback question record for the given tenant.
func (service *FeedbackQuestionService) doesFeedbackQuestionExist(tenantID, feedbackQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackQuestion{},
		repository.Filter("`id` = ?", feedbackQuestionID))
	if err := util.HandleError("Invalid feedback question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if order already exist.
func (service *FeedbackQuestionService) doesQuestionOrderExist(tenantID, feedbackQuestionID uuid.UUID, order uint,
	questionType string) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackQuestion{},
		repository.Filter("`order`=? AND `type`=? AND `id`!=? AND `is_active` = true", order, questionType, feedbackQuestionID))
	if err := util.HandleIfExistsError("Order already exist", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if order already exist for specified group.
func (service *FeedbackQuestionService) doesFeedbackQuestionOrderExistForGroup(tenantID, feedbackQuestionID, groupID uuid.UUID,
	order uint, questionType string) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackQuestion{},
		repository.Filter("`order`=? AND `type`=? AND `id`!=? AND `group_id`=? AND `is_active`=true",
			order, questionType, feedbackQuestionID, groupID))
	if err := util.HandleIfExistsError("Order already exist", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is feedback question/keyword record for the given tenant.
func (service *FeedbackQuestionService) validateQuestionUniqueness(feedbackQuestion *general.FeedbackQuestion) error {

	if feedbackQuestion.FeedbackQuestionGroup != nil {
		exists, err := repository.DoesRecordExistForTenant(service.DB, feedbackQuestion.TenantID, general.FeedbackQuestion{},
			repository.Filter("`type`=? AND (`question`=? || `keyword`=?) AND `id`!=? AND `is_active`=true AND `group_id`=?",
				feedbackQuestion.Type, feedbackQuestion.Question, feedbackQuestion.Keyword,
				feedbackQuestion.ID, feedbackQuestion.FeedbackQuestionGroup.ID))
		if err := util.HandleIfExistsError("Record already exists with the same type and question for group.",
			exists, err); err != nil {
			return errors.NewValidationError(err.Error())
		}
		return nil
	}

	exists, err := repository.DoesRecordExistForTenant(service.DB, feedbackQuestion.TenantID, general.FeedbackQuestion{},
		repository.Filter("`type`=? AND (`question`=? || `keyword`=?) AND `id`!=? AND `is_active`=true",
			feedbackQuestion.Type, feedbackQuestion.Question, feedbackQuestion.Keyword, feedbackQuestion.ID))
	if err := util.HandleIfExistsError("Record already exists with the same type and question.",
		exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}

	return nil
}

// returns error if there is no feedback option record for the given tenant.
func (service *FeedbackQuestionService) validateOptionsFieldUniqueness(feedbackOption general.FeedbackOption,
	tenantID, questionID uuid.UUID) error {

	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackOption{},
		repository.Filter("`question_id`=? AND `key`=? AND `value`=? AND `id`!=?", questionID,
			feedbackOption.Key, feedbackOption.Value, feedbackOption.ID))
	if err := util.HandleIfExistsError("Record already exists with the same key and value.",
		exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

// checkValidOptionFieldsInJSON will check if json consist repeated order
func (service *FeedbackQuestionService) checkValidOptionFieldsInJSON(feedbackOptions *[]general.FeedbackOption) error {

	feedbackOptionMap := make(map[uint]uint)
	for _, feedbackOption := range *feedbackOptions {
		feedbackOptionMap[feedbackOption.Order]++
		if feedbackOptionMap[feedbackOption.Order] > 1 {
			return errors.NewValidationError("Same order cannot be repeated for multiple orders")
		}
	}
	return nil
}

// addSearchQueries will append search queries from queryParams to queryProcessor
func (service *FeedbackQuestionService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["keyword"]; ok {
		util.AddToSlice("`keyword`", "LIKE ?", "AND", "%"+requestForm.Get("keyword")+"%", &columnNames, &conditions, &operators, &values)
	}

	if questionType, ok := requestForm["questionType"]; ok {
		util.AddToSlice("`type`", "=?", "AND", questionType, &columnNames, &conditions, &operators, &values)
	}

	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("`is_active`", "=?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}
	if groupID, ok := requestForm["groupID"]; ok {
		util.AddToSlice("`group_id`", "=?", "AND", groupID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
