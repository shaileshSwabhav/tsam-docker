package service

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// DiscussionService Give Access to Update, Add, Delete Discussion
type DiscussionService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewDiscussionService returns new instance of DiscussionService
func NewDiscussionService(db *gorm.DB, repo repository.Repository) *DiscussionService {
	return &DiscussionService{
		DB:           db,
		Repository:   repo,
		associations: []string{"Channel", "Author", "BestReply"},
	}
}

// AddDiscussion Add New Discussion
func (service *DiscussionService) AddDiscussion(discussion *community.Discussion) error {

	// check if all foreign keys exist.
	err := service.doForeignKeysExist(discussion.CreatedBy, discussion)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(discussion)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.Add(uow, discussion)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// UpdateDiscussion Update Discussion
func (service *DiscussionService) UpdateDiscussion(discussion *community.Discussion) error {

	// Check if discussion exist.
	err := service.doesDiscussionExist(discussion.TenantID, discussion.ID)
	if err != nil {
		return err
	}

	// check if all foreign keys exist.
	err = service.doForeignKeysExist(discussion.UpdatedBy, discussion)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(discussion)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Update(uow, discussion)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("An error occured", http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// DeleteDiscussion Delete Discussion By Discussion id
func (service *DiscussionService) DeleteDiscussion(discussion *community.Discussion) error {

	// Check if discussion exist.
	err := service.doesDiscussionExist(discussion.TenantID, discussion.ID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(discussion.TenantID, discussion.DeletedBy); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.UpdateWithMap(uow, discussion, map[string]interface{}{
		"DeletedBy": discussion.DeletedBy,
		"DeletedAt": time.Now(),
	})
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("unable to delete discussions", http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetDiscussions gets particular discussions.
func (service *DiscussionService) GetDiscussions(tenantID uuid.UUID,
	discussions *[]community.DiscussionDTO, parser *web.Parser, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, discussions, "`created_at`",
		service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount),
		repository.PreloadAssociations(service.associations))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("unable to get discussions", http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetDiscussion Get Discussion By Discussion id
func (service *DiscussionService) GetDiscussion(discussion *community.DiscussionDTO,
	queryProcessor ...repository.QueryProcessor) error {
	tenantID := discussion.TenantID

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetRecordForTenant(uow, tenantID, discussion,
		repository.Filter("`id` = ?", discussion.ID), repository.PreloadAssociations(service.associations))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// // GetDiscussionsByChannelID Get Discussion By Discussion id
// func (service *DiscussionService) GetDiscussionsByChannelID(allDiscussions *[]community.Discussion,
// 	channelID uuid.UUID, queryProcessor ...repository.QueryProcessor) error {

// 	// // Check if channel exist
// 	// err := service.doesChannelExist(channelID)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	uow := repository.NewUnitOfWork(service.DB, true)

// 	queryProcessor = append(queryProcessor, repository.Filter("`channel_id` =?", channelID),
// 		repository.PreloadAssociations(service.associations))
// 	queryProcessor = append(queryProcessor, repository.OrderBy("`created_at`", true))

// 	err := service.Repository.GetAll(uow, allDiscussions, queryProcessor...)
// 	if err != nil {
// 		uow.RollBack()
// 		return errors.NewHTTPError("unable to get discussions", http.StatusInternalServerError)
// 	}
// 	for index, discussion := range *allDiscussions {
// 		err = service.setHasBestReply(&discussion)
// 		if err != nil {
// 			return err
// 		}
// 		(*allDiscussions)[index] = discussion
// 	}
// 	uow.Commit()
// 	return nil
// }

// // GetDiscussionsByTalentID Get Discussion By Discussion id
// func (service *DiscussionService) GetDiscussionsByTalentID(allDiscussions *[]community.Discussion,
// 	talentID uuid.UUID, queryProcessor ...repository.QueryProcessor) error {

// 	// // Check if talent exist.
// 	// err := service.doesTalentExist(talentID)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	uow := repository.NewUnitOfWork(service.DB, true)

// 	queryProcessor = append(queryProcessor, repository.Filter("`talent_id` =?", talentID),
// 		repository.PreloadAssociations(service.associations))
// 	queryProcessor = append(queryProcessor, repository.OrderBy("`created_at`", true))

// 	err := service.Repository.GetAll(uow, allDiscussions, queryProcessor...)
// 	if err != nil {
// 		uow.RollBack()
// 		return errors.NewHTTPError("unable to get discussions", http.StatusInternalServerError)
// 	}
// 	for index, discussion := range *allDiscussions {
// 		err = service.setHasBestReply(&discussion)
// 		if err != nil {
// 			return err
// 		}
// 		(*allDiscussions)[index] = discussion
// 	}
// 	uow.Commit()
// 	return nil
// }

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// validateFieldUniqueness checks field uniqueness.
func (service *DiscussionService) validateFieldUniqueness(discussion *community.Discussion) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, discussion.TenantID, community.Discussion{},
		repository.Filter("`question`=? AND `id`!= ?", discussion.Question,
			discussion.ID))
	err = util.HandleIfExistsError("Similar question exists", exists, err)
	if err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

// func (service *DiscussionService) setHasBestReply(discussion *community.Discussion) error {
// 	var count int
// 	uow := repository.NewUnitOfWork(service.DB, true)
// 	err := service.Repository.GetCount(uow, community.Reply{}, &count,
// 		repository.Filter("discussion_id=?", discussion.ID), repository.Filter("best_reply=?", true))
// 	if err != nil {
// 		uow.RollBack()
// 		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
// 	}
// 	uow.Commit()
// 	if count != 0 {
// 		discussion.HasBestReply = true
// 		return nil
// 	}
// 	fmt.Println("NO BEST REPLY **************************")
// 	discussion.HasBestReply = false
// 	return nil
// }

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *DiscussionService) doForeignKeysExist(credentialID uuid.UUID, discussion *community.Discussion) error {
	tenantID := discussion.TenantID

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Check if channel exists.
	if err := service.doesChannelExist(tenantID, discussion.ChannelID); err != nil {
		return err
	}

	// Check if talent exists.
	if err := service.doesTalentExist(tenantID, discussion.AuthorID); err != nil {
		return err
	}

	// Check if best reply exists.
	if discussion.BestReplyID != nil {
		if err := service.doesReplyExist(tenantID, *discussion.BestReplyID); err != nil {
			return err
		}
	}
	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *DiscussionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *DiscussionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDiscussionExist will check if discussionID is valid
func (service *DiscussionService) doesDiscussionExist(tenantID, discussionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Discussion{},
		repository.Filter("`id`=?", discussionID))
	if err := util.HandleError("Invalid discussion id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesChannelExist will check if channelID is valid
func (service *DiscussionService) doesChannelExist(tenantID, channelID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Channel{},
		repository.Filter("`id`=?", channelID))
	if err := util.HandleError("Invalid channel id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTalentExist will check if talentID is valid
func (service *DiscussionService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Talent{},
		repository.Filter("`id`=?", talentID))
	if err := util.HandleError("Invalid talent id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesReplyExist will check if replyID is valid
func (service *DiscussionService) doesReplyExist(tenantID, replyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Reply{},
		repository.Filter("`id`=?", replyID))
	if err := util.HandleError("Invalid reply id", exists, err); err != nil {
		return err
	}
	return nil
}

// addSearchQueries adds search criteria to get specific discussions.
func (service *DiscussionService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	fmt.Println("=========================In discussion search============================", requestForm)
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	question := requestForm.Get("question")
	if !util.IsEmpty(question) {
		util.AddToSlice("`question`", "LIKE ?", "AND", "%"+question+"%", &columnNames, &conditions, &operators, &values)
	}
	if channelID, ok := requestForm["channelID"]; ok {
		util.AddToSlice("`channel_id`", "= ?", "AND", channelID, &columnNames, &conditions, &operators, &values)
	}
	if authorID, ok := requestForm["authorID"]; ok {
		util.AddToSlice("`author_id`", "= ?", "AND", authorID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
