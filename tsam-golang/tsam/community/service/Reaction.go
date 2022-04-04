package service

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ReactionService gives access to CRUD operations.
type ReactionService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

const reactionAddedToReply string = "New reaction added"

// NewReactionService returns an instance of ReactionService.
func NewReactionService(db *gorm.DB, repo repository.Repository) *ReactionService {
	return &ReactionService{
		DB:         db,
		Repository: repo,
	}
}

// AddReaction Add New Reaction to Database
func (service *ReactionService) AddReaction(reaction *community.Reaction) error {

	// check if all foreign keys exist.
	err := service.doForeignKeysExist(reaction.CreatedBy, reaction)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	fmt.Println("BEfore validate =====================")
	// Validates and returns error if fields have improper value.
	err = service.validateFields(uow, reaction)
	if err != nil {
		return err
	}
	fmt.Println("After validate =====================")

	err = service.Repository.Add(uow, reaction)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// eventpool := event.NewPool()
	// reactionNotification := notification.Reaction{
	// 	Notification: notification.Notification{
	// 		TypeID:     345,
	// 		SeenTime:   time.Now(),
	// 		NotifierID: reaction.CreatedBy,
	// 		NotifiedID: []uuid.UUID{},
	// 	},
	// }
	// eventpool.AssignJob()

	// notificationType := &community.NotificationType{}
	// err = service.Repository.GetRecord(uow, notificationType,
	// 	repository.Filter("`notification_type_name`=?", reactionAddedToReply))
	// if err != nil {
	// 	uow.RollBack()
	// 	return errors.NewHTTPError("notification type error: "+err.Error(), http.StatusInternalServerError)
	// }

	// notificationService := NewNotificationService(service.DB, service.Repository)

	// notification := &community.Notification{
	// 	DiscussionID:       tempReply.DiscussionID,
	// 	ReplyID:            reaction.ReplyID,
	// 	NotifierTalentID:   reaction.TalentID,
	// 	NotifiedTalentID:   tempReply.ReplierID,
	// 	NotificationTypeID: notificationType.ID,
	// }

	// err = notificationService.AddNotificationWithoutTransaction(notification, uow)
	// if err != nil {
	// 	uow.RollBack()
	// 	return errors.NewHTTPError("Can't add notification: "+err.Error(), http.StatusInternalServerError)
	// }
	// uow.Commit()
	return nil
}

// UpdateReaction Update Reaction
func (service *ReactionService) UpdateReaction(reaction *community.Reaction) error {

	// Check if reaction exist.
	err := service.doesReactionExist(reaction.TenantID, reaction.ID)
	if err != nil {
		return err
	}

	// check if all foreign keys exist.
	err = service.doForeignKeysExist(reaction.UpdatedBy, reaction)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFields(uow, reaction)
	if err != nil {
		return err
	}

	err = service.Repository.Update(uow, reaction)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// DeleteReaction delete specified reaction.
func (service *ReactionService) DeleteReaction(reaction *community.Reaction) error {

	// Check if reaction exist
	err := service.doesReactionExist(reaction.TenantID, reaction.ID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(reaction.TenantID, reaction.DeletedBy); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Delete(uow, reaction)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// #Niranjan : need to check this.
// validateFields checks field uniqueness.
func (service *ReactionService) validateFields(uow *repository.UnitOfWork, reaction *community.Reaction) error {
	fmt.Println("Stage 1===========================================")

	if reaction.ReplyID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, reaction.TenantID, community.Reaction{},
			repository.Filter("`created_by` = ? AND `is_liked` = ? AND `reply_id` = ? AND `id`!= ?",
				reaction.CreatedBy, reaction.IsLiked, reaction.ReplyID, reaction.ID))
		err = util.HandleIfExistsError("Already reacted", exists, err)
		if err != nil {
			return errors.NewValidationError(err.Error())
		}
		tempReply := community.Reply{}
		err = service.Repository.GetForTenant(uow, reaction.TenantID, *reaction.ReplyID,
			&tempReply, repository.Select("`replier_id`"))
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusNotFound)
		}

		if tempReply.ReplierID == reaction.CreatedBy {
			return errors.NewHTTPError("Can't react to your own post", http.StatusBadRequest)
		}
	}

	if reaction.DiscussionID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, reaction.TenantID, community.Reaction{},
			repository.Filter("`created_by` = ? AND `is_liked` = ? AND `discussion_id` = ? AND `id`!= ?",
				reaction.CreatedBy, reaction.IsLiked, reaction.DiscussionID, reaction.ID))
		err = util.HandleIfExistsError("Already reacted", exists, err)
		if err != nil {
			return errors.NewValidationError(err.Error())
		}
		tempDiscussion := community.Discussion{}
		err = service.Repository.GetForTenant(uow, reaction.TenantID, *reaction.DiscussionID,
			&tempDiscussion, repository.Select("`author_id`"))
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusNotFound)
		}

		if tempDiscussion.AuthorID == reaction.CreatedBy {
			return errors.NewHTTPError("Can't react to your own post", http.StatusBadRequest)
		}
	}
	return nil
}

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *ReactionService) doForeignKeysExist(credentialID uuid.UUID, reaction *community.Reaction) error {
	tenantID := reaction.TenantID

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *ReactionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *ReactionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesReactorExist will check if reactorID is valid credentialID.
func (service *ReactionService) doesReactorExist(tenantID, reactorID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Credential{},
		repository.Filter("`id`=?", reactorID))
	if err := util.HandleError("Invalid reactor ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesReplyExist will check if replyID is valid
func (service *ReactionService) doesReplyExist(tenantID, replyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Reply{},
		repository.Filter("`id`=?", replyID))
	if err := util.HandleError("Invalid reply ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDiscussionExist will check if discussionID is valid
func (service *ReactionService) doesDiscussionExist(tenantID, discussionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Discussion{},
		repository.Filter("`id`=?", discussionID))
	if err := util.HandleError("Invalid discussion id", exists, err); err != nil {
		return err
	}
	return nil
}

// doesReactionExist will check if replyID is valid
func (service *ReactionService) doesReactionExist(tenantID, reactionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, community.Reaction{},
		repository.Filter("`id`=?", reactionID))
	if err := util.HandleError("Invalid reaction ID", exists, err); err != nil {
		return err
	}
	return nil
}
