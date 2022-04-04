package service

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	html "github.com/grokify/html-strip-tags-go"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// NotificationService gives access to CRUD operations.
type NotificationService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewNotificationService returns an instance of NotificationService.
func NewNotificationService(db *gorm.DB, repo repository.Repository) *NotificationService {
	return &NotificationService{
		DB:         db,
		Repository: repo,
	}
}

// AddNotification adds new notification.
func (notificationService *NotificationService) AddNotification(notification *community.Notification) error {

	// Check if all foreign key exist
	err := notificationService.doesForeignKeyExist(notification)
	if err != nil {
		return err
	}

	// notification.ID = util.GenerateUUID()
	unseen := false
	notification.IsSeen = &unseen

	uow := repository.NewUnitOfWork(notificationService.DB, false)

	err = notificationService.Repository.Add(uow, notification)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Unable to add notification", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// AddNotificationWithoutTransaction adds new notification without transaction.
func (notificationService *NotificationService) AddNotificationWithoutTransaction(notification *community.Notification,
	uow *repository.UnitOfWork) error {

	// Check if all foreign key exist
	err := notificationService.doesForeignKeyExist(notification)
	if err != nil {
		return err
	}

	if notification.NotifierID == notification.SubscriberID {
		return nil
	}

	// notification.ID = util.GenerateUUID()
	unseen := false
	notification.IsSeen = &unseen

	err = notification.ValidateNotification()
	if err != nil {
		return err
	}

	err = notificationService.Repository.Add(uow, notification)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Unable to add notification", http.StatusInternalServerError)
	}
	return nil
}

