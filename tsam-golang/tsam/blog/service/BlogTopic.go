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

// BlogTopicService provide method to update, delete, add, get method for blog topic.
type BlogTopicService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewBlogTopicService creates new instance of BlogTopic.
func NewBlogTopicService(db *gorm.DB, repository repository.Repository) *BlogTopicService {
	return &BlogTopicService{
		DB:         db,
		Repository: repository,
	}
}

// AddBlogTopic adds new blog topic to database.
func (service *BlogTopicService) AddBlogTopic(blogTopic *blog.BlogTopic) error {

	// Validate tenant id.
	err := service.doesTenantExist(blogTopic.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate if same blog topic name exists.
	err = service.doesNameExist(blogTopic)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(blogTopic.CreatedBy, blogTopic.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	//  Creating unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add blog topic.
	err = service.Repository.Add(uow, blogTopic)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateBlogTopic updates blog topic in database.
func (service *BlogTopicService) UpdateBlogTopic(blogTopic *blog.BlogTopic) error {

	// Validate tenant ID.
	err := service.doesTenantExist(blogTopic.TenantID)
	if err != nil {
		return err
	}

	// Validate blog topic ID.
	err = service.doesBlogTopicExist(blogTopic.ID, blogTopic.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(blogTopic.UpdatedBy, blogTopic.TenantID)
	if err != nil {
		return err
	}

	// Validate if same blog topic name exists.
	err = service.doesNameExist(blogTopic)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update blog topic.
	err = service.Repository.Update(uow, blogTopic)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteBlogTopic delete blog topic from database.
func (service *BlogTopicService) DeleteBlogTopic(blogTopic *blog.BlogTopic) error {
	credentialID := blogTopic.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(blogTopic.TenantID)
	if err != nil {
		return err
	}

	// Validate blog topic ID.
	err = service.doesBlogTopicExist(blogTopic.ID, blogTopic.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, blogTopic.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update blog topic for updating deleted_by and deleted_at fields of blog topic.
	if err := service.Repository.UpdateWithMap(uow, blogTopic, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", blogTopic.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog Topic could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetBlogTopic returns one blog topic by id.
func (service *BlogTopicService) GetBlogTopic(blogTopic *blog.BlogTopic) error {

	// Validate tenant ID.
	err := service.doesTenantExist(blogTopic.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get blog topic by id from database.
	err = service.Repository.GetForTenant(uow, blogTopic.TenantID, blogTopic.ID, blogTopic)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBlogTopics returns all blog topics.
func (service *BlogTopicService) GetBlogTopics(blogTopics *[]blog.BlogTopic,
	form url.Values, tenantID uuid.UUID, limit, offset int, totalCount *int) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all blog topic.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, blogTopics, "`name`",
		service.addSearchQueries(form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBlogTopicList returns all the blog topic(without pagination).
func (service *BlogTopicService) GetBlogTopicList(blogTopics *[]blog.BlogTopic,
	tenantID uuid.UUID) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all blog topic.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, blogTopics, "`name`")
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds all search queries by comparing with the blog topic data.
func (service *BlogTopicService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
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

// doesNameExist returns true if the name already exists for the blog topic in database.
func (service *BlogTopicService) doesNameExist(blogTopic *blog.BlogTopic) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, blogTopic.TenantID, &blog.BlogTopic{},
		repository.Filter("`name`=? AND `id`!=?", blogTopic.Name, blogTopic.ID))
	if err := util.HandleIfExistsError("Name : "+blogTopic.Name+" exists", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *BlogTopicService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBlogTopicExist validates if blog topic exists or not in database.
func (service *BlogTopicService) doesBlogTopicExist(blogTopicID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, blog.BlogTopic{},
		repository.Filter("`id` = ?", blogTopicID))
	if err := util.HandleError("Invalid blog topic ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *BlogTopicService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
