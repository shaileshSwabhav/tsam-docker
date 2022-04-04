package controller

import (
	"net/http"

	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/security"

	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/web"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/community/service"
	"github.com/techlabs/swabhav/tsam/log"
)

// ChannelController Give access to Add, Update, Delete, Get and Get all Controllers Function
type ChannelController struct {
	service *service.ChannelService
	log     log.Logger
	auth    *security.Authentication
}

// NewChannelController returns new instance of Channel Controller
func NewChannelController(service *service.ChannelService, log log.Logger, auth *security.Authentication) *ChannelController {
	return &ChannelController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes will register the login routes.

// RegisterRoutes Register All Endpoint To Router.
func (controller *ChannelController) RegisterRoutes(router *mux.Router) {
	// get := router.HandleFunc("/tenant/{tenantID}/channel/{channelID}",
	// 	controller.GetChannel).Methods(http.MethodGet)
	unguarded := router.PathPrefix("/tenant/{tenantID}").Subrouter()
	guarded := router.PathPrefix("/tenant/{tenantID}").Subrouter()

	unguarded.HandleFunc("/channels",
		controller.GetChannels).Methods(http.MethodGet)

	guarded.HandleFunc("/channel",
		controller.AddChannel).Methods(http.MethodPost)

	guarded.HandleFunc("/channels/{channelID}",
		controller.UpdateChannel).Methods(http.MethodPut)

	guarded.HandleFunc("/channels/{channelID}",
		controller.DeleteChannel).Methods(http.MethodDelete)

	guarded.Use(controller.auth.Middleware)
	// *exclude = append(*exclude, getAll)
	controller.log.Info("Channel Routes Registered")
}

// AddChannel godoc
// AddChannel Adds New Channel
// @Description Adds New Channel
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param channel body community.Channel true "Add Channel"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /tenant/{tenantID}/channels [POST]
func (controller *ChannelController) AddChannel(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add Channel Called==============================")
	channel := &community.Channel{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, channel)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	channel.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	channel.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = channel.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddChannel(channel)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Channel successfully added")
}

// UpdateChannel godoc
// UpdateChannel Update Channel
// @Description Update Channel By Channel ID
// @Tags community-forum
// @Accept  json
// @Produce  plain
// @Param channelID path string true "Channel ID For Update" Format(uuid.UUID)
// @Param channel body community.Channel true "Update Channel"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /tenant/{tenantID}/channel/{channelID} [PUT]
func (controller *ChannelController) UpdateChannel(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Update Channel Called==============================")
	channel := community.Channel{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &channel)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data",
			http.StatusBadRequest))
		return
	}

	channel.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// UpdatedBy will be set with the given credentialID.
	channel.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	channel.ID, err = parser.GetUUID("channelID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("invalid Channel ID", http.StatusBadRequest))
		return
	}

	err = channel.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data",
			http.StatusBadRequest))
		return
	}

	err = controller.service.UpdateChannel(&channel)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Channel successfully updated")
}

// DeleteChannel godoc
// DeleteChannel Delete Channel By ChannelID
// @Description Delete Channel By ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param channelID path string true "channel ID" Format(uuid.UUID)
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /tenant/{tenantID}/channel/{channelID} [Delete]
func (controller *ChannelController) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Delete Channel Called==============================")

	channel := &community.Channel{}
	parser := web.NewParser(r)

	var err error

	channel.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	channel.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	channel.ID, err = parser.GetUUID("channelID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse channel id", http.StatusBadRequest))
		return
	}

	err = controller.service.DeleteChannel(channel)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Channel deleted successfully")
}

// GetChannels godoc
// GetChannels Return All Channels
// @Description Return All Channels
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param limit path int true "total number of result" Format(int)
// @Param offset path int true "page number" Format(int)
// @Success 200 {array} []community.Channel
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /tenant/{tenantID}/channel/{limit}/{offset} [GET]
func (controller *ChannelController) GetChannels(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Channels Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	var totalCount int
	channels := &[]community.Channel{}
	err = controller.service.GetChannels(tenantID, channels, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, channels)
}

// // GetChannel godoc
// // GetChannels Return Channel By ChannelID
// // @Description Return Channel By ID
// // @Tags community-forum
// // @Accept  json
// // @Produce  json
// // @Param channelID path string true "channel ID" Format(uuid.UUID)
// // @Success 200 {object} community.Channel
// // @Failure 400 {object} errors.ValidationError
// // @Failure 500 {object} errors.HTTPError
// // @Router /tenant/{tenantID}/channel/{channelID} [GET]
// func (controller *ChannelController) GetChannel(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================Get Channel Called==============================")
// 	channelID, err := util.ParseUUID(mux.Vars(r)["channelID"])

// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse Channel ID",
// 			http.StatusBadRequest))
// 		return
// 	}

// 	channel := community.Channel{}

// 	err = controller.ChannelService.GetChannel(&channel, channelID)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, channel)
// }
