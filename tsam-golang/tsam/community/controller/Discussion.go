package controller

import (
	"net/http"

	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/security"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/community/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/web"
)

// DiscussionController gives access to CRUD operations for entity
type DiscussionController struct {
	service *service.DiscussionService
	log     log.Logger
	auth    *security.Authentication
}

// NewDiscussionController returns new instance of DiscussionController
func NewDiscussionController(discussionService *service.DiscussionService,
	log log.Logger, auth *security.Authentication) *DiscussionController {
	return &DiscussionController{
		service: discussionService,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *DiscussionController) RegisterRoutes(router *mux.Router) {

	unguarded := router.PathPrefix("/tenant/{tenantID}").Subrouter()
	guarded := router.PathPrefix("/tenant/{tenantID}").Subrouter()

	guarded.HandleFunc("/discussion",
		controller.AddDiscussion).Methods(http.MethodPost)

	guarded.HandleFunc("/discussions/{discussionID}",
		controller.UpdateDiscussion).Methods(http.MethodPut)
	guarded.HandleFunc("/discussions/{discussionID}",
		controller.DeleteDiscussion).Methods(http.MethodDelete)

	unguarded.HandleFunc("/discussions/{discussionID}", controller.GetDiscussion).
		Methods(http.MethodGet)

		// channels are not included as we need to get discussions by talent
	unguarded.HandleFunc("/discussions",
		controller.GetDiscussions).Methods(http.MethodGet)

	guarded.Use(controller.auth.Middleware)

	controller.log.Info("Discussion Routes Registered")
}

// AddDiscussion godoc
// AddDiscussion Add New Discussion
// @Description Add New Discussion
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param discussion body community.Discussion true "Add Discussion"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /discussion [POST]
func (controller *DiscussionController) AddDiscussion(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add Discussion Called==============================")
	discussion := &community.Discussion{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, discussion)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data",
			http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	discussion.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	discussion.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// // Parse and set channelID.
	// discussion.ChannelID, err = parser.GetUUID("channelID")
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, errors.NewHTTPError("unable to parse channel id", http.StatusBadRequest))
	// 	return
	// }

	err = discussion.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(),
			http.StatusBadRequest))
		return
	}

	err = controller.service.AddDiscussion(discussion)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(),
			http.StatusBadRequest))
		return
	}
	web.RespondJSON(w, http.StatusOK, "Discussion added successfully.")
}

// UpdateDiscussion godoc
// UpdateDiscussion Update Discussion By Discussion ID
// @Description Update Discussion By Discussion ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param discussion body community.Discussion true "Update Discussion Data"
// @Param discussion path string true "Discussion ID"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /discussion/{discussionID} [PUT]
func (controller *DiscussionController) UpdateDiscussion(w http.ResponseWriter,
	r *http.Request) {
	controller.log.Info("==============================Update Discussion Called==============================")
	discussion := &community.Discussion{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, discussion)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data",
			http.StatusBadRequest))
		return
	}

	discussion.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// UpdatedBy will be set with the given credentialID.
	discussion.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// // Parse and set channelID.
	// discussion.ChannelID, err = parser.GetUUID("channelID")
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, errors.NewHTTPError("unable to parse channel id", http.StatusBadRequest))
	// 	return
	// }

	discussion.ID, err = parser.GetUUID("discussionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse discusssion id",
			http.StatusBadRequest))
		return
	}

	err = discussion.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data",
			http.StatusBadRequest))
		return
	}

	err = controller.service.UpdateDiscussion(discussion)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(),
			http.StatusBadRequest))
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Discussion updated")
}

// DeleteDiscussion godoc
// DeleteDiscussion Delete Discussion By Discussion ID
// @Description DeleteDiscussion Discussion By Discussion ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param discussion path string true "Discussion ID"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /discussion/{discussionID} [DELETE]
func (controller *DiscussionController) DeleteDiscussion(w http.ResponseWriter,
	r *http.Request) {
	controller.log.Info("==============================Delete Discussion Called==============================")
	discussion := &community.Discussion{}
	var err error
	parser := web.NewParser(r)

	discussion.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	discussion.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// // Parse and set channelID.
	// discussion.ChannelID, err = parser.GetUUID("channelID")
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, errors.NewHTTPError("unable to parse channel id", http.StatusBadRequest))
	// 	return
	// }

	discussion.ID, err = parser.GetUUID("discussionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse discusssion id",
			http.StatusBadRequest))
		return
	}

	err = controller.service.DeleteDiscussion(discussion)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to delete discussion",
			http.StatusBadRequest))
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Discussion deleted")
}

