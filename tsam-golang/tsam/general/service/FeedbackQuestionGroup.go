package service

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// FeedbackQuestionGroupService provides method to update, delete, add, get for Feedback Question Group.
type FeedbackQuestionGroupService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewFeedbackQuestionGroupService returns new instance of FeedbackQuestionGroupService.
func NewFeedbackQuestionGroupService(db *gorm.DB, repository repository.Repository) *FeedbackQuestionGroupService {
	return &FeedbackQuestionGroupService{
		DB:         db,
		Repository: repository,
	}
}

// AddFeedbackQuestionGroup adds new feedbackQuestionGroup to database.
func (service *FeedbackQuestionGroupService) AddFeedbackQuestionGroup(feedbackQuestionGroup *general.FeedbackQuestionGroup, uows ...*repository.UnitOfWork) error {
	// Validate tenant id.
	err := service.doesTenantExist(feedbackQuestionGroup.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate if same city name exists.
	err = service.doesSameFieldConflictExist(feedbackQuestionGroup)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(feedbackQuestionGroup.CreatedBy, feedbackQuestionGroup.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	//  Creating unit of work.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Add feedbackQuestionGroup to database.
	err = service.Repository.Add(uow, feedbackQuestionGroup)
	if err != nil {
		uow.RollBack()
		return err
	}

	// If calling function does not pass the uow only then make commit.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddFeedbackQuestionGroups adds multiple feedbackQuestionGroups.
func (service *FeedbackQuestionGroupService) AddFeedbackQuestionGroups(feedbackQuestionGroups *[]general.FeedbackQuestionGroup,
	feedbackQuestionGroupIDs *[]uuid.UUID, tenantID, credentialID uuid.UUID) error {
	// Check for same group name conflict.
	for i := 0; i < len(*feedbackQuestionGroups); i++ {
		for j := 0; j < len(*feedbackQuestionGroups); j++ {
			if i != j && (*feedbackQuestionGroups)[i].GroupName == (*feedbackQuestionGroups)[j].GroupName {
				log.NewLogger().Error("GroupName:" + (*feedbackQuestionGroups)[j].GroupName + " exists")
				return errors.NewValidationError("GroupName:" + (*feedbackQuestionGroups)[j].GroupName + " exists")
			}
		}
	}

	// Add individual feedbackQuestionGroup.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, feedbackQuestionGroup := range *feedbackQuestionGroups {
		feedbackQuestionGroup.TenantID = tenantID
		feedbackQuestionGroup.CreatedBy = credentialID
		err := service.AddFeedbackQuestionGroup(&feedbackQuestionGroup, uow)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
		*feedbackQuestionGroupIDs = append(*feedbackQuestionGroupIDs, feedbackQuestionGroup.ID)
	}

	uow.Commit()
	return nil
}

// UpdateFeedbackQuestionGroup updates feedbackQuestionGroup in database.
func (service *FeedbackQuestionGroupService) UpdateFeedbackQuestionGroup(feedbackQuestionGroup *general.FeedbackQuestionGroup) error {
	// Validate tenant ID.
	err := service.doesTenantExist(feedbackQuestionGroup.TenantID)
	if err != nil {
		return err
	}

	// Validate feedbackQuestionGroup ID.
	err = service.doesFeedbackQuestionGroupExist(feedbackQuestionGroup.ID, feedbackQuestionGroup.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(feedbackQuestionGroup.UpdatedBy, feedbackQuestionGroup.TenantID)
	if err != nil {
		return err
	}

	// Validate if same city name exists.
	err = service.doesSameFieldConflictExist(feedbackQuestionGroup)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update feedbackQuestionGroup.
	err = service.Repository.Update(uow, feedbackQuestionGroup)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteFeedbackQuestionGroup deletes feedbackQuestionGroup from database.
func (service *FeedbackQuestionGroupService) DeleteFeedbackQuestionGroup(feedbackQuestionGroup *general.FeedbackQuestionGroup) error {
	credentialID := feedbackQuestionGroup.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(feedbackQuestionGroup.TenantID)
	if err != nil {
		return err
	}

	// Validate feedbackQuestionGroup ID.
	err = service.doesFeedbackQuestionGroupExist(feedbackQuestionGroup.ID, feedbackQuestionGroup.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, feedbackQuestionGroup.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update feedbackQuestionGroup for updating deleted_by and deleted_at fields of feedbackQuestionGroup
	if err := service.Repository.UpdateWithMap(uow, feedbackQuestionGroup, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", feedbackQuestionGroup.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("FeedbackQuestionGroup could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetFeedbackQuestionGroup returns one feedbackQuestionGroup.
func (service *FeedbackQuestionGroupService) GetFeedbackQuestionGroup(feedbackQuestionGroup *general.FeedbackQuestionGroupDTO,
	tenantID, feedbackQuestionGroupID uuid.UUID) error {
	// Validate tenant ID
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	// validate feedbackQuestionGroup ID
	exists, err := repository.DoesRecordExist(service.DB, feedbackQuestionGroup,
		repository.Filter("`id` = ?", feedbackQuestionGroupID))
	if err != nil {
		return err
	}
	if !exists {
		return errors.NewValidationError("FeedbackQuestionGroup not found")
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetForTenant(uow, tenantID, feedbackQuestionGroupID, feedbackQuestionGroup)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetFeedbackQuestionGroups returns all feedbackQuestionGroups.
func (service *FeedbackQuestionGroupService) GetFeedbackQuestionGroups(feedbackQuestionGroup *[]general.FeedbackQuestionGroupDTO, tenantID uuid.UUID,
	form url.Values, limit, offset int, totalCount *int) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all feedbackQuestionGroups from database.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feedbackQuestionGroup, "`order`",
		service.addSearchQueries(form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFeedbackQuestionGroupList returns feedbackQuestionGroup list.
func (service *FeedbackQuestionGroupService) GetFeedbackQuestionGroupList(feedbackQuestionGroup *[]general.FeedbackQuestionGroupDTO,
	tenantID uuid.UUID) error {
	// Validate tenant ID
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get feedbackQuestionGroup list.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feedbackQuestionGroup, "`order`")
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFeedbackQuestionGroupListByType returns feedbackQuestionGroup list.
func (service *FeedbackQuestionGroupService) GetFeedbackQuestionGroupListByType(feedbackQuestionGroup *[]general.FeedbackQuestionGroupDTO,
	tenantID uuid.UUID, questionType string) error {
	// Validate tenant ID
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get feedbackQuestionGroup list.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feedbackQuestionGroup, "`order`",
		repository.Filter("`type` = ?", questionType))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFeedbackQuestionGroupByType will return group-wise feedback question based on question-type
func (service *FeedbackQuestionGroupService) GetFeedbackQuestionGroupByType(feedbackQuestionGroup *[]general.FeedbackQuestionGroupDTO,
	tenantID uuid.UUID, questionType string) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get feedbackQuestionGroup by question type.
	// New push
	// err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feedbackQuestionGroup, "`order`",
	// 	repository.Filter("`type` = ?", questionType),
	// 	repository.PreloadWithCustomCondition(map[string][]repository.QueryProcessor{
	// 		"FeedbackQuestions": {
	// 			repository.OrderBy("feedback_questions.`order`"),
	// 			repository.Filter("feedback_questions.`is_active` = '1'"),
	// 		},
	// 	}),
	// 	repository.PreloadWithCustomCondition(map[string][]repository.QueryProcessor{
	// 		"FeedbackQuestions.Options": {repository.OrderBy("feedback_options.`order`")},
	// 	}))
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }

	//  #shailesh, test this #niranjan.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feedbackQuestionGroup, "`order`",
		repository.Filter("`type` = ?", questionType),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "FeedbackQuestions",
			Queryprocessors: []repository.QueryProcessor{
				repository.OrderBy("feedback_questions.`order`"),
				repository.Filter("feedback_questions.`is_active` = '1'")}},
			repository.Preload{
				Schema:          "FeedbackQuestions.Options",
				Queryprocessors: []repository.QueryProcessor{repository.OrderBy("feedback_options.`order`")},
			}))
	if err != nil {
		uow.RollBack()
		return err
	}

	// // new preload with custom condition
	// err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feedbackQuestionGroup, "`order`",
	// 	repository.Filter("`type` = ?", questionType),
	// 	repository.PreloadForCustomCondition("FeedbackQuestions", []repository.QueryProcessor{
	// 		repository.Filter("feedback_questions.`is_active` = '1'"), repository.OrderBy("feedback_questions.`order`")}),
	// 	repository.PreloadForCustomCondition("FeedbackQuestions.Options",
	// 		[]repository.QueryProcessor{repository.OrderBy("feedback_options.`order`")}))
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }

	for index := range *feedbackQuestionGroup {
		// get max score.
		err = service.Repository.Scan(uow, &(*feedbackQuestionGroup)[index],
			repository.Model(&general.FeedbackOption{}),
			repository.Select("MAX(feedback_options.`key`) AS max_score"),
			repository.Join("INNER JOIN `feedback_questions` ON "+
				"`feedback_options`.`question_id` = `feedback_questions`.`id`"),
			repository.Filter("feedback_questions.`group_id` = ?", (*feedbackQuestionGroup)[index].ID),
			repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
			repository.Filter("feedback_options.`tenant_id`=? AND feedback_options.`deleted_at` IS NULL", tenantID))
		if err != nil {
			uow.RollBack()
			return err
		}

		// get min score.
		err = service.Repository.Scan(uow, &(*feedbackQuestionGroup)[index],
			repository.Model(&general.FeedbackOption{}),
			repository.Select("MIN(feedback_options.`key`) AS min_score"),
			repository.Join("INNER JOIN `feedback_questions` ON "+
				"`feedback_options`.`question_id` = `feedback_questions`.`id`"),
			repository.Filter("feedback_questions.`group_id` = ?", (*feedbackQuestionGroup)[index].ID),
			repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
			repository.Filter("feedback_options.`tenant_id`=? AND feedback_options.`deleted_at` IS NULL", tenantID))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetGroupwiseKeywordName will get all the feedback question groups for faculty-session-feedback.
func (service *FeedbackQuestionGroupService) GetGroupwiseKeywordName(tenantID uuid.UUID,
	keywordNames *[]dashboard.GroupWiseKeywordName) error {

	uow := repository.NewUnitOfWork(service.DB, true)

	tempFeedbackQuestionGroup := []general.FeedbackQuestionGroup{}

	// Get feedback question keyword group wise
	err := service.Repository.GetAllInOrderForTenant(uow, tenantID, &tempFeedbackQuestionGroup, "`order`",
		repository.Table("feedback_question_groups"), repository.Select("`id`, `group_name`"),
		repository.Filter("`type` = ?", "Faculty_Session_Feedback"))
	if err != nil {
		uow.RollBack()
		return err
	}

	// tempGroupwiseKeywords := make([]dashboard.GroupWiseKeywordName, len(tempFeedbackQuestionGroup))
	for index := range tempFeedbackQuestionGroup {

		tempGroupwiseKeywords := dashboard.GroupWiseKeywordName{}
		tempGroupwiseKeywords.GroupName = (tempFeedbackQuestionGroup)[index].GroupName

		err = service.getKeywords(uow, &tempGroupwiseKeywords.Keywords,
			tenantID, (tempFeedbackQuestionGroup)[index].ID)
		if err != nil {
			return err
		}

		*keywordNames = append(*keywordNames, tempGroupwiseKeywords)
	}

	return nil
}

// getKeywords will return keyword names from feedback_questions table
func (service *FeedbackQuestionGroupService) getKeywords(uow *repository.UnitOfWork,
	keywords *[]dashboard.KeywordName, tenantID uuid.UUID, groupID ...uuid.UUID) error {

	var queryProcessor repository.QueryProcessor

	if groupID != nil {
		queryProcessor = repository.Filter("`group_id` = ?", groupID)
	}

	err := service.Repository.GetAllForTenant(uow, tenantID, keywords,
		repository.Table("feedback_questions"), repository.Select("feedback_questions.`keyword` AS name"),
		repository.Filter("feedback_questions.`has_options` = true AND feedback_questions.`type`='Faculty_Session_Feedback'"),
		repository.Filter("feedback_questions.`is_active` = '1'"), queryProcessor, repository.OrderBy("feedback_questions.`order`"))
	if err != nil {
		return err
	}

	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// addSearchQueries adds search queries.
func (service *FeedbackQuestionGroupService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["groupName"]; ok {
		util.AddToSlice("`group_name`", "LIKE ?", "AND", "%"+requestForm.Get("groupName")+"%", &columnNames, &conditions, &operators, &values)
	}
	if _, ok := requestForm["groupDescription"]; ok {
		util.AddToSlice("`group_description`", "LIKE ?", "AND", "%"+requestForm.Get("groupDescription")+"%", &columnNames, &conditions, &operators, &values)
	}
	if _, ok := requestForm["order"]; ok {
		orderInNumber, _ := strconv.Atoi(requestForm.Get("order"))
		util.AddToSlice("`order`", "=?", "AND", orderInNumber, &columnNames, &conditions, &operators, &values)
	}
	if questionType, ok := requestForm["type"]; ok {
		util.AddToSlice("`type`", "= ?", "AND", questionType, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesSameFieldConflictExist returns true if the order and group name already exist for the feedbackQuestionGroup in database.
func (service *FeedbackQuestionGroupService) doesSameFieldConflictExist(feedbackQuestionGroup *general.FeedbackQuestionGroup) error {
	// Check for same group name and conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, feedbackQuestionGroup.TenantID, &general.FeedbackQuestionGroup{},
		repository.Filter("(`order`=? OR `group_name`=?) AND `type` = ? AND `id`!=?", feedbackQuestionGroup.Order, feedbackQuestionGroup.GroupName,
			feedbackQuestionGroup.Type, feedbackQuestionGroup.ID))
	if err := util.HandleIfExistsError("Order: "+strconv.FormatUint(uint64(feedbackQuestionGroup.Order), 10)+" Group name: "+
		feedbackQuestionGroup.GroupName+" for Type: "+feedbackQuestionGroup.Type+" exists", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *FeedbackQuestionGroupService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesFeedbackQuestionGroupExist validates if feedbackQuestionGroup exists or not in database.
func (service *FeedbackQuestionGroupService) doesFeedbackQuestionGroupExist(feedbackQuestionGroupID uuid.UUID, tenantID uuid.UUID) error {
	// Check feedbackQuestionGroup exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackQuestionGroup{},
		repository.Filter("`id` = ?", feedbackQuestionGroupID))
	if err := util.HandleError("Invalid feedbackQuestionGroup ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *FeedbackQuestionGroupService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
