package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/config"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/homePage/contact/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/homePage/contact"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

type ContactInfoController struct {
	service *service.ContactInfoService
	log     log.Logger
	auth    *security.Authentication
}

func NewContactInfoController(service *service.ContactInfoService, log log.Logger, auth *security.Authentication) *ContactInfoController {
	return &ContactInfoController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

func (controller *ContactInfoController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	sendEmail := router.HandleFunc("/contact/send-mail", controller.ContactUs).Methods(http.MethodPost)

	*exclude = append(*exclude, sendEmail)
	controller.log.Info("Send Contact Mail Routes Registered")
}

func (controller *ContactInfoController) ContactUs(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================SendContactMail API Called==============================")
	emailBody := contact.ContactInfo{}

	email := controller.auth.Config.GetString(config.SenderMail)
	pass := controller.auth.Config.GetString(config.SenderPass)
	err := web.UnmarshalJSON(r, &emailBody)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse requested data", http.StatusBadRequest))
		return
	}

	err = controller.service.ContactUs(&emailBody, email, pass)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("response couldn't be send, try again", http.StatusBadRequest))
		return
	}

	web.RespondJSON(w, http.StatusOK, "Response Sent Successfully")
}
