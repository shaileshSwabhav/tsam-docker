package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	notificationdemo "github.com/techlabs/swabhav/tsam/notification_demo"
	"github.com/techlabs/swabhav/tsam/notification_demo/service"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// import (
// 	"log"

// 	"github.com/gorilla/mux"
// 	"github.com/techlabs/swabhav/tsam/security"
// )

type BlogNotificationController struct {
	blogNotificationService *service.NotificationService
	log                     log.Logger
	auth                    *security.Authentication
}

func NewAndroidUserController(blogNotificationService *service.NotificationService, log log.Logger, auth *security.Authentication) *BlogNotificationController {
	return &BlogNotificationController{
		blogNotificationService: blogNotificationService,
		log:                     log,
		auth:                    auth,
	}
}

func (controller *BlogNotificationController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	//get
	router.HandleFunc("/tenant/{tenantID}/notifications/credential/{credentialID}", controller.GetAllNotifications).Methods(http.MethodGet)

}

func (controller *BlogNotificationController) GetAllNotifications(w http.ResponseWriter, r *http.Request) {

	controller.log.Info("********************************GetAllNotifications call**************************************")
	parser := web.NewParser(r)

	notifications := []notificationdemo.Notification_Test_DTO{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	notifiedUserID, err := parser.GetUUID("credentialID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}
	err = controller.blogNotificationService.GetAllNotifications(tenantID, notifiedUserID, &notifications)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, notifications)
}

// 	// add
// 	// router.HandleFunc("/tenant/{tenantID}/andr", controller.AddAndroidUser).Methods(http.MethodPost)

// 	// login
// 	// login := router.HandleFunc("/tenant/{tenantID}/android-login",
// 	// 	controller.Login).Methods(http.MethodPost)
// 	// router.HandleFunc("/tenant/{tenantID}/company/{companyID}", controller.UpdateAndroidUser).Methods(http.MethodPut)

// 	// delete
// 	// router.HandleFunc("/tenant/{tenantID}/company/{companyID}", controller.DeleteAndroidUser).Methods(http.MethodDelete)

// 	// get
// 	// router.HandleFunc("/tenant/{tenantID}/company/{companyID}", controller.GetCompanyByI).Methods(http.MethodGet)
// 	// router.HandleFunc("/tenant/{tenantID}/company", controller.GetAllCompanies).Methods(http.MethodGet)

// 	// *exclude = append(*exclude, addAndroidUser, login)
// 	// controller.log.Info("Blog Notification Routes Registered")
// }
