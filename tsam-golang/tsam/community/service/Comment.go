package service

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

const commentAddedNotification = "New comment added"

// CommentService gives access to CRUD operations on comment.
type CommentService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewCommentService returns new instance of CommentService
func NewCommentService(db *gorm.DB, repo repository.Repository) *CommentService {
	return &CommentService{
		DB:         db,
		Repository: repo,
	}
}

// AddComment adds new comment.
func (commentService *CommentService) AddComment(comment *community.Reply) error {
	// comment.ID = util.GenerateUUID()

	// Check if reply exist.
	err := commentService.doesReplyExist(*comment.ReplyID)
	if err != nil {
		return err
	}

	// // Check if talent exist.
	// err = commentService.doesTalentExist(comment.Replier.ID)
	// if err != nil {
	// 	return err
	// }

	// // Check if discussion exist.
	// err = commentService.doesDiscussionExist(*comment.DiscussionID)
	// if err != nil {
	// 	return err
	// }

	uow := repository.NewUnitOfWork(commentService.DB, false)

	// reply := &community.Reply{}
	// // ReplyValidation checks for nil pointer errors initially.
	// err = commentService.Repository.Get(uow, *comment.ReplyID, reply)
	// if err != nil {
	// 	uow.RollBack()
	// 	log.NewLogger().Error(err.Error())
	// 	return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	// }

	// talent := &community.Talent{}
	// err = commentService.Repository.Get(uow, comment.Replier.ID, talent)
	// if err != nil {
	// 	uow.RollBack()
	// 	log.NewLogger().Error(err.Error())
	// 	return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	// }

	discussion := &community.Discussion{}
	// ReplyValidation checks for nil pointer errors initially.
	err = commentService.Repository.Get(uow, comment.DiscussionID, discussion,
		repository.Select("`id`, `talent_id`"))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// comment.ReplierID = &comment.Replier.ID
	// comment.Replier = nil
	err = commentService.Repository.Add(uow, comment)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	notificationType := &community.NotificationType{}
	err = commentService.Repository.GetRecord(uow, notificationType,
		repository.Filter("`notification_type_name`=?", commentAddedNotification))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("&&&&&&&&&&&&&&&&&&&&& ", notificationType)

	notificationService := NewNotificationService(commentService.DB, commentService.Repository)

	authorID := discussion.AuthorID
	authorNotification := &community.Notification{
		DiscussionID:       discussion.ID,
		ReplyID:            *comment.ReplyID,
		SubscriberID:       authorID,
		NotifierID:         comment.ReplierID,
		NotificationTypeID: notificationType.ID,
	}
	// ***********Change notification service to accept reply as argument & check if it is a reply
	// or comment inside it and call add accordingly.
	err = notificationService.AddNotificationWithoutTransaction(authorNotification, uow)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}
	var talentIDs []string
	err = repository.PluckColumn(notificationService.DB, "replies", "DISTINCT replier_id", &talentIDs,
		repository.Filter("`reply_id`=?", comment.ReplyID))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	for _, id := range talentIDs {
		talentID, err := util.ParseUUID(id)
		if err != nil {
			return err
		}
		if talentID == authorID {
			continue
		}
		// Only the talentID is different so simple authorNotification = notification should work!
		notification := &community.Notification{
			DiscussionID:       discussion.ID,
			ReplyID:            *comment.ReplyID,
			SubscriberID:       talentID,
			NotifierID:         comment.ReplierID,
			NotificationTypeID: notificationType.ID,
		}
		err = notificationService.AddNotificationWithoutTransaction(notification, uow)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
	}
	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doesTalentExist will check if talentID is valid
func (commentService *CommentService) doesTalentExist(talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(commentService.DB, community.Talent{},
		repository.Filter("`id`=?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesReplyExist will check if replyID is valid
func (commentService *CommentService) doesReplyExist(replyID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(commentService.DB, community.Reply{},
		repository.Filter("`id`=?", replyID))
	if err := util.HandleError("Invalid reply ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesDiscussionExist will check if discussionID is valid
func (commentService *CommentService) doesDiscussionExist(discussionID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(commentService.DB, community.Discussion{},
		repository.Filter("`id`=?", discussionID))
	if err := util.HandleError("Invalid discussion ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