// // GetDiscussionsByChannelID godoc
// // GetDiscussionsByChannelID Return Discussion By CHannel ID
// // @Description Return Discussion By Channel ID
// // @Tags community-forum
// // @Accept  json
// // @Produce  json
// // @Param limit path int true "total number of result" Format(int)
// // @Param offset path int true "page number" Format(int)
// // @Param channelID path string true "Get Discussion By Channel id" Format(uuid.UUID)
// // @Success 200 {array} []community.Discussion
// // @Failure 400 {object} errors.ValidationError
// // @Failure 500 {object} errors.HTTPError
// // @Router /discussion/channel/{channelID}/{limit}/{offset} [GET]
// func (controller *DiscussionController) GetDiscussionsByChannelID(w http.ResponseWriter,
// 	r *http.Request) {
// 	controller.log.Info("==============================GetDiscussionsByChannelID Called==============================")
// 	params := mux.Vars(r)

// 	tenantID, err := util.ParseUUID(params["tenantID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}

// 	channelID, err := util.ParseUUID(params["channelID"])
// 	if err != nil {
// 		web.RespondError(w, errors.NewHTTPError("unable to parse channel ID", http.StatusBadRequest))
// 		return
// 	}

// 	discussions := []community.DiscussionDTO{}
// 	err = controller.DiscussionService.GetDiscussionsByChannelID(tenantID, channelID, discussions)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, discussions)
// }

// GetDiscussion godoc
// GetDiscussion Return Discussion By Discussion ID
// @Description Return Discussion By Channel ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param discussionID path string true "Get Discussion By Discussion id" Format(uuid.UUID)
// @Success 200 {array} []community.Discussion
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /discussion/{discussionID} [GET]
func (controller *DiscussionController) GetDiscussion(w http.ResponseWriter,
	r *http.Request) {
	controller.log.Info("==============================GetDiscussion Called==============================")
	discussion := &community.DiscussionDTO{}
	var err error
	parser := web.NewParser(r)

	discussion.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// // Parse and set channelID.
	// discussion.ChannelID, err = parser.GetUUID("channelID")
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, errors.NewHTTPError("unable to parse channel id", http.StatusBadRequest))
	// 	return
	// }

	discussion.ID, err = parser.GetUUID("discussionID")
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse discussion ID",
			http.StatusBadRequest))
		return
	}
	err = controller.service.GetDiscussion(discussion)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, discussion)
}

// GetDiscussions godoc
// GetDiscussions Return Discussion By Discussion ID
// @Description Return Discussion By Talent ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param talentID path string true "Get Discussion By Talent id" Format(uuid.UUID)
// @Success 200 {array} []community.Discussion
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /discussion/talent/{talentID} [GET]
func (controller *DiscussionController) GetDiscussions(w http.ResponseWriter,
	r *http.Request) {
	controller.log.Info("==============================GetDiscussions Called==============================")
	parser := web.NewParser(r)

	// pagination := &general.Pagination{}
	// pagination.Limit, pagination.Offset = web.GetLimitAndOffset(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// talentID, err := util.ParseUUID(params["talentID"])
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, errors.NewHTTPError("unable to parse talent ID", http.StatusBadRequest))
	// 	return
	// }
	var totalCount int
	discussions := &[]community.DiscussionDTO{}

	// err = controller.DiscussionService.GetDiscussions(pagination, r.Form, discussions)
	err = controller.service.GetDiscussions(tenantID, discussions, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Some error occurred", http.StatusBadRequest))
		return
	}
	// web.RespondJSON(w, http.StatusOK, discussions)
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, discussions)
}
