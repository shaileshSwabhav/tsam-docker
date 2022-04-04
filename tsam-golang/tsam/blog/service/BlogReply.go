package service

import (
	// "encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/blog"
	"github.com/techlabs/swabhav/tsam/models/general"
	// notfication_service "github.com/techlabs/swabhav/tsam/notification_demo/service"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// BlogReplyService provide method to update, delete, add, get method for blog reply.
type BlogReplyService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewBlogReplyService creates new instance of BlogReplyService.
func NewBlogReplyService(db *gorm.DB, repository repository.Repository) *BlogReplyService {
	return &BlogReplyService{
		DB:         db,
		Repository: repository,
	}
}

// BlogReplyAssociationNames provides preload associations array for blog reply.
var BlogReplyAssociationNames []string = []string{
	"Replier", "Replier.Role", "Replies.Replier", "Replies.Replier.Role", "Reactions", "Reactions.Reactor",
	"Reactions.Reactor.Role", "Replies.Reactions", "Replies.Reactions.Reactor", "Replies.Reactions.Reactor.Role",
}

// AddBlogReply adds new blog reply to database.
func (service *BlogReplyService) AddBlogReply(blogReply *blog.BlogReply) error {

	// Validate tenant id.
	err := service.doesTenantExist(blogReply.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate if same blog reply name exists.
	err = service.doesReplyExist(blogReply)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(blogReply.CreatedBy, blogReply.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// If blog id exists then validate blog id.
	if blogReply.BlogID != nil {
		err = service.doesBlogExist(*blogReply.BlogID, blogReply.TenantID)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// If reply id exists then validate reply id.
	if blogReply.ReplyID != nil {
		err = service.doesBlogReplyExist(*blogReply.ReplyID, blogReply.TenantID)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Validate replier ID.
	err = service.doesCredentialExist(*blogReply.ReplierID, blogReply.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	//  Creating unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add blog reply.
	err = service.Repository.Add(uow, blogReply)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	//Add Notification
	blogNotification := blog.BlogNotification{
		BlogID:      blogReply.BlogID,
		BlogReplyID: &blogReply.ID,
	}
	blogNotification.TenantID = blogReply.TenantID

	// blog := blog.Blog{}
	// err = service.Repository.GetRecordForTenant(uow, blogReply.TenantID, &blog,
	// 	repository.Select("author_id"),
	// 	repository.Filter("id=?", blogReply.BlogID))
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	uow.RollBack()
	// 	return err
	// }
	// notification := general.Notification_Test{

	// 	NotifierID: *blogReply.ReplierID,
	// 	NotifiedID: blog.AuthorID,
	// }
	// notification.TenantID = blogReply.TenantID
	// notificationService := notfication_service.NewNotificationService(service.DB, service.Repository)
	// childTable, err := json.Marshal(&blogNotification)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	uow.RollBack()
	// 	return err
	// }
	// err = notificationService.AddNotification(&childTable, &notification, uow)

	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateBlogReply updates blog reply in database.
func (service *BlogReplyService) UpdateBlogReply(blogReply *blog.BlogReply) error {

	// Validate tenant ID.
	err := service.doesTenantExist(blogReply.TenantID)
	if err != nil {
		return err
	}

	// Validate blog reply ID.
	err = service.doesBlogReplyExist(blogReply.ID, blogReply.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(blogReply.UpdatedBy, blogReply.TenantID)
	if err != nil {
		return err
	}

	// Validate if same blog reply reply exists.
	err = service.doesReplyExist(blogReply)
	if err != nil {
		return err
	}

	// If blog id exists then validate blog id.
	if blogReply.BlogID != nil {
		err = service.doesBlogExist(*blogReply.BlogID, blogReply.TenantID)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// If reply id exists then validate reply id.
	if blogReply.ReplyID != nil {
		err = service.doesBlogReplyExist(*blogReply.ReplyID, blogReply.TenantID)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Validate replier ID.
	err = service.doesCredentialExist(*blogReply.ReplierID, blogReply.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update blog reply.
	err = service.Repository.Update(uow, blogReply)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteBlogReply delete blog reply from database.
func (service *BlogReplyService) DeleteBlogReply(blogReply *blog.BlogReply) error {
	credentialID := blogReply.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(blogReply.TenantID)
	if err != nil {
		return err
	}

	// Validate blog reply ID.
	err = service.doesBlogReplyExist(blogReply.ID, blogReply.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, blogReply.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	//  Delete reply association from database.
	if err := service.DeleteReplyAssociation(uow, blogReply, credentialID); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Failed to delete blog reply", http.StatusInternalServerError)
	}

	// Update blog reply for updating deleted_by and deleted_at fields of blog reply.
	if err := service.Repository.UpdateWithMap(uow, blogReply, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", blogReply.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog reply could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetBlogReply returns one blog reply by id.
func (service *BlogReplyService) GetBlogReply(blogReply *blog.BlogReplyDTO) error {

	// Validate tenant ID.
	err := service.doesTenantExist(blogReply.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get blog reply by id from database.
	err = service.Repository.GetForTenant(uow, blogReply.TenantID, blogReply.ID, blogReply)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBlogReplies returns all blog replies.
func (service *BlogReplyService) GetBlogReplies(blogReplies *[]blog.BlogReplyDTO,
	form url.Values, tenantID, blogID uuid.UUID) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate blog id.
	err = service.doesBlogExist(blogID, tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all blog replies.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, blogReplies, "`created_at`",
		service.addSearchQueries(form),
		repository.Filter("`blog_id`=?", blogID),
		repository.Filter("reply_id IS NULL AND is_verified=1"),
		repository.PreloadAssociations(BlogReplyAssociationNames),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "Replies",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("replies.`is_verified`=1"),
				repository.OrderBy("replies.`created_at`")},
		}))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteReplyAssociation deletes reply's associations.
func (service *BlogReplyService) DeleteReplyAssociation(uow *repository.UnitOfWork, blogReply *blog.BlogReply, credentialID uuid.UUID) error {

	//********************************************** Sub Reply Reactions ***************************************************

	subReplies := []blog.BlogReply{}

	// Get sub replies.
	if err := service.Repository.GetAllForTenant(uow, blogReply.TenantID, &subReplies,
		repository.Filter("`reply_id`=?", blogReply.ID)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Collect all sub replies' ids in a variable.
	var subReplyIDs []uuid.UUID
	for _, subReply := range subReplies {
		subReplyIDs = append(subReplyIDs, subReply.ID)
	}

	// Deleting sub replies' reactions.
	if err := service.Repository.UpdateWithMap(uow, &blog.BlogReaction{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`reply_id` IN (?) AND `tenant_id`=?", subReplyIDs, blogReply.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog reply could not be deleted", http.StatusInternalServerError)
	}

	// Deleting sub replies.
	if err := service.Repository.UpdateWithMap(uow, &blog.BlogReply{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`reply_id`=? AND `tenant_id`=?", blogReply.ID, blogReply.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog reply could not be deleted", http.StatusInternalServerError)
	}

	//********************************************** Blog Reply Reaction ***************************************************
	if err := service.Repository.UpdateWithMap(uow, &blog.BlogReaction{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`reply_id`=? AND `tenant_id`=?", blogReply.ID, blogReply.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog reply could not be deleted", http.StatusInternalServerError)
	}

	return nil
}

// addSearchQueries adds all search queries by comparing with the blog reply data.
func (service *BlogReplyService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	// Name.
	if _, ok := requestForm["name"]; ok {
		util.AddToSlice("`name`", "LIKE ?", "AND", "%"+requestForm.Get("name")+"%", &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesReplyExist returns true if the reply already exists for the blog reply in database.
func (service *BlogReplyService) doesReplyExist(blogReply *blog.BlogReply) error {

	// Create query processor according to parent reply or blog exists or not.
	var queryProcessor repository.QueryProcessor

	// If blog id is present.
	if blogReply.BlogID != nil {
		queryProcessor = repository.Filter("`reply`=? AND `id`!=? AND reply_id IS NULL AND `blog_id`=?",
			blogReply.Reply, blogReply.ID, blogReply.BlogID)
	}

	// If reply id is present.
	if blogReply.ReplyID != nil {
		queryProcessor = repository.Filter("`reply`=? AND `id`!=? AND blog_id IS NULL AND `reply_id`=?",
			blogReply.Reply, blogReply.ID, blogReply.ReplyID)
	}

	// Check for same reply conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, blogReply.TenantID, &blog.BlogReply{}, queryProcessor)
	if err := util.HandleIfExistsError("You have already replied with the same reply", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *BlogReplyService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBlogReplyExist validates if blog reply exists or not in database.
func (service *BlogReplyService) doesBlogReplyExist(blogReplyID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, blog.BlogReply{},
		repository.Filter("`id` = ?", blogReplyID))
	if err := util.HandleError("Invalid blog reply ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBlogExist validates if blog exists or not in database.
func (service *BlogReplyService) doesBlogExist(blogID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, blog.Blog{},
		repository.Filter("`id` = ?", blogID))
	if err := util.HandleError("Invalid blog ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *BlogReplyService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
