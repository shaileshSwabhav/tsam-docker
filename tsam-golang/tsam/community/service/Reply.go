package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

const (
	replyAddedNotification = "New reply added"
	markedBestReply        = "Marked as best reply"
)

// ReplyService provides CRUD operations for replies.
type ReplyService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewReplyService returns new instance of ReplyService.
func NewReplyService(db *gorm.DB, repo repository.Repository) *ReplyService {
	return &ReplyService{
		DB:           db,
		Repository:   repo,
		associations: []string{"Replier"},
		// Preload comments.
	}
}

// AddReply adds a new reply.
func (service *ReplyService) AddReply(reply *community.Reply) error {

	// check if all foreign keys exist.
	err := service.doForeignKeysExist(reply.CreatedBy, reply)
	if err != nil {
		return err
	}

	err = service.validateFields(reply)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, reply)
	if err != nil {
		uow.RollBack()
		return err
	}

	// notificationType := &community.NotificationType{}
	// err = service.Repository.GetRecord(uow, notificationType, repository.Filter("notification_type_name=?", replyAddedNotification))
	// if err != nil {
	// 	uow.RollBack()
	// 	return errors.NewHTTPError("notification type error: "+err.Error(), http.StatusInternalServerError)
	// }

	// notificationService := NewNotificationService(service.DB, service.Repository)

	// notification := &community.Notification{
	// 	DiscussionID:       *reply.DiscussionID,
	// 	ReplyID:            reply.ID,
	// 	NotifierTalentID:   *reply.ReplierID,
	// 	NotifiedTalentID:   discussion.AuthorID,
	// 	NotificationTypeID: notificationType.ID,
	// }
	// err = notificationService.AddNotificationWithoutTransaction(notification, uow)
	// if err != nil {
	// 	uow.RollBack()
	// 	return errors.NewHTTPError("Can't add notification: "+err.Error(), http.StatusInternalServerError)
	// }
	uow.Commit()
	return nil
}

