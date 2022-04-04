package controller

import (
	"net/http"

	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/util"

	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/web"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/community/service"
	"github.com/techlabs/swabhav/tsam/log"
)

// CommentController gives access to CRUD operations for entity
type CommentController struct {
	CommentService *service.CommentService
}

// NewCommentController returns new instance of CommentController
func NewCommentController(commentService *service.CommentService) *CommentController {
	return &CommentController{
		CommentService: commentService,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (commentController *CommentController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/comment/reply/{"+paramReplyID+"}", commentController.AddComment).Methods(http.MethodPost)
	log.NewLogger().Info("Comment Routes Registered")
}

// AddComment godoc
// AddComment Add New Comment
// @Description Add New Comment
// @Tags community-forum
// @Accept  json
// @Produce  json
// @Param reply body community.Reply true "Add New Comment"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /reply [POST]
func (commentController *CommentController) AddComment(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddComment called==============================")
	comment := &community.Reply{}

	replyID, err := util.ParseUUID(mux.Vars(r)[paramReplyID])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(),
			http.StatusBadRequest))
		return
	}

	err = web.UnmarshalJSON(r, comment)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(),
			http.StatusBadRequest))
		return
	}

	err = comment.Validate()
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(),
			http.StatusBadRequest))
		return
	}

	comment.ReplyID = &replyID

	err = commentController.CommentService.AddComment(comment)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, comment.ID)
}
