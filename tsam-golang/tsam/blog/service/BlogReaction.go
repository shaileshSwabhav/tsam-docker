package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/blog"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// BlogReactionService provide method to update, delete, add, get method for blog reaction.
type BlogReactionService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewBlogReactionService creates new instance of BlogReactionService.
func NewBlogReactionService(db *gorm.DB, repository repository.Repository) *BlogReactionService {
	return &BlogReactionService{
		DB:         db,
		Repository: repository,
	}
}

// BlogReactionAssociationNames provides preload associations array for blog reaction.
var BlogReactionAssociationNames []string = []string{
	"Reactor",
}

// AddBlogReaction adds new blog reaction to database.
func (service *BlogReactionService) AddBlogReaction(blogReaction *blog.BlogReaction) error {

	// Validate tenant id.
	err := service.doesTenantExist(blogReaction.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(blogReaction.CreatedBy, blogReaction.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// If blog id exists then validate blog id.
	if blogReaction.BlogID != nil {
		err = service.doesBlogExist(*blogReaction.BlogID, blogReaction.TenantID)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// If reply id exists then validate reply id.
	if blogReaction.ReplyID != nil {
		err = service.doesBlogReplyExist(*blogReaction.ReplyID, blogReaction.TenantID)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Validate reactor ID.
	err = service.doesCredentialExist(*blogReaction.ReactorID, blogReaction.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	//  Creating unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add blog reaction.
	err = service.Repository.Add(uow, blogReaction)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateBlogReaction updates blog reaction in database.
func (service *BlogReactionService) UpdateBlogReaction(blogReaction *blog.BlogReaction) error {

	// Validate tenant ID.
	err := service.doesTenantExist(blogReaction.TenantID)
	if err != nil {
		return err
	}

	// Validate blog reaction ID.
	err = service.doesBlogReactionExist(blogReaction.ID, blogReaction.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(blogReaction.UpdatedBy, blogReaction.TenantID)
	if err != nil {
		return err
	}

	// If blog id exists then validate blog id.
	if blogReaction.BlogID != nil {
		err = service.doesBlogExist(*blogReaction.BlogID, blogReaction.TenantID)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// If reply id exists then validate reply id.
	if blogReaction.ReplyID != nil {
		err = service.doesBlogReplyExist(*blogReaction.ReplyID, blogReaction.TenantID)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Validate reactor ID.
	err = service.doesCredentialExist(*blogReaction.ReactorID, blogReaction.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting blog reaction already present in database.
	tempBlogReaction := blog.BlogReaction{}

	// Get blog reaction for getting created_by field of blog reaction from database.
	if err := service.Repository.GetForTenant(uow, blogReaction.TenantID, blogReaction.ID, &tempBlogReaction); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp blog reaction to blog reaction to be updated.
	blogReaction.CreatedBy = tempBlogReaction.CreatedBy

	// Update blog reaction.
	err = service.Repository.Save(uow, &blogReaction)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteBlogReaction delete blog reaction from database.
func (service *BlogReactionService) DeleteBlogReaction(blogReaction *blog.BlogReaction) error {
	credentialID := blogReaction.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(blogReaction.TenantID)
	if err != nil {
		return err
	}

	// Validate blog reaction ID.
	err = service.doesBlogReactionExist(blogReaction.ID, blogReaction.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, blogReaction.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update blog reaction for updating deleted_by and deleted_at fields of blog reaction.
	if err := service.Repository.UpdateWithMap(uow, blogReaction, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", blogReaction.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog reaction could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetBlogReaction returns one blog reaction by id.
func (service *BlogReactionService) GetBlogReaction(blogReaction *blog.BlogReactionDTO) error {

	// Validate tenant ID.
	err := service.doesTenantExist(blogReaction.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get blog reaction by id from database.
	err = service.Repository.GetForTenant(uow, blogReaction.TenantID, blogReaction.ID, blogReaction)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBlogReactions returns all blog reactions.
func (service *BlogReactionService) GetBlogReactions(blogReactions *[]blog.BlogReactionDTO,
	form url.Values, tenantID uuid.UUID) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all blog reactions.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, blogReactions, "`created_at`",
		service.addSearchQueries(form),
		repository.PreloadAssociations(BlogReactionAssociationNames))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds all search queries by comparing with the blog reaction data.
func (service *BlogReactionService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	// Blog id.
	if _, ok := requestForm["blogID"]; ok {
		util.AddToSlice("blog_id", "= ?", "AND", requestForm.Get("blogID"), &columnNames, &conditions, &operators, &values)
	}

	// Reply id.
	if _, ok := requestForm["replyID"]; ok {
		util.AddToSlice("reply_id", "= ?", "AND", requestForm.Get("replyID"), &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}


// doesTenantExist validates if tenant exists or not in database.
func (service *BlogReactionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBlogReactionExist validates if blog reaction exists or not in database.
func (service *BlogReactionService) doesBlogReactionExist(blogReactionID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, blog.BlogReaction{},
		repository.Filter("`id` = ?", blogReactionID))
	if err := util.HandleError("Invalid blog reaction ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBlogReplyExist validates if blog reply exists or not in database.
func (service *BlogReactionService) doesBlogReplyExist(blogReplyID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, blog.BlogReply{},
		repository.Filter("`id` = ?", blogReplyID))
	if err := util.HandleError("Invalid blog reply ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBlogExist validates if blog exists or not in database.
func (service *BlogReactionService) doesBlogExist(blogID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, blog.Blog{},
		repository.Filter("`id` = ?", blogID))
	if err := util.HandleError("Invalid blog ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *BlogReactionService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