// UpdateReply updates reply in database.
func (service *ReplyService) UpdateReply(reply *community.Reply) error {

	// Check if discussion exist.
	err := service.doesReplyExist(reply.TenantID, reply.ID)
	if err != nil {
		return err
	}

	// check if all foreign keys exist.
	err = service.doForeignKeysExist(reply.UpdatedBy, reply)
	if err != nil {
		return err
	}

	err = service.validateFields(reply)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// tempReply := community.Reply{}
	// tempReply.ID = reply.ID
	// err = service.GetReply(&tempReply)
	// if err != nil {
	// 	return err
	// }

	// // tempTalent := community.Talent{}
	// // err = service.Repository.Get(uow, reply.Replier.ID, &tempTalent)
	// // if err != nil {
	// // 	uow.RollBack()
	// // 	if util.IsContain(err.Error(), "not found") {
	// // 		return errors.NewHTTPError("invalid user", http.StatusBadRequest)
	// // 	}
	// // 	return errors.NewHTTPError("unable to add Reply", http.StatusInternalServerError)
	// // }

	// discussionService := NewDiscussionService(service.DB, service.Repository)
	// discussion := &community.DiscussionDTO{}

	// // if reply.DiscussionID == nil {
	// // 	return errors.NewHTTPError("Discussion id must be specified", http.StatusBadRequest)
	// // }

	// discussion.ID = *reply.DiscussionID
	// err = discussionService.GetDiscussion(discussion)
	// if err != nil {
	// 	uow.RollBack()
	// 	return errors.NewHTTPError("Can't get discussion", http.StatusBadRequest)
	// }

	// reply.ReplierID = &reply.Replier.ID
	// reply.Replier = nil

	// // Add a notification if it marked as best reply.
	// isBestReply := false
	// if tempReply.BestReply != nil {
	// 	isBestReply = *(tempReply.BestReply)
	// }
	// if !isBestReply {
	// 	if reply.BestReply != nil {
	// 		if *(reply.BestReply) {
	// 			notificationType := &community.NotificationType{}
	// 			err = service.Repository.GetRecord(uow, notificationType,
	// 				repository.Filter("notification_type_name=?", markedBestReply))
	// 			if err != nil {
	// 				uow.RollBack()
	// 				return errors.NewHTTPError("notification type error: "+err.Error(), http.StatusInternalServerError)
	// 			}
	// 			notificationService := NewNotificationService(service.DB, service.Repository)

	// 			notification := &community.Notification{
	// 				DiscussionID:       *reply.DiscussionID,
	// 				ReplyID:            reply.ID,
	// 				NotifierTalentID:   discussion.AuthorID,
	// 				NotifiedTalentID:   *reply.ReplierID,
	// 				NotificationTypeID: notificationType.ID,
	// 			}
	// 			err = notificationService.AddNotificationWithoutTransaction(notification, uow)
	// 			if err != nil {
	// 				uow.RollBack()
	// 				return errors.NewHTTPError("Can't add notification: "+err.Error(), http.StatusInternalServerError)
	// 			}
	// 		}
	// 	}
	// }

	err = service.Repository.Update(uow, reply)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// DeleteReply Return Reply By ID
func (service *ReplyService) DeleteReply(reply *community.Reply) error {

	// Check if reply exist
	err := service.doesReplyExist(reply.TenantID, reply.ID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(reply.TenantID, reply.DeletedBy); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.UpdateWithMap(uow, reply, map[string]interface{}{
		"DeletedBy": reply.DeletedBy,
		"DeletedAt": time.Now(),
	})
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("unable to delete Reply", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// GetReply Return Reply By ID
func (service *ReplyService) GetReply(reply *community.Reply, queryProcessor ...repository.QueryProcessor) error {

	// Check if reply exist
	err := service.doesReplyExist(reply.TenantID, reply.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	queryProcessor = append(queryProcessor, repository.PreloadAssociations(service.associations))

	err = service.Repository.Get(uow, reply.ID, reply, queryProcessor...)
	if err != nil {
		uow.RollBack()
		if uow.DB.RecordNotFound() {
			return errors.NewHTTPError("Record not found.", http.StatusBadRequest)
		}
		return errors.NewHTTPError("unable to get reply", http.StatusInternalServerError)
	}
	// var noOfLikes *int
	// noOfLikes, err = service.getNumberOfLikesOfReply(reply)
	// if err != nil {
	// 	return err
	// }
	// reply.NumberOfLikes = *noOfLikes
	uow.Commit()
	return nil
}

// GetRepliesByDiscussionID Get Reply By Discussion ID
func (service *ReplyService) GetRepliesByDiscussionID(tenantID uuid.UUID,
	discussionID uuid.UUID, replies *[]community.ReplyDTO, parser *web.Parser, totalCount *int) error {
	// func (service *ReplyService) GetRepliesByDiscussionID(pagination *general.Pagination, formValues url.Values,
	// 	discussionID uuid.UUID, replies *[]community.ReplyDTO) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, replies,
		"`created_at`",
		service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount),
		repository.PreloadAssociations(service.associations),
	)
	if err != nil {
		uow.RollBack()
		return err
	}

	// queryProcessor = append(queryProcessor, repository.Filter("`discussion_id` = ? AND `reply_id` IS NULL", discussionID),
	// 	repository.PreloadAssociations(service.association))
	// // Swipe QueryProcessor 1st & last
	// queryProcessor[0], queryProcessor[len(queryProcessor)-1] = queryProcessor[len(queryProcessor)-1], queryProcessor[0]
	// queryProcessor = append(queryProcessor, repository.OrderBy("`created_at`", true), repository.OrderBy("`best_reply` DESC", true))

	// err = service.Repository.GetAll(uow, replies, queryProcessor...)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	uow.RollBack()
	// 	return errors.NewHTTPError("unable to get Reply", http.StatusInternalServerError)
	// }
	// if replies != nil {
	// 	for index, reply := range *replies {
	// 		noOfLikes, err := service.getNumberOfLikesOfReply(&reply)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		(*replies)[index].NumberOfLikes = *noOfLikes

	// 		// set isLiked flag
	// 		if talentID != uuid.Nil {
	// 			totalLikes := 0
	// 			uow := repository.NewUnitOfWork(service.DB, true)
	// 			err := service.Repository.GetCount(uow, &community.Like{}, &totalLikes,
	// 				repository.Filter("`talent_id` = ? AND `reply_id`=?", talentID, reply.ID))
	// 			if err != nil {
	// 				return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	// 			}
	// 			if totalLikes > 0 {
	// 				(*replies)[index].IsLiked = true
	// 				continue
	// 			}
	// 		}
	// 		(*replies)[index].IsLiked = false
	// 	}
	// }
	replies, err = service.getRepliesOfReply(uow, replies)
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

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *ReplyService) doForeignKeysExist(credentialID uuid.UUID, reply *community.Reply) error {
	tenantID := reply.TenantID

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Check if discussion exists.
	if err := service.doesDiscussionExist(tenantID, reply.DiscussionID); err != nil {
		return err
	}

	// Check if talent exists.
	if err := service.doesReplierExist(tenantID, reply.ReplierID); err != nil {
		return err
	}

	// Check if best reply exists.
	if reply.ReplyID != nil {
		if err := service.doesReplyExist(tenantID, *reply.ReplyID); err != nil {
			return err
		}
	}
	return nil
}

// addSearchQueries adds search criteria to get specific discussions.
func (service *ReplyService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	// question := requestForm.Get("question")
	// if !util.IsEmpty(question) {
	// 	util.AddToSlice("`question`", "LIKE ?", "AND", "%"+question+"%", &columnNames, &conditions, &operators, &values)
	// }
	// if channelID, ok := requestForm["channelID"]; ok {
	// 	util.AddToSlice("`channel_id`", "= ?", "AND", channelID, &columnNames, &conditions, &operators, &values)
	// }
	// if authorID, ok := requestForm["authorID"]; ok {
	// 	util.AddToSlice("`author_id`", "= ?", "AND", authorID, &columnNames, &conditions, &operators, &values)
	// }

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// Return Replies Of Reply
func (service *ReplyService) getRepliesOfReply(uow *repository.UnitOfWork, replyList *[]community.ReplyDTO) (*[]community.ReplyDTO, error) {
	for index, reply := range *replyList {
		replies := []community.ReplyDTO{}
		err := service.Repository.GetAllInOrder(uow, &replies, "created_at ASC",
			repository.PreloadAssociations(service.associations), repository.Filter("`reply_id` =?", reply.ID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return nil, errors.NewHTTPError("unable to get Reply", http.StatusInternalServerError)
		}
		if len(replies) > 0 {
			reply.Comments = &replies
		} else {
			reply.Comments = nil
		}
		rep, err := service.getRepliesOfReply(uow, &replies)
		if err != nil {
			return nil, err
		}
		if len(*rep) > 0 {
			reply.Comments = rep
		}
		(*replyList)[index] = reply
	}
	return replyList, nil
}

// pass uow #n
func (service *ReplyService) getNumberOfLikesOfReply(reply *community.ReplyDTO) (*int, error) {
	var totalLikes int
	uow := repository.NewUnitOfWork(service.DB, true)
	err := service.Repository.GetCount(uow, &community.Reaction{}, &totalLikes, repository.Filter("reply_id = ?", reply.ID))
	if err != nil {
		uow.RollBack()
		return nil, errors.NewHTTPError("unable to get likes", http.StatusInternalServerError)
	}
	uow.Commit()
	return &totalLikes, nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *ReplyService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *ReplyService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDiscussionExist will check if discussionID is valid
func (service *ReplyService) doesDiscussionExist(tenantID, discussionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Discussion{},
		repository.Filter("`id`=?", discussionID))
	if err := util.HandleError("Invalid discussion id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesReplyExist will check if replyID is valid
func (service *ReplyService) doesReplyExist(tenantID, replyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Reply{},
		repository.Filter("`id`=?", replyID))
	if err := util.HandleError("Invalid reply id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesReplierExist will check if credentialID is valid
func (service *ReplyService) doesReplierExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid replierID", exists, err); err != nil {
		return err
	}
	return nil
}

// validate fields will validate and check for uniqueness in particular fields.
func (service *ReplyService) validateFields(reply *community.Reply) error {
	if reply.IsBestReply != nil && *reply.IsBestReply {
		if reply.ReplyID != nil {
			return errors.NewHTTPError("Comment can't be marked as best answer",
				http.StatusBadRequest)
		}
		exists, err := repository.DoesRecordExistForTenant(service.DB, reply.TenantID, reply,
			repository.Filter("`discussion_id` = ? AND `is_best_reply` = ?", reply.DiscussionID, true))
		err = util.HandleIfExistsError("Discussion already has a best reply.", exists, err)
		if err != nil {
			return err
		}
	}
	return nil
}
