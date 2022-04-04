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

// BlogViewService provide method to update, delete, add, get method for blog viewer.
type BlogViewService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewBlogViewService creates new instance of BlogVieBlogViewServicewerService.
func NewBlogViewService(db *gorm.DB, repository repository.Repository) *BlogViewService {
	return &BlogViewService{
		DB:         db,
		Repository: repository,
	}
}

// AddBlogView adds new blog viewer to database.
func (service *BlogViewService) AddBlogView(blogView *blog.BlogView) error {

	// Validate tenant id.
	err := service.doesTenantExist(blogView.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(blogView.CreatedBy, blogView.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate blog ID.
	err = service.doesBlogExist(blogView.BlogID, blogView.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}


	// Validate viewer ID.
	err = service.doesCredentialExist(blogView.ViewerID, blogView.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if blog viewer is already present in database for the given blog id and viewer id.
	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting blog viewer already present in database.
	tempBlogView := blog.BlogView{}

	// Get blog viewer from database.
	isBlogViewNotFound := false
	if err := service.Repository.GetRecordForTenant(uow, blogView.TenantID, &tempBlogView,
		repository.Filter("`blog_id`=? AND `viewer_id`=?", blogView.BlogID, blogView.ViewerID)); err != nil {
		if err == gorm.ErrRecordNotFound {
			isBlogViewNotFound = true
		} else {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Internal server error")
		}
	}

	// If blog viewer is present then return.
	if !isBlogViewNotFound{
		return nil
	}

	// If blog viewer is not present then add blog viewer.
	err = service.Repository.Add(uow, blogView)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateBlogView updates blog viewer in database.
func (service *BlogViewService) UpdateBlogView(blogView *blog.BlogView) error {

	// Validate tenant ID.
	err := service.doesTenantExist(blogView.TenantID)
	if err != nil {
		return err
	}

	// Validate blog viewer ID.
	err = service.doesBlogViewExist(blogView.ID, blogView.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(blogView.UpdatedBy, blogView.TenantID)
	if err != nil {
		return err
	}

	// Validate blog ID.
	err = service.doesBlogExist(blogView.BlogID, blogView.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate viewer ID.
	err = service.doesCredentialExist(blogView.ViewerID, blogView.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update blog viewer.
	err = service.Repository.Update(uow, blogView)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteBlogView delete blog viewer from database.
func (service *BlogViewService) DeleteBlogView(blogView *blog.BlogView) error {
	credentialID := blogView.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(blogView.TenantID)
	if err != nil {
		return err
	}

	// Validate blog viewer ID.
	err = service.doesBlogViewExist(blogView.ID, blogView.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, blogView.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update blog viewer for updating deleted_by and deleted_at fields of blog viewer.
	if err := service.Repository.UpdateWithMap(uow, blogView, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", blogView.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog viewer could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetBlogView returns one blog viewer by id.
func (service *BlogViewService) GetBlogView(blogView *blog.BlogView) error {

	// Validate tenant ID.
	err := service.doesTenantExist(blogView.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get blog viewer by id from database.
	err = service.Repository.GetForTenant(uow, blogView.TenantID, blogView.ID, blogView)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBlogViews returns all blog viewers.
func (service *BlogViewService) GetBlogViews(blogViews *[]blog.BlogView,
	form url.Values, tenantID uuid.UUID, limit, offset int, totalCount *int) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all blog viewers.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, blogViews, "`created_at`",
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

// addSearchQueries adds all search queries by comparing with the blog viewer data.
func (service *BlogViewService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
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

// doesTenantExist validates if tenant exists or not in database.
func (service *BlogViewService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBlogViewExist validates if blog viewer exists or not in database.
func (service *BlogViewService) doesBlogViewExist(blogViewID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, blog.BlogView{},
		repository.Filter("`id` = ?", blogViewID))
	if err := util.HandleError("Invalid blog view ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBlogExist validates if blog exists or not in database.
func (service *BlogViewService) doesBlogExist(blogID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, blog.Blog{},
		repository.Filter("`id` = ?", blogID))
	if err := util.HandleError("Invalid blog ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *BlogViewService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
