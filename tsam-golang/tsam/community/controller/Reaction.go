package controller

import (
	"net/http"

	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/community/service"
)

// ReactionController Give Access to Update, Delete, Add, Get and Get All Reaction
type ReactionController struct {
	reactionService *service.ReactionService
	log             log.Logger
	auth            *security.Authentication
}

// NewReactionController returns new instance of ReactionController
func NewReactionController(service *service.ReactionService,
	log log.Logger, auth *security.Authentication) *ReactionController {
	return &ReactionController{
		reactionService: service,
		log:             log,
		auth:            auth,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *ReactionController) RegisterRoutes(route *mux.Router) {
	route.HandleFunc("/tenant/{tenantID}/reaction", controller.AddReaction).
		Methods(http.MethodPost)
	// route.HandleFunc("/reaction", controller.GetReactions).
	// 	Methods(http.MethodGet)
	// route.HandleFunc("/reaction/{reactionID}", controller.GetReaction).
	// 	Methods(http.MethodGet)
	route.HandleFunc("/tenant/{tenantID}/reaction/{reactionID}", controller.UpdateReaction).
		Methods(http.MethodPut)
	route.HandleFunc("/tenant/{tenantID}/reaction/{reactionID}", controller.DeleteReaction).
		Methods(http.MethodDelete)
	// route.HandleFunc("/reaction/count/replyID/{replyID}", controller.GetTotalReactionByReplyID).
	// 	Methods(http.MethodGet)
	controller.log.Info("Reaction routes registered")
}

// AddReaction Add New Reactions
func (controller *ReactionController) AddReaction(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddReaction Called==============================")

	reaction := community.Reaction{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &reaction)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	reaction.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	reaction.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	err = reaction.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	err = controller.reactionService.AddReaction(&reaction)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, reaction.ID)
}

// UpdateReaction updates reaction.
func (controller *ReactionController) UpdateReaction(w http.ResponseWriter,
	r *http.Request) {
	controller.log.Info("==============================Update Reaction Called==============================")
	reaction := &community.Reaction{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, reaction)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data",
			http.StatusBadRequest))
		return
	}

	reaction.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// UpdatedBy will be set with the given credentialID.
	reaction.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
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

	reaction.ID, err = parser.GetUUID("reactionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse reaction id",
			http.StatusBadRequest))
		return
	}

	err = reaction.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(),
			http.StatusBadRequest))
		return
	}

	err = controller.reactionService.UpdateReaction(reaction)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(),
			http.StatusBadRequest))
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Reaction updated")
}

// DeleteReaction Delete Reaction By Reaction ID
func (controller *ReactionController) DeleteReaction(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteReaction Called==============================")

	reaction := community.Reaction{}
	var err error
	parser := web.NewParser(r)

	reaction.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	reaction.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	reaction.ID, err = parser.GetUUID("reactionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse discusssion id",
			http.StatusBadRequest))
		return
	}

	err = controller.reactionService.DeleteReaction(&reaction)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
}

// // GetReaction Return Reaction By ID
// func (controller *ReactionController) GetReaction(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================GetReaction Called==============================")

// 	reaction := community.Reaction{}
// 	var err error

// 	reaction.ID, err = util.ParseUUID(mux.Vars(r)["reactionID"])
// 	if err != nil {
// 		web.RespondError(w, errors.NewHTTPError("unable to parse Reaction id", http.StatusBadRequest))
// 		return
// 	}
// 	err = controller.reactionService.GetReaction(&reaction)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, &reaction)
// }

// // GetReactions Return All Reaction
// func (controller *ReactionController) GetReactions(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================GetReactions Called==============================")

// 	reactions := []community.Reaction{}
// 	err := controller.reactionService.GetReactions(&reactions)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, &reactions)
// }

// // GetTotalReactionByReplyID Return Total Reaction By Reply ID
// func (controller *ReactionController) GetTotalReactionByReplyID(w http.ResponseWriter, r *http.Request) {

// 	replyID, err := util.ParseUUID(mux.Vars(r)["replyID"])
// 	if err != nil {
// 		web.RespondError(w, errors.NewHTTPError("unable to parse Reply ID", http.StatusBadRequest))
// 		return
// 	}
// 	var totalCount int
// 	err = controller.reactionService.GetReactionCountByReplyID(&totalCount, &replyID)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, &totalCount)
// }
