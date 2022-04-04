package controller

import (
	"net/http"

	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/community/service"
	"github.com/techlabs/swabhav/tsam/log"
)

// NotificationTypeController Give Access to Update, Delete, Add, Get and Get All Notification Type
type NotificationTypeController struct {
	NotificationTypeService *service.NotificationTypeService
}

// NewNotificationTypeController returns New Instance of NotificationTypeController
func NewNotificationTypeController(notificationService *service.NotificationTypeService) *NotificationTypeController {
	return &NotificationTypeController{
		NotificationTypeService: notificationService,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *NotificationTypeController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/tenant/{tenantID}/notification-type/credential/{credentialID}", controller.AddNotificationType).Methods(http.MethodPost)
	log.NewLogger().Info("Notification Type Routes Registered")
}

// AddNotificationType godoc
// AddNotificationType Add New Notification Type
// @Description Add New Notification Type
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param notification body community.NotificationType true "Add Notification Type"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /tenant/{tenantID}/notification-type/credential/{credentialID} [POST]
func (controller *NotificationTypeController) AddNotificationType(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddNotificationType called==============================")
	notificationType := &community.NotificationType{}
	params := mux.Vars(r)

	err := web.UnmarshalJSON(r, notificationType)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	notificationType.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	notificationType.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	err = notificationType.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	err = controller.NotificationTypeService.AddNotificationType(notificationType)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Notification type added.")
}
