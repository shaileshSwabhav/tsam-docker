package module

import (
	"github.com/techlabs/swabhav/tsam"
	communitycon "github.com/techlabs/swabhav/tsam/community/controller"
	communityser "github.com/techlabs/swabhav/tsam/community/service"
	"github.com/techlabs/swabhav/tsam/repository"
)

func registerCommunityRoutes(app *tsam.App, repository repository.Repository) {
	defer app.WG.Done()
	// Community Module
	log := app.Log

	discussionService := communityser.NewDiscussionService(app.DB, repository)
	discussionController := communitycon.NewDiscussionController(discussionService, log, app.Auth)
	replyService := communityser.NewReplyService(app.DB, repository)
	replyController := communitycon.NewReplyController(replyService, log, app.Auth)
	reactionService := communityser.NewReactionService(app.DB, repository)
	reactionController := communitycon.NewReactionController(reactionService, log, app.Auth)
	commentService := communityser.NewCommentService(app.DB, repository)
	commentController := communitycon.NewCommentController(commentService)
	notificationService := communityser.NewNotificationService(app.DB, repository)
	notificationController := communitycon.NewNotificationController(notificationService)
	notificationTypeService := communityser.NewNotificationTypeService(app.DB, repository)
	notificationTypeController := communitycon.NewNotificationTypeController(notificationTypeService)
	channelService := communityser.NewChannelService(app.DB, repository)
	channelController := communitycon.NewChannelController(channelService, log, app.Auth)

	// #niranjan new router
	app.RegisterControllerRoutes([]tsam.Controller{
		channelController,
		discussionController,
	})

	// replyController, reactionController,

	test := []interface{}{replyController,
		commentController, notificationController, notificationTypeController,
		reactionController, channelController, discussionController}
	_ = test
}
