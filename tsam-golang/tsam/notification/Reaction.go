package notification

import (
	"fmt"
	"runtime"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// Reaction is used to send Reaction notification data.
type Reaction struct {
	Notification
	ID           uuid.UUID
	DiscussionID *uuid.UUID
	ReplyID      *uuid.UUID
}

// Below function(line 35) can be named as NewReactionOnReply if needed.
// func NewReaction(seentime time.Time,notificationTypeID,notifierID,reactionID,discussionID uuid.UUID) *reaction{
// 	return &reaction{
// 		Notification: Notification{
// 			SeenTime: seentime,
// 			TypeID: notificationTypeID,
// 			NotifierID: notifierID,
// 		},
// 		ID: reactionID,
// 		DiscussionID: discussionID,
// 	}
// }

// NewReaction must be called with replyID = nil if the reaction is on the discussion.
func NewReaction(tenantID, notifierID, reactionID uuid.UUID,
	discussionID, replyID *uuid.UUID) *Reaction {
	r := &Reaction{
		ID:           reactionID,
		DiscussionID: discussionID,
		ReplyID:      replyID,
	}
	r.NotifierID = notifierID
	r.TenantID = tenantID
	return r
}

// Needs a lot of changing. (There is a transaction commit issue)
func (r *Reaction) Notify(db *gorm.DB, repo repository.Repository) {
	fmt.Println("In Fired notify-----------", db)
	uow := repository.NewUnitOfWork(db, false)
	fmt.Println("After transaction begins-----------")
	notifiedID := struct{ AuthorID uuid.UUID }{}

	err := repo.GetRecordForTenant(uow, r.TenantID, &notifiedID,
		repository.Table("community_discussions"),
		repository.Filter("`id` = ?", r.DiscussionID),
		repository.Select("`author_id`"))
	if err != nil {
		fmt.Println("error man--------------------------", err)
		return
	}
	notificationID := util.GenerateUUID()
	n := Notification{
		TenantBase: general.TenantBase{
			ID:         notificationID,
			IgnoreHook: true,
			CreatedBy:  r.CreatedBy,
			TenantID:   r.TenantID,
		},
		NotifierID: r.CreatedBy,
		NotifiedID: notifiedID.AuthorID,
	}

	err = repo.Add(uow, n)
	if err != nil {
		fmt.Println("add error man--------------------------", err)
		return
	}
	uow.Commit()

	rn := Reaction{
		Notification: n,
		ID:           notificationID,
		DiscussionID: r.DiscussionID,
	}
	_ = rn
	fmt.Println("Go routine Number is ---------------------------------------------------------", runtime.NumGoroutine())

	// Get replyID & ID 's replierID
	//  Get discussion author ID
	// keep them in one array
	// or if only notifications to direct contact
}

func (r *Reaction) addNotification() {

}
