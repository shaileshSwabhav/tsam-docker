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

// ReplyController gives access to CRUD operations for entity
type ReplyController struct {
	service *service.ReplyService
	// DiscussionService *service.DiscussionService
	log    log.Logger
	jwtKey string
	auth   *security.Authentication
}

// NewReplyController returns new instance of ReplyController
func NewReplyController(replyService *service.ReplyService, log log.Logger, auth *security.Authentication) *ReplyController {
	return &ReplyController{
		service: replyService,
		log:     log,
		auth:    auth,
		// DiscussionService: discussionService,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *ReplyController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/tenant/{tenantID}/discussions/{discussionID}/reply",
		controller.AddReply).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/discussions/{discussionID}/replies/{replyID}",
		controller.UpdateReply).Methods(http.MethodPut)
	router.HandleFunc("/tenant/{tenantID}/discussions/{discussionID}/replies/{replyID}",
		controller.DeleteReply).Methods(http.MethodDelete)
	router.HandleFunc("/tenant/{tenantID}/discussions/{discussionID}/replies/{replyID}",
		controller.GetReply).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/discussions/{discussionID}/replies",
		controller.GetRepliesByDiscussionID).Methods(http.MethodGet)
	controller.log.Info("Reply Routes Registered")
}

// AddReply godoc
// AddReply Add New Reply
// @Description Add New Reply
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param reply body community.Reply true "Add New Reply"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /reply [POST]
func (controller *ReplyController) AddReply(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddReply called==============================")
	reply := community.Reply{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &reply)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data",
			http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	reply.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	reply.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set discussionID.
	reply.DiscussionID, err = parser.GetUUID("discussionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse discussion id", http.StatusBadRequest))
		return
	}

	err = reply.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(),
			http.StatusBadRequest))
		return
	}

	err = controller.service.AddReply(&reply)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Reply added successfully.")
}

// UpdateReply godoc
// UpdateReply Update Reply By Reply ID
// @Description Update Reply By Reply ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param reply body community.Reply true "Add Reply"
// @Param replyID path string true "Reply ID For Update Reply" Format(uuid.UUID)
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /reply/{replyID} [PUT]
func (controller *ReplyController) UpdateReply(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateReply called==============================")
	reply := community.Reply{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &reply)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data",
			http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	reply.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field.
	reply.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set discussionID.
	reply.DiscussionID, err = parser.GetUUID("discussionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse discussion id", http.StatusBadRequest))
		return
	}

	reply.ID, err = parser.GetUUID("replyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse reply id", http.StatusBadRequest))
		return
	}

	err = reply.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to get Reply ID", http.StatusBadRequest))
		return
	}

	err = controller.service.UpdateReply(&reply)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Reply Updated.")
}

// DeleteReply godoc
// DeleteReply Delete Reply By Reply ID
// @Description Delete Reply By Reply ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param replyID path string true "Reply ID For Delete Reply" Format(uuid.UUID)
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /reply/{replyID} [DELETE]
func (controller *ReplyController) DeleteReply(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteReply called==============================")
	reply := community.Reply{}
	var err error
	parser := web.NewParser(r)

	reply.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	reply.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	reply.DiscussionID, err = parser.GetUUID("discussionID")
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse discussion id",
			http.StatusBadRequest))
		return
	}

	reply.ID, err = parser.GetUUID("replyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse discusssion id",
			http.StatusBadRequest))
		return
	}

	err = controller.service.DeleteReply(&reply)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, "Reply deleted.")
}

// GetReply godoc
// GetReply Get Reply By Reply ID
// @Description GetReply Reply By Reply ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param replyID path string true "Reply ID For Get Reply" Format(uuid.UUID)
// @Success 200 {object} community.Reply
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /reply/{replyID} [GET]
func (controller *ReplyController) GetReply(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetReply called==============================")
	reply := community.Reply{}
	var err error
	parser := web.NewParser(r)

	reply.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	reply.ID, err = parser.GetUUID("replyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse reply id",
			http.StatusBadRequest))
		return
	}

	err = controller.service.GetReply(&reply)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
}

// GetRepliesByDiscussionID godoc
// GetRepliesByDiscussionID Return multiple Reply By Discussion ID
// @Description Return multiple Reply By Discussion ID
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param limit path int true "total number of result" Format(int)
// @Param offset path int true "page number" Format(int)
// @Param discussionid path string true "Get ReplyDTO By Discussion id" Format(uuid.UUID)
// @Success 200 {array} []community.ReplyDTO
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /reply/discussion/{discussionid}/{limit}/{offset} [GET]
func (controller *ReplyController) GetRepliesByDiscussionID(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetRepliesByDiscussionID called==============================")
	parser := web.NewParser(r)
	// ************************* Parsing optional params **********************
	// Check if talent has already liked the reply & set flag accordingly.(Do it in service)

	// var talentID uuid.UUID
	// if parsedTalentID := r.FormValue("talentID"); len(parsedTalentID) > 0 {
	// 	talentID, _ = util.ParseUUID(parsedTalentID)
	// }

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	discussionID, err := parser.GetUUID("discussionID")
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to get Discussion ID",
			http.StatusBadRequest))
		return
	}

	var totalCount int
	replies := &[]community.ReplyDTO{}
	// pageNum, err := strconv.ParseInt(params["offset"], 10, 64)
	// if err != nil {
	// 	controller.log.Error("unable to parse offset value form user error :=", err.Error())
	// 	pageNum = 0
	// }
	// replyDTO := community.ReplyDTO{}
	// if pageNum == 0 {
	// 	err = controller.DiscussionService.GetDiscussion(&discussion,
	// 		discussionID)
	// 	if err != nil {
	// 		web.RespondError(w, err)
	// 		return
	// 	}
	// 	replyDTO.DiscussionDTO = &discussion
	// }
	err = controller.service.GetRepliesByDiscussionID(tenantID, discussionID, replies, parser, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// replyDTO.ReplyDTO = replies
	web.RespondJSON(w, http.StatusOK, &replies)
}
