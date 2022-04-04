package controller

import (
	"net/http"

	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/repository"

	"github.com/techlabs/swabhav/tsam/util"

	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/web"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/community/service"
	"github.com/techlabs/swabhav/tsam/log"
)

// NotificationController Give Access to Update, Delete, Add, Get and Get All Notification
type NotificationController struct {
	NotificationService *service.NotificationService
}

// NewNotificationController returns New Instance of NotificationController
func NewNotificationController(notificationService *service.NotificationService) *NotificationController {
	return &NotificationController{
		NotificationService: notificationService,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (notificationController *NotificationController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/tenant/{tenantID}/credential/{credentialID}/notification/limit/{limit}/offset/{offset}",
		notificationController.GetAllNotificationsByTalentID).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/discussion/{discussionID}/notification/credential/{credentialID}",
		notificationController.UpdateNotificationOfTalentToSeen).Methods(http.MethodPut)
	//  Below APIs not required yet.
	router.HandleFunc("/tenant/{tenantID}/notification/credential/{credentialID}",
		notificationController.AddNotification).Methods(http.MethodPost)

	// router.HandleFunc("/tenant/{tenantID}/notification/limit/{limit}/offset/{offset}",
	// 	notificationController.GetAllNotifications).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/notification/{notificationID}",
		notificationController.GetNotificationByID).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/notification/{notificationID}/credential/{credentialID}",
		notificationController.UpdateNotification).Methods(http.MethodPut)
	router.HandleFunc("/tenant/{tenantID}/notification/{notificationID}/credential/{credentialID}",
		notificationController.DeleteNotification).Methods(http.MethodDelete)
	log.NewLogger().Info("Notification Routes Registered")
}

// AddNotification godoc
// AddNotification Add New Notification
// @Description Add New Notification
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param notification body community.Notification true "Add Notification"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /notification [POST]
func (notificationController *NotificationController) AddNotification(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddNotification called==============================")

	notification := community.Notification{}

	err := web.UnmarshalJSON(r, &notification)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	err = notification.ValidateNotification()
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	err = notificationController.NotificationService.AddNotification(&notification)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, notification.ID)
}

// UpdateNotification godoc
// UpdateNotification Update Notification By Notification ID
// @Description Update Notification By Notification ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param notification body community.Notification true "Add Notification"
// @Param notificationID path string true "Notification ID For Update Notification" Format(uuid.UUID)
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /notification/{notificationID} [PUT]
func (notificationController *NotificationController) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateNotification called==============================")
	notification := community.Notification{}

	err := web.UnmarshalJSON(r, &notification)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	notification.ID, err = util.ParseUUID(mux.Vars(r)[paramNotificationID])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to get notification ID", http.StatusBadRequest))
		return
	}

	err = notification.ValidateNotification()
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	err = notificationController.NotificationService.UpdateNotification(&notification)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Notification Updated")
}

// DeleteNotification godoc
// DeleteNotification Delete Notification By Notification ID
// @Description Delete Notification By Notification ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param notificationID path string true "Notification ID For Delete Notification" Format(uuid.UUID)
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /notification/{notificationID} [DELETE]
func (notificationController *NotificationController) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteNotification called==============================")
	notification := community.Notification{}
	var err error

	notification.ID, err = util.ParseUUID(mux.Vars(r)[paramNotificationID])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to get Notification ID", http.StatusBadRequest))
		return
	}
	err = notificationController.NotificationService.DeleteNotification(&notification)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Notification deleted")
}

// GetNotificationByID godoc
// GetNotificationByID Get Notification By Notification ID
// @Description GetNotificationByID gets notification by notification ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param notificationID path string true "Notification ID For Get Notification" Format(uuid.UUID)
// @Success 200 {object} community.Notification
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /notification/{notificationID} [GET]
func (notificationController *NotificationController) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetNotificationByID called==============================")
	notification := &community.Notification{}
	var err error

	notification.ID, err = util.ParseUUID(mux.Vars(r)[paramNotificationID])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to get Notification ID", http.StatusBadRequest))
		return
	}
	// notification.ID = notificationID

	err = notificationController.NotificationService.GetNotificationByID(notification)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// notificationSlice := &[]*community.Notification{
	// 	notification,
	// }
	web.RespondJSON(w, http.StatusOK, notification)

}

// GetAllNotifications godoc
// GetAllNotifications returns all notifications
// @Description Return multiple Notification By Discussion ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param limit path int true "total number of result" Format(int)
// @Param offset path int true "page number" Format(int)
// @Success 200 {array} []community.Notification
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /notification/{limit}/{offset} [GET]
func (notificationController *NotificationController) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllNotifications called==============================")
	allNotifications := &[]community.Notification{}

	err := notificationController.NotificationService.GetAllNotifications(allNotifications, repository.Paging(w, r))
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, allNotifications)

}

// GetAllNotificationsByTalentID godoc
// GetAllNotificationsByTalentID returns notifications
// @Description returns all notifications of a specific talent
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param limit path int true "total number of result" Format(int)
// @Param offset path int true "page number" Format(int)
// @Success 200 {array} []community.Notification
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /notification/{limit}/{offset} [GET]
func (notificationController *NotificationController) GetAllNotificationsByTalentID(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllNotificationsByTalentID called==============================")
	allNotificationsDTO := &[]*community.NotificationDTO{}

	talentID, err := util.ParseUUID(mux.Vars(r)[paramTalentID])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to get Talent ID", http.StatusBadRequest))
		return
	}

	err = notificationController.NotificationService.GetAllNotificationsByID(w, r, talentID, allNotificationsDTO,
		repository.Paging(w, r))
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, allNotificationsDTO)

}

// UpdateNotificationOfTalentToSeen godoc
// UpdateNotificationOfTalentToSeen returns "updated to seen"
// @Description sets specific notifications of talent's isSeen & seenTime
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /notification/{limit}/{offset} [PUT]
func (notificationController *NotificationController) UpdateNotificationOfTalentToSeen(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateNotificationOfTalentToSeen called==============================")
	talentID, err := util.ParseUUID(mux.Vars(r)[paramTalentID])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to get Talent ID", http.StatusBadRequest))
		return
	}
	discussionID, err := util.ParseUUID(mux.Vars(r)[paramDiscussionID])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to get discussion ID", http.StatusBadRequest))
		return
	}
	err = notificationController.NotificationService.UpdateNotificationOfTalentToSeen(w, talentID, discussionID,
		repository.Paging(w, r))
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "updated to seen")

}