// UpdateNotification updates the notification by ID.
func (notificationService *NotificationService) UpdateNotification(notification *community.Notification) error {

	// Check if all foreign key exist
	err := notificationService.doesForeignKeyExist(notification)
	if err != nil {
		return err
	}

	// Check if notification exist
	err = notificationService.doesNotificationExist(notification.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(notificationService.DB, false)

	err = notificationService.Repository.Update(uow, notification)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Unable to update notification", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// DeleteNotification sets deletedAt of the notification from db.
func (notificationService *NotificationService) DeleteNotification(notification *community.Notification) error {

	// Check if notification exist
	err := notificationService.doesNotificationExist(notification.ID)
	if err != nil {
		return err
	}

	// uow := repository.NewUnitOfWork(notificationService.DB, true)
	// notification := &community.Notification{}
	// err := notificationService.GetNotificationByID(notification)
	// if err != nil {
	// 	uow.RollBack()
	// 	return errors.NewHTTPError("no record found", http.StatusBadRequest)
	// }
	// uow.Commit()

	uow := repository.NewUnitOfWork(notificationService.DB, false)
	err = notificationService.Repository.Delete(uow, notification)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Unable to get notification", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// GetNotificationByID returns the notification.
func (notificationService *NotificationService) GetNotificationByID(notification *community.Notification,
	queryProcessors ...repository.QueryProcessor) error {

	// Check if notification exist
	err := notificationService.doesNotificationExist(notification.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(notificationService.DB, true)

	err = notificationService.Repository.Get(uow, notification.ID, notification, queryProcessors...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		// if util.IsContain(err.Error(), "not found") {
		// 	return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		// }
		return errors.NewHTTPError("Unable to get notification", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// GetAllNotifications gets all notifications from DB with a query processor.
func (notificationService *NotificationService) GetAllNotifications(allNotifications *[]community.Notification,
	queryProcessors ...repository.QueryProcessor) error {

	uow := repository.NewUnitOfWork(notificationService.DB, true)

	err := notificationService.Repository.GetAllInOrder(uow, allNotifications, "`created_at` DESC", queryProcessors...)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil

}

// UpdateNotificationOfTalentToSeen updates specific notifications of talent to seen
func (notificationService *NotificationService) UpdateNotificationOfTalentToSeen(w http.ResponseWriter,
	notifiedID, discussionID uuid.UUID, queryProcessors ...repository.QueryProcessor) error {

	// Check if notified talent exist
	err := notificationService.doesTalentExist(notifiedID)
	if err != nil {
		return err
	}

	// Check if discussion exist
	err = notificationService.doesDiscussionExist(discussionID)
	if err != nil {
		return err
	}

	discussionNotifications := []community.Notification{}

	uow := repository.NewUnitOfWork(notificationService.DB, false)

	err = notificationService.GetAllNotifications(&discussionNotifications,
		repository.Filter("`discussion_id`=?", discussionID), repository.Filter("`notified_talent_id`=?", notifiedID))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	valueTime := time.Now()
	valueTrue := true

	for _, notification := range discussionNotifications {
		notification.IsSeen = &valueTrue
		notification.SeenTime = &valueTime
		// **** modify query****
		if err := notificationService.Repository.Update(uow, notification); err != nil {
			uow.RollBack()
			return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		}
	}

	// Need to commit before adding new headers.
	uow.Commit()
	err = notificationService.setUnseenCountOfNotificationsToHeader(w, notifiedID)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil

}

// GetAllNotificationsByID gets all notifications of a specific talent.
func (notificationService *NotificationService) GetAllNotificationsByID(w http.ResponseWriter,
	req *http.Request, notifiedID uuid.UUID, allNotificationsDTO *[]*community.NotificationDTO,
	queryProcessors ...repository.QueryProcessor) error {

	// Check if notification exist
	err := notificationService.doesTalentExist(notifiedID)
	if err != nil {
		return err
	}

	notifications := &[]community.Notification{}
	notificationsDTO := []*community.NotificationDTO{}

	uow := repository.NewUnitOfWork(notificationService.DB, true)
	lim := mux.Vars(req)["limit"]
	err = notificationService.setUnseenCountOfNotificationsToHeader(w, notifiedID)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	limit, _ := strconv.Atoi(lim)
	if limit == 0 {
		return nil
	}

	queryProcessors = append(queryProcessors, repository.Filter("`notified_talent_id`=?", notifiedID))
	err = notificationService.Repository.GetAllInOrder(uow, notifications, "`created_at` DESC", queryProcessors...)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	//  ******** Need to change pluck column to getRecord. So we can use a simple string variable ******
	for _, notification := range *notifications {
		notificationDTO := community.NewNotificationDTO()

		notificationDTO.ID = notification.ID
		notificationDTO.IsSeen = notification.IsSeen
		notificationDTO.Added = getDuration(notification.CreatedAt)

		var (
			question, notificationTypeName []string
			// firstName,
		)
		notifierTalent := &community.Talent{}
		notifiedTalent := &community.Talent{}
		err := repository.PluckColumn(notificationService.DB, "discussions", "question",
			&question, repository.Filter("`id` = ?", notification.DiscussionID))
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		}

		err = notificationService.Repository.Get(uow, notification.SubscriberID, notifiedTalent)
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		}

		err = notificationService.Repository.Get(uow, notification.NotifierID, notifierTalent)
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		}

		notificationDTO.NotifiedTalent = notifiedTalent
		notificationDTO.NotifierTalent = notifierTalent
		notificationDTO.Discussion.DiscussionID = notification.DiscussionID
		notificationDTO.Discussion.Question = shortenString(question[0])

		notificationReply := community.Reply{}
		err = notificationService.Repository.Get(uow, notification.ReplyID, &notificationReply)
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		}
		// err = repository.PluckColumn(notificationService.DB, "replies", "replier_id",
		// 	&replierID, repository.Filter("`id` = ?", notification.ReplyID))
		// if err != nil {
		// 	return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		// }
		// err = repository.PluckColumn(notificationService.DB, "talents", "first_name",
		// 	&firstName, repository.Filter("`id` = ?", notificationReply.ID))
		// if err != nil {
		// 	return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		// }
		if notification.ReplyID != uuid.Nil {
			notificationDTO.Reply.ReplyID = notification.ReplyID
			notificationDTO.Reply.Replier.ID = notificationReply.ReplierID
			// notificationDTO.Reply.Replier.FirstName = &firstName[0]
		}
		// Simple get can be used.
		err = repository.PluckColumn(notificationService.DB, "notification_types", "notification_type_name",
			&notificationTypeName, repository.Filter("`id` = ?", notification.NotificationTypeID))
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		}

		// Create a NotificationDTO with all required fields
		notificationDTO.NotificationType.ID = notification.NotificationTypeID
		notificationDTO.NotificationType.NotificationTypeName = notificationTypeName[0]

		if util.IsEmpty(notifierTalent.FirstName) && len(notificationTypeName) > 0 && notificationDTO.Discussion.Question != nil {
			notifierFullName := notifierTalent.FirstName
			if util.IsEmpty(notifierTalent.LastName) {
				notifierFullName = notifierTalent.FirstName + " " + notifierTalent.LastName
			}
			var notificationText string
			if notificationDTO.NotificationType.NotificationTypeName == "New reply added" {
				notificationText = fmt.Sprintf("<b>%s</b> by <b>%s</b> on discussion with question <b>%s</b>",
					notificationTypeName[0], notifierFullName,
					*notificationDTO.Discussion.Question)
			} else {
				htmlStrippedReply := html.StripTags(notificationReply.Reply)
				reply := *shortenString(htmlStrippedReply)
				notificationText =
					fmt.Sprintf("<b>%s</b> by <b>%s</b> on reply <b>%s</b> on discussion with question <b>%s</b>",
						notificationTypeName[0], notifierFullName, reply, *notificationDTO.Discussion.Question)
			}
			notificationDTO.NotificationText = &notificationText
		}
		notificationsDTO = append(notificationsDTO, notificationDTO)
	}

	*allNotificationsDTO = notificationsDTO
	uow.Commit()
	return nil

}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

func (notificationService *NotificationService) setUnseenCountOfNotificationsToHeader(w http.ResponseWriter, notifiedID uuid.UUID) error {
	var totalNotifications, totalUnseenNotifications int

	uow := repository.NewUnitOfWork(notificationService.DB, true)
	err := notificationService.Repository.GetCount(uow, &community.Notification{}, &totalUnseenNotifications, repository.Filter("is_seen =?", false),
		repository.Filter("notified_talent_id=?", notifiedID))
	if err != nil {
		return err
	}
	// Need to clean this part
	err = notificationService.Repository.GetCount(uow, &community.Notification{}, &totalNotifications,
		repository.Filter("notified_talent_id=?", notifiedID))
	if err != nil {
		return err
	}

	web.SetNewHeader(w, "Talent-Unseen-Notifications-Total-Count", strconv.Itoa(totalUnseenNotifications))
	web.SetNewHeader(w, "Talent-Notifications-Total-Count", strconv.Itoa(totalNotifications))

	// w.Header().Add("Access-Control-Expose-Headers", "Talent-Unseen-Notifications-Total-Count, Talent-Notifications-Total-Count")
	// w.Header().Set("Talent-Unseen-Notifications-Total-Count", strconv.Itoa(totalUnseenNotifications))
	// w.Header().Set("Talent-Notifications-Total-Count", strconv.Itoa(totalNotifications))
	return nil

}

// change to directly getting hours,minutes etc.
func getDuration(createdTime time.Time) *string {
	var result string
	min := 60
	hour := min * min
	day := hour * 24

	duration := int(time.Now().Sub(createdTime).Seconds())
	// duration := int(time.Since(createdTime).Seconds())
	if duration < min {
		result = fmt.Sprintf("%d seconds", duration)
		return &result
	}
	if duration < min+min {
		result = fmt.Sprintf("%d minute", duration/min)
		return &result
	}
	if duration < min*min {
		result = fmt.Sprintf("%d minutes", duration/min)
		return &result
	}
	if duration < hour+hour {
		result = fmt.Sprintf("%d hour", duration/hour)
		return &result
	}
	if duration < day {
		result = fmt.Sprintf("%d hours", duration/hour)
		return &result
	}
	if duration < day+day {
		result = fmt.Sprintf("%d day", duration/day)
		return &result
	}
	result = fmt.Sprintf("%d days", duration/day)
	return &result
}

func shortenString(str string) *string {
	if len(str) > 30 {
		str = str[:30] + "..."
	}
	return &str
}

func (notificationService *NotificationService) doesForeignKeyExist(notification *community.Notification) error {

	// Check if notifier talent exist
	err := notificationService.doesTalentExist(notification.NotifierID)
	if err != nil {
		return err
	}

	// Check if notified talent exist
	err = notificationService.doesTalentExist(notification.SubscriberID)
	if err != nil {
		return err
	}

	// Check if discussion exist
	err = notificationService.doesDiscussionExist(notification.DiscussionID)
	if err != nil {
		return err
	}

	// Check if reply exist
	err = notificationService.doesReplyExist(notification.ReplyID)
	if err != nil {
		return err
	}

	// Check if notification type exist
	err = notificationService.doesNotificationTypeExist(notification.NotificationTypeID)
	if err != nil {
		return err
	}

	return nil
}

// doesTalentExist will check if talentID is valid
func (notificationService *NotificationService) doesTalentExist(talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(notificationService.DB, community.Talent{},
		repository.Filter("`id`=?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesNotificationExist will check if notificationID is valid
func (notificationService *NotificationService) doesNotificationExist(notificationID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(notificationService.DB, community.Notification{},
		repository.Filter("`id`=?", notificationID))
	if err := util.HandleError("Invalid notification ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesNotificationTypeExist will check if notificationTypeID is valid
func (notificationService *NotificationService) doesNotificationTypeExist(notificationTypeID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(notificationService.DB, community.NotificationType{},
		repository.Filter("`id`=?", notificationTypeID))
	if err := util.HandleError("Invalid notification type ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesDiscussionExist will check if discussionID is valid
func (notificationService *NotificationService) doesDiscussionExist(discussionID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(notificationService.DB, community.Discussion{},
		repository.Filter("`id`=?", discussionID))
	if err := util.HandleError("Invalid discussion ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesReplyExist will check if replyID is valid
func (notificationService *NotificationService) doesReplyExist(replyID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(notificationService.DB, community.Reply{},
		repository.Filter("`id`=?", replyID))
	if err := util.HandleError("Invalid reply ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
